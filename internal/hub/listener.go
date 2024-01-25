// internal/hub/listener.go
package hub

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"strconv"
)

type HubConfig struct {
	Host string    `yaml:"host"`
	Port int       `yaml:"port"`
	TLS  TLSConfig `yaml:"tls"`
}

// GetListenAddress returns the address the hub should listen on
func (cfg HubConfig) GetListenAddress() string {
	protocol := "http"
	if cfg.TLS.Enabled {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s:%d", protocol, cfg.Host, cfg.Port)
}

type TLSConfig struct {
	Enabled bool   `yaml:"enabled"`
	Cert    string `yaml:"cert"`
	Key     string `yaml:"key"`
	CA      string `yaml:"ca"`
}

func NewListener(cfg HubConfig) (net.Listener, error) {
	address := cfg.Host + ":" + strconv.Itoa(cfg.Port)

	if cfg.TLS.Enabled {
		cert, err := tls.LoadX509KeyPair(cfg.TLS.Cert, cfg.TLS.Key)
		if err != nil {
			return nil, err
		}

		caCert, err := os.ReadFile(cfg.TLS.CA)
		if err != nil {
			return nil, err
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientCAs:    caCertPool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		}

		listener, err := tls.Listen("tcp", address, tlsConfig)
		if err != nil {
			return nil, err
		}

		return listener, nil
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	return listener, nil
}
