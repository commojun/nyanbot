package constant

import "os"

const (
	BaseMinuteDuration = 1
)

var (
	ChannelSecret      = os.Getenv("NYAN_CHANNEL_SECRET")
	ChannelAccessToken = os.Getenv("NYAN_ACCESS_TOKEN")
	DefaultRoomID      = os.Getenv("NYAN_DEFAULT_ROOM_ID")

	GoogleClientEmail  = os.Getenv("NYAN_GOOGLE_CLIENT_EMAIL")
	GooglePrivateKey   = os.Getenv("NYAN_GOOGLE_PRIVATE_KEY")
	GooglePrivateKeyID = os.Getenv("NYAN_GOOGLE_PRIVATE_KEY_ID")
	GoogleTokenURL     = os.Getenv("NYAN_GOOGLE_TOKEN_URL")

	SheetID = os.Getenv("NYAN_SHEET_ID")

	RedisHost = os.Getenv("NYAN_REDIS_HOST")
)
