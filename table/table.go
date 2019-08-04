package table

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/commojun/nyanbot/app/constant"
	"github.com/commojun/nyanbot/app/sheet"
)

type Table struct {
	Alarms []alarm `json:"alarm"`
}

type alarm struct {
	ID         string `json:"id"`
	Minute     string `json:"minute"`
	Hour       string `json:"hour"`
	DayOfMonth string `json:"day_of_month"`
	Month      string `json:"month"`
	DayOfWeek  string `json:"day_of_week"`
	WeekNum    string `json:"week_num"`
	Message    string `json:"message"`
	RoomID     string `json:"room_id"`
}

func New() (*Table, error) {
	table := Table{}
	return &table, nil
}

func (table *Table) LoadFromSheet() error {
	s, err := sheet.New()
	if err != nil {
		return err
	}

	res, err := s.Get(constant.SheetID, "alarm")

	header := res.Values[0]
	data := res.Values[1:]

	headerMap := make(map[string]int, len(header))
	for i, cell := range header {
		headerMap[cell.(string)] = i
	}

	alms := []alarm{}
	for _, row := range data {
		alm := alarm{}
		almType := reflect.TypeOf(alm)
		for i := 0; i < almType.NumField(); i++ {
			cIndex := headerMap[almType.Field(i).Tag.Get("json")]
			value := row[cIndex]
			reflect.ValueOf(&alm).Elem().Field(i).SetString(value.(string))
		}
		alms = append(alms, alm)
		table.Alarms = alms
	}

	jsonBytes, err := json.Marshal(table.Alarms)
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return err
	}

	fmt.Println(string(jsonBytes))

	return err
}
