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

generate-certs:
	rm *.pem
	openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/CN=*.arrebato"
	openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/CN=server.arrebato"
	openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem
	openssl req -newkey rsa:4096 -nodes -keyout client-key.pem -out client-req.pem -subj "/CN=*.client.arrebato"
	openssl x509 -req -in client-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem
	rm *-req.pem
