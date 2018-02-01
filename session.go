package dynamostore

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func getDynamoSession(region, accessKeyID, secretAccessKey string) *dynamodb.DynamoDB {
	awsSession := session.New(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewCredentials(&environmentCredentialsProvider{
			AccessKeyID:     accessKeyID,
			SecretAccessKey: secretAccessKey,
		}),
	})
	return dynamodb.New(awsSession)
}

type environmentCredentialsProvider struct {
	AccessKeyID     string
	SecretAccessKey string
}

func (e *environmentCredentialsProvider) Retrieve() (credentials.Value, error) {
	return credentials.Value{
		AccessKeyID:     e.AccessKeyID,
		SecretAccessKey: e.SecretAccessKey,
	}, nil
}

func (e *environmentCredentialsProvider) IsExpired() bool {
	return false
}
