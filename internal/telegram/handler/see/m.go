package see

import (
	"fmt"
	"strings"
	"time"

	"github.com/fernoe1/appointment-telegram-bot/internal/domain"
	"github.com/fernoe1/appointment-telegram-bot/internal/repository"
	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/constant"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type Manager struct {
	r *repository.R
}

func New(r *repository.R) *Manager {
	return &Manager{r}
}

func (m *Manager) CallbackHandler(ctx *th.Context, query telego.CallbackQuery) error {
	sess := m.r.Session(query.From.ID)
	if sess == nil {
		err := ctx.Bot().AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{
			CallbackQueryID: query.ID,
			Text:            "Время вашей сессии истекло. Пожалуйста, начните заново, снова используя команду /see.",
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

	if sess.Command != repository.See {
		err := ctx.Bot().DeleteMessage(ctx, &telego.DeleteMessageParams{
			ChatID:    query.Message.GetChat().ChatID(),
			MessageID: query.Message.GetMessageID(),
		})

		return err
	}

	if query.Data == constant.SeeInlineButtonTodayCallback {
		appts, err := m.r.AppointmentsOn(time.Now())
		if err != nil {

			return err
		}

		_, err = ctx.Bot().EditMessageText(ctx, &telego.EditMessageTextParams{
			ChatID:    query.Message.GetChat().ChatID(),
			MessageID: query.Message.GetMessageID(),
			Text:      formatAppointments(appts),
		})
	}

	if query.Data == constant.SeeInlineButtonWeekCallback {
		appts, err := m.r.AppointmentsFromToWeek(time.Now())
		if err != nil {

			return err
		}

		_, err = ctx.Bot().EditMessageText(ctx, &telego.EditMessageTextParams{
			ChatID:    query.Message.GetChat().ChatID(),
			MessageID: query.Message.GetMessageID(),
			Text:      formatAppointments(appts),
		})
	}

	if query.Data == constant.SeeInlineButtonAllCallback {
		appts, err := m.r.AllAppointments()
		if err != nil {

			return err
		}

		_, err = ctx.Bot().EditMessageText(ctx, &telego.EditMessageTextParams{
			ChatID:    query.Message.GetChat().ChatID(),
			MessageID: query.Message.GetMessageID(),
			Text:      formatAppointments(appts),
		})
	}

	return nil
}

func formatAppointments(appts []domain.Appointment) string {
	if len(appts) == 0 {
		return "Записей нет."
	}

	var b strings.Builder

	b.WriteString("Записи:\n\n")

	for _, a := range appts {
		fmt.Fprintf(
			&b,
			"Телеграм ID: %d\nДень: %s\nЧас: %d:00-%d:00\nТэг: %s\nНомер телефона:%s\n\n",
			a.TID,
			a.Date,
			a.Hour,
			a.Hour+1,
			a.Username,
			a.PhoneNumber,
		)
	}

	return b.String()
}
