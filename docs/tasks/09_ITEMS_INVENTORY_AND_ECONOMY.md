# Items, Inventory and Economy Tasks

## ECON-01: Implement Wallet Projection and Immutable Currency Ledger
**Priority:** Critical · **Phase:** 2 Castle/economy slice · **Component:** Backend / Database · **Effort:** L · **Parallelism:** Central financial-data owner

**Description:** Establish Gold, Stone and Crystals balances with append-only transaction history. **Business Requirements:** PRD §8; GDD §13; Economy §2,5. **Technical Requirements:** Integer units, nonnegative balance, immutable source/reason/reference, idempotency, transactional projection and reconciliation. **Dependencies:** BACK-04, BACK-06.

**Implementation Steps:** 1. Add wallet/ledger schema and indexes. 2. Implement grant/debit operation. 3. Add balance projection update in same transaction. 4. Add reconciliation job and mismatch alert. 5. Expose read/history. **Expected Files / Modules:** `internal/economy`, migrations and worker. **API / Database Impact:** `/wallet`, `/wallet/transactions`; wallets/currency_transactions.

**Acceptance Criteria:** No client balance mutation; duplicate key cannot double-credit/debit; projection equals ledger; negative balance rejected. **Testing Requirements:** Concurrent debit, duplicate, rollback, projection reconcile and pagination tests. **Completion Definition:** Ledger invariants pass integration and reconciliation checks.

## ECON-02: Implement Atomic Construction, Repair and Resource Spend Integration
**Priority:** Critical · **Phase:** 2-6 · **Component:** Backend · **Effort:** L · **Parallelism:** Central with CASTLE

**Description:** Provide the shared transactional debit contract used by castle mutations and other resource sinks. **Business Requirements:** PRD §8; GDD §8,13; Economy §4-5. **Technical Requirements:** Same transaction handle, scoped idempotency, clear insufficient balance, server-configured cost; this task establishes the primitive and tests, while consumer tasks integrate it. **Dependencies:** ECON-01, BACK-06.

**Implementation Steps:** 1. Define shared debit interface. 2. Require caller's transaction context. 3. Store operation reference on ledger. 4. Add reusable insufficient-funds and idempotency behavior. **Expected Files / Modules:** economy transaction service and contract tests. **API / Database Impact:** Supports construction/repair endpoints and ledger.

**Acceptance Criteria:** Debit cannot commit outside its caller's mutation transaction; competing spends serialize; duplicate key returns original outcome. **Testing Requirements:** Rollback, competing debit and idempotency contract tests. **Completion Definition:** Castle and repair consumers pass transaction review when integrated.

## ECON-03: Implement Server-Side Match Rewards
**Priority:** Critical · **Phase:** 6 · **Component:** Backend · **Effort:** L · **Parallelism:** Central with BATTLE finalizer

**Description:** Calculate fair ranked/friendly participation and win rewards from validated outcome. **Business Requirements:** PRD §8; GDD §14; Economy §3,12. **Technical Requirements:** Configured reward tables, fair completion, abandonment/farming limits, idempotency and provenance. **Dependencies:** ECON-01, BATTLE-05, FOUND-02.

**Implementation Steps:** 1. Define match mode/outcome reward policy. 2. Implement deterministic reward calculator. 3. Grant via ledger/inventory transaction. 4. Record reason and source match. 5. Add source/sink metrics. **Expected Files / Modules:** reward service/config and tests. **API / Database Impact:** reward ledger/inventory and result payload.

**Acceptance Criteria:** Client cannot supply rewards; repeated finalization grants once; loss reward is limited to fair completion; friendly rewards follow anti-farm config. **Testing Requirements:** win/loss/draw/abandonment, duplicate and farming tests. **Completion Definition:** Product-approved reward table and transaction suite pass.

## ECON-04: Implement Item Catalog, Ownership and Consumption
**Priority:** High · **Phase:** 2-6 · **Component:** Backend / Database · **Effort:** L · **Parallelism:** Catalog can parallelize with combat contracts

**Description:** Implement item definitions, inventory ownership, expiry, equipment and consumption for offensive/defensive/repair items. **Business Requirements:** PRD §8; GDD §9; Economy §7. **Technical Requirements:** Server-defined effects, grant/consume ledger, loadout slot/charge constraints, no duplicate stacking unless configured. **Dependencies:** ECON-01, BACK-06, FOUND-02.

**Implementation Steps:** 1. Add definitions, owned lots/transactions and equip schema. 2. Implement list/equip/unequip. 3. Implement consume interface with validation hooks. 4. Add expiry and balance checks. 5. Ensure match loadout snapshot. **Expected Files / Modules:** inventory module, migrations and handlers. **API / Database Impact:** `/inventory`, `/inventory/equip`, `/inventory/unequip`, `/inventory/use`.

**Acceptance Criteria:** Ownership/expiry enforced; invalid activation does not consume; use and effect commit once; client cannot set effect magnitude. **Testing Requirements:** Concurrent consumption, expiry, duplicate, slot and stacking tests. **Completion Definition:** Inventory integration works against combat use case.

## ECON-05: Implement Economy Transaction History and Support Inspection
**Priority:** Medium · **Phase:** 6-8 · **Component:** Backend / Admin · **Effort:** M · **Parallelism:** After ledger schema

**Description:** Provide player-readable transaction history and permission-gated support inspection. **Business Requirements:** PRD §8; UI/UX §5.13; Roadmap §3. **Technical Requirements:** Stable pagination, clear reason/source, PII minimization, read-only support scope. **Dependencies:** ECON-01, ADMIN-01.

**Implementation Steps:** 1. Define user-safe ledger projection. 2. Implement paginated history. 3. Add role-gated admin lookup. 4. Add reconciliation/refund cross-reference. **Expected Files / Modules:** economy read API and admin inspection. **API / Database Impact:** `/wallet/transactions`, admin player economy read.

**Acceptance Criteria:** History shows every balance change with source; support cannot mutate directly; inaccessible account data is denied. **Testing Requirements:** Pagination, redaction and role access tests. **Completion Definition:** Support workflow reviewed and ledger query is indexed.

## ECON-06: Add Reward Claim and Idempotency Framework
**Priority:** High · **Phase:** 7 · **Component:** Backend · **Effort:** M · **Parallelism:** Parallel with mission/season modules after BACK-06

**Description:** Standardize user-initiated reward claims across missions, achievements, seasons and events. **Business Requirements:** GDD §12,15; API §8-11. **Technical Requirements:** Eligibility check, one claim per source, grant and claim status in one transaction. **Dependencies:** ECON-01, BACK-06.

**Implementation Steps:** 1. Define reward source identifiers. 2. Add unique claim key. 3. Grant resources/items through ledger. 4. Return prior outcome for retries. 5. Add notification/outbox. **Expected Files / Modules:** shared reward claim service. **API / Database Impact:** claim routes and reward-source constraints.

**Acceptance Criteria:** Concurrent claims grant once; ineligible claim leaves no partial record; every grant links to its source. **Testing Requirements:** Duplicate and concurrent claim, rollback and foreign-source tests. **Completion Definition:** Mission/season claims both consume common service.

## ECON-07: Add Economy Configuration and Balance Simulation
**Priority:** High · **Phase:** 6-8 · **Component:** Backend / Analytics · **Effort:** L · **Parallelism:** Simulation can parallelize after config schema

**Description:** Make costs, rewards, repair curves and caps tunable and validate economic sustainability. **Business Requirements:** GDD §13,17; Economy §12; Roadmap §3. **Technical Requirements:** Versioned config, active-match snapshot, simulation for activity/win-rate/free-paid paths, inflation and repair affordability metrics. **Dependencies:** ECON-01, BACK-07.

**Implementation Steps:** 1. Define typed config schema and ranges. 2. Build deterministic simulation harness. 3. Run 1/7/30/90-day player profiles and 30-70% win rates. 4. Report sinks/sources, repair failures and progression times. 5. Block invalid config activation. **Expected Files / Modules:** economy config, simulator and reports. **API / Database Impact:** versioned game configuration.

**Acceptance Criteria:** Simulation reproducible from config version; core repair is achievable without payment in tested profiles; invalid bounds cannot publish. **Testing Requirements:** Golden simulation and config validation tests. **Completion Definition:** Economy owner signs off metrics before ranked launch.

## ECON-08: Implement Store Catalog, Cosmetic Entitlements and Season Pass
**Priority:** High · **Phase:** 8 Store/admin/launch · **Component:** Backend / Mobile · **Effort:** L · **Parallelism:** Catalog and UI can parallelize with fixed product schema

**Description:** Deliver remotely configurable cosmetic catalog, free/premium season tracks and entitlement checks. **Business Requirements:** PRD §5 P1,8; GDD §12,16; Economy §8-10; UI/UX §5.15-5.16. **Technical Requirements:** Product contents/price/duration/limits disclosed; cosmetics no stats; premium pass status verified server-side. **Dependencies:** ECON-04, PROG-05, ADMIN-05.

**Implementation Steps:** 1. Model products, tracks and entitlements. 2. Add catalog API. 3. Integrate native purchase UI. 4. Gate premium track by verified entitlement. 5. Support expiry and season claim deadline. **Expected Files / Modules:** store and mobile screens, catalog config. **API / Database Impact:** `/store/products`, season reward endpoints, products/entitlements.

**Acceptance Criteria:** Catalog can be updated without app release; cosmetic items do not change game stats; free track remains available; price/contents clear. **Testing Requirements:** Catalog version, entitlement gating, expiry and UI purchase states. **Completion Definition:** Sandbox store review signs off catalog and season pass.

## ECON-09: Implement Server-Side Native Purchase Verification
**Priority:** Critical · **Phase:** 8 · **Component:** Backend · **Effort:** XL · **Parallelism:** Provider-specific integration owner

**Description:** Verify Apple/Google signed purchase data and fulfill transaction exactly once. **Business Requirements:** PRD §8; Economy §13-14; API §10. **Technical Requirements:** Platform API verification, unique platform transaction, pending state, restore, refunds/revocations, idempotent fulfillment and provider outage retry. **Dependencies:** ECON-08, BACK-06, FOUND-02.

**Implementation Steps:** 1. Implement provider verification adapters and sandbox config. 2. Persist pending purchase. 3. Verify with platform server. 4. Fulfill currency/entitlement transactionally. 5. Process notifications/refunds/revocations and reconcile. 6. Add support diagnostics. **Expected Files / Modules:** `internal/purchase`, provider adapters and jobs. **API / Database Impact:** `/purchases/verify`, `/purchases`, `/purchases/restore`; purchases/entitlements.

**Acceptance Criteria:** Invalid/unverified purchase grants nothing; duplicate transaction grants once; refunds/revocations update entitlement; provider outage leaves retryable pending state. **Testing Requirements:** Sandbox, forged token, duplicate notification, refund, restore and outage tests. **Completion Definition:** Platform sandbox evidence and operational runbook.

## ECON-10: Implement Fairness, Anti-Farming and Pay-to-Win Controls
**Priority:** Critical · **Phase:** 6-8 · **Component:** Backend / Game Design · **Effort:** L · **Parallelism:** Central balance/security review

**Description:** Enforce free recovery and ranked constraints across paid and earned items. **Business Requirements:** PRD §2,4,8; GDD §9,16; Economy §10-11. **Technical Requirements:** Ranked slot/charge caps, free earn paths, no paid immunity/MMR, repeated-friendly-match diminishing rewards, collusion signals. **Dependencies:** ECON-03, ECON-04, ECON-07, CASTLE-05.

**Implementation Steps:** 1. Encode ranked eligibility per ability/item. 2. Add anti-farm windows and repeated-opponent telemetry. 3. Add suspicious flow monitoring and moderation review path. 4. Validate basic repair remains possible. 5. Publish clear player-facing item descriptions. **Expected Files / Modules:** combat/economy policy, analytics and admin investigation. **API / Database Impact:** config, reward caps and abuse events.

**Acceptance Criteria:** Purchase cannot bypass competitive limits or guarantee victory; suspicious friendly farming is limited and observable; no premium currency required for basic ranked recovery. **Testing Requirements:** free/paid progression comparisons, repeat-pair and ranked ability tests. **Completion Definition:** Fairness review and balance simulation approved before store launch.
