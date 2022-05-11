package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sawaryu/micro-socket/connection"
)

var connectionStore = connection.NewConnection()

func handleRequest(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	rc := event.RequestContext
	switch rk := rc.RouteKey; rk {
	case "$connect":
		// 接続時にconnectionIDを保存する
		err := connectionStore.AddConnectionID(ctx, rc.ConnectionID)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
	case "$disconnect":
		// 切断時にconnectionIDをリボークする
		err := connectionStore.MarkConnectionIDDisconnected(ctx, rc.ConnectionID)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
	case "$default":
		// manage every message sent by the clients
		log.Println("Default", rc.ConnectionID)
		err := connection.Echo(ctx, event, connectionStore)
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
