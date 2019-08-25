package text_message_action

import (
	"log"
	"strings"

	"github.com/commojun/nyanbot/app/linebot"
	origin "github.com/line/line-bot-sdk-go/linebot"
)

type TextMessageAction struct {
	Bot     *linebot.LineBot
	Event   *origin.Event
	Message *origin.TextMessage
}

func New(bot *linebot.LineBot, event *origin.Event, msg *origin.TextMessage) *TextMessageAction {
	return &TextMessageAction{
		Bot:     bot,
		Event:   event,
		Message: msg,
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
		export,
		ojiline,
		drawFortune,
	}
}

func (tma *TextMessageAction) Do() {
	actFlg := false
	var err error
	for _, action := range actions() {
		if strings.HasPrefix(tma.Message.Text, action.Prefix) {
			err = action.Do(tma)
			actFlg = true
			if err != nil {
				log.Printf("[text_message_action.Do] actionPrefix: %s, error: %s", action.Prefix, err)
			} else {
				log.Printf("[text_message_action.Do] actionPrefix: %s", action.Prefix)
			}
		}
	}
	if !actFlg {
		err = echo(tma)
		if err != nil {
			log.Printf("[text_message_action.Do] action: echo, error: %s, ", err)
		} else {
			log.Printf("[text_message_action.Do] action: echo")
		}
	}

	return
}

func echo(tma *TextMessageAction) error {
	err := tma.Bot.TextReply(tma.Message.Text, tma.Event.ReplyToken)
	if err != nil {
		return err
	}
	return nil
}
