package echo

import (
	"fmt"
	"log"
	"net/http"

	"github.com/commojun/nyanbot/app/linebot"
	linebotOrg "github.com/line/line-bot-sdk-go/linebot"
)

type Echo struct {
	Bot *linebot.LineBot
}

func New() (*Echo, error) {
	bot, err := linebot.New()
	if err != nil {
		return &Echo{}, err
	}

	return &Echo{Bot: bot}, nil
}

func (echo *Echo) StartServer() error {
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("ping\n")
		// この辺はlinebot側でラップしてしまったほうがよさそう
		events, err := echo.Bot.Client.ParseRequest(req)
		if err != nil {
			if err == linebotOrg.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebotOrg.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebotOrg.TextMessage:
					fmt.Printf("%v", message)
					fmt.Printf("%v", event.Source)
					if _, err = echo.Bot.Client.ReplyMessage(event.ReplyToken, linebotOrg.NewTextMessage(message.Text)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})

	if err := http.ListenAndServe(":1337", nil); err != nil {
		return err
	}

	return nil
}
