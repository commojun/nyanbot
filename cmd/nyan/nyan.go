package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/commojun/nyanbot"
	"github.com/commojun/nyanbot/app/alarm"
	"github.com/commojun/nyanbot/app/anniversary"
	"github.com/commojun/nyanbot/app/hello"
	"github.com/commojun/nyanbot/cache"
)

func main() {
	flag.Usage = flagUsage

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	switch os.Args[1] {
	case "server":
		Server()
	case "hello":
		Hello()
	case "alarm":
		Alarm()
	case "anniversary":
		Anniversary()
	default:
		flagUsage()
	}
}

func flagUsage() {
	usageText := `nyanbot

Usage:
nyan command [args]

server
hello
alarm
anniversary`

	fmt.Fprintf(os.Stderr, "%s\n\n", usageText)
}

func Server() {
	if err := cache.Initialize(); err != nil {
		log.Fatal(err)
	}

	server, err := nyanbot.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func Hello() {
	if err := cache.Initialize(); err != nil {
		log.Fatal(err)
	}

	hello, err := hello.New()
	if err != nil {
		log.Fatal(err)
	}

	err = hello.Say()
	if err != nil {
		log.Fatal(err)
	}
}

func Alarm() {
	if err := cache.Initialize(); err != nil {
		log.Fatal(err)
	}

	alm, err := alarm.New()
	if err != nil {
		log.Fatal(err)
	}

	err = alm.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func Anniversary() {
	if err := cache.Initialize(); err != nil {
		log.Fatal(err)
	}

	anniv, err := anniversary.New()
	if err != nil {
		log.Fatal(err)
	}

	err = anniv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
