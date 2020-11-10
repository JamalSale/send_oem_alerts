package main

import (
	"log"
)

const (
	confFile = "config.yml"
)

func main() {

	cfg, err := ProcessConfigParameters(confFile)
	if err != nil {
		log.Fatal("Unable to ProcessConfigParameters ", err)
	}

	// send alert to PagerDuty
	if cfg.PagerDuty.Enable && MessageIsCritical() {
		SendPagerDutyMessage(cfg.PagerDuty.ApiEndpoint)

	}

	// send alert to PagerDuty
	if cfg.Telegram.Enable {
		SendTelegramMessage()
	}

}
