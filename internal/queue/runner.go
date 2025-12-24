package queue

import (
	"context"
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"printy/internal/printer"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// Runner listens to the Redis queue and fetches print data from Postgres.
type Runner struct {
	redisClient  *redis.Client
	postgresDB   *sql.DB
	queueName    string
	pollInterval time.Duration
}

// NewRunner initializes Redis and Postgres clients using environment variables.
// Required env vars:
//   - POSTGRES_CONNECTION (e.g. postgres://user:pass@localhost:5432/db?sslmode=disable)
//
// Optional env vars:
//   - REDIS_HOST (default: localhost:6379)
//   - REDIS_PASSWORD (default: empty)
//   - REDIS_DB (default: 0)
func NewRunner() (*Runner, error) {
	pgDSN := os.Getenv("POSTGRES_CONNECTION")
	if pgDSN == "" {
		return nil, fmt.Errorf("POSTGRES_CONNECTION environment variable is not set")
	}

	postgresDB, err := sql.Open("postgres", pgDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := postgresDB.PingContext(ctx); err != nil {
		postgresDB.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	redisAddr := os.Getenv("REDIS_HOST")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisDB := 0
	if redisDBStr := os.Getenv("REDIS_DB"); redisDBStr != "" {
		if parsed, convErr := strconv.Atoi(redisDBStr); convErr == nil {
			redisDB = parsed
		} else {
			log.Printf("Invalid REDIS_DB value %q, defaulting to 0: %v", redisDBStr, convErr)
		}
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
		// Always use TLS; adjust config here if certificates are needed later.
		TLSConfig: &tls.Config{},
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		postgresDB.Close()
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &Runner{
		redisClient:  redisClient,
		postgresDB:   postgresDB,
		queueName:    "printing_queue",
		pollInterval: time.Second,
	}, nil
}

// Run starts polling the Redis queue every second until the context is done.
func (r *Runner) Run(ctx context.Context) {
	ticker := time.NewTicker(r.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Queue runner stopped: %v", ctx.Err())
			return
		case <-ticker.C:
			r.popAndProcess(ctx)
		}
	}
}

// Close releases database and Redis resources.
func (r *Runner) Close() {
	if r.postgresDB != nil {
		if err := r.postgresDB.Close(); err != nil {
			log.Printf("Error closing Postgres connection: %v", err)
		}
	}

	if r.redisClient != nil {
		if err := r.redisClient.Close(); err != nil {
			log.Printf("Error closing Redis client: %v", err)
		}
	}
}

func (r *Runner) popAndProcess(ctx context.Context) {
	popCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	id, err := r.redisClient.LPop(popCtx, r.queueName).Result()
	switch {
	case errors.Is(err, redis.Nil):
		return // Queue empty, nothing to do this tick.
	case err != nil:
		log.Printf("Failed to pop from Redis: %v", err)
		return
	}

	printID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Invalid print id %q: %v", id, err)
		return
	}

	queryCtx, queryCancel := context.WithTimeout(ctx, 5*time.Second)
	defer queryCancel()

	var (
		rasterData []byte
		imgWidth   int
		imgHeight  int
	)

	err = r.postgresDB.QueryRowContext(
		queryCtx,
		"SELECT raster_data, image_width, image_height FROM prints WHERE id = $1",
		printID,
	).Scan(&rasterData, &imgWidth, &imgHeight)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No print found with id %d", printID)
			return
		}
		log.Printf("Failed to fetch raster_data for id %d: %v", printID, err)
		return
	}

	log.Printf("Matched ticket id: %d (w=%d h=%d)", printID, imgWidth, imgHeight)

	printerName := os.Getenv("PRINTER_NAME")

	imagePrinter := printer.NewImagePrinter(printerName)
	err = imagePrinter.PrintRawRaster(rasterData, imgWidth, imgHeight)
	if err != nil {
		log.Printf("Failed to print raster data: %v", err)
		return
	}

	log.Printf("Printed raster data for ticket %d", printID)
}
