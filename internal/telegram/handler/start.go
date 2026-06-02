package handler

import (
	"github.com/fernoe1/appointment-telegram-bot/internal/repository"
	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/constant"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func onStart(r *repository.R) th.MessageHandler {
	return func(ctx *th.Context, message telego.Message) error {
		var (
			tid = message.From.ID
			cid = message.Chat.ID
		)

		sess := r.Session(cid)
		if sess != nil {
			_, err := ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
				ChatID: telego.ChatID{ID: cid},
				Text:   "У вас уже запущена сессия записи на приём. Пожалуйста, сначала завершите её или отмените её с помощью команды /cancel.",
			},
			)

			return err
		}

		exists, err := r.AppointmentByTID(tid)
		if err != nil {

			return err
		}

		if exists != nil {
			_, err := ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
				ChatID: telego.ChatID{ID: cid},
				Text:   "У вас уже есть запись, введите /edit, чтобы изменить вашу запись.",
			},
			)

			return err
		}

		r.SetSession(cid, &repository.Session{Command: repository.Start})

		_, err = ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
			ChatID:      telego.ChatID{ID: cid},
			Text:        "Добро пожаловать! Пожалуйста, воспользуйтесь календарём, чтобы выбрать удобную дату для записи.",
			ReplyMarkup: createCalendarButton(),
		},
		)

		return err
	}
}

func createCalendarButton() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(constant.CalendarInlineButton).WithCallbackData(constant.CalendarInlineButtonCallback),
		),
	)
}
