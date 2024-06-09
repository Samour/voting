package polls

import (
	"context"

	"github.com/Samour/voting/clients"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var tableName = "polls"

func insertNewPollItem(p *Poll) error {
	client := clients.DynamoDb()
	condition := "attribute_not_exists(PollId)"
	item, err := attributevalue.MarshalMap(p)
	if err != nil {
		return err
	}

	_, err = client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName:           &tableName,
		Item:                item,
		ConditionExpression: &condition,
	})

	return err
}

func getPollItem(id string) (*Poll, error) {
	client := clients.DynamoDb()
	item, err := client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"PollId": &types.AttributeValueMemberS{
				Value: id,
			},
			"Discriminator": &types.AttributeValueMemberS{
				Value: "poll",
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(item.Item) == 0 {
		return nil, nil
	}

	poll := &Poll{}
	err = attributevalue.UnmarshalMap(item.Item, poll)
	if err != nil {
		return nil, err
	}

	return poll, nil
}

func updatePollItem(p *Poll) error {
	client := clients.DynamoDb()
	condition := "attribute_exists(PollId)"
	item, err := attributevalue.MarshalMap(p)
	if err != nil {
		return err
	}

	_, err = client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName:           &tableName,
		Item:                item,
		ConditionExpression: &condition,
	})

	return err
}
