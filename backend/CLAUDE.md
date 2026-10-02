# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

Go backend (module name `backend`), a modular monolith on Echo + pgx. See the root `../CLAUDE.md` for monorepo-level commands (`make up`, docker compose, web).

## Commands (run from `backend/`)

- `make lint` — `go tool golangci-lint run ./...`
- `make fmt` — gofumpt; `make gen` — `go generate ./internal/...` then fmt; `make tidy` — `go mod tidy`
- `make test-unit` — `go test ./internal/...`; single test: `go test ./internal/trainers/... -run TestName`
- `make test-integration` (tag `integration`, `./internal/...`) and `make test-component` (tag `component`, `./tests/...`) — source `.env.test` and need a running Postgres

Dev tools (golangci-lint, gofumpt, sqlc, oapi-codegen) are pinned via `go tool`; no separate install. `.go-cache/` is build cache; ignore it when searching.

## Layout

- `cmd/main.go` — builds config (`internal/configs`, env-driven), pgx pool and `Svc`; injects `ExternalServices` (e.g. `TokenVerifier`).
- `internal/svc.go` — wires everything (see startup flow below). `Run` starts the Echo server and closes the pgx pool on context cancel.
- `internal/modules` — the `Module` interface (`Name`, `Init`, `RegisterHttp`, `RegisterContracts`) and `modules/contracts` (`Contracts` struct + `Verify`).
- `internal/{users,trainers,trainings}` — one module each: `module.go`, `domain/`, `ports/`, `adapters/`.
- `docker/` — container build files.
- `internal/<module>/adapters/db/` — persistence for a module (see Persistence below).

## Startup flow (`internal/svc.go`)

1. Create the Echo server via `golang-common` (`NewEchoServerWithTokenVerifier`), which provides a public router and a token-protected router.
2. For each module in the list: `Init`, then `RegisterContracts` (fills `contracts.Contracts`).
3. `contracts.Verify()` checks every contract was registered.
4. Each module's `RegisterHttp(ctx, publicRouter, protectedRouter)`.

Adding a module: implement `modules.Module`, add it to the list in `svc.go`, and add its client to `contracts.Contracts` and `Verify`.

## Inter-module rules

Modules talk only through contracts: an interface in `ports/module/client` (e.g. `trainers/ports/module/client.Trainers`), implemented in `ports/module/` and exposed via `RegisterContracts`. Depend on other modules' `client` interfaces, never their internals.

## Domain conventions (see `trainers/domain`)

Entities have unexported fields and are built only through factories (`HourFactory.NewAvailableHour(trainerUUID, hour)`) that validate input against a `*FactoryConfig` (defaults in the constructor, overridable via functional options; `Config()` exposes it). The factory returns a `common.Error` (400) whose `Details` list each failure (`invalid-trainer-uuid`, `invalid-training-hour`); the underlying time errors (`ErrPastHour`, `TooLateHourError`, ...) are only kept as detail message text, so `errors.Is` won't match them. State-changing methods (`ScheduleTraining`, `MakeAvailable`, `MakeNotAvailable`) take the caller's `trainerUUID` and return a 403 `common.Error` for a different trainer. Enum-like values use `common.Enum` / `common.MustEnum` from `golang-common`.

Loading from the database bypasses validation via `domain.UnmarshalHour`. Repositories are interfaces in the domain (`domain.HourRepository`) and take an `upsertFn` callback that mutates the aggregate inside a transaction.

## Persistence (`internal/<module>/adapters/db`)

- Each module owns a Postgres schema named after the module (`trainers`); SQL must be schema-qualified (`trainers.hours`) because migrations set no `search_path`.
- `migrations/*.sql` (golang-migrate, up/down pairs) are embedded in `module.go` and applied in `Init` through `golang-common/db.MigrateDatabaseUp`, which keeps a per-module `schema_migrations` table.
- `queries/*.sql` + `sqlc.yml` generate `dbmodels/` (do not hand-edit). Regenerate with `make gen` (`//go:generate go tool sqlc generate` in `sqlc_gen.go`). `sqlc.yml` maps `trainers.hours.status` to `domain.HourStatus`; `trainer_uuid` is `varchar(255)` so it maps to a Go `string`.
- The repository (`hour_repository.go`) wraps queries in `commonDb.UpdateInTx`; a missing row falls back to `HourFactory.NewNotAvailableHour`.

## OpenAPI (`/api/openapi`)

API specs live at the repo root in `api/openapi/<module>.yml` (e.g. `trainers.yml`), outside the Go module. `make gen` produces two outputs from the trainers spec, each configured by an `oapi-config.yml` next to its `oapi_gen.go` directive (which points back at the spec with `../` paths):
- `trainers/ports/http/openapi.gen.go` — Echo strict-server interface + models (`StrictServerInterface`, implemented by `Handler`).
- `trainers/ports/http/clients/client.gen.go` — HTTP client used by the component tests in `tests/`.

Never hand-edit `*.gen.go`; change the spec and rerun `make gen`.

## Testing conventions

Tests use `github.com/stretchr/testify`: `require` for checks that must stop the test (errors, setup, anything later assertions depend on) and `assert` for independent value checks. Prefer table-driven tests with `t.Run` + `t.Parallel()`; domain tests live in the external `domain_test` package (see `trainers/domain/factory_test.go`).

## Config

Env vars (see `.env.example`): `PORT`, `APP_ENV`, `CORS_ALLOWED_ORIGINS`, `DB_USERNAME`, `DB_PASSWORD`, `DB_NAME`, `DB_HOST`, `DB_PORT`. Copy `.env.example` to `.env` before `make up`.
