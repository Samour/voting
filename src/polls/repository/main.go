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

type Paged[D any] struct {
	Items            []D
	LastEvaluatedKey *string
}

func InsertNewPollItem(p *model.Poll) error {
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

func GetPollItem(id string) (*model.Poll, error) {
	client := clients.DynamoDb()
	item, err := client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"PollId": &types.AttributeValueMemberS{
				Value: id,
			},
			"Discriminator": &types.AttributeValueMemberS{
				Value: model.DiscriminatorPoll,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(item.Item) == 0 {
		return nil, nil
	}

	poll := &model.Poll{}
	err = attributevalue.UnmarshalMap(item.Item, poll)
	if err != nil {
		return nil, err
	}

	return poll, nil
}

func UpdatePollItem(p *model.Poll) error {
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

func RecordVote(v *model.Vote) error {
	client := clients.DynamoDb()

	voteItem, err := attributevalue.MarshalMap(v)
	if err != nil {
		return err
	}

	var writes []types.TransactWriteItem

	voteCondition := "attribute_not_exists(PollId)"
	writes = append(writes, types.TransactWriteItem{
		Put: &types.Put{
			TableName:           &tableName,
			Item:                voteItem,
			ConditionExpression: &voteCondition,
		},
	})

	statisticsUpdate := "SET Statistics.Votes = Statistics.Votes + :inc"
	statisticsCondition := "attribute_exists(PollId)"
	writes = append(writes, types.TransactWriteItem{
		Update: &types.Update{
			TableName: &tableName,
			Key: map[string]types.AttributeValue{
				"PollId": &types.AttributeValueMemberS{
					Value: v.PollId,
				},
				"Discriminator": &types.AttributeValueMemberS{
					Value: model.DiscriminatorPoll,
				},
			},
			UpdateExpression:    &statisticsUpdate,
			ConditionExpression: &statisticsCondition,
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":inc": &types.AttributeValueMemberN{
					Value: "1",
				},
			},
		},
	})

	_, err = client.TransactWriteItems(context.Background(), &dynamodb.TransactWriteItemsInput{
		TransactItems: writes,
	})
	return err
}

func GetPollVoteItems(id string, exclusiveStartKey *string) (*Paged[model.Vote], error) {
	client := clients.DynamoDb()

	var esk map[string]types.AttributeValue = nil
	if exclusiveStartKey != nil {
		esk = map[string]types.AttributeValue{
			"PollId": &types.AttributeValueMemberS{
				Value: id,
			},
			"Discriminator": &types.AttributeValueMemberS{
				Value: *exclusiveStartKey,
			},
		}
	}

	keyConditionExpression := "PollId = :poll_id and begins_with(Discriminator, :discriminator)"
	items, err := client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:              &tableName,
		KeyConditionExpression: &keyConditionExpression,
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":poll_id": &types.AttributeValueMemberS{
				Value: id,
			},
			":discriminator": &types.AttributeValueMemberS{
				Value: model.DiscriminatorVote,
			},
		},
		ExclusiveStartKey: esk,
	})
	if err != nil {
		return nil, err
	}

	results := make([]model.Vote, items.Count)
	err = attributevalue.UnmarshalListOfMaps(items.Items, &results)
	if err != nil {
		return nil, err
	}

	var lastEvaluatedKey *string = nil
	if items.LastEvaluatedKey != nil {
		lastVote := &model.Vote{}
		err := attributevalue.UnmarshalMap(items.LastEvaluatedKey, lastVote)
		if err != nil {
			return nil, err
		}
		lastEvaluatedKey = &lastVote.Discriminator
	}

	return &Paged[model.Vote]{
		Items:            results,
		LastEvaluatedKey: lastEvaluatedKey,
	}, nil
}

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
