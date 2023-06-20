package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	//	"github.com/zhengyi13/cert-exporter/prober"
	"github.com/zhengyi13/cert-exporter/config"
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

// TODO(zhengyi) rewrite this to take the list of things to probe, and then loop over them in the body
// metric {label_per_hostport?}
func (collector *expirationCollector) Collect(ch chan<- prometheus.Metric) {

	// TODO(zhengyi) Implement logic here to get the actual metric value
	//
	//
	var mValue int64

	// Write the latest value for the metric(s) to the prometheus metric channel
	m1 := prometheus.MustNewConstMetric(collector.expirationDate, prometheus.GaugeValue, float64(mValue))
	m1 = prometheus.NewMetricWithTimestamp(time.Now(), m1)
	ch <- m1
}

func main() {

	// TODO(zhengyi) read config
	myconfigfile := "/home/zhengyi/code/cert-exporter/config.yaml"
	
	hosts, err := config.GetConfig(myconfigfile)
	if err != nil {
		log.Fatalf("Failed to GetConfig: %v", err)
	}
	for _, host := range(*hosts) {
		log.Printf("Getting cert from %s:%d", host.Hostname, host.Port)
	}	
	// TODO(zhengyi) for hp in hostports; do ed.Collect(hp)? Nah, probably
	// just pass the resulting struct/list-o-struct to
	// collector.Collect and let *it* loop over the collection
	ed := newExpirationCollector()
	prometheus.MustRegister(ed)
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":9201", nil))
}
