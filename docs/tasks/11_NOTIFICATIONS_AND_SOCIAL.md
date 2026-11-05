# Notifications and Social Tasks

## SOC-01: Implement Player Search and Friend Requests
**Priority:** High · **Phase:** 5 Ranked and friend entry · **Component:** Backend / Mobile · **Effort:** M · **Parallelism:** Parallel with queue UI after privacy rules

**Description:** Support safe player discovery, incoming/outgoing friend requests, accept/reject/remove. **Business Requirements:** PRD §7; UI/UX §5.17; DB §10; API §6. **Technical Requirements:** Normalized query, pagination/rate limits, privacy filters, duplicate request prevention. **Dependencies:** PROG-01, BACK-05, FOUND-02.

**Implementation Steps:** 1. Add friendship/request schema and indexes. 2. Implement search/request state machine. 3. Apply rate limits and privacy. 4. Build friend screens and empty states. **Expected Files / Modules:** `internal/social`, mobile friends feature. **API / Database Impact:** `/players/search`, `/friends*`.

**Acceptance Criteria:** Duplicate requests are idempotent/conflict clearly; private/suspended users follow policy; pagination is stable. **Testing Requirements:** concurrent accept, duplicate, privacy and abuse tests. **Completion Definition:** Two accounts establish and remove friendship end-to-end.

## SOC-02: Implement Block and Privacy Enforcement
**Priority:** Critical · **Phase:** 5 · **Component:** Backend · **Effort:** M · **Parallelism:** Central matchmaking/social policy

**Description:** Allow block/unblock and expose a shared block policy for search, friend requests, invites, matchmaking and presence. **Business Requirements:** PRD §8; UI/UX §5.17; Security requirements. **Technical Requirements:** Directed block record, query filtering, race-safe candidate/invitation validation, privacy setting checks. **Dependencies:** SOC-01.

**Implementation Steps:** 1. Add blocked-user schema/routes. 2. Filter search/presence. 3. Reject invitations and candidate pairing. 4. Add client controls/confirmation. 5. Audit moderation interactions where required. **Expected Files / Modules:** social policy service and mobile controls. **API / Database Impact:** `/players/{id}/block` and unblock; blocked_users.

**Acceptance Criteria:** Blocked pair never receives invite or matchmaking assignment; concurrent block/match reservation has deterministic documented behavior. **Testing Requirements:** Every discovery/queue path and privacy regression test. **Completion Definition:** Security suite verifies block enforcement across modules.

## SOC-03: Implement Presence and Availability
**Priority:** Medium · **Phase:** 5-7 · **Component:** Backend / Realtime · **Effort:** M · **Parallelism:** Parallel with notification delivery

**Description:** Expose privacy-aware online/presence state for friend features. **Business Requirements:** PRD §5 P1; UI/UX §5.17; Architecture §3. **Technical Requirements:** Ephemeral Redis lease, heartbeat, stale expiry, no durable gameplay decision, user privacy preference. **Dependencies:** BACK-08, SOC-02.

**Implementation Steps:** 1. Define presence semantics/visibility. 2. Add heartbeat and TTL. 3. Publish friend updates. 4. Add preference and stale cleanup. **Expected Files / Modules:** presence adapter and friend UI. **API / Database Impact:** preference schema and realtime presence events.

**Acceptance Criteria:** Stale clients disappear within configured TTL; private user is not exposed; Redis loss does not damage persistent state. **Testing Requirements:** heartbeat expiry, privacy and redis restart tests. **Completion Definition:** Presence behavior validated with multiple devices.

## SOC-04: Implement Invitation Delivery and Expiry Notifications
**Priority:** High · **Phase:** 5 · **Component:** Backend / Realtime · **Effort:** M · **Parallelism:** After MATCH-07 contract

**Description:** Deliver invitation state in-app/realtime and notify recipient of expiry/acceptance. **Business Requirements:** PRD §7; UI/UX §5.17; API §6,12. **Technical Requirements:** Recipient authorization, event dedupe, expiry job, privacy preferences. **Dependencies:** MATCH-07, BACK-06, BACK-08.

**Implementation Steps:** 1. Add invitation event transitions. 2. Emit outbox/realtime events after commit. 3. Add expiry worker. 4. Add mobile incoming/outgoing invitation states. **Expected Files / Modules:** social event handlers and UI. **API / Database Impact:** invitation table and matchmaking events.

**Acceptance Criteria:** Duplicate notification does not duplicate invitation; expired invitation cannot start match; sender sees authoritative accepted/rejected result. **Testing Requirements:** event retry, expiry race, blocked recipient and privacy tests. **Completion Definition:** Friend challenge notification verified across two devices.

## SOC-05: Implement Player Reporting and Moderation Intake
**Priority:** High · **Phase:** 8 · **Component:** Backend / Mobile · **Effort:** M · **Parallelism:** Parallel with admin after report schema

**Description:** Allow users to report a player/match with reason and context. **Business Requirements:** PRD §5 P1; UI/UX §5.17; Roadmap §3. **Technical Requirements:** Rate limits, status workflow, evidence reference, privacy, no public exposure of reporter. **Dependencies:** BACK-05, BATTLE-08.

**Implementation Steps:** 1. Define report types and required fields. 2. Add create/status schema. 3. Add mobile report flow. 4. Route to admin triage. **Expected Files / Modules:** social report API, player_reports, mobile UI. **API / Database Impact:** report creation/read routes (recommended additions to source API).

**Acceptance Criteria:** Report is acknowledged without exposing reporter identity; duplicate/spam rates are limited; admin queue receives case. **Testing Requirements:** Authorization, validation, spam and privacy tests. **Completion Definition:** Full report-to-triage staging workflow passes.

## SOC-06: Implement Player Moderation Outcomes and Safety Controls
**Priority:** High · **Phase:** 8 · **Component:** Backend / Admin · **Effort:** L · **Parallelism:** Central with ADMIN-04

**Description:** Apply audited restrictions, suspensions, blocks and appeals/escalation. **Business Requirements:** PRD §8; Roadmap §3. **Technical Requirements:** Role authorization, reason, expiry, affected session revocation, matching/invitation eligibility and immutable audit. **Dependencies:** ADMIN-01, SOC-05.

**Implementation Steps:** 1. Model moderation action/status. 2. Apply restrictions in auth/social/matchmaking. 3. Revoke sessions when needed. 4. Add expiry/review job and appeal state. **Expected Files / Modules:** moderation service and admin workflow. **API / Database Impact:** user status/moderation/audit tables.

**Acceptance Criteria:** Suspended user cannot queue/invite; expiry is processed once; every action can be reviewed and appealed per policy. **Testing Requirements:** Cross-module enforcement, expiry, session revocation and audit tests. **Completion Definition:** Moderation scenario validated by QA and operations.

## SOC-07: Implement In-App and Push Notifications
**Priority:** High · **Phase:** 7 P1 · **Component:** Backend / Mobile · **Effort:** L · **Parallelism:** Push adapter can parallelize after outbox

**Description:** Provide notification inbox/preferences and push for invitations, match results, castle attacks and progression events. **Business Requirements:** PRD §5 P1; UI/UX §5.20-5.21; Roadmap Epic 8. **Technical Requirements:** Outbox-fed at-least-once delivery, dedupe, deep-link authorization, preference/quiet rules, provider tokens protected. **Dependencies:** BACK-06, BACK-09, MOB-01.

**Implementation Steps:** 1. Add notification/preference schema. 2. Implement inbox/read endpoints. 3. Integrate push adapter and token lifecycle. 4. Add dedupe/retry/dead-letter. 5. Implement mobile inbox and deep links. **Expected Files / Modules:** notification module, worker, mobile screen. **API / Database Impact:** `/notifications*`, `notification.received`, notification tables.

**Acceptance Criteria:** Duplicate delivery creates one inbox item; disabled preferences suppress appropriate push; opening deep link rechecks authorization. **Testing Requirements:** Provider outage, duplicate event, preference, token rotation and unauthorized deep-link tests. **Completion Definition:** Staging push and in-app notification flow accepted.
