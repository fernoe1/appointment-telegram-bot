package handler

import (
	"context"
	"fmt"

	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/util"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func onStart(ctx *th.Context, msg telego.Message) error {
	_, err := ctx.Bot().SendMessage(context.Background(),
		tu.Message(
			tu.ID(util.GetChatID(msg)),
			fmt.Sprintf("Hello %s",
				msg.From.Username,
			),
		).WithReplyMarkup(util.GiveMainMenu("")))
	if err != nil {

		return fmt.Errorf("telegram.handler.onStart: %w", err)
	}

	return nil
}
