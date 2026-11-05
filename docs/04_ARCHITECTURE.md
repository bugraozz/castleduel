# KALE DÜELLOSU — TECHNICAL ARCHITECTURE

## 1. Architecture goals

Build a secure, maintainable, scalable MVP without premature microservice complexity.

Recommended initial deployment: modular Go backend, PostgreSQL as source of truth, Redis for ephemeral coordination/cache, and a WebSocket gateway within the same deployable backend or a separable process when load requires it.

## 2. Recommended stack

### Mobile
- React Native + Expo + TypeScript.
- Expo Router.
- TanStack Query or equivalent for server state.
- Zustand or equivalent for transient client state.
- Typed API client generated or maintained from OpenAPI.
- WebSocket client with reconnect and event sequencing.
- Secure token storage using platform-backed secure storage.

### Backend
- Go.
- REST API under `/api/v1`.
- WebSocket for live match and presence events.
- PostgreSQL.
- Redis.
- Background worker for notifications, season closure, purchase reconciliation, and analytics exports.

### Operations
- Docker/container deployment.
- TLS reverse proxy/load balancer.
- CI/CD.
- Centralized logs, metrics, traces, error reporting.
- Automated backups and restore tests.

## 3. Logical modules

1. Auth: identity, sessions, token rotation.
2. Player: profile, statistics, preferences.
3. Castle: castle state and structures.
4. Construction: build/upgrade validation.
5. Economy: wallets and immutable ledger.
6. Inventory: item ownership and equipment.
7. Question: question catalog and versioning.
8. Duel: match lifecycle and round state.
9. Realtime: WebSocket authentication and event delivery.
10. Matchmaking: queues, candidate selection, reservations.
11. Combat: ability validation and damage resolution.
12. League: MMR and league mapping.
13. Season: season state, tracks, leaderboard, rewards.
14. Social: friends, blocks, invitations, reports.
15. Notification: in-app and push dispatch.
16. Purchase: native transaction verification and entitlements.
17. Mission: mission progress and achievements.
18. Admin: roles, moderation, configuration, audit.
19. Analytics: event ingestion and product metrics.

Modules should have explicit interfaces and own their domain rules. Avoid circular dependencies.

## 4. Deployment topology

Initial topology:
- Mobile clients
- HTTPS load balancer/API
- WebSocket endpoint
- Go application instance(s)
- PostgreSQL primary
- Redis
- Background worker
- Push provider integration
- Apple/Google purchase verification integration
- Monitoring/logging stack

Redis is not the durable source of truth for player assets, balances, match results, or MMR. PostgreSQL remains authoritative.

## 5. Domain ownership

| Domain | Source of truth | Critical consistency |
|---|---|---|
| Account | PostgreSQL | Strong |
| Castle | PostgreSQL | Strong |
| Wallet/ledger | PostgreSQL | Strong |
| Inventory | PostgreSQL | Strong |
| Duel result | PostgreSQL | Strong |
| Live queue | Redis + durable ticket metadata | Atomic reservation |
| Live match state | Server memory/Redis with durable checkpoints | Recoverable |
| MMR | PostgreSQL | Strong |
| Leaderboard | Derived/cacheable | Rebuildable |
| Notifications | PostgreSQL/outbox + provider | At-least-once delivery |
| Analytics | Event pipeline/warehouse | Eventual |

## 6. Real-time match architecture

1. Client authenticates REST and WebSocket sessions.
2. Client requests matchmaking.
3. Matchmaking service creates ticket and adds eligible queue entry.
4. Candidate pair is reserved atomically.
5. Both clients receive match-found event.
6. Both confirm before deadline.
7. Duel service creates match record and immutable question sequence.
8. Server broadcasts round state and deadline.
9. Clients submit answers over authenticated WebSocket or REST fallback.
10. Server validates identity, round, deadline, and duplicate status.
11. Round resolves after both answers or deadline.
12. Server sends resolved answer and score state.
13. After final round, combat and result processing execute transactionally.
14. Result, MMR, castle damage, inventory, and reward ledger are persisted.
15. Clients receive final result and refresh authoritative state.

## 7. WebSocket reliability

Every event should include:
- `event_id`
- `event_type`
- `match_id` where relevant
- `sequence`
- `server_time`
- `payload`
- schema version

Clients must:
- Authenticate connection.
- Reconnect with backoff.
- Send last acknowledged sequence or state version.
- Ignore duplicate event IDs.
- Detect sequence gaps and request state snapshot.
- Never assume local timer determines official deadline.

Server must:
- Validate ownership of match.
- Enforce event order.
- Rate-limit submissions.
- Persist critical transitions.
- Publish outbox events after transaction commit.

## 8. Transaction and idempotency strategy

Use idempotency keys for:
- Match confirmation.
- Answer submission.
- Ability use.
- Reward claim.
- Construction/upgrade.
- Repair.
- Purchase fulfillment.

For a reward or economy operation:
1. Begin transaction.
2. Check idempotency record.
3. Validate match/result and eligibility.
4. Insert ledger entries.
5. Update balance/materialized wallet.
6. Insert inventory/reward changes.
7. Insert outbox event.
8. Commit.

A duplicate request returns the original outcome without applying the effect again.

## 9. Security boundaries

- TLS in production.
- Access tokens short-lived; refresh rotation and revocation.
- Authorization checks on every resource.
- No secret or answer key in client bundle/payload.
- Server-generated match deadlines and timestamps.
- Input validation and rate limits.
- Admin MFA and least privilege.
- Audit logs for sensitive actions.
- Purchase verification server-side.
- Database credentials and provider secrets in secret manager/environment, never source control.

## 10. Scaling strategy

Start modular monolith. Scale stateless API/WebSocket instances horizontally.

When needed:
- Partition matchmaking by region/queue.
- Use Redis atomic scripts or equivalent distributed locking for ticket reservation.
- Separate WebSocket gateway.
- Separate background workers.
- Add read replicas for read-heavy reporting.
- Cache public profiles and leaderboards.
- Archive high-volume duel event data by retention policy.

Do not split services until throughput, team ownership, or reliability requirements justify it.

## 11. Observability

Metrics:
- API latency/error rate.
- WebSocket connections and reconnects.
- Matchmaking queue depth/wait time.
- Match confirmation failures.
- Match completion/abandonment.
- Duplicate/rejected answer counts.
- Reward transaction failures.
- Purchase verification latency/failure.
- Database pool saturation.
- Redis latency.
- Push delivery outcomes.

Logs must include correlation IDs and avoid passwords, tokens, answer keys, and unnecessary personal data.

## 12. Resilience and disaster recovery

- Automated PostgreSQL backups.
- Point-in-time recovery where supported.
- Restore drills.
- Redis data treated as ephemeral unless specific persistence is required.
- Durable match checkpoints.
- Outbox retry and dead-letter handling.
- Graceful shutdown and connection draining.
- Maintenance mode and client-compatible minimum version.

## 13. Environments

- Local development.
- Shared development.
- Staging.
- Production.

Each environment must have separate credentials, databases, payment configuration, push credentials, and analytics destinations.

## 14. Architecture decisions requiring validation

- WebSocket answer submission vs REST submission with WebSocket updates.
- Redis deployment and high availability.
- Regional topology and data residency.
- Push provider.
- Analytics warehouse.
- 3D castle rendering approach and asset pipeline.
- Hosting provider and cost envelope.
