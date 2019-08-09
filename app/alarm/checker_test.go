package alarm

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/commojun/nyanbot/app/time_util"
	"github.com/commojun/nyanbot/masterdata/table"
)

type AlarmTestCase struct {
	table.Alarm
	BaseMinuteDuration int    `json:"base_minute_duration,string"`
	Expected           bool   `json:"expected,string"`
	TimeStr            string `json:"time"`
	Time               time.Time
}

func ReadTestCase() ([]AlarmTestCase, error) {
	b, err := ioutil.ReadFile("testdata/alarm_test_case.json")
	if err != nil {
		return []AlarmTestCase{}, err
	}

	var cases []AlarmTestCase
	err = json.Unmarshal(b, &cases)
	if err != nil {
		return []AlarmTestCase{}, err
	}

	//stringからTime型へ
	for i, c := range cases {
		cases[i].Time, err = time_util.JSTParse(c.TimeStr)
		if err != nil {
			return []AlarmTestCase{}, err
		}
	}

	return cases, nil
}

func TestCheck(t *testing.T) {
	cases, err := ReadTestCase()
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range cases {
		alm := c.Alarm
		result, err := Check(&alm, c.Time, c.BaseMinuteDuration)
		if err != nil {
			t.Fatal(err)
		}
		if result != c.Expected {
			t.Fatalf("result is %t, but expected is %t for id %s", result, c.Expected, c.ID)
		}
	}
}
