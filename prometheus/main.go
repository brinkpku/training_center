package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			processDur.Observe(float64(rand.Intn(15)))
			actionCounter.WithLabelValues([2]string{"post", "get"}[rand.Intn(2)]).Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "global",
		Name:      "myapp_processed_ops_total",
		Help:      "The total number of processed events",
	})
	processDur = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "myapp_processed_duration",
		Help: "The latency of process",
	})
	actionCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "global",
		Subsystem: "detail",
		Name:      "myapp_action",
		Help:      "The action num",
	}, []string{"type"})
)

func main() {
	recordMetrics()

	registry := prometheus.NewRegistry()
	registry.MustRegister(opsProcessed, processDur, actionCounter)
	hanlder := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		Registry: registry,
	})
	// default handler
	// http.Handle("/metrics", promhttp.Handler())
	http.Handle("/metrics", hanlder)
	http.ListenAndServe(":2112", nil)
}
