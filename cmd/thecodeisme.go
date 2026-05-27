package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/client"
	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/handler"
	"github.com/mymmrac/telego"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	h := handler.MustNew(client.MustNew(os.Getenv("BOT_TOKEN"), telego.WithDefaultDebugLogger()))

	go func() { _ = h.Start() }()
	<-stop
	_ = h.Stop()
}
