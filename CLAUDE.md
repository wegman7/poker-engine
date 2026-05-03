# Poker Engine — CLAUDE.md

## What This Is

A Go microservice implementing Texas Hold'em poker game logic. It runs one engine per room, manages game state as a state machine, and communicates with clients over WebSockets. There is a separate Django backend that handles auth and room management; this service handles only game logic.

## Running the Project

```bash
# Development
go run ./cmd/app -env=dev

# Tests
go test ./internal/engine -v

# Docker (prod)
./exec.sh
```

Server listens on `:8080`.

## HTTP API

```
POST /start-engine
Body: { "roomName": "room1", "smallBlind": 1.0, "bigBlind": 2.0 }
```

## WebSocket

Connect: `ws://host/ws/engineconsumer/{roomName}?token={authToken}`

**Incoming commands** (JSON, sent by clients):
`join`, `leave`, `fold`, `check`, `call`, `bet`, `sitIn`, `sitOut`, `addChips`, `startGame`

**Outgoing**: serialized game state (JSON) sent on every engine tick.

## Architecture

### State Machine

The engine runs a tick loop (`tick()`) cycling through these states:

```
StateProcessSitCommands
StateStartHand
StatePauseAfterStartHand
StatePostBlinds
StateDealCards
StateProcessGameCommands   <- main betting loop
StateEndStreet
StateDealStreet
StateShowdown
StateEndHand
```

Two separate command queues:
- **Sit queue** — processed between hands: join, leave, sitIn, sitOut, addChips
- **Game queue** — processed during hands: fold, check, call, bet

### Player Ordering

Players are stored in **two circular linked lists**:
- `next` — all seated players (for seating/blind rotation)
- `nextInHand` — players still active in the current hand

The spotlight pointer tracks whose turn it is. Rotating spotlight handles street transitions and skipping all-in players.

### Pot & Side Pot Logic

Side pots are calculated in `state.go` when players are all-in. Each side pot tracks a `maxWin` per player. See the scratch work at the bottom of `todo.txt` for example calculations.

### Hand Evaluation

Uses `github.com/chehsunliu/poker` for hand ranking. Showdown in `engine.go` compares all active players and splits pots on ties.

## Key Files

| File | Purpose |
|---|---|
| `cmd/app/main.go` | Entry point, HTTP server |
| `internal/engine/engine.go` | Main game loop, state machine, command processing |
| `internal/engine/state.go` | Game state, player management, pot/side pot calculations |
| `internal/engine/player.go` | Player struct, per-action handlers |
| `internal/engine/websockets.go` | WebSocket connection handling, command deserialization |
| `internal/engine/serializeState.go` | Converts state to JSON for clients |
| `internal/engine/startEngineHandler.go` | HTTP handler for `/start-engine` |
| `internal/engine/util.go` | Auth0 token fetch, card comparison helpers |
| `config/config.go` | Environment config (timing, limits, URLs) |

## Configuration

Loaded from `.env` (dev) or `.env.prod` (prod):

| Var | Purpose |
|---|---|
| `BACKEND_URL` | Django backend URL |
| `EMAIL` / `PASSWORD` | Auth0 credentials used by the engine itself |
| `AUTH0_DOMAIN` | Auth0 tenant domain |
| `AUTH0_CLIENT_ID` | OAuth client ID |
| `AUTH0_AUDIENCE` | OAuth audience |

Config struct (from `config/config.go`):
- `DEBUG` — bool
- `ENGINE_LOOP_PAUSE` — 10ms in both dev and prod
- `PAUSE_SHORT` / `PAUSE_MEDIUM` / `PAUSE_LONG` — UI animation delays
- `MAX_PLAYERS` — 9

## Known Bugs / TODO

- **Bug**: When one player is all-in and another is not, the engine sometimes waits for the non-all-in player to act instead of going straight to showdown. Should auto-complete the street.

## Backlog

- Game messages in state updates (e.g., "Player X won $Y with Z hand")
- Time bank implementation (player decision timers)
- Decide where to track player time (frontend vs backend)

## Dependencies

- `github.com/chehsunliu/poker` — hand evaluation
- `github.com/gorilla/websocket` — WebSocket server
- `github.com/joho/godotenv` — env file loading
- Go 1.22.5
- Deployed via Docker to GCP
