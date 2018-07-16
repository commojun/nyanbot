package nyanbot

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os/user"
)

var configFile = "/.config/nyanbot/config.yml"

type Config struct {
	ChannelSecret      string `yaml:"channel_secret"`
	ChannelAccessToken string `yaml:"channel_access_token"`
	RoomId             string `yaml:"room_id"`
}

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
