package isitcg

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

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
		// Track unique IPs with HyperLogLog
		hllKey := fmt.Sprintf("products:hll:%s", results.ProductName)
		if err := c.db.PFAdd(ctx, hllKey, clientIP).Err(); err != nil {
			log.Printf("failed to add IP to HLL for product %s: %v", results.ProductName, err)
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
