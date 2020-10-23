package main

import (
	"encoding/json"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
	"sort"
)

const (
	crossMark    = "\xE2\x9D\x8C"
	warningsSign = "\xE2\x9A\xA0\xEF\xB8\x8F"
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

func FormatMessage() (result string, err error) {
	// emoji formatting
	switch os.Getenv("SEVERITY") {
	case "Critical":
		os.Setenv("SEVERITY", crossMark+"Critical")
	case "Warning":
		os.Setenv("SEVERITY", warningsSign+"Warning")
	}

	// unmarshal OemEnv struct to alertMsg map
	oemMsg := &oemEnv{}
	var alertMsg map[string]string
	err = envconfig.Process("", oemMsg)
	if err != nil {
		log.Fatal(err)
	}

	fields, _ := json.Marshal(oemMsg)

	err = json.Unmarshal(fields, &alertMsg)

	if err != nil {
		log.Fatal("Unable to unmarshal alertMsg", err)
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
