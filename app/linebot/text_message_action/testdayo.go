package text_message_action

import "context"

var (
	testdayo = Action{
		Prefix: "てすと",
		Do:     doTest,
	}
)

func doTest(ctx context.Context, tma *TextMessageAction) error {
	return tma.Bot.TextReply(ctx, "これはテストへの返信だよ！！", tma.Event.ReplyToken)
}
