# Real-Time Architecture

## Transport and ownership

Use authenticated WebSocket over TLS at `/api/v1/realtime`; REST is the normal account/content interface and may provide a fallback answer command. WebSocket is a transport, not a second gameplay implementation: both transports call the same typed application command, authorization, idempotency, and transaction logic.

The gateway may run in the API process initially. Multiple instances use Redis pub/sub or an equivalent adapter for delivery, while PostgreSQL and recoverable match checkpoints remain authoritative. Presence is advisory; it never determines a winner.

## Envelope and sequencing

Server event:

```json
{
  "event_id": "uuid",
  "event_type": "duel.question_received",
  "schema_version": 1,
  "sequence": 8,
  "server_time": "2026-11-05T12:00:00Z",
  "match_id": "uuid",
  "state_version": 14,
  "payload": {}
}
```

Client command includes `command_id`, `event_type`, `schema_version`, and payload. The authenticated principal is taken from the connection, never from a client-supplied user ID. Sequence is monotonic per match stream; event ID is globally unique. Clients deduplicate event IDs, detect gaps, acknowledge last applied sequence, and request a snapshot. Server checks match membership for every subscription and command.

## Lifecycle

1. Establish WSS, validate origin/client version, authenticate access token, and bind user/session.
2. Subscribe only to that user's queue, invitation, notification, and currently authorized match channels.
3. Maintain heartbeat and last-seen connection state; allow multiple device connections only under explicit session policy.
4. Join queue or receive match-found; confirmation deadline is server timestamp.
5. Start duel after both confirmation requirements pass; persist immutable rules, question-version, and loadout snapshots.
6. Server sends question and authoritative deadline. Do not send answer key.
7. Accept one answer per user/question until server receive time reaches deadline. Answer contents are checked against the selected version on server.
8. Resolve after both submissions or deadline. Persist before broadcasting result.
9. After final round, finalize duel/combat/rewards once; persist result and state version before `duel.finished`.
10. On disconnect, keep server timer running unless a documented pause policy applies. On reconnect, authenticate and send snapshot plus events after last acknowledged sequence.

## Event contracts

### Client to server

| Event | Payload | Validation / result |
|---|---|---|
| `matchmaking.queue_join` | `{mode:"ranked", region:"auto"}` | Server derives rating, castle power and eligibility; returns queued ticket once |
| `matchmaking.queue_cancel` | `{ticket_id, command_id}` | Owner/status check; race-safe against reservation |
| `matchmaking.match_confirm` | `{match_id, command_id}` | Participant, deadline and pending status; idempotent |
| `matchmaking.match_decline` | `{match_id, reason?, command_id}` | Participant; releases reservation per policy |
| `duel.answer_submit` | `{round_number, question_id, answer, command_id}` | Active member, expected round/question, deadline, schema and one answer only |
| `duel.ability_activate` | `{ability_key, target?, command_id}` | Match snapshot, owner, phase, eligibility, charges, target and cooldown |
| `realtime.ack` | `{match_id, sequence}` | Monotonic ack, authorized stream |
| `realtime.state_sync` | `{match_id, last_sequence, state_version}` | Membership; returns snapshot or bounded event replay |

### Server to client

| Event | Essential payload / semantics |
|---|---|
| `matchmaking.queued` / `searching` | Ticket ID, search start, current broad stage; no fabricated ETA |
| `matchmaking.match_found` / `confirmation_requested` | Match ID, opponent public summary, server deadline, confirmation state |
| `matchmaking.match_ready` | Duel ID, mode, state version, permitted loadout snapshot |
| `matchmaking.canceled` | Ticket/match and server-resolved cancellation reason |
| `duel.started` | Match rules snapshot summary, first round and state version |
| `duel.question_received` | Round, question ID/version, localized prompt/options, deadline; no answer key |
| `duel.answer_received` | Which participants submitted; never reveal answer content |
| `duel.answer_resolved` | Correct option and per-player correctness/score after resolution |
| `duel.score_updated` | Authoritative totals and breakdown |
| `duel.ability_resolved` | Applied effect, target, limits and result |
| `duel.attack_resolved` / `duel.castle_damaged` | Combat event ID, damage breakdown, affected castle/structure state version |
| `duel.player_disconnected` / `reconnected` | Public connection status; no sensitive reason |
| `duel.finished` | Immutable result, rating/reward/damage summaries, final state version |
| `notification.received` | Notification ID/type/deep link under privacy settings |
| `realtime.error` | Stable code, command ID, retryability and request ID |

### Example question and answer command

```json
{
  "event_type": "duel.question_received",
  "schema_version": 1,
  "sequence": 8,
  "server_time": "2026-11-05T12:00:00Z",
  "match_id": "uuid",
  "payload": {
    "round_number": 1,
    "question_id": "uuid",
    "question_version": 3,
    "type": "multiple_choice",
    "prompt": "Localized question text",
    "options": [{"id":"A","text":"..."},{"id":"B","text":"..."}],
    "deadline_at": "2026-11-05T12:00:10Z"
  }
}
```

```json
{
  "event_type": "duel.answer_submit",
  "schema_version": 1,
  "command_id": "uuid",
  "match_id": "uuid",
  "payload": {
    "round_number": 1,
    "question_id": "uuid",
    "answer": {"option_id":"B"}
  }
}
```

Client response duration is not trusted. Server records receipt timestamp and computes duration from the server-controlled round start, applying the documented latency policy. The latency policy must be defined and tested before ranking is enabled.

## Ordering, retries and recovery

- Acknowledge only after the client applies an event; reconnect with last applied sequence.
- Duplicate command ID returns the original command outcome. A new command with an already submitted answer returns `ANSWER_ALREADY_SUBMITTED`.
- If requested sequence history has expired or a gap cannot be filled, return a complete authorized snapshot with current state version.
- Persist important transitions before broadcast; use outbox for cross-instance delivery.
- Reconnect uses exponential backoff with jitter and a bounded retry policy; active match is never silently discarded.
- Server restart recovers unresolved matches from durable deadlines/checkpoints. Timers are scheduled/recomputed from persisted UTC deadlines, not process-local clocks.
- Late answers are rejected consistently even during recovery. Document behavior for server outage spanning a deadline; recommended policy is pause/cancel and refund participation, not infer answers.

## Timeouts and disconnections

Question timers continue on disconnect unless an approved pause rule says otherwise. Missing submission counts as timeout. Matchmaking confirmation expiry releases both reservations and applies an explicit queue re-entry policy. Abandonment rewards/MMR are configurable, but server records disconnects, grace period, and final reason. No client time value changes the deadline.

## Operations and tests

Monitor connection count, auth failures, message rates, payload rejects, sequence gaps, reconnect success, queue/event lag, timer lateness, active matches per instance, and broadcast failures. Test duplicate/reordered commands, simultaneous answers, timer boundary, disconnect/reconnect, stale token, unauthorized subscription, multi-instance delivery, process restart, and event replay/snapshot recovery.
