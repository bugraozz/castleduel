# Real-Time Battle System Tasks

## BATTLE-01: Define Durable Duel State Machine
**Priority:** Critical · **Phase:** 4 Authoritative duel · **Component:** Backend / Database · **Effort:** L · **Parallelism:** Centralized owner

**Description:** Model duel creation, preparation, active rounds, resolution, completion, cancellation and abandonment. **Business Requirements:** PRD §8; Architecture §6; DB §5. **Technical Requirements:** Explicit allowed transitions, exactly two eligible participants, rules/config/loadout snapshots, durable status and version. **Dependencies:** BACK-04, QUES-01, FOUND-02.

**Implementation Steps:** 1. Define state diagram and transition guards. 2. Add duel/participant/result schema. 3. Implement transition use cases. 4. Enforce active-user uniqueness and state version. **Expected Files / Modules:** `internal/duel`, duel migrations and transition tests. **API / Database Impact:** duel create/state routes, duels/duel_players.

**Acceptance Criteria:** Invalid transition is rejected; both participants are authorized; restart can determine current state from durable records. **Testing Requirements:** Transition table, duplicate start and active-match uniqueness tests. **Completion Definition:** State machine reviewed by realtime and matchmaking owners.

## BATTLE-02: Implement Match Rooms and Authenticated Participant Routing
**Priority:** High · **Phase:** 4 · **Component:** Backend / Realtime · **Effort:** M · **Parallelism:** Contract-coordinated with BACK-08

**Description:** Join each participant to only their active duel room and route ordered events. **Business Requirements:** Architecture §6-7; API §12. **Technical Requirements:** Membership checks on join/send, event sequence and deduplication, cross-instance routing, bounded room lifetime. **Dependencies:** BATTLE-01, BACK-08.

**Implementation Steps:** 1. Bind authorized duel subscription. 2. Add per-match sequencing/event IDs. 3. Publish after durable commit/outbox. 4. Remove stale subscriptions on expiry/revocation. **Expected Files / Modules:** duel realtime adapter and event schemas. **API / Database Impact:** event log/outbox and match channel subscription.

**Acceptance Criteria:** Non-participant cannot subscribe; both participants receive same sequence; duplicate delivery can be detected; room does not own authoritative match state. **Testing Requirements:** Spoofed match ID, multi-instance delivery, duplicate and disconnect tests. **Completion Definition:** Two clients receive ordered events from separate gateway instances.

## BATTLE-03: Implement Round Lifecycle and Server Timer
**Priority:** Critical · **Phase:** 4 · **Component:** Backend · **Effort:** L · **Parallelism:** Central with question selector

**Description:** Start synchronized rounds with server-controlled deadlines and resolve on submissions/deadline. **Business Requirements:** PRD §8; GDD §2-4; UI/UX §5.8. **Technical Requirements:** Persist round start/deadline UTC, configured question count/timer, scheduler recovery and timeout resolution. **Dependencies:** BATTLE-01, QUES-03, BACK-09, FOUND-02.

**Implementation Steps:** 1. Persist rounds/deadlines. 2. Emit question to both players. 3. Resolve when both answered or deadline reached. 4. Rebuild timers after restart. 5. Enforce approved policy for backend outage crossing deadline. **Expected Files / Modules:** duel round service, timer worker and tests. **API / Database Impact:** duel_questions, `duel.started`, `duel.question_received`.

**Acceptance Criteria:** Both receive identical question/version/deadline; local clock does not affect resolution; exactly one resolution event. **Testing Requirements:** timer boundary, both-answer race, restart and outage policy tests. **Completion Definition:** Simulated clock and process restart tests pass.

## BATTLE-04: Implement Idempotent Answer Submission
**Priority:** Critical · **Phase:** 4 · **Component:** Backend · **Effort:** M · **Parallelism:** One implementation shared by REST and WebSocket

**Description:** Accept a single valid answer from each participant for current round. **Business Requirements:** PRD §8; API §7,12. **Technical Requirements:** Authenticated membership, expected round/question, server receipt time, deadline, schema validation, unique answer and command ID. **Dependencies:** BATTLE-01, BATTLE-03, QUES-05, BACK-06.

**Implementation Steps:** 1. Implement command handler. 2. Apply round/deadline/member validation. 3. Persist server timestamp and answer once. 4. Return stable duplicate/late errors. 5. Route REST fallback to same use case. **Expected Files / Modules:** answer use case, handler, repository. **API / Database Impact:** `POST /duels/{id}/answers`, `duel.answer_submit`, duel_answers.

**Acceptance Criteria:** Client correctness and duration ignored; duplicate submission never replaces first; late answer cannot score; malformed payload is rejected. **Testing Requirements:** Near-simultaneous answers, idempotency, timeout boundary and REST/WS parity. **Completion Definition:** Security and integration tests pass.

## BATTLE-05: Implement Deterministic Scoring and Winner/Tie Resolution
**Priority:** Critical · **Phase:** 4 · **Component:** Backend · **Effort:** M · **Parallelism:** Central product policy review

**Description:** Calculate per-round score, streak, speed bonus, aggregate score, winner and draw/sudden death from server inputs. **Business Requirements:** GDD §4; PRD §4,8. **Technical Requirements:** Correctness outranks speed; fast incorrect earns no speed bonus; deterministic integer calculations; rules snapshot. **Dependencies:** BATTLE-04, FOUND-02.

**Implementation Steps:** 1. Implement pure scoring function. 2. Store response duration and breakdown. 3. Apply configured tie-break order. 4. Implement sudden-death/draw policy. 5. Publish resolved score events. **Expected Files / Modules:** scoring domain package, config schema and fixtures. **API / Database Impact:** answer scores/result payload; `duel.answer_resolved`, `duel.score_updated`.

**Acceptance Criteria:** Faster incorrect answer never beats correct answer; output is deterministic; unresolved equal outcomes follow approved tie policy. **Testing Requirements:** Property/table tests for correctness, streak, speed threshold, tie and replay. **Completion Definition:** Product-approved examples and unit suite accepted.

## BATTLE-06: Implement Ability Activation Window and Match Snapshot
**Priority:** High · **Phase:** 6 · **Component:** Backend · **Effort:** M · **Parallelism:** Depends on CASTLE/ECON item contracts

**Description:** Validate and record ability commands during allowed match/combat phases. **Business Requirements:** GDD §9; PRD §8. **Technical Requirements:** Match-start loadout snapshot; server eligibility/effect/cooldown/charges/target; duplicate command rejection; consume only on successful activation. **Dependencies:** BATTLE-01, CASTLE-05, ECON-04.

**Implementation Steps:** 1. Snapshot eligible loadout at match start. 2. Define phase/target validation. 3. Apply activation through combat service. 4. Store effect and consumption idempotently. 5. Emit result event. **Expected Files / Modules:** duel ability command, combat interface and schema. **API / Database Impact:** `/duels/{id}/abilities`, `duel.ability_activate`, usage records.

**Acceptance Criteria:** Mid-match loadout changes have no effect; invalid target/window consumes nothing; valid duplicate command does not consume twice. **Testing Requirements:** Ownership, cooldown, stacking, duplicate and ranked-disabled ability tests. **Completion Definition:** Ability integration approved by balance/security owner.

## BATTLE-07: Implement Disconnect, Reconnect and Presence Policy
**Priority:** High · **Phase:** 4 · **Component:** Backend / Realtime · **Effort:** L · **Parallelism:** Parallel with client recovery after state contract

**Description:** Track connection state without allowing client presence to control gameplay outcomes. **Business Requirements:** PRD §8; UI/UX §5.8. **Technical Requirements:** Heartbeat, server grace/abandonment policy, active timer behavior, notifications and state version. **Dependencies:** BATTLE-02, BATTLE-03, BACK-08, FOUND-02.

**Implementation Steps:** 1. Track connection lease by session. 2. Emit disconnected/reconnected events. 3. Apply approved grace/abandonment policy. 4. Ensure timers continue or pause only by approved rule. 5. Handle multiple devices/session revocation. **Expected Files / Modules:** presence adapter and duel state machine. **API / Database Impact:** participant connection state and realtime events.

**Acceptance Criteria:** A disconnect cannot fabricate an answer or pause a timer locally; reconnection restores active state; duplicate connection status events are harmless. **Testing Requirements:** transient loss, prolonged disconnect, concurrent device and timer tests. **Completion Definition:** Network fault test on staging succeeds.

## BATTLE-08: Finalize Match Result, Rating, Damage and Rewards Once
**Priority:** Critical · **Phase:** 6 · **Component:** Backend / Database · **Effort:** XL · **Parallelism:** Centralized transaction owner

**Description:** Atomically persist immutable result and all ranked outcome effects. **Business Requirements:** PRD §4,8; DB §15; Architecture §8. **Technical Requirements:** Unique finalization key; result, MMR/history, castle damage, XP/wallet/items/missions/statistics and outbox consistent; retry returns original result. **Dependencies:** BATTLE-05, CASTLE-04, ECON-03, PROG-03, BACK-06.

**Implementation Steps:** 1. Define finalization transaction contract across modules. 2. Calculate validated result and effects. 3. Persist result and unique result marker. 4. Apply MMR/castle/economy/progression writes and outbox. 5. Make transaction retry-safe and observable. **Expected Files / Modules:** duel finalizer, transaction orchestration and integration tests. **API / Database Impact:** `duel_results`, rating, damage, ledger, reward/mission and outbox records.

**Acceptance Criteria:** Any injected failure rolls back all effects or permits safe replay; final result immutable; duplicate finalize cannot duplicate reward, damage or MMR. **Testing Requirements:** Failure injection at each write, concurrent finalization and ledger reconciliation. **Completion Definition:** Transaction design reviewed and staging fault-injection suite passes.

## BATTLE-09: Implement Match Cancellation, Timeout and Abandonment Policy
**Priority:** High · **Phase:** 4-6 · **Component:** Backend · **Effort:** M · **Parallelism:** Policy centralized

**Description:** Handle pre-start cancellation, server errors, no-show and abandonment consistently. **Business Requirements:** PRD §8; Roadmap §4-5. **Technical Requirements:** Explicit result types, rating/reward eligibility, audit reason, no fabricated winner; queue and active-match transition safety. **Dependencies:** BATTLE-01, BATTLE-07, FOUND-02.

**Implementation Steps:** 1. Define failure/result reasons. 2. Implement transition and cleanup handlers. 3. Apply approved rating/reward policy. 4. Notify participants and record history. **Expected Files / Modules:** duel/matchmaking state machine and tests. **API / Database Impact:** result_type/status and events.

**Acceptance Criteria:** Every terminal state has stable reason and policy; cancellation cannot leave a player locked in match; duplicate timeout jobs are safe. **Testing Requirements:** no-show, server failure, pre-confirm cancellation and repeated timeout tests. **Completion Definition:** Product policy and test cases approved.

## BATTLE-10: Implement Match State Snapshot and Event Replay
**Priority:** High · **Phase:** 4 · **Component:** Backend / Realtime · **Effort:** L · **Parallelism:** Central event contract owner

**Description:** Recover current duel state and replay retained events after reconnect or sequence gaps. **Business Requirements:** PRD §8; Architecture §7; API §12. **Technical Requirements:** Durable state version, event sequence, bounded replay window, complete snapshot fallback and retention. **Dependencies:** BATTLE-01, BATTLE-02, BACK-06.

**Implementation Steps:** 1. Persist critical sequence transitions. 2. Implement authorized `/state` snapshot. 3. Implement replay after last acknowledged sequence. 4. Define fallback when history expires. 5. Add retention and metrics. **Expected Files / Modules:** duel event log/state sync handler. **API / Database Impact:** `GET /duels/{id}/state`, `realtime.state_sync`, event log.

**Acceptance Criteria:** Snapshot reflects committed server state; replay has no missing/duplicate semantic transitions; expired history falls back to snapshot. **Testing Requirements:** sequence gap, restart, duplicate event and expired replay tests. **Completion Definition:** Reconnect suite passes across process restart.

## BATTLE-11: Build End-to-End Multiplayer Battle Verification
**Priority:** Critical · **Phase:** 4-8 · **Component:** Testing · **Effort:** L · **Parallelism:** Independent QA after contract stabilization

**Description:** Verify full synchronized duel on two clients with adversarial and recovery cases. **Business Requirements:** PRD §11; Roadmap §4-5. **Technical Requirements:** Realistic latency, separate clients, observability and deterministic fixtures. **Dependencies:** BATTLE-03..BATTLE-10, MOB-07.

**Implementation Steps:** 1. Automate two-client matchmaking/test fixture. 2. Cover normal win/loss/draw. 3. Inject duplicate/late answer and disconnect. 4. Verify result/rewards/castle state. 5. Capture trace IDs and failures. **Expected Files / Modules:** integration/E2E test harness and fixtures. **API / Database Impact:** Exercises REST, WebSocket, DB and Redis.

**Acceptance Criteria:** Two independent clients show the same final result; no answer key leaks; recovery returns consistent state and duplicate effects remain one. **Testing Requirements:** Cases T-02..T-10 in release strategy. **Completion Definition:** CI plus staging device evidence attached to release gate.
