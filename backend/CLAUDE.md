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
- `internal/configs/modules` — the `Module` interface (`Name`, `Init`, `RegisterHttp`, `RegisterContracts`) and `configs/modules/contracts` (`Contracts` struct + `Verify`).
- `internal/configs/auth` — token verification (`clerk_verifier.go`) and role-based authorization: `Resource`, `Action`, `Permission` (with predefined `auth.HourView`, `auth.TrainingCancel`, ...), a static `Policy` per `shared.Role` and `auth.Authorize(ctx, policy, permission)`. Trainer: hour make-available/make-unavailable/view, training reschedule-request/reschedule-response/cancel/view. Attendee: hour schedule/cancel/view, all training actions. The trainers and trainings HTTP handlers call `auth.Authorize` first (401 without session, 403 when denied or the role is empty/unknown) and use the session it returns; users handlers are not covered. `svc.go` builds the `*Policy` once and passes it to `trainers.NewModule` and `trainings.NewModule`; a new module whose handlers authorize takes it the same way.
- `internal/{users,trainers,trainings}` — one module each: `module.go`, `domain/`, `app/{commands,queries}`, `ports/{http,module}`, `adapters/db`. `internal/shared` holds types used by several modules (e.g. `shared.Role`).
- `tests/<module>/` — component tests (see Testing conventions).
- `docker/` — container build files.
- `internal/<module>/adapters/db/` — persistence for a module (see Persistence below).

## Startup flow (`internal/svc.go`)

1. Create the Echo server via `golang-common` (`NewEchoServerWithTokenVerifier`), which provides a public router and a token-protected router.
2. For each module in the list: `Init`, then `RegisterContracts` (fills `contracts.Contracts`).
3. `contracts.Verify()` checks every contract was registered.
4. Each module's `RegisterHttp(ctx, publicRouter, protectedRouter)`.

Adding a module: implement `modules.Module`, add it to the list in `svc.go`, and add its client to `contracts.Contracts` and `Verify`.

## Inter-module rules

Modules talk only through contracts: an interface in `ports/module/client` (e.g. `trainers/ports/module/client.Trainers` with `ScheduleHour` / `CancelHourSchedule`, `(ctx, req) (resp, error)`), implemented in `ports/module/` on top of the module's own command handlers and exposed via `RegisterContracts`. `svc.go` passes the global `*contracts.Contracts` to every `NewModule`; a consuming module defines the narrow interface it needs in its `app/commands` (`ModulesContract`) and the global struct satisfies it. The consumer's handler is built in `Init`/`RegisterHttp`, which run after earlier modules registered their contracts, so keep that order in mind (trainers before trainings). Cross-module writes are not transactional: `ScheduleTraining` books the hour in trainers, stores the training, and releases the hour if the store fails.

## Domain conventions (see `trainers/domain`)

Entities have unexported fields and are built only through factories (`HourFactory.NewAvailableHour(trainerUUID, hour)`) that validate input against a `*FactoryConfig` (defaults in the constructor, overridable via functional options; `Config()` exposes it). The factory returns a `common.Error` (400) whose `Details` list each failure (`invalid-trainer-uuid`, `invalid-training-hour`); the underlying time errors (`ErrPastHour`, `TooLateHourError`, ...) are only kept as detail message text, so `errors.Is` won't match them. State-changing methods (`ScheduleTraining`, `MakeAvailable`, `MakeNotAvailable`) take the caller's `trainerUUID` and return a 403 `common.Error` for a different trainer. Enum-like values use `common.Enum` / `common.MustEnum` from `golang-common`.

`trainings/domain`: `Training` is built by `NewTraining` (400 with details `invalid-training-*`); `Cancel`, `ProposeReschedule`, `ApproveReschedule`, `RejectReschedule` take the acting `User` and return 403 for non-participants and 409 for wrong state; the proposer can't approve their own reschedule.

Loading from the database bypasses validation via `domain.UnmarshalHour` / `UnmarshalTraining` / `UnmarshalUser`. Repositories are interfaces in the domain (`domain.HourRepository`) and take an `upsertFn` callback that mutates the aggregate inside a transaction.

## Persistence (`internal/<module>/adapters/db`)

- Each module owns a Postgres schema named after the module (`trainers`); SQL must be schema-qualified (`trainers.hours`) because migrations set no `search_path`.
- `migrations/*.sql` (golang-migrate, up/down pairs) are embedded in `module.go` and applied in `Init` through `golang-common/db.MigrateDatabaseUp`, which keeps a per-module `schema_migrations` table.
- `queries/*.sql` + `sqlc.yml` generate `dbmodels/` (do not hand-edit). Regenerate with `make gen` (`//go:generate go tool sqlc generate` in `sqlc_gen.go`). `sqlc.yml` maps `trainers.hours.status` to `domain.HourStatus`; `trainer_uuid` is `varchar(255)` so it maps to a Go `string`.
- Trainings: schema `trainings`, table `trainings.trainings`; roles are not stored (attendee/trainer roles are fixed by `NewTraining`, so `toDomain` sets `shared.RoleAttendee`/`RoleTrainer`), and the reschedule proposer is only `move_proposed_by_id`, resolved to the trainer or attendee when loading. A single `UpsertTraining` query backs both `AddTraining` and `UpdateTraining`; the latter relies on the repeatable-read transaction from `commonDb.UpdateInTx` (retried on serialization failures) instead of `FOR UPDATE`. Partial unique indexes on `(hour, trainer_id)` and `(hour, attendee_id)` (`WHERE NOT canceled`) forbid double booking; a violation currently surfaces as a generic DB error. There is a single migration file (`00001`), so a dev database from an older layout must be recreated (`make down-clean`). Reads for the HTTP side go through a read model (`user_trainings_read_model.go`) instead of the repository.
- The repository (`hour_repository.go`) wraps queries in `commonDb.UpdateInTx`; a missing row falls back to `HourFactory.NewNotAvailableHour`.

## OpenAPI (`/api/openapi`)

API specs live at the repo root in `api/openapi/<module>.yml` (e.g. `trainers.yml`), outside the Go module. `make gen` produces two outputs from the trainers spec, each configured by an `oapi-config.yml` next to its `oapi_gen.go` directive (which points back at the spec with `../` paths):
- `trainers/ports/http/openapi.gen.go` — Echo strict-server interface + models (`StrictServerInterface`, implemented by `Handler`).
- `trainers/ports/http/clients/client.gen.go` — HTTP client used by the component tests in `tests/`.

`trainings/ports/http` follows the same layout from `trainings.yml` (`POST`/`GET /api/v1/trainings`, `POST /api/v1/trainings/{trainingUuid}/cancel` (204; the role comes from the session, trainer if it has the trainer role; a malformed uuid is a 400); the GET lists upcoming trainings including canceled ones, paginated with `page` (from 1) and `pageSize` (default 10, max 50); response `{items, pagination: {page, pageSize, total}}`). Handlers read the user from the session (`Subject`, `Role`), after `auth.Authorize`. The token carries no username: `ScheduleTraining` and `CancelTraining` look up the display name of each participant with the users contract (`GetUser`, 404 for an unknown user, so both trainer and attendee must exist in `users`), and the `trainerUsername` field of the schedule request is ignored. The session user is the attendee/trainer, never a value from the request.

Never hand-edit `*.gen.go`; change the spec and rerun `make gen`.

## Testing conventions

Test split: domain rules and query/command `Validate()` get unit tests; app handlers and HTTP handlers are covered only by component tests (`tests/<module>`), not by their own unit tests; SQL goes through `integration` tests next to the adapter. Component tests start the whole app in-process on `PORT` from `.env.test`, so packages must run serially (`-p 1`, already in `make test-component`); `tests.NewSession(userID, roles...)` registers a bearer token. If the dev database has a migration version the repo doesn't know (`no migration found for version N`), run the tests against a fresh database (`DB_NAME=<tmp> go test -p 1 -tags component ./tests/...`; `make` re-sources `.env.test` and overrides `DB_NAME`).

Tests use `github.com/stretchr/testify`: `require` for checks that must stop the test (errors, setup, anything later assertions depend on) and `assert` for independent value checks. Prefer table-driven tests with `t.Run` + `t.Parallel()`; domain tests live in the external `domain_test` package (see `trainers/domain/factory_test.go`).

## Config

Env vars (see `.env.example`): `PORT`, `APP_ENV`, `CORS_ALLOWED_ORIGINS`, `DB_USERNAME`, `DB_PASSWORD`, `DB_NAME`, `DB_HOST`, `DB_PORT`, `DB_SSLMODE` (default `disable`; the deploy sets `require` because RDS forces SSL). Clerk: `CLERK_SECRET_KEY`, `CLERK_AUTHORIZED_PARTIES` (allowed `azp` origins, `;`-separated) and `CLERK_WEBHOOK_SIGNING_SECRET`. Copy `.env.example` to `.env` before `make up`.

## Authentication (Clerk)

`cmd/main.go` `newTokenVerifier`: with `CLERK_SECRET_KEY` set it uses `internal/clerkauth.Verifier` (`clerk-sdk-go` `jwt.Verify`: signature via cached JWKS, expiry, issuer, plus `azp` against `CLERK_AUTHORIZED_PARTIES`); without it only `APP_ENV=dev` starts, with the fixed mock tokens (component tests always use the stub). The session maps `sub` to `Subject` and the active organization role (`org:trainer` -> `trainer`) to `Role`; there are no custom claims, and usernames and credit are never in the token (usernames come from the users contract).

The `users` module owns role and credit keyed by the Clerk user id. `POST /api/v1/webhooks/clerk` (`users/ports/http/clerk_webhook.go`, public route, registered only when `CLERK_WEBHOOK_SIGNING_SECRET` is set) verifies the Svix signature over the raw body and handles `user.created` (idempotent `EnsureUser`, default role attendee, mirrors the role to Clerk public metadata through `adapters.ClerkRoleMirror` when the secret key is set) and `user.updated` (rename, or create if unknown); other events, including `user.deleted`, are acknowledged and ignored. Locally Clerk needs a tunnel to reach the webhook; there is no get-or-create fallback yet, so without the webhook a new Clerk user has no row in `users` (use the tunnel, or call `POST /api/v1/users/onboard` once).
