package item_actions

import (
	"context"
	"dynamodb/definitions"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func PutItem[T any](table *definitions.Table, action definitions.IItemAction[T]) error {
	client := table.Client
	item, err := attributevalue.MarshalMap(action.GetData())
	if err != nil {
		return err
	}
	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(table.TableName), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table: %v\n", err)
	}
	return err
}
