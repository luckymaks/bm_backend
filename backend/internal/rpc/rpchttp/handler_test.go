package rpchttp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/luckymaks/bm_backend/backend/internal/rpc/rpchttp"
	bmv1 "github.com/luckymaks/bm_backend/backend/proto/bm/v1"
	"github.com/luckymaks/bm_backend/backend/proto/bm/v1/bmv1connect"
)

type mockApiService struct {
	bmv1connect.UnimplementedApiServiceHandler
	getHealthResp *bmv1.GetHealthResponse
	getHealthErr  error
}

func (m *mockApiService) GetHealth(ctx context.Context, req *bmv1.GetHealthRequest) (*bmv1.GetHealthResponse, error) {
	if m.getHealthErr != nil {
		return nil, m.getHealthErr
	}
	return m.getHealthResp, nil
}

func TestNewHandler(t *testing.T) {
	t.Parallel()

	t.Run("creates handler that responds to RPC requests", func(t *testing.T) {
		t.Parallel()
		mockSvc := &mockApiService{
			getHealthResp: &bmv1.GetHealthResponse{},
		}
		mockSvc.getHealthResp.SetStatus("ok")
		mockSvc.getHealthResp.SetRegion("us-east-1")

		handler := rpchttp.NewHandler(mockSvc)

		req := httptest.NewRequest(
			http.MethodPost,
			"/v1/rpc/bm.v1.ApiService/GetHealth",
			strings.NewReader("{}"),
		)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "ok")
	})

	t.Run("returns 404 for unknown paths", func(t *testing.T) {
		t.Parallel()
		mockSvc := &mockApiService{}

		handler := rpchttp.NewHandler(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/unknown/path", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestWithCORS(t *testing.T) {
	t.Parallel()

	inner := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := rpchttp.WithCORS(inner)

	t.Run("request with origin includes CORS headers", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/v1/rpc/test", nil)
		req.Header.Set("Origin", "https://example.com")

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("preflight request is handled", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodOptions, "/v1/rpc/test", nil)
		req.Header.Set("Origin", "https://example.com")
		req.Header.Set("Access-Control-Request-Method", "POST")

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNoContent, rec.Code)
	})
}
