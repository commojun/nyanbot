package key_value

import (
	"context"
	"fmt"
	"reflect"

	"github.com/commojun/nyanbot/app/sheet"
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

func LoadKVsFromSheet(ctx context.Context, s *sheet.Sheet, sheetID string) (*KVs, error) {
	kvs := &KVs{}

	kvsType := reflect.TypeOf(*kvs)
	for i := 0; i < kvsType.NumField(); i++ {
		// kvを作成
		kvName := kvsType.Field(i).Tag.Get("kvName")
		kv, err := getKVFromSheet(ctx, s, kvName, sheetID)
		if err != nil {
			return nil, err
		}
		// 生成したkvをkvsにセットする
		reflect.ValueOf(kvs).Elem().Field(i).Set(reflect.ValueOf(*kv))
	}

	return kvs, nil
}

func getKVFromSheet(ctx context.Context, s *sheet.Sheet, kvName string, sheetID string) (*map[string]string, error) {
	// シートからデータを取得
	res, err := s.Get(ctx, sheetID, kvName)
	if err != nil {
		return &map[string]string{}, err
	}

	// シートのヘッダ情報
	header := res.Values[0]
	headerMap := make(map[string]int, len(header))
	for i, cell := range header {
		headerMap[cell.(string)] = i
	}

	data := res.Values[1:]
	kv := map[string]string{}
	for _, row := range data {
		// keyが空の行はスキップする
		if row[headerMap["key"]] == "" {
			continue
		}

		kv[row[headerMap["key"]].(string)] = row[headerMap["value"]].(string)
	}

	return &kv, err
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
