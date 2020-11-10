package main

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Configure struct {
	Telegram struct {
		Token  string `yaml:"token"`
		ChatID int64  `yaml:"chatID"`
		Enable bool   `yaml:"enable"`
		Debug  bool   `yaml:"debug"`
	} `yaml:"telegram"`
	PagerDuty struct {
		RoutingKey             string   `yaml:"routingKey"`
		ApiEndpoint            string   `yaml:"apiEndpoint"`
		Enable                 bool     `yaml:"enable"`
		Debug                  bool     `yaml:"debug"`
		ClientUrl              string   `yaml:"clientUrl"`
		ImagesSrc              string   `yaml:"imagesSrc"`
		ImagesHref             string   `yaml:"imagesHref"`
		ImagesAlt              string   `yaml:"imagesAlt"`
		VoiceCallAlertWildcard []string `yaml:"voiceCallAlertWildcard"`
	} `yaml:"pagerDuty"`
}

func (c *Configure) ParseConfig(data []byte) error {
	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}
	if c.Telegram.Token == "" || c.Telegram.ChatID == 0 {
		return errors.New("telegram Token or ChatID is empty: must have some values")
	}
	return nil
}

func ProcessConfigParameters(configFile string) (cfg Configure, err error) {
	// read and decode config
	execPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadFile(execPath + "/" + configFile)
	if err != nil {
		log.Fatal(err)
	}

	var config Configure
	if err := config.ParseConfig(data); err != nil {
		log.Fatal("error in ParseConfig ", err)
	}
	return config, err
}
