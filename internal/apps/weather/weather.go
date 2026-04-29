package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const baseURL = "https://weather.tsukumijima.net/api/forecast/city/%s"

type Forecast struct {
	Title     string        `json:"title"`
	Forecasts []DayForecast `json:"forecasts"`
}

type DayForecast struct {
	DateLabel   string      `json:"dateLabel"`
	Telop       string      `json:"telop"`
	Temperature Temperature `json:"temperature"`
}

type Temperature struct {
	Min *Celsius `json:"min"`
	Max *Celsius `json:"max"`
}

type Celsius struct {
	Celsius string `json:"celsius"`
}

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

func formatForecast(f *Forecast) string {
	lines := []string{fmt.Sprintf("【%s】", f.Title)}
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
		lines = append(lines, fmt.Sprintf("%s: %s（最低%s℃ / 最高%s℃）", d.DateLabel, d.Telop, minTemp, maxTemp))
	}
	return strings.Join(lines, "\n")
}
