package main

import (
	"log"

	"github.com/alecthomas/kong"
	"github.com/commojun/nyanbot/app/alarm"
	"github.com/commojun/nyanbot/app/anniversary"
	"github.com/commojun/nyanbot/app/hello"
	"github.com/commojun/nyanbot/app/linebot"
	"github.com/commojun/nyanbot/app/server"
	"github.com/commojun/nyanbot/cache"
	"github.com/commojun/nyanbot/config"
)

type CLI struct {
	config.Config `embed:""`

	Server      ServerCmd      `cmd:"" help:"Start HTTP server"`
	Hello       HelloCmd       `cmd:"" help:"Send hello message"`
	Alarm       AlarmCmd       `cmd:"" help:"Run alarm checker"`
	Anniversary AnniversaryCmd `cmd:"" help:"Run anniversary checker"`
}

// ServerCmd: HTTP サーバーを起動
type ServerCmd struct{}

func (cmd *ServerCmd) Run(ctx *CLI) error {
	if err := cache.Initialize(ctx.Config); err != nil {
		log.Fatal(err)
	}

	srv, err := server.New(ctx.Config)
	if err != nil {
		log.Fatal(err)
	}

	return srv.Start()
}

// HelloCmd: テスト用 hello メッセージを送信
type HelloCmd struct{}

func (cmd *HelloCmd) Run(ctx *CLI) error {
	if err := cache.Initialize(ctx.Config); err != nil {
		log.Fatal(err)
	}

	bot, err := linebot.New(ctx.Config)
	if err != nil {
		log.Fatal(err)
	}

	h := hello.New(bot)
	return h.Say()
}

// AlarmCmd: アラームチェッカーを実行
type AlarmCmd struct{}

func (cmd *AlarmCmd) Run(ctx *CLI) error {
	if err := cache.Initialize(ctx.Config); err != nil {
		log.Fatal(err)
	}

	bot, err := linebot.New(ctx.Config)
	if err != nil {
		log.Fatal(err)
	}

	alm := alarm.New(bot)
	return alm.Run()
}

// AnniversaryCmd: 記念日チェッカーを実行
type AnniversaryCmd struct{}

func (cmd *AnniversaryCmd) Run(ctx *CLI) error {
	if err := cache.Initialize(ctx.Config); err != nil {
		log.Fatal(err)
	}

	bot, err := linebot.New(ctx.Config)
	if err != nil {
		log.Fatal(err)
	}

	anniv := anniversary.New(bot)
	return anniv.Run()
}

func main() {
	var cli CLI
	ctx := kong.Parse(&cli)
	err := ctx.Run(&cli)
	ctx.FatalIfErrorf(err)
}
