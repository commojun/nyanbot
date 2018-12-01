package alarm

import (
	"io/ioutil"
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

func LoadPushMessageCase() ([]PushMessageTestCase, error) {
	b, err := ioutil.ReadFile(projectRoot + "csv/push_message_test.csv")
	if err != nil {
		return []PushMessageTestCase{}, err
	}

	var cases []PushMessageTestCase
	err = csvutil.Unmarshal(b, &cases)
	if err != nil {
		return []PushMessageTestCase{}, err
	}

	//stringから各構造へ
	for i, c := range cases {
		cases[i].Time, err = JSTParse(c.TimeStr)
		if err != nil {
			return []PushMessageTestCase{}, err
		}
		cases[i].Expected, err = strconv.ParseBool(c.ExpectedStr)
		if err != nil {
			return []PushMessageTestCase{}, err
		}
	}

	return cases, nil
}

func TestCanSendPushMessage(t *testing.T) {
	cases, err := LoadPushMessageCase()
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range cases {
		pmsg := c.PushMessage
		result, err := CanSendPushMessage(pmsg, c.Time)
		if err != nil {
			t.Fatal(err)
		}
		if result != c.Expected {
			t.Fatalf("result is %t, but expected is %t for id %d", result, c.Expected, c.ID)
		}
	}
}
