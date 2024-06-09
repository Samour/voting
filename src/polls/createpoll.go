package polls

import (
	"context"

	"github.com/Samour/voting/clients"
	"github.com/Samour/voting/utils"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type PollItem struct {
	PollId        string
	Discriminator string
	Status        string
	Name          string
	Options       []string
}

func (p *PollItem) toDynamoDbMap() map[string]types.AttributeValue {
	var optionsList []types.AttributeValue
	for _, option := range p.Options {
		optionsList = append(optionsList, &types.AttributeValueMemberS{Value: option})
	}

	return map[string]types.AttributeValue{
		"PollId": &types.AttributeValueMemberS{
			Value: p.PollId,
		},
		"Discriminator": &types.AttributeValueMemberS{
			Value: p.Discriminator,
		},
		"Status": &types.AttributeValueMemberS{
			Value: p.Status,
		},
		"Name": &types.AttributeValueMemberS{
			Value: p.Name,
		},
		"Options": &types.AttributeValueMemberL{
			Value: optionsList,
		},
	}
}

func CreatePoll() (*string, error) {
	id := utils.IdGen()
	poll := PollItem{
		PollId:        id,
		Discriminator: "poll",
		Status:        "draft",
		Name:          "",
		Options:       []string{""},
	}

	client := clients.DynamoDb()
	tableName := "polls"
	condition := "attribute_not_exists(PollId)"
	_, err := client.PutItem(context.Background(), &dynamodb.PutItemInput{
		Item:                poll.toDynamoDbMap(),
		TableName:           &tableName,
		ConditionExpression: &condition,
	})
	if err != nil {
		return nil, err
	}

	return &id, nil
}
