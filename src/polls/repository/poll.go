package repository

import (
	"context"

	"github.com/Samour/voting/clients"
	"github.com/Samour/voting/polls/model"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var tableName = "polls"

func InsertNewPollItem(p interface{}) error {
	client := clients.DynamoDb()
	item, err := attributevalue.MarshalMap(p)
	if err != nil {
		return err
	}

	condition := "attribute_not_exists(PollId)"
	_, err = client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName:           &tableName,
		Item:                item,
		ConditionExpression: &condition,
	})

	return err
}

func GetPollItem(id string, discriminator string, p interface{}) error {
	client := clients.DynamoDb()
	item, err := client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"PollId": &types.AttributeValueMemberS{
				Value: id,
			},
			"Discriminator": &types.AttributeValueMemberS{
				Value: discriminator,
			},
		},
	})
	if err != nil {
		return err
	}

	if len(item.Item) == 0 {
		return nil
	}

	err = attributevalue.UnmarshalMap(item.Item, p)
	if err != nil {
		return err
	}

	return nil
}

func UpdatePollItem(p model.Poll) error {
	client := clients.DynamoDb()
	item, err := attributevalue.MarshalMap(p)
	if err != nil {
		return err
	}

	condition := "attribute_exists(PollId)"
	_, err = client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName:           &tableName,
		Item:                item,
		ConditionExpression: &condition,
	})

	return err
}

func ScanPollItems() ([]model.Poll, error) {
	client := clients.DynamoDb()
	filterExpression := "Discriminator = :discriminator"
	items, err := client.Scan(context.Background(), &dynamodb.ScanInput{
		TableName:        &tableName,
		FilterExpression: &filterExpression,
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":discriminator": &types.AttributeValueMemberS{
				Value: model.DiscriminatorPoll,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.Poll, items.Count)
	err = attributevalue.UnmarshalListOfMaps(items.Items, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
