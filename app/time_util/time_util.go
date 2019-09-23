package time_util

import (
	"fmt"
	"time"
)

var nyantime = struct {
	Layout         string
	LocationName   string
	LocationOffset int
}{
	Layout:         "2006-01-02 15:04:05",
	LocationName:   "JST",
	LocationOffset: 9 * 60 * 60,
}

func JSTParse(s string) (time.Time, error) {
	l := time.FixedZone(nyantime.LocationName, nyantime.LocationOffset)

	t, err := time.ParseInLocation(nyantime.Layout, s, l)
	if err != nil {
		return time.Time{}, err
	}
	fmt.Println(t)
	return t, nil
}

func LocalTime() time.Time {
	return time.Now().In(time.Local)
}
