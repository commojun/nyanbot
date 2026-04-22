package text_message_action

import (
	"context"
	"fmt"

	"github.com/commojun/nyanbot/app/fortune"
	"github.com/commojun/nyanbot/masterdata"
)

var (
	drawFortune = Action{
		Prefix: "おみくじ",
		Do:     doDrawFortune,
	}
)

func doDrawFortune(ctx context.Context, tma *TextMessageAction) error {
	userID := extractUserID(tma.Event.Source)

	nickname, err := masterdata.GetKeyVals().Nickname(userID)
	if err != nil {
		nickname = "あなた"
	}

	f := fortune.New()
	result := f.DrawByStringSeed(userID)

	msg := fmt.Sprintf("%sの今日の運勢\n>>>%s<<<", nickname, result)

	return tma.Bot.TextReply(ctx, msg, tma.Event.ReplyToken)
}
