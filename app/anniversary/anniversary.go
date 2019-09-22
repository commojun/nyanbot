package anniversary

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/commojun/nyanbot/app/linebot"
	"github.com/commojun/nyanbot/app/time_util"
	"github.com/commojun/nyanbot/masterdata/table"
)

type AnniversaryManager struct {
	Anniversaries []table.Anniversary
	Bot           *linebot.LineBot
}

var (
	TEMPLATE = "今日で%sから%d%sだよ！"
)

func New() (*AnniversaryManager, error) {
	annivs, err := table.Anniversaries()
	if err != nil {
		return &AnniversaryManager{}, err
	}

	bot, err := linebot.New()
	if err != nil {
		return &AnniversaryManager{}, err
	}

	am := AnniversaryManager{
		Anniversaries: annivs,
		Bot:           bot,
	}
	return &am, nil
}

func (am *AnniversaryManager) Run() error {
	now := time_util.LocalTime()
	for _, a := range am.Anniversaries {
		msg, check, err := MakeCheckMessage(&a, now)
		if err != nil {
			return err
		}
		if check {
			err := am.Bot.TextMessageWithRoomKey(msg, a.RoomKey)
			if err != nil {
				return err
			}
		} else {
			log.Printf("msg: %s", msg)
		}
	}
	return nil
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
