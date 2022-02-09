.DEFAULT_GOAL=build

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o arrebato main.go

proto-generate:
	rm -rf internal/proto
	buf generate

proto-lint:
	buf lint

proto-breaking:
	buf breaking --against '.git#branch=master'

docker-compose: build
	docker-compose down
	docker-compose build
	docker-compose up

test:
	go test -race -short ./...

test-e2e:
	cd e2e && go test -race ./...

lint:
	golangci-lint run --enable-all

format:
	gofumpt -l -w .

snapshot:
	goreleaser release --snapshot --rm-dist

install-tools: install-buf install-kustomize install-protoc-plugins install-golangci-lint install-gofumpt install-bbolt install-syft

kustomize:
	kustomize build deploy/kustomize -o install.yaml

install-kustomize:
	go install sigs.k8s.io/kustomize/kustomize/v4

install-syft:
	go install github.com/anchore/syft

install-protoc-plugins:
	go install \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc \
		google.golang.org/protobuf/cmd/protoc-gen-go

install-buf:
	go install github.com/bufbuild/buf/cmd/buf

install-golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

install-gofumpt:
	go install mvdan.cc/gofumpt

install-bbolt:
	go install go.etcd.io/bbolt/cmd/bbolt
