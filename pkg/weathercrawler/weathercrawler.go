package weathercrawler

import (
	"github.com/gocolly/colly/v2"
)

type WeatherCrawler struct {
	forecastClt *colly.Collector
}

func NewWeatherCrawler() *WeatherCrawler {
	return &WeatherCrawler{
		forecastClt: colly.NewCollector(),
	}
}

const (
	todayForeCast = "https://nchmf.gov.vn/Kttvsite/vi-VN/1/thoi-tiet-dat-lien-24h-12h2-15.html"
)

type ForeCastInfo struct {
	Title    string
	Location string
	Info     string
}

func (w *WeatherInfo) GetTodayForecast() *ForeCastInfo {

	return &ForeCastInfo{}
}
