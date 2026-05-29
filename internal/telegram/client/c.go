package client

import (
	"context"
	"log"

	"github.com/mymmrac/telego"
)

type Client struct {
	*telego.Bot
}

func MustNew(botToken string, options ...telego.BotOption) *Client {
	if botToken == "" {
		log.Fatal("telegram.client.MustNew: no token")
	}

	bot, err := telego.NewBot(botToken, options...)
	if err != nil {
		log.Panicf("telegram.client.MustNew->NewBot: %v", err)
	}

	return &Client{
		Bot: bot,
	}
}

func (c *Client) MustGetLongPollingUpdates() <-chan telego.Update {
	updates, err := c.UpdatesViaLongPolling(context.Background(), nil)
	if err != nil {
		log.Panicf("telegram.client.MustGetLongPollingUpdates: %v", err)
	}

	return updates
}
