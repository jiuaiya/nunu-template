.PHONY: bootstrap wire dev build test check vuln

bootstrap:
	go mod download
	go install github.com/google/wire/cmd/wire@v0.7.0

wire:
	go run github.com/google/wire/cmd/wire@v0.7.0 ./cmd/server/wire

dev:
	go run ./cmd/server

build:
	mkdir -p bin
	go build -trimpath -o bin/server ./cmd/server

test:
	go test -race -cover ./...

vuln:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

check:
	./scripts/check.sh
