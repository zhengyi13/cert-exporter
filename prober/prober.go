package prober

import (
	"crypto/tls"
	"fmt"
	"log"
)

type HostPort string

func Probe(hp HostPort) (timestamp int64, err error){
	log.Printf("probing %s", hp)
	config :=tls.Config{}
	conn, err := tls.Dial("tcp", string(hp), &config)
	// error case; die early
	if err != nil {
		return timestamp, fmt.Errorf("dial failure: %w", err)
	}
	defer conn.Close()
	// otherwise, get cert details
	// TODO(zhengyi) turn this into Promethus metrics
	log.Printf("client connected to %s", conn.RemoteAddr().String())
	s := conn.ConnectionState()
	for _, cert := range s.PeerCertificates {
		log.Printf("CN: %s", cert.Subject)
		timestamp = cert.NotAfter.Unix()
	}
	return timestamp, err
}