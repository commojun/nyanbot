package nyanbot

import (
	"io/ioutil"
	"log"

	"github.com/jszwec/csvutil"
)

type PushMessgae struct {
	ID         int    `csv:"id"`
	Minute     string `csv:"minute"`
	Hour       string `csv:"hour"`
	DayOfMonth string `csv:"day_of_month"`
	Month      string `csv:"month"`
	DayOfWeek  string `csv:"day_of_week"`
	WeekNum    string `csv:"week_num"`
	Message    string `csv:"message"`
}

func LoadPushMessages() []PushMessgae {
	var pmsgs []PushMessgae

	config := LoadConfig()

	b, err := ioutil.ReadFile(config.CsvFileDir + "push_message.csv")
	if err != nil {
		log.Fatal(err)
	}

	err = csvutil.Unmarshal(b, &pmsgs)
	if err != nil {
		log.Fatal(err)
	}

	return pmsgs
}
