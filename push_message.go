package nyanbot

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

func CanSendPushMessage(pmsg PushMessgae) bool {

	config := LoadConfig()
	// 現在の日付時刻を取得
	now := time.Now()
	weekNum := (int(now.Day()) / 7) + 1

	figs := []string{pmsg.Month, pmsg.WeekNum, pmsg.DayOfWeek, pmsg.DayOfMonth, pmsg.Hour}
	comparison := []int{int(now.Month()), weekNum, now.Day(), int(now.Weekday()), now.Hour()}

	for i, fig := range figs {
		if fig == "*" {
			continue
		}
		fignum, err := strconv.Atoi(fig)
		if err != nil {
			log.Fatal(err)
		}
		if fignum != comparison[i] {
			fmt.Printf("false because %d != %d \n", fignum, comparison[i])
			return false
		}
	}

	if pmsg.Minute != "*" {
		fignum, err := strconv.Atoi(pmsg.Minute)
		if err != nil {
			log.Fatal(err)
		}

		d := config.BaseMinuteDuration
		m1 := now.Minute()
		m2 := now.Add(time.Duration(-d) * time.Minute).Minute()
		// 過去d分間を見て当てはまらなかったらfalse
		if (fignum < (m1-d) && m1 < fignum) &&
			(fignum < m2 && (m2+d) < fignum) {
			fmt.Printf("false because not match ( %d <= %d <= %d ) \n", m2, fignum, m1)
			return false
		}
	}

	return true
}

func SendPushMessage() {
	config := LoadConfig()
	pmsgs := LoadPushMessages()

	bot, err := linebot.New(config.ChannelSecret, config.ChannelAccessToken)
	if err != nil {
		log.Fatal(err)
	}

	for _, pmsg := range pmsgs {
		if CanSendPushMessage(pmsg) == false {
			continue
		}
		bot.PushMessage(config.RoomId, linebot.NewTextMessage(pmsg.Message)).Do()
		if err != nil {
			log.Fatal(err)
		}
	}
}
