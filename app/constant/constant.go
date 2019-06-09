package constant

import "os"

const (
	BaseMinuteDuration = 1
)

var (
	UserDir = os.Getenv("NYAN_USER_DIR")

	ChannelSecret      = os.Getenv("NYAN_CHANNEL_SECRET")
	ChannelAccessToken = os.Getenv("NYAN_ACCESS_TOKEN")
	RoomId             = os.Getenv("NYAN_ROOM_ID")

	//csv
	CsvDir       = UserDir + "/csv"
	AlarmCsvPath = CsvDir + "/push_message.csv"
)
