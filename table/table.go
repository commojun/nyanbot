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
	Name   []room  `json:"room"`
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

type room struct {
	ID     string `json:"id"`
	RoomID string `json:"room_id"`
	Name   string `json:"name"`
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

	objs := make([]alarm, len(data))
	skipIds := []int{}
	for i, row := range data {

		// id列が空の行はスキップする
		if row[headerMap["id"]] == "" {
			skipIds = append(skipIds, i)
			continue
		}

		obj := objs[i]
		objType := reflect.TypeOf(obj)
		for j := 0; j < objType.NumField(); j++ {
			cIndex := headerMap[objType.Field(j).Tag.Get("json")]
			value := row[cIndex]
			reflect.ValueOf(&obj).Elem().Field(j).SetString(value.(string))
			objs[i] = obj
		}
	}

	// IDが空だった行を除外する
	for _, id := range skipIds {
		objs = append(objs[:id], objs[id+1:]...)
	}

	table.Alarms = objs

	jsonBytes, err := json.Marshal(table.Alarms)
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return err
	}

	fmt.Println(string(jsonBytes))

	return err
}
