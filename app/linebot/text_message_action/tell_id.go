package text_message_action

import (
	"context"
	"fmt"

	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

var (
	tellID = Action{
		Prefix: "ID",
		Do:     doTellID,
	}
)

func doTellID(ctx context.Context, tma *TextMessageAction) error {
	replyText := ""
	switch src := tma.Event.Source.(type) {
	case webhook.UserSource:
		replyText += fmt.Sprintf("あなたのID: %s\n", src.UserId)
	case webhook.GroupSource:
		replyText += fmt.Sprintf("あなたのID: %s\n", src.UserId)
		replyText += fmt.Sprintf("このグループのID: %s\n", src.GroupId)
	}
	replyText += "だよ！"
	return tma.Bot.TextReply(ctx, replyText, tma.Event.ReplyToken)
}
