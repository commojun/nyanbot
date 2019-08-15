package text_message_action

import origin "github.com/line/line-bot-sdk-go/linebot"

var (
	testdayo = Action{
		Prefix: "てすと",
		Do:     doTest,
	}
)

func doTest(tma *TextMessageAction) error {
	_, err := tma.BotClient.ReplyMessage(tma.Event.ReplyToken, origin.NewTextMessage("テストへの返信")).Do()
	if err != nil {
		return err
	}
	return nil
}
