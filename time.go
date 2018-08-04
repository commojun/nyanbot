package nyanbot

import "time"

var nyantime = struct {
	Layout      string
	LocationStr string
}{
	Layout:      "2006-01-02 15:04:05",
	LocationStr: "Asia/Tokyo",
}

func JSTParse(s string) (time.Time, error) {
	l, err := time.LoadLocation(nyantime.LocationStr)
	if err != nil {
		return time.Time{}, err
	}

	t, err := time.ParseInLocation(nyantime.Layout, s, l)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
