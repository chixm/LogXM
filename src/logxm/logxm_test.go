package logxm

import (
	"testing"
	"time"
)

func TestLogging(t *testing.T) {
	t.Log(`Logxm Test`)

	logger := New(nil)

	go func() {
		t := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-t.C:
				logger.Info(`This is Info level log.`)
				logger.Warn(`Warning level log.`)
				logger.Error(`Error level log.`)
				logger.Fatal(`Fatal level log. This method finishes application with error.`)
			}
		}
	}()

	time.Sleep(1 * time.Minute)
}
