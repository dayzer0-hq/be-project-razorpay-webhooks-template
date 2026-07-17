package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"
	"github.com/joho/godotenv"
	"github.com/dayzer0-hq/be-project-razorpay-webhooks-template/internal/db"
	"github.com/dayzer0-hq/be-project-razorpay-webhooks-template/internal/worker"
)
func main() {
	_ = godotenv.Load()
	if err := db.Init(context.Background()); err != nil { log.Fatal(err) }
	concurrency, _ := strconv.Atoi(os.Getenv("WEBHOOK_WORKERS"))
	if concurrency <= 0 { concurrency = 50 }
	w := worker.NewWorker(db.Pool, concurrency)
	log.Printf("worker started with %d goroutines", concurrency)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		w.Poll(context.Background())
	}
}
