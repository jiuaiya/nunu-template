# Self-maintained Nunu SaaS template

This repository is a Nunu project template, not an application runtime framework. Nunu clones it, rewrites the Go module path in Go files, runs `go mod tidy`, and removes the template Git metadata. Generated services do not depend on Nunu at runtime.

The minimum Go toolchain is the patched `1.26.5` release.

## Generate a service

Pin the CLI version and point it at this Git repository:

```bash
go run github.com/go-nunu/nunu@v1.1.3 new my-service -r <template-git-url>
cd my-service
cp .env.example .env.local
make check
```

Environment files are examples only; the service reads process environment variables and never loads `.env` implicitly.

## Included baseline

- Gin HTTP server with explicit constructors assembled by Wire at compile time
- GORM/MySQL and Redis adapters
- Zap structured logs without raw-body logging
- Prometheus process and HTTP request metrics
- request IDs, panic recovery, request-size and HTTP timeout limits
- liveness/readiness endpoints and graceful shutdown
- race tests, vet, formatting checks, and `govulncheck` in CI

No demo user domain, MongoDB, gRPC, cron runtime, JWT policy, RBAC policy, migrations, or object storage implementation is imposed by the template. Add those only when the product needs them.

## Component generation

Custom component templates are stored in `.nunu/create`:

```bash
nunu create all account -t .nunu/create
```

After generation, add the constructors and route registration explicitly in `cmd/server/wire/wire.go`, then run `make wire`. Generation never mutates application wiring automatically.

## Commands

```bash
make bootstrap
make wire
make dev
make test
make check
make build
```

Wire is pinned to `v0.7.0`. Its generated file is committed, so production builds do not need the Wire binary and have no Wire runtime dependency.

The default service endpoints are:

- `GET /health/live`
- `GET /health/ready`
- `GET /metrics`
- `GET /api/v1/ping`
