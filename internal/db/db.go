package db

import (
	"context"
	"fmt"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
)
var Pool *pgxpool.Pool
func Init(ctx context.Context) error {
	url := os.Getenv("DATABASE_URL")
	if url == "" { return fmt.Errorf("DATABASE_URL not set") }
	p, err := pgxpool.New(ctx, url)
	if err != nil { return fmt.Errorf("pgxpool.New: %w", err) }
	if err := p.Ping(ctx); err != nil { return fmt.Errorf("db ping: %w", err) }
	Pool = p; return nil
}
