package alarm

import (
	"fmt"
	"strconv"
	"time"

	"github.com/commojun/nyanbot/app/constant"
)

type Entity struct {
	ID         int    `csv:"id"`
	Minute     string `csv:"minute"`
	Hour       string `csv:"hour"`
	DayOfMonth string `csv:"day_of_month"`
	Month      string `csv:"month"`
	DayOfWeek  string `csv:"day_of_week"`
	WeekNum    string `csv:"week_num"`
	Message    string `csv:"message"`
}

func (entity *Entity) Check(now time.Time) (bool, error) {
	weekNum := (int(now.Day()) / 7) + 1

	figs := []string{entity.Month, entity.WeekNum, entity.DayOfWeek, entity.DayOfMonth, entity.Hour}
	comparison := []int{int(now.Month()), weekNum, int(now.Weekday()), now.Day(), now.Hour()}

	for i, fig := range figs {
		if fig == "*" {
			continue
		}
		fignum, err := strconv.Atoi(fig)
		if err != nil {
			return false, err
		}
		if fignum != comparison[i] {
			fmt.Printf("false because %d != %d \n", fignum, comparison[i])
			return false, nil
		}
	}

	if entity.Minute != "*" {
		fignum, err := strconv.Atoi(entity.Minute)
		if err != nil {
			return false, err
		}

		d := constant.BaseMinuteDuration
		m1 := now.Minute()
		m2 := now.Add(time.Duration(-d) * time.Minute).Minute()
		// 過去d分間を見て当てはまらなかったらfalse
		if (fignum <= (m1-d) || m1 < fignum) &&
			(fignum <= m2 || (m2+d) < fignum) {
			fmt.Printf("false because not match ( %d < %d <= %d ) \n", m2, fignum, m1)
			return false, nil
		}
	}

	return true, nil
}
