# Recommended Technology Stack

## Selection

| Area | Recommendation | Reason / constraint |
|---|---|---|
| Mobile | React Native, Expo, TypeScript | Directly matches product requirement; shared iOS/Android delivery and typed feature modules |
| Navigation | Expo Router | File-based routing aligned with Expo; guard active-match routes and deep links |
| Server state | TanStack Query | Cache, invalidation, retry, and mutation lifecycle for API-backed data |
| Transient client state | Zustand, limited to UI/session/live view state | Small surface, avoids duplicating server-authoritative state |
| Forms | React Hook Form + Zod | Typed validation and clear field-level UX; validate again on server |
| API contracts | OpenAPI 3.x, generated/checked TypeScript client | Avoid drift; contract validation in CI |
| Real-time | Native WebSocket client behind typed event adapter | Sequence, replay/snapshot recovery, reconnect and deduplication are app-specific requirements |
| Secure storage | Expo SecureStore for refresh credentials | Platform-backed storage; never bundle secrets |
| Local persistence | AsyncStorage only for non-sensitive preferences/cache | No local authoritative economy or match result |
| Animation/audio/haptics | React Native Reanimated; Expo AV / Haptics as needed | Native-feeling short feedback; honor reduced-motion and device capability |
| Push | Expo Notifications initially, provider abstraction | Reduce MVP integration cost; validate delivery and store/platform constraints |
| Backend | Go modular monolith | Existing architecture document recommendation; explicit modules, predictable deployment, suitable concurrent I/O |
| HTTP | Go standard library or established minimal router | Avoid framework lock-in; typed validation/middleware; select one during foundation |
| WebSocket | Maintained Go WebSocket library behind transport interface | Must support bounded messages, origin/auth checks, graceful shutdown |
| Database | PostgreSQL | Relational consistency, constraints, transactions, JSONB for versioned configuration only |
| Migrations | Versioned SQL migration tool selected in foundation | Reproducible, reviewed schema changes; safe forward deployment |
| Ephemeral coordination | Redis | Presence, queue acceleration, short-lived coordination; never durable ownership of value |
| Background processing | Go worker using PostgreSQL outbox initially | Reliable retries without a broker in MVP |
| Admin | Web app, TypeScript, React (specific framework selected after admin UX prototype) | Admin is required; keep separate privileged client and deployment boundary |
| Containers/CI | Docker, GitHub Actions or existing CI provider | Repeatable local and staged environments |
| Telemetry | OpenTelemetry-compatible traces/metrics/log correlation; crash reporting SDK | Product requires API, real-time, economy and client health visibility |

## Alternatives considered

- **Backend framework:** a Go standard-library/minimal-router service is preferred over a heavier framework for a small modular monolith. Choose a maintained router only if it materially improves routing, middleware, or testability.
- **Node.js/TypeScript backend:** could share language/contracts with mobile, but the project architecture source explicitly recommends Go. Do not switch unless team expertise or operational constraints justify and record an ADR.
- **Microservices:** rejected for initial release; no independent scale/team boundary is established.
- **Redis as source of truth:** rejected for inventory, rating, wallet, castle, match results, or purchases due to durability and consistency requirements.
- **GraphQL:** not required; REST and typed WebSocket events align to the documented API.
- **Local-first economy/game state:** rejected because authoritative validation and idempotency are mandatory.

## Compatibility and operational constraints

- Pin supported Expo SDK, React Native, Node.js, Go, PostgreSQL and Redis versions in repository/toolchain files. Upgrade through explicit compatibility checks rather than floating dependencies.
- Keep OpenAPI and WebSocket JSON schema/event versioning compatible across mobile rollout windows; support old clients only for a documented minimum version.
- Use UTC `timestamptz` on server and locale-aware formatting on clients.
- Use integer units for currency, score and damage; avoid floating-point gameplay calculations.
- Apply platform-native in-app purchase policies for digital goods, server verification, restoration, refund and revocation.
- Separate dev/staging/prod credentials, databases, store products, push topics, and analytics destinations.
- Keep localization keys and question versions independent from client release cadence.

## Development environment

Required: Git repository, Go toolchain, Node.js package manager compatible with the pinned Expo version, Docker with PostgreSQL and Redis, migration runner, test tooling, and platform simulators/devices. Local environment variables are documented in an example file; real secrets remain outside source control. Provide one documented command path for dependencies, migrations, backend/mobile startup, and tests.

## Dependency policy

Add libraries only for a documented requirement, choose maintained packages with compatible licensing and security posture, pin through lockfiles, and record upgrades. Keep native Expo module compatibility verified in CI. Add analytics, crash reporting, push, and purchase SDKs only with the corresponding privacy, platform, and operational review.
