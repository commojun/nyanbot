package key_value

import (
	"context"
	"fmt"
	"reflect"
)

var (
	Room     = "room"
	Nickname = "nickname"
	Test     = "testkv"
)

type KVs struct {
	Rooms     map[string]string `kvName:"room"`
	Nicknames map[string]string `kvName:"nickname"`
	Tests     map[string]string `kvName:"testkv"`
}

// SheetFetcher: スプレッドシートから指定シート名の行データを取得するための抽象
type SheetFetcher interface {
	Get(ctx context.Context, spreadsheetID string, sheetName string) ([][]any, error)
}

func LoadKVsFromSheet(ctx context.Context, s SheetFetcher, sheetID string) (*KVs, error) {
	kvs := &KVs{}

	kvsType := reflect.TypeOf(*kvs)
	for i := 0; i < kvsType.NumField(); i++ {
		kvName := kvsType.Field(i).Tag.Get("kvName")
		kv, err := getKVFromSheet(ctx, s, kvName, sheetID)
		if err != nil {
			return nil, err
		}
		reflect.ValueOf(kvs).Elem().Field(i).Set(reflect.ValueOf(*kv))
	}

	return kvs, nil
}

func getKVFromSheet(ctx context.Context, s SheetFetcher, kvName string, sheetID string) (*map[string]string, error) {
	rows, err := s.Get(ctx, sheetID, kvName)
	if err != nil {
		return &map[string]string{}, err
	}
	kv, err := parseKVRows(rows)
	if err != nil {
		return &map[string]string{}, err
	}
	return &kv, nil
}

// parseKVRows: Google Sheets の生データ ([][]any) を map[string]string にパースする純粋関数
// 1行目をヘッダとして "key" と "value" 列の位置を検出し、以降の行を取り込む。key が空の行はスキップ
func parseKVRows(rows [][]any) (map[string]string, error) {
	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows: header missing")
	}

	header := rows[0]
	headerMap := make(map[string]int, len(header))
	for i, cell := range header {
		name, ok := cell.(string)
		if !ok {
			return nil, fmt.Errorf("header cell %d is not a string", i)
		}
		headerMap[name] = i
	}

	keyIdx, ok := headerMap["key"]
	if !ok {
		return nil, fmt.Errorf("key column is missing in header")
	}
	valueIdx, ok := headerMap["value"]
	if !ok {
		return nil, fmt.Errorf("value column is missing in header")
	}

	kv := map[string]string{}
	for rowIdx, row := range rows[1:] {
		if keyIdx >= len(row) || row[keyIdx] == "" {
			continue
		}
		key, ok := row[keyIdx].(string)
		if !ok {
			return nil, fmt.Errorf("row %d key is not a string", rowIdx+1)
		}
		var value string
		if valueIdx < len(row) {
			v, ok := row[valueIdx].(string)
			if !ok {
				return nil, fmt.Errorf("row %d value is not a string", rowIdx+1)
			}
			value = v
		}
		kv[key] = value
	}
	return kv, nil
}

// RoomID: ルームキーからルームIDを取得
func (kvs *KVs) RoomID(key string) (string, error) {
	id, ok := kvs.Rooms[key]
	if !ok {
		return "", fmt.Errorf("room key not found: %s", key)
	}
	return id, nil
}

// Nickname: ユーザーIDからニックネームを取得
func (kvs *KVs) Nickname(userID string) (string, error) {
	name, ok := kvs.Nicknames[userID]
	if !ok {
		return "", fmt.Errorf("nickname not found: %s", userID)
	}
	return name, nil
}
