package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/fernoe1/appointment-telegram-bot/internal/domain"
	"github.com/fernoe1/appointment-telegram-bot/internal/repository"
	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/client"
	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/constant"
	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/handler"
	"github.com/fernoe1/appointment-telegram-bot/migrate"
	"github.com/fernoe1/appointment-telegram-bot/pkg"
	"github.com/mymmrac/telego"
)

type Server struct {
	handler *handler.Handler
	bot     *client.Client
	r       *repository.R

	remindedMu sync.Mutex
	reminded   map[string]struct{}
}

func MustNew() *Server {
	db, err := pkg.NewGormDB(os.Getenv("DB_URL"))
	if err != nil {
		log.Panicf("internal.server.MustNew->NewGormDB: %v", err)
	}

	if err = migrate.Run(db); err != nil {
		log.Panicf("internal.server.MustNew->migrate.Run: %v", err)
	}

	r := repository.New(db)
	tgClient := client.MustNew(os.Getenv("BOT_TOKEN"), telego.WithDefaultLogger(false, true))

	return &Server{
		handler:  handler.MustNew(tgClient, r),
		bot:      tgClient,
		r:        r,
		reminded: make(map[string]struct{}),
	}
}

func (s *Server) Start() error {
	go s.appointmentWorker()

	return s.handler.Start()
}

func (s *Server) Stop() error {

	return s.handler.Stop()
}

func (s *Server) appointmentWorker() {
	s.processAppointments()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.processAppointments()
	}
}

func (s *Server) processAppointments() {
	now := time.Now()

	if err := s.r.DeleteAppointmentsBefore(now); err != nil {
		log.Printf("internal.server.processAppointments->DeleteAppointmentsBefore: %v", err)
	}

	appts, err := s.r.AppointmentsFrom(now)
	if err != nil {
		log.Printf("internal.server.processAppointments->AppointmentsFrom: %v", err)

		return
	}

	for _, appt := range appts {
		apptTime, err := appointmentTime(appt)
		if err != nil {
			log.Printf("internal.server.processAppointments->appointmentTime: %v", err)

			continue
		}

		if now.After(apptTime) {
			if err := s.r.DeleteAppointment(appt.TID); err != nil {
				log.Printf("internal.server.processAppointments->DeleteAppointment: %v", err)
			}

			s.clearReminderFlags(appt.TID)

			continue
		}

		left := apptTime.Sub(now)

		if left <= 3*time.Hour && left >= 3*time.Hour-5*time.Minute {
			s.sendReminder(appt, 3)
		}

		if left <= 1*time.Hour && left >= 1*time.Hour-5*time.Minute {
			s.sendReminder(appt, 1)
		}
	}
}

func appointmentTime(appt domain.Appointment) (time.Time, error) {
	if t, err := time.ParseInLocation(
		"2006-01-02 15:04",
		fmt.Sprintf("%s %02d:00", appt.Date, appt.Hour),
		time.Local,
	); err == nil {
		return t, nil
	}

	if t, err := time.ParseInLocation("2006-01-02 15:04", appt.Date, time.Local); err == nil {

		return t, nil
	}

	if t, err := time.ParseInLocation(time.RFC3339, appt.Date, time.Local); err == nil {

		return t, nil
	}

	return time.Time{}, fmt.Errorf("invalid appointment date format: %q", appt.Date)
}

func (s *Server) sendReminder(appt domain.Appointment, hoursLeft int) {
	key := fmt.Sprintf("%d:%dh", appt.TID, hoursLeft)
	if s.wasReminded(key) {
		return
	}

	userText := fmt.Sprintf(
		"Напоминаем, до вашей записи осталось %d часов.",
		hoursLeft,
	)

	adminText := fmt.Sprintf(
		"Напоминаем, через %d часов у вас встреча с %s.\nДанные:\n Тэг: %s\nДата: %s",
		hoursLeft,
		appt.PhoneNumber,
		appt.Username,
		appt.Date,
	)

	_, err := s.bot.SendMessage(context.Background(), &telego.SendMessageParams{
		ChatID: telego.ChatID{ID: appt.CID},
		Text:   userText,
	})
	if err != nil {
		log.Printf("internal.server.sendReminder->SendMessage(user): %v", err)
	}

	_, err = s.bot.SendMessage(context.Background(), &telego.SendMessageParams{
		ChatID: telego.ChatID{ID: constant.AdminTID},
		Text:   adminText,
	})
	if err != nil {
		log.Printf("internal.server.sendReminder->SendMessage(admin): %v", err)
	}

	if err == nil {
		s.setReminded(key)
	}
}

func (s *Server) clearReminderFlags(tid int64) {
	s.remindedMu.Lock()
	defer s.remindedMu.Unlock()

	delete(s.reminded, fmt.Sprintf("%d:3h", tid))
	delete(s.reminded, fmt.Sprintf("%d:1h", tid))
}

func (s *Server) wasReminded(key string) bool {
	s.remindedMu.Lock()
	defer s.remindedMu.Unlock()

	_, ok := s.reminded[key]

	return ok
}

func (s *Server) setReminded(key string) {
	s.remindedMu.Lock()
	defer s.remindedMu.Unlock()

	s.reminded[key] = struct{}{}
}
