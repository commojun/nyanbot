package config

import (
	"io/ioutil"

	"github.com/commojun/nyanbot/app/constant"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	ChannelSecret      string `yaml:"channel_secret"`
	ChannelAccessToken string `yaml:"channel_access_token"`
	RoomId             string `yaml:"room_id"`
}

func Load() (*Config, error) {
	return LoadWithPath(constant.ConfigPath)
}

func LoadWithPath(path string) (*Config, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return &Config{}, err
	}

	var conf Config
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return &Config{}, err
	}

	return &conf, nil
}
