package alarm

import (
	"time"

	"github.com/commojun/nyanbot/app/linebot"
	"github.com/commojun/nyanbot/masterdata/table"
)

type AlarmManager struct {
	Alarms *[]table.Alarm
	Bot    *linebot.LineBot
}

func New() (*AlarmManager, error) {
	alms, err := table.Alarms()
	if err != nil {
		return &AlarmManager{}, err
	}

	bot, err := linebot.New()
	if err != nil {
		return &AlarmManager{}, err
	}

	am := AlarmManager{
		Alarms: alms,
		Bot:    bot,
	}

	return &am, nil
}

func (am *AlarmManager) Run() error {

	for _, a := range *am.Alarms {
		chk, err := Check(&a, time.Now())
		if err != nil {
			return err
		}
		if chk == false {
			continue
		}
		am.Bot.TextMessage(a.Message)
		if err != nil {
			return err
		}
	}

	return nil

}
