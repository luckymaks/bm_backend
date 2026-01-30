package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/luckymaks/bm_backend/backend/internal/model"
	"github.com/luckymaks/bm_backend/backend/internal/rpc"
	"github.com/luckymaks/bm_backend/backend/internal/rpc/rpcconfig"
	"github.com/luckymaks/bm_backend/backend/internal/rpc/rpchttp"
)

func main() {
	ctx := context.Background()

	cfgLoader, err := rpcconfig.NewLoader()
	if err != nil {
		panic("unable to load config: " + err.Error())
	}

	cfg, err := cfgLoader.Load(ctx)
	if err != nil {
		panic("unable to load config: " + err.Error())
	}

	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("unable to load AWS config: " + err.Error())
	}

	if cfg.Env.MainTableName == "" {
		slog.Error("DYNAMO_TABLE_NAME environment variable is required")
		os.Exit(1)
	}

	dynamoClient := dynamodb.NewFromConfig(awsCfg)

	mdl := model.New(dynamoClient, cfg.Env.MainTableName)
	rpcHandler := rpc.NewRPC(rpc.Config{
		DeploymentIdent: cfg.Env.DeploymentIdent,
		AWSRegion:       cfg.Env.AWSRegion,
	}, mdl)
	handler := rpchttp.NewHandler(rpcHandler)

	slog.Info("starting httpapi server", "port", "12001")

	//nolint:errcheck,gosec
	http.ListenAndServe(":12001", handler)
}
