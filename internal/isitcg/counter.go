package isitcg

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// normalizeProductName normalizes a product name for consistent hashing
func normalizeProductName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	// Collapse multiple spaces into one
	fields := strings.Fields(name)
	return strings.Join(fields, " ")
}

type Counter interface {
	Count(ctx context.Context, results Results, clientIP string) error
}

func NewRedisCounter(db *redis.Client) Counter {
	return &redisCounter{db: db}
}

type redisCounter struct {
	db *redis.Client
}

func (c *redisCounter) Count(ctx context.Context, results Results, clientIP string) error {
	if results.ProductName != "" {
		// Normalize and hash product name for consistent keys
		normalized := normalizeProductName(results.ProductName)
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(normalized)))
		hllKey := fmt.Sprintf("products:hll:%s", hash)

		// Track unique IPs with HyperLogLog
		if err := c.db.PFAdd(ctx, hllKey, clientIP).Err(); err != nil {
			log.Printf("failed to add IP to HLL for product %s: %v", results.ProductName, err)
		}

		// Set rolling TTL (90 days) - resets on each access
		if err := c.db.Expire(ctx, hllKey, 90*24*time.Hour).Err(); err != nil {
			log.Printf("failed to set TTL for HLL key %s: %v", hllKey, err)
		}

		// Update last-seen timestamp
		if err := c.db.ZAdd(ctx, "products:recent", redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: results.ProductName,
		}).Err(); err != nil {
			log.Printf("failed to update timestamp for product %s: %v", results.ProductName, err)
		}
	}

	for _, ingredient := range results.Remainder {
		if ingredient != "" {
			if err := c.db.ZIncrBy(ctx, "ingredients:unknown", 1, ingredient).Err(); err != nil {
				log.Printf("failed to increment unknown ingredient count for %s: %v", ingredient, err)
				return err
			}
		}
	}

	return nil
}
