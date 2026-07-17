CREATE TABLE IF NOT EXISTS merchant_webhooks (
  merchant_id TEXT PRIMARY KEY,
  url         TEXT NOT NULL,
  secret      TEXT NOT NULL,
  active      BOOLEAN DEFAULT TRUE,
  created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS webhook_deliveries (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_id      TEXT NOT NULL,
  merchant_id   TEXT NOT NULL,
  event_type    TEXT NOT NULL,
  payload       JSONB NOT NULL,
  status        TEXT NOT NULL DEFAULT 'pending',  -- pending | delivered | permanently_failed
  attempt_count INT DEFAULT 0,
  deliver_after TIMESTAMPTZ DEFAULT NOW(),
  created_at    TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE (event_id, merchant_id)
);

CREATE TABLE IF NOT EXISTS delivery_attempts (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  delivery_id  UUID REFERENCES webhook_deliveries(id),
  attempt_no   INT NOT NULL,
  http_status  INT,
  response     TEXT,
  attempted_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_deliveries_pending ON webhook_deliveries(status, deliver_after)
  WHERE status = 'pending';
