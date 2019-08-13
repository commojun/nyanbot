package api

import (
	"net/http"

	"github.com/commojun/nyanbot/app/linebot"
)

var (
	message = API{
		Name: "/message",
		Post: postMessage,
	}
)

func postMessage(req *http.Request, res *Response) error {
	bot, err := linebot.New()
	if err != nil {
		res.Status = http.StatusInternalServerError
		return err
	}

	var parsedReq struct {
		RoomKey string `json:"room_key"`
		Message string `json:"message"`
	}
	err = parseJSONRequest(req, &parsedReq)
	if err != nil {
		res.Status = http.StatusInternalServerError
		return err
	}

	// TODO 申し訳程度のTOKEN照合機能を追加

	err = bot.TextMessageWithRoomKey(parsedReq.Message, parsedReq.RoomKey)
	if err != nil {
		res.Status = http.StatusInternalServerError
		return err
	}
	return nil
}
