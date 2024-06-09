package polls

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

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
