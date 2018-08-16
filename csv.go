package nyanbot

import (
	"io/ioutil"

	"github.com/jszwec/csvutil"
)

type PushMessage struct {
	ID         int    `csv:"id"`
	Minute     string `csv:"minute"`
	Hour       string `csv:"hour"`
	DayOfMonth string `csv:"day_of_month"`
	Month      string `csv:"month"`
	DayOfWeek  string `csv:"day_of_week"`
	WeekNum    string `csv:"week_num"`
	Message    string `csv:"message"`
}

var samplePushMessageCSV = projectRoot + "csv/push_message_sample.csv"

func LoadPushMessages() ([]PushMessage, error) {
	var pmsgs []PushMessage

	config, err := LoadConfig()
	if err != nil {
		return []PushMessage{}, err
	}

	path := samplePushMessageCSV
	if config.CsvFileDir != "" {
		path = config.CsvFileDir + "push_message.csv"
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return []PushMessage{}, err
	}

	err = csvutil.Unmarshal(b, &pmsgs)
	if err != nil {
		return []PushMessage{}, err
	}

	return pmsgs, nil
}
