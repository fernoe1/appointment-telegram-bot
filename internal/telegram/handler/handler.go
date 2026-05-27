package handler

import (
	"log"

	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/client"
	th "github.com/mymmrac/telego/telegohandler"
)

type Handler struct {
	*th.BotHandler
}

func MustNew(c *client.Client) *Handler {
	if c == nil {
		log.Fatal("telegram.handler.MustNewBotHandler: client is nil")
	}

	bh, err := th.NewBotHandler(c.Bot, c.MustGetLongPollingUpdates())
	if err != nil {
		log.Panicf("telegram.handler.MustNewBotHandler->NewBotHandler: %v", err)
	}

	h := &Handler{
		BotHandler: bh,
	}

	return h
}
