# Mobile Architecture

## Runtime and project structure

Use React Native + Expo + TypeScript with Expo Router. Keep feature code cohesive and transport-independent where possible.

```text
app/                         # route composition, auth/match guards
src/
  features/{auth,castle,duel,matchmaking,friends,league,store,...}/
    api.ts hooks.ts schemas.ts components/ screens/
  shared/
    api/        # typed OpenAPI client, auth refresh, error mapping
    realtime/   # typed envelopes, sequence/reconnect/snapshot
    state/      # minimal session and transient UI stores
    ui/         # tokens and accessible primitives
    storage/    # secure secrets and non-sensitive preferences
    localization/ telemetry/ navigation/
```

Routes cover the five bottom tabs and contextual screens required in the source: onboarding/auth, home castle, duel selection/search/preparation/question/combat/result, workshop/defense/inventory, league/season, store, friends, profile/history, missions, notifications, settings/support/privacy.

## Navigation and active match safety

Guard authenticated routes and deep links. A match session is restored from server state on app resume; never navigate away or clear a duel because of a local timeout. Show reconnecting/paused/recovery states and a safe route back to the active match. Ask confirmation for intentional navigation when a match is live. Handle push/deep links only after auth and resource authorization.

## State ownership

- TanStack Query owns server-backed profile, castle, wallet, inventory, league, friends, catalog, and history; invalidate/refetch after authoritative mutations.
- A small Zustand store holds ephemeral UI state, current connection/session indicators, and presentation state keyed to server `state_version`.
- Do not maintain authoritative scores, timers, wallet balance, inventory, damage, reward claims, MMR, or correctness on client. Predicted countdown is display-only, anchored to server deadline.
- Persist only non-sensitive preferences/drafts as appropriate. Store refresh credential in platform secure storage; never log it or include it in analytics.

## API layer

Generate or type-check the client against OpenAPI. Centralize bearer attachment, refresh rotation, request IDs, idempotency keys, validation, and error envelope mapping. Prevent simultaneous refresh races; handle revoked/expired sessions explicitly. Avoid blind retries on non-idempotent operations. All economic and gameplay mutations render server-confirmed results; optimistic updates are limited to reversible low-risk UI actions.

## WebSocket layer

One session-aware connection manager handles access token renewal, heartbeat, bounded exponential reconnect with jitter, typed schema parsing, event ID dedupe, monotonic sequence, ack, and snapshot/state sync. Route events to feature handlers; commands include unique command IDs and never trust client identity/time. On sequence gap, request recovery before rendering potentially stale match result. WebSocket and REST answer fallback invoke the same server operation.

## UI system, accessibility and localization

Centralize color/type/spacing/radius/motion tokens and reusable buttons, cards, dialogs, loading/error/empty/offline states. Cover reduced motion, dynamic text, screen readers, contrast, safe areas, touch targets, color-independent status, and small/low-memory devices. Turkish is launch language; all copy uses translation keys, dates/numbers/currencies are locale-formatted, and layout supports RTL. Castle visual state and combat animation derive from server-confirmed events.

## Failure and offline behavior

Show explicit retry/recovery for API errors and distinguish offline from server rejection. Cache read-only content only where stale use is safe; do not allow offline economy/gameplay mutations to appear successful. Preserve pending command IDs and reconcile server state after reconnect. Support splash maintenance/update/network states, matchmaking cancellation, pending purchases, failed/refunded/restored purchases, and account deletion/logout flows.

## Test strategy

Unit-test parsers, reducers, validation, sequence handling and display-only timer drift. Component tests cover loading/empty/error/offline/success, accessibility and destructive confirmations. Contract tests run against OpenAPI and event schemas. E2E tests on representative iOS/Android devices cover registration to duel, ranked/friend flows, background/reconnect, build/repair, results, store lifecycle and account deletion. Visual QA covers required high-fidelity screens and reduced motion.
