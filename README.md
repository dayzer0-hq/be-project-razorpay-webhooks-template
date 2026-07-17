# Razorpay Webhook Delivery Engine — DayZer0 Template

Starter for the **"Build Razorpay's webhook delivery engine"** project at DayZer0.

## Get started

```bash
docker compose up -d
psql $DATABASE_URL < db/schema.sql
cp .env.example .env

# Terminal 1 — API server
go run ./cmd/server

# Terminal 2 — delivery worker
go run ./cmd/worker
```

## What you need to build

- `internal/handlers/webhooks.go` — Enqueue (idempotent insert), GetDelivery, ForceRetry
- `internal/worker/delivery.go` — processJob (deliver + retry scheduling), Poll

The `signPayload`, `deliver`, `BackoffSchedule`, and `MaxAttempts` are pre-implemented.

## Submitting

Open a PR to `main`. Include `Closes #N` for every issue. Raj reviews automatically.
