package timer

import "time"

type Stopwatch struct {
	start   time.Time
	elapsed time.Duration
	running bool
}

func Start() *Stopwatch {
	return &Stopwatch{start: time.Now(), running: true}
}

func (s *Stopwatch) Elapsed() time.Duration {
	if s.running {
		return s.elapsed + time.Since(s.start)
	}
	return s.elapsed
}

func (s *Stopwatch) Seconds() float64 {
	return s.Elapsed().Seconds()
}

func (s *Stopwatch) Stop() time.Duration {
	if s.running {
		s.elapsed += time.Since(s.start)
		s.running = false
	}
	return s.elapsed
}

func (s *Stopwatch) Reset() {
	s.start = time.Now()
	s.elapsed = 0
	s.running = true
}
