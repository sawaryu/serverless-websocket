package session

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type ISessionStore interface {
	FetchSessionIDs(ctx context.Context) []string
	AddSessionID(ctx context.Context, connectionID string) error
	MarkSessionDisconnected(ctx context.Context, connectionID string) error
}

type SessionStore struct {
	DDB *dynamodb.DynamoDB
}

type Session struct {
	SessionID string `dynamodbav:"session_id" json:"session_id"`
}

const (
	RegionName = "ap-northeast-1"
	TableName  = "Session"
)

func NewStore() ISessionStore {
	new_session, err := session.NewSession(aws.NewConfig().WithRegion(RegionName))
	if err != nil {
		log.Fatalln("couldn't connect to dynamodb", err.Error())
	}
	ddb := dynamodb.New(new_session)

	connection := SessionStore{
		DDB: ddb,
	}
	return &connection
}

// scan all connections (* 1MB: max scan size )
func (store *SessionStore) FetchSessionIDs(ctx context.Context) (sessionIDs []string) {
	sessions := []Session{}
	scanOut, err := store.DDB.Scan(&dynamodb.ScanInput{
		TableName: aws.String(TableName),
	})

	if err != nil {
		log.Fatalf("couldn't scan items: %s", err.Error())
	}

	for _, scanedSession := range scanOut.Items {
		var sessionTmp Session
		_ = dynamodbattribute.UnmarshalMap(scanedSession, &sessionTmp)
		sessions = append(sessions, sessionTmp)
	}

	for _, session := range sessions {
		sessionIDs = append(sessionIDs, session.SessionID)
	}

	return sessionIDs

}

// put new connected session
func (store *SessionStore) AddSessionID(ctx context.Context, sessionID string) error {
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item: map[string]*dynamodb.AttributeValue{
			"session_id": {
				S: aws.String(sessionID),
			},
		},
	}

	_, err := store.DDB.PutItem(putInput)
	if err != nil {
		log.Fatalf("couldn't calling input item: %s", err.Error())
	}

	return nil
}

// delete disconnected session
func (store *SessionStore) MarkSessionDisconnected(ctx context.Context, sessionID string) error {
	deleteInput := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"session_id": {
				S: aws.String(sessionID),
			},
		},
		TableName: aws.String(TableName),
	}

	_, err := store.DDB.DeleteItem(deleteInput)
	if err != nil {
		log.Fatalf("couldn't calling delete item: %s", err.Error())
	}

	return nil
}
