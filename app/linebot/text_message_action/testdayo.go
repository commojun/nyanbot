package text_message_action

var (
	testdayo = Action{
		Prefix: "てすと",
		Do:     doTest,
	}
)

func doTest(tma *TextMessageAction) error {
	err := tma.Bot.TextReply(tma.Ctx, "これはテストへの返信だよ！！", tma.Event.ReplyToken)
	if err != nil {
		return err
	}
	return nil
}
