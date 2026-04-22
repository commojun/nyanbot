package anniversary

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/commojun/nyanbot/app/time_util"
	"github.com/commojun/nyanbot/masterdata"
	"github.com/commojun/nyanbot/masterdata/table"
)

type LineBot interface {
	TextMessageWithRoomKey(ctx context.Context, msg string, roomKey string) error
}

type AnniversaryManager struct {
	Anniversaries []table.Anniversary
	Bot           LineBot
}

var (
	TEMPLATE = "今日は%sから%d%sだよ！"
)

func New(bot LineBot) *AnniversaryManager {
	annivs := masterdata.GetTables().Anniversaries

	am := AnniversaryManager{
		Anniversaries: annivs,
		Bot:           bot,
	}
	return &am
}

func (am *AnniversaryManager) Run(ctx context.Context) error {
	now := time_util.LocalTime()
	for _, a := range am.Anniversaries {
		if err := ctx.Err(); err != nil {
			return err
		}
		msg, check, err := MakeCheckMessage(&a, now)
		if err != nil {
			return err
		}
		if check {
			err := am.Bot.TextMessageWithRoomKey(ctx, msg, a.RoomKey)
			if err != nil {
				return err
			}
		} else {
			log.Printf("did not send: %s", msg)
		}
	}
	return nil
}

func RandomMsg() (string, error) {
	annivs := masterdata.GetTables().Anniversaries
	if len(annivs) == 0 {
		return "", fmt.Errorf("no anniversaries available")
	}

	rand.Seed(time.Now().UnixNano())
	a := annivs[rand.Intn(len(annivs))]

	msg, _, err := MakeCheckMessage(&a, time_util.LocalTime())
	if err != nil {
		return "", err
	}

	return msg, nil
}

func MakeCheckMessage(a *table.Anniversary, now time.Time) (string, bool, error) {
	aTimeString := fmt.Sprintf("%s 00:00:00", a.Date)
	aTime, err := time_util.JSTParse(aTimeString)
	if err != nil {
		return "", false, err
	}
	period, err := strconv.Atoi(a.Period)
	if err != nil {
		return "", false, err
	}
	duration := now.Sub(aTime)
	days := int(duration.Hours()) / 24

	msg := ""
	check := true
	if aTime.Month() == now.Month() && aTime.Day() == now.Day() {
		years := now.Year() - aTime.Year()
		msg = fmt.Sprintf(TEMPLATE, a.Name, years, "年")
	} else {
		msg = fmt.Sprintf(TEMPLATE, a.Name, days, "日")
		if days%period != 0 {
			check = false
		}
	}

	return msg, check, nil
}
