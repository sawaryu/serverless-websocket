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
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type ConnectionStorer interface {
	GetConnectionIDs(ctx context.Context) ([]string, error)
	AddConnectionID(ctx context.Context, connectionID string) error
	MarkConnectionIDDisconnected(ctx context.Context, connectionID string) error
}

type connectionStorerStruct struct {
	ConnectionStorer
	DDB *dynamodb.DynamoDB
}

type Thread struct {
	// 'dynamodbav' must be assign attribute name
	ConnectionID string `dynamodbav:"connectionId" json:"connectionId"`
}

const (
	RegionName    = "ap-northeast-1"
	TableName     = "ConnectionsTable"
	AttributeName = "connectionId"
)

// create instance having dynamoinstance and connection interface
func NewConnection() ConnectionStorer {
	new_session, err := session.NewSession(aws.NewConfig().WithRegion("ap-northeast-1"))
	if err != nil {
		log.Fatalln("cannot connect to dynamodb", err.Error())
	}
	ddb := dynamodb.New(new_session)

	connection := connectionStorerStruct{
		DDB: ddb,
	}
	return &connection
}

// select
func (con *connectionStorerStruct) GetConnectionIDs(ctx context.Context) ([]string, error) {
	// scan all connections (* 1MB: max scan size )
	var threads []Thread = []Thread{}
	scanOut, err := con.DDB.Scan(&dynamodb.ScanInput{
		TableName: aws.String(TableName),
	})

	if err != nil {
		log.Fatalln("cannot call scan output", err.Error())
	}

	// unmarshal
	for _, scanedThread := range scanOut.Items {
		var threadTmp Thread
		_ = dynamodbattribute.UnmarshalMap(scanedThread, &threadTmp)
		threads = append(threads, threadTmp)
	}

	// to []string
	var connectionIDs []string
	for _, thread := range threads {
		connectionIDs = append(connectionIDs, thread.ConnectionID)
	}

	return connectionIDs, nil

}

// insert
func (con *connectionStorerStruct) AddConnectionID(ctx context.Context, connectionID string) error {
	param := &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item: map[string]*dynamodb.AttributeValue{
			AttributeName: {
				S: aws.String(connectionID),
			},
		},
	}
	_, err := con.DDB.PutItem(param)
	if err != nil {
		log.Fatalln("cannot input the connection id", err.Error())
	}

	return nil
}

// delete
func (con *connectionStorerStruct) MarkConnectionIDDisconnected(ctx context.Context, connectionID string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			AttributeName: {
				N: aws.String(connectionID),
			},
		},
		TableName: aws.String(TableName),
	}

	_, err := con.DDB.DeleteItem(input)
	if err != nil {
		log.Fatalln("cannot calling DeleteItem", err)
	}

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
