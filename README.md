# Distributed Task Scheduler

A production-grade distributed task scheduling system built with Next.js, Go, Redis, and PostgreSQL. Supports cron-based scheduling, priority queues, real-time updates via WebSockets, OAuth authentication, and horizontal worker scaling.

---

## Architecture

```
Next.js (Frontend)
      ↓
Go Fiber (API Gateway)
      ↓
┌─────────────────────────────────┐
│ Scheduler Service (Go)          │
│ Worker Nodes (Go)               │
│ Message Queue (Redis)           │
│ Database (PostgreSQL)           │
└─────────────────────────────────┘
      ↓
Docker + Kubernetes
```

### Services

| Service | Description |
|---|---|
| `api` | REST API gateway — handles auth, CRUD, WebSocket hub |
| `scheduler` | Polls DB every 5s, dispatches due tasks to Redis queue |
| `worker` | Consumes jobs from Redis, executes tasks, writes logs |
| `postgres` | Stores tasks, executions, logs, users |
| `redis` | Job queue (sorted set by priority) + pub/sub for real-time events |

---

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | Next.js 14, Tailwind CSS (via globals.css) |
| Backend API | Go, Fiber v2 |
| Scheduler | Go |
| Workers | Go |
| Queue | Redis (sorted set + pub/sub) |
| Database | PostgreSQL |
| Auth | JWT + Google OAuth + GitHub OAuth |
| Real-time | WebSockets |
| Deployment | Docker, Kubernetes, Render, AWS EC2 |

---

## Features

- **Cron scheduling** — standard 5-field cron expressions
- **Priority queue** — high, normal, low priority via Redis sorted sets
- **Retry policies** — exponential backoff with configurable max retries
- **Distributed locking** — Redis-based locks prevent duplicate execution
- **Real-time updates** — WebSocket + Redis pub/sub replaces polling
- **Execution logs** — per-execution log entries with info/warn/error levels
- **OAuth** — Google and GitHub login alongside email/password
- **JWT auth** — HttpOnly cookies, 24h expiry
- **Worker scaling** — run multiple worker instances, Redis handles distribution

---

## Project Structure

```
/distributed-task-scheduler
  /backend
    /cmd
      /api            ← API server entry point
      /scheduler      ← Scheduler service entry point
      /worker         ← Worker service entry point
    /internal
      /auth           ← JWT, bcrypt, OAuth, middleware
      /db             ← PostgreSQL connection pool
      /handlers       ← Fiber route handlers
      /lock           ← Redis distributed locking
      /models         ← Go structs
      /queue          ← Redis queue (enqueue/dequeue)
      /retry          ← Exponential backoff
      /services       ← Business logic
      /websockets     ← WebSocket hub
    /migrations       ← SQL migration files
    /pkg/utils        ← Shared utilities (cron parser etc.)
    go.mod

  /frontend
    /app
      /dashboard      ← Task overview
      /tasks/new      ← Create task
      /tasks/[id]     ← Task detail + logs
      /login          ← Login page
      /register       ← Register page
    /components       ← Sidebar, StatusBadge, ConfirmDialog etc.
    /hooks            ← useWebSocket, useAuth
    /services         ← API calls, types

  /k8s
    /api
    /scheduler
    /worker
    /redis
    /postgres
    hpa.yaml

  /docker
    /api
    /scheduler
    /worker

  /monitoring
    /prometheus
    /grafana

  docker-compose.yml
```

---

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- Docker + Docker Compose
- `golang-migrate` CLI

### 1. Clone the repo

```bash
git clone https://github.com/yourusername/distributed-task-scheduler.git
cd distributed-task-scheduler
```

### 2. Set up environment variables

```bash
cp backend/.env.example backend/.env
```

Fill in all values in `backend/.env`. See [Environment Variables](#environment-variables) below.

### 3. Start infrastructure

```bash
docker-compose up -d
```

This starts PostgreSQL and Redis.

### 4. Run database migrations

```bash
cd backend
migrate -path ./migrations -database "YOUR_DATABASE_URL" up
```

### 5. Start the services

Open three terminals:

```bash
# Terminal 1 — API
cd backend/cmd/api && air

# Terminal 2 — Scheduler
cd backend/cmd/scheduler && air

# Terminal 3 — Worker
cd backend/cmd/worker && air
```

### 6. Start the frontend

```bash
cd frontend
npm install
npm run dev
```

Visit `http://localhost:3000`

---

## Environment Variables

Create `backend/.env` from `backend/.env.example`:

```bash
# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5433/task_scheduler?sslmode=disable

# Redis
REDIS_ADDR=localhost:6379

# App
APP_PORT=8080

# JWT
JWT_SECRET=your_super_secret_key_here

# Google OAuth
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback

# GitHub OAuth
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
GITHUB_REDIRECT_URL=http://localhost:8080/auth/github/callback
```

Create `frontend/.env.local`:

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_WEB_SOCKET_URL=localhost:8080
```

### Getting OAuth Credentials

**Google:**
1. Go to [console.cloud.google.com](https://console.cloud.google.com)
2. Create a project → APIs & Services → Credentials
3. Create OAuth 2.0 Client ID → Web application
4. Add `http://localhost:8080/auth/google/callback` as authorized redirect URI

**GitHub:**
1. Go to [github.com/settings/developers](https://github.com/settings/developers)
2. New OAuth App
3. Set callback URL to `http://localhost:8080/auth/github/callback`

---

## API Reference

### Auth
| Method | Endpoint | Description |
|---|---|---|
| POST | `/auth/register` | Register with email + password |
| POST | `/auth/login` | Login, returns HttpOnly cookie |
| POST | `/auth/logout` | Clear auth cookie |
| GET | `/auth/me` | Get current user from cookie |
| GET | `/auth/google` | Redirect to Google OAuth |
| GET | `/auth/google/callback` | Google OAuth callback |
| GET | `/auth/github` | Redirect to GitHub OAuth |
| GET | `/auth/github/callback` | GitHub OAuth callback |

### Tasks
| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/v1/tasks` | List all tasks |
| POST | `/api/v1/tasks` | Create a task |
| GET | `/api/v1/tasks/:id` | Get task with executions + logs |
| DELETE | `/api/v1/tasks/:id` | Delete task |
| POST | `/api/v1/tasks/:id/retry` | Retry a failed task |

### WebSocket
| Endpoint | Description |
|---|---|
| `ws://localhost:8080/ws` | Real-time task status updates |

---

## Task Priority

Tasks are queued using a Redis sorted set. Priority scores:

| Priority | Score | Behavior |
|---|---|---|
| `high` | 3 | Dequeued first |
| `normal` | 2 | Default |
| `low` | 1 | Dequeued last |

---

## Cron Schedule Format

Standard 5-field cron syntax:

```
┌─────── minute (0-59)
│ ┌───── hour (0-23)
│ │ ┌─── day of month (1-31)
│ │ │ ┌─ month (1-12)
│ │ │ │ ┌ day of week (0-6, Sunday=0)
│ │ │ │ │
* * * * *
```

Examples:

| Expression | Meaning |
|---|---|
| `* * * * *` | Every minute |
| `*/5 * * * *` | Every 5 minutes |
| `0 8 * * *` | Daily at 8AM |
| `0 9 * * 1` | Every Monday at 9AM |

---

## Retry Policy

Failed tasks are retried with exponential backoff:

```
delay = retry_delay_seconds × 2^retry_count
```

Example with `retry_delay_seconds = 60` and `max_retries = 3`:

| Attempt | Delay |
|---|---|
| 1st retry | 60s |
| 2nd retry | 120s |
| 3rd retry | 240s |
| Permanently failed | — |

---

## Running Multiple Workers

Each worker binary competes for jobs on the same Redis queue. Start as many as you need:

```bash
# Terminal 3
cd backend/cmd/worker && air

# Terminal 4
cd backend/cmd/worker && air

# Terminal 5
cd backend/cmd/worker && air
```

Redis `BZPOPMAX` ensures each job is only picked up by one worker.

---

## Deployment

See [Deployment Guide](./docs/deployment.md) for:
- Docker build instructions
- Kubernetes manifests
- Render deployment
- AWS EC2 deployment

---

## License

MIT
