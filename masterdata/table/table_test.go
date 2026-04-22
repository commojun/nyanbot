package table

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

// --- parseTableRows テスト ---

func TestParseTableRows_Alarm_Normal(t *testing.T) {
	rows := [][]any{
		{"id", "minute", "hour", "day_of_month", "month", "day_of_week", "week_num", "message", "room_key"},
		{"1", "30", "10", "*", "*", "*", "*", "おはよう", "roomA"},
		{"2", "0", "20", "*", "*", "*", "*", "おやすみ", "roomB"},
	}

	v, err := parseTableRows(reflect.TypeOf([]Alarm{}), rows)
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	got := v.Elem().Interface().([]Alarm)
	if len(got) != 2 {
		t.Fatalf("len=%d, want 2", len(got))
	}
	if got[0].ID != "1" || got[0].Minute != "30" || got[0].Message != "おはよう" || got[0].RoomKey != "roomA" {
		t.Errorf("1件目が不正: %+v", got[0])
	}
	if got[1].ID != "2" || got[1].Hour != "20" || got[1].Message != "おやすみ" {
		t.Errorf("2件目が不正: %+v", got[1])
	}
}

func TestParseTableRows_Anniversary_Normal(t *testing.T) {
	rows := [][]any{
		{"id", "date", "period", "name", "room_key"},
		{"1", "2020-01-01", "100", "記念日A", "roomA"},
		{"2", "2021-06-15", "365", "記念日B", "roomB"},
	}

	v, err := parseTableRows(reflect.TypeOf([]Anniversary{}), rows)
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	got := v.Elem().Interface().([]Anniversary)
	if len(got) != 2 {
		t.Fatalf("len=%d, want 2", len(got))
	}
	if got[0].Date != "2020-01-01" || got[0].Period != "100" {
		t.Errorf("1件目が不正: %+v", got[0])
	}
}

func TestParseTableRows_SkipsEmptyIDRow(t *testing.T) {
	rows := [][]any{
		{"id", "minute", "hour", "day_of_month", "month", "day_of_week", "week_num", "message", "room_key"},
		{"1", "0", "10", "*", "*", "*", "*", "msg1", "roomA"},
		{"", "0", "11", "*", "*", "*", "*", "skip", "roomX"}, // 空ID
		{"3", "0", "12", "*", "*", "*", "*", "msg3", "roomC"},
	}

	v, err := parseTableRows(reflect.TypeOf([]Alarm{}), rows)
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	got := v.Elem().Interface().([]Alarm)
	if len(got) != 2 {
		t.Fatalf("空IDはスキップされるべきで len=2 だが %d だった", len(got))
	}
	if got[0].ID != "1" || got[1].ID != "3" {
		t.Errorf("スキップ後の並びが不正: %+v", got)
	}
}

func TestParseTableRows_HeaderOnly_EmptySlice(t *testing.T) {
	rows := [][]any{
		{"id", "minute", "hour", "day_of_month", "month", "day_of_week", "week_num", "message", "room_key"},
	}

	v, err := parseTableRows(reflect.TypeOf([]Alarm{}), rows)
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	got := v.Elem().Interface().([]Alarm)
	if len(got) != 0 {
		t.Errorf("データ行なしの場合は空スライスであるべきだが len=%d だった", len(got))
	}
}

func TestParseTableRows_ShuffledColumnOrder(t *testing.T) {
	// 列順が構造体定義と異なっていても、json タグで紐づくので正しくマッピングされる
	rows := [][]any{
		{"room_key", "message", "week_num", "day_of_week", "month", "day_of_month", "hour", "minute", "id"},
		{"roomZ", "逆順テスト", "*", "*", "*", "*", "5", "45", "42"},
	}

	v, err := parseTableRows(reflect.TypeOf([]Alarm{}), rows)
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	got := v.Elem().Interface().([]Alarm)
	if len(got) != 1 {
		t.Fatalf("len=%d, want 1", len(got))
	}
	if got[0].ID != "42" || got[0].Minute != "45" || got[0].Hour != "5" || got[0].RoomKey != "roomZ" || got[0].Message != "逆順テスト" {
		t.Errorf("列順入替でマッピング失敗: %+v", got[0])
	}
}

func TestParseTableRows_EmptyRows_Error(t *testing.T) {
	_, err := parseTableRows(reflect.TypeOf([]Alarm{}), [][]any{})
	if err == nil {
		t.Error("空入力ではエラーが返るべきだが nil だった")
	}
}

func TestParseTableRows_MissingIDColumn_Error(t *testing.T) {
	rows := [][]any{
		{"minute", "hour", "day_of_month", "month", "day_of_week", "week_num", "message", "room_key"}, // id欠落
		{"30", "10", "*", "*", "*", "*", "msg", "room"},
	}
	_, err := parseTableRows(reflect.TypeOf([]Alarm{}), rows)
	if err == nil {
		t.Error("id列欠落でエラーが返るべきだが nil だった")
	}
}

func TestParseTableRows_MissingFieldColumn_Error(t *testing.T) {
	// message 列がヘッダに存在しない → Alarm 構造体の Message フィールドに対応する列がない
	rows := [][]any{
		{"id", "minute", "hour", "day_of_month", "month", "day_of_week", "week_num", "room_key"},
		{"1", "30", "10", "*", "*", "*", "*", "room"},
	}
	_, err := parseTableRows(reflect.TypeOf([]Alarm{}), rows)
	if err == nil {
		t.Error("構造体フィールドに対応する列欠落でエラーが返るべきだが nil だった")
	}
}

func TestParseTableRows_NonStringHeader_Error(t *testing.T) {
	rows := [][]any{
		{"id", "minute", 123, "day_of_month", "month", "day_of_week", "week_num", "message", "room_key"},
		{"1", "30", "10", "*", "*", "*", "*", "msg", "room"},
	}
	_, err := parseTableRows(reflect.TypeOf([]Alarm{}), rows)
	if err == nil {
		t.Error("非stringヘッダでエラーが返るべきだが nil だった")
	}
}

// --- LoadTablesFromSheet テスト（SheetFetcher mock） ---

type mockFetcher struct {
	data map[string][][]any
	err  error
}

func (m *mockFetcher) Get(ctx context.Context, sheetID string, sheetName string) ([][]any, error) {
	if m.err != nil {
		return nil, m.err
	}
	rows, ok := m.data[sheetName]
	if !ok {
		return nil, errors.New("sheet not found: " + sheetName)
	}
	return rows, nil
}

func TestLoadTablesFromSheet_Success(t *testing.T) {
	fetcher := &mockFetcher{
		data: map[string][][]any{
			"alarm": {
				{"id", "minute", "hour", "day_of_month", "month", "day_of_week", "week_num", "message", "room_key"},
				{"1", "30", "10", "*", "*", "*", "*", "alarm1", "roomA"},
				{"2", "0", "20", "*", "*", "*", "*", "alarm2", "roomB"},
			},
			"anniversary": {
				{"id", "date", "period", "name", "room_key"},
				{"1", "2020-01-01", "100", "記念日A", "roomX"},
			},
		},
	}

	ts, err := LoadTablesFromSheet(context.Background(), fetcher, "dummy-sheet-id")
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	if len(ts.Alarms) != 2 {
		t.Errorf("Alarms len=%d, want 2", len(ts.Alarms))
	}
	if len(ts.Anniversaries) != 1 {
		t.Errorf("Anniversaries len=%d, want 1", len(ts.Anniversaries))
	}
	if ts.Alarms[0].Message != "alarm1" {
		t.Errorf("Alarm message不正: %q", ts.Alarms[0].Message)
	}
	if ts.Anniversaries[0].Name != "記念日A" {
		t.Errorf("Anniversary name不正: %q", ts.Anniversaries[0].Name)
	}
}

func TestLoadTablesFromSheet_FetcherError(t *testing.T) {
	fetcher := &mockFetcher{err: errors.New("network down")}

	_, err := LoadTablesFromSheet(context.Background(), fetcher, "dummy-sheet-id")
	if err == nil {
		t.Error("fetcher エラーが伝播するべきだが nil だった")
	}
}

func TestLoadTablesFromSheet_ParseError(t *testing.T) {
	// alarm シートだけ壊れたヘッダを返す
	fetcher := &mockFetcher{
		data: map[string][][]any{
			"alarm":       {}, // 空 → parseTableRows でエラー
			"anniversary": {},
		},
	}

	_, err := LoadTablesFromSheet(context.Background(), fetcher, "dummy-sheet-id")
	if err == nil {
		t.Error("パースエラーが伝播するべきだが nil だった")
	}
}

func TestLoadTablesFromSheet_EmptyButValidSheets(t *testing.T) {
	fetcher := &mockFetcher{
		data: map[string][][]any{
			"alarm": {
				{"id", "minute", "hour", "day_of_month", "month", "day_of_week", "week_num", "message", "room_key"},
			},
			"anniversary": {
				{"id", "date", "period", "name", "room_key"},
			},
		},
	}

	ts, err := LoadTablesFromSheet(context.Background(), fetcher, "dummy-sheet-id")
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	if len(ts.Alarms) != 0 {
		t.Errorf("Alarms は空であるべきだが %d 件", len(ts.Alarms))
	}
	if len(ts.Anniversaries) != 0 {
		t.Errorf("Anniversaries は空であるべきだが %d 件", len(ts.Anniversaries))
	}
}
