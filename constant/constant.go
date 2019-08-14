package constant

import (
	"os"
	"strings"
)

const (
	BaseMinuteDuration = 1
)

var (
	ServerPort = os.Getenv("NYAN_SERVER_PORT")

	ChannelSecret      = os.Getenv("NYAN_CHANNEL_SECRET")
	ChannelAccessToken = os.Getenv("NYAN_ACCESS_TOKEN")
	DefaultRoomID      = os.Getenv("NYAN_DEFAULT_ROOM_ID")

	GoogleClientEmail  = os.Getenv("NYAN_GOOGLE_CLIENT_EMAIL")
	GooglePrivateKey   = strings.Replace(os.Getenv("NYAN_GOOGLE_PRIVATE_KEY"), `\n`, "\n", -1)
	GooglePrivateKeyID = os.Getenv("NYAN_GOOGLE_PRIVATE_KEY_ID")
	GoogleTokenURL     = os.Getenv("NYAN_GOOGLE_TOKEN_URL")

	SheetID = os.Getenv("NYAN_SHEET_ID")

	RedisHost = os.Getenv("NYAN_REDIS_HOST")
)
