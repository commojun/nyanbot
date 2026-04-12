package key_value

import (
	"reflect"

	"github.com/commojun/nyanbot/app/sheet"
)

var (
	Room     = "room"
	Nickname = "nickname"
	Test     = "testkv"
)

type KVs struct {
	Rooms    map[string]string `kvName:"room"`
	Nickname map[string]string `kvName:"nickname"`
	Tests    map[string]string `kvName:"testkv"`
}

func LoadKVsFromSheet(s *sheet.Sheet, sheetID string) (*KVs, error) {
	kvs := &KVs{}

	kvsType := reflect.TypeOf(*kvs)
	for i := 0; i < kvsType.NumField(); i++ {
		// kvを作成
		kvName := kvsType.Field(i).Tag.Get("kvName")
		kv, err := getKVFromSheet(s, kvName, sheetID)
		if err != nil {
			return nil, err
		}
		// 生成したkvをkvsにセットする
		reflect.ValueOf(kvs).Elem().Field(i).Set(reflect.ValueOf(*kv))
	}

	return kvs, nil
}

func getKVFromSheet(s *sheet.Sheet, kvName string, sheetID string) (*map[string]string, error) {
	// シートからデータを取得
	res, err := s.Get(sheetID, kvName)
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
