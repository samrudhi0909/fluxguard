package limiter

import (
	"context"
	_ "embed" // This allows us to embed files
	"time"

	"github.com/redis/go-redis/v9"
)

// Go Magic: This reads script.lua and puts it into this string variable at compile time!
//
//go:embed script.lua
var limiterScript string

// RateLimiterService defines the methods we expose
type RateLimiterService struct {
	rdb *redis.Client
}

// NewRateLimiter creates a new instance
func NewRateLimiter(rdb *redis.Client) *RateLimiterService {
	return &RateLimiterService{
		rdb: rdb,
	}
}

// Allow checks if a request should be allowed
func (s *RateLimiterService) Allow(ctx context.Context, userID string, capacity int32, rate int32) (bool, error) {
	now := time.Now().Unix() // Current time in seconds

	// We send the script to Redis
	// Keys: [userID]
	// Args: [capacity, rate, now, 1 (cost)]
	result, err := s.rdb.Eval(ctx, limiterScript, []string{userID}, capacity, rate, now, 1).Result()
	
	if err != nil {
		// FAIL-OPEN STRATEGY:
		// If Redis is down, we log the error but return TRUE.
		// Why? Because we don't want to block users just because our tool is broken.
		return true, err
	}

	// Redis returns 1 (true) or 0 (false)
	return result.(int64) == 1, nil
}