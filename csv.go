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

func LoadPushMessages() ([]PushMessage, error) {
	var pmsgs []PushMessage

	config, err := LoadConfig()
	if err != nil {
		return []PushMessage{}, err
	}

	b, err := ioutil.ReadFile(config.CsvFileDir + "push_message.csv")
	if err != nil {
		return []PushMessage{}, err
	}

	err = csvutil.Unmarshal(b, &pmsgs)
	if err != nil {
		return []PushMessage{}, err
	}

	return pmsgs, nil
}
