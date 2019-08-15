package api

import (
	"log"
	"net/http"

	"github.com/commojun/nyanbot/app/linebot"
	"github.com/commojun/nyanbot/app/linebot/text_message_action"
	origin "github.com/line/line-bot-sdk-go/linebot"
)

var (
	lineHook = API{
		Name: "/callback",
		Post: postLineHook,
	}
)

func postLineHook(req *http.Request, res *Response) error {
	bot, err := linebot.New()
	if err != nil {
		res.Status = http.StatusInternalServerError
		return err
	}

	bot.Events, err = bot.Client.ParseRequest(req)
	if err != nil {
		return err
	}
	if err != nil {
		if err == origin.ErrInvalidSignature {
			res.Status = http.StatusBadRequest
		} else {
			res.Status = http.StatusInternalServerError
		}
		return err
	}

	actByLineEvents(bot)
	return nil
}

func actByLineEvents(bot *linebot.LineBot) {
	for _, event := range bot.Events {
		var err error
		if event.Type == origin.EventTypeMessage {
			switch message := event.Message.(type) {
			case *origin.TextMessage:
				tma := text_message_action.New(bot, event, message)
				tma.Do()
			case *origin.ImageMessage:
				log.Printf("[linebot.ActByEvents] ImageMessage")
			default:
				log.Printf("[linebot.ActByEvents] message: %s", message)
			}
		} else {
			log.Printf("[linebot.ActByEvents] event: %s", event.Type)
		}
		if err != nil {
			log.Printf("[linebot.ActByEvents] event: %s, error: %s", event, err)
		}
	}
	return
}
