package key_value

import (
	"testing"
)

func TestRoomID(t *testing.T) {
	kvs := &KVs{
		Rooms: map[string]string{
			"main":  "room_id_001",
			"sub":   "room_id_002",
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
