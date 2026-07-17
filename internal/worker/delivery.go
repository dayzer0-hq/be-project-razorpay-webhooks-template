package worker

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// BackoffSchedule defines wait durations between attempts (index = attemptCount after failure).
// attempt 1: immediate, attempt 2: 15s, attempt 3: 1min, attempt 4: 5min, attempt 5: 30min → give up
var BackoffSchedule = []time.Duration{
	15 * time.Second,
	1 * time.Minute,
	5 * time.Minute,
	30 * time.Minute,
	2 * time.Hour,
}

const MaxAttempts = 5

// signPayload returns the HMAC-SHA256 signature string for the Razorpay webhook header.
func signPayload(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

type job struct {
	deliveryID string
	url        string
	secret     string
	payload    any
	attemptNo  int
}

// Worker is a bounded pool that processes delivery jobs.
type Worker struct {
	db   *pgxpool.Pool
	jobs chan job
	wg   sync.WaitGroup
}

func NewWorker(db *pgxpool.Pool, concurrency int) *Worker {
	w := &Worker{db: db, jobs: make(chan job, 1000)}
	for i := 0; i < concurrency; i++ {
		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			for j := range w.jobs {
				w.processJob(context.Background(), j)
			}
		}()
	}
	return w
}

func (w *Worker) processJob(ctx context.Context, j job) {
	// TODO: call deliver(ctx, j.url, j.secret, j.payload)
	// On success (2xx): UPDATE webhook_deliveries SET status='delivered'
	//                   INSERT delivery_attempts with http_status, response
	// On failure:
	//   if j.attemptNo < MaxAttempts:
	//     delay = BackoffSchedule[j.attemptNo-1]
	//     UPDATE webhook_deliveries SET deliver_after=NOW()+delay, attempt_count++
	//   else:
	//     UPDATE webhook_deliveries SET status='permanently_failed'
	log.Printf("TODO: process delivery %s attempt %d", j.deliveryID, j.attemptNo)
}

func deliver(ctx context.Context, url, secret string, payload any) (int, string, error) {
	body, err := json.Marshal(payload)
	if err != nil { return 0, "", err }

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil { return 0, "", err }

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Razorpay-Signature", signPayload(secret, body))
	req.Header.Set("User-Agent", "Razorpay-Webhooks/1.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil { return 0, "", err }
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
	return resp.StatusCode, string(respBytes), nil
}

// Poll fetches pending deliveries and enqueues them for processing.
// Run this on a ticker (every 5 seconds) from cmd/worker/main.go.
func (w *Worker) Poll(ctx context.Context) {
	// TODO: SELECT id, merchant_id, payload, attempt_count
	//       FROM webhook_deliveries
	//       WHERE status='pending' AND deliver_after <= NOW()
	//       LIMIT 100
	// TODO: For each: fetch merchant webhook URL + secret from merchant_webhooks
	//       Enqueue to w.jobs
}

func (w *Worker) Stop() { close(w.jobs); w.wg.Wait() }
