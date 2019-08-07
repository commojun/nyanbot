package linebot

import (
	"github.com/commojun/nyanbot/constant"
	origin "github.com/line/line-bot-sdk-go/linebot"
)

type LineBot struct {
	Client *origin.Client
	RoomId string
}

func New() (*LineBot, error) {
	botClient, err := origin.New(constant.ChannelSecret, constant.ChannelAccessToken)
	if err != nil {
		return &LineBot{}, err
	}

	var bot = &LineBot{
		Client: botClient,
		RoomId: constant.RoomId,
	}

	return bot, nil
}

func (bot *LineBot) TextMessage(msg string) error {
	textMsg := origin.NewTextMessage(msg)
	_, err := bot.Client.PushMessage(bot.RoomId, textMsg).Do()
	if err != nil {
		return err
	}

	return nil
}
