package text_message_action

import (
	"fmt"
)

var (
	tellID = Action{
		Prefix: "ID",
		Do:     doTellID,
	}
)

func doTellID(tma *TextMessageAction) error {
	replyText := ""
	if tma.Event.Source.UserID != "" {
		replyText += fmt.Sprintf("あなたのID: %s\n", tma.Event.Source.UserID)
	}
	if tma.Event.Source.RoomID != "" {
		replyText += fmt.Sprintf("この部屋のID: %s\n", tma.Event.Source.RoomID)
	}
	if tma.Event.Source.GroupID != "" {
		replyText += fmt.Sprintf("このグループのID: %s\n", tma.Event.Source.GroupID)
	}
	replyText += "だよ！"
	err := tma.Bot.TextReply(replyText, tma.Event.ReplyToken)
	if err != nil {
		return err
	}
	return nil
}
