package repository

import (
	"context"

	"github.com/Samour/voting/clients"
	"github.com/Samour/voting/polls/model"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func RecordVote(pollId string, v interface{}) error {
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
					Value: pollId,
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

func GetPollVoteItems(id string, exclusiveStartKey *string, out interface{}) (*string, error) {
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

	err = attributevalue.UnmarshalListOfMaps(items.Items, out)
	if err != nil {
		return nil, err
	}

	if items.LastEvaluatedKey != nil {
		lastVote := &model.FptpVote{}
		err := attributevalue.UnmarshalMap(items.LastEvaluatedKey, lastVote)
		if err != nil {
			return nil, err
		}
		return &lastVote.Discriminator, nil
	}

	return nil, nil
}
