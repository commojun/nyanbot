package hello

import (
	"github.com/commojun/nyanbot/app/config"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Hello struct {
	Config *config.Config
}

func New() (*Hello, error) {
	conf, err := config.Load()
	if err != nil {
		return &Hello{}, err
	}

	var hello = Hello{}
	hello.Config = conf
	return &hello, nil
}

func (hello *Hello) Say() error {
	bot, err := linebot.New(hello.Config.ChannelSecret, hello.Config.ChannelAccessToken)
	if err != nil {
		return err
	}

	bot.PushMessage(hello.Config.RoomId, linebot.NewTextMessage("Hello nyan!")).Do()
	if err != nil {
		return err
	}

	return nil
}
