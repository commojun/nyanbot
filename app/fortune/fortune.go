package fortune

import (
	"fmt"
	"math/rand"
	"time"
)

type Fortune struct {
	Results []string
}

var (
	results = []string{
		"大吉",
		"中吉",
		"小吉",
		"吉",
		"末吉",
		"凶",
		"大凶",
	}
)

func New() *Fortune {
	return &Fortune{
		Results: results,
	}
}

func (f *Fortune) DrawByStringSeed(seed string) string {
	//stringをbyteに変換してseedとする
	byte := []byte(seed)
	var num int
	for _, b := range byte {
		num += int(b)
	}
	return f.Draw(num)
}

func (f *Fortune) Draw(seed int) string {
	//日付と引数に依存したseedを作成
	t := time.Now()
	day := t.Truncate(time.Hour).Add(-time.Duration(t.Hour()) * time.Hour)
	fmt.Println(day)
	rand.Seed(day.Unix() + int64(seed))
	i := rand.Intn(len(f.Results))

	return f.Results[i]
}
