# KALE DÜELLOSU — UI/UX SPECIFICATION

## 1. Experience direction

Create a polished, mobile-first fantasy castle game with clear competitive feedback. Use a stylized 3D castle as the main visual anchor, supported by readable 2D interface elements.

The product must not resemble a generic quiz form. The castle, player identity, combat, and progression should be visible throughout the experience.

## 2. Visual design principles

- Strong visual hierarchy and readable question text.
- Clear distinction between ranked and friendly modes.
- High-contrast timer and answer feedback.
- Consistent resource and currency presentation.
- Tactile buttons with clear pressed/disabled/loading states.
- Castle visuals communicate damage and upgrades.
- Animations are brief and never block critical input.
- Support safe areas, small screens, large text, and reduced motion.

## 3. Design tokens

Create a centralized token system rather than scattering raw values.

Suggested direction:
- Background: deep navy/charcoal fantasy environment.
- Primary action: warm gold.
- Secondary action: muted blue.
- Success: green.
- Error/damage: red.
- Premium currency: crystal/violet accent.
- Surfaces: layered dark and light panels depending on scene.
- Typography: highly legible display face for headings and neutral sans-serif for body.
- Spacing: 4/8-point scale.
- Radius: small, medium, large, and pill tokens.
- Motion: short feedback, medium transitions, long celebratory animations.

Exact colors, fonts, and component dimensions must be validated in visual prototypes.

## 4. Navigation model

Bottom tabs:
- Castle
- Duel
- League
- Friends
- Profile

Contextual routes:
- Workshop
- Defense
- Inventory
- Store
- Season Pass
- Missions
- Notifications
- Settings
- Match History

Prevent accidental navigation away from an active match. If the app backgrounds, preserve match state and reconnect.

## 5. Screen specifications

### 5.1 Splash
Purpose: initialize app.
Components: logo, loading indicator, version/update state.
States: loading, mandatory update, maintenance, network failure.
Acceptance: user never sees an indefinite spinner without recovery guidance.

### 5.2 Onboarding
Use short illustrated steps:
1. Answer questions.
2. Win and earn attack opportunities.
3. Build and defend your castle.
4. Climb leagues.
Allow skip except mandatory legal consent.

### 5.3 Authentication
Fields: email, password, username where applicable.
Include inline validation, password visibility, loading, errors, recovery, and terms links.
Do not expose whether an arbitrary email exists in ways that enable account enumeration.

### 5.4 Home/Castle
Primary layout:
- Header: avatar, username/level, notifications.
- Resource bar: Gold, Stone, Crystals.
- Central castle scene with interactive structures.
- Castle durability indicator.
- Current league/MMR card.
- Primary CTA: Find Opponent.
- Secondary CTA: Challenge Friend.
- Missions/season progress cards.

Castle interactions:
- Tap structure to inspect level, durability, upgrade, and repair.
- Visual damaged state must match server data.
- Provide a clear path to Workshop and Defense.

### 5.5 Duel selection
Two equally understandable cards:
- Automatic Ranked Match.
- Friend Challenge.
Explain ranked consequences before queueing.
Show eligibility restrictions and active queue state.

### 5.6 Matchmaking search
Show search animation, elapsed time, current rating/league, and cancel.
Do not show fabricated exact ETA.
When matched, show opponent and a confirmation countdown.
If canceled or failed, explain state and allow retry.

### 5.7 Duel preparation
Show both players, league, castle preview, and permitted loadouts.
Display rules and match type.
Ready/decline actions must have clear timeout feedback.

### 5.8 Question duel
Priority order:
1. Question prompt.
2. Timer.
3. Answer options.
4. Round progress.
5. Both player status indicators.
6. Score/streak.

States:
- Waiting for question.
- Answerable.
- Answer submitted.
- Opponent submitted.
- Correct.
- Incorrect.
- Timeout.
- Reconnecting.
- Match paused by server.
- Match ended.

Answer buttons must lock after submission. Feedback must not leak the correct answer before both submissions are resolved if that would disadvantage a player.

### 5.9 Combat scene
Show attack source, target, ability effect, impact, and resulting durability.
Provide skip/reduce-motion option.
Combat display must be derived from server-resolved event data.

### 5.10 Result screen
Show winner, score breakdown, answer accuracy, damage, defenses, rewards, MMR change, league movement, and next actions.
Clearly distinguish earned rewards from promotional offers.
Provide rematch only where eligible.

### 5.11 Workshop
Tabs/categories for structures and upgrades.
Each item shows:
- Preview.
- Level/current state.
- Unlock requirements.
- Cost.
- Resulting stats.
- Build/upgrade action.
Insufficient resources must show exact missing amount.
Confirm expensive actions.

### 5.12 Defense
Show damaged structures, repair costs, active shields, available defenses, and attack history.
Repair confirmation must state resource cost and expected durability restored.

### 5.13 Inventory
Categories: offensive, defensive, repair, materials, cosmetics.
Each item shows quantity, rarity, effect, expiry, and usage eligibility.
Disable invalid use with explanation.

### 5.14 League
Show current tier, rating, progress to promotion, leaderboard, recent rating changes, and season time.
Never imply rank is guaranteed until server confirms.

### 5.15 Season
Show season theme, end date/time localized, free/premium tracks, season XP, missions, and claimable rewards.
Clearly label paid track and its benefits.

### 5.16 Store
Product cards must include price, contents, purchase conditions, and confirmation.
Use native purchase sheet where required.
Handle pending, success, canceled, failed, refunded, and restored states.

### 5.17 Friends
Search, incoming/outgoing requests, friend list, online state, profile, challenge, remove, block, report.
Do not reveal private presence if user privacy settings prohibit it.

### 5.18 Profile
Show player identity, stats, achievements, castle summary, and match history.
Edit only allowed fields.
Public profile must not show email, payment, or private account identifiers.

### 5.19 Missions
Group daily, weekly, and season objectives.
Show progress and claim state.
Claim button must show loading and server-confirmed result.

### 5.20 Notifications
Unread/read states, category icon, timestamp, deep link, and empty state.
Support mark-as-read and notification preferences.

### 5.21 Settings and support
Language, sound/music, haptics, reduced motion, notification preferences, privacy, blocked users, support, terms, privacy policy, account deletion, logout.

## 6. Interaction patterns

- Use confirmation dialogs for irreversible or costly actions.
- Use optimistic UI only for reversible low-risk actions.
- For economy, inventory, combat, and purchases, show server-confirmed state.
- Disable repeated submissions while requests are pending.
- Preserve draft/search state when appropriate.
- Provide clear retry paths for network errors.
- Never silently discard a player's active match.

## 7. Accessibility and localization

- Dynamic text sizing.
- Minimum touch target guidance.
- Screen-reader labels for controls and resource indicators.
- Color must not be the only status signal.
- Reduced motion.
- Turkish-first copy with translation keys.
- Locale-aware dates, numbers, and currencies.
- Layout architecture must support RTL languages.

## 8. Required UX deliverables

Before implementation, create:
- Information architecture diagram.
- User flow diagrams.
- Low-fidelity wireframes for all 21 screens.
- High-fidelity designs for Home, Matchmaking, Duel, Combat, Results, Workshop, League, and Store.
- Component library.
- Interaction/state specifications.
- Accessibility review.
- Design QA checklist.
