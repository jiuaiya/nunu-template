#!/usr/bin/env sh
set -eu

unformatted="$(gofmt -l .)"
if [ -n "$unformatted" ]; then
  printf '%s\n' 'The following Go files need gofmt:' >&2
  printf '%s\n' "$unformatted" >&2
  exit 1
fi

go vet ./...
go test -race ./...
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
