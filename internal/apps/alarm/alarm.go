package alarm

import (
	"context"
	"log"

	"github.com/commojun/nyanbot/internal/config"
	"github.com/commojun/nyanbot/internal/masterdata"
	"github.com/commojun/nyanbot/internal/masterdata/table"
	"github.com/commojun/nyanbot/internal/timeutil"
)

type LineBot interface {
	TextMessageWithRoomKey(ctx context.Context, msg string, roomKey string) error
}

type AlarmManager struct {
	Alarms                  []table.Alarm
	Bot                     LineBot
	AlarmBaseMinuteDuration int
}

func New(cfg config.Config, bot LineBot) *AlarmManager {
	alms := masterdata.GetTables().Alarms

	return &AlarmManager{
		Alarms:                  alms,
		Bot:                     bot,
		AlarmBaseMinuteDuration: cfg.AlarmBaseMinuteDuration,
	}
}

func (am *AlarmManager) Run(ctx context.Context) error {

	for _, a := range am.Alarms {
		if err := ctx.Err(); err != nil {
			return err
		}
		chk, err := Check(&a, timeutil.LocalTime(), am.AlarmBaseMinuteDuration)
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
