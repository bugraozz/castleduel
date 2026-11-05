# Mobile Foundation Tasks

## MOB-01: Initialize Expo TypeScript App and Navigation
**Implementation status:** Blocked pending FOUND-01 and FOUND-03 completion. The Expo TypeScript workspace and lint/type-check commands are verified; Expo Router, route ownership, auth/active-match guards, deep links, and route smoke tests remain unimplemented until prerequisites pass.
**Priority:** Critical · **Phase:** 1 Identity and app shell · **Component:** Mobile · **Effort:** M · **Parallelism:** Parallel after FOUND-01

**Description:** Establish Expo app, route groups and navigation structure for all documented screens. **Business Requirements:** PRD §6; UI/UX §4-5. **Technical Requirements:** Expo Router, typed routes, auth gates, deep-link mapping and minimum version/update state. **Dependencies:** FOUND-01, FOUND-03.

**Implementation Steps:** 1. Pin Expo SDK and compatible modules. 2. Add route groups for auth, tabs and contextual match screens. 3. Add auth and active-match route guards. 4. Add startup, mandatory update, maintenance and offline shells. **Expected Files / Modules:** `apps/mobile/app`, `src/navigation`, app config. **API / Database Impact:** Uses auth/profile routes when available.

**Acceptance Criteria:** All 21 documented screen destinations have route ownership or an explicitly grouped route; protected paths require auth; active match cannot be silently discarded by navigation. **Testing Requirements:** Navigation/auth/deep-link tests on iOS and Android. **Completion Definition:** App shell builds and route smoke suite passes.

## MOB-02: Build Accessible Theme and Reusable UI States
**Priority:** High · **Phase:** 1-2 · **Component:** Mobile · **Effort:** L · **Parallelism:** Parallel with backend after UX token review

**Description:** Create centralized tokens and reusable accessible components. **Business Requirements:** UI/UX §2-3,6-8. **Technical Requirements:** Theme tokens, scalable text, safe area, contrast, screen-reader semantics, reduced motion and shared loading/empty/error/offline/success patterns. **Dependencies:** MOB-01.

**Implementation Steps:** 1. Define token names and theme provider. 2. Implement buttons, cards, dialogs, resource indicators, forms and feedback states. 3. Add accessibility labels/roles and reduced motion. 4. Validate small screens and dynamic text. **Expected Files / Modules:** `src/shared/ui`, tokens and component tests. **API / Database Impact:** None.

**Acceptance Criteria:** No raw repeated color/spacing values in shared components; interactive targets and status are accessible; all relevant screens can render network states. **Testing Requirements:** Component accessibility and snapshot/visual checks. **Completion Definition:** Design QA accepts component library and accessibility checklist.

## MOB-03: Implement Typed API Client and Server-State Layer
**Priority:** High · **Phase:** 1-2 · **Component:** Mobile · **Effort:** M · **Parallelism:** Requires FOUND-04 contracts

**Description:** Centralize API communication, errors, retries and cache invalidation. **Business Requirements:** Architecture §2; UI/UX §6. **Technical Requirements:** OpenAPI-generated/checked client, TanStack Query, request IDs, idempotency keys, no automatic retry for unsafe mutations. **Dependencies:** MOB-01, FOUND-04.

**Implementation Steps:** 1. Generate client from OpenAPI. 2. Add transport/auth/error mapping. 3. Configure query defaults and per-feature keys. 4. Add mutation invalidation rules. 5. Add offline/maintenance presentation. **Expected Files / Modules:** `src/shared/api`, `src/shared/state`, generated contracts. **API / Database Impact:** Client coverage follows documented `/api/v1`.

**Acceptance Criteria:** Contract drift fails CI; stable errors render user-safe messages; unsafe mutations are never retried without idempotency. **Testing Requirements:** Mocked transport, schema mismatch and refresh/error tests. **Completion Definition:** Profile and castle query use the shared client.

## MOB-04: Implement Secure Authentication Screens and Session Handling
**Priority:** High · **Phase:** 1 · **Component:** Mobile · **Effort:** L · **Parallelism:** Can parallel with BACK after contract freeze

**Description:** Deliver registration/login/recovery/logout UI and secure session lifecycle. **Business Requirements:** PRD §5,8; UI/UX §5.3,5.21. **Technical Requirements:** SecureStore for refresh credentials, refresh race handling, generic auth errors, validation and explicit revoked/offline states. **Dependencies:** MOB-03, BACK-02.

**Implementation Steps:** 1. Build forms and field validation. 2. Store/rotate credentials securely. 3. Handle expiration, logout and account deletion. 4. Add terms/privacy links and loading states. **Expected Files / Modules:** `features/auth`, secure storage adapter, tests. **API / Database Impact:** `/auth/*`, `/account`.

**Acceptance Criteria:** No credential in AsyncStorage/log/analytics; refresh is serialized; logout/deletion clears local state; error copy avoids account enumeration. **Testing Requirements:** Auth E2E, secure storage adapter tests, refresh race and offline cases. **Completion Definition:** Registration-to-home flow succeeds on supported platforms.

## MOB-05: Add Typed Realtime Client and Recovery State
**Priority:** High · **Phase:** 4 · **Component:** Mobile / Realtime · **Effort:** L · **Parallelism:** Coordinate event schemas with BATTLE

**Description:** Build reconnecting WebSocket manager with sequencing and snapshot recovery. **Business Requirements:** PRD §8-9; UI/UX §4,5.8. **Technical Requirements:** Typed event validation, heartbeat, jittered backoff, command IDs, event dedupe, ack, gap detection, state sync. **Dependencies:** MOB-03, BACK-08, FOUND-04.

**Implementation Steps:** 1. Implement connection/auth lifecycle. 2. Validate inbound envelopes. 3. Track per-match sequence/version. 4. Reconnect and request replay/snapshot. 5. Expose connection state to UI. **Expected Files / Modules:** `src/shared/realtime`, duel feature adapter. **API / Database Impact:** WebSocket events and `/duels/{id}/state`.

**Acceptance Criteria:** Duplicate events apply once; sequence gap requests recovery; reconnect restores authoritative match; user sees disconnected/reconnecting state. **Testing Requirements:** Unit tests for ordering and E2E background/reconnect. **Completion Definition:** Test device can recover an active match after network interruption.

## MOB-06: Implement Onboarding and Home/Castle Shell
**Priority:** Medium · **Phase:** 2 · **Component:** Mobile · **Effort:** L · **Parallelism:** Parallel with castle API after design specs

**Description:** Deliver first-time walkthrough and main hub with castle/resources/league/mission summaries. **Business Requirements:** PRD §6-7; UI/UX §5.2,5.4. **Technical Requirements:** Server-backed values, clear Find Opponent and Challenge Friend actions, tutorial completion state. **Dependencies:** MOB-02, MOB-03, PROG-01, CASTLE-01.

**Implementation Steps:** 1. Build skip/consent-aware onboarding. 2. Integrate castle/player/wallet summaries. 3. Add navigation to duel, workshop, defense and missions. 4. Render loading/empty/offline/error. **Expected Files / Modules:** onboarding and home features. **API / Database Impact:** player/castle/wallet reads.

**Acceptance Criteria:** Server values are clearly distinguished from cached values; primary entry paths both visible; onboarding can be resumed safely. **Testing Requirements:** Screen states, localization and first-use E2E. **Completion Definition:** Design QA validates primary home on small/large screens.

## MOB-07: Implement Duel and Result User Flows
**Priority:** High · **Phase:** 4-6 · **Component:** Mobile · **Effort:** XL · **Parallelism:** UI can parallelize against frozen event schemas

**Description:** Implement preparation, question, combat and result experiences. **Business Requirements:** PRD §6-8; UI/UX §5.7-5.10. **Technical Requirements:** Lock answer after submission, server-deadline display, no answer leak, stateful reconnect, combat/results only from resolved data. **Dependencies:** MOB-05, BATTLE-04..BATTLE-10.

**Implementation Steps:** 1. Build preparation/readiness and rules view. 2. Build question state machine and display-only timer. 3. Render submitted/opponent/feedback/timeout/reconnect states. 4. Render combat, result breakdown, rewards/MMR and next action. **Expected Files / Modules:** `features/duel` screens and components. **API / Database Impact:** Duel REST and realtime contracts.

**Acceptance Criteria:** One submission per round UI; no local result calculation; active match survives background and reconnection; reward display waits for server confirmation. **Testing Requirements:** E2E two clients, timer boundary, duplicate submit, reconnect and accessibility. **Completion Definition:** Two physical devices complete and display matching server result.

## MOB-08: Implement Castle, Inventory and Economy Screens
**Priority:** High · **Phase:** 2,6 · **Component:** Mobile · **Effort:** L · **Parallelism:** Parallel after stable castle/economy API

**Description:** Implement workshop, defense and inventory flows with clear costs and safe confirmation. **Business Requirements:** PRD §6-8; UI/UX §5.11-5.13. **Technical Requirements:** Server-confirmed balances/state; show missing amount, durability, item eligibility/expiry; confirm costly actions. **Dependencies:** MOB-03, CASTLE-02..CASTLE-08, ECON-01..ECON-05.

**Implementation Steps:** 1. Add structure preview/cost/prerequisite panels. 2. Add build/upgrade/repair request state and idempotency. 3. Add defense history/shields. 4. Add categorized inventory and disabled-reason states. **Expected Files / Modules:** castle/workshop/defense/inventory features. **API / Database Impact:** `/castle*`, `/wallet*`, `/inventory*`.

**Acceptance Criteria:** Failed actions never optimistically spend; insufficient balance shows exact deficit; repeated taps do not duplicate spend; combat state reflects server response. **Testing Requirements:** UI mutation error, idempotent retry, stale version and accessibility tests. **Completion Definition:** Staging end-to-end castle upgrade/repair/inventory use validated.

## MOB-09: Add Localization, Analytics and Crash Reporting Integration
**Priority:** Medium · **Phase:** 2-8 · **Component:** Mobile · **Effort:** M · **Parallelism:** Parallel after privacy approval

**Description:** Add translation-key workflow and privacy-safe client observability. **Business Requirements:** PRD §1,9-10; UI/UX §7. **Technical Requirements:** Turkish-first, locale-aware dates/numbers/currencies, RTL readiness, consent and PII redaction. **Dependencies:** MOB-02, FOUND-02.

**Implementation Steps:** 1. Add localization resource workflow. 2. Remove hardcoded user-facing strings. 3. Integrate approved crash/analytics provider with consent policy. 4. Add release/version context and no-secret scrubbing. **Expected Files / Modules:** localization resources, telemetry adapter and privacy config. **API / Database Impact:** Optional locale/profile fields and analytics events.

**Acceptance Criteria:** Launch language defaults to Turkish; localization QA finds no unkeyed text in touched screens; telemetry excludes credentials, answers and private data. **Testing Requirements:** Locale/RTL formatting and telemetry redaction tests. **Completion Definition:** Privacy-reviewed staging telemetry and language QA sign-off.

## MOB-10: Implement Profile, League, Season and Mission Screens
**Priority:** High · **Phase:** 7 competitive P1 · **Component:** Mobile · **Effort:** L · **Parallelism:** Parallel after progression contracts

**Description:** Deliver player profile/statistics/history, current league, leaderboard, season tracks and mission/achievement experiences. **Business Requirements:** PRD §6 P1; UI/UX §5.14,5.15,5.18-5.19. **Technical Requirements:** Server-confirmed rank/progress/reward claims, stable leaderboard pagination, localized dates and clear premium-track labels. **Dependencies:** MOB-03, PROG-01, PROG-04..PROG-08.

**Implementation Steps:** 1. Build profile and match-history views. 2. Add league/MMR progress and leaderboard. 3. Add season dates/free-premium track. 4. Add daily/weekly/season missions and achievements. 5. Add claim loading/error/success states. **Expected Files / Modules:** profile, league, season, mission and achievement features. **API / Database Impact:** `/players/*`, `/leagues/*`, `/leaderboards/*`, `/seasons/*`, `/missions/*`, `/achievements/*`.

**Acceptance Criteria:** Client never predicts rank/reward eligibility; claims refresh authoritative state; paid track is clearly distinguished; pagination and empty/error/offline states work. **Testing Requirements:** Screen contract tests, claim retry, season boundary, large text and localization. **Completion Definition:** Staging progression E2E and UX review pass.

## MOB-11: Implement Friends, Notifications, Store and Settings Screens
**Priority:** High · **Phase:** 5-8 · **Component:** Mobile · **Effort:** L · **Parallelism:** Parallel after social/purchase contracts

**Description:** Complete social discovery, friend requests, block/report, invitation inbox, notifications, store, purchase states and account preferences. **Business Requirements:** PRD §6 P1; UI/UX §5.16-5.21. **Technical Requirements:** Privacy-aware presence, authorized deep links, native purchase sheet, pending/refunded/restored states, notification preferences, account deletion/support. **Dependencies:** MOB-03, SOC-01..SOC-07, ECON-08, ECON-09, BACK-03.

**Implementation Steps:** 1. Build friend search/request/list/profile actions. 2. Add invitation and notification inbox/deep links. 3. Add store catalog and native purchase lifecycle. 4. Add language/audio/haptics/privacy/settings/support and deletion. 5. Cover offline, pending, failure and restoration states. **Expected Files / Modules:** friends, notification, store, settings and account features. **API / Database Impact:** social, notification, store/purchase and account endpoints.

**Acceptance Criteria:** Block/report and privacy rules are visible; purchases are not shown as complete until server verifies; refund/restore states reconcile; destructive account actions require confirmation. **Testing Requirements:** Social privacy, deep-link authorization, purchase sandbox, deletion and preference tests. **Completion Definition:** Representative iOS/Android store and social E2E accepted.
