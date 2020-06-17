package profile

import "time"

type StopwatchMode uint8

const (
	ModeNormal StopwatchMode = iota
	ModeAutoRestart
)

type Stopwatch struct {
	start time.Time
	dots  map[string]time.Time
	mode  StopwatchMode
}

func NewStopwatch(mode StopwatchMode) *Stopwatch {
	return &Stopwatch{
		dots: make(map[string]time.Time),
		mode: mode,
	}
}

func (p *Stopwatch) Start() {
	p.start = time.Now()
}

func (p *Stopwatch) Dot(name string) {
	d := time.Now()
	p.dots[name] = d
}

func (p *Stopwatch) GetDotDuration(name string) time.Duration {
	if d, ok := p.dots[name]; ok {
		return d.Sub(p.start)
	}
	return 0
}
