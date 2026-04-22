package text_message_action

import (
	"context"

	"github.com/commojun/nyanbot/app/anniversary"
)

var (
	randomAnniversary = Action{
		Prefix: "今日は何の日",
		Do:     doRandomAnniversary,
	}
)

func doRandomAnniversary(ctx context.Context, tma *TextMessageAction) error {
	msg, err := anniversary.RandomMsg()
	if err != nil {
		return err
	}

	return tma.Bot.TextReply(ctx, msg, tma.Event.ReplyToken)
}
