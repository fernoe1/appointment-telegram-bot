package handler

import (
	"errors"
	"strconv"

	"github.com/fernoe1/appointment-telegram-bot/internal/repository"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"gorm.io/gorm"
)

func onDelete(r *repository.R) th.MessageHandler {
	return func(ctx *th.Context, message telego.Message) error {
		_, _, args := tu.ParseCommand(message.Text)
		if len(args) == 0 {
			_, err := ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
				ChatID: message.Chat.ChatID(),
				Text:   "Формат команды: /delete <telegram_id>",
			})

			return err
		}

		tid, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			_, _ = ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
				ChatID: message.Chat.ChatID(),
				Text:   "Недопустимый Telegram ID.",
			})

			return err
		}

		err = r.DeleteAppointment(tid)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_, _ = ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
				ChatID: message.Chat.ChatID(),
				Text:   "Telegram ID не найден.",
			})

			return err
		}

		if err == nil {
			_, _ = ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
				ChatID: message.Chat.ChatID(),
				Text:   "Запись успешно удалена.",
			})

			return nil
		}

		return err
	}
}
