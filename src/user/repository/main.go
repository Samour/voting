package repository

import (
	"context"
	"errors"

	"github.com/Samour/voting/clients"
	"github.com/Samour/voting/user/model"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var userTableName = "users"
var usernamePasswordCredentialTableName = "username-password-credentials"

type UsernameUnavailableError struct{}

func (UsernameUnavailableError) Error() string {
	return "UsernameUnavailableError"
}

func InsertNewUser(user model.User, credential model.UsernamePasswordCredential) error {
	client := clients.DynamoDb()
	userItem, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}
	credentialItem, err := attributevalue.MarshalMap(credential)
	if err != nil {
		return err
	}

	var transact []types.TransactWriteItem

	userCondition := "attribute_not_exists(UserId)"
	transact = append(transact, types.TransactWriteItem{
		Put: &types.Put{
			TableName:           &userTableName,
			Item:                userItem,
			ConditionExpression: &userCondition,
		},
	})

	credentialCondition := "attribute_not_exists(Username)"
	transact = append(transact, types.TransactWriteItem{
		Put: &types.Put{
			TableName:           &usernamePasswordCredentialTableName,
			Item:                credentialItem,
			ConditionExpression: &credentialCondition,
		},
	})

	_, err = client.TransactWriteItems(context.Background(), &dynamodb.TransactWriteItemsInput{
		TransactItems: transact,
	})

	return translateInsertNewUserError(err)
}

func translateInsertNewUserError(err error) error {
	if err == nil {
		return nil
	}

	var txnCancelled *types.TransactionCanceledException
	if !errors.As(err, &txnCancelled) {
		return err
	}
	if txnCancelled.CancellationReasons[1].Code == nil {
		return err
	}
	if *txnCancelled.CancellationReasons[1].Code == "ConditionalCheckFailed" {
		return UsernameUnavailableError{}
	}

	return err
}

func LoadUsernamePasswordCredential(username string) (model.UsernamePasswordCredential, error) {
	client := clients.DynamoDb()
	item, err := client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: &usernamePasswordCredentialTableName,
		Key: map[string]types.AttributeValue{
			"Username": &types.AttributeValueMemberS{
				Value: username,
			},
		},
	})
	if err != nil {
		return model.UsernamePasswordCredential{}, err
	}

	if len(item.Item) == 0 {
		return model.UsernamePasswordCredential{}, nil
	}

	credential := model.UsernamePasswordCredential{}
	err = attributevalue.UnmarshalMap(item.Item, &credential)
	return credential, err
}
