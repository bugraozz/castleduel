# KALE DÜELLOSU — ROADMAP, ADMIN, QA AND RELEASE PLAN

## 1. Delivery strategy

Build a modular MVP in vertical slices. Each phase must produce an integrated, testable experience. Do not build all screens as static mockups before validating the real-time game loop.

## 2. Development epics

### Epic 1 — Foundation (P0)
- Repository and environment setup.
- Mobile app shell and navigation.
- Go backend skeleton.
- PostgreSQL schema and migrations.
- Authentication and session management.
- Player profile.
- Configuration and logging.

Acceptance:
- User can register/login.
- Profile loads from API.
- Environments are separated.
- Migrations run reproducibly.

### Epic 2 — Persistent castle and economy (P0)
- Starter castle creation.
- Structure definitions.
- Castle state API.
- Wallet and ledger.
- Construction and upgrade.
- Repair.
- Inventory basics.

Acceptance:
- Castle survives logout/reinstall.
- Build/upgrade/repair are server-validated.
- Duplicate requests do not duplicate deductions or rewards.

### Epic 3 — Question management (P0)
- Question schema and versioning.
- Category/difficulty.
- Admin create/edit/publish.
- Question selection service.
- Localization support.

Acceptance:
- Published question sequence can be assigned to a duel.
- Answer key is never exposed before resolution.
- Editing a question does not mutate an active match snapshot.

### Epic 4 — Duel engine and real-time gateway (P0)
- Duel lifecycle.
- WebSocket auth.
- Synchronized round delivery.
- Server deadlines.
- Answer submission.
- Score calculation.
- Disconnect/reconnect.
- Result persistence.

Acceptance:
- Two physical devices complete a synchronized match.
- Duplicate and late answers are rejected.
- Server determines score and winner.
- Reconnect restores state.

### Epic 5 — Matchmaking and friends (P0)
- MMR baseline.
- Queue tickets.
- Atomic reservations.
- Progressive search expansion.
- Match confirmation.
- Queue cancellation.
- Friend search/requests.
- Direct invitation.
- Block filtering.

Acceptance:
- Two players can be matched automatically.
- One player cannot be reserved twice.
- Cancellation races are safe.
- Friend challenge uses the same duel engine.
- Blocked players are never paired.

### Epic 6 — Combat and abilities (P0)
- Attack points.
- Damage calculation.
- Castle damage events.
- Basic bomb and shield.
- Repair kit.
- Loadout validation.
- Combat result animation events.

Acceptance:
- Damage is persistent and transactionally applied.
- Defense affects damage according to server configuration.
- Abilities cannot be duplicated or activated outside valid windows.

### Epic 7 — Leagues and seasons (P1)
- Rating calculation.
- League mapping.
- Leaderboards.
- Season lifecycle.
- Free/premium reward tracks.
- Season reward claims.

Acceptance:
- Only eligible ranked matches affect MMR.
- Leaderboards are derived from validated data.
- Season closure is repeatable/idempotent.
- Castle progression remains intact after reset.

### Epic 8 — Missions, achievements, notifications (P1)
- Mission definitions.
- Progress from server events.
- Achievement unlocks.
- In-app notifications.
- Push notification dispatch.

Acceptance:
- Duplicate events do not duplicate progress/rewards.
- Notification preferences are respected.

### Epic 9 — Store and payments (P1)
- Catalog.
- Native purchase flow.
- Server verification.
- Entitlements.
- Refund/revocation reconciliation.
- Purchase restoration.

Acceptance:
- Invalid receipts grant nothing.
- Duplicate transactions grant once.
- Refund/revocation is reflected correctly.

### Epic 10 — Admin and operations (P1)
- Admin authentication/roles.
- Question management.
- Game configuration.
- Player moderation.
- Economy tools.
- Season tools.
- Audit logs.
- Monitoring dashboards.

Acceptance:
- Sensitive actions are role-restricted and audited.
- Configuration changes are versioned and reversible where safe.

## 3. Admin panel specification

### Roles
- Super Admin: full access.
- Game Designer: questions, balance, configuration.
- Moderator: reports, player restrictions.
- Support Agent: read-only player/match inspection and approved support actions.
- Analyst: aggregated metrics only.

Use least privilege and MFA for privileged accounts.

### Dashboard
Metrics:
- Registered users.
- DAU/WAU/MAU.
- Concurrent players.
- Matches started/completed.
- Match abandonment.
- Matchmaking wait percentiles.
- MMR distribution.
- Economy source/sink.
- Purchases, refunds, revenue.
- Crash/error rate.

### Question manager
- Create/edit/archive.
- Category/difficulty/language.
- Question type.
- Answer options and answer key.
- Media upload/reference.
- Explanation.
- Preview.
- Draft/publish.
- Version history.
- Duplicate detection.
- Bulk import with validation.

### Game configuration
Editable, versioned values:
- Round count and timer.
- Scoring and bonuses.
- Damage coefficients.
- Defense caps.
- Ability definitions.
- Reward tables.
- Repair cost curves.
- Matchmaking expansion.
- League thresholds.
- Season dates.

Configuration changes must be validated, audited, and activated safely. Active matches use a rules snapshot.

### Player administration
- Search by player ID/username.
- Profile and match history.
- Reports.
- Inventory/economy inspection.
- Restriction/suspension.
- Support notes.
- Audit trail.

Admin tools must not silently edit balances or match outcomes. Corrections require reason and audit record.

### Season administration
- Create season.
- Configure dates.
- Configure reward tracks.
- Preview rewards.
- Activate/close season.
- Reconcile results.
- Monitor distribution.
- Re-run idempotent jobs safely.

## 4. QA strategy

### Unit tests
- Scoring and tie resolution.
- MMR updates.
- Matchmaking candidate rules.
- Damage and defense calculations.
- Repair costs.
- Ability limits.
- Reward generation.
- Ledger arithmetic.
- Season reset.

### Integration tests
- Authentication/session rotation.
- PostgreSQL transactions.
- Redis queue reservation.
- WebSocket event delivery.
- Duel finalization.
- Castle damage and repair.
- Mission reward claims.
- Purchase verification.

### End-to-end tests
- First registration to first duel.
- Automatic matchmaking.
- Friend challenge.
- Disconnect and reconnect.
- Build/upgrade/repair.
- Attack/defense.
- League promotion.
- Season reward claim.
- Purchase and restore.
- Account deletion.

### Load tests
Simulate:
- Concurrent queue joins.
- Large matchmaking bursts.
- Concurrent active duels.
- WebSocket reconnect storms.
- Result processing bursts.
- Leaderboard reads.
- Notification fan-out.

Define expected load profiles before interpreting performance results.

### Security tests
- Token theft/replay.
- Unauthorized player resource access.
- Answer key leakage.
- Client clock manipulation.
- Duplicate answer/reward.
- Matchmaking double reservation.
- Inventory/currency tampering.
- Purchase forgery.
- Admin privilege escalation.
- Rate-limit bypass.

### Device tests
Test representative supported iPhone and Android devices, OS versions, screen sizes, low-memory devices, and poor network conditions.

## 5. Release gates

Do not release until:
- No known critical/high security defects remain.
- Core multiplayer flow passes on real devices.
- Match results and rewards are transactionally safe.
- Account deletion and privacy disclosures are implemented.
- Purchase verification and restoration are tested if monetization is enabled.
- Crash reporting and monitoring are active.
- Backups and restore process are tested.
- App Store and Google Play metadata/assets are ready.
- Support and moderation processes are operational.
- Legal terms and privacy policy are reviewed for launch territories.

## 6. App Store / Google Play preparation

- App name and subtitle.
- App icon.
- Screenshots for supported device sizes.
- Preview video if used.
- Description and keywords.
- Age rating questionnaire.
- Privacy disclosures/data safety.
- Terms and privacy policy URLs.
- Support URL.
- In-app purchase metadata.
- Account deletion path.
- Review account and instructions if login is required.
- TestFlight/internal testing and staged rollout.

Review current platform policies before submission; requirements can change.

## 7. Operational runbooks

Prepare runbooks for:
- Matchmaking outage.
- WebSocket outage.
- Database failover.
- Redis failure.
- Purchase verification outage.
- Season closure failure.
- Reward duplication incident.
- Suspicious cheating wave.
- Push provider outage.
- Emergency game configuration rollback.

## 8. Post-launch priorities

Use analytics and player feedback to prioritize:
- Question quality and repetition.
- Matchmaking wait and fairness.
- Castle repair frustration.
- Economy inflation/deflation.
- Ability balance.
- Retention.
- Accessibility.
- Additional content.
- Social safety.

## 9. Definition of Done

A feature is done only when:
- Requirements and edge cases are documented.
- UX states are designed.
- Backend authorization and validation exist.
- Database migrations are included.
- API/WebSocket contracts are documented.
- Unit and integration tests pass.
- Analytics events are implemented where relevant.
- Errors and loading states are handled.
- Accessibility and localization are considered.
- Security implications are reviewed.
- Staging acceptance criteria pass.
- Operational ownership is clear.
