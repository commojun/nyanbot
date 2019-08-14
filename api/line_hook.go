package api

import (
	"net/http"

	"github.com/commojun/nyanbot/app/linebot"
)

var (
	lineHook = API{
		Name: "/callback",
		Post: postLineHook,
	}
)

func postLineHook(req *http.Request, res *Response) error {
	bot, err := linebot.New()
	if err != nil {
		res.Status = http.StatusInternalServerError
		return err
	}

	bot.Events, err = bot.Client.ParseRequest(req)
	if err != nil {
		return err
	}
	if err != nil {
		if linebot.IsInvalidSignature(err) {
			res.Status = http.StatusBadRequest
		} else {
			res.Status = http.StatusInternalServerError
		}
		return err
	}

	bot.ActByEvents()
	return nil
}
