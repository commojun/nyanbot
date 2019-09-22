package table

type Anniversary struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Period  string `json:"period"`
	Name    string `json:"name"`
	RoomKey string `json:"room_key"`
}

func Anniversaries() ([]Anniversary, error) {
	anniversaries := []Anniversary{}
	err := RestoreFromRedis(&anniversaries, "anniversary")
	return anniversaries, err
}
