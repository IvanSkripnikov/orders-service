package logger

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	panicTriggerPush = promauto.NewCounter(prometheus.CounterOpts{
		Name: "triggered_send_to_panic_log_total",
		Help: "Total number of panic level messages sent to the log",
	})
	fatalTriggerPush = promauto.NewCounter(prometheus.CounterOpts{
		Name: "triggered_send_to_fatal_log_total",
		Help: "Total number of fatal level messages sent to the log",
	})
	errorTriggerPush = promauto.NewCounter(prometheus.CounterOpts{
		Name: "triggered_send_to_error_log_total",
		Help: "Total number of error level messages sent to the log",
	})
	warningTriggerPush = promauto.NewCounter(prometheus.CounterOpts{
		Name: "triggered_send_to_warning_log_total",
		Help: "Total number of warning level messages sent to the log",
	})
	infoTriggerPush = promauto.NewCounter(prometheus.CounterOpts{
		Name: "triggered_send_to_info_log_total",
		Help: "Total number of info level messages sent to the log",
	})
	debugTriggerPush = promauto.NewCounter(prometheus.CounterOpts{
		Name: "triggered_send_to_debug_log_total",
		Help: "Total number of debug level messages sent to the log",
	})
	traceTriggerPush = promauto.NewCounter(prometheus.CounterOpts{
		Name: "triggered_send_to_trace_log_total",
		Help: "Total number of trace level messages sent to the log",
	})
)
