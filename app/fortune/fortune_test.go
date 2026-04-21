package fortune

import (
	"testing"
)

func TestDrawByStringSeed(t *testing.T) {
	f := New()

	t.Run("同一seed文字列は同一結果を返す（冪等性）", func(t *testing.T) {
		r1 := f.DrawByStringSeed("にゃんこ")
		r2 := f.DrawByStringSeed("にゃんこ")
		if r1 != r2 {
			t.Errorf("同一seedで異なる結果: %q != %q", r1, r2)
		}
	})

	t.Run("異なるseed文字列でも有効な運勢が返る", func(t *testing.T) {
		seeds := []string{"alice", "bob", "charlie"}
		for _, seed := range seeds {
			got := f.DrawByStringSeed(seed)
			if !isValidResult(got) {
				t.Errorf("seed %q の結果が無効: %q", seed, got)
			}
		}
	})

	t.Run("空文字列でもパニックせず有効な運勢が返る", func(t *testing.T) {
		got := f.DrawByStringSeed("")
		if !isValidResult(got) {
			t.Errorf("空文字列の結果が無効: %q", got)
		}
	})

	t.Run("同一実行内で複数のseed文字列がすべて有効な運勢を返す", func(t *testing.T) {
		seeds := []string{"大吉希望", "test123", "日本語seed", "!@#$"}
		for _, seed := range seeds {
			got := f.DrawByStringSeed(seed)
			if !isValidResult(got) {
				t.Errorf("seed %q の結果が無効: %q", seed, got)
			}
		}
	})
}

// isValidResult は結果が有効な運勢（results スライスに含まれる）かどうかを確認する
func isValidResult(s string) bool {
	for _, r := range results {
		if s == r {
			return true
		}
	}
	return false
}
