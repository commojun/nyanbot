package table

type Anniversary struct {
	ID      string `json:"id"`
	Year    string `json:"year"`
	Month   string `json:"month"`
	Day     string `json:"day"`
	Name    string `json:"name"`
	Period  string `json:"period"`
	RoomKey string `json:"room_key"`
}

func Anniversaries() (*[]Anniversary, error) {
	anniversaries := &[]Anniversary{}
	err := RestoreFromRedis(anniversaries, "anniversary")
	return anniversaries, err
}
