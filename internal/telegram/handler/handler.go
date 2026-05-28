package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/client"
	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/constant"
	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/handler/manager"
	"github.com/mymmrac/telego"
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

	tm := manager.NewTimeManager()
	cm := manager.NewCalendarManager(tm)

	h.RegisterHandlers(cm.CallbackQueryForCalendar, tm.CallbackQueryForTime)

	return h
}

func (h *Handler) RegisterHandlers(callbackCalendarHandler, callbackTimeHandler th.CallbackQueryHandler) {
	h.HandleMessage(onStart, th.CommandEqual("start"))
	h.Handle(func(ctx *th.Context, upd telego.Update) error {
		contact := upd.Message.Contact
		fmt.Printf("Received contact: %s (%s)\n", contact.PhoneNumber, contact.UserID)
		return nil
	}, func(ctx context.Context, update telego.Update) bool {
		return update.Message != nil && update.Message.Contact != nil
	})

	h.HandleCallbackQuery(callbackCalendarHandler, th.AnyCallbackQueryWithMessage(),
		th.CallbackDataContains(constant.CallbackCalendar))
	h.HandleCallbackQuery(callbackTimeHandler, th.AnyCallbackQueryWithMessage(),
		th.CallbackDataPrefix("time"))
}
