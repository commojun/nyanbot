package key_value

import (
	"reflect"

	"github.com/commojun/nyanbot/app/constant"
	"github.com/commojun/nyanbot/app/sheet"
)

type KVs struct {
	Rooms  map[string]string `kvName:"room"`
	TestKV map[string]string `kvName:"testkv"`
}

func New() (*KVs, error) {
	kvs := KVs{}
	return &kvs, nil
}

func (kvs *KVs) LoadKVsFromSheet() error {
	// sheetServiceは使い回す
	s, err := sheet.New()
	if err != nil {
		return err
	}

	kvsType := reflect.TypeOf(*kvs)
	for i := 0; i < kvsType.NumField(); i++ {
		// 各kvを作成
		kvName := kvsType.Field(i).Tag.Get("kvName")
		kv, err := getKVFromSheet(s, kvName)
		if err != nil {
			return err
		}
		// 生成したkvをkvsにセットする
		reflect.ValueOf(kvs).Elem().Field(i).Set(reflect.ValueOf(*kv))
	}

	return err
}

func getKVFromSheet(s *sheet.Sheet, kvName string) (*map[string]string, error) {
	// シートからデータを取得
	res, err := s.Get(constant.SheetID, kvName)
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
