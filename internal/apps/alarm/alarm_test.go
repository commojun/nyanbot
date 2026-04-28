package alarm

import (
	"context"
	"errors"
	"testing"

	"github.com/commojun/nyanbot/internal/masterdata/table"
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

// wildcardAlarm: すべての時刻フィールドが "*" で常時マッチするアラーム
func wildcardAlarm(id, message, roomKey string) table.Alarm {
	return table.Alarm{
		ID:         id,
		Month:      "*",
		WeekNum:    "*",
		DayOfWeek:  "*",
		DayOfMonth: "*",
		Hour:       "*",
		Minute:     "*",
		Message:    message,
		RoomKey:    roomKey,
	}
}

func TestAlarmManager_Run_EmptyAlarms(t *testing.T) {
	bot := &mockBot{}
	am := &AlarmManager{Alarms: nil, Bot: bot}

	if err := am.Run(context.Background()); err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	if len(bot.sent) != 0 {
		t.Errorf("送信回数は 0 であるべきだが %d だった", len(bot.sent))
	}
}

func TestAlarmManager_Run_CancelledContext(t *testing.T) {
	bot := &mockBot{}
	am := &AlarmManager{
		Alarms: []table.Alarm{wildcardAlarm("1", "msg1", "room1")},
		Bot:    bot,
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

func TestAlarmManager_Run_SendsWildcardAlarm(t *testing.T) {
	bot := &mockBot{}
	am := &AlarmManager{
		Alarms: []table.Alarm{wildcardAlarm("1", "hello", "roomA")},
		Bot:    bot,
	}

	if err := am.Run(context.Background()); err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	if len(bot.sent) != 1 {
		t.Fatalf("送信回数は 1 であるべきだが %d だった", len(bot.sent))
	}
	if bot.sent[0].msg != "hello" || bot.sent[0].roomKey != "roomA" {
		t.Errorf("送信内容が不正: %+v", bot.sent[0])
	}
}

func TestAlarmManager_Run_SkipsNonMatchingAlarm(t *testing.T) {
	// Month=13 は存在しないので絶対にマッチしない
	nonMatching := table.Alarm{
		ID:         "1",
		Month:      "13",
		WeekNum:    "*",
		DayOfWeek:  "*",
		DayOfMonth: "*",
		Hour:       "*",
		Minute:     "*",
		Message:    "never",
		RoomKey:    "roomX",
	}
	bot := &mockBot{}
	am := &AlarmManager{Alarms: []table.Alarm{nonMatching}, Bot: bot}

	if err := am.Run(context.Background()); err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	if len(bot.sent) != 0 {
		t.Errorf("マッチしないアラームは送信されないべきだが %d 件送信された", len(bot.sent))
	}
}

func TestAlarmManager_Run_MixedMatchNonMatch(t *testing.T) {
	// マッチ/非マッチを混ぜた時に、マッチした分だけ送信され順序も保たれる
	nonMatching := table.Alarm{
		ID:         "2",
		Month:      "13", // 絶対にマッチしない
		WeekNum:    "*",
		DayOfWeek:  "*",
		DayOfMonth: "*",
		Hour:       "*",
		Minute:     "*",
		Message:    "never",
		RoomKey:    "roomX",
	}
	bot := &mockBot{}
	am := &AlarmManager{
		Alarms: []table.Alarm{
			wildcardAlarm("1", "first", "roomA"),
			nonMatching,
			wildcardAlarm("3", "third", "roomC"),
		},
		Bot: bot,
	}

	if err := am.Run(context.Background()); err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	if len(bot.sent) != 2 {
		t.Fatalf("マッチ2件のみ送信されるべきだが %d 件だった", len(bot.sent))
	}
	if bot.sent[0].msg != "first" || bot.sent[0].roomKey != "roomA" {
		t.Errorf("1件目が不正: %+v", bot.sent[0])
	}
	if bot.sent[1].msg != "third" || bot.sent[1].roomKey != "roomC" {
		t.Errorf("2件目が不正: %+v", bot.sent[1])
	}
}

func TestAlarmManager_Run_ContinuesOnBotError(t *testing.T) {
	// Bot がエラーを返してもループは継続する（現仕様）
	bot := &mockBot{err: errors.New("bot error")}
	am := &AlarmManager{
		Alarms: []table.Alarm{
			wildcardAlarm("1", "msg1", "room1"),
			wildcardAlarm("2", "msg2", "room2"),
		},
		Bot: bot,
	}

	if err := am.Run(context.Background()); err != nil {
		t.Fatalf("Bot エラーは握り潰されるべきだが Run が %v を返した", err)
	}
	if len(bot.sent) != 2 {
		t.Errorf("エラーでも後続の送信は実行されるべきだが %d 件だった", len(bot.sent))
	}
}
