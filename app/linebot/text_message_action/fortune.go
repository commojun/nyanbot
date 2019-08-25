package text_message_action

import (
	"fmt"

	"github.com/commojun/nyanbot/app/fortune"
	"github.com/commojun/nyanbot/app/redis"
	"github.com/commojun/nyanbot/masterdata/key_value"
)

var (
	drawFortune = Action{
		Prefix: "おみくじ",
		Do:     doDrawFortune,
	}
)

func doDrawFortune(tma *TextMessageAction) error {
	// 名前取得
	redisClient := redis.NewClient()
	nickname, err := redisClient.HGet(key_value.Nickname, tma.Event.Source.UserID).Result()
	if err != nil {
		nickname = "あなた"
	}

	f := fortune.New()
	result := f.DrawByStringSeed(tma.Event.Source.UserID)

	msg := fmt.Sprintf("%sの今日の運勢\n>>>%s<<<", nickname, result)

	err = tma.Bot.TextReply(msg, tma.Event.ReplyToken)
	if err != nil {
		return err
	}

	return nil
}
