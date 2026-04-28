package table

import (
	"context"
	"fmt"
	"reflect"
)

type Tables struct {
	Alarms        []Alarm       `tableName:"alarm"`
	Anniversaries []Anniversary `tableName:"anniversary"`
}

// SheetFetcher: スプレッドシートから指定シート名の行データを取得するための抽象
type SheetFetcher interface {
	Get(ctx context.Context, spreadsheetID string, sheetName string) ([][]any, error)
}

func LoadTablesFromSheet(ctx context.Context, s SheetFetcher, sheetID string) (*Tables, error) {
	ts := &Tables{}

	tsType := reflect.TypeOf(*ts)
	for i := 0; i < tsType.NumField(); i++ {
		tName := tsType.Field(i).Tag.Get("tableName")
		tType := tsType.Field(i).Type
		tValue, err := getTableFromSheet(ctx, s, tType, tName, sheetID)
		if err != nil {
			return nil, err
		}
		reflect.ValueOf(ts).Elem().Field(i).Set(tValue.Elem())
	}

	return ts, nil
}

func getTableFromSheet(ctx context.Context, s SheetFetcher, tType reflect.Type, tName string, sheetID string) (reflect.Value, error) {
	rows, err := s.Get(ctx, sheetID, tName)
	if err != nil {
		return reflect.Value{}, err
	}
	return parseTableRows(tType, rows)
}

// parseTableRows: Google Sheets の生データ ([][]any) を構造体スライスにマッピングする純粋関数
// tType は対象スライス型 (例: []Alarm, []Anniversary)。構造体フィールドは json タグでヘッダ列と対応づける
func parseTableRows(tType reflect.Type, rows [][]any) (reflect.Value, error) {
	if len(rows) == 0 {
		return reflect.New(tType), fmt.Errorf("no rows: header missing")
	}

	header := rows[0]
	headerMap := make(map[string]int, len(header))
	for i, cell := range header {
		name, ok := cell.(string)
		if !ok {
			return reflect.Value{}, fmt.Errorf("header cell %d is not a string", i)
		}
		headerMap[name] = i
	}

	idIdx, ok := headerMap["id"]
	if !ok {
		return reflect.Value{}, fmt.Errorf("id column is missing in header")
	}

	tValue := reflect.New(tType)
	objType := tType.Elem()
	for rowIdx, row := range rows[1:] {
		// id列が空の行はスキップ
		if idIdx >= len(row) || row[idIdx] == "" {
			continue
		}

		objValue := reflect.New(objType)
		for j := 0; j < objType.NumField(); j++ {
			colName := objType.Field(j).Tag.Get("json")
			cIndex, ok := headerMap[colName]
			if !ok {
				return reflect.Value{}, fmt.Errorf("column %q (field %q) not found in header", colName, objType.Field(j).Name)
			}
			if cIndex >= len(row) {
				// 行末が足りない場合は空文字扱い
				continue
			}
			value, ok := row[cIndex].(string)
			if !ok {
				return reflect.Value{}, fmt.Errorf("row %d column %q is not a string", rowIdx+1, colName)
			}
			objValue.Elem().Field(j).SetString(value)
		}
		tValue.Elem().Set(reflect.Append(tValue.Elem(), objValue.Elem()))
	}

	return tValue, nil
}
