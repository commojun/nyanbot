package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Songmu/retry"
	"github.com/commojun/nyanbot"
	"github.com/commojun/nyanbot/app/alarm"
	"github.com/commojun/nyanbot/app/anniversary"
	"github.com/commojun/nyanbot/app/hello"
	"github.com/commojun/nyanbot/app/redis"
	"github.com/commojun/nyanbot/masterdata/key_value"
	"github.com/commojun/nyanbot/masterdata/table"
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
	case "export":
		Export()
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
export
alarm
anniversary`

	fmt.Fprintf(os.Stderr, "%s\n\n", usageText)
}

func Server() {
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
	hello, err := hello.New()
	if err != nil {
		log.Fatal(err)
	}

	err = hello.Say()
	if err != nil {
		log.Fatal(err)
	}
}

func Export() {
	// Redisに接続できるか
	rc := redis.NewClient()
	err := retry.Retry(10, 10*time.Second, func() error {
		log.Println("attempt to connect to redis")
		err := rc.Keys("*").Err()
		return err
	})
	if err != nil {
		log.Println("failed to connect to redis")
		log.Fatal(err)
	}
	log.Println("redis connection ok")

	log.Println("Table initialize")
	_, err = table.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("KeyValue initialize")
	_, err = key_value.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Initialize done.")
}

func Alarm() {
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
	anniv, err := anniversary.New()
	if err != nil {
		log.Fatal(err)
	}

	err = anniv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
