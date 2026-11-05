# Infrastructure and Deployment

## Environments

| Environment | Purpose | Data and integrations |
|---|---|---|
| Local | developer loop | Docker PostgreSQL/Redis, synthetic data, store/push sandboxes or stubs |
| Shared development | integrated team testing | isolated credentials/data, test notifications and purchases |
| Staging | production-like release gate | separate DB/Redis, sandbox store, realistic migrations/load checks |
| Production | live service | restricted credentials, production store/push, backups, monitoring and incident ownership |

Never reuse production data or secrets in lower environments without approved anonymization.

## MVP topology

Mobile clients connect through a TLS load balancer to horizontally scalable Go API/WebSocket instances. PostgreSQL is managed or operated with backups/PITR; Redis is managed or monitored. A worker deployment processes outbox jobs. Object storage/CDN serves approved question media and user-facing assets; do not expose answer keys. Integrate push and Apple/Google purchase verification. Central telemetry receives structured logs, metrics, traces, and crash reports.

No separate message broker, Kubernetes requirement, service mesh, analytics warehouse, or microservices are assumed for MVP.

## Docker and local development

Provide version-pinned backend/mobile dependencies and Docker Compose for PostgreSQL, Redis, and optional local object storage. Include health checks, non-production credentials, migration commands, seed data, and documented reset of named local volumes only. No production secret defaults. API and worker run as separate commands from one application image or build.

## CI/CD and deployment flow

1. Validate formatting, lint, static analysis, unit and contract tests.
2. Build reproducible backend and mobile artifacts; scan dependencies, container and secrets.
3. Run database integration tests and migration from empty schema; verify forward-compatible migration.
4. Publish versioned artifact with commit provenance.
5. Deploy backend to staging; run migration, smoke, WebSocket, reconnect and critical economy checks.
6. Promote with approval to production; use expand/contract database changes and backward-compatible API/event changes.
7. Roll out mobile builds through internal/beta, then staged store release; maintain minimum-version/maintenance behavior.
8. Observe health and error budgets; rollback app/backend/config safely. Do not roll back a destructive migration without a tested recovery plan.

## Secrets and access

Use a managed secret store or environment injection with least privilege and rotation. Separate database, Redis, store, push, analytics and signing credentials per environment. No private keys, provider credentials, passwords, tokens, or production endpoints with secrets in source or mobile bundle. Restrict CI production deploy permissions; audit access.

## Observability and SLOs

Track documented targets: normal API p95 under 300 ms, real-time event p95 under 250 ms in supported regions under agreed load, and initial 99.5% crash-free sessions. Define workload, regions, exclusions and measurement windows before claiming target attainment. Monitor API/WebSocket latency/errors, active connections, reconnects, matchmaking wait/success, match completion/abandonment, timer lateness, result/reward failures, ledger reconciliation, purchase verification, push outcomes, DB pool/locks, Redis latency, queue/outbox age, and client crashes. Correlate by request/match/event IDs without secrets or answer keys.

## Backup, recovery and operations

Automated encrypted PostgreSQL backups and PITR where available; monitor backup completion. Run scheduled restore drills and record achieved RPO/RTO. Redis is ephemeral coordination; rebuild queues/presence from durable state and expire stale tickets. Persist enough match state for recovery; define behavior for outages that cross a question deadline. Graceful shutdown drains connections and gives active commands bounded completion time.

Runbooks: matchmaking outage, WebSocket outage/reconnect storm, PostgreSQL failover/restore, Redis loss, purchase verification outage, season close/reward retry, duplicate reward/ledger mismatch, cheating wave, push provider outage, emergency configuration rollback.

## Scaling options

Scale API/WebSocket instances horizontally first. Add region-specific queues, a dedicated gateway/worker pools, PostgreSQL read replicas, leaderboard caches, event partitioning/archive, object CDN, or analytics warehouse only in response to measured needs. Preserve durable ownership and event contracts when extracting modules.
