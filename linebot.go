package nyanbot

import (
	"fmt"
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

func Message() {
	config := LoadConfig()

	bot, err := linebot.New(config.ChannelSecret, config.ChannelAccessToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.PushMessage(config.RoomId, linebot.NewTextMessage("ニャンだよ")).Do()
	if err != nil {
		log.Fatal(err)
	}
}

func PushMessageFromCSV() {
	config := LoadConfig()
	pmsgs := LoadPushMessages()

	bot, err := linebot.New(config.ChannelSecret, config.ChannelAccessToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(pmsgs))

	for _, pmsg := range pmsgs {
		bot.PushMessage(config.RoomId, linebot.NewTextMessage(pmsg.Message)).Do()
		if err != nil {
			log.Fatal(err)
		}
	}
}
