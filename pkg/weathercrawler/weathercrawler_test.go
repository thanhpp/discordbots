package weathercrawler_test

import (
	"fmt"
	"testing"

	"github.com/thanhpp/discordbots/pkg/weathercrawler"
)

func TestWeatherCrawler(t *testing.T) {
	w := weathercrawler.NewWeatherCrawler()

	fmt.Printf("%+v", w.GetInfoNow())
}
