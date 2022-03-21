// Package cmd contains the command-line implementation for the CLI.
package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"os"
	"strings"

	"github.com/davidsbond/arrebato/pkg/arrebato"
)

func loadClient(ctx context.Context) (*arrebato.Client, error) {
	config := arrebato.Config{
		Addresses: strings.Split(os.Getenv("ARREBATO_SERVERS"), ","),
		ClientID:  os.Getenv("ARREBATO_CLIENT_ID"),
	}

	if _, ok := os.LookupEnv("ARREBATO_TLS"); !ok {
		return arrebato.Dial(ctx, config)
	}

	ca, err := ioutil.ReadFile(os.Getenv("ARREBATO_TLS_CA"))
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(ca)

	cert, err := tls.LoadX509KeyPair(os.Getenv("ARREBATO_TLS_CERT"), os.Getenv("ARREBATO_TLS_KEY"))
	if err != nil {
		return nil, err
	}

	config.TLS = &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
		MinVersion:   tls.VersionTLS13,
	}

	return arrebato.Dial(ctx, config)
}
