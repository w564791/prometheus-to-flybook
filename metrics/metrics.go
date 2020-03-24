package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	FlybookMetricsRecived = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "alert_received",
			Help: "Alert received",
			//ConstLabels: prometheus.Labels{"version":"0.1"},

		})
	FlybookMetricsSend = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "alert_send",
			Help: "Alert send",
			//ConstLabels: prometheus.Labels{"version":"0.1"},

		})
	FlybookMetricsCode = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "alert_send_code",
			Help: "Alert send code",
			//ConstLabels: prometheus.Labels{"code":"0.1"},

		}, []string{"code"})
)
