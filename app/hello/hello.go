package hello

import (
	"github.com/commojun/nyanbot/app/linebot"
)

type Hello struct {
}

func New() (*Hello, error) {
	return &Hello{}, nil
}

func (hello *Hello) Say() error {
	bot, err := linebot.New()
	if err != nil {
		return err
	}

	err = bot.TextMessage("Hello!")
	if err != nil {
		return err
	}

	return nil
}
