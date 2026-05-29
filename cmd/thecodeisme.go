package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fernoe1/appointment-telegram-bot/internal/server"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	srv := server.MustNew()

	go func() { _ = srv.Start() }()
	<-stop
	_ = srv.Stop()
}
