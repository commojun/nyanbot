package text_message_action

import (
	"math/rand"
	"time"

	"github.com/commojun/nyanbot/app/ojisan"
	"github.com/commojun/nyanbot/app/redis"
	"github.com/commojun/nyanbot/masterdata/key_value"
)

var (
	ojiline = Action{
		Prefix: "おじさん",
		Do:     doOjisan,
	}
)

func doOjisan(tma *TextMessageAction) error {
	// 名前取得
	redisClient := redis.NewClient()
	nickname, err := redisClient.HGet(key_value.Nickname, tma.Event.Source.UserID).Result()
	if err != nil {
		nickname = "にゃんこ"
	}

	rand.Seed(time.Now().UnixNano())
	emojiNum := rand.Intn(9)
	level := rand.Intn(4)

	oji := ojisan.New(nickname, emojiNum, level)
	msg, err := oji.Generate()
	if err != nil {
		return err
	}

	err = tma.Bot.TextReply(msg, tma.Event.ReplyToken)
	if err != nil {
		return err
	}

	return nil
}
