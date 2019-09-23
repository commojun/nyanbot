package text_message_action

import "github.com/commojun/nyanbot/app/anniversary"

var (
	randomAnniversary = Action{
		Prefix: "今日は何の日",
		Do:     doRandomAnniversary,
	}
)

func doRandomAnniversary(tma *TextMessageAction) error {
	am, err := anniversary.New()
	if err != nil {
		return err
	}

	msg, err := am.RandomMsg()
	if err != nil {
		return err
	}

	err = tma.Bot.TextReply(msg, tma.Event.ReplyToken)
	if err != nil {
		return err
	}

	return nil
}
