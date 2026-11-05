# Kale Düellosu

Kale Düellosu is a Turkish-first, competitive mobile game. Product requirements are in `docs/`; architecture and implementation tasks are in `docs/architecture/` and `docs/tasks/`.

## Repository layout

- `apps/mobile/`: Expo + React Native + TypeScript client.
- `services/backend/`: Go modular monolith API and worker.
- `deploy/`: local dependency infrastructure.
- `docs/`: product source, architecture, and task backlog.

This is a monorepo so the mobile client, backend, and shared API contracts can evolve together. It does not imply a microservice deployment.

## Toolchain

- Node.js `24.11.1` (see `.nvmrc`).
- Go `1.26.10` (see `services/backend/go.mod`; raised from 1.25.3 after the vulnerability scan found fixes only in the supported Go 1.26 patch line).
- Docker Engine `28.3.2` and Docker Compose `2.39.1` were used for validation. PostgreSQL and Redis images are pinned by digest in `deploy/compose.yaml`.
- npm `11.2.0` (pinned in the root package manifest).
- npm lockfile is committed for reproducible mobile dependencies.

The Go module currently uses the local module path `castleduello.local/backend` because this workspace has no Git remote or canonical module namespace yet. Replace it with the approved repository path before publishing reusable Go packages.

## Local development

1. Copy `.env.example` to `.env`. Its values are for a disposable local environment only.
2. Start PostgreSQL and Redis:

   ```powershell
   docker compose --env-file .env -f deploy\compose.yaml up -d
   ```

3. Start the API or worker from `services/backend` after dependencies report healthy:

   ```powershell
   Set-Location services\backend
   $env:APP_ENV = "development"
   $env:DATABASE_URL = "postgres://castle_duello:local-development-only@localhost:5432/castle_duello?sslmode=disable"
   $env:REDIS_URL = "redis://localhost:6379/0"
   go run .\cmd\api
   ```

4. In another terminal, start Expo:

   ```powershell
   npm run mobile:start
   ```

The API exposes `GET /health/live` and `GET /health/ready`. Readiness checks PostgreSQL and Redis; liveness does not depend on them. The worker currently provides process startup, dependency validation, and graceful shutdown only; domain jobs are added with their owning backlog tasks.

## Validation

```powershell
Set-Location services\backend
gofmt -w cmd internal
go test ./...
go vet ./...

Set-Location ..\..
npm ci
npm run mobile:typecheck
npm run mobile:lint
docker compose --env-file .env -f deploy\compose.yaml config

go -C services\backend run golang.org/x/vuln/cmd/govulncheck@v1.7.0 ./...
go -C services\backend run github.com/google/go-licenses@v1.6.0 report --ignore castleduello.local/backend ./cmd/api ./cmd/worker
npm audit --audit-level=high
Push-Location apps\mobile
npx --yes --package=license-checker-rseidelsohn@5.0.1 license-checker-rseidelsohn --json --out "$env:TEMP\npm-licenses.json"
Pop-Location
docker run --rm -v "${PWD}:/repo" zricethezav/gitleaks:v8.24.2 detect --source=/repo --no-git --redact --no-banner
```

The root repository has a local clean-clone smoke check; its canonical origin/module path remains pending repository ownership. Docker-backed readiness requires Docker Engine. Never use the example local credentials outside a disposable development environment.

The current Expo toolchain reports 12 moderate npm advisories in `glob`/`minimatch` transitive dependencies. The available automatic fix downgrades Expo to SDK 53; keep SDK 58 and monitor for a compatible upstream fix. The vulnerability CI gate blocks high and critical npm advisories while reporting moderate findings.

Stop local services with `docker compose --env-file .env -f deploy\compose.yaml down`. This preserves the local PostgreSQL volume; remove it only when intentionally resetting local data.

## Conventions

- Keep all API and gameplay decisions server-authoritative.
- Use typed configuration, UTC timestamps, integer economy/gameplay units, and structured logs.
- Never commit `.env`, provider credentials, tokens, or production secrets.
- Preserve the original product documents; update architecture/task documents as implementation decisions are verified.
