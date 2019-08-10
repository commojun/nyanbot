package table

type Alarm struct {
	ID         string `json:"id"`
	Minute     string `json:"minute"`
	Hour       string `json:"hour"`
	DayOfMonth string `json:"day_of_month"`
	Month      string `json:"month"`
	DayOfWeek  string `json:"day_of_week"`
	WeekNum    string `json:"week_num"`
	Message    string `json:"message"`
	RoomKey    string `json:"room_key"`
}

func Alarms() (*[]Alarm, error) {
	alarms := &[]Alarm{}
	err := RestoreFromRedis(alarms, "alarm")
	return alarms, err
}
