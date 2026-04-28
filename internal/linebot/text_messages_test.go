package linebot

import (
	"context"
	"strings"
	"testing"

	"github.com/commojun/nyanbot/internal/masterdata"
	"github.com/commojun/nyanbot/internal/masterdata/key_value"
	"github.com/commojun/nyanbot/internal/masterdata/table"
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
	tma := NewTextMessageAction(bot, event, msg)

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
	tma := NewTextMessageAction(bot, event, msg)

	tma.Do(context.Background())

	if len(bot.replies) != 1 {
		t.Fatalf("返信回数は 1 であるべきだが %d だった", len(bot.replies))
	}
	if bot.replies[0].msg != "これはテストへの返信だよ！！" {
		t.Errorf("testdayo の返信が不正: got %q", bot.replies[0].msg)
	}
}

func TestTextMessageAction_Do_TellIDPrefix_UserSource(t *testing.T) {
	bot := &mockBot{}
	event, msg := newEventWithText("ID")
	tma := NewTextMessageAction(bot, event, msg)

	tma.Do(context.Background())

	if len(bot.replies) != 1 {
		t.Fatalf("返信回数は 1 であるべきだが %d だった", len(bot.replies))
	}
	got := bot.replies[0].msg
	if !strings.Contains(got, "U-test") {
		t.Errorf("UserID が含まれるべきだが含まれなかった: %q", got)
	}
	if strings.Contains(got, "このグループのID") {
		t.Errorf("UserSource では Group ID の記述は含まれないべき: %q", got)
	}
}

func TestTextMessageAction_Do_TellIDPrefix_GroupSource(t *testing.T) {
	bot := &mockBot{}
	event := webhook.MessageEvent{
		ReplyToken: "reply-token-grp",
		Source: webhook.GroupSource{
			UserId:  "U-grouser",
			GroupId: "G-12345",
		},
	}
	msg := webhook.TextMessageContent{Id: "msg-grp", Text: "ID"}
	tma := NewTextMessageAction(bot, event, msg)

	tma.Do(context.Background())

	if len(bot.replies) != 1 {
		t.Fatalf("返信回数は 1 であるべきだが %d だった", len(bot.replies))
	}
	got := bot.replies[0].msg
	if !strings.Contains(got, "U-grouser") {
		t.Errorf("UserID が含まれるべき: %q", got)
	}
	if !strings.Contains(got, "G-12345") {
		t.Errorf("GroupID が含まれるべき: %q", got)
	}
}

func TestTextMessageAction_Do_OnlyFirstMatchFires(t *testing.T) {
	// 先勝ち仕様: 複数 action がマッチ可能でも最初の1つだけ fire する
	bot := &mockBot{}
	event, msg := newEventWithText("てすと") // testdayo にマッチ
	tma := NewTextMessageAction(bot, event, msg)

	tma.Do(context.Background())

	if len(bot.replies) != 1 {
		t.Errorf("先勝ち仕様で返信は1回のみであるべきだが %d 回だった", len(bot.replies))
	}
	// testdayoのメッセージ内容が返る（echo のテキストコピーではない）
	if bot.replies[0].msg != "これはテストへの返信だよ！！" {
		t.Errorf("testdayo が実行されるべきだが echo が実行された: %q", bot.replies[0].msg)
	}
}

// --- masterdata 依存の action テスト ---

func setupMasterData() func() {
	original := getMasterDataSnapshot()
	md := &masterdata.MasterData{
		Tables: &table.Tables{
			Anniversaries: []table.Anniversary{
				{
					ID:      "a1",
					Date:    "2020-01-01",
					Period:  "1",
					Name:    "テスト記念日",
					RoomKey: "test-room",
				},
			},
		},
		KeyVals: &key_value.KVs{
			Nicknames: map[string]string{
				"U-test": "にゃんこもじゅん",
			},
		},
	}
	masterdata.SetTestData(md)
	return func() { masterdata.SetTestData(original) }
}

func getMasterDataSnapshot() *masterdata.MasterData {
	return &masterdata.MasterData{
		Tables:  masterdata.GetTables(),
		KeyVals: masterdata.GetKeyVals(),
	}
}

func TestTextMessageAction_Do_FortunePrefix(t *testing.T) {
	restore := setupMasterData()
	defer restore()

	bot := &mockBot{}
	event, msg := newEventWithText("おみくじ")
	tma := NewTextMessageAction(bot, event, msg)

	tma.Do(context.Background())

	if len(bot.replies) != 1 {
		t.Fatalf("返信回数は 1 であるべきだが %d だった", len(bot.replies))
	}
	got := bot.replies[0].msg
	if !strings.Contains(got, "にゃんこもじゅん") {
		t.Errorf("ニックネームが含まれるべき: %q", got)
	}
	if !strings.Contains(got, "今日の運勢") {
		t.Errorf("運勢フォーマットが含まれるべき: %q", got)
	}
}

func TestTextMessageAction_Do_FortunePrefix_UnknownUser(t *testing.T) {
	restore := setupMasterData()
	defer restore()

	bot := &mockBot{}
	event := webhook.MessageEvent{
		ReplyToken: "reply-token-unknown",
		Source:     webhook.UserSource{UserId: "U-unknown"},
	}
	msg := webhook.TextMessageContent{Id: "msg-u", Text: "おみくじ"}
	tma := NewTextMessageAction(bot, event, msg)

	tma.Do(context.Background())

	if len(bot.replies) != 1 {
		t.Fatalf("返信回数は 1 であるべきだが %d だった", len(bot.replies))
	}
	got := bot.replies[0].msg
	// 未知ユーザーは "あなた" にフォールバック
	if !strings.Contains(got, "あなた") {
		t.Errorf("未知ユーザーは 'あなた' フォールバックするべき: %q", got)
	}
}

func TestTextMessageAction_Do_OjisanPrefix(t *testing.T) {
	restore := setupMasterData()
	defer restore()

	bot := &mockBot{}
	event, msg := newEventWithText("おじさん")
	tma := NewTextMessageAction(bot, event, msg)

	tma.Do(context.Background())

	if len(bot.replies) != 1 {
		t.Fatalf("返信回数は 1 であるべきだが %d だった", len(bot.replies))
	}
	// ojisan の生成物は乱数依存なので内容詳細は検証しない。空文字でないことのみ確認
	if bot.replies[0].msg == "" {
		t.Errorf("おじさんメッセージは空であってはならない")
	}
}

func TestTextMessageAction_Do_AnniversaryPrefix(t *testing.T) {
	restore := setupMasterData()
	defer restore()

	bot := &mockBot{}
	event, msg := newEventWithText("今日は何の日")
	tma := NewTextMessageAction(bot, event, msg)

	tma.Do(context.Background())

	if len(bot.replies) != 1 {
		t.Fatalf("返信回数は 1 であるべきだが %d だった", len(bot.replies))
	}
	got := bot.replies[0].msg
	if !strings.Contains(got, "テスト記念日") {
		t.Errorf("記念日名が含まれるべき: %q", got)
	}
}

func TestTextMessageAction_Do_BotErrorDoesNotPanic(t *testing.T) {
	// Bot がエラーを返しても Do は panic せずログ出力して抜ける
	bot := &mockBot{err: &botError{}}
	event, msg := newEventWithText("てすと")
	tma := NewTextMessageAction(bot, event, msg)

	tma.Do(context.Background())

	if len(bot.replies) != 1 {
		t.Errorf("エラーでも送信試行は行われるべき")
	}
}

type botError struct{}

func (e *botError) Error() string { return "bot failure" }
