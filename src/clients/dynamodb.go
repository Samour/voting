package clients

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var dynamodbClient *dynamodb.Client

func DynamoDb() *dynamodb.Client {
	return dynamodbClient
}

type staticCredentials struct{}

func (c staticCredentials) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     "dummy",
		SecretAccessKey: "dummy",
		SessionToken:    "dummy",
		Source:          "static",
	}, nil
}

func WarmDynamoDbClient() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithCredentialsProvider(staticCredentials{}))
	if err != nil {
		panic(err)
	}

	dynamodbClient = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String("http://localhost:8000")
	})
}
