package repository

import (
	"time"

	"github.com/fernoe1/appointment-telegram-bot/internal/domain"
	"github.com/jellydator/ttlcache/v3"
	"gorm.io/gorm"
)

type Command string

const (
	Start Command = "start"
	Edit  Command = "edit"
	Admin Command = "admin"
)

type Session struct {
	Command Command
	// change theses w metadata in future
	Day  time.Time
	Hour int
}

type R struct {
	d *gorm.DB
	s *ttlcache.Cache[int64, Session]
}

func New(d *gorm.DB) *R {
	s := ttlcache.New(
		ttlcache.WithTTL[int64, Session](10 * time.Minute),
	)

	go s.Start()

	return &R{
		d: d,
		s: s,
	}
}

func (r *R) Session(CID int64) *Session {
	if item := r.s.Get(CID); item != nil {
		val := item.Value()

		return &val
	}

	return nil
}

func (r *R) SetSession(CID int64, session *Session) {
	r.s.Set(CID, *session, ttlcache.DefaultTTL)
}

func (r *R) UserByTID(TID int64) (*domain.User, error) {
	var u domain.User

	err := r.d.Where("telegram_id = ?", TID).First(&u).Error

	if err != nil {

		return nil, err
	}

	return &u, err
}

func (r *R) FullDays() (map[time.Time]struct{}, error) {
	var bookedDays []domain.BookedDay

	err := r.d.Where("full = ?", true).Find(&bookedDays).Error
	if err != nil {

		return nil, err
	}

	fullDays := make(map[time.Time]struct{})

	for _, day := range bookedDays {
		t, err := time.Parse(domain.AppointmentDateLayout, day.Date)
		if err != nil {

			continue
		}

		fullDays[t] = struct{}{}
	}

	return fullDays, nil
}

func (r *R) TimeSlotExists(day time.Time, hour int) (bool, error) {
	var count int64

	date := day.Format(domain.AppointmentDateLayout)

	err := r.d.Model(&domain.Appointment{}).
		Joins("JOIN booked_days ON booked_days.id = appointments.booked_day_id").
		Where("booked_days.date = ? AND appointments.hour = ?", date, hour).
		Count(&count).
		Error
	if err != nil {

		return false, err
	}

	return count > 0, nil
}

func (r *R) AppointmentCountByDay(day time.Time) (int64, error) {
	var count int64

	date := day.Format(domain.AppointmentDateLayout)

	err := r.d.
		Model(&domain.Appointment{}).
		Joins("JOIN booked_days ON booked_days.id = appointments.booked_day_id").
		Where("booked_days.date = ?", date).
		Count(&count).
		Error
	if err != nil {

		return 0, err
	}

	return count, nil
}
