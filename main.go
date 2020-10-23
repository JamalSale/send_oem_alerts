package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
)

const (
	confFile  = "config.yml"
	parseMode = "html"
)

type config struct {
	Telegram struct {
		Token  string `yaml:"token"`
		ChatID int64  `yaml:"chatID"`
		Debug  bool   `yaml:"debug"`
	} `yaml:"telegram"`
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

	var cfg config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal("Unable to decode config parameters ", err)
	}

	result, err := FormatMessage()
	if err != nil {
		log.Panic(err)
	}

	// init TgBot
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Panic(err)
	}

	// debug mode
	if cfg.Telegram.Debug {
		bot.Debug = true
		log.Print(result)
		defer log.Printf("Authorized on account: %s", bot.Self.UserName)
		defer log.Printf("Message send from: %s ", bot.Self.FirstName)

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

}
