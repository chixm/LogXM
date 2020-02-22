package logxm

import (
	"testing"
	"time"
)

func TestLogging(t *testing.T) {
	t.Log(`Logxm Test`)

	c := StandardConfig()

	logger := New(c)

	go func() {
		t := time.NewTicker(1 * time.Second)
		select {
		case <-t.C:
			logger.Info(`Log Every Second`)
		}
	}()

	time.Sleep(1 * time.Minute)
}
