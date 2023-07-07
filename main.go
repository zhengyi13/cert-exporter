package main

import (
	"flag"
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

// newExpirationCollector just gives us our basic collector.
func newExpirationCollector(c *config.Config) *expirationCollector {
	return &expirationCollector{
		expirationDate: prometheus.NewDesc("expiration_date",
			"when this cert will expire in unix timestamp format",
			[]string{"hostport"}, nil),
		probeConfig: c,
	}
}

// Describe is required. Not sure ATM why.
func (collector *expirationCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.expirationDate
}

// Collect actually pulls together the metric and associated labels whenever the /metrics URL is hit.
func (collector *expirationCollector) Collect(ch chan<- prometheus.Metric) {
	for _, hp := range(*collector.probeConfig) {
		ts, err := prober.Probe(hp)
		if err != nil {
			log.Printf("Probing hp (%s, %d) failed: %v", hp.Hostname, hp.Port, err)
			continue
		}
		// Write the latest value for the metric(s) to the prometheus metric channel
		m := prometheus.MustNewConstMetric(
			collector.expirationDate,
			prometheus.GaugeValue,
			float64(ts),
			hp.String(),
		)
		m = prometheus.NewMetricWithTimestamp(time.Now(), m)
		ch <- m
	}
}

func main() {

	configFile := flag.String("config", "", "path to config file")
	flag.Parse()

	config, err := config.GetConfig(*configFile)
	if err != nil {
		log.Fatalf("Unable to load config; bailing out: %v", err)
	}

	ed := newExpirationCollector(config)

	if err != nil {
		log.Fatalf("Failed to GetConfig: %v", err)
	}

	prometheus.MustRegister(ed)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9201", nil))
}
