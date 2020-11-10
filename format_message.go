package main

import (
	"encoding/json"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
	"sort"
	s "strings"
	"time"
)

const (
	// emoji formatting
	emojiCrossMark    = "\xE2\x9D\x8C"
	emojiWarningsSign = "\xE2\x9A\xA0\xEF\xB8\x8F"

	// PagerDuty Events API v2 const
	layoutISO   = "02.01.2006 15:04:05 MST"
	eventAction = "trigger"
	client      = "OEM Monitoring"
)

// Oracle OEM metric variables from:
// https://docs.oracle.com/cd/E73210_01/EMADM/GUID-B48F6A84-EE89-498D-94E0-5DE1E7A0CFBC.htm#EMADM9066

type oemEnv struct {
	Severity             string `envconfig:"SEVERITY"`
	HostName             string `envconfig:"HOST_NAME"`
	TargetType           string `envconfig:"TARGET_TYPE"`
	TargetName           string `envconfig:"TARGET_NAME"`
	Message              string `envconfig:"MESSAGE"`
	Metric               string `envconfig:"METRIC_COLUMN"`
	MetricValue          string `envconfig:"VALUE"`
	IncidentCreationTime string `envconfig:"INCIDENT_CREATION_TIME"`
	EventReportedTime    string `envconfig:"EVENT_REPORTED_TIME"`
	JobName              string `envconfig:"SOURCE_OBJ_NAME"`
	JoOwner              string `envconfig:"SOURCE_OBJ_OWNER"`
	JoType               string `envconfig:"SOURCE_OBJ_SUB_TYPE"`
	JobStatus            string `envconfig:"EXECUTION_STATUS"`
	JobError             string `envconfig:"JOB_ERROR"`
}

// PagerDuty Events API v2 structure:
// https://developer.pagerduty.com/docs/events-api-v2/trigger-events/

type pagerDuty struct {
	Payload struct {
		Summary       string    `json:"summary"`
		Timestamp     time.Time `json:"timestamp,omitempty"`
		Source        string    `json:"source"`
		Severity      string    `json:"severity"`
		Component     string    `json:"component,omitempty"`
		Group         string    `json:"group,omitempty"`
		Class         string    `json:"class,omitempty"`
		CustomDetails oemEnv    `json:"custom_details,omitempty"`
	} `json:"payload"`
	RoutingKey string `json:"routing_key"`

	Images *pagerDutyImages `json:"images,omitempty"`
	Links  *pagerDutyLinks  `json:"links,omitempty"`

	EventAction string `json:"event_action"`
	Client      string `json:"client,omitempty"`
	ClientUrl   string `json:"client_url,omitempty"`
}

type pagerDutyImages struct {
	Src  string `json:"src"`
	Href string `json:"href"`
	Alt  string `json:"alt"`
}

type pagerDutyLinks struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

func MessageIsCritical() (critical bool) {

	cfg, err := ProcessConfigParameters(confFile)
	if err != nil {
		log.Fatal("Unable to ProcessConfigParameters ", err)
	}

	for i := range cfg.PagerDuty.VoiceCallAlertWildcard {
		if s.Contains(os.Getenv("MESSAGE"), cfg.PagerDuty.VoiceCallAlertWildcard[i]) {
			return true
		}

	}
	return false
}

func FormatPagerDutyMessage(oemDataMap map[string]string) (pagerDutyMsg *pagerDuty, err error) {

	cfg, err := ProcessConfigParameters(confFile)
	if err != nil {
		log.Fatal("Unable to ProcessConfigParameters ", err)
	}

	//  get OEM message map
	oemDataMap, err = processOemEnvVariables()
	if err != nil {
		log.Fatal("Unable to processOemEnvVariables ", err)
	}
	// convert map to json
	oemDataJson, _ := json.Marshal(oemDataMap)

	// convert json to struct
	oemMsg := oemEnv{}
	if err := json.Unmarshal(oemDataJson, &oemMsg); err != nil {
		log.Fatal("Unable to unmarshal oemDataJson ", err)
	}

	pagerDutyMsg = &pagerDuty{}

	data, _ := json.Marshal(pagerDutyMsg)

	err = json.Unmarshal(data, &pagerDutyMsg)

	// oem fields to pagerDuty fields mapping
	for key := range oemDataMap {
		switch key {
		case "HostName":
			pagerDutyMsg.Payload.Source = oemMsg.HostName
		case "Severity":
			pagerDutyMsg.Payload.Severity = oemMsg.Severity
		case "EventReportedTime":
			timestamp, _ := time.Parse(layoutISO, oemMsg.EventReportedTime)
			pagerDutyMsg.Payload.Timestamp = timestamp
		case "Message":
			pagerDutyMsg.Payload.CustomDetails.Message = oemMsg.Message
		case "Metric":
			pagerDutyMsg.Payload.Class = oemMsg.Metric
		case "MetricValue":
			pagerDutyMsg.Payload.CustomDetails.MetricValue = oemMsg.MetricValue
		case "TargetType":
			pagerDutyMsg.Payload.Group = oemMsg.TargetType
		case "TargetName":
			pagerDutyMsg.Payload.Component = oemMsg.TargetName

		}

	}

	pagerDutyMsg.Payload.Summary = oemMsg.Severity + " alert on " + oemMsg.HostName
	pagerDutyMsg.RoutingKey = cfg.PagerDuty.RoutingKey
	pagerDutyMsg.EventAction = eventAction
	pagerDutyMsg.Client = client
	pagerDutyMsg.ClientUrl = cfg.PagerDuty.ClientUrl

	return pagerDutyMsg, err
}

func processOemEnvVariables() (result map[string]string, err error) {

	oemMsg := &oemEnv{}
	err = envconfig.Process("", oemMsg)
	if err != nil {
		log.Fatal(err)
	}
	message, err := json.Marshal(oemMsg)

	if err != nil {
		log.Fatal("Unable to marshal oemMsg: ", err)
	}

	var oemMsgMap map[string]string

	err = json.Unmarshal(message, &oemMsgMap)

	return oemMsgMap, err
}

func FormatTelegramMessage() (result string, err error) {
	// emoji formatting
	switch os.Getenv("SEVERITY") {
	case "Critical":
		os.Setenv("SEVERITY", emojiCrossMark+"Critical")
	case "Warning":
		os.Setenv("SEVERITY", emojiWarningsSign+"Warning")
	}

	// unmarshal OemEnv struct to alertMsg map
	var alertMsg map[string]string

	alertMsg, err = processOemEnvVariables()
	if err != nil {
		log.Fatal("Unable to processOemEnvVariables ", err)
	}

	// we need to sort keys for order iteration over map
	keys := make([]string, 0)
	for k := range alertMsg {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		// filter out all empty values
		if len(alertMsg[k]) > 0 {
			// Convert each key/value pair to string adding html bold tags
			result += fmt.Sprintf("<b>%s</b> : %s \n", k, alertMsg[k])

		}

	}

	return result, err
}
