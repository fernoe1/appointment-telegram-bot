package handler

import (
	"fmt"

	"github.com/fernoe1/appointment-telegram-bot/internal/repository"
	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/constant"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func onSee(r *repository.R, adminTID int64) th.MessageHandler {
	return func(ctx *th.Context, message telego.Message) error {
		if message.From.ID != adminTID {
			appt, err := r.AppointmentByTID(message.From.ID)
			if err != nil {

				return err
			}

			if appt == nil {
				_, err = ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
					ChatID: telego.ChatID{ID: message.Chat.ID},
					Text:   "У вас нет записи. Пожалуйста, используйте команду /start, чтобы записаться.",
				},
				)

				return nil
			}

			_, err = ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
				ChatID: telego.ChatID{ID: message.Chat.ID},
				Text:   fmt.Sprintf("Ваша запись назначена на %s в %d:00-%d-00.", appt.Date, appt.Hour, appt.Hour+1),
			},
			)

			return err
		}

		var (
			cid = message.Chat.ID
		)

		_, err := ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
			ChatID:      telego.ChatID{ID: cid},
			Text:        "Выберите один из вариантов.",
			ReplyMarkup: createSeeButtons(),
		},
		)

		return err
	}
}

func createSeeButtons() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(constant.SeeInlineButtonToday).WithCallbackData(constant.SeeInlineButtonTodayCallback),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(constant.SeeInlineButtonWeek).WithCallbackData(constant.SeeInlineButtonWeekCallback),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(constant.SeeInlineButtonAll).WithCallbackData(constant.SeeInlineButtonAllCallback),
		),
	)
}
