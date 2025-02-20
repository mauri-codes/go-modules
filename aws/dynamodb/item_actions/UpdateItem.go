package item_actions

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mauri-codes/go-modules/aws/dynamodb/definitions"
)

func UpdateItem[T any](table *definitions.Table, action definitions.IItemAction[T]) error {
	client := table.Client
	var err error
	var itemKeys map[string]types.AttributeValue
	pk, err := attributevalue.Marshal(action.GetHashKeyValue())
	if err != nil {
		return err
	}
	if action.GetSortKeyValue() != "" {
		sk, err := attributevalue.Marshal(action.GetSortKeyValue())
		if err != nil {
			return err
		}
		itemKeys = map[string]types.AttributeValue{table.HashKey: pk, table.SortKey: sk}
	} else {
		itemKeys = map[string]types.AttributeValue{table.HashKey: pk}
	}
	exp := action.GetExpression()
	_, err = client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName:                 aws.String(table.TableName),
		Key:                       itemKeys,
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
		UpdateExpression:          exp.Update(),
		ReturnValues:              types.ReturnValueUpdatedNew,
	})
	if err != nil {
		log.Printf("Couldn't update item: %v\n", err)
	}
	return err
}
