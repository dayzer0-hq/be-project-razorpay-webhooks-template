## What this PR does

## Issues closed

Closes #

## Checklist

- [ ] POST /v1/webhooks is idempotent (same event_id + merchant_id = same delivery_id)
- [ ] Backoff schedule correct: 15s → 1min → 5min → 30min → permanently_failed
- [ ] X-Razorpay-Signature is computed and verifiable
- [ ] Worker resumes pending deliveries after restart
- [ ] `go test -race ./...` passes
