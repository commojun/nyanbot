package linebot

import (
	"log"

	"github.com/commojun/nyanbot/app/linebot/text_message_action"
	"github.com/commojun/nyanbot/app/redis"
	"github.com/commojun/nyanbot/constant"
	"github.com/commojun/nyanbot/masterdata/key_value"
	origin "github.com/line/line-bot-sdk-go/linebot"
)

type LineBot struct {
	Client        *origin.Client
	DefaultRoomID string
	Events        []*origin.Event
}

func New() (*LineBot, error) {
	botClient, err := origin.New(constant.ChannelSecret, constant.ChannelAccessToken)
	if err != nil {
		return &LineBot{}, err
	}

	var bot = &LineBot{
		Client:        botClient,
		DefaultRoomID: constant.DefaultRoomID,
	}

	return bot, nil
}

func IsInvalidSignature(err error) bool {
	return err == origin.ErrInvalidSignature
}

func (bot *LineBot) TextMessage(msg string) error {
	err := bot.textMessageWithRoomID(msg, bot.DefaultRoomID)
	if err != nil {
		return err
	}
	return nil
}

func (bot *LineBot) TextMessageWithRoomKey(msg string, roomKey string) error {
	redisClient := redis.NewClient()
	roomID, err := redisClient.HGet(key_value.Room, roomKey).Result()
	if err != nil {
		return err
	}

	err = bot.textMessageWithRoomID(msg, roomID)
	if err != nil {
		return err
	}
	return nil
}

func (bot *LineBot) textMessageWithRoomID(msg string, roomID string) error {
	textMsg := origin.NewTextMessage(msg)
	_, err := bot.Client.PushMessage(roomID, textMsg).Do()
	if err != nil {
		return err
	}

	return nil
}

func (bot *LineBot) ActByEvents() {
	for _, event := range bot.Events {
		var err error
		if event.Type == origin.EventTypeMessage {
			switch message := event.Message.(type) {
			case *origin.TextMessage:
				tma := text_message_action.New(bot.Client, event, message)
				tma.Do()
			default:
				log.Printf("[linebot.ActByEvents] message: %s, event: %s", message, event)
			}
		} else {
			log.Printf("[linebot.ActByEvents] event: %s", event)
		}
		if err != nil {
			log.Printf("[linebot.ActByEvents] error: %s, event: %s", err, event)
		}
	}
	return
}
