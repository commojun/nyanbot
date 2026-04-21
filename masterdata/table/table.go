package table

import (
	"context"
	"reflect"

	"github.com/commojun/nyanbot/app/sheet"
)

type Tables struct {
	Alarms        []Alarm       `tableName:"alarm"`
	Anniversaries []Anniversary `tableName:"anniversary"`
}

func LoadTablesFromSheet(ctx context.Context, s *sheet.Sheet, sheetID string) (*Tables, error) {
	ts := &Tables{}

	tsType := reflect.TypeOf(*ts)
	for i := 0; i < tsType.NumField(); i++ {
		// tableを生成
		tName := tsType.Field(i).Tag.Get("tableName")
		tType := tsType.Field(i).Type
		tValue, err := getTableFromSheet(ctx, s, tType, tName, sheetID)
		if err != nil {
			return nil, err
		}
		// 生成したtableをtablesにセットする
		reflect.ValueOf(ts).Elem().Field(i).Set(tValue.Elem())
	}

	return ts, nil
}

func getTableFromSheet(ctx context.Context, s *sheet.Sheet, tType reflect.Type, tName string, sheetID string) (reflect.Value, error) {
	// シートからデータを取得
	res, err := s.Get(ctx, sheetID, tName)
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
