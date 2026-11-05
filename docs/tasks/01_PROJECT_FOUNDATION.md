# Project Foundation Tasks

## FOUND-01: Establish Repository Structure and Toolchain
**Implementation status:** In progress. The root Git repository is initialized, the mobile scaffold's initial commit is preserved on `mobile-history`, and a root clone is being validated. Node/npm/Go and PostgreSQL/Redis image versions are pinned. The canonical remote was not supplied, so `castleduello.local/backend` remains provisional and no origin is configured.
**Priority:** Critical · **Phase:** 0 Foundation · **Component:** Infrastructure · **Effort:** M · **Parallelism:** Central coordination

**Description:** Inspect the existing repository, retain the eight source documents, decide a mobile/backend/admin workspace structure, and pin supported toolchains. **Business Requirements:** PRD §5, Roadmap §2 Epic 1. **Technical Requirements:** Record the modular-monolith decision and package boundaries; do not introduce app features. **Dependencies:** None.

**Implementation Steps:** 1. Inventory files and existing conventions. 2. Choose monorepo layout and package ownership. 3. Pin Go, Node, Expo SDK, package manager, Docker, PostgreSQL and Redis versions. 4. Document setup and compatibility. **Expected Files / Modules:** repository root, `apps/mobile`, `apps/admin` (if approved), `services/api`, `services/worker`, shared contracts, toolchain files, README. **API / Database Impact:** None; layout reserves migration and contract locations.

**Acceptance Criteria:** A new developer can identify each deployable and its owner; all tool versions are pinned; original source docs remain unchanged; local setup steps are deterministic. **Testing Requirements:** Clean-clone toolchain/setup smoke check. **Completion Definition:** Reviewed structure decision and successful setup evidence are recorded.

## FOUND-02: Approve Product and Architecture Decisions
**Implementation status:** Blocked pending product/legal/engineering owners to approve the release-blocking choices below; gameplay and compliance defaults have not been assumed.
**Priority:** Critical · **Phase:** 0 Foundation · **Component:** Backend · **Effort:** M · **Parallelism:** Centralized product/architecture coordination

**Description:** Resolve or explicitly schedule decisions that change ranked gameplay, data contracts, or compliance. **Business Requirements:** PRD §12; GDD §2,4-9; Economy §3,10-14. **Technical Requirements:** Version decisions in ADRs/configuration; no undocumented defaults. **Dependencies:** FOUND-01.

**Implementation Steps:** 1. Decide score/speed/tie/draw policy and ranked draw rating. 2. Decide both-player versus winner-only castle damage, damage floor/protection, ability roster and ranked limits. 3. Decide friendly rewards, queue rating-gap expansion/regions, season policy, and loadout snapshot. 4. Decide answer transport/latency policy, auth provider/token format, data retention/age/launch regions, hosting and load profile. 5. Record owner and deadline for unresolved items. **Expected Files / Modules:** `docs/decisions/` ADRs and configuration ownership notes. **API / Database Impact:** Defines rules snapshots, config schema and contract error/policy values.

**Acceptance Criteria:** Every listed release-blocking choice has an approved value or a tracked decision gate; active matches can be tied to a rules version. **Testing Requirements:** Decision-to-contract review and ranked policy examples. **Completion Definition:** Product and engineering owners sign off the ADR index; unresolved legal questions block store submission.

## FOUND-03: Create Local Development and CI Baseline
**Implementation status:** In progress. Digest-pinned Compose services, CI dependency/vulnerability, license inventory, and Gitleaks scans are configured. Local Compose and actionlint checks pass; Gitleaks found no leaks and the Go reachable-code scan is clean. npm still reports 12 moderate `glob`/`minimatch` advisories; its suggested fix downgrades Expo SDK 58 to 53, so it was not applied. License inventories report first-party `UNLICENSED` metadata; no legal allowlist was invented. A GitHub Actions run and legal license policy approval remain outstanding.
**Priority:** Critical · **Phase:** 0 Foundation · **Component:** Infrastructure · **Effort:** L · **Parallelism:** Can parallelize with design after FOUND-01

**Description:** Make local dependencies and quality checks reproducible. **Business Requirements:** Roadmap Epic 1; Architecture §13. **Technical Requirements:** Docker Compose for PostgreSQL/Redis, health checks, non-production config, migration/test commands, CI formatting/lint/static/unit checks. **Dependencies:** FOUND-01.

**Implementation Steps:** 1. Add local service definitions and safe example env file. 2. Add API/worker/mobile build and test commands. 3. Add CI jobs and cache only non-secret dependencies. 4. Add license/dependency and secret scanning. 5. Document startup, shutdown, migrations and troubleshooting. **Expected Files / Modules:** Compose file, CI workflows, `.env.example`, scripts and contributor docs. **API / Database Impact:** Creates local DB/Redis only; no production schema yet.

**Acceptance Criteria:** Clean clone can start dependencies, run checks and stop without production credentials; CI fails on formatting/type/test failures. **Testing Requirements:** CI on clean runner; container health and migration smoke test. **Completion Definition:** Passing CI run and reproducible local setup instructions.

## FOUND-04: Define Shared API, Event and Error Conventions
**Priority:** High · **Phase:** 0 Foundation · **Component:** Backend · **Effort:** M · **Parallelism:** Central contract owner

**Description:** Establish versioned REST/WebSocket conventions before clients and handlers diverge. **Business Requirements:** API §1,12-14; Architecture §7-8. **Technical Requirements:** OpenAPI, event envelope/schema version, UTC timestamp, UUID, error envelope, pagination, request IDs, idempotency and compatibility policy. **Dependencies:** FOUND-02, FOUND-03.

**Implementation Steps:** 1. Define auth notation and token handling. 2. Define stable error codes and retryability. 3. Specify Idempotency-Key/command ID scope and replay behavior. 4. Define sequence, ack and snapshot semantics. 5. Add contract linting and version rules. **Expected Files / Modules:** OpenAPI base, JSON event schemas, shared typed client generation and contract CI. **API / Database Impact:** Establishes contract and future idempotency storage requirement.

**Acceptance Criteria:** Example REST mutation and WebSocket command validate; duplicate/changed-payload semantics are explicit; generated client is type-safe. **Testing Requirements:** Schema validation and compatibility/lint tests. **Completion Definition:** Contract templates accepted by mobile/backend owners.

## FOUND-05: Establish Migration, Configuration and Secret Practices
**Priority:** High · **Phase:** 0 Foundation · **Component:** Database / Infrastructure · **Effort:** M · **Parallelism:** Coordinate with backend foundation

**Description:** Define reproducible schema change and per-environment configuration lifecycle. **Business Requirements:** Database §1,16; Architecture §8-13; Roadmap §7. **Technical Requirements:** Versioned migrations, expand/contract deploy policy, typed config validation, secret injection and rotation. **Dependencies:** FOUND-03, FOUND-04.

**Implementation Steps:** 1. Select migration tool. 2. Add migration checks for empty DB and rollback safety. 3. Define config layering and immutable game config versions. 4. Add secret scanning and environment separation. 5. Document safe local data reset. **Expected Files / Modules:** migration runner, config package, deployment templates and secret handling guide. **API / Database Impact:** None beyond migration framework.

**Acceptance Criteria:** CI validates migration ordering; no secret is committed; production config cannot silently default to development values. **Testing Requirements:** Empty-schema migration, invalid config startup failure, secret scan. **Completion Definition:** Migration/config patterns demonstrated in a test migration and documented.
