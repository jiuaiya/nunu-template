# Backend Agent Instructions

These instructions apply to this repository and to projects generated from it.

- Treat Nunu as a generator only. Generated source belongs to the application and must not import `github.com/go-nunu/nunu` at runtime.
- Use the pinned Wire generator for compile-time dependency assembly. Commit `wire_gen.go`; do not introduce runtime DI containers, package-level service locators, or hidden `init` registration.
- HTTP handlers validate protocol input and call services. They must not issue SQL or Redis commands directly.
- Services own use-case orchestration and depend on repository interfaces where substitution is useful.
- Repository and platform packages are the only layers allowed to import GORM, database drivers, or Redis clients.
- Never log raw request or response bodies, credentials, tokens, cookies, authorization headers, or DSNs.
- Every outbound operation must accept `context.Context` and have a bounded timeout at the caller or client level.
- Schema changes use reviewed migrations. Production code must not call `AutoMigrate`.
- Add tests for behavior changes. Run `make check` before handoff.
- Run `govulncheck ./...` after dependency changes and commit `go.mod` with `go.sum`.
