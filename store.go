package dynamostore

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fsm/fsm"
)

type DynamoStore struct {
	DynamoSession *dynamodb.DynamoDB
	DynamoTable   string
	Network       string
}

func (d *DynamoStore) FetchTraverser(uuid string) (fsm.Traverser, error) {
	// Fetch Item
	result, err := d.DynamoSession.GetItem(
		&dynamodb.GetItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"network": {
					S: aws.String(d.Network),
				},
				"uuid": {
					S: aws.String(uuid),
				},
			},
			TableName: aws.String(d.DynamoTable),
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
		network:       d.Network,
		uuid:          uuid,
		currentState:  *result.Item["currentState"].S,
		dynamoSession: d.DynamoSession,
		dynamoTable:   d.DynamoTable,
		dynamoData:    data,
	}, nil
}

func (d *DynamoStore) CreateTraverser(uuid string) (fsm.Traverser, error) {
	// Create element in Dynamo
	_, err := d.DynamoSession.PutItem(
		&dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"network": {
					S: aws.String(d.Network),
				},
				"uuid": {
					S: aws.String(uuid),
				},
				"data": {
					M: make(map[string]*dynamodb.AttributeValue, 0),
				},
			},
			TableName: aws.String(d.DynamoTable),
		},
	)
	if err != nil {
		return nil, err
	}

	// Create Traverser
	return &dynamoTraverser{
		network:       d.Network,
		uuid:          uuid,
		dynamoSession: d.DynamoSession,
		dynamoTable:   d.DynamoTable,
		dynamoData:    make(map[string]interface{}, 0),
	}, nil
}
