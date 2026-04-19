package main

import (
	"context"
	"log"

	"github.com/alecthomas/kong"
	"github.com/commojun/nyanbot/app/alarm"
	"github.com/commojun/nyanbot/app/anniversary"
	"github.com/commojun/nyanbot/app/hello"
	"github.com/commojun/nyanbot/app/linebot"
	"github.com/commojun/nyanbot/app/server"
	"github.com/commojun/nyanbot/config"
	"github.com/commojun/nyanbot/masterdata"
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

func (cmd *ServerCmd) Run(cliCtx *CLI) error {
	ctx := context.Background()
	if err := masterdata.Initialize(ctx, cliCtx.Config); err != nil {
		log.Fatal(err)
	}

	srv, err := server.New(cliCtx.Config)
	if err != nil {
		log.Fatal(err)
	}

	return srv.Start()
}

// HelloCmd: テスト用 hello メッセージを送信
type HelloCmd struct{}

func (cmd *HelloCmd) Run(cliCtx *CLI) error {
	ctx := context.Background()
	if err := masterdata.Initialize(ctx, cliCtx.Config); err != nil {
		log.Fatal(err)
	}

	bot, err := linebot.New(cliCtx.Config)
	if err != nil {
		log.Fatal(err)
	}

	h := hello.New(bot)
	return h.Say(ctx)
}

// AlarmCmd: アラームチェッカーを実行
type AlarmCmd struct{}

func (cmd *AlarmCmd) Run(cliCtx *CLI) error {
	ctx := context.Background()
	if err := masterdata.Initialize(ctx, cliCtx.Config); err != nil {
		log.Fatal(err)
	}

	bot, err := linebot.New(cliCtx.Config)
	if err != nil {
		log.Fatal(err)
	}

	alm := alarm.New(bot)
	return alm.Run(ctx)
}

// AnniversaryCmd: 記念日チェッカーを実行
type AnniversaryCmd struct{}

func (cmd *AnniversaryCmd) Run(cliCtx *CLI) error {
	ctx := context.Background()
	if err := masterdata.Initialize(ctx, cliCtx.Config); err != nil {
		log.Fatal(err)
	}

	bot, err := linebot.New(cliCtx.Config)
	if err != nil {
		log.Fatal(err)
	}

	anniv := anniversary.New(bot)
	return anniv.Run(ctx)
}

func main() {
	var cli CLI
	ctx := kong.Parse(&cli)
	err := ctx.Run(&cli)
	ctx.FatalIfErrorf(err)
}
