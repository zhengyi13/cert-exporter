package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/zhengyi13/cert-exporter/config"
	"github.com/zhengyi13/cert-exporter/prober"
)

type expirationCollector struct {
	expirationDate *prometheus.Desc
	probeConfig *config.Config
	
}

func newExpirationCollector() *expirationCollector {
	return &expirationCollector{
		expirationDate: prometheus.NewDesc("expiration_date", "when this cert will expire", nil, nil),
	}
}

func (collector *expirationCollector) LoadConfig(filename string) error {
	// do stuff to load config into collector.RelevantFields
	c, err := config.GetConfig(filename)
	if err != nil {
		return fmt.Errorf("you still need to implement this %s", "bro")
	}
	collector.probeConfig = c
	return nil
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
	for _, hp := range(*collector.probeConfig) {
		ts, err := prober.Probe(hp)
		if err != nil {
			log.Printf("Probing hp (%s, %d) failed: %v", hp.Hostname, hp.Port, err)
			continue
		}
		// Write the latest value for the metric(s) to the prometheus metric channel
		m := prometheus.MustNewConstMetric(collector.expirationDate, prometheus.GaugeValue, float64(ts))
		m = prometheus.NewMetricWithTimestamp(time.Now(), m)
		ch <- m

	}
	
}

func main() {

	configFile := flag.String("config", "", "path to config file")
	flag.Parse()
	
	hosts, err := config.GetConfig(*configFile)
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
