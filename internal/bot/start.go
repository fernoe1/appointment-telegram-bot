package bot

import (
	tele "gopkg.in/telebot.v4"
)

func (b Bot) onStart(c tele.Context) error {
	return c.Send(
		b.Text(c, "start", c.Sender()),
	)
}
