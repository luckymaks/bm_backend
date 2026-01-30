package rpc

import (
	"context"

	"github.com/cockroachdb/errors"

	"github.com/luckymaks/bm_backend/backend/internal/model"
	"github.com/luckymaks/bm_backend/backend/internal/model/modelconv"
	bmv1 "github.com/luckymaks/bm_backend/backend/proto/bm/v1"
)

func (r *RPC) CreateItem(
	ctx context.Context, req *bmv1.CreateItemRequest,
) (*bmv1.CreateItemResponse, error) {
	input, err := modelconv.FromProto[model.CreateItemInput](req)
	if err != nil {
		return nil, modelconv.ToConnectError(err)
	}

	output, err := r.model.CreateItem(ctx, input)
	if err != nil {
		return nil, modelconv.ToConnectError(err)
	}

	resp, err := modelconv.ToProto(output, &bmv1.CreateItemResponse{})
	if err != nil {
		return nil, errors.Wrap(err, "convert response")
	}

	return resp, nil
}

func (r *RPC) GetItem(
	ctx context.Context, req *bmv1.GetItemRequest,
) (*bmv1.GetItemResponse, error) {
	input, err := modelconv.FromProto[model.GetItemInput](req)
	if err != nil {
		return nil, modelconv.ToConnectError(err)
	}

	output, err := r.model.GetItem(ctx, input)
	if err != nil {
		return nil, modelconv.ToConnectError(err)
	}

	resp, err := modelconv.ToProto(output, &bmv1.GetItemResponse{})
	if err != nil {
		return nil, errors.Wrap(err, "convert response")
	}

	return resp, nil
}
