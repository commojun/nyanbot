package alarm

import (
	"log"
	"strconv"
	"time"

	"github.com/commojun/nyanbot/masterdata/table"
)

func Check(alm *table.Alarm, now time.Time, baseMinuteDuration int) (bool, error) {
	weekNum := (int(now.Day()) / 7) + 1

	labels := []string{"month", "week_num", "day_of_week", "day_of_month", "hour"}
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
			log.Printf("[ID:%s] false because the %s %d doesn't match %d \n", alm.ID, labels[i], fignum, comparison[i])
			return false, nil
		}
	}

	if alm.Minute != "*" {
		fignum, err := strconv.Atoi(alm.Minute)
		if err != nil {
			return false, err
		}

		d := baseMinuteDuration
		m1 := now.Minute()
		m2 := now.Add(time.Duration(-d) * time.Minute).Minute()
		// 過去d分間を見て当てはまらなかったらfalse
		if (fignum <= (m1-d) || m1 < fignum) &&
			(fignum <= m2 || (m2+d) < fignum) {
			log.Printf("[ID:%s] false because the minute %d doesn't match ( %d < x <= %d ) \n", alm.ID, fignum, m2, m1)
			return false, nil
		}
	}

	log.Printf("[ID:%s] true room:%s", alm.ID, alm.RoomKey)
	return true, nil
}
