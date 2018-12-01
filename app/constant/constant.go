package constant

import "os"

const (
	BaseMinuteDuration = 5
)

var (
	ConfigDir = os.Getenv("NYAN_CONF_DIR")

	//config
	ConfigPath = ConfigDir + "/config.yml"

	//csv
	CsvDir       = ConfigDir + "/csv"
	AlarmCsvPath = CsvDir + "/push_message.csv"
)
