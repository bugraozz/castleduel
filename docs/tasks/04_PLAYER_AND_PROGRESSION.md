# Player and Progression Tasks

## PROG-01: Implement Player Profiles, Privacy and Statistics
**Priority:** High · **Phase:** 1-2 · **Component:** Backend / Database · **Effort:** M · **Parallelism:** Parallel with mobile profile UI after contract

**Description:** Persist editable profiles, public-safe profiles, privacy settings and rebuildable statistics. **Business Requirements:** PRD §3,8; UI/UX §5.18; DB §3. **Technical Requirements:** Normalized unique username; allowlisted public fields; counters derived from validated events. **Dependencies:** BACK-02, BACK-04.

**Implementation Steps:** 1. Add profile/statistics schema and indexes. 2. Implement self read/update and public profile. 3. Add privacy policy filtering. 4. Add statistics projection/rebuild path. **Expected Files / Modules:** `internal/player`, profile migration and tests. **API / Database Impact:** `/players/me`, `/players/me/statistics`, `/players/{id}/public`; player_profiles/statistics.

**Acceptance Criteria:** Username uniqueness is case/normalization safe; public response contains no private identifiers; stats rebuild equals event-derived values. **Testing Requirements:** Authorization, privacy, uniqueness and projection tests. **Completion Definition:** Profile routes pass contract and security tests.

## PROG-02: Implement Player XP and Level Progression
**Priority:** Medium · **Phase:** 2 / P1 progression · **Component:** Backend · **Effort:** M · **Parallelism:** Independent after reward events stabilize

**Description:** Award account XP and calculate player levels from validated gameplay/missions. **Business Requirements:** GDD §10; PRD §5,8. **Technical Requirements:** Versioned progression thresholds, idempotent source event, anti-farming cap, immutable grant provenance. **Dependencies:** PROG-01, ECON-03, BACK-06.

**Implementation Steps:** 1. Define XP curve/config schema. 2. Add grant operation and progression projection. 3. Add duplicate/diminishing-return rules. 4. Expose player level/XP progress. **Expected Files / Modules:** progression service/config and migration. **API / Database Impact:** profile XP fields and validated grant events.

**Acceptance Criteria:** Duplicate event grants XP once; threshold transitions are deterministic; level cannot be client-set. **Testing Requirements:** Boundary, replay, cap and rebuild tests. **Completion Definition:** XP history and level displayed from server state.

## PROG-03: Implement Match Rating and Rating History
**Priority:** Critical · **Phase:** 4-6 · **Component:** Backend / Database · **Effort:** L · **Parallelism:** Centralized with duel finalization

**Description:** Maintain server-authoritative MMR separately from XP and castle level. **Business Requirements:** PRD §4,8; GDD §10; DB §9. **Technical Requirements:** Approved rating algorithm/config; only eligible ranked match; draw/abandon policy; unique per-duel update; before/after history. **Dependencies:** BATTLE-05, FOUND-02.

**Implementation Steps:** 1. Add rating state/history. 2. Implement deterministic rating function. 3. Integrate once-only finalization. 4. Add draw/abandon/friendly eligibility. 5. Expose current rating and history. **Expected Files / Modules:** `internal/progression/rating`, migration and tests. **API / Database Impact:** rating fields/history; `/leagues/me`, match result payload.

**Acceptance Criteria:** Friendly/ineligible matches do not change MMR; duplicate result does not change it twice; rating trace includes before/after and config version. **Testing Requirements:** Unit simulations for win/loss/draw and integration concurrency/replay. **Completion Definition:** Approved algorithm and match-result integration tested.

## PROG-04: Implement League Mapping and Competitive Profile
**Priority:** High · **Phase:** 7 competitive P1 · **Component:** Backend · **Effort:** M · **Parallelism:** Parallel after PROG-03

**Description:** Map rating to configurable tiers/divisions and show promotion/demotion state. **Business Requirements:** PRD §5 P1; GDD §11; UI/UX §5.14. **Technical Requirements:** Server-configured thresholds, stable tie rules, no client hardcoding. **Dependencies:** PROG-03.

**Implementation Steps:** 1. Add league configuration and version. 2. Implement rating-to-league mapping. 3. Add league/profile API. 4. Add historical changes in results. **Expected Files / Modules:** league package, seed/config and tests. **API / Database Impact:** leagues table; `/leagues/me`.

**Acceptance Criteria:** Boundary ratings map consistently; config changes do not rewrite historical match snapshots; UI receives current tier and progress. **Testing Requirements:** Threshold, config-version and API contract tests. **Completion Definition:** League progression validated against approved tier table.

## PROG-05: Implement Seasons and Leaderboard Projection
**Priority:** High · **Phase:** 7 competitive P1 · **Component:** Backend / Database · **Effort:** L · **Parallelism:** Can parallelize with mission work after PROG-03

**Description:** Operate seasonal ranked state, stable leaderboards and rating history. **Business Requirements:** PRD §5 P1,8; GDD §12; API §9. **Technical Requirements:** UTC start/end, stable pagination/tie-break, validated ranked records, archival, separate season reset from permanent castle state. **Dependencies:** PROG-03, PROG-04.

**Implementation Steps:** 1. Add season and ranking schema. 2. Implement active season lookup and ranking projection. 3. Define tie-break and stable cursor. 4. Add global/friend/season reads. 5. Archive closed-season leaderboard. **Expected Files / Modules:** season/leaderboard modules and migrations. **API / Database Impact:** `/seasons/current`, `/leaderboards/*`, season tables.

**Acceptance Criteria:** Only eligible matches contribute; pagination is stable; rating reset never alters permanent castle/inventory/cosmetics. **Testing Requirements:** Concurrent updates, stable pagination and season boundary integration tests. **Completion Definition:** Staging leaderboard agrees with validated match data.

## PROG-06: Implement Idempotent Season Closure and Rewards
**Priority:** High · **Phase:** 7-8 · **Component:** Backend · **Effort:** L · **Parallelism:** Central job/reward coordination

**Description:** Close seasons, reconcile pending matches and distribute free/premium-track rewards exactly once. **Business Requirements:** GDD §12; Economy §9; Roadmap §2 Epic 7. **Technical Requirements:** Resumable job, eligibility snapshot, unique claim/grant, preserve permanent progression. **Dependencies:** PROG-05, ECON-03, ECON-06, BACK-09.

**Implementation Steps:** 1. Freeze eligible results at end. 2. Reconcile active matches per policy. 3. Calculate final ranks/reward eligibility. 4. Grant and record reward once. 5. Archive and apply configured rating soft reset. **Expected Files / Modules:** season closure job, reward config and audit. **API / Database Impact:** season rewards/claims; reward ledger/inventory records.

**Acceptance Criteria:** Re-running closure produces no duplicate grants; interrupted job resumes; season reset leaves castle and permanent entitlements unchanged. **Testing Requirements:** failure injection, duplicate run, pending-match and boundary tests. **Completion Definition:** Restore/replay drill succeeds in staging.

## PROG-07: Implement Missions and Claim Flow
**Priority:** Medium · **Phase:** 7 P1 retention · **Component:** Backend · **Effort:** L · **Parallelism:** Parallel with achievements after event taxonomy

**Description:** Add daily/weekly/season missions driven only by validated server events. **Business Requirements:** PRD §5 P1; GDD §15; UI/UX §5.19. **Technical Requirements:** Period keys, UTC boundaries, duplicate-event dedupe, claim idempotency, reward transaction. **Dependencies:** ECON-03, ECON-06, BACK-06.

**Implementation Steps:** 1. Add definitions/player progress schema. 2. Define supported validated events and period rotation. 3. Update progress idempotently. 4. Implement list/claim. 5. Add configurable caps and expiry. **Expected Files / Modules:** mission module, job and tests. **API / Database Impact:** `/missions`, `/missions/{id}/claim`; mission tables.

**Acceptance Criteria:** Duplicate events do not increment twice; claim requires completion and grants once; timezone boundaries are explicit. **Testing Requirements:** event replay, midnight/season boundary, concurrent claim tests. **Completion Definition:** Mission progress and reward verified in staging.

## PROG-08: Implement Achievements and Rewards
**Priority:** Medium · **Phase:** 7 P1 retention · **Component:** Backend · **Effort:** M · **Parallelism:** Parallel with PROG-07

**Description:** Unlock cumulative achievements from validated account events. **Business Requirements:** PRD §5 P1; GDD §15; API §11. **Technical Requirements:** Versioned definitions, monotonic progress, unique unlock and idempotent claim. **Dependencies:** ECON-03, ECON-06, BACK-06.

**Implementation Steps:** 1. Define achievement rule schema. 2. Implement event progress evaluator. 3. Add unlock and claim states. 4. Expose list/detail API. **Expected Files / Modules:** achievement module/config and tests. **API / Database Impact:** `/achievements`, `/achievements/{id}/claim`; achievement tables.

**Acceptance Criteria:** Replayed events cannot duplicate progress or grants; rules are server-controlled; UI can distinguish locked/unlocked/claimed. **Testing Requirements:** threshold, event replay and concurrent claim tests. **Completion Definition:** Initial achievement catalog approved and tested.

## PROG-09: Implement Player Match History
**Priority:** Medium · **Phase:** 6 · **Component:** Backend · **Effort:** M · **Parallelism:** After BATTLE-08

**Description:** Provide paginated private history and safe public summary. **Business Requirements:** PRD §5,8; UI/UX §5.18; DB §5. **Technical Requirements:** Stable cursor, ownership checks, match-mode visibility, retained immutable result. **Dependencies:** BATTLE-08, PROG-01.

**Implementation Steps:** 1. Add indexed participant-history query. 2. Define private/public fields. 3. Add pagination and filters. 4. Connect result details and damage/rewards. **Expected Files / Modules:** player history handlers and query tests. **API / Database Impact:** `/players/me/matches` (recommended addition) and public match summary as policy allows.

**Acceptance Criteria:** Only authorized data is returned; ordering is stable under new matches; no private opponent data leaks. **Testing Requirements:** Pagination, privacy and access control tests. **Completion Definition:** Mobile match-history screen consumes stable contract.
