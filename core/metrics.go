package core

// IStats interface to any metric collector
type Metrics interface {
	increment(int, string, string) error
	decrement(int, string, string) error
	set(int, string, string, int) error
}

const (
	METRIC_PROC_IN = iota + 1
	METRIC_PROC_OUT
	METRIC_PACKET_DROP
	METRIC_CONNECTION_TRANSIT
)

type MetricsVoid struct{}

func (o *MetricsVoid) decrement(metric int, pipelineNamestring string, name string) error {
	return nil
}
func (o *MetricsVoid) increment(metric int, pipelineNamestring string, name string) error {
	return nil
}
func (o *MetricsVoid) set(metric int, pipelineNamestring string, name string, v int) error { return nil }
