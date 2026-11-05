# KALE DÜELLOSU — API AND WEBSOCKET SPECIFICATION

## 1. API conventions

Base path: `/api/v1`  
Format: JSON  
Authentication: Bearer access token  
Time: ISO-8601 UTC  
IDs: UUID  
Errors: consistent machine-readable envelope.

Example error:
```json
{
  "error": {
    "code": "INSUFFICIENT_RESOURCES",
    "message": "Not enough Stone to complete this upgrade.",
    "details": {"required": 120, "available": 85},
    "request_id": "uuid"
  }
}
```

Use pagination for lists. Apply rate limits. Mutating endpoints should accept `Idempotency-Key` where duplicate execution is harmful.

## 2. Authentication endpoints

| Method | Route | Purpose |
|---|---|---|
| POST | `/auth/register` | Register |
| POST | `/auth/login` | Login |
| POST | `/auth/refresh` | Rotate session |
| POST | `/auth/logout` | Revoke session |
| POST | `/auth/password/forgot` | Start recovery |
| POST | `/auth/password/reset` | Complete recovery |
| DELETE | `/account` | Request account deletion |

Register request:
```json
{
  "email": "player@example.com",
  "password": "strong-password",
  "username": "CastlePlayer",
  "locale": "tr-TR"
}
```

Login response:
```json
{
  "access_token": "opaque-or-signed-token",
  "expires_in": 900,
  "refresh_token": "secure-refresh-token",
  "player": {
    "id": "uuid",
    "username": "CastlePlayer",
    "player_level": 1
  }
}
```

Do not log tokens or passwords.

## 3. Player endpoints

- `GET /players/me`
- `PATCH /players/me`
- `GET /players/me/statistics`
- `GET /players/{playerId}/public`

Public profile response must exclude private account fields.

## 4. Castle endpoints

- `GET /castle`
- `GET /castle/structures`
- `POST /castle/structures`
- `POST /castle/structures/{structureId}/upgrade`
- `POST /castle/repair`
- `POST /castle/structures/{structureId}/repair`
- `GET /castle/damage-history`

Build request:
```json
{
  "structure_key": "outer_wall",
  "slot_key": "north_wall"
}
```

Upgrade/repair requests must include an idempotency key. Backend validates prerequisites, state version, costs, and balances.

## 5. Matchmaking endpoints

- `POST /matchmaking/queue`
- `GET /matchmaking/queue`
- `DELETE /matchmaking/queue`
- `POST /matchmaking/matches/{matchId}/confirm`
- `POST /matchmaking/matches/{matchId}/decline`

Queue request:
```json
{
  "mode": "ranked",
  "region": "auto"
}
```

Queue response:
```json
{
  "ticket_id": "uuid",
  "status": "SEARCHING",
  "started_at": "2026-11-05T12:00:00Z"
}
```

The server derives MMR, league, and castle power from trusted records; clients must not submit these as authoritative values.

## 6. Friend endpoints

- `GET /friends`
- `GET /friends/requests`
- `GET /players/search?q=`
- `POST /friends/requests`
- `POST /friends/requests/{requestId}/accept`
- `POST /friends/requests/{requestId}/reject`
- `DELETE /friends/{friendId}`
- `POST /players/{playerId}/block`
- `DELETE /players/{playerId}/block`
- `POST /duel-invitations`
- `POST /duel-invitations/{invitationId}/accept`
- `POST /duel-invitations/{invitationId}/reject`

Invitation expiry, block rules, rate limits, and match eligibility are server-enforced.

## 7. Duel endpoints

- `GET /duels/{duelId}`
- `GET /duels/{duelId}/state`
- `POST /duels/{duelId}/answers`
- `POST /duels/{duelId}/abilities`
- `GET /duels/{duelId}/result`
- `POST /duels/{duelId}/rematch`

Answer request:
```json
{
  "round_number": 2,
  "question_id": "uuid",
  "answer": {"option_id": "B"},
  "client_submission_id": "uuid"
}
```

Server validates authenticated player membership, question/round, deadline, answer format, and duplicate submission. Client-supplied correctness or response time is ignored.

## 8. Inventory and economy endpoints

- `GET /inventory`
- `POST /inventory/equip`
- `POST /inventory/unequip`
- `POST /inventory/use`
- `GET /wallet`
- `GET /wallet/transactions`
- `POST /rewards/{rewardId}/claim`

Every use/debit/claim must be validated and idempotent.

## 9. League and season endpoints

- `GET /leagues/me`
- `GET /leaderboards/global`
- `GET /leaderboards/friends`
- `GET /leaderboards/seasons/{seasonId}`
- `GET /seasons/current`
- `GET /seasons/{seasonId}/rewards`
- `POST /seasons/{seasonId}/rewards/{rewardId}/claim`

Leaderboard pagination must use stable ordering and tie-break rules.

## 10. Store and purchase endpoints

- `GET /store/products`
- `POST /purchases/verify`
- `GET /purchases`
- `POST /purchases/restore`

Purchase verification request should contain platform, signed transaction/receipt data, and product identifier. The backend verifies with the relevant platform and grants entitlement exactly once.

## 11. Missions and achievements

- `GET /missions`
- `POST /missions/{missionId}/claim`
- `GET /achievements`
- `POST /achievements/{achievementId}/claim`

## 12. WebSocket contract

Endpoint: `/api/v1/realtime` (or separately versioned gateway path).

Connection must authenticate the user and authorize subscriptions to their own queue, invitations, and active match only.

Envelope:
```json
{
  "event_id": "uuid",
  "event_type": "duel.question_received",
  "sequence": 12,
  "server_time": "2026-11-05T12:00:00Z",
  "match_id": "uuid",
  "payload": {}
}
```

### Client-to-server events
- `matchmaking.queue_join`
- `matchmaking.queue_cancel`
- `matchmaking.match_confirm`
- `matchmaking.match_decline`
- `duel.answer_submit`
- `duel.ability_activate`
- `realtime.ack`
- `realtime.state_sync`

### Server-to-client events
- `matchmaking.queued`
- `matchmaking.searching`
- `matchmaking.match_found`
- `matchmaking.confirmation_requested`
- `matchmaking.match_ready`
- `matchmaking.canceled`
- `duel.started`
- `duel.question_received`
- `duel.answer_received`
- `duel.answer_resolved`
- `duel.score_updated`
- `duel.ability_resolved`
- `duel.attack_resolved`
- `duel.castle_damaged`
- `duel.player_disconnected`
- `duel.player_reconnected`
- `duel.finished`
- `notification.received`

### Example question event
```json
{
  "event_id": "uuid",
  "event_type": "duel.question_received",
  "sequence": 8,
  "server_time": "2026-11-05T12:00:00Z",
  "match_id": "uuid",
  "payload": {
    "round_number": 1,
    "question_id": "uuid",
    "question_version": 3,
    "type": "multiple_choice",
    "prompt": "Which planet is known as the Red Planet?",
    "options": [
      {"id": "A", "text": "Venus"},
      {"id": "B", "text": "Mars"},
      {"id": "C", "text": "Jupiter"},
      {"id": "D", "text": "Mercury"}
    ],
    "deadline_at": "2026-11-05T12:00:10Z"
  }
}
```

The correct answer is deliberately omitted.

### Example resolved event
```json
{
  "event_id": "uuid",
  "event_type": "duel.answer_resolved",
  "sequence": 10,
  "server_time": "2026-11-05T12:00:10Z",
  "match_id": "uuid",
  "payload": {
    "round_number": 1,
    "correct_option_id": "B",
    "players": [
      {"user_id": "uuid-a", "is_correct": true, "score_delta": 100},
      {"user_id": "uuid-b", "is_correct": false, "score_delta": 0}
    ]
  }
}
```

## 13. Error codes

Define stable codes including:
- `UNAUTHORIZED`
- `FORBIDDEN`
- `VALIDATION_ERROR`
- `RESOURCE_NOT_FOUND`
- `CONFLICT`
- `ALREADY_IN_MATCH`
- `MATCH_NOT_READY`
- `ANSWER_TOO_LATE`
- `ANSWER_ALREADY_SUBMITTED`
- `INSUFFICIENT_RESOURCES`
- `ITEM_NOT_AVAILABLE`
- `PURCHASE_VERIFICATION_FAILED`
- `RATE_LIMITED`
- `SERVER_MAINTENANCE`

## 14. API acceptance criteria

- OpenAPI specification exists and is version-controlled.
- Request/response schemas are validated.
- Authentication and authorization are tested.
- Mutations have idempotency where required.
- Pagination and rate limits are consistent.
- WebSocket events are versioned and documented.
- Client can recover state after sequence gaps or reconnect.
