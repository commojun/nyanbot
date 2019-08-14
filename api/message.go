package api

import (
	"fmt"
	"net/http"

	"github.com/commojun/nyanbot/app/linebot"
	"github.com/commojun/nyanbot/constant"
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
		Token   string `json:"token"`
	}
	err = parseJSONRequest(req, &parsedReq)
	if err != nil {
		res.Status = http.StatusInternalServerError
		return err
	}

	if parsedReq.Token != constant.MessageToken {
		res.Status = http.StatusInternalServerError
		return fmt.Errorf("Token does not match")
	}

	err = bot.TextMessageWithRoomKey(parsedReq.Message, parsedReq.RoomKey)
	if err != nil {
		res.Status = http.StatusInternalServerError
		return err
	}
	return nil
}
