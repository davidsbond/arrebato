//go:build tools

package tools

import (
	_ "github.com/anchore/syft"
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "go.etcd.io/bbolt/cmd/bbolt"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
	_ "mvdan.cc/gofumpt"
	_ "sigs.k8s.io/kustomize/kustomize/v4"
)
