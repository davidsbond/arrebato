.DEFAULT_GOAL=build

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o arrebato main.go

proto-generate:
	rm -rf internal/proto
	buf generate

proto-lint:
	buf lint

proto-breaking:
	buf breaking --against 'https://github.com/davidsbond/arrebato.git#branch=master'

proto-format:
	buf format -w

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
	goreleaser release --snapshot --rm-dist --skip-sign

release:
	goreleaser release --rm-dist

install-tools: install-buf install-kustomize install-protoc-plugins install-golangci-lint install-gofumpt install-bbolt install-syft install-kubeval install-goreleaser

kustomize:
	kustomize build deploy/kustomize -o install.yaml

install-kustomize:
	go install sigs.k8s.io/kustomize/kustomize/v4

install-syft:
	go install github.com/anchore/syft

install-goreleaser:
	go install github.com/goreleaser/goreleaser

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

install-kubeval:
	go install github.com/instrumenta/kubeval

generate-certs:
	rm *.pem
	openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/CN=*.test.com"
	openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/CN=*.test.com"
	openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile hack/v3.ext
	openssl req -newkey rsa:4096 -nodes -keyout client-key.pem -out client-req.pem -subj "/CN=*.test.com"
	openssl x509 -req -in client-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -extfile hack/v3.ext
	rm *-req.pem

generate-serf-key:
	rm -rf serf.key
	hexdump -n 16 -e '4/4 "%08X" 1 "\n"' /dev/urandom > serf.key

kubeval: kustomize
	kubeval install.yaml

update-distroless:
	./scripts/update_distroless.sh
