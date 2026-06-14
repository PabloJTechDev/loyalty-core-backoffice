# loyalty-core-backoffice

Go microservice — backoffice core for the loyalty platform.

Part of the **backoffice vertical**: `web-backoffice` → `bff-backoffice` → **`core-backoffice`**

---

## Responsibilities

- Expose operational **capabilities** available in the loyalty program
- Expose **operational alerts** for support and monitoring teams
- Act as the authoritative backend for backoffice domain metadata

The heavy lifting (points data, customer profiles) lives in `core-points`. This service handles the operational metadata layer.

---

## Endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/health` | Service health + status |
| `GET` | `/metrics` | Prometheus metrics |
| `GET` | `/v1/backoffice-capabilities` | List available backoffice capabilities |
| `GET` | `/v1/operational-alerts` | List active operational alerts |

---

## Architecture

Clean Architecture with Go journey packages under `internal/`:

```
internal/
  shared/          → WriteJSON, WithMetrics, LogEvent, Prometheus metrics
  health/          → handler + types  (GET /health)
  capabilities/    → handler + types  (GET /v1/backoffice-capabilities)
  alerts/          → handler + types  (GET /v1/operational-alerts)
main.go            → route wiring, graceful shutdown
```

Each journey package follows the same layering:

| File | Responsibility |
|---|---|
| `types.go` | Response structs and domain defaults |
| `handler.go` | HTTP handler (calls defaults or use cases, writes JSON) |

Dependency injection is done by constructor functions in `main.go` — no framework.

---

## Tech stack

- **Go** (stdlib `net/http`)
- **Prometheus** client (custom registry, counters + histograms)
- Structured JSON logging

---

## Running locally

```bash
cp .env.example .env
# Edit .env: set PORT (default 3005)

go run .
# or: go build -o core-backoffice . && ./core-backoffice
```

### Environment variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `3005` | HTTP listen port |
| `APP_ENV` | `development` | Environment tag |

---

## Health check

```bash
curl http://localhost:3005/health
# {"status":"ok","service":"core-backoffice"}
```

---

## Part of loyalty-platform

See the [monorepo root](https://github.com/PabloJTechDev/loyalty-platform) for the full architecture, port map, and Docker Compose setup.
