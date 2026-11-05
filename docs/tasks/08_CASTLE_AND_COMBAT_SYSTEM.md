# Castle and Combat System Tasks

## CASTLE-01: Create Persistent Starter Castle
**Priority:** Critical · **Phase:** 2 Castle/economy slice · **Component:** Backend / Database · **Effort:** M · **Parallelism:** Parallel after user creation hook

**Description:** Provision one primary castle for each eligible account and expose its state. **Business Requirements:** PRD §5,8; GDD §6. **Technical Requirements:** Unique user ownership, starter configuration, current/max durability, version and safe idempotent provisioning. **Dependencies:** BACK-02, BACK-04.

**Implementation Steps:** 1. Add castle schema and constraints. 2. Provision on onboarding/account lifecycle transaction or retryable job. 3. Implement owner-only castle read. 4. Add invariant/reconciliation for missing/duplicate castle. **Expected Files / Modules:** `internal/castle`, migrations and provisioning hook. **API / Database Impact:** `GET /castle`; castles table.

**Acceptance Criteria:** One primary castle per active user; retry cannot create duplicate; castle persists across logout/reinstall. **Testing Requirements:** Concurrent provisioning, authorization and deletion-policy tests. **Completion Definition:** New account obtains castle and fetches it after a new session.

## CASTLE-02: Model Structures, Slots and Versioned Definitions
**Priority:** High · **Phase:** 2 · **Component:** Backend / Database · **Effort:** M · **Parallelism:** Content/config can parallelize

**Description:** Define structures, unlock requirements, slots, durability and defensive attributes. **Business Requirements:** GDD §6; DB §4. **Technical Requirements:** Server-controlled definitions, max level, slot uniqueness and non-destructive combat status. **Dependencies:** CASTLE-01, FOUND-05.

**Implementation Steps:** 1. Add definition and owned-structure schema. 2. Seed reviewed structure catalog. 3. Validate levels/slots/unlock rules. 4. Expose read representation with current/max durability and effects. **Expected Files / Modules:** castle definitions/config, schema and seed data. **API / Database Impact:** `/castle/structures`; structure_definitions/castle_structures.

**Acceptance Criteria:** Invalid slot/level/definition is rejected; structure destruction cannot remove level/ownership; definitions are versioned for match calculations. **Testing Requirements:** Constraints, seed migration and definition compatibility tests. **Completion Definition:** Initial structure list approved and API schema published.

## CASTLE-03: Implement Construction and Upgrade Transactions
**Priority:** High · **Phase:** 2 · **Component:** Backend / Database · **Effort:** L · **Parallelism:** Central with ECON ledger owner

**Description:** Build/unlock or upgrade a structure with validated cost, prerequisites and slot. **Business Requirements:** PRD §7-8; UI/UX §5.11; DB §15. **Technical Requirements:** Atomic resource debit/state update, idempotency, expected version and max-level checks. **Dependencies:** CASTLE-02, ECON-01, ECON-02, BACK-06.

**Implementation Steps:** 1. Define cost/prerequisite evaluator. 2. Lock/check wallet and castle versions. 3. Insert debit ledger and mutate structure in one transaction. 4. Return authoritative resulting stats/cost. **Expected Files / Modules:** construction use case and tests. **API / Database Impact:** `POST /castle/structures`, `.../{id}/upgrade`; ledger/structure tables.

**Acceptance Criteria:** Concurrent purchases cannot overspend or occupy one slot twice; retries charge once; insufficient funds gives clear error. **Testing Requirements:** Simultaneous build, stale version, duplicate request and prerequisite tests. **Completion Definition:** Transaction integration suite passes.

## CASTLE-04: Implement Idempotent Castle Damage Resolution
**Priority:** Critical · **Phase:** 6 Combat and rewards · **Component:** Backend / Database · **Effort:** L · **Parallelism:** Central with BATTLE finalizer

**Description:** Apply damage from a resolved duel using server-calculated attack/defense and immutable combat event. **Business Requirements:** GDD §5-7; PRD §8. **Technical Requirements:** Integer damage, configured floor/caps/protection, structure targeting, unique combat ID, atomic aggregate/structure updates. **Dependencies:** CASTLE-02, CASTLE-05, BATTLE-05, FOUND-02.

**Implementation Steps:** 1. Define damage formula and snapshot version. 2. Validate target and protection. 3. Apply damage once and append history. 4. Update aggregate durability consistently. 5. Emit post-commit event. **Expected Files / Modules:** combat resolver, castle damage repository and migration. **API / Database Impact:** castle_damage_events; `duel.attack_resolved`, `duel.castle_damaged`.

**Acceptance Criteria:** Repeated combat ID does not apply damage twice; damage cannot exceed remaining durability or approved bounds; core progression is not permanently deleted. **Testing Requirements:** concurrent attacks, cap/floor, shield/protection, destroyed structure and replay tests. **Completion Definition:** Formula and transaction reviewed by product/balance owner.

## CASTLE-05: Implement Defense, Shield and Protection Rules
**Priority:** High · **Phase:** 6 · **Component:** Backend · **Effort:** M · **Parallelism:** Depends on abilities/config

**Description:** Resolve defensive structure effects, shield windows and new-player/repeated-attack protection. **Business Requirements:** GDD §6-7,9; PRD §8. **Technical Requirements:** Server-configured caps, expiry, target eligibility, no infinite immunity; fair, visible explanation. **Dependencies:** CASTLE-02, ECON-04, FOUND-02.

**Implementation Steps:** 1. Define ordered defense calculation. 2. Snapshot active shields and configuration. 3. Add protection expiry and refresh policy. 4. Expose current defenses in castle state. **Expected Files / Modules:** defense evaluator, config and read model. **API / Database Impact:** castle protection state, ability usage/effects.

**Acceptance Criteria:** Defense never exceeds approved maximum; expired/invalid shield has no effect; calculations deterministic and auditable. **Testing Requirements:** Stacking/caps, expiry, simultaneous attack and new-player protection tests. **Completion Definition:** Balance policy signed off and test fixtures cover combinations.

## CASTLE-06: Implement Repair and Restoration
**Priority:** High · **Phase:** 2-6 · **Component:** Backend / Database · **Effort:** M · **Parallelism:** Central with economy ledger

**Description:** Repair whole castle or selected structure using earnable resources or valid repair kit. **Business Requirements:** GDD §8; PRD §8; Economy §4. **Technical Requirements:** Formula/config snapshot, no over-repair, atomic debit/durability update, idempotency, no mandatory wait. **Dependencies:** CASTLE-02, ECON-01, ECON-02, ECON-04.

**Implementation Steps:** 1. Implement repair quote and cost calculation. 2. Validate balance or consumable. 3. Update durability and ledger/inventory atomically. 4. Record repair history and return restored amount. **Expected Files / Modules:** repair service and migration. **API / Database Impact:** `/castle/repair`, `/{structureId}/repair`; castle_repairs, ledger.

**Acceptance Criteria:** Repair cannot exceed max or spend beyond balance; duplicate call repairs/charges once; insufficient resources return clear required/available values. **Testing Requirements:** concurrent repairs, cost boundaries, idempotency and repair-kit use. **Completion Definition:** Staging repair flow reconciles wallet and durability.

## CASTLE-07: Implement Offensive Abilities and Target Resolution
**Priority:** High · **Phase:** 6 · **Component:** Backend · **Effort:** L · **Parallelism:** Central with BATTLE-06

**Description:** Implement documented offensive ability family (Destruction Bomb, Critical Strike, Siege Attack, Piercing Attack) through data-driven server definitions. **Business Requirements:** GDD §9; PRD §8. **Technical Requirements:** Effect value/cap, phase, target, ranked eligibility, counters, ownership/charges and no client-defined effect. **Dependencies:** CASTLE-04, ECON-04, BATTLE-06, FOUND-02.

**Implementation Steps:** 1. Define ability catalog and interaction rules. 2. Implement validation/effect resolver. 3. Apply target-specific damage/defense changes. 4. Record use/result/idempotency. 5. Add balance simulation fixtures. **Expected Files / Modules:** ability definitions, combat resolver and admin config. **API / Database Impact:** ability definitions/usage, duel ability event.

**Acceptance Criteria:** Each ability has documented target and counter; invalid use consumes nothing; ranked caps enforced server-side. **Testing Requirements:** Every ability condition, stacking/counter, duplicate and ranked restriction. **Completion Definition:** Ability catalog passes product/security balance checklist.

## CASTLE-08: Implement Damage History and Castle State Recovery
**Priority:** Medium · **Phase:** 6 · **Component:** Backend / Mobile · **Effort:** M · **Parallelism:** Parallel after damage event contract

**Description:** Expose attack history and reconcile castle aggregate/structure state. **Business Requirements:** PRD §7; UI/UX §5.12; DB §4. **Technical Requirements:** Owner authorization, stable pagination, immutable event history, aggregate consistency repair job with alert. **Dependencies:** CASTLE-04, CASTLE-06.

**Implementation Steps:** 1. Add damage-history endpoint and indexes. 2. Add reconcile query/job. 3. Expose attacker-safe public summary. 4. Connect Defense screen. **Expected Files / Modules:** castle history handlers, reconciliation job and mobile screen. **API / Database Impact:** `/castle/damage-history`; damage events.

**Acceptance Criteria:** History matches applied damage; unauthorized user cannot inspect private details; reconciliation detects but does not silently rewrite audit history. **Testing Requirements:** Pagination/privacy and aggregate mismatch alert tests. **Completion Definition:** Attack-to-repair history verified end-to-end.

## CASTLE-09: Implement Castle Construction, Damage and Repair E2E Tests
**Priority:** High · **Phase:** 6-8 · **Component:** Testing · **Effort:** M · **Parallelism:** QA after API freeze

**Description:** Verify persistent castle lifecycle across duel, upgrade, attacks and recovery. **Business Requirements:** PRD §11; Roadmap §4. **Technical Requirements:** Independent account fixtures and transaction-state assertions. **Dependencies:** CASTLE-03..CASTLE-08, ECON-01..ECON-04.

**Implementation Steps:** 1. Provision castle and structures. 2. Build/upgrade/repair. 3. Apply concurrent attacks and verify caps. 4. Restart/reload and reconcile. **Expected Files / Modules:** castle integration/E2E tests. **API / Database Impact:** Exercises castle/economy/duel endpoints.

**Acceptance Criteria:** State persists across session; no duplicated damage/debit; restoration is possible after loss. **Testing Requirements:** Fault injection and concurrency. **Completion Definition:** CI and staging evidence attached to release checklist.
