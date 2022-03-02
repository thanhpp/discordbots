package weathercrawler

import (
	"strings"

	"github.com/gocolly/colly/v2"
)

type WeatherCrawler struct {
	collector *colly.Collector
}

func NewWeatherCrawler() *WeatherCrawler {
	c := colly.NewCollector()
	return &WeatherCrawler{
		collector: c,
	}
}

// #wrapper > div > div.uk-container > section > div > div > article > div > div > div.content-news.fix-content-news > div > div > div:nth-child(2) > div > ul > li:nth-child(2) > div > div.uk-width-3-4
type WeatherInfo struct {
	LastUpdated string `selector:"div.time-update"`
	Temperature string `selector:"ul > li:nth-child(1) > div > div.uk-width-3-4"`
	Status      string `selector:"ul > li:nth-child(2) > div > div.uk-width-3-4"`
	Humidity    string `selector:"ul > li:nth-child(3) > div > div.uk-width-3-4"`
}

func (wi *WeatherInfo) beautifyString(s string) string {
	return strings.TrimSpace(strings.ReplaceAll(s, ":", ""))
}

func (wi *WeatherInfo) Beautify() {
	wi.LastUpdated = strings.TrimSpace(strings.ReplaceAll(strings.Join(strings.Fields(wi.LastUpdated), " "), "Cập nhật: ", ""))
	wi.Temperature = wi.beautifyString(wi.Temperature)
	wi.Status = wi.beautifyString(wi.Status)
	wi.Humidity = wi.beautifyString(wi.Humidity)
}

func (wi WeatherInfo) IsRaining() bool {
	return strings.Contains(strings.ToLower(wi.Status), "có mưa")
}

const (
	hanoiWeather = "https://nchmf.gov.vn/Kttvsite/vi-VN/1/ha-noi-w28.html"
)

func (w *WeatherCrawler) GetInfoNow() *WeatherInfo {
	var infoNow = new(WeatherInfo)

	// #wrapper > div > div.uk-container > section > div > div > article > div > div > div.content-news.fix-content-news > div > div > div:nth-child(2) > div
	w.collector.OnHTML("div.content-news.fix-content-news > div > div > div:nth-child(2) > div", func(h *colly.HTMLElement) {
		if err := h.Unmarshal(infoNow); err != nil {
			return
		}
	})
	w.collector.Visit(hanoiWeather)

	infoNow.Beautify()

	return infoNow
}
