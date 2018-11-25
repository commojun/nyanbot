package nyanbot

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	ChannelSecret      string `yaml:"channel_secret"`
	ChannelAccessToken string `yaml:"channel_access_token"`
	RoomId             string `yaml:"room_id"`
	CsvFileDir         string `yaml:"csv_file_dir"`
	BaseMinuteDuration int    `yaml:"base_minute_duration"`
}

var projectRoot = os.Getenv("GOPATH") + "/src/github.com/commojun/nyanbot/"

// デフォルトではテスト用の設定ファイルを参照
var ConfigFile = projectRoot + "config/config_test.yml"

// グローバル変数に設定を入れる
var Conf = Config{}

func LoadConfig() (Config, error) {
	buf, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return Config{}, err
	}

	var conf Config
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return Config{}, err
	}

	Conf = conf

	return conf, nil
}
