package table

import (
	"reflect"

	"encoding/json"

	"github.com/commojun/nyanbot/app/constant"
	"github.com/commojun/nyanbot/app/redis"
	"github.com/commojun/nyanbot/app/sheet"
)

type Tables struct {
	Alarms        []Alarm       `tableName:"alarm"`
	Anniversaries []Anniversary `tableName:"anniversary"`
}

type Alarm struct {
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

type Anniversary struct {
	ID      string `json:"id"`
	Year    string `json:"year"`
	Month   string `json:"month"`
	Day     string `json:"day"`
	Name    string `json:"name"`
	Period  string `json:"period"`
	RoomKey string `json:"room_key"`
}

func New() (*Tables, error) {
	// Redisから読み込むことを試みる
	ts := Tables{}
	return &ts, nil
}

func (ts *Tables) SaveToRedis() error {
	redisClient := redis.NewClient()

	tsValue := reflect.ValueOf(*ts)
	for i := 0; i < tsValue.NumField(); i++ {
		tName := tsValue.Type().Field(i).Tag.Get("tableName")
		// tableをJSON文字列に変換する
		jsonByte, err := json.Marshal(tsValue.Field(i).Interface())
		if err != nil {
			return err
		}
		// Redisに書き込む
		err = redisClient.Set(tName, string(jsonByte), 0).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func RestoreFromRedis() {}

func (ts *Tables) LoadTablesFromSheet() error {
	// sheetServiceは使い回す
	s, err := sheet.New()
	if err != nil {
		return err
	}

	tsType := reflect.TypeOf(*ts)
	for i := 0; i < tsType.NumField(); i++ {
		// tableを生成
		tName := tsType.Field(i).Tag.Get("tableName")
		tType := tsType.Field(i).Type
		tValue, err := getTableFromSheet(s, tType, tName)
		if err != nil {
			return err
		}
		// 生成したtableをtablesにセットする
		reflect.ValueOf(ts).Elem().Field(i).Set(tValue.Elem())
	}

	return err
}

func getTableFromSheet(s *sheet.Sheet, tType reflect.Type, tName string) (reflect.Value, error) {
	// シートからデータを取得
	res, err := s.Get(constant.SheetID, tName)
	if err != nil {
		return reflect.Value{}, err
	}

	// テーブルのヘッダ情報
	header := res.Values[0]
	headerMap := make(map[string]int, len(header))
	for i, cell := range header {
		headerMap[cell.(string)] = i
	}

	data := res.Values[1:]
	tValue := reflect.New(tType)
	for _, row := range data {
		// id列が空の行はスキップする
		if row[headerMap["id"]] == "" {
			continue
		}

		// 行データの作成
		objType := tType.Elem()
		objValue := reflect.New(objType)
		for j := 0; j < objType.NumField(); j++ {
			cIndex := headerMap[objType.Field(j).Tag.Get("json")]
			value := row[cIndex]
			objValue.Elem().Field(j).SetString(value.(string))
		}
		// 行を追加
		tValue.Elem().Set(reflect.Append(tValue.Elem(), objValue.Elem()))
	}

	return tValue, err
}
