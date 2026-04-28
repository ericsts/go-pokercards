# 🃏 Planning Poker

> Versão em português: [README.md](README.md)

A real-time **Planning Poker** (Scrum Poker) application for agile teams to estimate stories together. Built as a hands-on learning project: **Go** backend, **Vue 3** frontend, **WebSocket** communication, **Docker Compose** containers, and automated deploys with **GitHub Actions**.

---

## Table of contents

- [How it works](#how-it-works)
- [Tech stack](#tech-stack)
- [Architecture](#architecture)
- [Project structure](#project-structure)
- [Running locally](#running-locally)
- [Running with Docker](#running-with-docker)
- [Deploying to a VPS](#deploying-to-a-vps)
- [CI/CD with GitHub Actions](#cicd-with-github-actions)
- [WebSocket protocol](#websocket-protocol)
- [Design decisions](#design-decisions)

---

## How it works

1. The **room creator** types their name on the home page and clicks **"Create room"**.
2. They receive a shareable link to send to their team.
3. Each participant opens the link, enters their name, and joins the room.
4. The creator can assign any player as the **Scrum Master**.
5. Everyone secretly selects a card (votes are hidden until revealed).
6. The Scrum Master clicks **"Reveal cards"** — all votes appear simultaneously.
7. The most-voted card and the numeric average are displayed in the centre of the table.
8. The Scrum Master clicks **"New round"** to reset and start again.

---

## Tech stack

| Layer | Technology | Why |
|---|---|---|
| Backend | Go 1.22+ | Native concurrency, great fit for WebSocket servers |
| WebSocket | gorilla/websocket | The de-facto standard Go WebSocket library |
| Unique IDs | google/uuid | v4 UUIDs for rooms and players |
| Frontend | Vue 3 + Vite | Fine-grained reactivity with Composition API, fast builds |
| State | Pinia | Official Vue 3 store, simple and ergonomic |
| Routing | Vue Router 4 | HTML5 history mode, `/room/:id` routes |
| Containers | Docker + Compose | Reproducible builds, easy local dev |
| Proxy | nginx | Serves the frontend and forwards `/api/` to the backend |
| CI/CD | GitHub Actions | Automatic deploy on every push to `main` |

---

## Architecture

```
Browser  ──HTTP──►  nginx (port 3001 on the VPS)
                      │
                      ├── /          ──► Frontend (Vue 3 SPA, static files)
                      └── /api/      ──► Go backend (port 8080, internal to Docker)
                            │
                            ├── POST /api/rooms         create room
                            ├── GET  /api/rooms/:id     get room info
                            ├── GET  /api/cards         list valid card values
                            └── GET  /api/rooms/:id/ws  WebSocket upgrade
```

### Backend packages

```
internal/
  room/       Pure business logic. No networking dependency.
              Room, Player, View + Vote, Reveal, Reset, SetMaster methods.

  session/    WebSocket infrastructure.
              Client  — wraps *websocket.Conn with a dedicated writer goroutine.
              Session — groups a Room + Client map + embedded mutex.
              Store   — thread-safe map of roomID → *Session.

  handler/    HTTP and WebSocket handlers. Uses session.Store.
              Routes registered on Go 1.22's http.ServeMux with pattern params.
```

**Why separate `room` from `session`?**
The `room` package has zero WebSocket imports. This makes all game logic testable with plain unit tests and no network setup. The `session` package is the only one that knows about gorilla/websocket.

### Frontend components

```
views/
  HomeView    Room creation form.
  RoomView    Full room: table, deck, controls.

components/
  RoundTable  Pure-CSS round table. Positions players with trigonometry.
  PlayerSeat  Each player's seat: avatar, card, name, Scrum Master crown.
  CardDeck    Selectable card deck. Shows a badge for the chosen card.
  PokerCard   Single card: front (value), back (voting), empty, winner.

stores/
  room.js     Pinia store: opens/closes WebSocket, processes messages, exposes reactive state.
```

---

## Project structure

```
go-pokercards/
├── .github/
│   └── workflows/
│       └── deploy.yml          CI/CD pipeline
├── backend/
│   ├── cmd/server/main.go      Entry point
│   ├── internal/
│   │   ├── room/               Business logic + tests
│   │   ├── session/            WebSocket layer (Client, Session, Store) + tests
│   │   └── handler/            HTTP + WS handlers + tests
│   ├── go.mod
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── views/              HomeView, RoomView
│   │   ├── components/         RoundTable, PlayerSeat, CardDeck, PokerCard
│   │   ├── stores/room.js      WebSocket + global state
│   │   └── router/index.js
│   ├── nginx.conf              Internal container proxy config
│   ├── Dockerfile
│   └── package.json
├── nginx/
│   └── poker.ericsantos.eu.conf   nginx server block for the VPS
└── docker-compose.yml
```

---

## Running locally

### Prerequisites

- Go 1.22+
- Node.js 20+

### Backend

```bash
cd backend
go run ./cmd/server
# Server listens on http://localhost:8080
```

### Frontend

`vite.config.js` is already configured to proxy `/api/` → `localhost:8080`, so no extra setup is needed.

```bash
cd frontend
npm install
npm run dev
# App available at http://localhost:5173
```

### Running Go tests

```bash
cd backend
go test ./...
```

---

## Running with Docker

```bash
# From the project root
docker compose up --build
# Open http://localhost:3001
```

To stop:

```bash
docker compose down
```

---

## Deploying to a VPS

### One-time setup on the server

```bash
# 1. Install Docker (if not already present)
apt update && apt install -y docker.io docker-compose-plugin git

# 2. Clone the repository
git clone https://github.com/<your-username>/go-pokercards.git /opt/go-pokercards

# 3. First build
cd /opt/go-pokercards
docker compose up -d --build
```

### nginx on the VPS

Copy the config and enable the subdomain:

```bash
cp nginx/poker.ericsantos.eu.conf /etc/nginx/sites-available/
ln -s /etc/nginx/sites-available/poker.ericsantos.eu.conf /etc/nginx/sites-enabled/

# Issue a TLS certificate with Certbot
certbot --nginx -d poker.ericsantos.eu

nginx -s reload
```

> The `nginx/poker.ericsantos.eu.conf` file already contains the HTTP→HTTPS redirect and the `Upgrade`/`Connection` headers required for WebSocket to work through nginx.

---

## CI/CD with GitHub Actions

On every `git push origin main`:

1. **`test` job** — runs `go test ./...` on a GitHub-hosted runner.
2. **`deploy` job** (only runs if tests pass) — SSHes into the VPS, runs `git pull`, and rebuilds the containers.

### Required GitHub secrets

Go to **Settings → Secrets and variables → Actions** in your repository and add:

| Secret | Value |
|---|---|
| `DROPLET_HOST` | VPS IP address or hostname |
| `DROPLET_USER` | SSH username (e.g. `root`) |
| `DROPLET_SSH_KEY` | Full content of your SSH private key (e.g. `~/.ssh/id_rsa`) |

> The project path on the VPS is hardcoded as `/opt/go-pokercards` in the workflow. If you cloned it elsewhere, edit the `cd` line in [.github/workflows/deploy.yml](.github/workflows/deploy.yml).

### Generating a dedicated deploy key (recommended)

```bash
# Create a key pair just for deployments
ssh-keygen -t ed25519 -C "github-actions-deploy" -f ~/.ssh/deploy_key

# Authorise the public key on the VPS
ssh-copy-id -i ~/.ssh/deploy_key.pub root@<vps-ip>

# Copy the private key content into the DROPLET_SSH_KEY secret
cat ~/.ssh/deploy_key
```

---

## WebSocket protocol

The connection is opened at `GET /api/rooms/:id/ws?name=Alice[&player_id=uuid]`.

### Server → Client

| `type` | When | Extra fields |
|---|---|---|
| `init` | Immediately after connecting | `player_id`, `room` |
| `state` | After any state change | `room` |
| `error` | Invalid action attempted | `message` |

### Client → Server

| `action` | Who can send | Extra fields |
|---|---|---|
| `vote` | Any player | `value` (e.g. `"8"`) |
| `reveal` | Scrum Master only | — |
| `reset` | Scrum Master only | — |
| `set_master` | Room creator only | `player_id` |

### `room` object

```jsonc
{
  "id": "a1b2c3d4",
  "creator_id": "uuid",
  "master_id":  "uuid",
  "revealed":   false,
  "round":      1,
  "players": [
    { "id": "uuid", "name": "Alice", "has_voted": true,  "vote": "" },
    { "id": "uuid", "name": "Bob",   "has_voted": false, "vote": "" }
  ]
}
// When revealed=true, the "vote" field is populated with each player's chosen value.
```

---

## Design decisions

### Dedicated writer goroutine per connection

gorilla/websocket **does not allow concurrent writes** on the same connection. The standard solution is a dedicated goroutine per client that drains a channel and writes serially:

```go
// session/client.go
func (c *Client) writePump() {
    for msg := range c.send {   // blocks until a message arrives
        c.conn.WriteJSON(msg)
    }
}
```

Broadcasts simply send to the channel (`c.send <- msg`) with no additional synchronisation.

### `PrepareClient` before `RegisterClient`

A new client must receive `init` (containing their own player ID) **before** any broadcast. The handler therefore creates the client, sends init, and only then registers it in the session:

```go
client := sess.PrepareClient(conn)  // writer goroutine starts here
client.Send(initMsg)                // init is queued before any broadcast can arrive
sess.RegisterClient(playerID, client)
sess.BroadcastState()               // now visible to all connected clients
```

### Pure-CSS round table

No layout library. Each player's position is calculated with basic trigonometry:

```js
// components/RoundTable.vue
function seatStyle(i, total) {
  const angle = (i / total) * 2 * Math.PI - Math.PI / 2  // start at the top
  const r = 210  // radius in px
  return {
    transform: `translate(calc(-50% + ${Math.round(r * Math.cos(angle))}px),
                          calc(-50% + ${Math.round(r * Math.sin(angle))}px))`
  }
}
```

### Votes hidden until reveal

The backend's `ToView()` method never exposes a vote value while `Revealed == false` — it only sets `has_voted: true`. This ensures that even if someone inspects network traffic they cannot see other players' votes before the reveal.

### In-memory state

There is no database. All room state lives in the Go process. This keeps the project simple for learning purposes, but it means rooms are lost on restart. Replacing the in-memory `Store` with a Redis-backed implementation would be a natural next step for a production system.
