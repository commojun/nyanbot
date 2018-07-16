package main

import (
	"github.com/line/line-bot-sdk-go/linebot"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os/user"
)

type Config struct {
	ChannelSecret      string `yaml:"channel_secret"`
	ChannelAccessToken string `yaml:"channel_access_token"`
	RoomId             string `yaml:"room_id"`
}

func main() {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	buf, err := ioutil.ReadFile(u.HomeDir + "/.config/nyanbot/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := linebot.New(config.ChannelSecret, config.ChannelAccessToken)
	if err != nil {
	}

	bot.PushMessage(config.RoomId, linebot.NewTextMessage("ニャンだよ")).Do()
	if err != nil {
	}
}
