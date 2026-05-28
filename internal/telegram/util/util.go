package util

import (
	"fmt"

	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/constant"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func GetChatID(messageOrCallbackQueryOrUpdate interface{}) int64 {
	switch morqu := messageOrCallbackQueryOrUpdate.(type) {
	case telego.Message:
		return morqu.Chat.ID
	case telego.CallbackQuery:
		return morqu.Message.Message().Chat.ID
	case telego.Update:
		return morqu.Message.Chat.ID
	default:
		return 0
	}
}

func GiveTimes(hour int, day string) *telego.InlineKeyboardMarkup {
	buttons := make([][]telego.InlineKeyboardButton, 0, 9)

	if hour < 8 {
		hour = 8
	}

	for hour < 18 {
		if hour == 13 {
			hour++
			continue
		}

		buttons = append(buttons, []telego.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%d:00 - %d:00", hour, hour+1),
				CallbackData: fmt.Sprintf("time/%s/%d", day, hour),
			},
		})

		hour++
	}

	return &telego.InlineKeyboardMarkup{InlineKeyboard: buttons}
}

func GiveMainMenu(text string) *telego.InlineKeyboardMarkup {
	if text == "" {
		return tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(constant.CallbackCalendarName).WithCallbackData(constant.CallbackCalendar),
			),
		)
	}

	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(text).WithCallbackData(constant.CallbackCalendar),
		),
	)
}
