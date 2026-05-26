package bot

import (
	tele "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/layout"
	"gopkg.in/telebot.v4/middleware"
)

type Bot struct {
	*tele.Bot
	*layout.Layout
}

func New(path string) (*Bot, error) {
	lt, err := layout.New(path)
	if err != nil {
		return nil, err
	}

	b, err := tele.NewBot(lt.Settings())
	if err != nil {
		return nil, err
	}

	if cmds := lt.Commands(); cmds != nil {
		if err := b.SetCommands(cmds); err != nil {
			return nil, err
		}
	}

	return &Bot{b, lt}, nil
}

func (b *Bot) Start() {
	b.Use(middleware.Logger())
	b.Use(middleware.AutoRespond())
	b.Use(b.Middleware("en", func(r tele.Recipient) string {
		return "kz"
	}))

	b.Handle("/start", b.onStart)

	b.Bot.Start()
}
