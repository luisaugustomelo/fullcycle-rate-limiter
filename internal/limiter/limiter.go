package limiter

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type RateLimiter struct {
	Store     StoreStrategy
	IPLimit   int
	BlockTime int
	Tokens    map[string]int
}

func NewRateLimiter(store StoreStrategy) *RateLimiter {
	ipLimit, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_IP"))
	blockTime, _ := strconv.Atoi(os.Getenv("RATE_BLOCK_DURATION_SECONDS"))

	// Load tokens prefixed with RATE_LIMIT_TOKEN_
	tokens := make(map[string]int)
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "RATE_LIMIT_TOKEN_") {
			parts := strings.SplitN(env, "=", 2)
			token := strings.TrimPrefix(parts[0], "RATE_LIMIT_TOKEN_")
			val, _ := strconv.Atoi(parts[1])
			tokens[token] = val
		}
	}

	return &RateLimiter{
		Store:     store,
		IPLimit:   ipLimit,
		BlockTime: blockTime,
		Tokens:    tokens,
	}
}

func (rl *RateLimiter) GetLimit(token string) int {
	if limit, ok := rl.Tokens[token]; ok {
		return limit
	}
	return rl.IPLimit
}

func (rl *RateLimiter) GetKey(ip, token string) string {
	if token != "" {
		return fmt.Sprintf("token:%s", token)
	}
	return fmt.Sprintf("ip:%s", ip)
}
