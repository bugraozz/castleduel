# Admin Panel and Content Management Tasks

## ADMIN-01: Implement Admin Identity, MFA and RBAC
**Priority:** Critical · **Phase:** 8 Store/admin/launch · **Component:** Admin / Backend · **Effort:** L · **Parallelism:** Central security owner

**Description:** Create privileged admin access independent from player sessions with least-privilege roles. **Business Requirements:** Roadmap §3; Architecture §9. **Technical Requirements:** Super Admin, Game Designer, Moderator, Support Agent, Analyst; MFA for privileged roles; action/object permission checks. **Dependencies:** BACK-02, FOUND-02.

**Implementation Steps:** 1. Model admin identity/role/permission. 2. Implement MFA and admin session policy. 3. Add authorization middleware. 4. Add role assignment/revocation audit. 5. Build denial tests. **Expected Files / Modules:** admin auth package, schema, admin app shell. **API / Database Impact:** `/admin/v1/*`; admin_users, role/permission tables.

**Acceptance Criteria:** Player token cannot access admin; every role has explicit allow/deny matrix; privileged action requires MFA. **Testing Requirements:** Role matrix, session expiry, escalation and revocation tests. **Completion Definition:** Security review approves permission matrix.

## ADMIN-02: Build Admin Application Shell and Audit Context
**Priority:** High · **Phase:** 8 · **Component:** Admin · **Effort:** M · **Parallelism:** UI can parallelize after ADMIN-01 contracts

**Description:** Create separate admin web client with safe navigation, session controls and audit-aware mutation UX. **Business Requirements:** Roadmap §3; UI needs in admin workflows. **Technical Requirements:** Typed API client, CSRF/session security, reason fields for sensitive actions, no direct DB access. **Dependencies:** ADMIN-01, FOUND-04.

**Implementation Steps:** 1. Create app shell and protected routes. 2. Add API/session layer. 3. Add role-aware navigation. 4. Require reason/confirmation for high-risk actions. 5. Add access logs and timeout. **Expected Files / Modules:** `apps/admin`, admin API client and UI tests. **API / Database Impact:** Admin API only.

**Acceptance Criteria:** Role sees only authorized routes; session timeout and CSRF controls work; all mutations include actor and audit context. **Testing Requirements:** E2E auth/role navigation and CSRF tests. **Completion Definition:** Admin shell passes security checklist.

## ADMIN-03: Implement Question Management and Publishing
**Priority:** High · **Phase:** 3 / 8 · **Component:** Admin / Backend · **Effort:** L · **Parallelism:** Content team after QUES-01

**Description:** Provide question creation, review, version, localization, media, preview and publication controls. **Business Requirements:** Roadmap §3; PRD §5 P1. **Technical Requirements:** Draft/published immutability, reviewer permissions, duplicate detection, import validation and audit. **Dependencies:** ADMIN-01, ADMIN-02, QUES-01, QUES-02.

**Implementation Steps:** 1. Build question list/detail/editor. 2. Add localization/category/difficulty and answer key form. 3. Add preview/version diff. 4. Add publish/archive flow and import validation. **Expected Files / Modules:** Admin question feature and question endpoints. **API / Database Impact:** Admin question routes, question_versions and audit logs.

**Acceptance Criteria:** Published content is immutable; only approved role publishes; import errors are reported without silent drops. **Testing Requirements:** Role, version, import and answer-key access tests. **Completion Definition:** Content owner completes representative publish/revision workflow.

## ADMIN-04: Implement Player Moderation and Report Triage
**Priority:** High · **Phase:** 8 · **Component:** Admin / Backend · **Effort:** L · **Parallelism:** Moderation workflow design can parallelize

**Description:** Provide report queue, player search, restriction/suspension, block context and escalation. **Business Requirements:** PRD §5 P1; UI/UX §5.17; Roadmap §3. **Technical Requirements:** Moderator-only actions, reason/expiry, evidence access control, appeal/escalation state and immutable audit. **Dependencies:** ADMIN-01, SOC-05, SOC-06.

**Implementation Steps:** 1. Build report queue and filters. 2. Add restricted player/match inspection. 3. Add suspend/restrict/resolve action with reason. 4. Add audit and review status. **Expected Files / Modules:** moderation feature and API. **API / Database Impact:** reports, moderation actions and admin audit.

**Acceptance Criteria:** Analyst/support cannot moderate; every action records actor/reason/target/time; private details are redacted by role. **Testing Requirements:** Permission, audit, expiry, appeal and privacy tests. **Completion Definition:** Moderation runbook and workflow acceptance.

## ADMIN-05: Implement Versioned Game and Economy Configuration
**Priority:** Critical · **Phase:** 8 · **Component:** Admin / Backend · **Effort:** L · **Parallelism:** Central config owner

**Description:** Manage scoring, timers, damage, defenses, abilities, rewards, costs, matchmaking, leagues and seasons safely. **Business Requirements:** GDD §2,7,9,11-13; Roadmap §3. **Technical Requirements:** Typed validation/ranges, draft/preview/activate, effective timestamp, immutable versions, match snapshot, rollback by new activation. **Dependencies:** ADMIN-01, FOUND-05, ECON-07.

**Implementation Steps:** 1. Define typed config schemas and constraints. 2. Build draft editor and validation. 3. Add simulation/preview. 4. Add approval/activation and audit. 5. Ensure active match snapshot remains immutable. **Expected Files / Modules:** config schema, admin UI, audit and runtime loader. **API / Database Impact:** `game_configurations` history and config admin routes.

**Acceptance Criteria:** Invalid values cannot activate; version history shows actor/reason; rollback does not rewrite active/historical match rules. **Testing Requirements:** Schema bounds, activation race, snapshot and rollback tests. **Completion Definition:** Staging config change and recovery verified.

## ADMIN-06: Implement Match Inspection and Support Read Views
**Priority:** High · **Phase:** 8 · **Component:** Admin / Backend · **Effort:** M · **Parallelism:** Parallel after BATTLE event schema

**Description:** Inspect duel lifecycle, rounds, connection state, answer timing, finalization and damage for reports/support. **Business Requirements:** Roadmap §3; PRD §5 P1. **Technical Requirements:** Read-only by default, answer-key access restricted, redacted player fields, filters and stable timeline. **Dependencies:** ADMIN-01, BATTLE-10, PROG-09.

**Implementation Steps:** 1. Add match search and participant filters. 2. Build ordered event/result timeline. 3. Add reward/purchase references. 4. Add audited correction workflow only if approved. **Expected Files / Modules:** admin match inspection API/UI. **API / Database Impact:** Read models for duel/event history.

**Acceptance Criteria:** Support can diagnose result without editing it; sensitive data is permission-gated; corrections are impossible without elevated permission/reason/audit. **Testing Requirements:** Role redaction, event ordering and access tests. **Completion Definition:** Support team completes simulated investigation.

## ADMIN-07: Implement Economy and Purchase Support Tools
**Priority:** High · **Phase:** 8 · **Component:** Admin / Backend · **Effort:** M · **Parallelism:** After ECON ledger and purchase schema

**Description:** Inspect wallet ledger, inventory, purchase verification/refund and reconciliation without silent balance edits. **Business Requirements:** Roadmap §3; Economy §5,13. **Technical Requirements:** Read-only by default; correction uses compensating transaction; purchase platform state and idempotent replay. **Dependencies:** ADMIN-01, ECON-01, ECON-09.

**Implementation Steps:** 1. Add player economy timeline. 2. Add purchase status/verification detail. 3. Add approved compensation action. 4. Require reason and audit; alert on mismatch. **Expected Files / Modules:** admin economy/purchase pages and endpoints. **API / Database Impact:** Ledger, purchase, entitlement, audit.

**Acceptance Criteria:** Admin cannot directly edit balance; compensation preserves history; duplicate fulfillment remains prevented. **Testing Requirements:** Role, compensating ledger, purchase replay and audit tests. **Completion Definition:** Finance/support workflow review completed.

## ADMIN-08: Implement Season, Reward and Operations Dashboard
**Priority:** High · **Phase:** 8 · **Component:** Admin / Backend · **Effort:** L · **Parallelism:** After PROG-05/06

**Description:** Create, preview, activate and close seasons; monitor distribution and job progress. **Business Requirements:** Roadmap §3; GDD §12; API §9. **Technical Requirements:** Idempotent close/replay, eligibility snapshot, approval, role gates, dashboards for match/economy/reliability. **Dependencies:** ADMIN-05, PROG-05, PROG-06, BACK-07.

**Implementation Steps:** 1. Build season config/preview. 2. Add activation/closure controls. 3. Show reward/job status and errors. 4. Add operational metric dashboard. 5. Add audited safe retry controls. **Expected Files / Modules:** admin season/ops features, worker controls. **API / Database Impact:** season state, reward claim and job records.

**Acceptance Criteria:** Season closure can be monitored/retried safely; metrics have defined source/time window; no unapproved reward grant path exists. **Testing Requirements:** Role, duplicate closure, failure recovery and dashboard query tests. **Completion Definition:** Staging season dry run and runbook reviewed.
