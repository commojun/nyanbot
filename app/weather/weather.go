package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "https://weather.tsukumijima.net/api/forecast/city/%s"

// Forecast はtsukumijima天気APIのレスポンス構造体
type Forecast struct {
	Title     string        `json:"title"`
	Forecasts []DayForecast `json:"forecasts"`
}

// DayForecast は1日分の天気予報
type DayForecast struct {
	DateLabel   string      `json:"dateLabel"`  // "今日", "明日", "明後日"
	Telop       string      `json:"telop"`      // "晴れ" 等
	Temperature Temperature `json:"temperature"`
}

// Temperature は最高・最低気温
type Temperature struct {
	Min *Celsius `json:"min"`
	Max *Celsius `json:"max"`
}

// Celsius は摂氏気温
type Celsius struct {
	Celsius string `json:"celsius"`
}

// Fetch は指定したcityIDの天気予報を取得してテキストを返す
func Fetch(ctx context.Context, cityID string) (string, error) {
	url := fmt.Sprintf(baseURL, cityID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("リクエスト生成失敗: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("APIリクエスト失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("APIエラー: ステータスコード %d", resp.StatusCode)
	}

	var forecast Forecast
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		return "", fmt.Errorf("レスポンスのパース失敗: %w", err)
	}

	return formatForecast(&forecast), nil
}

// formatForecast は天気予報データをLINE送信用テキストに変換する
func formatForecast(f *Forecast) string {
	msg := fmt.Sprintf("【%s】\n", f.Title)
	for _, d := range f.Forecasts {
		if d.DateLabel != "今日" && d.DateLabel != "明日" {
			continue
		}
		minTemp := "?"
		maxTemp := "?"
		if d.Temperature.Min != nil {
			minTemp = d.Temperature.Min.Celsius
		}
		if d.Temperature.Max != nil {
			maxTemp = d.Temperature.Max.Celsius
		}
		msg += fmt.Sprintf("%s: %s（最低%s℃ / 最高%s℃）\n", d.DateLabel, d.Telop, minTemp, maxTemp)
	}
	return msg
}
