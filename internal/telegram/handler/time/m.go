package time

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fernoe1/appointment-telegram-bot/internal/domain"
	"github.com/fernoe1/appointment-telegram-bot/internal/repository"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Manager struct {
	r *repository.R
}

func New(r *repository.R) *Manager {
	return &Manager{r}
}

func (m *Manager) CallbackHandler(ctx *th.Context, query telego.CallbackQuery) error {
	tn := time.Now()
	p := strings.Split(query.Data, "/")
	if len(p) < 3 {

		return nil
	}

	day, err := parseDay(p[1])
	if err != nil {

		return err
	}

	hour, err := strconv.Atoi(p[2])
	if err != nil {

		return err
	}

	if sameDay(day, tn) && tn.Hour() > hour {
		err = ctx.Bot().AnswerCallbackQuery(ctx,
			&telego.AnswerCallbackQueryParams{
				CallbackQueryID: query.ID,
				Text:            "Это время уже прошло.",
				ShowAlert:       true,
			},
		)

		return err
	}

	if hour < 8 || hour >= 18 || hour == 13 {
		err = ctx.Bot().AnswerCallbackQuery(ctx,
			&telego.AnswerCallbackQueryParams{
				CallbackQueryID: query.ID,
				Text:            "Неверно выбрано время.",
				ShowAlert:       true,
			},
		)

		return err
	}

	taken, err := m.r.TimeSlotExists(day, hour)
	if err != nil {

		return err
	}

	if taken {
		err = ctx.Bot().AnswerCallbackQuery(ctx,
			&telego.AnswerCallbackQueryParams{
				CallbackQueryID: query.ID,
				Text:            "Этот час уже забронирован.",
				ShowAlert:       true,
			},
		)

		return err
	}

	appointmentCount, err := m.r.AppointmentCountByDay(day)
	if err != nil {

		return err
	}

	if appointmentCount > domain.MaxAppointmentsPerDay {
		err = ctx.Bot().AnswerCallbackQuery(ctx,
			&telego.AnswerCallbackQueryParams{
				CallbackQueryID: query.ID,
				Text:            "Выбранный день недоступен, пожалуйста, начните заново с помощью /start.",
				ShowAlert:       true,
			},
		)

		return err
	}

	_, err = ctx.Bot().EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:    query.Message.GetChat().ChatID(),
		MessageID: query.Message.GetMessageID(),
		Text: fmt.Sprintf("Выбранный день: %s, выбранное время: %d:00 - %d:00.",
			day.Format("02.01.2006"),
			hour,
			hour+1,
		),
	})
	if err != nil {

		return err
	}

	sess := m.r.Session(query.Message.GetChat().ID)
	if sess != nil && sess.Command == repository.Start {
		sess.Day = day
		sess.Hour = hour
	}

	_, err = ctx.Bot().SendMessage(ctx, tu.Message(telego.ChatID{ID: query.Message.GetChat().ID},
		"Пожалуйста, поделитесь своим контактом, чтобы мы могли связаться с вами.").
		WithReplyMarkup(
			tu.Keyboard(
				tu.KeyboardRow(
					tu.KeyboardButton("Поделиться контактом").WithRequestContact(),
				),
			),
		),
	)

	return err
}

func parseDay(day string) (time.Time, error) {
	return time.Parse("02.01.2006", day)
}

func sameDay(x, y time.Time) bool {
	return x.Year() == y.Year() && x.Month() == y.Month() && x.Day() == y.Day()
}
