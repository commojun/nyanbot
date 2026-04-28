package hello

import (
	"context"

	"github.com/commojun/nyanbot/internal/linebot"
)

type Hello struct {
	Bot *linebot.LineBot
}

func New(bot *linebot.LineBot) *Hello {
	return &Hello{Bot: bot}
}

func (hello *Hello) Say(ctx context.Context) error {

	err := hello.Bot.TextMessage(ctx, "Hello!")
	if err != nil {
		return err
	}

	return nil
}
