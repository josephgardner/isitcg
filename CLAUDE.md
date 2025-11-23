# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands

```bash
# Build
go build -o isitcg .

# Run (requires Redis)
PORT=5000 REDIS_URL=redis://localhost:6379 go run .

# Run tests
go test ./...
go test -v ./internal/isitcg/...  # Internal package only
```

**Environment Variables:**
- `REDIS_URL` - Redis connection string (required for analytics)
- `PORT` - HTTP port (defaults to 5000)

## Architecture

### Core Flow
1. User submits product name + comma-separated ingredients via form
2. Product data is JSON-serialized, DEFLATE-compressed, and Base64-encoded into a URL hash
3. `/view/{hash}` decompresses and analyzes ingredients against rules
4. Results show matched rules ranked by priority (danger > warning > good)
5. Analytics tracked in Redis sorted sets

### Key Components

**`internal/isitcg/`**
- `ingredients.go` - `IngredientHandler` interface with fuzzy matching logic
- `rule.go` - Loads rules from `ingredientrules.json` (24 categories, 1000+ ingredients)
- `compress.go` - DEFLATE + Base64 URL encoding for shareable results
- `counter.go` - Redis `ZIncrBy` tracking for products and unknown ingredients
- `results.go` - Result aggregation with priority ranking

**Root files**
- `main.go` - Entry point, initializes Redis client and dependencies
- `router.go` - gorilla/mux routes (GET/POST `/`, `/view/{hash}`, `/edit/{hash}`)
- `render.go` - HTML template rendering

### Data Structures

**Rule** - Each rule has name, description, result type (danger/warning/good), rank, and ingredient list

**Results** - Contains ProductName, overall Result, sorted MatchingRules, and Remainder (unmatched ingredients)

### Redis Keys
- `products` (ZSET) - Product names with search counts
- `ingredients:unknown` (ZSET) - Unmatched ingredients with counts

## Deployment

Deployed on Dokku (Digital Ocean droplet at 143.198.164.100).

```bash
# Push to dokku
git push dokku main

# Dokku commands
dokku logs isitcg
dokku ps:rebuild isitcg
dokku redis:info isitcg-redis
```
