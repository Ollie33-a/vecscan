package utils

import (
    "context"
    "time"
)

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
    tokens    chan struct{}
    stopChan  chan struct{}
    rate      int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(ratePerSecond int) *RateLimiter {
    rl := &RateLimiter{
        tokens:   make(chan struct{}, ratePerSecond),
        stopChan: make(chan struct{}),
        rate:     ratePerSecond,
    }
    
    // Start token refill goroutine
    go rl.refill()
    
    return rl
}

func (rl *RateLimiter) refill() {
    ticker := time.NewTicker(time.Second / time.Duration(rl.rate))
    defer ticker.Stop()
    
    for {
        select {
        case <-rl.stopChan:
            return
        case <-ticker.C:
            select {
            case rl.tokens <- struct{}{}:
            default:
                // Bucket full
            }
        }
    }
}

// Wait blocks until a token is available
func (rl *RateLimiter) Wait() {
    <-rl.tokens
}

// WaitCtx blocks until a token is available or context is cancelled
func (rl *RateLimiter) WaitCtx(ctx context.Context) bool {
    select {
    case <-ctx.Done():
        return false
    case <-rl.tokens:
        return true
    }
}

// Stop stops the rate limiter
func (rl *RateLimiter) Stop() {
    close(rl.stopChan)
}