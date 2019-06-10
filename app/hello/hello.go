package hello

import (
	"github.com/commojun/nyanbot/app/linebot"
)

type Hello struct {
	Bot *linebot.LineBot
}

func New() (*Hello, error) {
	bot, err := linebot.New()
	if err != nil {
		return &Hello{}, err
	}
	return &Hello{Bot: bot}, nil
}

func (hello *Hello) Say() error {

	err := hello.Bot.TextMessage("Hello!")
	if err != nil {
		return err
	}

	return nil
}
