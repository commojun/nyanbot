package nyanbot

import (
	"io/ioutil"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/jszwec/csvutil"
)

type PushMessageTestCase struct {
	PushMessage
	TimeStr     string `csv:"time"`
	ExpectedStr string `csv:"expected"`
	Time        time.Time
	Expected    bool
}

func LoadPushMessageCase() []PushMessageTestCase {
	config := LoadConfig()

	b, err := ioutil.ReadFile(config.CsvFileDir + "push_message_test.csv")
	if err != nil {
		log.Fatal(err)
	}

	var cases []PushMessageTestCase
	err = csvutil.Unmarshal(b, &cases)
	if err != nil {
		log.Fatal(err)
	}

	//stringから各構造へ
	for i, c := range cases {
		cases[i].Time, err = JSTParse(c.TimeStr)
		if err != nil {
			log.Fatal(err)
		}
		cases[i].Expected, err = strconv.ParseBool(c.ExpectedStr)
		if err != nil {
			log.Fatal(err)
		}
	}

	return cases
}

func TestCanSendPushMessage(t *testing.T) {
	cases := LoadPushMessageCase()

	for _, c := range cases {
		pmsg := c.PushMessage
		result := CanSendPushMessage(pmsg, c.Time)
		if result != c.Expected {
			t.Fatalf("result is %t, but expected is %t for id %d", result, c.Expected, c.ID)
		}
	}
}
