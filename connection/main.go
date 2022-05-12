package connection

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type ConnectionStorer interface {
	GetConnectionIDs(ctx context.Context) ([]string, error)
	AddConnectionID(ctx context.Context, connectionID string) error
	MarkConnectionIDDisconnected(ctx context.Context, connectionID string) error
}

type ConnectionStore struct {
	ConnectionStorer
	DDB *dynamodb.DynamoDB
}

type Thread struct {
	// 'dynamodbav' must be assign attribute name
	ConnectionID string `dynamodbav:"connectionId" json:"connectionId"`
}

const (
	RegionName = "ap-northeast-1"
	TableName  = "connectionsTable"
)

// create instance having dynamoinstance and connection interface
func NewStore() ConnectionStorer {
	new_session, err := session.NewSession(aws.NewConfig().WithRegion(RegionName))
	if err != nil {
		log.Fatalln("cannot connect to dynamodb", err.Error())
	}
	ddb := dynamodb.New(new_session)

	connection := ConnectionStore{
		DDB: ddb,
	}
	return &connection
}

// select
func (store *ConnectionStore) GetConnectionIDs(ctx context.Context) ([]string, error) {
	// scan all connections (* 1MB: max scan size )
	var threads []Thread = []Thread{}
	scanOut, err := store.DDB.Scan(&dynamodb.ScanInput{
		TableName: aws.String(TableName),
	})

	if err != nil {
		log.Fatalf("cannot calling scan items: %s", err.Error())
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
func (store *ConnectionStore) AddConnectionID(ctx context.Context, connectionID string) error {
	param := &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item: map[string]*dynamodb.AttributeValue{
			"connectionId": {
				S: aws.String(connectionID),
			},
		},
	}
	_, err := store.DDB.PutItem(param)
	if err != nil {
		log.Fatalf("cannot calling input item: %s", err.Error())
	}

	return nil
}

// delete
func (store *ConnectionStore) MarkConnectionIDDisconnected(ctx context.Context, connectionID string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"connectionId": {
				S: aws.String(connectionID),
			},
		},
		TableName: aws.String(TableName),
	}

	_, err := store.DDB.DeleteItem(input)
	if err != nil {
		log.Fatalf("cannot calling delete item: %s", err.Error())
	}

	return nil
}
