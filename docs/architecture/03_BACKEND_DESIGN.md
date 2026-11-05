# Backend Design

## Runtime shape

One Go modular monolith exposes REST (`/api/v1`), authenticated WebSocket (`/api/v1/realtime`), and worker entry points. Transport handlers parse/validate, authenticate, apply request limits, call typed use cases, and map domain errors to stable API envelopes. Domain packages contain deterministic rules and do not depend on transport or external providers.

```text
cmd/api, cmd/worker
internal/
  auth player question duel realtime matchmaking social
  castle construction combat economy inventory progression league season
  purchase notification admin analytics
  platform/{postgres,redis,push,store,clock,telemetry}
  api/{http,websocket,contracts,middleware}
```

Folder names are a recommended starting point; modules may be combined where ownership remains clear. Do not create a package per endpoint.

## Module responsibilities and interfaces

| Module | Example use cases | Owns / emits |
|---|---|---|
| Auth | Register, login, refresh rotation, logout, deletion request | Users, identities, sessions; `account.*` |
| Player | Read/update profile, privacy and statistics | Profile/preferences; `player.*` |
| Question | Search/review/publish version, select match sequence | Question/version/category; `question.published` |
| Matchmaking | Join/cancel, candidate selection, reserve/confirm/decline | Tickets/reservations; `matchmaking.*` |
| Duel | Create, start, submit answer, resolve round, recover, finalize | Duel/players/rounds/answers/results; `duel.*` |
| Combat | Validate abilities, calculate damage/defense, apply event | Ability definitions/usage, damage event; `combat.*` |
| Castle | Read/build/upgrade/repair, protection and damage history | Castle/structures; `castle.*` |
| Economy / Inventory | Wallet, ledger, item ownership/equip/use, rewards | Wallet projection, ledgers, inventory; `economy.*` |
| Progression | XP, missions, achievements, MMR, league/season closure | Progression and claims; `progression.*` |
| Social | Friend requests, blocks, invitations, reports | Social graph and reports; `social.*` |
| Purchase | Catalog, verify, fulfill, restore, revoke/refund | Purchases and entitlements; `purchase.*` |
| Admin | Moderation, content, config, inspection, audited corrections | Admin roles, config versions, audit logs; `admin.*` |
| Notification | Inbox, preference checks, push delivery/retry | Notifications/outbox delivery state |

Cross-module calls should be typed interfaces in the application layer. Event consumers handle side effects asynchronously only when eventual consistency is acceptable. Combat and final match transactions must not rely on asynchronous delivery for a critical balance or durability result.

## Request lifecycle

1. Attach request ID and trace context; enforce body size, origin, TLS, and rate limits.
2. Authenticate; resolve actor and enforce object-level authorization.
3. Decode a versioned contract and validate syntactic constraints.
4. Call an application use case with server-derived identity and a scoped idempotency key.
5. Execute domain checks and transaction with constraints/locks.
6. Persist outbox event in the same commit when external effects are needed.
7. Map expected domain failures to stable error codes; log unexpected failures with correlation ID and no secrets.
8. Return authoritative state/version; never return a success-shaped fallback.

## Transaction and concurrency rules

- Use database uniqueness and foreign-key constraints as final protection, not only preflight checks.
- Lock the smallest required rows or use optimistic `version` checks. Return `CONFLICT` on stale expected state.
- Matchmaking reservation must lock/atomically claim both tickets and ensure neither player has a second active reservation or duel.
- A player answer is unique by `(duel_question_id, user_id)`; compare server receive time to stored deadline under transaction.
- Finalization has one result row and one unique operation ID; rating, rewards, castle damage, inventory effects and outbox insert commit once. A duplicate request returns the stored result.
- Build/upgrade/repair debit and castle mutation commit together.
- Item consumption and effect record commit together; reject duplicate activation without consuming twice.
- Purchase transaction ID is unique per platform; verified entitlement grant and purchase record commit together.
- Ledger entries are immutable. Update a wallet projection in the same transaction, and reconcile projection against ledger asynchronously with alerting.

## Idempotency convention

For harmful-to-repeat mutations, require `Idempotency-Key` (REST) or `command_id` (WebSocket). Scope uniqueness by authenticated user, operation type, and key; persist request hash and outcome. Same key/same request returns original outcome; same key/different request returns `IDEMPOTENCY_KEY_REUSED`. Retain records at least as long as retries/replay are possible, then retain a compact dedupe tombstone under policy.

## Error handling

Envelope: `{ "error": { "code": "...", "message": "...", "details": {}, "request_id": "..." } }`.
Stable codes include `UNAUTHORIZED`, `FORBIDDEN`, `VALIDATION_ERROR`, `RESOURCE_NOT_FOUND`, `CONFLICT`, `ALREADY_IN_MATCH`, `MATCH_NOT_READY`, `ANSWER_TOO_LATE`, `ANSWER_ALREADY_SUBMITTED`, `INSUFFICIENT_RESOURCES`, `ITEM_NOT_AVAILABLE`, `PURCHASE_VERIFICATION_FAILED`, `RATE_LIMITED`, `SERVER_MAINTENANCE`, `IDEMPOTENCY_KEY_REUSED`, `SEQUENCE_GAP`, and `INTERNAL_ERROR`. Do not expose secrets, answer keys, SQL errors, or private player details.

## Background jobs

Run outbox publishing, push delivery, purchase reconciliation, season close/final reward, leaderboard projection, mission/achievement event processing, wallet reconciliation, expired invitation/ticket cleanup, retention/partition maintenance, and analytics export as observable idempotent jobs. Each job needs bounded retry, backoff, dead-letter/alert policy, and safe replay. Season closure is resumable and idempotent.

## Security boundaries

Use short-lived access credentials and rotating/revocable refresh sessions; password hashing with a modern memory-hard algorithm where passwords are supported; least-privilege DB credentials; per-object checks; CSRF protections for browser-admin flows; admin MFA; strict secret management; secure headers; request and WebSocket message size limits; rate limits and abuse monitoring. Admin corrections require role, reason, before/after record, and audit entry in one transaction.

## Test seams

Inject clock, ID generator, repositories, transaction runner, queue coordinator, and provider clients. Unit-test pure scoring, rating, damage, costs, and eligibility. Use PostgreSQL/Redis integration tests for locking, uniqueness, transactions, and queue recovery. Contract tests validate OpenAPI and WebSocket schemas. Do not make correctness tests depend on real store or push providers.

## REST endpoint registry

All routes below are under `/api/v1` unless marked `/admin/v1`. Protected means a valid player bearer access credential plus object-level authorization; admin routes require a separate admin session, MFA and RBAC. Mutations with financial/gameplay impact require `Idempotency-Key`. Apply route-specific rate limits, bounded pagination and stable error envelopes. Exact field schemas belong in OpenAPI; the minimum payload and validation are summarized here.

| Method and route | Auth | Request / response | Validation, error and retry contract |
|---|---|---|---|
| `POST /auth/register` | Public | Email/password/username/locale → player and session | Normalized uniqueness, password policy; validation/conflict/rate-limited; registration throttle |
| `POST /auth/login` | Public | Credentials → access/refresh tokens and minimal player | Generic invalid credentials; no enumeration; per-account/IP throttle |
| `POST /auth/refresh`, `POST /auth/logout` | Refresh credential / protected | Refresh → rotated session; logout → revoked state | Previous refresh token is single-use; replay/expiry denied; logout idempotent |
| `POST /auth/password/forgot`, `POST /auth/password/reset` | Public / one-time reset token | Email or reset token/new password → generic accepted/completion | Enumeration-safe response, expiry and strength checks; one-time token and strict throttle |
| `DELETE /account` | Protected | Confirm deletion → request/status | Recent auth, approved retention and active-match policy; idempotent and throttled |
| `GET/PATCH /players/me`, `GET /players/me/statistics`, `GET /players/{playerId}/public` | Protected or public by policy | Read/permitted profile fields → safe profile/stats | Allowlisted fields, privacy and moderation filtering; object scope |
| `GET /castle`, `GET /castle/structures`, `GET /castle/damage-history` | Protected | Optional history cursor → authoritative state | Owner scope; stable bounded cursor |
| `POST /castle/structures`, `POST /castle/structures/{id}/upgrade` | Protected | Structure/slot or expected version → state | Unlock, slot, balance and version; idempotency required |
| `POST /castle/repair`, `POST /castle/structures/{id}/repair` | Protected | Target, durability/item choice, expected version → cost/result | Ownership, max durability, resources/item; idempotency required |
| `POST/GET/DELETE /matchmaking/queue` | Protected | Mode/region or ticket → trusted ticket state | Server derives rating/power; one ticket/player; join/cancel idempotent and rate-limited |
| `POST /matchmaking/matches/{id}/confirm`, `.../decline` | Protected participant | Command ID → confirmation state | Participant, reservation and deadline; idempotent |
| `GET /friends`, `GET /friends/requests`, `GET /players/search?q=` | Protected | Query/cursor → privacy-filtered page | Input bounds, privacy/block, stable pagination; search throttled |
| `POST /friends/requests`, `POST /friends/requests/{id}/accept`, `.../reject`, `DELETE /friends/{id}` | Protected | Recipient/decision → friendship state | Duplicate/self/block/privacy checks; idempotent |
| `POST/DELETE /players/{id}/block` | Protected | Target → block state | Self/block policy; no sensitive disclosure; idempotent |
| `POST /duel-invitations`, `POST /duel-invitations/{id}/accept`, `.../reject` | Protected | Recipient/mode → invitation or duel status | Friend/block/expiry/eligibility; shared match engine; idempotent and spam-limited |
| `GET /duels/{id}`, `GET /duels/{id}/state` | Protected participant | Optional last sequence → authorized state/snapshot | Membership, sequence bounds; read throttled |
| `POST /duels/{id}/answers` | Protected participant | Round/question/answer/command ID → accepted state | Membership, current round, server deadline, schema and uniqueness; strict per-round rate |
| `POST /duels/{id}/abilities` | Protected participant | Ability/target/command ID → resolved activation | Match snapshot, phase, ownership, charges and target; idempotent |
| `GET /duels/{id}/result`, `POST /duels/{id}/rematch` | Protected participant | None or eligible mode → immutable result/rematch status | Participant, terminal state and rematch policy; rematch idempotent |
| `GET /inventory`, `POST /inventory/equip`, `POST /inventory/unequip`, `POST /inventory/use` | Protected | Item/slot/target/command → server inventory/effect state | Ownership, expiry, loadout, quantity and phase; mutations idempotent |
| `GET /wallet`, `GET /wallet/transactions`, `GET /rewards/{id}`, `POST /rewards/{id}/claim` | Protected | Cursor/source/claim key → balance, history or reward | Owner scope, eligibility; claim idempotent |
| `GET /leagues/me`, `GET /leaderboards/global`, `GET /leaderboards/friends`, `GET /leaderboards/seasons/{id}` | Protected | Cursor/scope → stable ranked page | Validated data, tie-break and privacy; bounded reads |
| `GET /seasons/current`, `GET /seasons/{id}/rewards`, `POST /seasons/{id}/rewards/{id}/claim` | Protected | None or claim key → season/reward state | Season window, entitlement and eligibility; claim idempotent |
| `GET /store/products`, `POST /purchases/verify`, `GET /purchases`, `POST /purchases/restore` | Protected | Platform transaction/product → verification/entitlement status | Server platform verification, transaction uniqueness and refund state; retry-safe |
| `GET /missions`, `POST /missions/{id}/claim`, `GET /achievements`, `POST /achievements/{id}/claim` | Protected | None or claim key → progress/result | Validated event progress and eligibility; claim idempotent |
| `GET /notifications`, `POST /notifications/{id}/read`, `GET/PATCH /notification-preferences` | Protected | Cursor/updates → inbox/preferences | Owner scope and preference validation; read mutation idempotent |
| `POST /players/{id}/reports` | Protected | Reason and optional match/context → report ID/status | Valid reason, privacy and strict abuse limit |

### Admin endpoint family

Use `/admin/v1` and role-specific scopes for question draft/review/publish/import; player and report search/moderation; match inspection; economy/purchase read and compensating actions; game configuration draft/validate/activate/rollback; season create/preview/close/retry; audit and aggregated dashboards. All writes include reason and optimistic version; sensitive operations require explicit confirmation/approval and append-only audit. No endpoint performs untracked direct balance or result edits. Exact request/response fields must be specified in OpenAPI before frontend work.

| Method and route | Auth | Request → response | Validation / audit |
|---|---|---|---|
| `GET/POST /admin/v1/questions` | Designer/reviewer | Filter or draft content → page/draft | Draft validation; creation audited |
| `PATCH /admin/v1/questions/{id}` | Designer | Changed draft fields and expected version → version | Published content immutable; stale version conflict |
| `POST /admin/v1/questions/{id}/publish` | Publisher role | `{reason,expected_version}` → published version/audit ID | Full content QA; idempotent and audited |
| `POST /admin/v1/questions/import` | Designer | Validated batch/file reference → per-row report | Size/type/rules; partial vs atomic behavior explicit |
| `GET /admin/v1/players`, `GET /admin/v1/players/{id}` | Support/moderator | Search/filter → redacted profile | Role-specific fields and search rate limit |
| `POST /admin/v1/players/{id}/restrictions` | Moderator | `{action,reason,expires_at,expected_version}` → action/audit ID | Allowed action, expiry and active-match policy; idempotent |
| `GET /admin/v1/reports` | Moderator | Status/cursor → report page | Reporter/evidence access restricted |
| `GET /admin/v1/matches/{id}` | Support/moderator | Match ID → authorized event timeline | No answer-key disclosure before resolution; read-only by default |
| `GET /admin/v1/players/{id}/economy`, `GET /admin/v1/purchases/{id}` | Support Agent with economy-read permission | Player/purchase ID → ledger/verification state | Redacted read; no direct balance mutation |
| `POST /admin/v1/players/{id}/economy-adjustments` | Super Admin or explicitly delegated Support Agent permission | `{currency,amount,reason,reference,expected_version}` → transaction/audit ID | Compensating ledger entry only; approval and idempotency required |
| `GET/POST /admin/v1/configurations` | Designer/admin | Read or create draft config → version | Typed schema and range validation |
| `POST /admin/v1/configurations/{id}/validate`, `.../activate`, `.../rollback` | Designer/admin | `{reason,expected_version,approval_id?}` → validation/version/audit | Simulation and approval as policy; activation versioned and idempotent |
| `GET/POST /admin/v1/seasons`, `POST /admin/v1/seasons/{id}/close` | Season admin | Season config/action → season or job status | Date/reward validation; closure resumable and idempotent |
| `GET /admin/v1/audit-logs`, `GET /admin/v1/metrics` | Auditor/analyst | Filter/cursor/window → audit page or aggregates | Least privilege, privacy filters, bounded query |

### Minimum JSON payload and response shapes

These are the initial typed contract shapes; IDs are UUID strings and timestamps are ISO-8601 UTC. Pagination uses `{items, next_cursor}` and a bounded `limit`. Error responses use the common envelope above.

| Operation family | Request payload | Success response |
|---|---|---|
| Register | `{email, password, username, locale}` | `{access_token, expires_in, refresh_token, player:{id,username,player_level}}` |
| Login / refresh | Login: `{email, password}`; refresh: `{refresh_token}` | `{access_token, expires_in, refresh_token, player}` |
| Profile patch | `{username?, avatar_key?, bio?, locale?, privacy_settings?}` | `{player:{id,username,player_level,locale,...permitted_fields}}` |
| Castle build | `{structure_key, slot_key, expected_version}` | `{structure:{id,key,level,durability,max_durability}, castle_version, wallet}` |
| Castle upgrade / repair | `{expected_version, target_durability?}` or `{expected_version, repair_item_id?}` | `{castle, structure?, durability_restored, cost, wallet, operation_id}` |
| Queue join / status | Join: `{mode:"ranked", region:"auto"}`; status/cancel uses ticket ID | `{ticket_id,status,started_at,search_stage?}`; server derives rating/power |
| Match confirmation | `{command_id}` on confirm/decline route | `{match_id,status,confirmation_deadline,duel_id?}` |
| Friend request | `{recipient_id}`; accept/reject has no body or `{command_id}` | `{request_id,status,friendship?}` |
| Block/unblock | `{target_id}` for block; target in route for delete | `{target_id,blocked}` |
| Duel invitation | `{recipient_id,mode:"friendly"}` | `{invitation_id,status,expires_at}`; accept returns `{invitation_id,status,duel_id?}` |
| Answer submission | `{round_number,question_id,answer:{...},client_submission_id}` | `{accepted:true,round_number,received_at,state_version}`; no correctness until resolution |
| Ability activation | `{ability_key,target?,command_id}` | `{accepted:true,usage_id,effect_status,state_version}` |
| Duel state/result | State query `{last_sequence?,state_version?}`; result query has no body | `{state,state_version,sequence}` or immutable `{result,score,damage,rewards,rating_change}` |
| Inventory equip/use | `{item_id,slot_type,slot_number}` or `{item_id,target?,command_id}` | `{inventory,loadout?,effect?,state_version}` |
| Reward/mission/achievement claim | `{claim_idempotency_key}` or header key | `{claim_id,status,rewards,wallet,inventory}` |
| Leaderboard query | Cursor/limit query params; no mutation body | `{items:[{rank,player_id,username,rating}],next_cursor}` |
| Purchase verify/restore | `{platform,product_id,signed_transaction}`; restore includes signed transaction set/cursor | `{purchase_id,status,entitlements,verified_at}`; pending stays pending |
| Notification read/preferences | Read: no body; preferences patch `{push_enabled, categories, quiet_hours}` | `{notification_id,read_at}` or `{preferences}` |
| Player report | `{reason,details?,match_id?}` | `{report_id,status:"received"}`; reporter identity is not exposed |
| Admin mutation | Resource-specific typed fields plus `{reason,expected_version,approval_id?}` | `{resource,version,audit_id,status}`; no generic untyped arbitrary write |

### Validation and error behavior by operation family

- **Identity/profile:** `VALIDATION_ERROR`, `CONFLICT`, `UNAUTHORIZED`, `RATE_LIMITED`; normalize usernames on server and redact account existence from recovery/login responses.
- **Castle/economy/inventory:** `FORBIDDEN`, `RESOURCE_NOT_FOUND`, `CONFLICT`, `INSUFFICIENT_RESOURCES`, `ITEM_NOT_AVAILABLE`, `IDEMPOTENCY_KEY_REUSED`; no partial mutation.
- **Matchmaking/social:** `ALREADY_IN_MATCH`, `MATCH_NOT_READY`, `CONFLICT`, `RATE_LIMITED`, `ACCOUNT_RESTRICTED`; never expose blocked-user private state.
- **Duel/realtime:** `FORBIDDEN`, `ANSWER_TOO_LATE`, `ANSWER_ALREADY_SUBMITTED`, `SEQUENCE_GAP`, `MATCH_NOT_READY`; include command/request ID and whether retry is safe.
- **Progression/store:** `PURCHASE_PENDING`, `PURCHASE_VERIFICATION_FAILED`, `CONFLICT`, `RESOURCE_NOT_FOUND`; never grant before verified eligibility.
- **Admin:** `FORBIDDEN`, `VALIDATION_ERROR`, `CONFLICT`, `APPROVAL_REQUIRED`; denied sensitive writes are logged without leaking private data.

Mutations are rejected if the idempotency key is absent where required or reused with a different normalized payload. Rate limits are configurable by route and actor; report, auth, search, invitation, answer, and purchase verification routes need stricter limits than reads.

### Contract-wide error and retry rules

Use source error codes plus `IDEMPOTENCY_KEY_REUSED`, `SEQUENCE_GAP`, `ACCOUNT_RESTRICTED`, and `PURCHASE_PENDING`. A client may retry only safe reads or mutations with a stable idempotency key. Return the original result for the same key and payload; reject the same key with changed payload. Every error includes `request_id`; details must be safe and field-specific.
