package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/sawaryu/serverless-websocket/connection"
)

func handleRequest(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionStore := connection.NewStore()
	rc := event.RequestContext

	switch rk := rc.RouteKey; rk {
	case "$connect":
		err := connectionStore.AddConnectionID(ctx, rc.ConnectionID)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
	case "$disconnect":
		err := connectionStore.MarkConnectionDisconnected(ctx, rc.ConnectionID)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
	case "$default":
		err := handleDefault(ctx, event, connectionStore)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
	default:
		log.Fatalf("unknown route key: %v", rk)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}

// holds the api gateway for the entire lifespan of the lambda function
var apigateway *apigatewaymanagementapi.ApiGatewayManagementApi

func handleDefault(ctx context.Context, event events.APIGatewayWebsocketProxyRequest, store connection.IConnectionStore) error {
	if apigateway == nil {
		awsSession, err := session.NewSession()
		if err != nil {
			log.Fatalf("couldn't create new aws session: %s", err.Error())
		}
		domainName := event.RequestContext.DomainName
		stage := event.RequestContext.Stage
		endpoint := fmt.Sprintf("https://%s/%s", domainName, stage)
		apigateway = apigatewaymanagementapi.New(awsSession, aws.NewConfig().WithEndpoint(endpoint))
	}

	body := event.Body
	response := fmt.Sprintf("Echo me: %v", body)

	connectionIDs := store.FetchConnectionIDs(ctx)
	for _, connID := range connectionIDs {
		input := &apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: aws.String(connID),
			Data:         []byte(response),
		}
		_, _ = apigateway.PostToConnection(input)
	}
	return nil
}
