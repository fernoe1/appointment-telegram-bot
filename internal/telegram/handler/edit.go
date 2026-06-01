package handler

import (
	"github.com/fernoe1/appointment-telegram-bot/internal/repository"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func onEdit(r *repository.R) th.MessageHandler {
	return func(ctx *th.Context, message telego.Message) error {
		var (
			tid = message.From.ID
			cid = message.Chat.ID
		)

		sess := r.Session(cid)
		if sess != nil {
			_, err := ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
				ChatID: telego.ChatID{ID: cid},
				Text:   "У вас уже запущена сессия записи на приём. Пожалуйста, сначала завершите её.",
			},
			)

			return err
		}

		exists, err := r.AppointmentByTID(tid)
		if err != nil {

			return err
		}

		if exists == nil {
			_, err := ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
				ChatID: telego.ChatID{ID: cid},
				Text:   "У вас нет активной записи. Введите /start, чтобы записаться.",
			},
			)

			return err
		}

		r.SetSession(cid, &repository.Session{Command: repository.Edit})

		_, err = ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
			ChatID:      telego.ChatID{ID: cid},
			Text:        "Пожалуйста, воспользуйтесь календарём, чтобы изменить дату.",
			ReplyMarkup: createCalendarButton(),
		},
		)

		return err
	}
}
