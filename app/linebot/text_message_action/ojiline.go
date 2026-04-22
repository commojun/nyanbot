package text_message_action

import (
	"context"
	"math/rand"
	"time"

	"github.com/commojun/nyanbot/app/ojisan"
	"github.com/commojun/nyanbot/masterdata"
)

var (
	ojiline = Action{
		Prefix: "おじさん",
		Do:     doOjisan,
	}
)

func doOjisan(ctx context.Context, tma *TextMessageAction) error {
	userID := extractUserID(tma.Event.Source)

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

	return tma.Bot.TextReply(ctx, msg, tma.Event.ReplyToken)
}
