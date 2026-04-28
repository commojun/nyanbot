package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/commojun/nyanbot/internal/config"
	"github.com/commojun/nyanbot/internal/linebot"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type LineBot interface {
	ParseWebhookRequest(req *http.Request) (*webhook.CallbackRequest, error)
	TextReply(ctx context.Context, msg string, replyToken string) error
}

func makeLineHookAPI(cfg config.Config) API {
	return API{
		Name: "/callback",
		Post: func(ctx context.Context, req *http.Request, res *Response) error {
			bot, err := linebot.New(cfg)
			if err != nil {
				res.Status = http.StatusInternalServerError
				return err
			}
			return handleCallback(ctx, bot, req, res)
		},
	}
}

func handleCallback(ctx context.Context, bot LineBot, req *http.Request, res *Response) error {
	cb, err := bot.ParseWebhookRequest(req)
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
}

func actByLineEvents(ctx context.Context, bot LineBot, events []webhook.EventInterface) {
	for _, event := range events {
		if ctx.Err() != nil {
			return
		}
		switch e := event.(type) {
		case webhook.MessageEvent:
			switch msg := e.Message.(type) {
			case webhook.TextMessageContent:
				tma := linebot.NewTextMessageAction(bot, e, msg)
				tma.Do(ctx)
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
