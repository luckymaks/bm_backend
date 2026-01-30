package modelconv_test

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/require"

	"github.com/luckymaks/bm_backend/backend/internal/model"
	"github.com/luckymaks/bm_backend/backend/internal/model/modelconv"
	bmv1 "github.com/luckymaks/bm_backend/backend/proto/bm/v1"
)

func TestFromProto(t *testing.T) {
	t.Parallel()

	t.Run("converts proto request to model input", func(t *testing.T) {
		t.Parallel()
		id := "test-123"
		data := "hello world"
		req := bmv1.CreateItemRequest_builder{
			Id:   &id,
			Data: &data,
		}.Build()

		input, err := modelconv.FromProto[model.CreateItemInput](req)

		require.NoError(t, err)
		require.Equal(t, "test-123", input.ID)
		require.Equal(t, "hello world", input.Data)
	})

	t.Run("handles empty fields", func(t *testing.T) {
		t.Parallel()
		req := &bmv1.GetItemRequest{}

		input, err := modelconv.FromProto[model.GetItemInput](req)

		require.NoError(t, err)
		require.Empty(t, input.ID)
	})
}

func TestToProto(t *testing.T) {
	t.Parallel()

	t.Run("converts model output to proto response", func(t *testing.T) {
		t.Parallel()
		output := &model.GetItemOutput{
			ID:        "test-456",
			Data:      "stored data",
			CreatedAt: "2024-01-15T10:00:00Z",
		}

		resp, err := modelconv.ToProto(output, &bmv1.GetItemResponse{})

		require.NoError(t, err)
		require.Equal(t, "test-456", resp.GetId())
		require.Equal(t, "stored data", resp.GetData())
		require.Equal(t, "2024-01-15T10:00:00Z", resp.GetCreatedAt())
	})

	t.Run("converts create item output", func(t *testing.T) {
		t.Parallel()
		output := &model.CreateItemOutput{
			ID: "new-item",
		}

		resp, err := modelconv.ToProto(output, &bmv1.CreateItemResponse{})

		require.NoError(t, err)
		require.Equal(t, "new-item", resp.GetId())
	})
}

func TestToConnectError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		err          error
		expectedCode connect.Code
	}{
		{
			name:         "maps validation error to invalid argument",
			err:          model.ErrValidation,
			expectedCode: connect.CodeInvalidArgument,
		},
		{
			name:         "maps wrapped validation error to invalid argument",
			err:          errors.Wrap(model.ErrValidation, "id cannot be empty"),
			expectedCode: connect.CodeInvalidArgument,
		},
		{
			name:         "maps not found error",
			err:          model.ErrNotFound,
			expectedCode: connect.CodeNotFound,
		},
		{
			name:         "maps already exists error",
			err:          model.ErrAlreadyExists,
			expectedCode: connect.CodeAlreadyExists,
		},
		{
			name:         "maps unknown error to internal",
			err:          errors.New("database failure"),
			expectedCode: connect.CodeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := modelconv.ToConnectError(tt.err)

			var connectErr *connect.Error
			require.ErrorAs(t, result, &connectErr)
			require.Equal(t, tt.expectedCode, connectErr.Code())
		})
	}
}
