package rpc_test

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/luckymaks/bm_backend/backend/internal/model"
	bmv1 "github.com/luckymaks/bm_backend/backend/proto/bm/v1"
)

func TestCreateItem(t *testing.T) {
	t.Parallel()

	t.Run("successfully maps model response to proto", func(t *testing.T) {
		t.Parallel()
		id := "test-123"
		data := "hello world"

		ctx, r, mockModel := setup(t)

		mockModel.EXPECT().CreateItem(mock.Anything, model.CreateItemInput{
			ID:   id,
			Data: data,
		}).Return(&model.CreateItemOutput{
			ID: id,
		}, nil)

		resp, err := r.CreateItem(ctx, bmv1.CreateItemRequest_builder{
			Id:   &id,
			Data: &data,
		}.Build())

		require.NoError(t, err)
		require.Equal(t, id, resp.GetId())
		mockModel.AssertExpectations(t)
	})

	t.Run("maps model validation errors to invalid argument", func(t *testing.T) {
		t.Parallel()
		ctx, r, mockModel := setup(t)
		emptyID := ""
		data := "test"

		modelErr := errors.Wrap(model.ErrValidation, "id cannot be empty")
		mockModel.EXPECT().CreateItem(mock.Anything, model.CreateItemInput{
			ID:   emptyID,
			Data: data,
		}).Return(nil, modelErr)

		resp, err := r.CreateItem(ctx, bmv1.CreateItemRequest_builder{
			Id:   &emptyID,
			Data: &data,
		}.Build())

		require.Nil(t, resp)
		require.Error(t, err)
		var connectErr *connect.Error
		require.ErrorAs(t, err, &connectErr)
		require.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
		mockModel.AssertExpectations(t)
	})
}

func TestGetItem(t *testing.T) {
	t.Parallel()

	t.Run("successfully returns item", func(t *testing.T) {
		t.Parallel()
		id := "test-456"

		ctx, r, mockModel := setup(t)

		mockModel.EXPECT().GetItem(mock.Anything, model.GetItemInput{
			ID: id,
		}).Return(&model.GetItemOutput{
			ID:        id,
			Data:      "stored data",
			CreatedAt: "2024-01-15T10:00:00Z",
		}, nil)

		resp, err := r.GetItem(ctx, bmv1.GetItemRequest_builder{
			Id: &id,
		}.Build())

		require.NoError(t, err)
		require.Equal(t, id, resp.GetId())
		require.Equal(t, "stored data", resp.GetData())
		require.Equal(t, "2024-01-15T10:00:00Z", resp.GetCreatedAt())
		mockModel.AssertExpectations(t)
	})

	t.Run("maps not found errors", func(t *testing.T) {
		t.Parallel()
		ctx, r, mockModel := setup(t)
		id := "nonexistent"

		mockModel.EXPECT().GetItem(mock.Anything, model.GetItemInput{
			ID: id,
		}).Return(nil, model.ErrNotFound)

		resp, err := r.GetItem(ctx, bmv1.GetItemRequest_builder{
			Id: &id,
		}.Build())

		require.Nil(t, resp)
		require.Error(t, err)
		var connectErr *connect.Error
		require.ErrorAs(t, err, &connectErr)
		require.Equal(t, connect.CodeNotFound, connectErr.Code())
		mockModel.AssertExpectations(t)
	})
}
