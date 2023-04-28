package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

var m *metrics

type metrics struct {
	endpoint_calls *prometheus.CounterVec
}

func NewMetrics(reg prometheus.Registerer) {
	m = &metrics{
		endpoint_calls: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "minitwit",
			Name:      "endpoint_calls",
			Help:      "The number of requests made to endpoints.",
		}, []string{"name", "method"}),
	}
	reg.MustRegister(m.endpoint_calls)
}

func CountEndpoint(name string, method string) {
	m.endpoint_calls.With(prometheus.Labels{"name": name, "method": method}).Inc()
}
