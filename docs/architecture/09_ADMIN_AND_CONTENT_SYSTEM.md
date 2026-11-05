# Admin and Content Management

## Architecture

Provide a separate authenticated admin web client and privileged API surface backed by the same domain modules, never direct database writes. All writes require authorization and audited use cases. Require MFA for privileged roles and separate admin credentials/session policy from player accounts. Restrict network/access as appropriate for the hosting provider.

## Roles and permissions

| Role | Allowed scope | Explicit limits |
|---|---|---|
| Super Admin | all admin capabilities, role/config emergency control | MFA, reason, full audit; break-glass use reviewed |
| Game Designer | question content, balance/config drafts, season/reward setup | cannot moderate accounts or grant arbitrary value unless separately authorized |
| Moderator | reports, block/suspend/restrict, safety notes | no economy/config changes |
| Support Agent | player/match/economy read inspection and approved support workflow | no silent edits; corrections use reasoned compensating actions if granted |
| Analyst | aggregated analytics and leaderboard inspection | no private account fields or mutation |

Permission checks are action- and object-scoped, not just role labels. Role assignment/revocation is audited. Require approval or dual control for sensitive production economic/config operations when operationally feasible.

## Content lifecycle

Question editor supports category, difficulty, type, language, prompt, options, answer key, explanation, media, accessibility alternative, draft/review/publish/archive, duplicate detection, preview, bulk import validation, version history and publication audit. Store each published localized version immutably; active matches retain their version snapshot. Only authorized reviewers can publish; define two-person review for high-risk/bulk changes.

Question QA checks correct answer, ambiguity, locale/translation, media rights, accessibility, duplicate exposure, category/difficulty and answer-key access. Content access/logging must not expose keys to client analytics or broad admin roles.

## Player, match and safety tools

Search by ID/normalized username; inspect public-safe profile, match timeline, connection/answer/result events, reports, restrictions, and economy/transaction history under role-based redaction. Moderator workflow supports report triage, block/suspend/restrict, expiry/reason, evidence and appeal/escalation fields. Support notes are access-controlled and retained under policy. Never permit untracked result or balance editing.

## Configuration and economy operations

Version game configuration, reward tables, league thresholds, matchmaking expansion, scoring, damage, abilities, repair costs, season dates and catalog. Validate constraints and simulate preview before activation. Activation is staged and timestamped; active duels keep rules snapshot. Rollback activates prior safe version rather than rewriting history. Any grant/reversal is represented by a ledger or entitlement operation with reference, reason, idempotency, actor and audit.

## Audit and reporting

Record admin actor, role, action, target, before/after, reason, timestamp, request ID, and approval where applicable. Audit rows are append-only and access-controlled. Provide filters for player reports, suspicious economy flows, match inspection, purchase/refund status, leaderboard, season reward status, configuration history, and operational health. Aggregate analytics must avoid unnecessary personal data.

## Admin MVP acceptance

The P1 admin panel supports authenticated least-privilege access, question lifecycle, player moderation, match inspection, economy/configuration and season tools, audit logs, and analytics visibility. Validate denial cases for every role, publish/version behavior, audit integrity, and idempotent season/reward actions.
