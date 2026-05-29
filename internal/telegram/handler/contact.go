package handler

import (
	"github.com/fernoe1/appointment-telegram-bot/internal/repository"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func onContact(r *repository.R) th.Handler {
	return func(ctx *th.Context, upd telego.Update) error {
		sess := r.Session(upd.Message.Chat.ID)
		if sess == nil {

			return nil
		}

		if sess.Command == repository.Start {
			_, err := ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
				ChatID:      upd.Message.Chat.ChatID(),
				Text:        "Запись была успешно оформлена.",
				ReplyMarkup: &telego.ReplyKeyboardRemove{RemoveKeyboard: true},
			})

			return err
		}

		return nil
	}
}
