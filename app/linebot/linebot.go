package linebot

import (
	"github.com/commojun/nyanbot/app/config"
	origin "github.com/line/line-bot-sdk-go/linebot"
)

type LineBot struct {
	Client *origin.Client
	Config *config.Config
}

func New() (*LineBot, error) {
	conf, err := config.Load()
	if err != nil {
		return &LineBot{}, err
	}

	botClient, err := origin.New(conf.ChannelSecret, conf.ChannelAccessToken)
	if err != nil {
		return &LineBot{}, err
	}

	var bot = &LineBot{
		Client: botClient,
		Config: conf,
	}

	return bot, nil
}

func (bot *LineBot) TextMessage(msg string) error {
	textMsg := origin.NewTextMessage(msg)
	_, err := bot.Client.PushMessage(bot.Config.RoomId, textMsg).Do()
	if err != nil {
		return err
	}

	return nil
}
