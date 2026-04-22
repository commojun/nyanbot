package alarm

import (
	"context"
	"log"

	"github.com/commojun/nyanbot/app/time_util"
	"github.com/commojun/nyanbot/constant"
	"github.com/commojun/nyanbot/masterdata"
	"github.com/commojun/nyanbot/masterdata/table"
)

type LineBot interface {
	TextMessageWithRoomKey(ctx context.Context, msg string, roomKey string) error
}

type AlarmManager struct {
	Alarms []table.Alarm
	Bot    LineBot
}

func New(bot LineBot) *AlarmManager {
	alms := masterdata.GetTables().Alarms

	am := AlarmManager{
		Alarms: alms,
		Bot:    bot,
	}
	return &am
}

func (am *AlarmManager) Run(ctx context.Context) error {

	for _, a := range am.Alarms {
		if err := ctx.Err(); err != nil {
			return err
		}
		chk, err := Check(&a, time_util.LocalTime(), constant.BaseMinuteDuration)
		if err != nil {
			log.Printf("[ID:%s] error: %s", a.ID, err)
			continue
		}
		if chk == false {
			continue
		}
		err = am.Bot.TextMessageWithRoomKey(ctx, a.Message, a.RoomKey)
		if err != nil {
			log.Printf("[ID:%s] error: %s", a.ID, err)
		}
	}

	return nil

}
