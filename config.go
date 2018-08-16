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

var projectRoot = os.Getenv("GOPATH") + "/src/github.com/junpooooow/nyanbot/"

// デフォルトではテスト用の設定ファイルを参照
var ConfigFile = projectRoot + "config/config_test.yml"

func LoadConfig() (Config, error) {
	buf, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
