package linebot

import (
	"github.com/commojun/nyanbot/cache"
	"github.com/commojun/nyanbot/config"
	origin "github.com/line/line-bot-sdk-go/linebot"
)

type LineBot struct {
	Client        *origin.Client
	DefaultRoomID string
	Events        []*origin.Event
}

func New(cfg config.Config) (*LineBot, error) {
	botClient, err := origin.New(cfg.ChannelSecret, cfg.ChannelAccessToken)
	if err != nil {
		return &LineBot{}, err
	}

	var bot = &LineBot{
		Client:        botClient,
		DefaultRoomID: cfg.DefaultRoomID,
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
	roomID, err := cache.GetRoomID(roomKey)
	if err != nil {
		return err
	}
	return bot.TextMessageWithRoomID(msg, roomID)
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
