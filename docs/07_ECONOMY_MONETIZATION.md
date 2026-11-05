# KALE DÜELLOSU — ECONOMY AND MONETIZATION DESIGN

## 1. Objectives

- Support sustainable long-term progression.
- Make regular play rewarding.
- Ensure castle repairs are achievable after losses.
- Create meaningful resource sinks.
- Monetize cosmetics and optional content.
- Prevent paid advantages from deciding ranked outcomes.
- Make all economic flows measurable and auditable.

## 2. Currency model

### Gold
Common currency for basic construction, upgrades, repairs, and standard items.

### Stone
Construction material for walls, towers, gates, and structural upgrades.

### Crystals
Premium currency for cosmetics, premium offers, and optional convenience. Crystals must not be mandatory for ranked participation or basic repair.

Do not use multiple premium currencies in MVP unless a clear product need exists.

## 3. Resource sources

| Source | Gold | Stone | Crystals | XP/Season XP |
|---|---:|---:|---:|---:|
| Ranked win | Yes | Yes | Rare/limited if configured | Yes |
| Ranked loss with fair completion | Small | Small/optional | No by default | Limited |
| Friendly match | Configurable, anti-farm limited | Configurable | No | Limited |
| Daily mission | Yes | Yes | Rare milestone only | Yes |
| Weekly mission | Yes | Yes | Optional milestone | Yes |
| Achievement | Yes | Yes | Selected milestones | Yes |
| League reward | Yes | Yes | Optional | Season XP |
| Season track | Yes | Yes | Premium track where disclosed | Yes |
| Purchase | Optional packages | Optional packages | Yes | No direct MMR |

Exact reward quantities must be tuned through simulation and playtests.

## 4. Resource sinks

- Construction and upgrades.
- Castle and structure repairs.
- Eligible consumables.
- Cosmetic unlocks.
- Event participation.
- Optional convenience actions.

Core repair must always remain possible through earnable resources. Avoid a situation where a defeated player cannot play effectively until paying.

## 5. Economy transaction rules

- Every grant/debit creates a ledger record.
- No client-side balance mutation.
- Transaction and wallet update occur atomically.
- Use unique idempotency keys.
- Record source/reason/reference.
- Provide support/admin audit trail.
- Reconcile materialized wallet balances against ledger.
- Negative balances are prohibited unless a specifically designed debt mechanic is approved (not recommended).

## 6. Progression pacing

Define separate progression tracks:
- Player XP/level.
- Castle level.
- Competitive MMR/league.
- Season XP.
- Achievements.

Do not gate competitive matchmaking behind spending or excessive castle upgrades.

New players should receive starter materials and a protected onboarding period.

## 7. Item rarity

Suggested rarity labels:
- Common
- Uncommon
- Rare
- Epic
- Legendary

Rarity is a collection/presentation attribute, not an automatic measure of combat superiority. Every gameplay item must have explicit effect values and balance limits.

## 8. Store catalog

Categories:
- Castle skins.
- Castle decorations.
- Structure appearance variants.
- Player avatars and frames.
- Attack animations.
- Shield/defense visual effects.
- Season pass.
- Optional resource bundles.
- Limited-time cosmetic bundles.

Product details must disclose contents, price, duration, and restrictions.

## 9. Season pass

Two tracks:
- Free track available to all players.
- Premium track unlocked through native purchase.

Progress through Season XP and missions.
Rewards may include cosmetics, resources, profile items, and non-ranked-limiting items.
Premium pass entitlement must be verified server-side.
Season end and reward claim deadlines must be clear.

## 10. Pay-to-win safeguards

- Cosmetics have no gameplay stats.
- Purchased consumables have free earn paths.
- Ranked item slots and charges are capped.
- No purchased immunity to all attacks.
- No purchase directly increases MMR.
- No paid-only question advantages.
- No paid access to easier questions.
- Ability power must be normalized or constrained in ranked play.
- Matchmaking can use castle power as a secondary signal.
- Publish clear item descriptions.

## 11. Anti-farming and abuse controls

- Limit repeated friendly-match rewards.
- Detect collusive win trading.
- Apply diminishing returns to repetitive low-quality matches.
- Do not grant full rewards for early abandonment.
- Enforce daily/weekly caps where necessary.
- Track suspicious resource flows.
- Review unusual account clusters and repeated opponent pairs.

## 12. Economy simulation

Before launch, run simulations for:
- New player after 1, 7, 30, and 90 days.
- Average player activity.
- High-activity player.
- Win rates from 30% to 70%.
- Repeated losses and recovery.
- Repair-heavy play.
- Premium and free progression.
- Resource inflation.
- Ability consumption.
- Season pass completion.

Measure:
- Time to first upgrade.
- Time to each castle tier.
- Average repair affordability.
- Resource balance distribution.
- Source/sink ratio.
- Percentage of players unable to afford basic repair.
- Free vs paid progression gap.

## 13. Purchase lifecycle

1. Client requests product catalog.
2. User selects product.
3. Native store purchase sheet handles transaction.
4. Client sends signed transaction/receipt to backend.
5. Backend verifies with platform.
6. Backend records transaction uniquely.
7. Entitlement/currency is granted transactionally.
8. Client refreshes account state.
9. Refund/revocation notifications are reconciled.

Do not grant items based solely on client callback.

## 14. Compliance and consumer clarity

- Use native in-app purchase mechanisms for digital goods where required by platform policy.
- Display price and product contents clearly.
- Support purchase restoration where applicable.
- Handle refunds and revoked entitlements.
- Publish terms, privacy, and support paths.
- Assess local rules for each launch territory.
- Do not add real-money betting or cash-out mechanics to the core game.
