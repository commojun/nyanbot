package nyanbot

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
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
