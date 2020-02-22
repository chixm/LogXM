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
				logger.Info(`Log Every Second`)
			}
		}
	}()

	time.Sleep(1 * time.Minute)
}
