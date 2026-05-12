package timer

import (
	"testing"
	"time"
)

func TestStopwatch_Elapsed(t *testing.T) {
	sw := Start()
	time.Sleep(10 * time.Millisecond)
	d := sw.Elapsed()
	if d < 10*time.Millisecond {
		t.Errorf("expected at least 10ms, got %v", d)
	}
}

func TestStopwatch_Seconds(t *testing.T) {
	sw := Start()
	time.Sleep(10 * time.Millisecond)
	s := sw.Seconds()
	if s < 0.01 {
		t.Errorf("expected at least 0.01s, got %f", s)
	}
}

func TestStopwatch_Stop(t *testing.T) {
	sw := Start()
	time.Sleep(10 * time.Millisecond)
	d := sw.Stop()
	if d < 10*time.Millisecond {
		t.Errorf("expected at least 10ms, got %v", d)
	}
	// After stop, time should remain stable.
	stopped := sw.Elapsed()
	time.Sleep(10 * time.Millisecond)
	if sw.Elapsed() != stopped {
		t.Errorf("expected elapsed stable after stop, got %v (was %v)", sw.Elapsed(), stopped)
	}
}

func TestStopwatch_Reset(t *testing.T) {
	sw := Start()
	time.Sleep(10 * time.Millisecond)
	sw.Stop()
	sw.Reset()
	d := sw.Elapsed()
	if d > time.Millisecond {
		t.Errorf("expected near-zero after reset, got %v", d)
	}
}

func TestStopwatch_ZeroStart(t *testing.T) {
	sw := Start()
	d := sw.Elapsed()
	// Should be negligible at start.
	if d > time.Millisecond {
		t.Errorf("expected near-zero at start, got %v", d)
	}
}
