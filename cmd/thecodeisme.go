package main

import (
	"log"

	"github.com/fernoe1/appointment-telegram-bot/internal/bot"
)

func main() {
	b, err := bot.New("./config/config.local.yaml")
	if err != nil {
		log.Fatal(err)
	}

	b.Start()
}
