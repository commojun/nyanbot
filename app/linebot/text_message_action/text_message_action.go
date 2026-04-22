package text_message_action

import (
	"context"
	"log"
	"strings"

	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type LineBot interface {
	TextReply(ctx context.Context, msg string, replyToken string) error
}

type TextMessageAction struct {
	Bot     LineBot
	Event   webhook.MessageEvent
	Message webhook.TextMessageContent
}

func New(bot LineBot, event webhook.MessageEvent, msg webhook.TextMessageContent) *TextMessageAction {
	return &TextMessageAction{
		Bot:     bot,
		Event:   event,
		Message: msg,
	}
}

type Action struct {
	Prefix string
	Do     func(context.Context, *TextMessageAction) error
}

func actions() []Action {
	return []Action{
		testdayo,
		tellID,
		ojiline,
		drawFortune,
		randomAnniversary,
		getWeather,
	}
}

func (tma *TextMessageAction) Do(ctx context.Context) {
	// 先勝ち: 最初にマッチしたactionのみ実行する (LINE の ReplyToken が1回限りのため)
	for _, action := range actions() {
		if strings.HasPrefix(tma.Message.Text, action.Prefix) {
			err := action.Do(ctx, tma)
			if err != nil {
				log.Printf("[text_message_action.Do] actionPrefix: %s, error: %s", action.Prefix, err)
			} else {
				log.Printf("[text_message_action.Do] actionPrefix: %s", action.Prefix)
			}
			return
		}
	}
	// どのactionにもマッチしなかった場合のみechoにフォールバック
	err := echo(ctx, tma)
	if err != nil {
		log.Printf("[text_message_action.Do] action: echo, error: %s, ", err)
	} else {
		log.Printf("[text_message_action.Do] action: echo")
	}
}

func echo(ctx context.Context, tma *TextMessageAction) error {
	return tma.Bot.TextReply(ctx, tma.Message.Text, tma.Event.ReplyToken)
}

func extractUserID(src webhook.SourceInterface) string {
	switch s := src.(type) {
	case webhook.UserSource:
		return s.UserId
	case webhook.GroupSource:
		return s.UserId
	case webhook.RoomSource:
		return s.UserId
	}
	return ""
}
