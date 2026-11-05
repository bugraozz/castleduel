# Matchmaking and Invitation Tasks

## MATCH-01: Implement Matchmaking Ticket Lifecycle
**Priority:** Critical · **Phase:** 5 Ranked and friend entry · **Component:** Backend / Database · **Effort:** M · **Parallelism:** Central with BATTLE

**Description:** Create, read and cancel ranked queue tickets with eligibility checks. **Business Requirements:** PRD §5,7-8; API §5; Roadmap Epic 5. **Technical Requirements:** One active ticket per player; trusted server rating/level/castle power/region; queue cancellation race-safe. **Dependencies:** BACK-04, PROG-03, BACK-06.

**Implementation Steps:** 1. Add ticket schema and partial unique constraint. 2. Implement queue eligibility and join. 3. Implement status lookup. 4. Implement cancellation transaction. 5. Add stale ticket expiration. **Expected Files / Modules:** matchmaking module, migrations and handlers. **API / Database Impact:** `POST/GET/DELETE /matchmaking/queue`; matchmaking_tickets.

**Acceptance Criteria:** Duplicate join cannot create a second ticket; client-supplied rating is ignored; cancel racing with reservation has one deterministic winner. **Testing Requirements:** Concurrent join/cancel and expiration tests. **Completion Definition:** Ticket state integration suite passes.

## MATCH-02: Implement Rating-Bounded Candidate Selection
**Priority:** High · **Phase:** 5 · **Component:** Backend · **Effort:** L · **Parallelism:** Algorithm can parallelize; policy centrally approved

**Description:** Find fair opponent candidates with progressive search expansion. **Business Requirements:** PRD §4,8; GDD §11; DB §6. **Technical Requirements:** MMR primary, castle power secondary, level/region/availability constraints, maximum rating gap, blocked-user exclusion, queue-age expansion. **Dependencies:** MATCH-01, PROG-03, SOC-02, FOUND-02.

**Implementation Steps:** 1. Define eligibility and fairness score. 2. Implement expansion schedule/config. 3. Filter blocks, active matches and region constraints. 4. Record wait/gap outcome. 5. Add policy telemetry. **Expected Files / Modules:** candidate selector/config and query indexes. **API / Database Impact:** ticket selection and matchmaking_history.

**Acceptance Criteria:** Search never exceeds approved maximum gap; eligible oldest compatible candidate is considered by deterministic policy; blocked/active players are excluded. **Testing Requirements:** Table-driven boundaries, queue aging, region and block tests. **Completion Definition:** Simulated distribution review approved before ranked activation.

## MATCH-03: Implement Atomic Pair Reservation
**Priority:** Critical · **Phase:** 5 · **Component:** Backend / Database · **Effort:** L · **Parallelism:** Centralized concurrency-critical work

**Description:** Atomically reserve two candidate tickets and both players. **Business Requirements:** PRD §8; DB §6,14; Security requirements. **Technical Requirements:** DB constraints/locks or Redis script plus durable claim; reservation lease; one active reservation/match per user; recovery after process/Redis failure. **Dependencies:** MATCH-01, MATCH-02, BACK-04.

**Implementation Steps:** 1. Design reservation state and uniqueness constraints. 2. Lock/claim both tickets atomically. 3. Persist reservation and lease. 4. Handle retry and expired lease recovery. 5. Emit match-found event after commit. **Expected Files / Modules:** reservation repository/service, migration and race tests. **API / Database Impact:** matchmaking_matches and active participation constraint.

**Acceptance Criteria:** In simultaneous queue stress, neither user is reserved twice and no ticket pair creates two matches; crash recovery releases/continues per policy. **Testing Requirements:** High-concurrency multi-worker test and Redis-loss integration. **Completion Definition:** Load test evidence and invariants reviewed by database owner.

## MATCH-04: Implement Match Confirmation, Decline and Cancellation
**Priority:** High · **Phase:** 5 · **Component:** Backend · **Effort:** M · **Parallelism:** After MATCH-03

**Description:** Let both reserved players confirm/decline before a server deadline. **Business Requirements:** PRD §7; UI/UX §5.6-5.7; API §5. **Technical Requirements:** Idempotent confirmation, expiry, explicit return-to-queue policy, no partial active duel. **Dependencies:** MATCH-03, BACK-06, BACK-08.

**Implementation Steps:** 1. Persist per-player confirmation. 2. Add confirm/decline operations. 3. Implement deadline worker and cleanup. 4. Apply approved eligibility/queue return rule. 5. Notify both participants. **Expected Files / Modules:** reservation state machine and handlers. **API / Database Impact:** confirm/decline routes and realtime events.

**Acceptance Criteria:** One confirmation cannot start match; both are required; late confirmation is rejected; repeated operation returns same outcome. **Testing Requirements:** deadline boundary, duplicate, decline/cancel race and notification tests. **Completion Definition:** Full confirmation flow passes two-client integration test.

## MATCH-05: Create Duel from Reserved Pair
**Priority:** Critical · **Phase:** 5 · **Component:** Backend · **Effort:** M · **Parallelism:** Contract-coordinated with BATTLE

**Description:** Convert confirmed reservation to one duel using the shared duel engine. **Business Requirements:** PRD §7-8; Roadmap Epic 5. **Technical Requirements:** Atomic one-time creation, mode/rules/loadout snapshots, participants exactly two, clear rollback/retry. **Dependencies:** MATCH-04, BATTLE-01, QUES-03.

**Implementation Steps:** 1. Implement transaction that creates duel/players/questions. 2. Mark reservation consumed. 3. Publish ready/start event after commit. 4. Recover duplicate callback using reservation ID idempotency. **Expected Files / Modules:** matchmaking-to-duel application service. **API / Database Impact:** duels, duel_players, duel_questions and match association.

**Acceptance Criteria:** Replayed confirmation cannot create another duel; ranked match uses trusted snapshots; both clients join same duel state. **Testing Requirements:** Transaction failure injection, duplicate confirmation and exactly-two participant tests. **Completion Definition:** Ranked path creates one playable duel in staging.

## MATCH-06: Implement Queue Status, Fair Search Feedback and Cancellation UX Contract
**Priority:** Medium · **Phase:** 5 · **Component:** Mobile / Backend · **Effort:** M · **Parallelism:** UI can parallelize after event schema

**Description:** Expose safe queue state and show honest search progress. **Business Requirements:** UI/UX §5.5-5.6. **Technical Requirements:** Server queue stage/elapsed duration; no fabricated ETA; cancellation confirms server state. **Dependencies:** MATCH-01, MATCH-04, MOB-05.

**Implementation Steps:** 1. Add `searching` updates or polling state. 2. Implement status/reconnect recovery. 3. Build elapsed timer and cancel UI. 4. Show match-found confirmation deadline and retry reason. **Expected Files / Modules:** matchmaking routes/events and mobile feature. **API / Database Impact:** queue status response and queue events.

**Acceptance Criteria:** Client recovers queue after reconnect; cancel result reflects server; no unsupported exact ETA is displayed. **Testing Requirements:** Queue background/reconnect/cancel E2E and event contract test. **Completion Definition:** UX review and staging queue flow pass.

## MATCH-07: Implement Direct Friend Invitations
**Priority:** Critical · **Phase:** 5 · **Component:** Backend / Social · **Effort:** L · **Parallelism:** Parallel with ranked selector after shared duel API

**Description:** Send, accept, reject, expire and rate-limit direct player challenges. **Business Requirements:** PRD §2,5,7; API §6; UI/UX §5.17. **Technical Requirements:** Friend/block/privacy eligibility, expiry, invitation idempotency, unranked rules and anti-farm reward configuration. **Dependencies:** SOC-01, SOC-02, MATCH-05, BACK-06.

**Implementation Steps:** 1. Add invitation schema/state machine. 2. Implement sender/recipient authorization and limits. 3. Implement accept/decline/expiry. 4. Create friendly duel through shared engine. 5. Persist match history and apply approved reward policy. **Expected Files / Modules:** social invitation service, handlers and tests. **API / Database Impact:** `/duel-invitations*`; duel_invitations and duel association.

**Acceptance Criteria:** Accepted invitation starts one friendly duel; expired/blocked/ineligible invite cannot be accepted; friendly match does not affect MMR by default. **Testing Requirements:** duplicate acceptance, expiry, block race and reward farming tests. **Completion Definition:** Two friends complete direct challenge on independent devices.

## MATCH-08: Add Matchmaking Recovery, Metrics and Abuse Controls
**Priority:** High · **Phase:** 5-8 · **Component:** Backend / Testing · **Effort:** M · **Parallelism:** After MATCH-03

**Description:** Handle stale tickets/reservations, duplicate queues and matchmaking abuse observably. **Business Requirements:** PRD §8; Economy §11; Roadmap §4-5. **Technical Requirements:** Cleanup jobs, per-account/device limits, blocked-pair guarantees, queue wait/gap metrics and replayable history. **Dependencies:** MATCH-03, BACK-07, BACK-09.

**Implementation Steps:** 1. Add reconciliation/expiry job. 2. Add duplicate/abuse limits. 3. Instrument queue depth/wait/failures. 4. Add operator diagnostics and alert thresholds. **Expected Files / Modules:** matchmaking jobs, metrics, admin diagnostics. **API / Database Impact:** history and reservation timestamps/status.

**Acceptance Criteria:** Orphan reservations are resolved within configured lease; abuse limits return stable errors; alerts detect queue stalls. **Testing Requirements:** Process restart, stale ticket, abuse-limit and blocked candidate tests. **Completion Definition:** Runbook and staging fault-injection evidence.
