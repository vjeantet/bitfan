package core

import (
	"net/http"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "bitfan"

type metricsPrometheus struct {
	agent_packet_in           *prometheus.CounterVec
	agent_packet_out          *prometheus.CounterVec
	connection_packet_transit *prometheus.GaugeVec
	goroutines                prometheus.GaugeFunc
}

func NewPrometheus(address string, path string) Metrics {
	stats := &metricsPrometheus{
		goroutines: prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "runtime",
				Name:      "goroutines_count",
				Help:      "Number of goroutines that currently exist.",
			},
			func() float64 { return float64(runtime.NumGoroutine()) },
		),

		agent_packet_in: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "agent",
			Name:      "packet_consumption",
			Help:      "packets consumed by processors",
		},
			[]string{"pipeline", "agent"},
		),

		agent_packet_out: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "agent",
			Name:      "packet_production",
			Help:      "packets produced by processors",
		},
			[]string{"pipeline", "agent"},
		),

		connection_packet_transit: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "connection",
			Name:      "transit",
			Help:      "packets in transit to processors",
		},
			[]string{"pipeline", "agent"},
		),
	}

	prometheus.MustRegister(stats.agent_packet_in)
	prometheus.MustRegister(stats.agent_packet_out)
	prometheus.MustRegister(stats.connection_packet_transit)
	prometheus.MustRegister(stats.goroutines)

	http.Handle(path, prometheus.Handler())
	go http.ListenAndServe(address, nil)

	return stats
}

func (s *metricsPrometheus) set(metric int, pipelineName string, name string, v int) error {
	switch metric {
	case METRIC_CONNECTION_TRANSIT:
		s.connection_packet_transit.WithLabelValues(pipelineName, name).Set(float64(v))
	}

	return nil
}

func (s *metricsPrometheus) increment(metric int, pipelineName string, name string) error {

	switch metric {
	case METRIC_PROC_OUT:
		s.agent_packet_out.WithLabelValues(pipelineName, name).Inc()
	case METRIC_PROC_IN:
		s.agent_packet_in.WithLabelValues(pipelineName, name).Inc()
	case METRIC_CONNECTION_TRANSIT:
		s.connection_packet_transit.WithLabelValues(pipelineName, name).Inc()
	}

	return nil
}

func (s *metricsPrometheus) decrement(metric int, pipelineName string, name string) error {
	s.connection_packet_transit.WithLabelValues(pipelineName, name).Dec()
	return nil
}
