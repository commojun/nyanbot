package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/commojun/nyanbot/internal/apps/alarm"
	"github.com/commojun/nyanbot/internal/apps/anniversary"
	"github.com/commojun/nyanbot/internal/apps/hello"
	"github.com/commojun/nyanbot/internal/config"
	"github.com/commojun/nyanbot/internal/handler"
	"github.com/commojun/nyanbot/internal/linebot"
	"github.com/commojun/nyanbot/internal/masterdata"
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := masterdata.Initialize(ctx, cliCtx.Config); err != nil {
		return err
	}

	srv, err := handler.NewServer(cliCtx.Config)
	if err != nil {
		return err
	}

	return srv.Start(ctx)
}

// HelloCmd: テスト用 hello メッセージを送信
type HelloCmd struct{}

func (cmd *HelloCmd) Run(cliCtx *CLI) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := masterdata.Initialize(ctx, cliCtx.Config); err != nil {
		return err
	}

	bot, err := linebot.New(cliCtx.Config)
	if err != nil {
		return err
	}

	h := hello.New(bot)
	return h.Say(ctx)
}

// AlarmCmd: アラームチェッカーを実行
type AlarmCmd struct{}

func (cmd *AlarmCmd) Run(cliCtx *CLI) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := masterdata.Initialize(ctx, cliCtx.Config); err != nil {
		return err
	}

	bot, err := linebot.New(cliCtx.Config)
	if err != nil {
		return err
	}

	alm := alarm.New(cliCtx.Config, bot)
	return alm.Run(ctx)
}

// AnniversaryCmd: 記念日チェッカーを実行
type AnniversaryCmd struct{}

func (cmd *AnniversaryCmd) Run(cliCtx *CLI) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := masterdata.Initialize(ctx, cliCtx.Config); err != nil {
		return err
	}

	bot, err := linebot.New(cliCtx.Config)
	if err != nil {
		return err
	}

	anniv := anniversary.New(bot)
	return anniv.Run(ctx)
}

func main() {
	var cli CLI
	kongCtx := kong.Parse(&cli)
	err := kongCtx.Run(&cli)
	if errors.Is(err, context.Canceled) {
		log.Println("received shutdown signal, exiting normally")
		return
	}
	kongCtx.FatalIfErrorf(err)
}
