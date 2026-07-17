package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"github.com/dayzer0-hq/be-project-razorpay-webhooks-template/internal/db"
	"github.com/dayzer0-hq/be-project-razorpay-webhooks-template/internal/handlers"
)
func main() {
	_ = godotenv.Load()
	if err := db.Init(context.Background()); err != nil { log.Fatal(err) }
	port := os.Getenv("PORT"); if port == "" { port = "8080" }
	log.Printf("server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.NewRouter()))
}
