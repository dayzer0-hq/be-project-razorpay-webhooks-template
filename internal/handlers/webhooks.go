package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/v1/webhooks", Enqueue)
	r.Get("/v1/webhooks/{deliveryID}", GetDelivery)
	r.Post("/v1/webhooks/{deliveryID}/retry", ForceRetry)
	return r
}

type enqueueRequest struct {
	EventType  string         `json:"event_type"`
	MerchantID string         `json:"merchant_id"`
	EventID    string         `json:"event_id"`
	Payload    map[string]any `json:"payload"`
}

// Enqueue handles POST /v1/webhooks
// Idempotent on event_id + merchant_id.
// Returns 202 with { delivery_id }.
func Enqueue(w http.ResponseWriter, r *http.Request) {
	var req enqueueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest); return
	}
	// TODO: Validate req fields non-empty
	// TODO: INSERT INTO webhook_deliveries (event_id, merchant_id, event_type, payload)
	//       ON CONFLICT (event_id, merchant_id) DO NOTHING
	//       RETURNING id
	// If 0 rows inserted (duplicate), look up existing delivery_id and return it
	// Return 202 { delivery_id }
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"error": "TODO: implement Enqueue"})
}

// GetDelivery handles GET /v1/webhooks/{deliveryID}
func GetDelivery(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "deliveryID")
	// TODO: SELECT delivery + all delivery_attempts WHERE delivery_id = id
	// Return { delivery_id, status, attempt_count, attempts: [...] }
	_ = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"error": "TODO: implement GetDelivery"})
}

// ForceRetry handles POST /v1/webhooks/{deliveryID}/retry (admin use)
func ForceRetry(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "deliveryID")
	// TODO: UPDATE webhook_deliveries SET status='pending', deliver_after=NOW()
	//       WHERE id = deliveryID AND status = 'permanently_failed'
	// Return 204 on success, 404 if not found or status != permanently_failed
	_ = id
	w.WriteHeader(http.StatusNoContent)
}
