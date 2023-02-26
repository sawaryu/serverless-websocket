package connection

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type IConnectionStore interface {
	FetchConnectionIDs(ctx context.Context) ([]string, error)
	AddConnectionID(ctx context.Context, connectionID string) error
	MarkConnectionDisconnected(ctx context.Context, connectionID string) error
}

type ConnectionStore struct {
	*dynamodb.DynamoDB
}

type Connection struct {
	ConnectionID string `dynamodbav:"connection_id" json:"connection_id"`
}

const (
	RegionName = "ap-northeast-1"
	TableName  = "Connection"
)

func NewStore() IConnectionStore {
	new_connection, err := session.NewSession(aws.NewConfig().WithRegion(RegionName))
	if err != nil {
		log.Fatalln("couldn't connect to dynamodb", err.Error())
	}
	ddb := dynamodb.New(new_connection)

	connection := ConnectionStore{
		ddb,
	}
	return &connection
}

// scan all connections (* 1MB: max scan size )
func (store *ConnectionStore) FetchConnectionIDs(ctx context.Context) (connectionIDs []string, err error) {
	var connections []Connection
	scanOut, err := store.Scan(&dynamodb.ScanInput{
		TableName: aws.String(TableName),
	})

	if err != nil {
		return connectionIDs, err
	}

	for _, itemConnection := range scanOut.Items {
		var connectionTmp Connection
		err = dynamodbattribute.UnmarshalMap(itemConnection, &connectionTmp)
		if err != nil {
			return connectionIDs, err
		}
		connections = append(connections, connectionTmp)
	}

	for _, connection := range connections {
		connectionIDs = append(connectionIDs, connection.ConnectionID)
	}

	return connectionIDs, err

}

// put new connected connection
func (store *ConnectionStore) AddConnectionID(ctx context.Context, connectionID string) error {
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item: map[string]*dynamodb.AttributeValue{
			"connection_id": {
				S: aws.String(connectionID),
			},
		},
	}

	_, err := store.PutItem(putInput)
	if err != nil {
		return err
	}

	return nil
}

// delete disconnected connection
func (store *ConnectionStore) MarkConnectionDisconnected(ctx context.Context, connectionID string) error {
	deleteInput := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"connection_id": {
				S: aws.String(connectionID),
			},
		},
		TableName: aws.String(TableName),
	}

	_, err := store.DeleteItem(deleteInput)
	if err != nil {
		return err
	}

	return nil
}
