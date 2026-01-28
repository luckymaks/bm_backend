package rpc_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/luckymaks/bm_backend/backend/internal/rpc"
	"github.com/stretchr/testify/assert"
)

func TestWithCORS(t *testing.T) {
	t.Parallel()
	
	innerHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	
	handler := rpc.WithCORS(innerHandler)
	
	t.Run("adds CORS headers to response", func(t *testing.T) {
		t.Parallel()
		
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Origin", "https://example.com")
		rec := httptest.NewRecorder()
		
		handler.ServeHTTP(rec, req)
		
		assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
	})
	
	t.Run("handles preflight OPTIONS request", func(t *testing.T) {
		t.Parallel()
		
		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		req.Header.Set("Origin", "https://example.com")
		req.Header.Set("Access-Control-Request-Method", "POST")
		rec := httptest.NewRecorder()
		
		handler.ServeHTTP(rec, req)
		
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
	
	t.Run("allows connect-rpc headers", func(t *testing.T) {
		t.Parallel()
		
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set("Origin", "https://example.com")
		req.Header.Set("Connect-Protocol-Version", "1")
		rec := httptest.NewRecorder()
		
		handler.ServeHTTP(rec, req)
		
		assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, http.StatusOK, rec.Code)
	})
	
	t.Run("sets max age header", func(t *testing.T) {
		t.Parallel()
		
		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		req.Header.Set("Origin", "https://example.com")
		req.Header.Set("Access-Control-Request-Method", "POST")
		rec := httptest.NewRecorder()
		
		handler.ServeHTTP(rec, req)
		
		assert.Equal(t, "7200", rec.Header().Get("Access-Control-Max-Age"))
	})
}
