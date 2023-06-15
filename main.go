package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zhengyi13/cert-exporter/prober"
)

type expirationCollector struct {
	expirationDate *prometheus.Desc
}

func newExpirationCollector() *expirationCollector {
	return &expirationCollector{
		expirationDate: prometheus.NewDesc("expiration_date", "when this cert will expire", nil, nil),
	}
}

func (collector *expirationCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.expirationDate
}

func (collector *expirationCollector) Collect(ch chan<- prometheus.Metric) {

	// Implement logic here to get the actual metric value
	// TODO write that here:
	//
	//
	var mValue int64

	// Write the latest value for the metric(s) to the prometheus metric channel
	m1 := prometheus.MustNewConstMetric(collector.expirationDate, prometheus.GaugeValue, float64(mValue))
	m1 = prometheus.NewMetricWithTimestamp(time.Now(), m1)
	ch <- m1
}

func main() {
	log.Println("exporter main")


	ed := newExpirationCollector()
	prometheus.MustRegister(ed)

	http.Handle("/metrics", promhttp.Handler())



	var h prober.HostPort
	h = "foo"
	prober.Probe(h)
	log.Fatal(http.ListenAndServe(":9201", nil))
}
