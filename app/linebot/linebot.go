package linebot

import (
	"github.com/commojun/nyanbot/config"
	"github.com/commojun/nyanbot/masterdata"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

type LineBot struct {
	Client        *messaging_api.MessagingApiAPI
	ChannelSecret string
	DefaultRoomID string
}

func New(cfg config.Config) (*LineBot, error) {
	client, err := messaging_api.NewMessagingApiAPI(cfg.ChannelAccessToken)
	if err != nil {
		return &LineBot{}, err
	}

	return &LineBot{
		Client:        client,
		ChannelSecret: cfg.ChannelSecret,
		DefaultRoomID: cfg.DefaultRoomID,
	}, nil
}

func (bot *LineBot) TextMessage(msg string) error {
	return bot.TextMessageWithRoomID(msg, bot.DefaultRoomID)
}

func (bot *LineBot) TextMessageWithRoomKey(msg string, roomKey string) error {
	roomID, err := masterdata.GetKeyVals().RoomID(roomKey)
	if err != nil {
		return err
	}
	return bot.TextMessageWithRoomID(msg, roomID)
}

func (bot *LineBot) TextMessageWithRoomID(msg string, roomID string) error {
	_, err := bot.Client.PushMessage(&messaging_api.PushMessageRequest{
		To:       roomID,
		Messages: []messaging_api.MessageInterface{messaging_api.TextMessage{Text: msg}},
	}, "")
	return err
}

func (bot *LineBot) TextReply(msg string, replyToken string) error {
	_, err := bot.Client.ReplyMessage(&messaging_api.ReplyMessageRequest{
		ReplyToken: replyToken,
		Messages:   []messaging_api.MessageInterface{messaging_api.TextMessage{Text: msg}},
	})
	return err
}
