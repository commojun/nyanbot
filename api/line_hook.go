package api

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/commojun/nyanbot/app/linebot"
	"github.com/commojun/nyanbot/app/linebot/text_message_action"
	"github.com/commojun/nyanbot/config"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

func makeLineHookAPI(cfg config.Config) API {
	return API{
		Name: "/callback",
		Post: func(ctx context.Context, req *http.Request, res *Response) error {
			bot, err := linebot.New(cfg)
			if err != nil {
				res.Status = http.StatusInternalServerError
				return err
			}

			cb, err := webhook.ParseRequest(bot.ChannelSecret, req)
			if err != nil {
				if errors.Is(err, webhook.ErrInvalidSignature) {
					res.Status = http.StatusBadRequest
				} else {
					res.Status = http.StatusInternalServerError
				}
				return err
			}

			actByLineEvents(ctx, bot, cb.Events)
			return nil
		},
	}
}

func actByLineEvents(ctx context.Context, bot *linebot.LineBot, events []webhook.EventInterface) {
	for _, event := range events {
		switch e := event.(type) {
		case webhook.MessageEvent:
			switch msg := e.Message.(type) {
			case webhook.TextMessageContent:
				tma := text_message_action.New(ctx, bot, e, msg)
				tma.Do()
			case webhook.ImageMessageContent:
				log.Printf("[linebot.ActByEvents] ImageMessage")
			default:
				log.Printf("[linebot.ActByEvents] message: %+v", msg)
			}
		default:
			log.Printf("[linebot.ActByEvents] event: %+v", e)
		}
	}
}
