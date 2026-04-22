package weather_notify

import (
	"context"
	"log"

	"github.com/commojun/nyanbot/app/weather"
	"github.com/commojun/nyanbot/config"
)

type LineBot interface {
	TextMessageWithRoomKey(ctx context.Context, msg string, roomKey string) error
}

// WeatherNotifier は朝の天気予報をプッシュ通知する
type WeatherNotifier struct {
	CityID  string
	RoomKey string
	Bot     LineBot
}

// New は WeatherNotifier を生成する
func New(bot LineBot, cfg config.Config) *WeatherNotifier {
	return &WeatherNotifier{
		CityID:  cfg.WeatherCityID,
		RoomKey: cfg.WeatherRoomKey,
		Bot:     bot,
	}
}

// Run は天気予報を取得してLINEに送信する
func (wn *WeatherNotifier) Run(ctx context.Context) error {
	if wn.CityID == "" {
		log.Println("[weather_notify] NYAN_WEATHER_CITY_ID が未設定のためスキップ")
		return nil
	}
	if wn.RoomKey == "" {
		log.Println("[weather_notify] NYAN_WEATHER_ROOM_KEY が未設定のためスキップ")
		return nil
	}

	msg, err := weather.Fetch(ctx, wn.CityID)
	if err != nil {
		return err
	}

	return wn.Bot.TextMessageWithRoomKey(ctx, msg, wn.RoomKey)
}
