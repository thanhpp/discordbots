package weathercrawler_test

import (
	"fmt"
	"testing"

	"github.com/thanhpp/discordbots/pkg/weathercrawler"
)

func TestWeatherCrawlerInfoNow(t *testing.T) {
	w := weathercrawler.NewWeatherCrawler()
	for i := 0; i < 10; i++ {
		fmt.Printf("%+v \n", w.GetInfoNow())
	}
}
