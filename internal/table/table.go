// Package table provides the Builder type which is used to produce nice looking tables for printing to the command-line.
package table

import (
	"context"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

type (
	// The Builder type is used to produce nice looking tables.
	Builder struct {
		rows [][]string
	}
)

// NewBuilder returns a new instance of the Builder type that provides methods for building tables. The provided headers
// will show at the top of the table. Headers are automatically made upper case.
func NewBuilder(headers ...string) *Builder {
	for i, header := range headers {
		headers[i] = strings.ToUpper(header)
	}

	return &Builder{
		rows: [][]string{
			headers,
		},
	}
}

// AddRow adds a row to the table, formatting items as strings.
func (b *Builder) AddRow(items ...interface{}) {
	row := make([]string, len(items))
	for i, item := range items {
		row[i] = fmt.Sprint(item)
	}

	b.rows = append(b.rows, row)
}

// Build the table, flushing it to the configured io.Writer.
func (b *Builder) Build(ctx context.Context, output io.Writer) error {
	writer := tabwriter.NewWriter(output, 4, 4, 4, ' ', 0)

	for _, row := range b.rows {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		for i, item := range row {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			if _, err := writer.Write([]byte(item)); err != nil {
				return err
			}

			if i != len(row)-1 {
				if _, err := writer.Write([]byte{'\t'}); err != nil {
					return err
				}
			}
		}

		if _, err := writer.Write([]byte{'\n'}); err != nil {
			return err
		}
	}

	return writer.Flush()
}
