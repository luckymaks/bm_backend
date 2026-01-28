package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/luckymaks/bm_backend/backend/internal/rpc"
)

var dynamoClient *dynamodb.Client
var tableName string

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("unable to load AWS config: " + err.Error())
	}
	dynamoClient = dynamodb.NewFromConfig(cfg)
	tableName = os.Getenv("MAIN_TABLE_NAME")
	
	e := echo.New()
	e.Use(middleware.Recover())
	
	e.GET("/", handleHello)
	e.GET("/health", handleHealth)
	e.POST("/items", handleCreateItem)
	e.GET("/items/:id", handleGetItem)

	server := &http.Server{
		Addr:    ":12001",
		Handler: rpc.WithCORS(e),
	}
	//nolint:errcheck
	server.ListenAndServe()
}

func handleHello(c echo.Context) error {
	return c.String(http.StatusOK, "hello, "+c.RealIP())
}

func handleHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
		"region": os.Getenv("AWS_REGION"),
	})
}

type CreateItemRequest struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

func handleCreateItem(c echo.Context) error {
	var req CreateItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	
	_, err := dynamoClient.PutItem(c.Request().Context(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"pk":        &types.AttributeValueMemberS{Value: "ITEM#" + req.ID},
			"sk":        &types.AttributeValueMemberS{Value: "ITEM#" + req.ID},
			"data":      &types.AttributeValueMemberS{Value: req.Data},
			"createdAt": &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)},
		},
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	
	return c.JSON(http.StatusCreated, map[string]string{"id": req.ID})
}

func handleGetItem(c echo.Context) error {
	id := c.Param("id")
	
	result, err := dynamoClient.GetItem(c.Request().Context(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "ITEM#" + id},
			"sk": &types.AttributeValueMemberS{Value: "ITEM#" + id},
		},
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	
	if result.Item == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "item not found"})
	}
	
	item := map[string]string{
		"id": id,
	}
	if v, ok := result.Item["data"].(*types.AttributeValueMemberS); ok {
		item["data"] = v.Value
	}
	if v, ok := result.Item["createdAt"].(*types.AttributeValueMemberS); ok {
		item["createdAt"] = v.Value
	}
	
	return c.JSON(http.StatusOK, item)
}
