package core

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExtractIP(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]string
		remote   string
		wantIP   string
	}{
		{
			name:   "X-Forwarded-For single IP",
			headers: map[string]string{"X-Forwarded-For": "1.2.3.4"},
			remote:  "10.0.0.1:12345",
			wantIP:  "1.2.3.4",
		},
		{
			name:   "X-Forwarded-For multiple IPs takes first",
			headers: map[string]string{"X-Forwarded-For": "1.2.3.4, 5.6.7.8"},
			remote:  "10.0.0.1:12345",
			wantIP:  "1.2.3.4",
		},
		{
			name:   "X-Forwarded-For with whitespace",
			headers: map[string]string{"X-Forwarded-For": "  1.2.3.4  , 5.6.7.8"},
			remote:  "10.0.0.1:12345",
			wantIP:  "1.2.3.4",
		},
		{
			name:   "X-Real-IP when no X-Forwarded-For",
			headers: map[string]string{"X-Real-IP": "1.2.3.4"},
			remote:  "10.0.0.1:12345",
			wantIP:  "1.2.3.4",
		},
		{
			name:   "X-Forwarded-For takes precedence over X-Real-IP",
			headers: map[string]string{
				"X-Forwarded-For": "1.2.3.4",
				"X-Real-IP":       "5.6.7.8",
			},
			remote: "10.0.0.1:12345",
			wantIP: "1.2.3.4",
		},
		{
			name:   "RemoteAddr when no proxy headers",
			headers: map[string]string{},
			remote:  "10.0.0.1:12345",
			wantIP:  "10.0.0.1",
		},
		{
			name:   "RemoteAddr without port",
			headers: map[string]string{},
			remote:  "10.0.0.1",
			wantIP:  "10.0.0.1",
		},
		{
			name:   "IPv6 RemoteAddr with port",
			headers: map[string]string{},
			remote:  "[::1]:12345",
			wantIP:  "[::1]",
		},
		{
			name:   "empty X-Forwarded-For falls through to RemoteAddr",
			headers: map[string]string{"X-Forwarded-For": ""},
			remote:  "10.0.0.1:12345",
			wantIP:  "10.0.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = tt.remote
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			got := extractIP(req)
			if got != tt.wantIP {
				t.Errorf("extractIP() = %q, want %q", got, tt.wantIP)
			}
		})
	}
}

func TestRateLimiterAllowsBurst(t *testing.T) {
	limiter := newIPLimiters(5, 10)
	ip := "192.168.1.1"

	for i := 0; i < 10; i++ {
		l := limiter.getLimiter(ip)
		if !l.Allow() {
			t.Errorf("request %d within burst should be allowed", i+1)
		}
	}

	l := limiter.getLimiter(ip)
	if l.Allow() {
		t.Error("request beyond burst should be rejected")
	}
}

func TestRateLimiterSeparateIPs(t *testing.T) {
	limiter := newIPLimiters(5, 2)
	ip1 := "192.168.1.1"
	ip2 := "192.168.1.2"

	l1 := limiter.getLimiter(ip1)
	l1.Allow()
	l1.Allow()

	if l1.Allow() {
		t.Error("ip1 should be rate limited after burst")
	}

	l2 := limiter.getLimiter(ip2)
	if !l2.Allow() {
		t.Error("ip2 should not be affected by ip1's rate limit")
	}
}

func TestRateLimiterSweep(t *testing.T) {
	limiter := newIPLimiters(5, 10)
	limiter.getLimiter("old-ip")
	limiter.getLimiter("new-ip")

	limiter.mu.Lock()
	limiter.limiters["old-ip"].lastSeen = limiter.limiters["old-ip"].lastSeen.Add(-2 * rateLimiterMaxAge)
	limiter.mu.Unlock()

	limiter.sweepOlderThan(rateLimiterMaxAge)

	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	if _, exists := limiter.limiters["old-ip"]; exists {
		t.Error("old-ip should have been evicted")
	}
	if _, exists := limiter.limiters["new-ip"]; !exists {
		t.Error("new-ip should still exist")
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	AuthRateLimiter = newIPLimiters(1, 2)

	handler := RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "1.2.3.4:5678"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("request %d should succeed, got %d", i+1, w.Code)
		}
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "1.2.3.4:5678"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("request beyond burst should return 429, got %d", w.Code)
	}

	AuthRateLimiter = newIPLimiters(5, 10)
}
