package weathercrawler_test

import (
	"fmt"
	"testing"

	"github.com/thanhpp/discordbots/pkg/weathercrawler"
)

func TestWeatherCrawler(t *testing.T) {
	for i := 0; i < 10; i++ {
		w := weathercrawler.NewWeatherCrawler()

		fmt.Printf("%+v \n", w.GetInfoNow())
	}
}
