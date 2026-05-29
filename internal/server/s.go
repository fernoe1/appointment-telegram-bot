package server

import (
	"fmt"
	"log"

	"github.com/fernoe1/appointment-telegram-bot/internal/repository"
	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/client"
	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/handler"
	"github.com/fernoe1/appointment-telegram-bot/migrate"
	"github.com/fernoe1/appointment-telegram-bot/pkg"
	"github.com/mymmrac/telego"
)

type Server struct {
	handler *handler.Handler
}

func MustNew() *Server {
	db, err := pkg.NewGormDB(dbURL)
	if err != nil {
		log.Panicf("internal.server.MustNew->NewGormDB: %v", err)
	}

	if err = migrate.Run(db); err != nil {
		log.Panicf("internal.server.MustNew->migrate.Run: %v", err)
	}

	r := repository.New(db)
	tgClient := client.MustNew(botToken, telego.WithDefaultDebugLogger())

	return &Server{
		handler: handler.MustNew(tgClient, r),
	}
}

func (s *Server) Start() error {
	if s == nil || s.handler == nil {
		return fmt.Errorf("internal.server.Start: server is nil")
	}
	return s.handler.Start()
}

func (s *Server) Stop() error {
	if s == nil || s.handler == nil {
		return nil
	}
	return s.handler.Stop()
}
