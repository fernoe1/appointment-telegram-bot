package manager

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fernoe1/appointment-telegram-bot/internal/telegram/util"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type TimeManager struct {
	sync.RWMutex
	currentHour int
}

func NewTimeManager() *TimeManager {
	tn := time.Now()
	hour := tn.Hour()

	tm := TimeManager{currentHour: hour}
	go tm.changeHour(tn)

	return &tm
}

func (tm *TimeManager) updateHour() {
	tm.Lock()
	defer tm.Unlock()

	tm.currentHour = time.Now().Hour()
}

func (tm *TimeManager) CurrentHour() int {
	tm.RLock()
	defer tm.RUnlock()

	return tm.currentHour
}

func (tm *TimeManager) changeHour(currentDate time.Time) {
	nextHour := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), currentDate.Hour()+1, 0, 0, 0,
		currentDate.Location())
	duration := nextHour.Sub(currentDate)

	time.AfterFunc(duration, func() {
		tm.updateHour()
		tm.changeHour(nextHour)
	})
}

func (tm *TimeManager) CallbackQueryForTime(ctx *th.Context, query telego.CallbackQuery) error {
	parts := strings.Split(query.Data, "/")
	fmt.Println(parts)

	_, err := ctx.Bot().EditMessageText(ctx, &telego.EditMessageTextParams{
		ChatID:    telego.ChatID{ID: util.GetChatID(query)},
		MessageID: query.Message.Message().MessageID,
		Text:      fmt.Sprintf("selected: %s %s", parts[1], parts[2]),
	})

	keyboard := tu.Keyboard(
		tu.KeyboardRow(
			tu.KeyboardButton("Contact").WithRequestContact(),
		),
	).WithResizeKeyboard().WithOneTimeKeyboard()

	msg := tu.Message(telego.ChatID{ID: util.GetChatID(query)}, "Please share your contact").WithReplyMarkup(keyboard)
	_, err = ctx.Bot().SendMessage(ctx, msg)

	return err
}
