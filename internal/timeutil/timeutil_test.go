package timeutil

import (
	"testing"
	"time"
)

func TestJSTParse(t *testing.T) {
	t.Run("正常系: 正しいフォーマットの文字列をパース", func(t *testing.T) {
		s := "2026-04-21 09:00:00"
		got, err := JSTParse(s)
		if err != nil {
			t.Fatalf("予期しないエラー: %v", err)
		}
		// タイムゾーンがJST(+09:00)であること
		_, offset := got.Zone()
		if offset != 9*60*60 {
			t.Errorf("タイムゾーンオフセットが不正: got %d, want %d", offset, 9*60*60)
		}
		name, _ := got.Zone()
		if name != "JST" {
			t.Errorf("タイムゾーン名が不正: got %s, want JST", name)
		}
		// 年月日時分秒が正しくパースされていること
		if got.Year() != 2026 || got.Month() != time.April || got.Day() != 21 {
			t.Errorf("日付が不正: got %v", got)
		}
		if got.Hour() != 9 || got.Minute() != 0 || got.Second() != 0 {
			t.Errorf("時刻が不正: got %v", got)
		}
	})

	t.Run("正常系: 別の日時をパース", func(t *testing.T) {
		s := "2018-07-12 07:30:00"
		got, err := JSTParse(s)
		if err != nil {
			t.Fatalf("予期しないエラー: %v", err)
		}
		if got.Year() != 2018 || got.Month() != time.July || got.Day() != 12 {
			t.Errorf("日付が不正: got %v", got)
		}
		if got.Hour() != 7 || got.Minute() != 30 || got.Second() != 0 {
			t.Errorf("時刻が不正: got %v", got)
		}
	})

	t.Run("異常系: 不正なフォーマット", func(t *testing.T) {
		_, err := JSTParse("not-a-date")
		if err == nil {
			t.Error("エラーが返るべきだが nil だった")
		}
	})

	t.Run("異常系: 空文字列", func(t *testing.T) {
		_, err := JSTParse("")
		if err == nil {
			t.Error("エラーが返るべきだが nil だった")
		}
	})

	t.Run("異常系: 日付部分のみ（時刻なし）", func(t *testing.T) {
		_, err := JSTParse("2026-04-21")
		if err == nil {
			t.Error("エラーが返るべきだが nil だった")
		}
	})
}
