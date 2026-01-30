package rpchttp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/require"

	"github.com/luckymaks/bm_backend/backend/internal/model"
	"github.com/luckymaks/bm_backend/backend/internal/rpc/rpchttp"
	bmv1 "github.com/luckymaks/bm_backend/backend/proto/bm/v1"
	"github.com/luckymaks/bm_backend/backend/proto/bm/v1/bmv1connect"
)

var _ bmv1connect.ApiServiceHandler = (*stubApiService)(nil)

type stubApiService struct {
	bmv1connect.UnimplementedApiServiceHandler
	err error
}

func (s *stubApiService) GetHealth(
	_ context.Context,
	_ *bmv1.GetHealthRequest,
) (*bmv1.GetHealthResponse, error) {
	if s.err != nil {
		return nil, s.err
	}

	resp := &bmv1.GetHealthResponse{}
	resp.SetStatus("ok")
	return resp, nil
}

func TestErrorInterceptor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		err          error
		expectedCode connect.Code
	}{
		{
			name:         "maps ErrNotFound to CodeNotFound",
			err:          errors.Wrap(model.ErrNotFound, "item not found"),
			expectedCode: connect.CodeNotFound,
		},
		{
			name:         "maps ErrAlreadyExists to CodeAlreadyExists",
			err:          errors.Wrap(model.ErrAlreadyExists, "item exists"),
			expectedCode: connect.CodeAlreadyExists,
		},
		{
			name:         "maps ErrValidation to CodeInvalidArgument",
			err:          errors.Wrap(model.ErrValidation, "id empty"),
			expectedCode: connect.CodeInvalidArgument,
		},
		{
			name:         "preserves existing ConnectRPC error",
			err:          connect.NewError(connect.CodePermissionDenied, errors.New("forbidden")),
			expectedCode: connect.CodePermissionDenied,
		},
		{
			name:         "unknown error passes through as-is",
			err:          errors.New("something went wrong"),
			expectedCode: connect.CodeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stub := &stubApiService{err: tt.err}
			_, handler := bmv1connect.NewApiServiceHandler(
				stub,
				connect.WithInterceptors(rpchttp.NewErrorInterceptor()),
			)

			srv := httptest.NewServer(handler)
			t.Cleanup(srv.Close)

			client := bmv1connect.NewApiServiceClient(http.DefaultClient, srv.URL)
			_, err := client.GetHealth(context.Background(), &bmv1.GetHealthRequest{})

			require.Error(t, err)
			require.Equal(t, tt.expectedCode, connect.CodeOf(err))
		})
	}

	t.Run("no error passes through", func(t *testing.T) {
		t.Parallel()

		stub := &stubApiService{err: nil}
		_, handler := bmv1connect.NewApiServiceHandler(
			stub,
			connect.WithInterceptors(rpchttp.NewErrorInterceptor()),
		)

		srv := httptest.NewServer(handler)
		t.Cleanup(srv.Close)

		client := bmv1connect.NewApiServiceClient(http.DefaultClient, srv.URL)
		_, err := client.GetHealth(context.Background(), &bmv1.GetHealthRequest{})

		require.NoError(t, err)
	})
}
