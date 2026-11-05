# KALE DÜELLOSU — PRODUCT REQUIREMENTS DOCUMENT (PRD)

**Version:** 1.0  
**Status:** MVP baseline  
**Platforms:** iOS and Android  
**Product language at launch:** Turkish; architecture supports localization  
**Primary game mode:** Ranked two-player real-time trivia duel  
**Secondary game mode:** Unranked friend challenge

## 1. Product summary

Kale Düellosu is a competitive mobile game in which two players answer synchronized timed questions. Correct answers and match performance create attack opportunities and rewards. Winners can damage opponents' persistent castles; players build, upgrade, defend, and repair their own castles between matches. A rating-based matchmaking system, leagues, seasons, social features, missions, inventory, and optional purchases support long-term engagement.

The product must feel like a castle strategy game powered by knowledge duels—not a quiz app with decorative castle artwork.

## 2. Product goals

1. Deliver fair, low-latency two-player trivia matches.
2. Make every ranked match matter through persistent but recoverable castle consequences.
3. Give players meaningful construction, defensive, and offensive decisions.
4. Support both automatic matchmaking and friend challenges in the first public release.
5. Create long-term goals through player progression, leagues, seasons, missions, and achievements.
6. Monetize primarily through cosmetics and optional convenience while preserving competitive fairness.
7. Operate securely with server-authoritative gameplay and auditable economy transactions.

### Non-goals for MVP

- Real-money wagering, cash prizes, or gambling mechanics.
- Large-scale guild wars or clan territory systems.
- Open-world exploration.
- Player-to-player trading or a player-driven marketplace.
- User-generated questions without moderation.
- Cross-platform web/desktop gameplay.

## 3. Target users and personas

- **Casual competitor:** wants short, understandable matches and visible progress.
- **Trivia enthusiast:** values question variety, accuracy, and fair scoring.
- **Builder/collector:** wants to customize and improve a persistent castle.
- **Social player:** wants to challenge friends and compare profiles.
- **Ranked competitor:** wants transparent MMR, leagues, seasons, and leaderboards.

## 4. Product principles

- Correctness outranks speed: a fast incorrect answer never beats a correct answer.
- MMR is the primary matchmaking signal.
- Castle power is a secondary fairness signal, not a replacement for MMR.
- Losing must have consequences but must not erase long-term progression.
- Paid content must not guarantee victories.
- The server is authoritative for match state, scoring, damage, rewards, inventory, rating, and purchases.
- Permanent castle progression survives season resets.
- Every reward and purchase operation must be idempotent.

## 5. MVP scope

### P0 — playable competitive foundation
- Account registration, login, session refresh, account deletion.
- Player profile and public profile.
- Persistent castle, structures, durability, damage, and repairs.
- Question bank with categories, difficulty, localization, and versioning.
- Two-player real-time duel with synchronized questions.
- Server-side answer validation and scoring.
- Automatic rating-based matchmaking.
- Friend search, requests, invitations, and unranked friend matches.
- Combat resolution and basic attack/defense abilities.
- Inventory and resource ledger.
- Match result persistence and reconnect recovery.

### P1 — first public release completeness
- Leagues, MMR history, leaderboards, seasons, and rewards.
- Daily/weekly/season missions and achievements.
- Store catalog, cosmetics, season pass, and native purchase verification.
- Push and in-app notifications.
- Reporting, blocking, moderation, support tools.
- Admin panel for questions, game configuration, players, economy, and seasons.
- Analytics, crash reporting, monitoring, backups, and operational runbooks.

### P2 — post-launch
- Events and limited-time game modes.
- Additional question formats and languages.
- More castle themes, structures, and visual effects.
- Advanced social features and clubs.
- Match replays or detailed battle reports.
- More sophisticated anti-abuse models.

## 6. Navigation and screen inventory

Bottom navigation:
1. Home/Castle
2. Duel
3. League
4. Friends
5. Profile

Contextual screens:
- Onboarding and authentication
- Matchmaking search
- Duel preparation
- Question duel
- Combat/result
- Castle workshop
- Castle defense
- Inventory
- Store and season pass
- Missions and achievements
- Notifications
- Settings, support, privacy

### Screen requirements

| Screen | Purpose | Main elements | Key actions |
|---|---|---|---|
| Splash | Initialize app | Branding, update/session state | Continue |
| Onboarding | Explain core loop | Duel, castle, league walkthrough | Continue/skip |
| Authentication | Create/access account | Email/platform sign-in, recovery | Register/login |
| Home/Castle | Main hub | Castle, durability, resources, league, missions | Duel, build, repair |
| Duel selection | Select match type | Automatic match, friend challenge | Find/challenge |
| Matchmaking | Find fair opponent | Search state, timer, cancel | Cancel |
| Preparation | Confirm match | Both profiles/castles, readiness | Ready/decline |
| Question duel | Answer questions | Question, timer, options, status, score | Submit answer |
| Combat | Show battle outcome | Abilities, attack, shield, damage | Resolve/continue |
| Results | Explain outcome | Winner, stats, damage, rewards, rating | Rematch/home |
| Workshop | Develop castle | Structures, costs, previews | Build/upgrade |
| Defense | Protect/repair castle | Damage report, defenses, repair | Equip/repair |
| Inventory | Manage items | Categories, quantities, rarity | Equip/use |
| League | Competitive progression | Tier, MMR, standings | View leaderboard |
| Season | Seasonal progression | Dates, tracks, missions, rewards | Claim |
| Store | Purchases | Product catalog, price, details | Purchase |
| Friends | Social graph | Search, requests, list, status | Add/challenge/block |
| Profile | Identity/statistics | Record, achievements, castle | Inspect/edit |
| Missions | Retention objectives | Daily/weekly/season tasks | Claim |
| Notifications | Events and alerts | Read/unread list | Open/dismiss |
| Settings | Account/preferences | Audio, language, privacy, support | Update |

Every screen must include loading, empty, error, offline, and success states where relevant. Destructive actions require confirmation.

## 7. Core user journeys

### 7.1 First-time user
1. Launch app and accept required terms/privacy notices.
2. Register and verify account where required.
3. Select unique moderated username and avatar.
4. Complete a short tutorial duel against a controlled tutorial opponent or scripted training flow.
5. Receive starter castle, resources, and basic items.
6. Land on Home/Castle with clear next action.

### 7.2 Automatic ranked duel
1. Player selects Find Opponent.
2. Backend validates eligibility and creates one matchmaking ticket.
3. Queue searches compatible players by MMR and constraints.
4. Match reservation atomically locks both players.
5. Both players receive a match-found prompt and confirmation window.
6. When both confirm, server creates/starts duel and sends match state.
7. Players answer synchronized questions.
8. Server resolves score, combat, castle damage, rating, and rewards.
9. Result is persisted and delivered to both clients.
10. Home screen reflects updated castle, wallet, inventory, and rating.

### 7.3 Friend challenge
1. Player searches/adds friend or selects an existing friend.
2. Sender chooses friendly duel and sends invitation.
3. Recipient accepts before expiry.
4. Both enter same duel engine.
5. Default friend matches do not affect MMR or ranked leaderboard.
6. Match result is recorded in history; reward policy is configurable and abuse-limited.

### 7.4 Castle upgrade
1. Player opens Workshop.
2. Selects unlocked structure or upgrade.
3. Client previews costs and requirements.
4. Backend validates current balance, prerequisites, slots, and state version.
5. Transaction deducts resources and creates/upgrades structure.
6. Updated castle state is returned and visualized.

### 7.5 Attack and repair
1. Ranked match winner receives eligible attack outcome.
2. Combat service calculates damage after defenses and ability effects.
3. Castle and structure damage events are recorded transactionally.
4. Defender receives notification.
5. Defender opens Defense and repairs using valid resources/items.
6. Repair transaction updates durability and ledger.

## 8. Functional requirements

### Accounts and profiles
- Unique user ID and normalized username.
- Secure authentication and refresh-token rotation.
- User can view/edit permitted profile fields.
- Public profile excludes private account data.
- User can request account deletion and data handling according to policy.
- Moderation can suspend or restrict accounts with audited reason.

### Duel
- Only two eligible players per standard duel.
- Both receive identical question IDs and versions in identical round order.
- Question deadlines are server-controlled.
- One answer submission per player per round.
- Late, duplicate, malformed, or unauthorized answers are rejected.
- Match state is recoverable after reconnect.
- Final result is immutable except through audited administrative correction.

### Matchmaking
- Both automatic matchmaking and friend challenge are MVP requirements.
- A player may have at most one active queue ticket and one active duel.
- Queue reservations must be atomic.
- Blocked users must never be paired.
- Search widens gradually but respects configured maximum rating gap.
- Confirmation timeout returns unconfirmed players to eligible state according to policy.
- Queue cancellation must be race-safe.

### Castle
- Every account has one persistent primary castle in MVP.
- Castle and structures have current/max durability.
- Damage is persistent and repairable.
- Core castle progression cannot be permanently destroyed by combat.
- Construction and upgrades validate prerequisites and balances on server.
- Season resets never delete castle structures or permanent cosmetics.

### Economy and rewards
- Gold, Stone, and Crystals are separate balances.
- All balance changes use immutable transaction ledger entries.
- Rewards are calculated and granted server-side.
- Duplicate claims are prevented by unique idempotency keys.
- Losing players may receive limited participation rewards only when fair-play criteria are met.

### Abilities
- Offensive and defensive abilities have server-defined effect values, limits, target rules, and acquisition paths.
- Ability loadouts are validated before match start.
- Consumables are deducted only when successfully activated under valid conditions.
- No infinite stacking or client-defined effects.

### Competitive progression
- MMR is separate from player XP and castle level.
- Leagues map to rating thresholds and configurable promotion/demotion rules.
- Season reset may adjust rating but never removes permanent castle progress or purchased entitlements.
- Leaderboards use validated server data.

### Store
- Product catalog is remotely configurable.
- Native purchase verification is server-side.
- Refund/revocation events update entitlements.
- Cosmetics are permanent unless explicitly sold as time-limited access.
- Purchases never directly grant guaranteed competitive wins.

## 9. Non-functional requirements

- Secure transport (TLS) for all production traffic.
- p95 normal API latency target under 300 ms under agreed load profile.
- Real-time event p95 target under 250 ms in supported regions under normal conditions.
- Match state recovery after transient disconnect.
- 99.5% crash-free session target as an initial quality objective.
- Observability for API, WebSocket, matchmaking, economy, and purchase flows.
- Automated database backups and tested restore procedures.
- Accessibility: scalable text, sufficient contrast, reduced motion support.
- Localization-ready UI and content.
- No secrets embedded in the mobile bundle.

## 10. Product metrics

North-star candidate: completed fair duels per weekly active player.

Track:
- DAU/WAU/MAU.
- D1/D7/D30 retention.
- Match completion and abandonment.
- Matchmaking success and wait time percentiles.
- Duel frequency per active player.
- Question accuracy and response-time distribution.
- MMR distribution and match rating gap.
- Castle upgrade/repair rates.
- Resource sources and sinks.
- League progression.
- Store conversion, revenue, refunds.
- Crash-free sessions and API/WebSocket reliability.

## 11. Acceptance criteria for MVP

A release candidate is acceptable only when:
- Two independent devices can complete a synchronized duel.
- Automatic matchmaking and friend challenges both work.
- Match state survives supported reconnect scenarios.
- Server determines answer correctness, scores, winner, damage, rewards, and MMR.
- Castle damage persists across sessions and can be repaired.
- Construction and inventory changes are transactionally safe.
- Duplicate answer/reward/purchase requests do not duplicate effects.
- Leagues and seasons preserve permanent castle progress.
- Reporting/blocking and account deletion flows are available.
- Critical flows have automated tests and operational monitoring.

## 12. Open decisions to validate before production

These are configurable product decisions, not assumed final truths:
- Final duel length and question timer.
- Whether combat is automatic after each round or resolved after all questions.
- Exact damage formula and structure targeting.
- Whether friend matches grant any resources.
- Exact league thresholds and season duration.
- Resource reward amounts and construction costs.
- Ranked ability restrictions.
- Regional matchmaking policy.
- Age rating and parental/child safety requirements.
- Launch language and territory-specific compliance.
