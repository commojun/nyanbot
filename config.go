package nyanbot

import (
	"io/ioutil"
	"log"
	"os/user"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	ChannelSecret      string `yaml:"channel_secret"`
	ChannelAccessToken string `yaml:"channel_access_token"`
	RoomId             string `yaml:"room_id"`
	CsvFileDir         string `yaml:"csv_file_dir"`
	BaseMinuteDuration int    `yaml:"base_minute_duration"`
}

var configFile = "/.config/nyanbot/config.yml"

func LoadConfig() Config {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	buf, err := ioutil.ReadFile(u.HomeDir + configFile)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
