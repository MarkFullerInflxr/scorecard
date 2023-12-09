package utils

import "time"

type Stopwatch struct {
	start time.Time
	end   time.Time
}

func Start() Stopwatch {
	stopwatch := Stopwatch{}
	stopwatch.start = time.Now()
	return stopwatch
}

func (s *Stopwatch) Stop() {
	s.end = time.Now()
}

func (s *Stopwatch) Elapsed() time.Duration {
	return s.end.Sub(s.start)
}
