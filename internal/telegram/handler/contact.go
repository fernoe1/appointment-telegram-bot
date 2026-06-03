package handler

import (
	"fmt"

	"github.com/fernoe1/appointment-telegram-bot/internal/domain"
	"github.com/fernoe1/appointment-telegram-bot/internal/repository"
	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/constant"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func onContact(r *repository.R) th.Handler {
	return func(ctx *th.Context, upd telego.Update) error {
		var (
			tid         = upd.Message.From.ID
			cid         = upd.Message.Chat.ID
			username    = upd.Message.From.Username
			phoneNumber = upd.Message.Contact.PhoneNumber
		)

		sess := r.Session(upd.Message.Chat.ID)
		if sess == nil {

			return nil
		}

		if sess.Command == repository.Start {
			err := r.CreateAppointment(
				tid,
				cid,
				username,
				phoneNumber,
				sess.Day,
				sess.Hour,
			)

			if err != nil {
				_, _ = ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
					ChatID:      upd.Message.Chat.ChatID(),
					Text:        "Произошла ошибка. Пожалуйста, выполните команду /start ещё раз, чтобы начать заново.",
					ReplyMarkup: &telego.ReplyKeyboardRemove{RemoveKeyboard: true},
				})

				r.DeleteSession(cid)
				return err
			}

			_, err = ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
				ChatID:      upd.Message.Chat.ChatID(),
				Text:        "Запись была успешно оформлена.",
				ReplyMarkup: &telego.ReplyKeyboardRemove{RemoveKeyboard: true},
			})

			_, _ = ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
				ChatID: telego.ChatID{ID: constant.AdminTID},
				Text: fmt.Sprintf("Была создана новая запись.\nДанные:\nТэг: %s\nКонтакт: %s\nДата: %s\nВремя: %d",
					username,
					phoneNumber,
					sess.Day.Format(domain.AppointmentDateLayout),
					sess.Hour,
				),
			})

			if err != nil {

				r.DeleteSession(cid)
				return err
			}
		}

		if sess.Command == repository.Edit {
			err := r.UpdateAppointment(
				tid,
				cid,
				username,
				phoneNumber,
				sess.Day,
				sess.Hour,
			)

			if err != nil {
				_, _ = ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
					ChatID:      upd.Message.Chat.ChatID(),
					Text:        "Произошла ошибка. Пожалуйста, выполните команду /edit ещё раз, чтобы начать заново.",
					ReplyMarkup: &telego.ReplyKeyboardRemove{RemoveKeyboard: true},
				})

				r.DeleteSession(cid)
				return err
			}

			_, err = ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
				ChatID:      upd.Message.Chat.ChatID(),
				Text:        "Запись была успешно изменена.",
				ReplyMarkup: &telego.ReplyKeyboardRemove{RemoveKeyboard: true},
			})

			if err != nil {

				r.DeleteSession(cid)
				return err
			}
		}

		r.DeleteSession(cid)

		return nil
	}
}
