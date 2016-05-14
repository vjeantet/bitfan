package metrics

import "github.com/quipo/statsd"

// IStats interface to any metric collector
type IStats interface {
	Increment(string, int64) error
	Decrement(string, int64) error
}

type StatsVoid struct {
}

func (o *StatsVoid) Increment(name string, v int64) error { return nil }
func (o *StatsVoid) Decrement(name string, v int64) error { return nil }

type statsStatd struct {
	*statsd.StatsdBuffer
}

func (s *statsStatd) Increment(name string, v int64) error {
	return s.Incr(name, v)
}
func (s *statsStatd) Decrement(name string, v int64) error {
	return s.Decr(name, v)
}

// func NewStatd() {
// 	interval := time.Millisecond * 100 // aggregate stats and flush every ..

// 	statsdclient := statsd.NewStatsdClient("192.168.59.103:8125", "veino.")
// 	statsdclient.CreateSocket()
// 	stats := statsd.NewStatsdBuffer(interval, statsdclient)
// 	statsi := &statsStatd{stats}
// 	SetIStat(statsi)
// }
