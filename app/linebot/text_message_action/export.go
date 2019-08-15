package text_message_action

import (
	"github.com/commojun/nyanbot/masterdata"
)

var (
	export = Action{
		Prefix: "export",
		Do:     doExport,
	}
)

func doExport(tma *TextMessageAction) error {
	// いったん返答
	initialMsg := "データのエクスポートを開始するよ！"
	err := tma.Bot.TextReply(initialMsg, tma.Event.ReplyToken)
	if err != nil {
		return err
	}

	err = masterdata.Initialize()
	if err != nil {
		errMsg := "データのエクスポートに失敗しちゃったよ。"
		tma.Bot.TextMessageWithRoomID(errMsg, tma.Event.Source.UserID)
		return err
	}

	successMsg := "データのエクスポートが完了したよ！"
	err = tma.Bot.TextMessageWithRoomID(successMsg, tma.Event.Source.UserID)
	if err != nil {
		return err
	}

	return nil
}
