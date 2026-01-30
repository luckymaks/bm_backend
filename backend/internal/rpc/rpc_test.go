package rpc_test

import (
	"context"
	"testing"

	"github.com/luckymaks/bm_backend/backend/internal/rpc"
	"github.com/luckymaks/bm_backend/backend/internal/rpc/mocks"
)

type testFixtures struct {
	ctx       context.Context
	rpc       *rpc.RPC
	mockModel *mocks.MockModelClient
}

func setup(t *testing.T) (context.Context, *rpc.RPC, *mocks.MockModelClient) {
	t.Helper()
	f := setupWithConfig(t, rpc.Config{})
	return f.ctx, f.rpc, f.mockModel
}

func setupWithConfig(t *testing.T, cfg rpc.Config) *testFixtures {
	t.Helper()

	mockModel := mocks.NewMockModelClient(t)

	r := rpc.NewRPC(cfg, mockModel)

	return &testFixtures{
		ctx:       t.Context(),
		rpc:       r,
		mockModel: mockModel,
	}
}
