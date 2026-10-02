# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

Monorepo with a Go backend (`backend/`), a Next.js frontend (`web/`), and OpenAPI specs (`api/openapi/`) that backend codegen consumes. Early-stage: module wiring exists, handlers/contracts are mostly empty, domain logic is starting to land (trainers hours).

## Commands

Root (`Makefile`, uses docker compose): `make up` (build + start backend and Postgres, tails backend logs), `make down`, `make down-clean` (also wipes Postgres volume). Copy `backend/.env.example` to `backend/.env` first. The backend runs in a container with `reflex` hot-reload on `backend/internal` and `backend/cmd`; it listens on port 4000, Postgres on 5432.

Backend (run from `backend/`, see `backend/Makefile`):
- `make lint` — `go tool golangci-lint run ./...`
- `make fmt` — gofumpt (tools are pinned via `go tool`, no separate install)
- `make gen` — `go generate ./internal/...` then fmt
- `make test-unit` — `go test ./internal/...`; single test: `go test ./internal/trainers/... -run TestName`
- `make test-integration` (tag `integration`, sources `.env.test`) and `make test-component` (tag `component`, `./tests/...`) — need a running Postgres
- CI runs lint, test-unit, and `go build ./...`

Web (run from `web/`): `npm run dev`, `npm run lint`, `npm run format:check` / `npm run format`, `npx tsc --noEmit` (CI's "test" step), `npm run build`.

## Backend architecture

Modular monolith; Go module name is `backend`. Modules live in `backend/internal/{users,trainers,trainings}`, each with `module.go` implementing `modules.Module` (`Name`, `Init`, `RegisterHttp`, `RegisterContracts`), plus `domain/` and `ports/`.

- `internal/svc.go` wires everything: builds the Echo server (from the external `github.com/quangbach27/golang-common` library, including token-verifier auth and public/protected routers), then for each module runs `Init` + `RegisterContracts`, calls `contracts.Verify()`, and finally `RegisterHttp`. New modules must be added to the list in `svc.go` and to `contracts.Contracts`/`Verify`.
- Inter-module communication goes only through contracts: each module exposes an interface in `ports/module/client` (e.g. `trainers/ports/module/client.Trainers`), implemented in `ports/module/`, and registered in the shared `modules/contracts.Contracts` struct. Modules should depend on other modules' `client` interfaces, never on their internals.
- `cmd/main.go` creates the config (`internal/configs`, env-driven), pgx pool, and `Svc`. `ExternalServices` (e.g. `TokenVerifier`) is injected from main.

Domain conventions (see `trainers/domain`): entities have unexported fields and are only constructed through a factory (e.g. `HourFactory.NewAvailableHour`), which validates input against a `*FactoryConfig` (defaults set in `NewHourFactory`, overridable via functional options). Validation failures are sentinel errors (`ErrPastHour`) or typed errors carrying context (`TooLateHourError`). Enum-like values use `common.Enum`/`common.MustEnum` from `golang-common`.

## Frontend

Next.js 16 (App Router, `web/src/app`), React 19 with the React Compiler (`babel-plugin-react-compiler`, so avoid manual `useMemo`/`useCallback`), Tailwind 4, ESLint + Prettier (with the Tailwind class-sorting plugin). CORS origin for local dev is `http://localhost:3000` (backend `CORS_ALLOWED_ORIGINS`). Path alias `@/*` → `web/src/*`.

Structure (`web/src`):
- `app/(app)/` — route group whose `layout.tsx` wraps pages in `components/layout/app-layout` (header / scrollable `main` / footer in a full-height flex column). Routes: `trainings`, `set-schedule`. Layouts use Next's global `LayoutProps<'/'>` type helper rather than hand-written prop types.
- `features/<feature>/components/` — feature-specific UI (e.g. `features/set-schedule`); pages in `app/` should stay thin and compose these.
- `components/layout/` — app shell pieces; `components/ui/` — shadcn primitives (style `base-maia`, built on `@base-ui/react`, lucide icons). Add new primitives with the `shadcn` CLI per `web/components.json` rather than hand-writing them; `cn()` lives in `lib/utils.ts`.
