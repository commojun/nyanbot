package constant

import "os"

const (
	BaseMinuteDuration = 1
)

var (
	UserDir = os.Getenv("NYAN_USER_DIR")

	//config
	ConfigDir  = UserDir + "/config"
	ConfigPath = ConfigDir + "/config.yml"

	//csv
	CsvDir       = UserDir + "/csv"
	AlarmCsvPath = CsvDir + "/push_message.csv"
)
