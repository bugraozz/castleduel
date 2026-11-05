# Backend Foundation Tasks

## BACK-01: Initialize Go API and Worker Applications
**Implementation status:** In progress. API and worker startup, typed configuration, bounded PostgreSQL/Redis clients, structured logs, health/readiness endpoints, and shutdown handling are implemented. PostgreSQL/Redis integration tests and a Linux SIGTERM smoke test verified readiness and graceful shutdown. `pgx`, `x/crypto`, and Go were upgraded after vulnerability scanning; reachable-code scan is clean. Migration infrastructure remains blocked on FOUND-04/FOUND-05, and the GitHub Actions run remains unverified.
**Priority:** Critical · **Phase:** 1 Identity and app shell · **Component:** Backend · **Effort:** M · **Parallelism:** Central backend owner

**Description:** Create the deployable API and worker entry points with clean module boundaries. **Business Requirements:** PRD §9; Architecture §2-4. **Technical Requirements:** Typed configuration, graceful shutdown, health/readiness endpoints, dependency injection and transport/domain separation. **Dependencies:** FOUND-01, FOUND-03.

**Implementation Steps:** 1. Create API/worker commands. 2. Add health endpoints and structured startup/shutdown. 3. Wire PostgreSQL/Redis clients with bounded pools/timeouts. 4. Add dependency injection and test seams. **Expected Files / Modules:** `cmd/api`, `cmd/worker`, `internal/platform`, module package skeletons. **API / Database Impact:** Health endpoints only; initial migration runner wiring.

**Acceptance Criteria:** Both processes start with validated config, report readiness only after required dependencies are healthy, and shut down cleanly. **Testing Requirements:** Unit config tests, PostgreSQL/Redis integration startup, graceful shutdown test. **Completion Definition:** CI starts both binaries and checks health.

## BACK-02: Implement Authentication and Session Rotation
**Priority:** Critical · **Phase:** 1 Identity and app shell · **Component:** Backend · **Effort:** L · **Parallelism:** Central security review

**Description:** Implement account registration/login, secure sessions and refresh rotation. **Business Requirements:** PRD §5,8; API §2. **Technical Requirements:** Password hashing if password auth is used, normalized unique identity, short-lived access token, hashed refresh token, rotation/revocation, generic recovery response and rate limits. **Dependencies:** BACK-01, FOUND-02, FOUND-05.

**Implementation Steps:** 1. Add users/identity/session migrations. 2. Implement register/login/logout/refresh handlers. 3. Enforce uniqueness and credential throttling. 4. Rotate and revoke refresh credentials transactionally. 5. Add password recovery adapter/workflow if launch identity supports it. **Expected Files / Modules:** `internal/auth`, auth migrations, handlers and tests. **API / Database Impact:** `/auth/register`, `/login`, `/refresh`, `/logout`, `/password/*`; users, identities, sessions.

**Acceptance Criteria:** Rotated token invalidates predecessor; logout revokes session; plaintext passwords/tokens are never persisted/logged; enumeration-safe failures and throttling are tested. **Testing Requirements:** Unit and DB integration tests for concurrent refresh, replay, revoked/expired sessions, duplicate username. **Completion Definition:** Security-reviewed implementation passes auth integration suite.

## BACK-03: Implement Account Deletion and Privacy-Safe Profiles
**Priority:** High · **Phase:** 1 Identity and app shell · **Component:** Backend · **Effort:** M · **Parallelism:** Can parallelize after BACK-02 schema

**Description:** Provide permitted profile access/update and account deletion workflow. **Business Requirements:** PRD §8; UI/UX §5.18,5.21; Roadmap §5. **Technical Requirements:** Public response allowlist; deletion/anonymization policy with legal retention exception, confirmation and audit; prevent active-match deletion from corrupting result. **Dependencies:** BACK-02, PROG-01, FOUND-02.

**Implementation Steps:** 1. Implement self/public profile reads and allowed patch. 2. Validate privacy settings. 3. Implement deletion request, session revocation and scheduled anonymization/retention state. 4. Define active match and purchase retention behavior. **Expected Files / Modules:** `internal/player`, `internal/auth/deletion`, migrations and policy docs. **API / Database Impact:** `GET/PATCH /players/me`, `GET /players/{id}/public`, `DELETE /account`.

**Acceptance Criteria:** Public profile excludes private identity/payment fields; deletion request is confirmed and auditable; credentials revoke; retained records follow approved policy. **Testing Requirements:** Authorization, privacy allowlist, active-match/deletion and retention integration tests. **Completion Definition:** Privacy review and end-to-end account deletion evidence.

## BACK-04: Define Database Repositories and Transaction Boundaries
**Priority:** Critical · **Phase:** 1 Identity and app shell · **Component:** Database · **Effort:** L · **Parallelism:** Central backend/database coordination

**Description:** Establish PostgreSQL access and transaction patterns used across domain modules. **Business Requirements:** Database §1,14-16; Architecture §8. **Technical Requirements:** Migrations, repository boundaries, row locking/version checks, transaction runner, constraints and UTC/integer conventions. **Dependencies:** BACK-01, FOUND-05.

**Implementation Steps:** 1. Implement connection pool and migration startup policy. 2. Add transaction abstraction with rollback tests. 3. Add typed query/repository conventions. 4. Add query timeout and pool telemetry. 5. Prohibit cross-module repository access. **Expected Files / Modules:** `internal/platform/postgres`, migrations, transaction helpers. **API / Database Impact:** Foundation for all schema and transactional use cases.

**Acceptance Criteria:** Transaction rollback preserves all-or-nothing behavior; migrations run from empty DB; stale version and constraint failures map to domain conflict. **Testing Requirements:** PostgreSQL integration tests for rollback, locking and migration. **Completion Definition:** Pattern adopted by auth/profile mutation and documented.

## BACK-05: Add API Middleware, Validation and Error Mapping
**Priority:** High · **Phase:** 1 Identity and app shell · **Component:** Backend · **Effort:** M · **Parallelism:** After FOUND-04

**Description:** Apply consistent request IDs, auth, validation, size bounds, rate limits and error responses. **Business Requirements:** API §1,13-14; Security requirements. **Technical Requirements:** Stable envelope; no internal details; per-route limits; schema validation before use case. **Dependencies:** BACK-01, FOUND-04.

**Implementation Steps:** 1. Add request correlation and structured logs. 2. Add auth middleware and object-scoped authorization helpers. 3. Add body size/timeouts and validation. 4. Map domain errors to documented codes. 5. Add route rate-limit configuration and tests. **Expected Files / Modules:** `internal/api/middleware`, error package, route tests. **API / Database Impact:** Consistent behavior across `/api/v1`.

**Acceptance Criteria:** Every error has stable code and request ID; no password/token/internal SQL is logged or returned; oversized/rate-limited requests are rejected. **Testing Requirements:** Middleware contract, invalid payload, authz, rate-limit and log redaction tests. **Completion Definition:** Contract tests pass on representative endpoints.

## BACK-06: Add Durable Idempotency and Outbox Infrastructure
**Priority:** Critical · **Phase:** 1-2 · **Component:** Backend / Database · **Effort:** L · **Parallelism:** Central transaction owner

**Description:** Provide reusable deduplication and post-commit event delivery primitives. **Business Requirements:** Architecture §8; DB §13-15; PRD §4,8. **Technical Requirements:** Actor/operation/key scope, request hash, stored outcome, changed-payload rejection, transactional outbox and deduplicating consumer. **Dependencies:** BACK-04, FOUND-04.

**Implementation Steps:** 1. Add idempotency/outbox migrations. 2. Implement transaction-scoped reserve/complete. 3. Persist outcome and reject key reuse with a different payload. 4. Implement worker claim/retry/backoff and event dedupe. 5. Add retention and failure alert policy. **Expected Files / Modules:** `internal/platform/idempotency`, `outbox`, worker consumers. **API / Database Impact:** `idempotency_records`, `outbox_events`; mutating request headers/commands.

**Acceptance Criteria:** Concurrent identical requests apply one effect and return same result; different payload with same key is rejected; outbox event is delivered at least once. **Testing Requirements:** Concurrency, crash-after-commit, retry, duplicate consumer and dead-letter tests. **Completion Definition:** At least one representative mutation uses the shared primitive.

## BACK-07: Add Observability and Operational Diagnostics
**Priority:** High · **Phase:** 1-2 · **Component:** Infrastructure · **Effort:** M · **Parallelism:** Parallel after BACK-01

**Description:** Instrument API, DB, Redis, worker, matchmaking and gameplay-critical operations. **Business Requirements:** PRD §9-10; Roadmap §4,7. **Technical Requirements:** Correlation IDs, metrics, traces, redaction, health/readiness; never log tokens, answer keys or unnecessary PII. **Dependencies:** BACK-01, BACK-05.

**Implementation Steps:** 1. Add structured logging and trace propagation. 2. Add API latency/error and dependency metrics. 3. Add job/outbox age and retry metrics. 4. Document dashboards/alerts. **Expected Files / Modules:** telemetry package, dashboards and runbook references. **API / Database Impact:** None.

**Acceptance Criteria:** A request can be traced across handler/repository/worker; sensitive fields are redacted; health distinguishes liveness/readiness. **Testing Requirements:** Log redaction and metric emission tests. **Completion Definition:** Staging telemetry receives a synthetic request and alert thresholds are documented.

## BACK-08: Implement Authenticated Realtime Gateway
**Priority:** High · **Phase:** 4 Authoritative duel · **Component:** Backend / Realtime · **Effort:** L · **Parallelism:** Contract-coordinated with BATTLE

**Description:** Establish secure WebSocket connection, subscriptions, envelopes and multi-instance delivery adapter. **Business Requirements:** Architecture §6-7; API §12. **Technical Requirements:** TLS at edge, origin/version checks, token auth, identity-bound channels, message limits, heartbeat, event schema validation and graceful drain. **Dependencies:** BACK-02, BACK-05, FOUND-04.

**Implementation Steps:** 1. Implement handshake/auth. 2. Restrict subscriptions by principal and active membership. 3. Add typed inbound/outbound envelopes and rate limits. 4. Add Redis-backed cross-instance delivery abstraction. 5. Add heartbeat, close/revocation and shutdown handling. **Expected Files / Modules:** `internal/realtime`, transport adapter, schemas. **API / Database Impact:** `/api/v1/realtime`; event auth and subscription rules.

**Acceptance Criteria:** Unauthorized subscriptions/commands are rejected; expired tokens close; payload and rate limits work; multi-instance delivery is tested. **Testing Requirements:** WebSocket integration tests for auth, spoofing, oversized payload, reconnect, fan-out and shutdown. **Completion Definition:** Staging two-client connection and security test evidence.

## BACK-09: Define Background Job Reliability and Scheduling
**Priority:** Medium · **Phase:** 2-8 · **Component:** Backend · **Effort:** M · **Parallelism:** Parallel after BACK-06

**Description:** Standardize scheduled/retryable jobs for expiry, reconciliation, notifications and season operations. **Business Requirements:** Architecture §8; Roadmap §2-7. **Technical Requirements:** Idempotent handlers, bounded retries, lease/claim semantics, dead-letter alerts and safe replay. **Dependencies:** BACK-06, BACK-07.

**Implementation Steps:** 1. Define job envelope, schedule and lease policy. 2. Add retry/backoff and dead-letter state. 3. Add job health metrics. 4. Provide operator replay controls with audit for sensitive jobs. **Expected Files / Modules:** Worker scheduler and job packages. **API / Database Impact:** Job/outbox status fields and operation audit.

**Acceptance Criteria:** Crash/retry cannot duplicate side effects; stuck jobs alert; operator replay preserves idempotency. **Testing Requirements:** Lease expiry, repeated execution, backoff and failed-job recovery tests. **Completion Definition:** One scheduled expiry and one outbox job use the common pattern.
