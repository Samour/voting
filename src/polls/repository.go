package polls

import (
	"context"

	"github.com/Samour/voting/clients"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func insertNewPollItem(p *PollItem) error {
	client := clients.DynamoDb()
	tableName := "polls"
	condition := "attribute_not_exists(PollId)"
	_, err := client.PutItem(context.Background(), &dynamodb.PutItemInput{
		Item:                p.toDynamoDbMap(),
		TableName:           &tableName,
		ConditionExpression: &condition,
	})

	return err
}
