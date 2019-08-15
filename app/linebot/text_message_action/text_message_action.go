package text_message_action

import (
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
		testdayo,
		tellID,
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
				log.Printf("[text_message_action.Do] actionPrefix: %s", action.Prefix)
				actFlg = true
			}
		}
	}
	if !actFlg {
		err = echo(tma)
		if err != nil {
			log.Printf("[text_message_action.Do] error: %s, echo", err)
		} else {
			log.Printf("[text_message_action.Do] echo")
		}
	}

	return
}

func echo(tma *TextMessageAction) error {
	_, err := tma.BotClient.ReplyMessage(tma.Event.ReplyToken, origin.NewTextMessage(tma.Message.Text)).Do()
	if err != nil {
		return err
	}
	return nil
}
