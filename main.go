package main

import (
	"log"
	"github.com/zhengyi13/cert-exporter/prober"
)

func main() {
	log.Println("exporter main")
	var h prober.HostPort
	h = "foo"
	prober.Probe(h)
}
