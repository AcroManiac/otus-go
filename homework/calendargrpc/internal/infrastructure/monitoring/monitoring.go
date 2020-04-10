package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

func NewSummaryVec(ns, name, help string) *prometheus.SummaryVec {
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: ns,
		Name:      name,
		Help:      help,
	},
		[]string{"service"})

	prometheus.MustRegister(summaryVec)
	return summaryVec
}

func NewCounterVec(ns, name, help string) *prometheus.CounterVec {
	counterVec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: ns,
		Name:      name,
		Help:      help,
	},
		[]string{"service"})

	prometheus.MustRegister(counterVec)
	return counterVec
}
