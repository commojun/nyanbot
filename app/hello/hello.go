package hello

import (
	"github.com/commojun/nyanbot/app/linebot"
)

type Hello struct {
	Bot *linebot.LineBot
}

func New(bot *linebot.LineBot) *Hello {
	return &Hello{Bot: bot}
}

func (hello *Hello) Say() error {

	err := hello.Bot.TextMessage("Hello!")
	if err != nil {
		return err
	}

	return nil
}
