package dynamostore

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fsm/fsm"
)

// New returns an instance of a dynamoStore
func New() fsm.Store {
	return &dynamoStore{
		session: getDynamoSession(
			os.Getenv("DYNAMO_REGION"),
			os.Getenv("DYNAMO_ACCESS_KEY_ID"),
			os.Getenv("DYNAMO_SECRET_ACCESS_KEY"),
		),
		tableName: os.Getenv("DYNAMO_TABLE_NAME"),
	}
}

type dynamoStore struct {
	session   *dynamodb.DynamoDB
	tableName string
}

func (d *dynamoStore) FetchTraverser(uuid string) (fsm.Traverser, error) {
	// Fetch Item
	result, err := d.session.GetItem(
		&dynamodb.GetItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"uuid": {
					S: aws.String(uuid),
				},
			},
			TableName: aws.String(d.tableName),
		},
	)
	// Checking for errors with the request
	if err != nil {
		return nil, err
	}
	// Dynamo actually doesn't return an error when the traverser doesn't exist
	// It just returns an empty map.  So we have to check this here to see
	// if the traverser doesn't exist.
	if len(result.Item) == 0 {
		return nil, errors.New("Traverser does not exist")
	}

	// Get Data
	data := make(map[string]interface{}, 0)
	err = dynamodbattribute.ConvertFromMap(result.Item["data"].M, &data)
	if err != nil {
		return nil, err
	}

	// Create Traverser
	return &dynamoTraverser{
		uuid:          uuid,
		currentState:  *result.Item["currentState"].S,
		dynamoSession: d.session,
		dynamoTable:   d.tableName,
		dynamoData:    data,
	}, nil
}

func (d *dynamoStore) CreateTraverser(uuid string) (fsm.Traverser, error) {
	// Create element in Dynamo
	_, err := d.session.PutItem(
		&dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"uuid": {
					S: aws.String(uuid),
				},
				"data": {
					M: make(map[string]*dynamodb.AttributeValue, 0),
				},
			},
			TableName: aws.String(d.tableName),
		},
	)
	if err != nil {
		return nil, err
	}

	// Create Traverser
	return &dynamoTraverser{
		uuid:          uuid,
		dynamoSession: d.session,
		dynamoTable:   d.tableName,
		dynamoData:    make(map[string]interface{}, 0),
	}, nil
}
