package model

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/smithy-go/ptr"
	"github.com/cockroachdb/errors"
)

type CreateItemInput struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

type CreateItemOutput struct {
	ID string `json:"id"`
}

type GetItemInput struct {
	ID string `json:"id"`
}

type GetItemOutput struct {
	ID        string `json:"id"`
	Data      string `json:"data"`
	CreatedAt string `json:"created_at"`
}

type itemRecord struct {
	PK        string `dynamodbav:"pk"`
	SK        string `dynamodbav:"sk"`
	Data      string `dynamodbav:"data"`
	CreatedAt string `dynamodbav:"createdAt"`
}

func (m *Model) CreateItem(ctx context.Context, input CreateItemInput) (*CreateItemOutput, error) {
	if input.ID == "" {
		return nil, errors.Wrap(ErrValidation, "id cannot be empty")
	}

	item := itemRecord{
		PK:        fmt.Sprintf("ITEM#%s", input.ID),
		SK:        fmt.Sprintf("ITEM#%s", input.ID),
		Data:      input.Data,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal item")
	}

	_, err = m.dynamo.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: ptr.String(m.tableName),
		Item:      av,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to store item")
	}

	return &CreateItemOutput{
		ID: input.ID,
	}, nil
}

func (m *Model) GetItem(ctx context.Context, input GetItemInput) (*GetItemOutput, error) {
	if input.ID == "" {
		return nil, errors.Wrap(ErrValidation, "id cannot be empty")
	}

	key, err := attributevalue.MarshalMap(map[string]string{
		"pk": fmt.Sprintf("ITEM#%s", input.ID),
		"sk": fmt.Sprintf("ITEM#%s", input.ID),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal key")
	}

	result, err := m.dynamo.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: ptr.String(m.tableName),
		Key:       key,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get item")
	}

	if result.Item == nil {
		return nil, ErrNotFound
	}

	var record itemRecord
	if err := attributevalue.UnmarshalMap(result.Item, &record); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal item")
	}

	return &GetItemOutput{
		ID:        input.ID,
		Data:      record.Data,
		CreatedAt: record.CreatedAt,
	}, nil
}
