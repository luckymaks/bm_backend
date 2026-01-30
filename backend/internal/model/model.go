package model

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBClient interface {
	PutItem(
		ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options),
	) (*dynamodb.PutItemOutput, error)
	GetItem(
		ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options),
	) (*dynamodb.GetItemOutput, error)
	UpdateItem(
		ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options),
	) (*dynamodb.UpdateItemOutput, error)
	DeleteItem(
		ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options),
	) (*dynamodb.DeleteItemOutput, error)
	Scan(
		ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options),
	) (*dynamodb.ScanOutput, error)
	Query(
		ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options),
	) (*dynamodb.QueryOutput, error)
}

type Model struct {
	dynamo    DynamoDBClient
	tableName string
}

type Option func(*Model)

func New(dynamo DynamoDBClient, tableName string, opts ...Option) *Model {
	m := &Model{
		dynamo:    dynamo,
		tableName: tableName,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}
