package linebot

import (
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
	err := bot.TextMessageWithRoomID(msg, bot.DefaultRoomID)
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

	err = bot.TextMessageWithRoomID(msg, roomID)
	if err != nil {
		return err
	}
	return nil
}

func (bot *LineBot) TextMessageWithRoomID(msg string, roomID string) error {
	textMsg := origin.NewTextMessage(msg)
	_, err := bot.Client.PushMessage(roomID, textMsg).Do()
	if err != nil {
		return err
	}

	return nil
}

func (bot *LineBot) TextReply(msg string, replyToken string) error {
	textMsg := origin.NewTextMessage(msg)
	_, err := bot.Client.ReplyMessage(replyToken, textMsg).Do()
	if err != nil {
		return err
	}

	return nil
}
