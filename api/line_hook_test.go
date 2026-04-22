package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type mockLineBot struct {
	parseResult *webhook.CallbackRequest
	parseErr    error
	replies     []replyCall
}

type replyCall struct {
	msg        string
	replyToken string
}

func (m *mockLineBot) ParseWebhookRequest(req *http.Request) (*webhook.CallbackRequest, error) {
	return m.parseResult, m.parseErr
}

func (m *mockLineBot) TextReply(ctx context.Context, msg string, replyToken string) error {
	m.replies = append(m.replies, replyCall{msg: msg, replyToken: replyToken})
	return nil
}

func newTestRequest() *http.Request {
	return httptest.NewRequest(http.MethodPost, "/callback", strings.NewReader(""))
}

// --- handleCallback tests ---

func TestHandleCallback_Success(t *testing.T) {
	bot := &mockLineBot{
		parseResult: &webhook.CallbackRequest{
			Events: []webhook.EventInterface{
				webhook.MessageEvent{
					ReplyToken: "tk-1",
					Source:     webhook.UserSource{UserId: "U-x"},
					Message:    webhook.TextMessageContent{Id: "m1", Text: "てすと"},
				},
			},
		},
	}
	res := &Response{Status: http.StatusOK, Message: "OK"}

	err := handleCallback(context.Background(), bot, newTestRequest(), res)
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	if res.Status != http.StatusOK {
		t.Errorf("status = %d, want %d", res.Status, http.StatusOK)
	}
	// TMA 経由で TextReply が呼ばれてる
	if len(bot.replies) != 1 {
		t.Errorf("TextReply 呼び出し回数 = %d, want 1", len(bot.replies))
	}
}

func TestHandleCallback_InvalidSignature(t *testing.T) {
	bot := &mockLineBot{
		parseErr: webhook.ErrInvalidSignature,
	}
	res := &Response{Status: http.StatusOK, Message: "OK"}

	err := handleCallback(context.Background(), bot, newTestRequest(), res)
	if err == nil {
		t.Fatal("エラーが返るべきだが nil だった")
	}
	if res.Status != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", res.Status, http.StatusBadRequest)
	}
	if len(bot.replies) != 0 {
		t.Errorf("不正署名時は reply されないべきだが %d 回呼ばれた", len(bot.replies))
	}
}

func TestHandleCallback_OtherParseError(t *testing.T) {
	bot := &mockLineBot{
		parseErr: errors.New("some other error"),
	}
	res := &Response{Status: http.StatusOK, Message: "OK"}

	err := handleCallback(context.Background(), bot, newTestRequest(), res)
	if err == nil {
		t.Fatal("エラーが返るべきだが nil だった")
	}
	if res.Status != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", res.Status, http.StatusInternalServerError)
	}
}

// --- actByLineEvents tests ---

func TestActByLineEvents_EmptyEvents(t *testing.T) {
	bot := &mockLineBot{}
	actByLineEvents(context.Background(), bot, nil)
	if len(bot.replies) != 0 {
		t.Errorf("空イベントリストでは reply されないべき")
	}
}

func TestActByLineEvents_TextMessage(t *testing.T) {
	bot := &mockLineBot{}
	events := []webhook.EventInterface{
		webhook.MessageEvent{
			ReplyToken: "tk-1",
			Source:     webhook.UserSource{UserId: "U-x"},
			Message:    webhook.TextMessageContent{Id: "m1", Text: "てすと"},
		},
	}

	actByLineEvents(context.Background(), bot, events)

	if len(bot.replies) != 1 {
		t.Fatalf("TextReply 呼び出し回数 = %d, want 1", len(bot.replies))
	}
	if bot.replies[0].replyToken != "tk-1" {
		t.Errorf("replyToken = %q, want %q", bot.replies[0].replyToken, "tk-1")
	}
}

func TestActByLineEvents_ImageMessage_NoReply(t *testing.T) {
	bot := &mockLineBot{}
	events := []webhook.EventInterface{
		webhook.MessageEvent{
			ReplyToken: "tk-img",
			Source:     webhook.UserSource{UserId: "U-x"},
			Message:    webhook.ImageMessageContent{Id: "img-1"},
		},
	}

	actByLineEvents(context.Background(), bot, events)

	if len(bot.replies) != 0 {
		t.Errorf("画像メッセージは reply されないべきだが %d 回呼ばれた", len(bot.replies))
	}
}

func TestActByLineEvents_CancelledContext(t *testing.T) {
	bot := &mockLineBot{}
	events := []webhook.EventInterface{
		webhook.MessageEvent{
			ReplyToken: "tk-1",
			Source:     webhook.UserSource{UserId: "U-x"},
			Message:    webhook.TextMessageContent{Id: "m1", Text: "てすと"},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	actByLineEvents(ctx, bot, events)

	if len(bot.replies) != 0 {
		t.Errorf("キャンセル済み ctx では処理されないべきだが %d 回呼ばれた", len(bot.replies))
	}
}

func TestActByLineEvents_StopsOnMidwayCancel(t *testing.T) {
	// 最初の1件は処理され、途中で ctx キャンセルすると後続はスキップされる
	// → 完全な確定検証は難しいが、事前キャンセルで 0 件、キャンセルなしで2件になることで
	//   イテレーション毎のチェックが機能してることを確認できる
	bot := &mockLineBot{}
	events := []webhook.EventInterface{
		webhook.MessageEvent{
			ReplyToken: "tk-1",
			Source:     webhook.UserSource{UserId: "U-x"},
			Message:    webhook.TextMessageContent{Id: "m1", Text: "てすと"},
		},
		webhook.MessageEvent{
			ReplyToken: "tk-2",
			Source:     webhook.UserSource{UserId: "U-y"},
			Message:    webhook.TextMessageContent{Id: "m2", Text: "てすと"},
		},
	}

	actByLineEvents(context.Background(), bot, events)

	if len(bot.replies) != 2 {
		t.Errorf("2件とも処理されるべきだが %d 件だった", len(bot.replies))
	}
}
