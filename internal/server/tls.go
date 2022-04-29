package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

type (
	// The TLSConfig type contains configuration values for TLS between the server nodes and clients.
	TLSConfig struct {
		// Location of the TLS certificate file to use for transport credentials.
		CertFile string

		// Location of the TLS Key file to use for transport credentials.
		KeyFile string

		// The certificate of the CA that signs client certificates.
		CAFile string
	}
)

func (c TLSConfig) enabled() bool {
	return c.CertFile != "" && c.KeyFile != "" && c.CAFile != ""
}

func (c TLSConfig) tlsConfig() (*tls.Config, error) {
	ca, err := ioutil.ReadFile(c.CAFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA file: %w", err)
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(ca)

	cert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS files: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
		MinVersion:   tls.VersionTLS13,
	}, nil
}
