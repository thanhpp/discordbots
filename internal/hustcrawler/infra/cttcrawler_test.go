package infra_test

import (
	"fmt"
	"testing"

	"github.com/thanhpp/discordbots/internal/hustcrawler/infra"
)

func TestCTTCrawler(t *testing.T) {
	cttCrawler := infra.NewCTTCrawler()

	fmt.Println(cttCrawler.GetLatestNew())
}
