package model_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/cockroachdb/errors"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
	tcdynamodb "github.com/testcontainers/testcontainers-go/modules/dynamodb"

	"github.com/luckymaks/bm_backend/backend/internal/model"
)

var (
	sharedClient   *dynamodb.Client
	tableCounter   atomic.Uint64
	sharedCtx      context.Context
	sharedEndpoint string
)

var errContainerInit error

func TestMain(m *testing.M) {
	cleanup, err := safeInitSharedContainer()
	if err != nil {
		errContainerInit = err
	} else if cleanup != nil {
		defer cleanup()
	}
	m.Run()
}

func safeInitSharedContainer() (cleanup func(), err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Newf("Docker container init panic: %v", r)
		}
	}()
	return initSharedContainer()
}

func initSharedContainer() (cleanup func(), err error) {
	sharedCtx = context.Background()

	container, err := tcdynamodb.Run(sharedCtx, "amazon/dynamodb-local:2.2.1")
	if err != nil {
		return nil, errors.Wrap(err, "failed to start dynamodb container")
	}

	host, err := container.Host(sharedCtx)
	if err != nil {
		_ = container.Terminate(sharedCtx)
		return nil, errors.Wrap(err, "failed to get container host")
	}

	port, err := container.MappedPort(sharedCtx, nat.Port("8000/tcp"))
	if err != nil {
		_ = container.Terminate(sharedCtx)
		return nil, errors.Wrap(err, "failed to get container port")
	}

	sharedEndpoint = "http://" + host + ":" + port.Port()

	cfg, err := config.LoadDefaultConfig(sharedCtx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     "test",
				SecretAccessKey: "test",
			},
		}),
	)
	if err != nil {
		_ = container.Terminate(sharedCtx)
		return nil, errors.Wrap(err, "failed to load aws config")
	}

	sharedClient = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(sharedEndpoint)
	})

	return func() {
		_ = container.Terminate(sharedCtx)
	}, nil
}

func setup(t *testing.T) (context.Context, *model.Model) {
	t.Helper()

	if errContainerInit != nil {
		t.Skipf("skipping: Docker container not available: %v", errContainerInit)
	}
	if sharedClient == nil {
		t.Fatal("shared container not initialized - ensure TestMain calls initSharedContainer")
	}

	tableNum := tableCounter.Add(1)
	tableName := fmt.Sprintf("test-table-%d", tableNum)

	_, err := sharedClient.CreateTable(sharedCtx, &dynamodb.CreateTableInput{
		TableName:   aws.String(tableName),
		BillingMode: types.BillingModePayPerRequest,
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("pk"), KeyType: types.KeyTypeHash},
			{AttributeName: aws.String("sk"), KeyType: types.KeyTypeRange},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("pk"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("sk"), AttributeType: types.ScalarAttributeTypeS},
		},
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_, _ = sharedClient.DeleteTable(sharedCtx, &dynamodb.DeleteTableInput{
			TableName: aws.String(tableName),
		})
	})

	mdl := model.New(sharedClient, tableName)

	return sharedCtx, mdl
}

func setupWithGSI(t *testing.T) (context.Context, *model.Model) {
	t.Helper()

	if errContainerInit != nil {
		t.Skipf("skipping: Docker container not available: %v", errContainerInit)
	}
	if sharedClient == nil {
		t.Fatal("shared container not initialized - ensure TestMain calls initSharedContainer")
	}

	tableNum := tableCounter.Add(1)
	tableName := fmt.Sprintf("test-table-%d", tableNum)

	_, err := sharedClient.CreateTable(sharedCtx, &dynamodb.CreateTableInput{
		TableName:   aws.String(tableName),
		BillingMode: types.BillingModePayPerRequest,
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("pk"), KeyType: types.KeyTypeHash},
			{AttributeName: aws.String("sk"), KeyType: types.KeyTypeRange},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("pk"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("sk"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("gsi1pk"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("gsi1sk"), AttributeType: types.ScalarAttributeTypeS},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("gsi1"),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("gsi1pk"), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String("gsi1sk"), KeyType: types.KeyTypeRange},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
			},
		},
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_, _ = sharedClient.DeleteTable(sharedCtx, &dynamodb.DeleteTableInput{
			TableName: aws.String(tableName),
		})
	})

	mdl := model.New(sharedClient, tableName)

	return sharedCtx, mdl
}
