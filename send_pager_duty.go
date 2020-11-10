package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func SendPagerDutyMessage(url string) {

	oemData, err := processOemEnvVariables()
	if err != nil {
		log.Panic(err)
	}

	pagerDutyMessage, err := FormatPagerDutyMessage(oemData)
	if err != nil {
		log.Panic(err)
	}

	sendData, err := json.Marshal(pagerDutyMessage)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(sendData))
	if err != nil {
		print(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	log.Println(string(body))
}
