package isitcg

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type Counter interface {
	Count(ctx context.Context, results Results) error
}

func NewRedisCounter(db *redis.Client) Counter {
	return &redisCounter{db: db}
}

type redisCounter struct {
	db *redis.Client
}

func (c *redisCounter) Count(ctx context.Context, results Results) error {
	if results.ProductName != "" {
		if err := c.db.ZIncrBy(ctx, "products", 1, results.ProductName).Err(); err != nil {
			log.Printf("failed to count product name %s: %v", results.ProductName, err)
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
