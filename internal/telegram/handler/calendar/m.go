package calendar

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fernoe1/appointment-telegram-bot/internal/domain"
	"github.com/fernoe1/appointment-telegram-bot/internal/repository"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	dbf "github.com/thevan4/telegram-calendar/day_button_former"
	"github.com/thevan4/telegram-calendar/generator"
	tckm "github.com/thevan4/telegram-calendar/manager"
)

type Manager struct {
	km tckm.KeyboardManager
	r  *repository.R
}

func New(r *repository.R) *Manager {
	m := tckm.NewManager(
		generator.ChangeYearsForwardForChoose(1),
		generator.ChangeHomeButtonForBeauty("🏠︎"),
		generator.NewButtonsTextWrapper(
			dbf.ChangePostfixForCurrentDay(""),
			dbf.ChangePostfixForNonSelectedDay("❌"),
			dbf.ChangeTimezone(time.Local),
		),
	)

	return &Manager{km: m, r: r}
}

func (m *Manager) CallbackHandler(ctx *th.Context, query telego.CallbackQuery) error {
	sess := m.r.Session(query.From.ID)
	if sess == nil {
		err := ctx.Bot().AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
			CallbackQueryID: query.ID,
			Text:            "Время вашей сессии истекло. Пожалуйста, начните заново, снова используя команду /start.",
			ShowAlert:       true,
		})

		if err != nil {

			return err
		}

		err = ctx.Bot().DeleteMessage(ctx, &telego.DeleteMessageParams{
			ChatID:    query.Message.GetChat().ChatID(),
			MessageID: query.Message.GetMessageID(),
		})

		if err != nil {

			return err
		}

		return err
	}

	tn := time.Now()
	fullDays, err := m.r.FullDays()
	if err != nil {

		return err
	}

	if tn.Hour() >= 18 {
		fullDays[normalize(tn)] = struct{}{}
	}

	m.km.ApplyNewOptions(
		generator.ApplyNewOptionsForButtonsTextWrapper(
			dbf.ChangeUnselectableDaysBeforeDate(tn.AddDate(0, 0, -1)),
			dbf.ChangeUnselectableDays(fullDays),
		),
	)

	response := m.km.GenerateCalendarKeyboard(query.Data, tn)

	// handle unselectable day
	if response.IsUnselectableDay {
		err = ctx.Bot().AnswerCallbackQuery(ctx,
			&telego.AnswerCallbackQueryParams{
				CallbackQueryID: query.ID,
				Text: fmt.Sprintf(
					"На дату %s запись недоступна",
					response.SelectedDay.Format("02.01.2006"),
				),
				ShowAlert: true,
			},
		)

		return err
	}

	// handle selection
	if !response.SelectedDay.IsZero() {
		selectedDay := response.SelectedDay.Format(domain.AppointmentDateLayout)

		apptCount, err := m.r.AppointmentCountByDay(response.SelectedDay)
		if err != nil {

			return err
		}

		if apptCount >= domain.MaxAppointmentsPerDay {
			err = ctx.Bot().AnswerCallbackQuery(ctx,
				&telego.AnswerCallbackQueryParams{
					CallbackQueryID: query.ID,
					Text:            "Выбранная дата стала недоступной. Пожалуйста, выберите другую дату.",
					ShowAlert:       true,
				},
			)

			return err
		}

		timeButtons, err := m.createTimeButtons(response.SelectedDay)

		_, err = ctx.Bot().EditMessageText(ctx,
			&telego.EditMessageTextParams{
				ChatID:    telego.ChatID{ID: query.Message.GetChat().ID},
				MessageID: query.Message.Message().MessageID,
				Text: fmt.Sprintf("Выбрана дата: %s, теперь, выберите подходящее время.",
					selectedDay),
				ReplyMarkup: timeButtons,
			},
		)
		if err != nil {

			return err
		}

		sess.Day = response.SelectedDay
		m.r.SetSession(query.Message.GetChat().ID, sess)

		return nil
	}

	// handle void buttons
	if len(response.InlineKeyboardMarkup.InlineKeyboard) == 0 {
		err = ctx.Bot().AnswerCallbackQuery(ctx,
			&telego.AnswerCallbackQueryParams{
				CallbackQueryID: query.ID,
			},
		)

		return err
	}

	// handle movement

	b, err := json.Marshal(response.InlineKeyboardMarkup)
	if err != nil {

		return err
	}

	replyKeyboard := new(telego.InlineKeyboardMarkup)
	err = json.Unmarshal(b, replyKeyboard)
	if err != nil {

		return err
	}

	_, err = ctx.Bot().EditMessageReplyMarkup(ctx,
		&telego.EditMessageReplyMarkupParams{
			ChatID:      telego.ChatID{ID: query.Message.GetChat().ID},
			MessageID:   query.Message.Message().MessageID,
			ReplyMarkup: replyKeyboard,
		},
	)
	if err != nil {

		return err
	}

	return nil
}

func normalize(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func (m *Manager) createTimeButtons(day time.Time) (*telego.InlineKeyboardMarkup, error) {
	buttons := make([][]telego.InlineKeyboardButton, 0, 9)
	tn := time.Now()

	slots, err := m.r.AvailableTimeSlots(day)
	if err != nil {

		return nil, err
	}

	if slots == nil {

		return nil, nil
	}

	for _, h := range slots {
		if sameDay(tn, day) {
			if h <= tn.Hour() {
				continue
			}
		}

		buttons = append(buttons, []telego.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%d:00 - %d:00", h, h+1),
				CallbackData: fmt.Sprintf("time/%s/%d", day.Format(domain.AppointmentDateLayout), h),
			},
		})
	}

	return &telego.InlineKeyboardMarkup{InlineKeyboard: buttons}, nil
}

func sameDay(x, y time.Time) bool {
	return x.Year() == y.Year() && x.Month() == y.Month() && x.Day() == y.Day()
}
