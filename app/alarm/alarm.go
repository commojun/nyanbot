package alarm

import (
	"log"
	"time"

	"github.com/commojun/nyanbot/app/linebot"
	"github.com/commojun/nyanbot/constant"
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
		chk, err := Check(&a, time.Now(), constant.BaseMinuteDuration)
		if err != nil {
			log.Printf("[ID:%s] error: %s", a.ID, err)
			continue
		}
		if chk == false {
			continue
		}
		err = am.Bot.TextMessageWithRoomKey(a.Message, a.RoomKey)
		if err != nil {
			log.Printf("[ID:%s] error: %s", a.ID, err)
		}
	}

	return nil

}
