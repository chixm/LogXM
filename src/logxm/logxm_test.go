package logxm

import (
	"testing"
	"time"
)

func TestLogging(t *testing.T) {
	t.Log(`Logxm Test`)

	c := StandardConfig()

	SetupLog(c)

	TerminateLogging()

	time.Sleep(1 * time.Minute)
}
