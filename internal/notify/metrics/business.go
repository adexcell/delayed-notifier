// Package metrics contains metrics for usecases
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// NotificationsCreatedTotal counts how many delayed notifications were created.
	NotificationsCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "delayed_notifier_notifications_created_total",
		Help: "The total number of created delayed notifications",
	})

	// NotificationsProcessedTotal counts how many notifications were successfully processed.
	NotificationsProcessedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "delayed_notifier_notifications_processed_total",
		Help: "The total number of successfully processed notifications",
	})

	// NotificationsFailedTotal counts how many notifications were failed.
	NotificationsFailedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "delayed_notifier_notifications_failed_total",
		Help: "The total number of failed notifications",
	})
)
