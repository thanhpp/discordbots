package logger_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/thanhpp/discordbots/pkg/logger"
)

func Set() {
	logger.Set(&logger.LogConfig{
		Color:      false,
		Level:      "debug",
		LoggerName: "thanhpp's bot",
	}, true)
}

func TestOutputC(t *testing.T) {
	Set()

	var (
		shutDownTick = time.After(time.Second * 5)
	)

	c := logger.Get().OutputC()
	if c == nil {
		t.Error("OutputC is nil")
	}

	go func() {
		for {
			select {
			case msg := <-c:
				fmt.Println(msg)
			}
		}
	}()

	logger.Get().Debug("Test Debug")
	logger.Get().Debug("Test Debug")

	<-shutDownTick
}
