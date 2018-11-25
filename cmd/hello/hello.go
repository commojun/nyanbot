package main

import (
	"log"

	"github.com/commojun/nyanbot"
	flags "github.com/jessevdk/go-flags"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Options struct {
	Config string `short:"c" long:"config" description:"path to config file"`
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	if opts.Config != "" {
		nyanbot.ConfigFile = opts.Config
	}

	_, err = nyanbot.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := linebot.New(nyanbot.Conf.ChannelSecret, nyanbot.Conf.ChannelAccessToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.PushMessage(nyanbot.Conf.RoomId, linebot.NewTextMessage("Hello nyan!")).Do()
	if err != nil {
		log.Fatal(err)
	}
}
