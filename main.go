package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

const (
	confFile  = "config.yml"
	parseMode = "html"
)

type Config struct {
	Telegram struct {
		Token  string `yaml:"token"`
		ChatID int64  `yaml:"chatID"`
		Debug  bool   `yaml:"debug"`
	} `yaml:"telegram"`
}

// Oracle OEM metric variables from:
// https://docs.oracle.com/cd/E73210_01/EMADM/GUID-B48F6A84-EE89-498D-94E0-5DE1E7A0CFBC.htm#EMADM9066

type OemEnv struct {
	Severity             string `envconfig:"SEVERITY"`
	HostName             string `envconfig:"HOST_NAME"`
	TargetType           string `envconfig:"TARGET_TYPE"`
	TargetName           string `envconfig:"TARGET_NAME"`
	Message              string `envconfig:"MESSAGE"`
	Metric               string `envconfig:"METRIC_COLUMN"`
	MetricValue          string `envconfig:"VALUE"`
	IncidentCreationTime string `envconfig:"INCIDENT_CREATION_TIME"`
	EventReportedTime    string `envconfig:"EVENT_REPORTED_TIME"`
}

func main() {
	// read and decode config
	execPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open(execPath + "/" + confFile)
	if err != nil {
		log.Fatal("Unable to read config:  ", err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal("Unable to unmarshal OemEnv struct ", err)
	}

	// unmarshal OemEnv struct to alertMsg map
	oemMsg := &OemEnv{}
	var alertMsg map[string]string
	err = envconfig.Process("", oemMsg)
	if err != nil {
		log.Fatal(err)
	}

	fields, _ := json.Marshal(oemMsg)

	err = json.Unmarshal(fields, &alertMsg)

	if err != nil {
		log.Fatal("Unable to unmarshal config ", err)
	}

	// we need to sort keys for order iteration over map
	var result string
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

	// init TgBot
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Panic(err)
	}

	// debug mode: print message fields
	if cfg.Telegram.Debug == true {
		bot.Debug = true
		log.Println(result)
		for field, val := range alertMsg {
			fmt.Println("alertMsg: ", field, val)
		}
	}

	// if there is empty messages - dont try to send it
	if len(result) == 0 {
		log.Println("No messages to send")
		os.Exit(0)
	}
	// send message
	msg := tgbotapi.NewMessage(cfg.Telegram.ChatID, result)
	msg.ParseMode = parseMode
	_, err = bot.Send(msg)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Message send from: %s ", bot.Self.FirstName)

	log.Printf("Authorized on account: %s", bot.Self.UserName)

}
