package connection

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

type ConnectionStorer interface {
	GetConnectionIDs(ctx context.Context) ([]string, error)
	AddConnectionID(ctx context.Context, connectionID string) error
	MarkConnectionIDDisconnected(ctx context.Context, connectionID string) error
}

type connectionStorerStruct struct {
	ConnectionStorer
}

func NewConnection() ConnectionStorer {
	var connection = connectionStorerStruct{}
	return &connection
}

func (con *connectionStorerStruct) GetConnectionIDs(ctx context.Context) ([]string, error) {
	return []string{}, nil
}
func (con *connectionStorerStruct) AddConnectionID(ctx context.Context, connectionID string) error {
	return nil
}
func (con *connectionStorerStruct) MarkConnectionIDDisconnected(ctx context.Context, connectionID string) error {
	return nil
}

// holds the api gateway for the entire lifespan of the lambda function
var apigateway *apigatewaymanagementapi.ApiGatewayManagementApi

func Echo(ctx context.Context, event events.APIGatewayWebsocketProxyRequest, store ConnectionStorer) error {
	if apigateway == nil {
		sess, err := session.NewSession()
		if err != nil {
			log.Fatalln("Unable to create AWS session", err.Error())
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

	connections, err := store.GetConnectionIDs(ctx)
	for _, conn := range connections {
		input := &apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: aws.String(conn),
			Data:         []byte(resp),
		}
		_, err = apigateway.PostToConnection(input)
	}
	return nil
}
