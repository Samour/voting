package repository

import (
	"context"

	"github.com/Samour/voting/clients"
	"github.com/Samour/voting/polls/model"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetPollResultItem(id string) (*model.PollResult, error) {
	client := clients.DynamoDb()

	item, err := client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"PollId": &types.AttributeValueMemberS{
				Value: id,
			},
			"Discriminator": &types.AttributeValueMemberS{
				Value: model.DiscriminatorResult,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(item.Item) == 0 {
		return nil, nil
	}

	result := &model.PollResult{}
	err = attributevalue.UnmarshalMap(item.Item, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func InsertNewPollResultItem(r *model.PollResult) error {
	client := clients.DynamoDb()
	item, err := attributevalue.MarshalMap(r)
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
