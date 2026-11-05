# KALE DÜELLOSU — GAME DESIGN DOCUMENT

**Version:** 1.0 — Proposed baseline  
**Purpose:** Define consistent rules for duels, castle combat, progression, abilities, and economy.

## 1. Design pillars

1. Knowledge matters more than reaction speed.
2. Every match creates a clear outcome and understandable consequences.
3. Castle damage is meaningful but recoverable.
4. Defense and timing create strategic choices.
5. Long-term progression is persistent.
6. Monetization must not invalidate skill.

## 2. Standard duel format

Proposed ranked defaults:
- 5 rounds.
- 10-second answer window per round.
- Same question ID/version for both players in each round.
- Four options for standard multiple-choice questions.
- One submission per player per round.
- Server timestamp is authoritative.
- Optional short inter-round transition.
- Total target match duration: approximately 2–3 minutes, excluding queue time.

All values are remotely configurable and should be validated through playtests.

## 3. Question selection

Each question has:
- ID and immutable version.
- Language.
- Category.
- Difficulty (Easy, Medium, Hard, Expert).
- Type.
- Prompt and media reference.
- Options and correct-answer key.
- Explanation.
- Active/published status.
- Exposure/repetition metadata.

A match receives a preselected sequence to avoid mid-match content changes. Do not send the correct answer to clients before the round is resolved. Both players receive the same content version.

Question selection should balance category and difficulty across a match. Avoid immediate repeats for a player and detect overexposed questions. Image questions must have accessible text alternatives where appropriate.

## 4. Scoring rules

Proposed initial formula per player:
- Correct answer: 100 points.
- Incorrect or timeout: 0 points.
- Three consecutive correct answers: +50 streak bonus.
- Speed bonus: +25 only if both players answer correctly and the response-time difference exceeds a configurable meaningful threshold.

Important:
- A fast incorrect answer earns no speed bonus.
- Correctness is the primary factor.
- Speed bonus must not make a lower-correctness player win.
- Use deterministic server-side calculations.
- Store raw response duration and scoring breakdown for audit.

### Winner determination
Recommended lexicographic comparison:
1. Total score.
2. Number of correct answers.
3. Aggregate response time for correct answers.
4. If still tied, a sudden-death question or declared draw.

For MVP, choose one tie policy and expose it in game configuration. Recommended default: one sudden-death question with the same synchronization rules. If the sudden-death question also ties, declare a draw. Ranked MMR policy for draws should be configured separately.

## 5. Attack opportunity model

Recommended model: a player earns attack energy for correct answers; the duel winner receives a modest end-of-match advantage. This gives both players meaningful participation while preserving a clear winner.

Proposed baseline:
- Each correct answer grants 1 attack point.
- Incorrect/timeout grants 0.
- Winner receives +1 victory attack point.
- Maximum base attack points per player in a five-question match: 6.
- Combat power is converted to damage using configurable values and equipped abilities.

The product team must decide whether both players can damage each other in the same match or only the winner can attack. Recommended default: both players may inflict limited damage based on earned attack points, while the winner receives a bonus. This reduces the frustration of a total zero-reward loss. If the product chooses winner-only damage, loser participation rewards must still be carefully balanced and explained.

## 6. Castle model

Each player has one primary persistent castle in MVP.

Castle attributes:
- Level.
- Overall durability and maximum durability.
- Structure slots.
- Structures and their levels.
- Defense rating.
- Cosmetic theme.
- Current protection state.
- Last attacked timestamp.

Structures:
- Main Keep.
- Outer Wall.
- Gate.
- Defensive Tower.
- Archer Tower.
- Magic Tower.
- Watchtower.
- Resource Building.
- Special Defense Structure.

Each structure definition specifies:
- Unlock level.
- Construction cost.
- Upgrade cost curve.
- Durability.
- Defensive effect.
- Slot requirements.
- Visual asset key.
- Max level.

A structure's destruction in combat should disable or reduce its effect, not delete its level or permanently erase investment. Repair restores functionality.

## 7. Damage and combat resolution

All values below are proposed tuning parameters.

Suggested formula:

`rawDamage = baseDamage + (attackPoints * damagePerAttackPoint) + abilityBonus`

`afterDefense = max(minDamage, rawDamage * (1 - defenseReduction))`

`finalDamage = min(afterDefense, targetRemainingDurability)`

Where:
- `baseDamage` is a configurable match outcome value.
- `damagePerAttackPoint` converts earned attack points to damage.
- `abilityBonus` comes from validated offensive abilities.
- `defenseReduction` is calculated from shields, walls, traps, and caps.
- `minDamage` prevents absolute immunity only if the design requires it.

Use integer damage units to avoid floating-point discrepancies. Apply damage in a single transaction with a unique combat event ID. A repeated event must not apply damage twice.

Recommended protection:
- No single match can reduce the core castle below a configured critical durability floor.
- New-player protection applies for a configurable number of matches or account progression window.
- Repeated attacks against the same defender may trigger diminishing damage or temporary protection.
- Protection rules must be visible and understandable.

## 8. Repair rules

Repair options:
- Repair entire castle.
- Repair a selected structure.
- Use a repair kit.
- Complete a repair mission (optional).

Proposed cost formula:
`repairCost = ceil(durabilityMissing * costPerDurability * structureLevelMultiplier)`

The formula and multipliers must be configurable.

Rules:
- Repair cannot exceed max durability.
- Balance deduction and durability update occur in one database transaction.
- A repair request uses an idempotency key.
- Insufficient resources return a clear error.
- Repair does not remove structure upgrades.
- No mandatory real-time waiting should block basic recovery in MVP.

## 9. Ability system

### Offensive
- Destruction Bomb: adds fixed or percentage bonus damage, subject to caps.
- Critical Strike: boosts the next valid attack.
- Siege Attack: targets an eligible structure.
- Piercing Attack: reduces a defined defense effect.

### Defensive
- Shield: absorbs or reduces damage for a defined attack window.
- Reinforced Walls: passive capped damage reduction.
- Trap: reduces or counters a qualifying attack.
- Repair Kit: restores durability.

Each ability definition includes:
- ID, category, rarity, description.
- Effect type and magnitude.
- Duration/cooldown.
- Charges/usage limits.
- Target eligibility.
- Unlock/acquisition path.
- Ranked eligibility.
- Stackability and counter rules.

### Loadout
Proposed initial limit:
- 2 offensive slots.
- 2 defensive slots.
- Limited consumable charges per match.
- No duplicate stacking unless explicitly configured.

Players select loadouts before match start. The server snapshots the loadout at match start so players cannot alter it mid-match.

### Ranked fairness
- Every ability must have a counter or a clear cap.
- Premium purchase must not bypass ranked limits.
- Cosmetic variants have no gameplay effect.
- Matchmaking may consider castle/ability power as a secondary fairness signal.
- Competitive balancing can disable or normalize specific items.

## 10. Progression systems

Keep these separate:
- Account/player XP and level.
- Castle level.
- Competitive MMR and league.
- Season XP and season track.
- Achievement progress.

A player may have a high-level castle and a lower league, or vice versa.

### Account XP
Earned from valid participation, wins, missions, and achievements. Use diminishing returns or daily caps to reduce farming.

### Castle XP/level
Earned through construction, upgrades, and eligible gameplay milestones. Castle level unlocks slots and structures.

### MMR
Represents competitive skill. Updated only from eligible ranked matches.

### Season XP
Earned from ranked play and season missions. Resets or expires at season end according to published season rules.

## 11. Leagues

Initial tier names:
1. Stone
2. Bronze
3. Silver
4. Gold
5. Diamond
6. Champion

Each league defines:
- Minimum rating.
- Optional divisions.
- Promotion/demotion rules.
- Season reward.
- Visual badge.
- Leaderboard scope.

Do not hardcode thresholds in client code. Store configuration server-side.

## 12. Seasons

Proposed duration: 30 days, configurable.

Season includes:
- Start/end timestamps.
- Ranked leaderboard.
- Free reward track.
- Premium reward track.
- Missions.
- Cosmetic theme.
- Final rewards.

At season end:
- Freeze eligible results.
- Reconcile pending matches.
- Calculate final ranks.
- Issue rewards idempotently.
- Archive leaderboard.
- Apply rating reset/soft reset.
- Preserve castle, inventory, purchased cosmetics, achievements, and permanent account progression.

## 13. Economy

Currencies:
- Gold: common construction, repair, standard items.
- Stone: walls and structures.
- Crystals: premium currency, cosmetics, optional convenience products.

Sources:
- Duel participation and wins.
- Missions.
- Achievements.
- League and season rewards.
- Events.
- Approved purchases.

Sinks:
- Construction.
- Upgrades.
- Repairs.
- Eligible consumables.
- Cosmetics and event participation.

Every source/sink must be tracked in an immutable ledger. Economy analytics must detect inflation, hoarding, exploit loops, and progression dead ends.

## 14. Reward policy

Winning should grant more than losing, but losing players should retain a small, fair participation reward when they complete the match honestly.

Rewards may include:
- Gold.
- Stone.
- XP.
- Ability fragments/items.
- Season XP.
- Mission progress.

Reward rules must:
- Exclude or reduce rewards for abandonment.
- Prevent repeated claims.
- Apply daily anti-farming limits if needed.
- Be configurable by match type.
- Avoid granting ranked rewards for private matches unless explicitly intended.

## 15. Missions and achievements

Daily examples:
- Complete 3 duels.
- Win 1 duel.
- Answer 10 questions correctly.
- Repair a structure.
- Use a defense item.

Weekly examples:
- Win 10 ranked matches.
- Achieve a 5-answer streak.
- Upgrade a structure.
- Successfully defend against attacks.

Achievements:
- First victory.
- 10/100 victories.
- First structure upgrade.
- First defense.
- Reach each league.
- Correct-answer milestones.

Mission progress must derive from validated server events. Claims must be idempotent.

## 16. Monetization and fairness

Primary monetization:
- Castle skins and decorations.
- Avatar/profile cosmetics.
- Attack and defense visual effects.
- Premium season pass.
- Optional resource packs.
- Limited-time cosmetic bundles.

Consumable purchases may exist only with:
- Free earn paths.
- Match limits.
- Ranked caps/normalization.
- Server validation.
- Clear item descriptions.

No paid item may guarantee a win, bypass MMR, or make a player immune to normal competitive counterplay.

## 17. Balancing and simulation requirements

Before launch, simulate:
- 100-match progression for new, average, and highly active players.
- Resource income vs construction/repair costs.
- Win/loss streaks.
- Repeated attacks on one player.
- Low-activity player recovery.
- Premium vs free progression.
- Ability combinations and counterplay.
- MMR inflation and league distribution.

Do not finalize economy values from intuition alone. Use telemetry and playtests.
