package dynamostore

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type dynamoTraverser struct {
	network       string
	uuid          string
	currentState  string
	dynamoSession *dynamodb.DynamoDB
	dynamoTable   string
	dynamoData    map[string]interface{}
}

func (d *dynamoTraverser) UUID() string {
	return d.uuid
}

func (d *dynamoTraverser) SetUUID(newUUID string) {
	// Update UUID on Dynamo
	d.dynamoSession.UpdateItem(
		&dynamodb.UpdateItemInput{
			ExpressionAttributeNames: map[string]*string{
				"#K": aws.String("uuid"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":v": &dynamodb.AttributeValue{
					S: aws.String(newUUID),
				},
			},
			Key: map[string]*dynamodb.AttributeValue{
				"network": {
					S: aws.String(d.network),
				},
				"uuid": {
					S: aws.String(d.uuid),
				},
			},
			TableName:        aws.String(d.dynamoTable),
			UpdateExpression: aws.String("SET #K = :v"),
		},
	)

	// Update local UUID
	d.uuid = newUUID
}

func (d *dynamoTraverser) CurrentState() string {
	return d.currentState
}

func (d *dynamoTraverser) SetCurrentState(state string) {
	// Update state on Dynamo
	d.dynamoSession.UpdateItem(
		&dynamodb.UpdateItemInput{
			ExpressionAttributeNames: map[string]*string{
				"#K": aws.String("currentState"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":v": &dynamodb.AttributeValue{
					S: aws.String(state),
				},
			},
			Key: map[string]*dynamodb.AttributeValue{
				"network": {
					S: aws.String(d.network),
				},
				"uuid": {
					S: aws.String(d.uuid),
				},
			},
			TableName:        aws.String(d.dynamoTable),
			UpdateExpression: aws.String("SET #K = :v"),
		},
	)

	// Update local state
	d.currentState = state
}

func (d *dynamoTraverser) Upsert(key string, value interface{}) error {
	// Convert value to a DynamoAttribute
	item, err := dynamodbattribute.ConvertTo(value)
	if err != nil {
		return err
	}

	// Update attribute in Dynamo
	res, err := d.dynamoSession.UpdateItem(
		&dynamodb.UpdateItemInput{
			ExpressionAttributeNames: map[string]*string{
				"#D": aws.String("data"),
				"#K": aws.String(key),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":v": item,
			},
			Key: map[string]*dynamodb.AttributeValue{
				"network": {
					S: aws.String(d.network),
				},
				"uuid": {
					S: aws.String(d.uuid),
				},
			},
			ReturnValues:     aws.String("ALL_NEW"),
			TableName:        aws.String(d.dynamoTable),
			UpdateExpression: aws.String("SET #D.#K = :v"),
		},
	)

	// Update local
	if err == nil {
		data := make(map[string]interface{}, 0)
		err = dynamodbattribute.ConvertFromMap(res.Attributes["data"].M, &data)
		if err == nil {
			d.dynamoData = data
		}
	}
	return err
}

func (d *dynamoTraverser) Fetch(key string) (interface{}, error) {
	if val, ok := d.dynamoData[key]; ok {
		return val, nil
	}
	return nil, errors.New("Key not set")
}

func (d *dynamoTraverser) Delete(key string) error {
	// Update attribute in Dynamo
	res, err := d.dynamoSession.UpdateItem(
		&dynamodb.UpdateItemInput{
			ExpressionAttributeNames: map[string]*string{
				"#D": aws.String("data"),
				"#K": aws.String(key),
			},
			Key: map[string]*dynamodb.AttributeValue{
				"network": {
					S: aws.String(d.network),
				},
				"uuid": {
					S: aws.String(d.uuid),
				},
			},
			ReturnValues:     aws.String("ALL_NEW"),
			TableName:        aws.String(d.dynamoTable),
			UpdateExpression: aws.String("REMOVE #D.#K"),
		},
	)
	if err != nil {
		return err
	}

	// Update local
	data := make(map[string]interface{}, 0)
	err = dynamodbattribute.ConvertFromMap(res.Attributes["data"].M, &data)
	if err == nil {
		d.dynamoData = data
	}

	// OK
	return nil
}
