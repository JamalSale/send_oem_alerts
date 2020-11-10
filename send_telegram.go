package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
)

const (
	parseMode = "html"
)

func SendTelegramMessage() {
	cfg, err := ProcessConfigParameters(confFile)
	if err != nil {
		log.Fatal("Unable to ProcessConfigParameters ", err)
	}

	// prepare data for Telegram
	result, err := FormatTelegramMessage()
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
		log.Print("FormatTelegramMessage: result ", result)
		defer log.Printf("Authorized on account: %s", bot.Self.UserName)
		defer log.Printf("Message send from: %s ", bot.Self.FirstName)
		// write env to log
		x, err := os.Create("debug.log")
		if err != nil {
			log.Println(err)
		}
		for _, env := range os.Environ() {
			_, err := x.WriteString(env + "\n")
			if err != nil {
				log.Println(err)
				err = x.Close()
				if err != nil {
					log.Fatal(err)
				}
				return
			}
		}

	}

	// if there is empty messages - dont try to send it
	if len(result) == 0 {
		log.Println("No messages to send")
		os.Exit(0)
	}

	// send message to Telegram
	msg := tgbotapi.NewMessage(cfg.Telegram.ChatID, result)
	msg.ParseMode = parseMode
	_, err = bot.Send(msg)

	if err != nil {
		log.Fatal(err)
	}

}
