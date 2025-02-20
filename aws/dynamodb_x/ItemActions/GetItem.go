package item_actions

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	definitions "github.com/mauri-codes/go-modules/aws/dynamodb_x/Definitions"
)

func GetItem[T any](table *definitions.Table, action definitions.IItemAction[T]) (T, error) {
	client := table.Client
	var item T
	var itemKeys map[string]types.AttributeValue
	pk, err := attributevalue.Marshal(action.GetHashKeyValue())
	if err != nil {
		return item, err
	}
	if action.GetSortKeyValue() != "" {
		sk, err := attributevalue.Marshal(action.GetSortKeyValue())
		if err != nil {
			return item, err
		}
		itemKeys = map[string]types.AttributeValue{"pk": pk, "sk": sk}
	} else {
		itemKeys = map[string]types.AttributeValue{"pk": pk}
	}
	response, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: itemKeys, TableName: aws.String(table.TableName),
	})
	if err != nil {
		log.Printf("GetItem Error: %v\n", err)
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &item)
		if err != nil {
			log.Printf("Couldn't unmarshal response: %v\n", err)
		}
	}
	return item, err
}
