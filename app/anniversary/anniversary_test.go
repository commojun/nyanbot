package anniversary

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/commojun/nyanbot/app/time_util"
	"github.com/commojun/nyanbot/masterdata/table"
)

type mockBot struct {
	sent []sentMsg
	err  error
}

type sentMsg struct {
	msg     string
	roomKey string
}

func (m *mockBot) TextMessageWithRoomKey(ctx context.Context, msg string, roomKey string) error {
	m.sent = append(m.sent, sentMsg{msg: msg, roomKey: roomKey})
	return m.err
}

var jst = time.FixedZone("JST", 9*60*60)

func TestMakeCheckMessage(t *testing.T) {
	t.Run("記念日当日（月日一致・複数年後）: check=true, 年数メッセージ", func(t *testing.T) {
		a := &table.Anniversary{
			ID:     "1",
			Date:   "2020-04-21",
			Period: "100",
			Name:   "テスト記念日",
		}
		now := time.Date(2026, time.April, 21, 12, 0, 0, 0, jst)
		msg, check, err := MakeCheckMessage(a, now)
		if err != nil {
			t.Fatalf("予期しないエラー: %v", err)
		}
		if !check {
			t.Errorf("check は true であるべきだが false だった")
		}
		expected := fmt.Sprintf(TEMPLATE, a.Name, 6, "年")
		if msg != expected {
			t.Errorf("メッセージが不正: got %q, want %q", msg, expected)
		}
	})

	t.Run("記念日当日（同年同月同日・0年）: check=true, 0年メッセージ", func(t *testing.T) {
		a := &table.Anniversary{
			ID:     "2",
			Date:   "2026-04-21",
			Period: "100",
			Name:   "ゼロ年記念",
		}
		now := time.Date(2026, time.April, 21, 0, 0, 0, 0, jst)
		msg, check, err := MakeCheckMessage(a, now)
		if err != nil {
			t.Fatalf("予期しないエラー: %v", err)
		}
		if !check {
			t.Errorf("check は true であるべきだが false だった")
		}
		expected := fmt.Sprintf(TEMPLATE, a.Name, 0, "年")
		if msg != expected {
			t.Errorf("メッセージが不正: got %q, want %q", msg, expected)
		}
	})

	t.Run("月日不一致・Period倍数一致: check=true, 日数メッセージ", func(t *testing.T) {
		// 2020-01-01 から 100日後 = 2020-04-10
		a := &table.Anniversary{
			ID:     "3",
			Date:   "2020-01-01",
			Period: "100",
			Name:   "百日記念",
		}
		now := time.Date(2020, time.April, 10, 0, 0, 0, 0, jst)
		msg, check, err := MakeCheckMessage(a, now)
		if err != nil {
			t.Fatalf("予期しないエラー: %v", err)
		}
		if !check {
			t.Errorf("check は true であるべきだが false だった")
		}
		expected := fmt.Sprintf(TEMPLATE, a.Name, 100, "日")
		if msg != expected {
			t.Errorf("メッセージが不正: got %q, want %q", msg, expected)
		}
	})

	t.Run("月日不一致・Period倍数不一致: check=false", func(t *testing.T) {
		a := &table.Anniversary{
			ID:     "4",
			Date:   "2020-01-01",
			Period: "100",
			Name:   "任意記念日",
		}
		// 14日後: 14 % 100 = 14 ≠ 0
		now := time.Date(2020, time.January, 15, 0, 0, 0, 0, jst)
		_, check, err := MakeCheckMessage(a, now)
		if err != nil {
			t.Fatalf("予期しないエラー: %v", err)
		}
		if check {
			t.Errorf("check は false であるべきだが true だった")
		}
	})

	t.Run("異常系: 不正な日付フォーマット", func(t *testing.T) {
		a := &table.Anniversary{
			ID:     "5",
			Date:   "invalid-date",
			Period: "100",
			Name:   "エラーケース",
		}
		now := time.Date(2026, time.April, 21, 0, 0, 0, 0, jst)
		_, _, err := MakeCheckMessage(a, now)
		if err == nil {
			t.Error("エラーが返るべきだが nil だった")
		}
	})

	t.Run("異常系: 不正なPeriod値", func(t *testing.T) {
		a := &table.Anniversary{
			ID:     "6",
			Date:   "2020-01-01",
			Period: "not-a-number",
			Name:   "エラーケース2",
		}
		now := time.Date(2026, time.April, 21, 0, 0, 0, 0, jst)
		_, _, err := MakeCheckMessage(a, now)
		if err == nil {
			t.Error("エラーが返るべきだが nil だった")
		}
	})
}

func TestAnniversaryManager_Run_EmptyList(t *testing.T) {
	bot := &mockBot{}
	am := &AnniversaryManager{Anniversaries: nil, Bot: bot}

	if err := am.Run(context.Background()); err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	if len(bot.sent) != 0 {
		t.Errorf("送信回数は 0 であるべきだが %d だった", len(bot.sent))
	}
}

func TestAnniversaryManager_Run_CancelledContext(t *testing.T) {
	bot := &mockBot{}
	// 今日と同じ月日の記念日 → check=true で送信対象になる
	today := time_util.LocalTime()
	am := &AnniversaryManager{
		Anniversaries: []table.Anniversary{
			{
				ID:      "1",
				Date:    today.Format("2006-01-02"),
				Period:  "100",
				Name:    "テスト",
				RoomKey: "room1",
			},
		},
		Bot: bot,
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := am.Run(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("context.Canceled が返るべきだが %v だった", err)
	}
	if len(bot.sent) != 0 {
		t.Errorf("キャンセル済み ctx では送信されないべきだが %d 件送信された", len(bot.sent))
	}
}

func TestAnniversaryManager_Run_SendsOnMatch(t *testing.T) {
	bot := &mockBot{}
	today := time_util.LocalTime()
	// 今日と同じ月日（n年前）の記念日 → check=true
	pastDate := time.Date(today.Year()-3, today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	am := &AnniversaryManager{
		Anniversaries: []table.Anniversary{
			{
				ID:      "1",
				Date:    pastDate.Format("2006-01-02"),
				Period:  "100",
				Name:    "三周年",
				RoomKey: "roomA",
			},
		},
		Bot: bot,
	}

	if err := am.Run(context.Background()); err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	if len(bot.sent) != 1 {
		t.Fatalf("送信回数は 1 であるべきだが %d だった", len(bot.sent))
	}
	if bot.sent[0].roomKey != "roomA" {
		t.Errorf("roomKey 不正: got %q", bot.sent[0].roomKey)
	}
}

func TestAnniversaryManager_Run_SkipsOnNoMatch(t *testing.T) {
	bot := &mockBot{}
	today := time_util.LocalTime()
	// 昨日の日付 + Period=9999 → check=false
	yesterday := today.AddDate(0, 0, -1)
	am := &AnniversaryManager{
		Anniversaries: []table.Anniversary{
			{
				ID:      "1",
				Date:    yesterday.Format("2006-01-02"),
				Period:  "9999",
				Name:    "昨日",
				RoomKey: "roomB",
			},
		},
		Bot: bot,
	}

	if err := am.Run(context.Background()); err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	if len(bot.sent) != 0 {
		t.Errorf("マッチしない日は送信されないべきだが %d 件送信された", len(bot.sent))
	}
}

func TestAnniversaryManager_Run_MixedMatchAndNoMatch(t *testing.T) {
	// 3件の記念日を用意:
	// 1: 今日と同じ月日 (match)
	// 2: 昨日 + Period=9999 (no match)
	// 3: 今日と同じ月日 (match)
	bot := &mockBot{}
	today := time_util.LocalTime()
	yesterday := today.AddDate(0, 0, -1)
	past := time.Date(today.Year()-5, today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	am := &AnniversaryManager{
		Anniversaries: []table.Anniversary{
			{ID: "1", Date: past.Format("2006-01-02"), Period: "100", Name: "五周年", RoomKey: "roomA"},
			{ID: "2", Date: yesterday.Format("2006-01-02"), Period: "9999", Name: "昨日", RoomKey: "roomB"},
			{ID: "3", Date: today.Format("2006-01-02"), Period: "100", Name: "本日", RoomKey: "roomC"},
		},
		Bot: bot,
	}

	if err := am.Run(context.Background()); err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	if len(bot.sent) != 2 {
		t.Fatalf("マッチ2件のみ送信されるべきだが %d 件だった", len(bot.sent))
	}
	if bot.sent[0].roomKey != "roomA" {
		t.Errorf("1件目の roomKey が不正: %q", bot.sent[0].roomKey)
	}
	if bot.sent[1].roomKey != "roomC" {
		t.Errorf("2件目の roomKey が不正: %q", bot.sent[1].roomKey)
	}
}

func TestAnniversaryManager_Run_PropagatesBotError(t *testing.T) {
	botErr := errors.New("bot failure")
	bot := &mockBot{err: botErr}
	today := time_util.LocalTime()
	am := &AnniversaryManager{
		Anniversaries: []table.Anniversary{
			{
				ID:      "1",
				Date:    today.Format("2006-01-02"),
				Period:  "100",
				Name:    "本日",
				RoomKey: "roomC",
			},
		},
		Bot: bot,
	}

	err := am.Run(context.Background())
	if !errors.Is(err, botErr) {
		t.Errorf("bot エラーが伝播するべきだが %v だった", err)
	}
}

func TestAnniversaryManager_Run_PropagatesMakeMessageError(t *testing.T) {
	bot := &mockBot{}
	am := &AnniversaryManager{
		Anniversaries: []table.Anniversary{
			{
				ID:      "1",
				Date:    "invalid-date",
				Period:  "100",
				Name:    "壊れた記念日",
				RoomKey: "roomD",
			},
		},
		Bot: bot,
	}

	err := am.Run(context.Background())
	if err == nil {
		t.Error("MakeCheckMessage エラーが伝播するべきだが nil だった")
	}
	if len(bot.sent) != 0 {
		t.Errorf("エラー時は送信されないべきだが %d 件送信された", len(bot.sent))
	}
}
