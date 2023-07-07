package prober

import (
	"crypto/tls"
	"fmt"
	"log"
	"github.com/zhengyi13/cert-exporter/config"
)

func Probe(hp config.HostPort) (timestamp int64, err error){
	log.Printf("probing %s", hp)
	conn, err := tls.Dial("tcp", hp.String(), &tls.Config{})
	// error case; die early
	if err != nil {
		return timestamp, fmt.Errorf("dial failure: %w", err)
	}
	defer conn.Close()
	// otherwise, get cert details
	log.Printf("client connected to %s", conn.RemoteAddr().String())
	s := conn.ConnectionState()
	for _, cert := range s.PeerCertificates {
		log.Printf("CN: %s", cert.Subject)
		timestamp = cert.NotAfter.Unix()
		log.Printf("Expiration: %d; float: %f", timestamp, float64(timestamp))
	}

	return timestamp, err
}
