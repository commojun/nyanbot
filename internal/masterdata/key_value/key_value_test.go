package key_value

import (
	"context"
	"errors"
	"testing"
)

func TestRoomID(t *testing.T) {
	kvs := &KVs{
		Rooms: map[string]string{
			"main": "room_id_001",
			"sub":  "room_id_002",
		},
	}

	t.Run("存在するキーでルームIDを取得", func(t *testing.T) {
		got, err := kvs.RoomID("main")
		if err != nil {
			t.Fatalf("予期しないエラー: %v", err)
		}
		if got != "room_id_001" {
			t.Errorf("RoomID: got %q, want %q", got, "room_id_001")
		}
	})

	t.Run("別の存在するキーでルームIDを取得", func(t *testing.T) {
		got, err := kvs.RoomID("sub")
		if err != nil {
			t.Fatalf("予期しないエラー: %v", err)
		}
		if got != "room_id_002" {
			t.Errorf("RoomID: got %q, want %q", got, "room_id_002")
		}
	})

	t.Run("存在しないキーはエラーを返す", func(t *testing.T) {
		_, err := kvs.RoomID("not_exist")
		if err == nil {
			t.Error("エラーが返るべきだが nil だった")
		}
	})

	t.Run("空文字列キーはエラーを返す", func(t *testing.T) {
		_, err := kvs.RoomID("")
		if err == nil {
			t.Error("エラーが返るべきだが nil だった")
		}
	})
}

func TestNickname(t *testing.T) {
	kvs := &KVs{
		Nicknames: map[string]string{
			"user_001": "にゃんこ",
			"user_002": "たろう",
		},
	}

	t.Run("存在するユーザーIDでニックネームを取得", func(t *testing.T) {
		got, err := kvs.Nickname("user_001")
		if err != nil {
			t.Fatalf("予期しないエラー: %v", err)
		}
		if got != "にゃんこ" {
			t.Errorf("Nickname: got %q, want %q", got, "にゃんこ")
		}
	})

	t.Run("別の存在するユーザーIDでニックネームを取得", func(t *testing.T) {
		got, err := kvs.Nickname("user_002")
		if err != nil {
			t.Fatalf("予期しないエラー: %v", err)
		}
		if got != "たろう" {
			t.Errorf("Nickname: got %q, want %q", got, "たろう")
		}
	})

	t.Run("存在しないユーザーIDはエラーを返す", func(t *testing.T) {
		_, err := kvs.Nickname("not_exist")
		if err == nil {
			t.Error("エラーが返るべきだが nil だった")
		}
	})

	t.Run("空文字列ユーザーIDはエラーを返す", func(t *testing.T) {
		_, err := kvs.Nickname("")
		if err == nil {
			t.Error("エラーが返るべきだが nil だった")
		}
	})
}

// --- parseKVRows テスト ---

func TestParseKVRows_Normal(t *testing.T) {
	rows := [][]any{
		{"key", "value"},
		{"main", "room_id_001"},
		{"sub", "room_id_002"},
	}
	kv, err := parseKVRows(rows)
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	if len(kv) != 2 {
		t.Fatalf("len=%d, want 2", len(kv))
	}
	if kv["main"] != "room_id_001" || kv["sub"] != "room_id_002" {
		t.Errorf("値が不正: %+v", kv)
	}
}

func TestParseKVRows_SkipsEmptyKey(t *testing.T) {
	rows := [][]any{
		{"key", "value"},
		{"foo", "1"},
		{"", "orphan"}, // 空キー行
		{"bar", "2"},
	}
	kv, err := parseKVRows(rows)
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	if len(kv) != 2 {
		t.Fatalf("空キーはスキップされるべきで len=2 だが %d だった", len(kv))
	}
	if kv["foo"] != "1" || kv["bar"] != "2" {
		t.Errorf("スキップ後の内容が不正: %+v", kv)
	}
}

func TestParseKVRows_ShuffledColumnOrder(t *testing.T) {
	rows := [][]any{
		{"value", "key"}, // 列順逆
		{"VAL1", "K1"},
	}
	kv, err := parseKVRows(rows)
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	if kv["K1"] != "VAL1" {
		t.Errorf("列順入替でマッピング失敗: %+v", kv)
	}
}

func TestParseKVRows_HeaderOnly_EmptyMap(t *testing.T) {
	rows := [][]any{
		{"key", "value"},
	}
	kv, err := parseKVRows(rows)
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	if len(kv) != 0 {
		t.Errorf("空マップであるべきだが len=%d", len(kv))
	}
}

func TestParseKVRows_EmptyRows_Error(t *testing.T) {
	_, err := parseKVRows([][]any{})
	if err == nil {
		t.Error("空入力でエラーが返るべきだが nil だった")
	}
}

func TestParseKVRows_MissingKeyColumn_Error(t *testing.T) {
	rows := [][]any{
		{"value"}, // key 欠落
		{"val1"},
	}
	_, err := parseKVRows(rows)
	if err == nil {
		t.Error("key列欠落でエラーが返るべきだが nil だった")
	}
}

func TestParseKVRows_MissingValueColumn_Error(t *testing.T) {
	rows := [][]any{
		{"key"}, // value 欠落
		{"k1"},
	}
	_, err := parseKVRows(rows)
	if err == nil {
		t.Error("value列欠落でエラーが返るべきだが nil だった")
	}
}

// --- LoadKVsFromSheet テスト（SheetFetcher mock） ---

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

func TestLoadKVsFromSheet_Success(t *testing.T) {
	fetcher := &mockFetcher{
		data: map[string][][]any{
			"room": {
				{"key", "value"},
				{"main", "R-001"},
				{"sub", "R-002"},
			},
			"nickname": {
				{"key", "value"},
				{"user_a", "にゃんこ"},
			},
			"testkv": {
				{"key", "value"},
				{"foo", "bar"},
			},
		},
	}

	kvs, err := LoadKVsFromSheet(context.Background(), fetcher, "dummy-sheet-id")
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	if len(kvs.Rooms) != 2 {
		t.Errorf("Rooms len=%d, want 2", len(kvs.Rooms))
	}
	if len(kvs.Nicknames) != 1 {
		t.Errorf("Nicknames len=%d, want 1", len(kvs.Nicknames))
	}
	if len(kvs.Tests) != 1 {
		t.Errorf("Tests len=%d, want 1", len(kvs.Tests))
	}
	if kvs.Rooms["main"] != "R-001" {
		t.Errorf("Rooms['main'] = %q, want %q", kvs.Rooms["main"], "R-001")
	}
	if kvs.Nicknames["user_a"] != "にゃんこ" {
		t.Errorf("Nicknames['user_a'] = %q", kvs.Nicknames["user_a"])
	}
}

func TestLoadKVsFromSheet_FetcherError(t *testing.T) {
	fetcher := &mockFetcher{err: errors.New("sheet error")}

	_, err := LoadKVsFromSheet(context.Background(), fetcher, "dummy-sheet-id")
	if err == nil {
		t.Error("fetcher エラーが伝播するべきだが nil だった")
	}
}

func TestLoadKVsFromSheet_ParseError(t *testing.T) {
	// room シートが空 → parseKVRows でエラー
	fetcher := &mockFetcher{
		data: map[string][][]any{
			"room":     {},
			"nickname": {},
			"testkv":   {},
		},
	}

	_, err := LoadKVsFromSheet(context.Background(), fetcher, "dummy-sheet-id")
	if err == nil {
		t.Error("パースエラーが伝播するべきだが nil だった")
	}
}
