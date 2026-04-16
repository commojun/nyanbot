package text_message_action

import (
	"math/rand"
	"time"

	"github.com/commojun/nyanbot/app/ojisan"
	"github.com/commojun/nyanbot/masterdata"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

var (
	ojiline = Action{
		Prefix: "おじさん",
		Do:     doOjisan,
	}
)

func doOjisan(tma *TextMessageAction) error {
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
		nickname = "にゃんこ"
	}

	rand.Seed(time.Now().UnixNano())
	emojiNum := rand.Intn(9)
	level := rand.Intn(4)

	oji := ojisan.New(nickname, emojiNum, level)
	msg, err := oji.Generate()
	if err != nil {
		return err
	}

	return tma.Bot.TextReply(msg, tma.Event.ReplyToken)
}
