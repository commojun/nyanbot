package alarm

import (
	"fmt"
	"strconv"
	"time"

	"github.com/commojun/nyanbot/constant"
	"github.com/commojun/nyanbot/masterdata/table"
)

func Check(alm *table.Alarm, now time.Time) (bool, error) {
	weekNum := (int(now.Day()) / 7) + 1

	figs := []string{alm.Month, alm.WeekNum, alm.DayOfWeek, alm.DayOfMonth, alm.Hour}
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

	if alm.Minute != "*" {
		fignum, err := strconv.Atoi(alm.Minute)
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
