package linebot

import (
	"context"
	"net/http"

	"github.com/commojun/nyanbot/config"
	"github.com/commojun/nyanbot/masterdata"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
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

func (bot *LineBot) TextMessage(ctx context.Context, msg string) error {
	return bot.TextMessageWithRoomID(ctx, msg, bot.DefaultRoomID)
}

func (bot *LineBot) TextMessageWithRoomKey(ctx context.Context, msg string, roomKey string) error {
	roomID, err := masterdata.GetKeyVals().RoomID(roomKey)
	if err != nil {
		return err
	}
	return bot.TextMessageWithRoomID(ctx, msg, roomID)
}

func (bot *LineBot) TextMessageWithRoomID(ctx context.Context, msg string, roomID string) error {
	_, err := bot.Client.WithContext(ctx).PushMessage(&messaging_api.PushMessageRequest{
		To:       roomID,
		Messages: []messaging_api.MessageInterface{messaging_api.TextMessage{Text: msg}},
	}, "")
	return err
}

func (bot *LineBot) TextReply(ctx context.Context, msg string, replyToken string) error {
	_, err := bot.Client.WithContext(ctx).ReplyMessage(&messaging_api.ReplyMessageRequest{
		ReplyToken: replyToken,
		Messages:   []messaging_api.MessageInterface{messaging_api.TextMessage{Text: msg}},
	})
	return err
}

func (bot *LineBot) ParseWebhookRequest(req *http.Request) (*webhook.CallbackRequest, error) {
	return webhook.ParseRequest(bot.ChannelSecret, req)
}
