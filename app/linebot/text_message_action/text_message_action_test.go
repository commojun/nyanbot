package text_message_action

import (
	"context"
	"strings"
	"testing"

	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type mockBot struct {
	replies []replyCall
	err     error
}

type replyCall struct {
	msg        string
	replyToken string
}

func (m *mockBot) TextReply(ctx context.Context, msg string, replyToken string) error {
	m.replies = append(m.replies, replyCall{msg: msg, replyToken: replyToken})
	return m.err
}

func newEventWithText(text string) (webhook.MessageEvent, webhook.TextMessageContent) {
	event := webhook.MessageEvent{
		ReplyToken: "reply-token-xyz",
		Source:     webhook.UserSource{UserId: "U-test"},
	}
	msg := webhook.TextMessageContent{
		Id:   "msg-1",
		Text: text,
	}
	return event, msg
}

func TestTextMessageAction_Do_EchoFallback(t *testing.T) {
	bot := &mockBot{}
	event, msg := newEventWithText("プレフィックスにマッチしないテキスト")
	tma := New(bot, event, msg)

	tma.Do(context.Background())

	if len(bot.replies) != 1 {
		t.Fatalf("返信回数は 1 であるべきだが %d だった", len(bot.replies))
	}
	got := bot.replies[0]
	if got.msg != "プレフィックスにマッチしないテキスト" {
		t.Errorf("echo メッセージが不正: got %q", got.msg)
	}
	if got.replyToken != "reply-token-xyz" {
		t.Errorf("replyToken が不正: got %q", got.replyToken)
	}
}

func TestTextMessageAction_Do_TestdayoPrefix(t *testing.T) {
	bot := &mockBot{}
	event, msg := newEventWithText("てすと")
	tma := New(bot, event, msg)

	tma.Do(context.Background())

	if len(bot.replies) != 1 {
		t.Fatalf("返信回数は 1 であるべきだが %d だった", len(bot.replies))
	}
	if bot.replies[0].msg != "これはテストへの返信だよ！！" {
		t.Errorf("testdayo の返信が不正: got %q", bot.replies[0].msg)
	}
}

func TestTextMessageAction_Do_TellIDPrefix(t *testing.T) {
	bot := &mockBot{}
	event, msg := newEventWithText("ID")
	tma := New(bot, event, msg)

	tma.Do(context.Background())

	if len(bot.replies) != 1 {
		t.Fatalf("返信回数は 1 であるべきだが %d だった", len(bot.replies))
	}
	got := bot.replies[0].msg
	if !strings.Contains(got, "U-test") {
		t.Errorf("UserID が含まれるべきだが含まれなかった: %q", got)
	}
}

func TestTextMessageAction_Do_NoMatchNotDoubled(t *testing.T) {
	// echo fallback 時に通常 action が走らないことを確認
	bot := &mockBot{}
	event, msg := newEventWithText("hoge")
	tma := New(bot, event, msg)

	tma.Do(context.Background())

	if len(bot.replies) != 1 {
		t.Errorf("echo のみで 1 回だけ返信するべきだが %d 回だった", len(bot.replies))
	}
}

func TestTextMessageAction_Do_BotErrorDoesNotPanic(t *testing.T) {
	// Bot がエラーを返しても Do は panic せずログ出力して抜ける
	bot := &mockBot{err: &botError{}}
	event, msg := newEventWithText("てすと")
	tma := New(bot, event, msg)

	tma.Do(context.Background())

	if len(bot.replies) != 1 {
		t.Errorf("エラーでも送信試行は行われるべき")
	}
}

type botError struct{}

func (e *botError) Error() string { return "bot failure" }
