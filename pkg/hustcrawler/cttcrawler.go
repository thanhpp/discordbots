package hustcrawler

import (
	"github.com/gocolly/colly/v2"
)

const (
	cttDomain       = `https://ctt.hust.edu.vn`
	cttUndergradURL = `https://ctt.hust.edu.vn/DisplayWeb/DisplayListBaiViet?tag=ĐTĐH`
)

type CTTCrawler struct {
	collector *colly.Collector
}

func NewCTTCrawler() *CTTCrawler {
	c := colly.NewCollector()
	cttCrawler := &CTTCrawler{
		collector: c,
	}

	return cttCrawler
}

type CTTNew struct {
	Title string `selector:"ul > li:nth-child(1) > div > div > a > p.title"`
	Date  string `selector:"ul > li:nth-child(1) > div > div > p"`
	URL   string `selector:"ul > li:nth-child(1) > div > div > a.contentTitle" attr:"href"`
}

func (c CTTCrawler) GetLatestNew() *CTTNew {
	var (
		latestNew = new(CTTNew)
	)

	c.collector.OnHTML("div.col-md-9.col-xs-12", func(h *colly.HTMLElement) {
		if err := h.Unmarshal(latestNew); err != nil {
			return
		}
		latestNew.URL = cttDomain + latestNew.URL
	})

	c.collector.Visit(cttUndergradURL)

	return latestNew
}
