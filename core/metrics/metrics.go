package metrics

// IStats interface to any metric collector
type Metrics interface {
	Increment(int, string, string) error
	Decrement(int, string, string) error
	Set(int, string, string, int) error
}

const (
	PROC_IN = iota + 1
	PROC_OUT
	PACKET_DROP
	CONNECTION_TRANSIT
)

func New() *MetricsVoid {
	return &MetricsVoid{}
}

type MetricsVoid struct{}

func (o *MetricsVoid) Decrement(metric int, pipelineNamestring string, name string) error {
	return nil
}
func (o *MetricsVoid) Increment(metric int, pipelineNamestring string, name string) error {
	return nil
}
func (o *MetricsVoid) Set(metric int, pipelineNamestring string, name string, v int) error { return nil }
