package isitcg

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Submission represents data pushed to analytics inbox
type Submission struct {
	ProductName        string   `json:"product_name"`
	Ingredients        string   `json:"ingredients"`
	IPHash             string   `json:"ip_hash"`
	Timestamp          int64    `json:"timestamp"`
	UnknownIngredients []string `json:"unknown_ingredients"`
}

// Size limits
const (
	maxProductName = 100
	maxIngredients = 5000
	maxUnknown     = 50
	maxInboxSize   = 10000 // Max submissions in inbox before oldest are dropped
)

type Counter interface {
	Count(ctx context.Context, product Product, results Results, clientIP string) error
}

func NewRedisCounter(db *redis.Client) Counter {
	return &redisCounter{db: db}
}

type redisCounter struct {
	db *redis.Client
}

func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen]
	}
	return s
}

func truncateSlice(s []string, maxLen int) []string {
	if len(s) > maxLen {
		return s[:maxLen]
	}
	return s
}

func (c *redisCounter) Count(ctx context.Context, product Product, results Results, clientIP string) error {
	if product.Name == "" || product.Ingredients == "" {
		return nil
	}

	// Hash client IP for privacy
	ipHash := fmt.Sprintf("%x", sha256.Sum256([]byte(clientIP)))

	// Build submission
	submission := Submission{
		ProductName:        truncateString(product.Name, maxProductName),
		Ingredients:        truncateString(product.Ingredients, maxIngredients),
		IPHash:             ipHash,
		Timestamp:          time.Now().Unix(),
		UnknownIngredients: truncateSlice(results.Remainder, maxUnknown),
	}

	// Marshal to JSON
	data, err := json.Marshal(submission)
	if err != nil {
		log.Printf("failed to marshal submission: %v", err)
		return err
	}

	// Push to inbox
	if err := c.db.LPush(ctx, "submissions:inbox", data).Err(); err != nil {
		log.Printf("failed to push to submissions inbox: %v", err)
		return err
	}

	// Trim to max size (keep newest, drop oldest)
	if err := c.db.LTrim(ctx, "submissions:inbox", 0, maxInboxSize-1).Err(); err != nil {
		log.Printf("failed to trim submissions inbox: %v", err)
		// Don't return error - push succeeded
	}

	return nil
}
