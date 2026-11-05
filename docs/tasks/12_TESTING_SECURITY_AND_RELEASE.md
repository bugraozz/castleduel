# Testing, Security and Release Tasks

## RELEASE-01: Establish Layered Automated Test Strategy
**Priority:** Critical · **Phase:** 0-8 · **Component:** Testing · **Effort:** M · **Parallelism:** Test strategy owner coordinates all modules

**Description:** Add unit, integration, contract, E2E and test-data conventions across mobile/backend/admin. **Business Requirements:** Roadmap §4-5; Architecture §11. **Technical Requirements:** Deterministic fixtures, injected clocks, isolated DB/Redis, no production secrets/data. **Dependencies:** FOUND-03, FOUND-04.

**Implementation Steps:** 1. Select runners and coverage reporting. 2. Add DB/Redis test fixtures. 3. Add OpenAPI/WebSocket contract test job. 4. Define flaky-test handling and ownership. **Expected Files / Modules:** test packages, CI workflows, fixtures and test documentation. **API / Database Impact:** None.

**Acceptance Criteria:** Each layer can run locally and in CI; tests are isolated and repeatable; failures report actionable context. **Testing Requirements:** CI clean-run validation. **Completion Definition:** One representative unit, integration, contract and E2E test runs in CI.

## RELEASE-02: Test Concurrent Matchmaking and Assignment Invariants
**Priority:** Critical · **Phase:** 5-8 · **Component:** Testing / Backend · **Effort:** L · **Parallelism:** Independent QA after MATCH-03

**Description:** Prove two simultaneous queue entrants cannot create duplicate assignments. **Business Requirements:** PRD §8; Roadmap §4-5. **Technical Requirements:** Multiple workers/processes, repeated queue joins, cancellation race, blocked pair and database invariant checks. **Dependencies:** MATCH-01..MATCH-04, RELEASE-01.

**Implementation Steps:** 1. Generate concurrent queue load. 2. Race two workers for same tickets. 3. Race cancellation and reservation. 4. Assert one ticket/reservation/duel per participant. 5. Capture metrics and invariant failures. **Expected Files / Modules:** matchmaking stress suite. **API / Database Impact:** Exercises queue/reservation schema.

**Acceptance Criteria:** Across configured stress iterations, no player receives two active matches and every reservation reaches valid terminal state. **Testing Requirements:** T-11, T-12. **Completion Definition:** CI stress profile and staging load evidence pass.

## RELEASE-03: Test Answer Timing, Duplicates and Server Authority
**Priority:** Critical · **Phase:** 4-8 · **Component:** Testing / Security · **Effort:** L · **Parallelism:** Independent QA after BATTLE-04

**Description:** Verify correct behavior at timer boundaries and against forged/replayed answer commands. **Business Requirements:** PRD §8; GDD §4; Security requirements. **Technical Requirements:** Inject server time, submit near deadline, duplicate and malformed commands, forged duration/correctness, concurrent answers. **Dependencies:** BATTLE-03..BATTLE-05, QUES-05.

**Implementation Steps:** 1. Add deterministic clock test fixture. 2. Test just-before/at/after deadline. 3. Test duplicate and changed payload with same ID. 4. Submit client correctness/time fields. 5. Verify answer key unavailable pre-resolution. **Expected Files / Modules:** duel security/integration suite. **API / Database Impact:** REST and WebSocket answer commands.

**Acceptance Criteria:** Server receive time controls eligibility; first valid answer cannot be replaced; forged fields have no effect; one score event. **Testing Requirements:** T-02, T-03, T-06, T-07. **Completion Definition:** Adversarial answer suite passes on both transports.

## RELEASE-04: Test Disconnect, Reconnect and Server Recovery
**Priority:** Critical · **Phase:** 4-8 · **Component:** Testing / Realtime · **Effort:** L · **Parallelism:** Independent QA after state sync

**Description:** Verify active match restoration after network and process interruption. **Business Requirements:** PRD §8-9; Roadmap §4-5. **Technical Requirements:** Disconnect during question, resume after temporary interruption, sequence gap/replay, server restart and deadline crossing policy. **Dependencies:** BATTLE-07, BATTLE-10, MOB-05.

**Implementation Steps:** 1. Disconnect each participant at different phases. 2. Reconnect with last ack. 3. Force sequence gap and expired replay. 4. Restart API/worker during active round. 5. Verify timer, answer, result and client UI converge. **Expected Files / Modules:** chaos/recovery tests and runbook evidence. **API / Database Impact:** State sync, event log, durable deadline.

**Acceptance Criteria:** Clients recover authoritative state without duplicate effects; outage/deadline follows approved policy; active match is not silently discarded. **Testing Requirements:** T-04, T-05. **Completion Definition:** Staging network/process fault test passes.

## RELEASE-05: Test Atomic Rewards, Ledger and Claims
**Priority:** Critical · **Phase:** 6-8 · **Component:** Testing / Backend · **Effort:** L · **Parallelism:** Independent after ECON/BATTLE transaction freeze

**Description:** Verify result/reward transaction consistency under retries and injected failures. **Business Requirements:** PRD §4,8; DB §15; Economy §5. **Technical Requirements:** Interrupt after each write, retry same operation, compare ledger projection and inventory. **Dependencies:** BATTLE-08, ECON-01..ECON-06.

**Implementation Steps:** 1. Inject failure during finalization. 2. Retry finalize and claims. 3. Verify unique source key. 4. Reconcile wallet and inventory totals. 5. Test abandonment/farming reward policy. **Expected Files / Modules:** economy transaction/fault injection suite. **API / Database Impact:** Match rewards, mission, season and wallet claims.

**Acceptance Criteria:** Reward transaction is all-or-nothing or safely replayable; no duplicate grant; ledger and projection reconcile. **Testing Requirements:** T-09, T-15. **Completion Definition:** Failure-injection report has no invariant violation.

## RELEASE-06: Test Concurrent Castle Damage, Repair and Inventory Use
**Priority:** Critical · **Phase:** 6-8 · **Component:** Testing / Backend · **Effort:** L · **Parallelism:** Independent after CASTLE/ECON implementation

**Description:** Validate serialization of concurrent operations against shared player/castle assets. **Business Requirements:** PRD §8; GDD §7-9; Roadmap §4. **Technical Requirements:** Parallel attacks, repairs, builds, debits and item uses with stale versions/idempotency. **Dependencies:** CASTLE-03..CASTLE-07, ECON-01..ECON-04.

**Implementation Steps:** 1. Race two attacks on same castle. 2. Race repair and attack. 3. Race two upgrades/resource spends. 4. Race use of last item charge. 5. Assert durability, ledger and inventory invariants. **Expected Files / Modules:** castle/economy concurrency tests. **API / Database Impact:** Damage, repair, construction and inventory endpoints.

**Acceptance Criteria:** No negative balance, duplicated damage/consumption, lost update or durability above maximum; all results have audit events. **Testing Requirements:** T-08, T-10. **Completion Definition:** Stress and transaction tests pass on PostgreSQL.

## RELEASE-07: Execute Security and Anti-Cheat Test Plan
**Priority:** Critical · **Phase:** 4-8 · **Component:** Security / Testing · **Effort:** L · **Parallelism:** Independent security testing after contracts stabilize

**Description:** Validate identity, authorization, realtime, economy, purchase and admin threat controls. **Business Requirements:** PRD §4,8; Architecture security; Roadmap §4-5. **Technical Requirements:** Use authorized test environment/accounts; no production probing. **Dependencies:** BACK-02, BACK-05, BACK-08, ECON-09, ADMIN-01.

**Implementation Steps:** 1. Test token replay/revocation and object access. 2. Test WebSocket spoofing/subscription/rate limits. 3. Test duplicate economy/purchase operations. 4. Test admin escalation/audit. 5. Run dependency/secret/container scans. **Expected Files / Modules:** security test report and remediation issues. **API / Database Impact:** All protected APIs and admin routes.

**Acceptance Criteria:** No known Critical/High security issue remains at release gate; findings have owner and fix evidence; logs contain no secrets. **Testing Requirements:** Security checklist in `07_SECURITY_AND_ANTICHEAT.md`. **Completion Definition:** Security and engineering owners sign release report.

## RELEASE-08: Define Load Profile and Run Performance Tests
**Priority:** High · **Phase:** 7-8 · **Component:** Testing / Infrastructure · **Effort:** L · **Parallelism:** Can prepare harness while feature work continues

**Description:** Measure API, WebSocket, matchmaking and finalization against an agreed launch load. **Business Requirements:** PRD §9; Architecture §11; Roadmap §4. **Technical Requirements:** Define concurrency/region/network profile before interpreting thresholds; test p95 <300ms API and <250ms realtime under normal supported-region profile. **Dependencies:** BACK-07, MATCH-08, BATTLE-08.

**Implementation Steps:** 1. Agree player/concurrency/traffic model. 2. Build load harness. 3. Test queue bursts, concurrent duels, reconnect storm and leaderboard reads. 4. Capture DB/Redis saturation and event lag. 5. Tune or revise capacity plan. **Expected Files / Modules:** load tests, profile and results report. **API / Database Impact:** Production-like topology.

**Acceptance Criteria:** Report includes p95/p99, error rate and resource saturation; stated product targets are measured, not assumed. **Testing Requirements:** T-18. **Completion Definition:** Capacity baseline and scaling actions approved.

## RELEASE-09: Implement Production Monitoring and Alerting
**Priority:** High · **Phase:** 7-8 · **Component:** Infrastructure · **Effort:** M · **Parallelism:** Parallel with load test

**Description:** Provide dashboards and actionable alerts for gameplay, economy, purchases, mobile crashes and dependencies. **Business Requirements:** PRD §9-10; Roadmap §3-7. **Technical Requirements:** Correlated telemetry, privacy redaction, alert owners/runbooks and SLO definitions. **Dependencies:** BACK-07, MOB-09.

**Implementation Steps:** 1. Define service indicators and thresholds. 2. Build dashboards for API/realtime/MM/economy/purchase/client health. 3. Add alert routing/escalation. 4. Verify synthetic failures trigger alerts. **Expected Files / Modules:** dashboards, alert rules and on-call ownership. **API / Database Impact:** Telemetry only.

**Acceptance Criteria:** Reward duplication, ledger mismatch, queue stall, provider outage and crash-rate spikes are detected; every alert has an owner and runbook. **Testing Requirements:** Alert simulation and log redaction. **Completion Definition:** Staging game-day confirms alert delivery and response.

## RELEASE-10: Implement Backup, Restore and Disaster Recovery Drill
**Priority:** Critical · **Phase:** 8 · **Component:** Infrastructure / Database · **Effort:** L · **Parallelism:** Infrastructure owner

**Description:** Protect durable game data and verify recoverability. **Business Requirements:** PRD §9,11; Roadmap §5,7. **Technical Requirements:** Encrypted automated backups, PITR where available, documented RPO/RTO, restore access control and validation. **Dependencies:** BACK-04, infrastructure environment.

**Implementation Steps:** 1. Configure backup/PITR. 2. Monitor backup completion. 3. Restore into isolated environment. 4. Validate ledger/result/castle invariants. 5. Document failover and recovery procedure. **Expected Files / Modules:** backup config, restore scripts and runbook. **API / Database Impact:** Production PostgreSQL.

**Acceptance Criteria:** Restore completes within approved RTO; data loss stays within RPO; restored application passes invariant checks. **Testing Requirements:** T-18 recovery drill. **Completion Definition:** Dated restore report and operational sign-off.

## RELEASE-11: Complete Device, Accessibility and UX QA
**Priority:** High · **Phase:** 8 · **Component:** Mobile / Testing · **Effort:** L · **Parallelism:** Parallel after core flows stabilize

**Description:** Validate representative iOS/Android devices, network conditions, accessibility and reduced motion. **Business Requirements:** PRD §9,11; UI/UX §2,7-8. **Technical Requirements:** Supported OS/screen matrix, large text, screen reader, low-memory, safe areas, localization, offline/error states. **Dependencies:** MOB-02, MOB-07, MOB-08.

**Implementation Steps:** 1. Define device matrix. 2. Test all required flows and sizes. 3. Run accessibility audit. 4. Test poor network/background and reduced motion. 5. Resolve blockers and record known limitations. **Expected Files / Modules:** QA matrix and defects. **API / Database Impact:** None.

**Acceptance Criteria:** No blocking usability/accessibility defects in critical flows; active match and purchases recover safely. **Testing Requirements:** T-20 and device checklist. **Completion Definition:** Product/design/accessibility sign-off.

## RELEASE-12: Prepare App Store and Google Play Compliance
**Priority:** Critical · **Phase:** 8 · **Component:** Release / Legal · **Effort:** L · **Parallelism:** Begins early; centralized release owner

**Description:** Prepare store metadata, legal/privacy disclosures and platform-specific review requirements. **Business Requirements:** PRD §5,8; Economy §14; Roadmap §6. **Technical Requirements:** Current policies, age rating, privacy/data safety, support, terms, deletion flow, IAP metadata and review account. **Dependencies:** BACK-03, ECON-09, RELEASE-11, FOUND-02.

**Implementation Steps:** 1. Confirm territories and age/child obligations. 2. Prepare metadata/assets and policy URLs. 3. Complete privacy/data disclosures. 4. Test account deletion and purchase restoration. 5. Prepare reviewer instructions and sandbox accounts. **Expected Files / Modules:** store submission package and compliance evidence. **API / Database Impact:** account deletion and purchase operations.

**Acceptance Criteria:** Store forms and current policy checks complete; all disclosed data collection matches implementation; review path works. **Testing Requirements:** Store sandbox/review-account and privacy-flow tests. **Completion Definition:** Release owner records submission-readiness approval.

## RELEASE-13: Implement Release Pipeline and Staged Rollout
**Priority:** High · **Phase:** 8 · **Component:** Infrastructure · **Effort:** L · **Parallelism:** Can prepare after FOUND-03

**Description:** Deliver signed, traceable builds through staging, beta and production with safe rollback. **Business Requirements:** Roadmap §5-6; Architecture §8. **Technical Requirements:** Artifact provenance, approval gates, backward-compatible migrations/contracts, client minimum-version handling and staged rollout. **Dependencies:** FOUND-03, BACK-07, RELEASE-01.

**Implementation Steps:** 1. Build/version artifacts. 2. Deploy staging and run smoke gates. 3. Configure signing/secrets. 4. Add approval and staged promotion. 5. Document rollback/config rollback and maintenance mode. **Expected Files / Modules:** CI/CD workflows and deployment manifests. **API / Database Impact:** Environment rollout.

**Acceptance Criteria:** Production deploy is traceable to commit; failed smoke gate blocks promotion; rollback does not require unsafe destructive migration. **Testing Requirements:** Staging deploy, rollback and min-version test. **Completion Definition:** Dry-run promotion and rollback evidence.

## RELEASE-14: Complete Operational Runbooks and Launch Readiness Review
**Priority:** Critical · **Phase:** 8 · **Component:** Infrastructure / Product · **Effort:** M · **Parallelism:** Central release coordination

**Description:** Prepare team and systems for incidents and verify all release gates. **Business Requirements:** Roadmap §5,7,9; PRD §11. **Technical Requirements:** Named owners, severity/escalation, runbooks for every listed outage and operational evidence. **Dependencies:** RELEASE-07..RELEASE-13, ADMIN-08.

**Implementation Steps:** 1. Document matchmaking/WebSocket/DB/Redis/purchase/season/push/cheating/reward incident response. 2. Run tabletop/game day. 3. Verify monitoring, backups, support and moderation. 4. Review unresolved Critical/High defects. 5. Sign final go/no-go checklist. **Expected Files / Modules:** runbooks, release checklist and evidence bundle. **API / Database Impact:** Operational only.

**Acceptance Criteria:** Each incident has owner, detection, mitigation, recovery and communication path; no open Critical/High security issues; core gameplay and transaction gates pass. **Testing Requirements:** Tabletop and live staging recovery exercise. **Completion Definition:** Recorded release readiness approval by product, engineering, security and operations.

## Test scenario catalog

| Test ID | Scenario | Expected behavior | Minimum verification |
|---|---|---|---|
| T-01 | Register/login and profile retrieval on two accounts | Unique normalized identities; safe public profile | API + mobile E2E |
| T-02 | Two players submit answers nearly simultaneously | One answer per player; one deterministic round resolution | Concurrency integration |
| T-03 | Late answer, duplicate answer, replayed command | Stable rejection/original result; no extra score | API/WebSocket security |
| T-04 | Player disconnects during question then reconnects | Authoritative state/deadline restored; no local-clock advantage | Realtime E2E |
| T-05 | Server process restarts during active round | Durable state recovers or follows approved cancellation policy | Fault injection |
| T-06 | Client submits fabricated correctness, score or response time | Server ignores forged values and calculates independently | Security integration |
| T-07 | Client uses manipulated local clock or submits after timeout | Server deadline remains authoritative | Deterministic clock test |
| T-08 | Player uses an item they do not own / repeats activation | Rejected; no charge consumed or duplicate effect | Inventory/combat integration |
| T-09 | Match finishes and reward transaction is interrupted/retried | All-or-nothing or safe replay; exactly one reward | DB failure injection |
| T-10 | Two attacks/repairs modify one castle simultaneously | No lost update, over-repair, duplicated damage or double spend | PostgreSQL concurrency |
| T-11 | Two players enter matchmaking simultaneously | One valid reservation; no duplicate pairing | Multi-worker load test |
| T-12 | One player attempts multiple queue/match assignments | Unique active ticket/reservation/duel enforced | DB constraint stress |
| T-13 | Blocked users or expired direct invitations | Never paired/invited into a duel; expiry is final | Social/matchmaking integration |
| T-14 | Season closes while results/rewards are pending | Reconcile, archive and grant exactly once; castle persists | Job replay integration |
| T-15 | Duplicate mission/achievement/reward claim events | Progress/grant once; ineligible claim denied | Event/claim integration |
| T-16 | Forged, duplicated, refunded or restored purchase | Verify server-side; fulfill/revoke exactly once | Store sandbox |
| T-17 | Admin role attempts unauthorized moderation/economy/config action | Denied and audited; no hidden mutation | RBAC integration |
| T-18 | Queue burst, concurrent duels, reconnect storm, restore drill | Measured SLO/load and RPO/RTO targets met or capacity actions recorded | Load/ops exercise |
| T-19 | Duplicate push/outbox delivery and changed preferences | One inbox item; preference respected; retry safe | Worker/provider tests |
| T-20 | Account deletion, privacy, accessibility and offline states | Data policy enforced; critical UX remains recoverable and accessible | Privacy/device E2E |
