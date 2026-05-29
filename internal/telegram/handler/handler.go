package handler

import (
	"context"
	"log"

	"github.com/fernoe1/appointment-telegram-bot/internal/repository"
	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/client"
	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/constant"
	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/handler/calendar"
	TIME "github.com/fernoe1/appointment-telegram-bot/internal/telegram/handler/time"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type Handler struct {
	*th.BotHandler
}

func MustNew(c *client.Client, r *repository.R) *Handler {
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

	cm := calendar.New(r)
	tm := TIME.New(r)

	h.RegisterHandlers(cm.CallbackHandler, tm.CallbackHandler, r)

	return h
}

func (h *Handler) RegisterHandlers(
	callbackCalendarHandler,
	callbackTimeHandler th.CallbackQueryHandler,
	r *repository.R,
) {
	h.HandleMessage(onStart(r), th.CommandEqual("start"))
	h.Handle(onContact(r), func(ctx context.Context, update telego.Update) bool {
		return update.Message != nil && update.Message.Contact != nil
	})

	h.HandleCallbackQuery(callbackCalendarHandler, th.AnyCallbackQueryWithMessage(),
		th.CallbackDataContains(constant.CalendarInlineButtonCallback))
	h.HandleCallbackQuery(callbackTimeHandler, th.AnyCallbackQueryWithMessage(),
		th.CallbackDataPrefix(constant.TimeInlineButtonCallback))

}
