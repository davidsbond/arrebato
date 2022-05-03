// Package cmd contains the command-line implementation for the CLI.
package cmd

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"

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

func withClient(fn func(ctx context.Context, client *arrebato.Client, args []string) error) func(command *cobra.Command, args []string) error {
	return func(command *cobra.Command, args []string) error {
		client, err := loadClient(command.Context())
		if err != nil {
			return err
		}

		defer closeIt(client)
		return fn(command.Context(), client, args)
	}
}

func validTopicNames(command *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client, err := loadClient(command.Context())
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	defer closeIt(client)
	topics, err := client.Topics(command.Context())
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	names := make([]string, len(topics))
	for i, topic := range topics {
		names[i] = topic.Name
	}

	sort.Slice(names, func(i, j int) bool {
		if strings.HasPrefix(names[i], toComplete) {
			return true
		}

		return names[i] < names[j]
	})

	return names, cobra.ShellCompDirectiveNoSpace
}

func validNodeNames(command *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client, err := loadClient(command.Context())
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	defer closeIt(client)
	nodes, err := client.Nodes(command.Context())
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	names := make([]string, len(nodes))
	for i, node := range nodes {
		names[i] = node.Name
	}

	sort.Slice(names, func(i, j int) bool {
		if strings.HasPrefix(names[i], toComplete) {
			return true
		}

		return names[i] < names[j]
	})

	return names, cobra.ShellCompDirectiveNoSpace
}

func closeIt(c io.Closer) {
	if err := c.Close(); err != nil {
		hclog.Default().Error("failed to close", "error", err)
	}
}

func confirmf(input io.Reader, output io.Writer, prompt string, args ...interface{}) (bool, error) {
	r := bufio.NewReader(input)

	fmt.Fprintf(output, prompt, args...)
	fmt.Fprint(output, " [y/N]: ")

	value, err := r.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read input: %w", err)
	}

	if len(value) < 2 {
		return false, nil
	}

	return strings.ToLower(strings.TrimSpace(value))[0] == 'y', nil
}
