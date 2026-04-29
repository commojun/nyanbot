package linebot

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/commojun/nyanbot/internal/apps/anniversary"
	"github.com/commojun/nyanbot/internal/apps/fortune"
	"github.com/commojun/nyanbot/internal/apps/ojisan"
	"github.com/commojun/nyanbot/internal/apps/weather"
	"github.com/commojun/nyanbot/internal/masterdata"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

// TextBot: テキストメッセージへの返信に必要なインターフェース
type TextBot interface {
	TextReply(ctx context.Context, msg string, replyToken string) error
}

type TextMessageAction struct {
	Bot     TextBot
	Event   webhook.MessageEvent
	Message webhook.TextMessageContent
}

func NewTextMessageAction(bot TextBot, event webhook.MessageEvent, msg webhook.TextMessageContent) *TextMessageAction {
	return &TextMessageAction{
		Bot:     bot,
		Event:   event,
		Message: msg,
	}
}

type textAction struct {
	Prefix string
	Do     func(context.Context, *TextMessageAction) error
}

func actions() []textAction {
	return []textAction{
		{Prefix: "てすと", Do: doTest},
		{Prefix: "ID", Do: doTellID},
		{Prefix: "おじさん", Do: doOjisan},
		{Prefix: "おみくじ", Do: doDrawFortune},
		{Prefix: "今日は何の日", Do: doRandomAnniversary},
		{Prefix: "天気", Do: doGetWeather},
	}
}

func (tma *TextMessageAction) Do(ctx context.Context) {
	// 先勝ち: 最初にマッチしたactionのみ実行する (LINE の ReplyToken が1回限りのため)
	for _, action := range actions() {
		if strings.HasPrefix(tma.Message.Text, action.Prefix) {
			err := action.Do(ctx, tma)
			if err != nil {
				log.Printf("[TextMessageAction.Do] actionPrefix: %s, error: %s", action.Prefix, err)
			} else {
				log.Printf("[TextMessageAction.Do] actionPrefix: %s", action.Prefix)
			}
			return
		}
	}
	// どのactionにもマッチしなかった場合のみechoにフォールバック
	err := echo(ctx, tma)
	if err != nil {
		log.Printf("[TextMessageAction.Do] action: echo, error: %s, ", err)
	} else {
		log.Printf("[TextMessageAction.Do] action: echo")
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

func doTest(ctx context.Context, tma *TextMessageAction) error {
	return tma.Bot.TextReply(ctx, "これはテストへの返信だよ！！", tma.Event.ReplyToken)
}

func doTellID(ctx context.Context, tma *TextMessageAction) error {
	replyText := ""
	switch src := tma.Event.Source.(type) {
	case webhook.UserSource:
		replyText += fmt.Sprintf("あなたのID: %s\n", src.UserId)
	case webhook.GroupSource:
		replyText += fmt.Sprintf("あなたのID: %s\n", src.UserId)
		replyText += fmt.Sprintf("このグループのID: %s\n", src.GroupId)
	}
	replyText += "だよ！"
	return tma.Bot.TextReply(ctx, replyText, tma.Event.ReplyToken)
}

func doDrawFortune(ctx context.Context, tma *TextMessageAction) error {
	userID := extractUserID(tma.Event.Source)

	nickname, err := masterdata.GetKeyVals().Nickname(userID)
	if err != nil {
		nickname = "あなた"
	}

	f := fortune.New()
	result := f.DrawByStringSeed(userID)

	msg := fmt.Sprintf("%sの今日の運勢\n>>>%s<<<", nickname, result)

	return tma.Bot.TextReply(ctx, msg, tma.Event.ReplyToken)
}

func doOjisan(ctx context.Context, tma *TextMessageAction) error {
	userID := extractUserID(tma.Event.Source)

	nickname, err := masterdata.GetKeyVals().Nickname(userID)
	if err != nil {
		nickname = "にゃんこ"
	}

	rand.Seed(time.Now().UnixNano()) //nolint:staticcheck
	emojiNum := rand.Intn(9)
	level := rand.Intn(4)

	oji := ojisan.New(nickname, emojiNum, level)
	msg, err := oji.Generate()
	if err != nil {
		return err
	}

	return tma.Bot.TextReply(ctx, msg, tma.Event.ReplyToken)
}

func doRandomAnniversary(ctx context.Context, tma *TextMessageAction) error {
	msg, err := anniversary.RandomMsg()
	if err != nil {
		return err
	}

	return tma.Bot.TextReply(ctx, msg, tma.Event.ReplyToken)
}

func doGetWeather(ctx context.Context, tma *TextMessageAction) error {
	const cityID = "140010" // 横浜

	msg, err := weather.Fetch(ctx, cityID)
	if err != nil {
		return tma.Bot.TextReply(ctx, "天気の取得に失敗したよ…", tma.Event.ReplyToken)
	}

	return tma.Bot.TextReply(ctx, msg, tma.Event.ReplyToken)
}
