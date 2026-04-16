package text_message_action

import (
	"fmt"

	"github.com/commojun/nyanbot/app/fortune"
	"github.com/commojun/nyanbot/masterdata"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

var (
	drawFortune = Action{
		Prefix: "おみくじ",
		Do:     doDrawFortune,
	}
)

func doDrawFortune(tma *TextMessageAction) error {
	// UserID取得
	var userID string
	switch src := tma.Event.Source.(type) {
	case webhook.UserSource:
		userID = src.UserId
	case webhook.GroupSource:
		userID = src.UserId
	}

	// 名前取得
	nickname, err := masterdata.GetKeyVals().Nickname(userID)
	if err != nil {
		nickname = "あなた"
	}

	f := fortune.New()
	result := f.DrawByStringSeed(userID)

	msg := fmt.Sprintf("%sの今日の運勢\n>>>%s<<<", nickname, result)

	return tma.Bot.TextReply(msg, tma.Event.ReplyToken)
}
