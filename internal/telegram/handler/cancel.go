package handler

import (
	"github.com/fernoe1/appointment-telegram-bot/internal/repository"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func onCancel(r *repository.R) th.MessageHandler {
	return func(ctx *th.Context, message telego.Message) error {
		var (
			cid = message.Chat.ID
		)

		r.DeleteSession(cid)

		_, err := ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
			ChatID: telego.ChatID{ID: message.Chat.ID},
			Text:   "Процесс записи успешно отменён.",
		})

		return err
	}
}
