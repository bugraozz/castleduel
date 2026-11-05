# Security and Anti-Cheat

## Threat model

Protect player accounts, private data, fair match outcomes, durable castles, inventory/currency, purchases and privileged administration. Assume hostile or modified clients, automation, stolen credentials, malicious insiders, replayed/duplicated messages, network delay, provider outages, and race conditions.

| Asset | Threat | Primary controls |
|---|---|---|
| Account/session | credential stuffing, token theft, enumeration | password hashing, throttling, refresh rotation/revocation, generic recovery responses, secure storage |
| Questions/results | answer-key extraction, fake answer/result | omit key until resolution; server selects versions, validates answer and determines result |
| Match | clock manipulation, replay, double assignment, automation | server deadlines/receive time, command IDs, constraints, rate/anomaly detection, atomic reservations |
| Wallet/inventory/castle | forged balances, duplicate spend/attack, race | server-only commands, ledger, transactions, row versions/locks, unique idempotency |
| Purchases | fake receipt, replay, refund abuse | Apple/Google server verification, unique transaction, reconciliation and revocation |
| Admin | privilege escalation, silent correction | MFA, RBAC, least privilege, reason-required audited actions, separate admin client |
| Service | API abuse, DoS, scraping | per-route/IP/account limits, body/message caps, query bounds, WAF/edge controls and alerts |

## Server-authoritative rules

Never trust client-supplied correctness, score, answer time, timer, round status, winner, damage, defense, rating, rewards, wallet, inventory, loadout eligibility, or purchase outcome. Server selects question IDs/versions and deadlines; only reveal answer key after round closes. Compute response duration from server receipt and authoritative start time under a documented latency policy. Snapshot rules, configuration, question versions and loadout at match start.

## Identity, authorization and privacy

Authenticate each REST call and WebSocket connection. Authorize resource ownership and match membership at every operation; do not rely on hidden IDs. Short-lived access tokens, rotating hashed refresh tokens, revocation on logout/password reset, and secure client storage. Avoid account enumeration during login/recovery. Public profiles exclude email, provider IDs, payment data and private identifiers. Apply user block/privacy settings to search, presence, invitations and notifications. Define deletion/anonymization and legally required records before launch.

## Replay, races and abuse

Use idempotency/command IDs scoped to actor and operation; same key with changed payload is rejected. Enforce one active queue ticket/reservation/duel per player with DB constraints and atomic reservation. Use unique answer per user/question, unique combat event, result, reward source, purchase transaction, mission claim and ledger operation. Apply rate limits and anomaly detection to answer rates, repeated opponents, friendly-match farming, invitation spam, account clusters, and resource flows. Anti-automation signals should trigger review or graduated controls, not silently rewrite valid outcomes.

## WebSocket hardening

TLS only in production; validate origin and protocol/client versions; authenticate before subscription; allowlist event types; cap payload/frame sizes and message rates; bind identity to connection; reject spoofed user IDs and unauthorized match IDs; validate schema/version/sequence/command IDs; rate-limit reconnects; close expired/revoked sessions; avoid answer keys in logs; verify event delivery authorization on every reconnect.

## Economy and purchase security

Ledger and wallet changes are server-side, transactional, nonnegative unless explicitly designed otherwise, and auditable. Client can request an action but cannot specify effects/costs. Consumable use validates ownership, expiry, loadout, target, match phase and charge before atomic consumption. Verify native purchases with platform servers; grant only after signed transaction validation and unique transaction persistence. Handle pending, duplicate, refund, revocation and restore states. Reconcile platform notifications idempotently.

## Admin controls

Roles: Super Admin, Game Designer, Moderator, Support Agent, Analyst, per roadmap. Enforce least privilege; MFA required for privileged users. Support/analyst access should be read-only or aggregated. High-risk moderation, economy, configuration and correction operations require reason, role-specific permission, before/after values, actor, timestamp, request ID, and immutable audit. Avoid direct balance editing; use compensating entries. Provide emergency access review and revoke.

## Security testing checklist

- Unauthorized access to another player's profile, castle, inventory, duel or purchase.
- Session refresh replay, rotation, logout revocation and account enumeration.
- Missing/early answer-key leak and modified-client answer/correctness/time injection.
- Client clock manipulation; late, duplicate, reordered, malformed and cross-match commands.
- Concurrent queue joins, cancellation versus reservation, duplicate active match, blocked-pair candidate.
- Duplicate finalization, rewards, ledger operations, item activation and castle attacks.
- Negative balance, concurrent repair/build/item consumption and stale version conflict.
- Forged/replayed store receipt, duplicate notification, refund/revoke and restore paths.
- WebSocket origin, event spoofing, subscription authorization, oversized payload, rate-limit and reconnect abuse.
- Admin role escalation, CSRF, MFA bypass, unreasoned correction and audit tampering.
- Dependency/container/secrets scanning, privacy-sensitive logging, backup access and restore security.

## Incident readiness

Alert on result/reward duplication, ledger-wallet mismatch, answer-key exposure suspicion, unusual admin actions, purchase verification anomalies and cheating clusters. Runbooks cover queue/realtime outages, reward incident containment, credential compromise, purchase outage and data restoration. Preserve forensic evidence with access controls and defined retention.
