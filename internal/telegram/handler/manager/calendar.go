package manager

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/util"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	dbf "github.com/thevan4/telegram-calendar/day_button_former"
	"github.com/thevan4/telegram-calendar/generator"
	"github.com/thevan4/telegram-calendar/manager"
)

type CalendarManager struct {
	manager     manager.KeyboardManager
	timeManager *TimeManager
}

func NewCalendarManager(tm *TimeManager) *CalendarManager {
	tn := time.Now()
	prevDay := tn.AddDate(0, 0, -1)
	m := manager.NewManager(
		generator.ChangeYearsForwardForChoose(1),
		generator.ChangeHomeButtonForBeauty("🏠︎"),
		generator.NewButtonsTextWrapper(
			dbf.ChangeUnselectableDaysBeforeDate(prevDay),
			dbf.ChangePostfixForCurrentDay(""),
			dbf.ChangePostfixForNonSelectedDay("❌"),
		),
	)

	cm := &CalendarManager{manager: m, timeManager: tm}
	go cm.atNextMidnightChangeUnselectableDaysBeforeDate(tn)

	return cm
}

func (cm *CalendarManager) atNextMidnightChangeUnselectableDaysBeforeDate(currentDate time.Time) {
	nextMidnight := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day()+1, 0, 0, 0, 0,
		currentDate.Location())
	duration := nextMidnight.Sub(currentDate)

	time.AfterFunc(duration, func() {
		prevDay := time.Date(nextMidnight.Year(), nextMidnight.Month(), nextMidnight.Day()-1, 0, 0, 0, 0,
			nextMidnight.Location())
		cm.manager.ApplyNewOptions(generator.ApplyNewOptionsForButtonsTextWrapper(
			dbf.ChangeUnselectableDaysBeforeDate(prevDay),
		))

		cm.atNextMidnightChangeUnselectableDaysBeforeDate(nextMidnight)
	})
}

func (cm *CalendarManager) CallbackQueryForCalendar(ctx *th.Context, query telego.CallbackQuery) error {
	now := time.Now()
	generateCalendarKeyboardResponse := cm.manager.GenerateCalendarKeyboard(query.Data, now)

	if generateCalendarKeyboardResponse.IsUnselectableDay {
		err := ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Day "+
			generateCalendarKeyboardResponse.SelectedDay.Format("02.01.2006")+
			" is unselectable").WithShowAlert())

		return err
	}

	if !generateCalendarKeyboardResponse.SelectedDay.IsZero() {
		selectedDay := generateCalendarKeyboardResponse.SelectedDay.Format("02.01.2006")

		msg := &telego.EditMessageTextParams{
			ChatID:    telego.ChatID{ID: util.GetChatID(query)},
			MessageID: query.Message.Message().MessageID,
			Text: fmt.Sprintf("Selected date %s. Please select time now.",
				selectedDay),
			ReplyMarkup: util.GiveTimes(1, selectedDay),
		}
		_, err := ctx.Bot().EditMessageText(ctx, msg)

		return err
	}

	if len(generateCalendarKeyboardResponse.InlineKeyboardMarkup.InlineKeyboard) == 0 {
		err := ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))

		return err
	}

	b, err := json.Marshal(generateCalendarKeyboardResponse.InlineKeyboardMarkup)
	if err != nil {

		return err
	}

	replyKeyboard := new(telego.InlineKeyboardMarkup)
	err = json.Unmarshal(b, replyKeyboard)
	if err != nil {

		return err
	}

	responseMsg := &telego.EditMessageReplyMarkupParams{
		ChatID:      telego.ChatID{ID: util.GetChatID(query)},
		MessageID:   query.Message.Message().MessageID,
		ReplyMarkup: replyKeyboard,
	}

	_, err = ctx.Bot().EditMessageReplyMarkup(ctx, responseMsg)
	if err != nil {

		return err
	}

	return nil
}
