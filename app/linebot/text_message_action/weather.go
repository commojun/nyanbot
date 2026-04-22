package text_message_action

import (
	"context"

	"github.com/commojun/nyanbot/app/weather"
)

var (
	getWeather = Action{
		Prefix: "天気",
		Do:     doGetWeather,
	}
)

func doGetWeather(ctx context.Context, tma *TextMessageAction) error {
	// 環境変数からcityIDを取得する手段がないため、東京（130010）をデフォルトとする
	// WeatherCityID は weather_notify 経由でのみ利用可能
	const defaultCityID = "130010"

	msg, err := weather.Fetch(ctx, defaultCityID)
	if err != nil {
		return tma.Bot.TextReply(ctx, "天気の取得に失敗したよ…", tma.Event.ReplyToken)
	}

	return tma.Bot.TextReply(ctx, msg, tma.Event.ReplyToken)
}
