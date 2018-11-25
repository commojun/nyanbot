package main

import (
	"log"

	"github.com/commojun/nyanbot"
	flags "github.com/jessevdk/go-flags"
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

	err = nyanbot.SendPushMessage()
	if err != nil {
		log.Fatal(err)
	}
}
