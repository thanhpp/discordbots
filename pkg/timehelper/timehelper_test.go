package timehelper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/thanhpp/discordbots/pkg/timehelper"
)

func TestUntilNextTimeCond(t *testing.T) {
	var (
		from = time.Date(2022, time.March, 1, 0, 0, 1, 00, time.UTC)
	)

	timeCond, err := timehelper.NewTimeCond(0, 0, 0)
	if err != nil {
		t.Error(err)
	}

	d, err := timehelper.UntilNextTimeCond(from, timeCond)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(d.Seconds())
}
