package text_message_action

import (
	"fmt"
	"log"
	"strings"

	origin "github.com/line/line-bot-sdk-go/linebot"
)

type TextMessageAction struct {
	BotClient *origin.Client
	Event     *origin.Event
	Message   *origin.TextMessage
}

func New(cli *origin.Client, event *origin.Event, msg *origin.TextMessage) *TextMessageAction {
	return &TextMessageAction{
		BotClient: cli,
		Event:     event,
		Message:   msg,
	}
}

type Action struct {
	Prefix string
	Do     func(*TextMessageAction) error
}

func actions() []Action {
	return []Action{
		{
			Prefix: "てすと",
			Do:     testdayo,
		},
		{
			Prefix: "ID",
			Do:     tellID,
		},
	}
}

func (tma *TextMessageAction) Do() {
	actFlg := false
	var err error
	for _, action := range actions() {
		if strings.HasPrefix(tma.Message.Text, action.Prefix) {
			err = action.Do(tma)
			if err != nil {
				log.Printf("[text_message_action.Do] error: %s, actionPrefix: %s", err, action.Prefix)
			} else {
				actFlg = true
			}
		}
	}
	if !actFlg {
		err = echo(tma)
		if err != nil {
			log.Printf("[text_message_action.Do] error: %s, echo", err)
		}
	}

	return
}

func testdayo(tma *TextMessageAction) error {
	_, err := tma.BotClient.ReplyMessage(tma.Event.ReplyToken, origin.NewTextMessage("テストへの返信")).Do()
	if err != nil {
		return err
	}
	return nil
}

func tellID(tma *TextMessageAction) error {
	replyText := ""
	if tma.Event.Source.UserID != "" {
		replyText += fmt.Sprintf("あなたのID:%s\n", tma.Event.Source.UserID)
	}
	if tma.Event.Source.RoomID != "" {
		replyText += fmt.Sprintf("この部屋のID:%s\n", tma.Event.Source.RoomID)
	}
	if tma.Event.Source.GroupID != "" {
		replyText += fmt.Sprintf("あなたのID:%s\n", tma.Event.Source.GroupID)
	}
	replyText += "だよ！"
	_, err := tma.BotClient.ReplyMessage(tma.Event.ReplyToken, origin.NewTextMessage(replyText)).Do()
	if err != nil {
		return err
	}
	return nil
}

func echo(tma *TextMessageAction) error {
	_, err := tma.BotClient.ReplyMessage(tma.Event.ReplyToken, origin.NewTextMessage(tma.Message.Text)).Do()
	if err != nil {
		return err
	}
	return nil
}
