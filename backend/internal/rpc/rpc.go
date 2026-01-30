package rpc

import (
	"context"
	"os"

	"github.com/luckymaks/bm_backend/backend/internal/model"
	bmv1 "github.com/luckymaks/bm_backend/backend/proto/bm/v1"
	"github.com/luckymaks/bm_backend/backend/proto/bm/v1/bmv1connect"
)

var _ bmv1connect.ApiServiceHandler = (*RPC)(nil)

type ModelClient interface {
	CreateItem(ctx context.Context, input model.CreateItemInput) (*model.CreateItemOutput, error)
	GetItem(ctx context.Context, input model.GetItemInput) (*model.GetItemOutput, error)
}

type RPC struct {
	bmv1connect.UnimplementedApiServiceHandler
	model           ModelClient
	deploymentIdent string
	awsRegion       string
}

type Config struct {
	DeploymentIdent string
	AWSRegion       string
}

func NewRPC(cfg Config, mdl ModelClient) *RPC {
	return &RPC{
		model:           mdl,
		deploymentIdent: cfg.DeploymentIdent,
		awsRegion:       cfg.AWSRegion,
	}
}

func (r *RPC) GetHealth(
	ctx context.Context,
	req *bmv1.GetHealthRequest,
) (*bmv1.GetHealthResponse, error) {
	region := r.awsRegion
	if region == "" {
		region = os.Getenv("AWS_REGION")
	}

	resp := &bmv1.GetHealthResponse{}
	resp.SetStatus("ok")
	resp.SetRegion(region)
	return resp, nil
}
