# Requirements Traceability Matrix

This matrix maps source-document requirements to architecture, persistence, API/events, tasks, and verification. Requirement IDs are stable references for planning; source section numbers follow the documents in `docs/`. Task IDs are defined in `docs/tasks/`. Test IDs refer to the end-to-end/concurrency/security scenarios in `docs/tasks/12_TESTING_SECURITY_AND_RELEASE.md`.

## Feature traceability

| Req ID | Source document / section | Requirement or feature | Architecture / data | API / realtime | Implementation tasks | Verification |
|---|---|---|---|---|---|---|
| R-01 | PRD §1-5; UI/UX §1-4 | iOS/Android, Turkish-first, localization-ready competitive castle game | Mobile architecture; profile locale | `/players/me`; localization resources | MOB-01, MOB-02, MOB-09 | T-01, T-20 |
| R-02 | PRD §5,8; UI/UX §5.1-5.3 | Register/login/refresh/logout/account deletion, safe profile | Auth/Player; users, identities, sessions | `/auth/*`, `DELETE /account` | BACK-02, BACK-03, PROG-01 | T-01, T-20 |
| R-03 | PRD §5,8; UI/UX §5.4,5.18 | Player profile, public profile/privacy, statistics | Player; profiles/statistics | `/players/me`, `/players/{id}/public` | PROG-01, PROG-02, MOB-10 | T-01 |
| R-04 | PRD §5,7-8; UI/UX §5.4,5.10-5.12 | Persistent castle, structures, upgrades, damage, defense, repair | Castle/Construction; castles/structures/damage/repairs | `/castle*`; `duel.castle_damaged` | CASTLE-01..CASTLE-08, MOB-08 | T-08, T-10 |
| R-05 | GDD §2-5; PRD §8 | Synchronized two-player timed questions, score, tie and result | Duel/Question; duels/rounds/answers/results | `/duels/*`; `duel.started`, `question_received`, `answer_submit`, `answer_resolved`, `finished` | QUES-01..QUES-05, BATTLE-01..BATTLE-08, MOB-07 | T-02, T-03, T-06, T-07 |
| R-06 | GDD §3; DB §5; API §7,12 | Identical immutable localized question/version; key withheld until resolution | Question; question_versions/duel_questions | question delivery and resolution events | QUES-01..QUES-04, BATTLE-04 | T-02 |
| R-07 | PRD §5,7-8; Roadmap §2 Epic 5 | Ranked automatic matchmaking at initial release, rating/level/fairness criteria | Matchmaking; tickets/reservations/history | `/matchmaking/*`; `matchmaking.*` | MATCH-01..MATCH-05 | T-11, T-12 |
| R-08 | PRD §5,7-8; UI/UX §5.5-5.7 | Friend search, social requests, direct invite and unranked duel at initial release | Social/Matchmaking; friends/requests/invitations | `/friends*`, `/players/search`, `/duel-invitations*` | SOC-01..SOC-04, MATCH-06, MATCH-07, MOB-11 | T-13 |
| R-09 | PRD §8; GDD §3-4; UI/UX §5.8 | Server deadlines, answer validation, scoring, speed rule, tie handling | Duel; answers/results/rules snapshot | answer command and resolve events | BATTLE-04..BATTLE-08 | T-02, T-03, T-06 |
| R-10 | PRD §8; Architecture §6-8; API §12 | Authenticated realtime, ordering, reconnect and state recovery | Realtime; event log/checkpoints/outbox | `/realtime`, ack/state sync and duel events | BACK-08, BATTLE-02, BATTLE-09..BATTLE-11 | T-04, T-05 |
| R-11 | PRD §8; GDD §5-9 | Attack points, damage, defense, protection, abilities, destruction recovery | Combat/Castle; ability and damage records | `/duels/{id}/abilities`; ability/attack/damage events | CASTLE-05..CASTLE-08, ECON-04 | T-08, T-10 |
| R-12 | PRD §8; GDD §10-12; API §9 | XP, castle level, MMR, leagues, seasons, leaderboards, rating history | Progression/League/Season; rating and season tables | `/leagues*`, `/leaderboards*`, `/seasons*` | PROG-02..PROG-06, MOB-10 | T-14 |
| R-13 | GDD §15; PRD §5 P1; UI/UX §5.19 | Daily/weekly/season missions, achievements and idempotent claims | Progression; mission/achievement tables | `/missions*`, `/achievements*` | PROG-07, PROG-08, MOB-10 | T-15 |
| R-14 | PRD §5,8; GDD §13-16; Economy §2-12 | Gold/Stone/Crystals, fair rewards, sources/sinks, free recovery | Economy; wallets/ledger/reward config | `/wallet*`, `/rewards*`; result events | ECON-01..ECON-05, PROG-03 | T-09 |
| R-15 | GDD §9; PRD §8 | Ability definitions, loadouts, ownership, consumption, limits, server effects | Combat/Inventory; definitions/loadouts/usage | `/inventory/equip`, `/inventory/use`, `duel.ability_activate` | ECON-04, CASTLE-07 | T-08 |
| R-16 | PRD §5 P1; GDD §16; Economy §8-13; UI/UX §5.15-5.16 | Catalog, cosmetics/pass, native purchases, verification, restore/refund/revocation | Purchase; products/purchases/entitlements | `/store/products`, `/purchases/*` | ECON-08..ECON-10, MOB-11 | T-16 |
| R-17 | PRD §5 P1; UI/UX §5.17; DB §10; API §6 | Friends, presence/privacy, block/report, invitation safety | Social; friendships/blocks/reports | `/friends*`, `/players/{id}/block`, reports, notifications | SOC-01..SOC-06, MOB-11 | T-13, T-20 |
| R-18 | PRD §5 P1; UI/UX §5.20-5.21 | In-app/push notifications and preferences | Notification; inbox/preferences/outbox | notification API and `notification.received` | SOC-07, MOB-11 | T-19 |
| R-19 | PRD §5 P1; Roadmap §3 | Admin roles, question management, moderation, config/economy/seasons, audit | Admin; RBAC/config/audit | `/admin/v1/*` (admin-only contract) | ADMIN-01..ADMIN-08 | T-17 |
| R-20 | PRD §5,9-10; Roadmap §4,7 | Analytics, crash reporting, observability, backups/restore, runbooks | Infra/Analytics | telemetry and operational endpoints | RELEASE-09..RELEASE-12 | T-18 |
| R-21 | PRD §8-9; UI/UX §6-7 | Accessibility, loading/error/offline states, localization, safe active-match UX | Mobile design system | API/event error and state recovery | MOB-02..MOB-11 | T-20 |
| R-22 | PRD §8; GDD §13-16; Economy §5,11 | Immutable auditable ledger, no duplicate/double spend, farming controls | Economy; currency/inventory transactions | mutation idempotency | ECON-02..ECON-07 | T-09, T-10 |
| R-23 | PRD §9; Architecture §11; Roadmap §4-5 | API/realtime latency, crash-free goals, load and device QA | Infrastructure/observability | measured contracts and load profiles | RELEASE-08..RELEASE-10 | T-18 |
| R-24 | PRD §2,8; GDD §16; Economy §10,14 | No gambling; fairness, paid-content boundaries, store policy | Purchase/combat policy | catalog/verification | ECON-08..ECON-10, RELEASE-13 | T-16 |
| R-25 | Roadmap §5-6; PRD §8 | Account deletion, privacy disclosures, store readiness and support | Auth/Admin/Operations | account deletion/admin operations | BACK-03, MOB-11, RELEASE-13..RELEASE-14 | T-20 |

## Cross-cutting acceptance gate

No requirement is complete based only on documentation. Every row must be linked to implemented task evidence, automated tests where feasible, and staging validation. Resolve all configurable decisions noted in `01_SYSTEM_ARCHITECTURE.md` before enabling affected ranked rules. If a source requirement changes, update the source reference, architecture contract, task ID, API/schema, and test mapping together.

## UI screen traceability

The UI/UX specification's screen inventory and explicit state requirements are covered as follows. Shared loading/empty/error/offline/success, accessibility and localization requirements apply to every applicable row, not only the screen where they are first mentioned.

| UI/UX source section | Screens / behavior | Architecture and data/API | Task IDs |
|---|---|---|---|
| §5.1-5.3 | Splash, onboarding, authentication | Mobile startup/auth; `/auth/*`, `/account` | MOB-01, MOB-04, BACK-02, BACK-03 |
| §5.4 | Home/Castle, resources, missions, league, duel CTAs | Mobile home; castle/wallet/profile/season APIs | MOB-06, MOB-08, MOB-10 |
| §5.5-5.7 | Duel selection, matchmaking, preparation/readiness | Matchmaking and invitation state | MATCH-01..MATCH-07, MOB-05, MOB-11 |
| §5.8-5.10 | Question duel, combat scene, results and rematch | Server-authoritative duel/realtime | BATTLE-01..BATTLE-11, MOB-07 |
| §5.11-5.13 | Workshop, defense, inventory | Castle, combat, economy and item ownership | MOB-08, CASTLE-01..CASTLE-09, ECON-01..ECON-04 |
| §5.14-5.15 | League, leaderboard, season/pass | Progression and season services | PROG-03..PROG-06, MOB-10 |
| §5.16 | Store and native purchase lifecycle | Store/purchase verification and entitlement | ECON-08..ECON-10, MOB-11 |
| §5.17 | Friends, presence, challenge, block and report | Social safety and direct invitations | SOC-01..SOC-06, MATCH-07, MOB-11 |
| §5.18-5.19 | Profile, statistics, match history, missions/achievements | Player/progression | PROG-01..PROG-09, MOB-10 |
| §5.20-5.21 | Notifications, settings, privacy, support, deletion, logout | Notification/auth/privacy | SOC-07, BACK-03, MOB-04, MOB-11 |
| §6-8 | Confirmation, optimistic UI limits, offline/error recovery, accessibility/localization | Mobile state/API/realtime | MOB-02..MOB-11, RELEASE-11 |
