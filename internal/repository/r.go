package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/fernoe1/appointment-telegram-bot/internal/domain"
	"github.com/jellydator/ttlcache/v3"
	"gorm.io/gorm"
)

type Command string

const (
	Start  Command = "start"
	Edit   Command = "edit"
	See    Command = "see"
	Delete Command = "delete"
)

type Session struct {
	Command Command
	// change these w metadata in future
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

func (r *R) DeleteSession(CID int64) {
	r.s.Delete(CID)
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

func (r *R) AppointmentByTID(tid int64) (*domain.Appointment, error) {
	var appt domain.Appointment

	err := r.d.Where("t_id = ?", tid).First(&appt).Error

	if err == gorm.ErrRecordNotFound {

		return nil, nil
	}

	if err != nil {

		return nil, err
	}

	return &appt, err
}

func (r *R) FullDays() (map[time.Time]struct{}, error) {
	var appts []domain.Appointment

	err := r.d.Find(&appts).Error

	if err == gorm.ErrRecordNotFound {

		return nil, nil
	}

	if err != nil {

		return nil, err
	}

	dateCount := make(map[string]int)
	for _, appt := range appts {
		dateCount[appt.Date]++
	}
	for _, item := range r.s.Items() {
		dateCount[item.Value().Day.Format(domain.AppointmentDateLayout)]++
	}

	fullDays := make(map[time.Time]struct{})
	for dateStr, count := range dateCount {
		if count >= domain.MaxAppointmentsPerDay {
			date, err := time.Parse(domain.AppointmentDateLayout, dateStr)
			if err != nil {

				return nil, err
			}

			fullDays[date.Local().Add(-5*time.Hour)] = struct{}{}
		}
	}

	return fullDays, nil
}

func (r *R) AvailableTimeSlots(day time.Time) ([]int, error) {
	date := day.Format(domain.AppointmentDateLayout)

	var bookedHours []int
	err := r.d.Model(&domain.Appointment{}).Where("date = ?", date).Pluck("hour", &bookedHours).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {

		return nil, nil
	}

	if err != nil {

		return nil, err
	}

	booked := make(map[int]struct{}, len(bookedHours))
	for _, hour := range bookedHours {
		booked[hour] = struct{}{}
	}

	var slots []int
	for hour := 8; hour < 18; hour++ {
		if hour == 13 {
			continue
		}

		if _, exists := booked[hour]; !exists {
			slots = append(slots, hour)
		}
	}

	return slots, nil
}

func (r *R) TimeSlotExists(day time.Time, hour int) (bool, error) {
	var c int64

	date := day.Format(domain.AppointmentDateLayout)

	err := r.d.Where("date = ? AND hour = ?", date, hour).Find(&domain.Appointment{}).Count(&c).Error

	if err == gorm.ErrRecordNotFound {

		return false, nil
	}

	if err != nil {

		return false, err
	}

	for _, item := range r.s.Items() {
		if item.Value().Day.Equal(day) && item.Value().Hour == hour {
			c++
		}
	}

	return c > 0, nil
}

func (r *R) AppointmentCountByDay(day time.Time) (int64, error) {
	var c int64

	date := day.Format(domain.AppointmentDateLayout)

	err := r.d.Where("date = ?", date).Find(&domain.Appointment{}).Count(&c).Error

	if err == gorm.ErrRecordNotFound {

		return 0, nil
	}

	if err != nil {

		return 0, err
	}

	fmt.Println(day)

	for _, item := range r.s.Items() {
		fmt.Println(item.Value().Day)

		if item.Value().Day.Equal(day) {
			c++
		}
	}

	return c, nil
}

func (r *R) CreateAppointment(tid, cid int64, username, phoneNumber string, day time.Time, hour int) error {
	return r.d.Create(
		&domain.Appointment{
			TID:         tid,
			CID:         cid,
			Username:    username,
			PhoneNumber: phoneNumber,
			Date:        day.Format(domain.AppointmentDateLayout),
			Hour:        hour,
		},
	).Error
}

func (r *R) DeleteAppointment(tid int64) error {
	return r.d.Delete(&domain.Appointment{TID: tid}).Error
}

func (r *R) UpdateAppointment(tid, cid int64, username, phoneNumber string, day time.Time, hour int) error {
	return r.d.Updates(
		&domain.Appointment{
			TID:         tid,
			CID:         cid,
			Username:    username,
			PhoneNumber: phoneNumber,
			Date:        day.Format(domain.AppointmentDateLayout),
			Hour:        hour,
		},
	).Error
}

func (r *R) AppointmentsOn(day time.Time) ([]domain.Appointment, error) {
	var appts []domain.Appointment

	date := day.Format(domain.AppointmentDateLayout)

	err := r.d.Where("date = ?", date).Order("hour ASC").Find(&appts).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {

		return appts, nil
	}

	if err != nil {

		return nil, err
	}

	return appts, nil
}

func (r *R) AppointmentsFromToWeek(day time.Time) ([]domain.Appointment, error) {
	var appts []domain.Appointment

	weekday := int(day.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	start := day.AddDate(0, 0, -(weekday - 1))
	end := start.AddDate(0, 0, 6)

	err := r.d.
		Where(
			"date >= ? AND date <= ?",
			start.Format(domain.AppointmentDateLayout),
			end.Format(domain.AppointmentDateLayout),
		).Order("date ASC, hour ASC").Find(&appts).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {

		return nil, nil
	}

	if err != nil {

		return nil, err
	}

	return appts, nil
}

func (r *R) AllAppointments() ([]domain.Appointment, error) {
	var appts []domain.Appointment

	err := r.d.Order("date ASC, hour ASC").Find(&appts).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {

		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return appts, nil
}

func (r *R) AppointmentsFrom(day time.Time) ([]domain.Appointment, error) {
	var appts []domain.Appointment

	date := day.Format(domain.AppointmentDateLayout)

	err := r.d.
		Where("date >= ?", date).
		Order("date ASC, hour ASC").
		Find(&appts).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {

		return nil, nil
	}

	if err != nil {

		return nil, err
	}

	return appts, nil
}

func (r *R) DeleteAppointmentsBefore(day time.Time) error {
	date := day.Format(domain.AppointmentDateLayout)

	return r.d.Where("date < ?", date).Delete(&domain.Appointment{}).Error
}
