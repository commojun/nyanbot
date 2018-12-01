package nyanbot

import (
	"fmt"
	"strconv"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

func CanSendPushMessage(pmsg PushMessage, now time.Time) (bool, error) {

	config, err := LoadConfig()
	if err != nil {
		return false, err
	}

	weekNum := (int(now.Day()) / 7) + 1

	figs := []string{pmsg.Month, pmsg.WeekNum, pmsg.DayOfWeek, pmsg.DayOfMonth, pmsg.Hour}
	comparison := []int{int(now.Month()), weekNum, int(now.Weekday()), now.Day(), now.Hour()}

	for i, fig := range figs {
		if fig == "*" {
			continue
		}
		fignum, err := strconv.Atoi(fig)
		if err != nil {
			return false, err
		}
		if fignum != comparison[i] {
			fmt.Printf("false because %d != %d \n", fignum, comparison[i])
			return false, nil
		}
	}

	if pmsg.Minute != "*" {
		fignum, err := strconv.Atoi(pmsg.Minute)
		if err != nil {
			return false, err
		}

		d := config.BaseMinuteDuration
		m1 := now.Minute()
		m2 := now.Add(time.Duration(-d) * time.Minute).Minute()
		// 過去d分間を見て当てはまらなかったらfalse
		if (fignum <= (m1-d) || m1 < fignum) &&
			(fignum <= m2 || (m2+d) < fignum) {
			fmt.Printf("false because not match ( %d < %d <= %d ) \n", m2, fignum, m1)
			return false, nil
		}
	}

	return true, nil
}

func SendPushMessage() error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}
	pmsgs, err := LoadPushMessages()
	if err != nil {
		return err
	}

	bot, err := linebot.New(config.ChannelSecret, config.ChannelAccessToken)
	if err != nil {
		return err
	}

	for _, pmsg := range pmsgs {
		can, err := CanSendPushMessage(pmsg, time.Now())
		if err != nil {
			return err
		}
		if can == false {
			continue
		}
		bot.PushMessage(config.RoomId, linebot.NewTextMessage(pmsg.Message)).Do()
		if err != nil {
			return err
		}
	}

	return nil
}
