package hustcrawler_test

import (
	"fmt"
	"testing"

	"github.com/thanhpp/discordbots/pkg/hustcrawler"
)

func TestCTTCrawler(t *testing.T) {
	cttCrawler := hustcrawler.NewCTTCrawler()

	fmt.Printf("%+v", cttCrawler.GetLatestNew())
}
