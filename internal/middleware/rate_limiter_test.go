package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luisaugustomelo/fullcycle-rate-limiter/internal/limiter"
)

func setupTestRouter() *gin.Engine {
	_ = os.Setenv("RATE_LIMIT_IP", "2")                   
	_ = os.Setenv("RATE_BLOCK_DURATION_SECONDS", "10")  
	_ = os.Setenv("RATE_LIMIT_TOKEN_TESTTOKEN", "2")     
	redisStore := limiter.NewRedisStrategy("localhost:6379")

	_ = redisStore.FlushDB()
	rl := limiter.NewRateLimiter(redisStore)

	r := gin.Default()
	r.Use(RateLimitMiddleware(rl))
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Request accepted"})
	})

	return r
}

func TestRateLimiter_IPBased(t *testing.T) {
	router := setupTestRouter()

	// 1st Request (should pass)
	req1 := httptest.NewRequest("GET", "/", nil)
	resp1 := httptest.NewRecorder()
	router.ServeHTTP(resp1, req1)
	if resp1.Code != 200 {
		t.Fatalf("Expected 200 on request #1, got %d", resp1.Code)
	}

	// 2nd Request (should pass)
	req2 := httptest.NewRequest("GET", "/", nil)
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)
	if resp2.Code != 200 {
		t.Fatalf("Expected 200 on request #2, got %d", resp2.Code)
	}

	// 3rd Request (should be blocked)
	req3 := httptest.NewRequest("GET", "/", nil)
	resp3 := httptest.NewRecorder()
	router.ServeHTTP(resp3, req3)
	if resp3.Code != http.StatusTooManyRequests {
		t.Fatalf("Expected 429 on request #3, got %d", resp3.Code)
	}
}


func TestRateLimiter_TokenBased(t *testing.T) {
	_ = os.Setenv("RATE_LIMIT_TOKEN_TESTTOKEN", "2")

	router := setupTestRouter()

	for i := 1; i <= 2; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("API_KEY", "TESTTOKEN")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected 200 on token request #%d, got %d", i, w.Code)
		}
	}

	// 3rd request with token should fail
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("API_KEY", "TESTTOKEN")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected 429 on token request #3, got %d", w.Code)
	}
}

func TestRateLimiter_BlockExpires(t *testing.T) {
	router := setupTestRouter()

	// Consome o limite
	for i := 1; i <= 2; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// Vai ser bloqueado
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("Expected 429 immediately after limit exceeded, got %d", w.Code)
	}

	// Espera o tempo de expiração
	time.Sleep(11 * time.Second) // Espera +1s de margem

	// Deve passar agora
	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Expected 200 after block expired, got %d", w.Code)
	}
}


func TestRateLimiter_ExceedIPLimit(t *testing.T) {
	router := setupTestRouter()

	// 1st Request (should pass)
	req1 := httptest.NewRequest("GET", "/", nil)
	resp1 := httptest.NewRecorder()
	router.ServeHTTP(resp1, req1)
	if resp1.Code != 200 {
		t.Fatalf("Expected 200, got %d", resp1.Code)
	}

	// 2nd Request (should pass)
	req2 := httptest.NewRequest("GET", "/", nil)
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)
	if resp2.Code != 200 {
		t.Fatalf("Expected 200, got %d", resp2.Code)
	}

	// 3rd Request (should be blocked)
	req3 := httptest.NewRequest("GET", "/", nil)
	resp3 := httptest.NewRecorder()
	router.ServeHTTP(resp3, req3)
	if resp3.Code != http.StatusTooManyRequests {
		t.Fatalf("Expected 429, got %d", resp3.Code)
	}
}

func TestRateLimiter_ExceedTokenLimit(t *testing.T) {
	router := setupTestRouter()

	// 1st Request with token (should pass)
	req1 := httptest.NewRequest("GET", "/", nil)
	req1.Header.Set("API_KEY", "TESTTOKEN")
	resp1 := httptest.NewRecorder()
	router.ServeHTTP(resp1, req1)
	if resp1.Code != 200 {
		t.Fatalf("Expected 200, got %d", resp1.Code)
	}

	// 2nd Request with token (should pass)
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.Header.Set("API_KEY", "TESTTOKEN")
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)
	if resp2.Code != 200 {
		t.Fatalf("Expected 200, got %d", resp2.Code)
	}

	// 3rd Request with token (should be blocked)
	req3 := httptest.NewRequest("GET", "/", nil)
	req3.Header.Set("API_KEY", "TESTTOKEN")
	resp3 := httptest.NewRecorder()
	router.ServeHTTP(resp3, req3)
	if resp3.Code != http.StatusTooManyRequests {
		t.Fatalf("Expected 429, got %d", resp3.Code)
	}
}