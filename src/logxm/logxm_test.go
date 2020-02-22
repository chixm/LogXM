package logxm

import (
	"testing"
	"time"
)

func TestLogging(t *testing.T) {
	t.Log(`Logxm Test`)

	c := StandardConfig()

	New(c)

	//TerminateLogging(0)

	time.Sleep(1 * time.Minute)
}
