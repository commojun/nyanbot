package main

import (
	"log"
	"os"

	"github.com/commojun/nyanbot/app/echo"
	"github.com/commojun/nyanbot/app/linebot"
	"github.com/commojun/nyanbot/config"
)

func main() {
	bot, err := linebot.New(config.Config{
		ChannelSecret:      os.Getenv("NYAN_CHANNEL_SECRET"),
		ChannelAccessToken: os.Getenv("NYAN_ACCESS_TOKEN"),
		DefaultRoomID:      os.Getenv("NYAN_DEFAULT_ROOM_ID"),
	})
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New(bot)
	err = e.StartServer()
	if err != nil {
		log.Fatal(err)
	}
}
