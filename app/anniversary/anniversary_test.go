package anniversary

import (
	"fmt"
	"testing"
	"time"

	"github.com/commojun/nyanbot/masterdata/table"
)

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
