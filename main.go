package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/sawaryu/micro-socket/connection"
)

func handleRequest(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionStore := connection.NewStore()

	rc := event.RequestContext
	switch rk := rc.RouteKey; rk {
	case "$connect":
		// save connection id
		err := connectionStore.AddConnectionID(ctx, rc.ConnectionID)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
	case "$disconnect":
		// delete connection id
		err := connectionStore.MarkConnectionIDDisconnected(ctx, rc.ConnectionID)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
	case "$default":
		// get all current connection ids
		// manage every message sent by the clients
		log.Println("Default", rc.ConnectionID)
		err := handleDefault(ctx, event, connectionStore)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
	default:
		log.Fatalf("Unknown RouteKey %v", rk)
	}
	// API Gateway is expecting an "everything is ok" answer unless something happens
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}

// holds the api gateway for the entire lifespan of the lambda function
var apigateway *apigatewaymanagementapi.ApiGatewayManagementApi

func handleDefault(ctx context.Context, event events.APIGatewayWebsocketProxyRequest, store connection.ConnectionStorer) error {
	if apigateway == nil {
		sess, err := session.NewSession()
		if err != nil {
			log.Fatalf("unable to create aws session: %s", err.Error())
		}
		dname := event.RequestContext.DomainName
		stage := event.RequestContext.Stage
		endpoint := fmt.Sprintf("https://%v/%v", dname, stage)
		apigateway = apigatewaymanagementapi.New(sess, aws.NewConfig().WithEndpoint(endpoint))
	}

	body := event.Body
	resp := fmt.Sprintf("Echo me: %v", body)
	// if the body contains an integer, than a delay in the response is introduced
	delay, err := strconv.Atoi(body)
	if err != nil {
		delay = 0
	}
	time.Sleep(time.Duration(delay) * time.Second)

	connections, _ := store.GetConnectionIDs(ctx)
	for _, conn := range connections {
		input := &apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: aws.String(conn),
			Data:         []byte(resp),
		}
		_, _ = apigateway.PostToConnection(input)
	}
	return nil
}
