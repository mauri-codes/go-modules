package dynamo

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func Query[T any](table *Table, action IItemAction[T]) ([]T, error) {
	client := table.Client
	hash := expression.Key(table.HashKey).Equal(expression.Value(action.GetHashKeyValue()))
	var expr expression.Expression
	var err error
	if action.GetSortKeyValue() != "" {
		var sort expression.KeyConditionBuilder
		if action.GetSortKeyAction() == BEGINS_WITH {
			sort = expression.Key(table.SortKey).BeginsWith(action.GetSortKeyValue())
		} else {
			sort = expression.Key(table.SortKey).Equal(expression.Value(action.GetSortKeyValue()))
		}
		expr, err = expression.NewBuilder().WithKeyCondition(hash.And(sort)).Build()
	} else {
		expr, err = expression.NewBuilder().WithKeyCondition(hash).Build()
	}
	var queryResponse *dynamodb.QueryOutput
	var data []T
	if err != nil {
		log.Printf("Couldn't build expression for query: %v\n", err)
	} else {
		queryPaginator := dynamodb.NewQueryPaginator(client, &dynamodb.QueryInput{
			TableName:                 aws.String(table.TableName),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			KeyConditionExpression:    expr.KeyCondition(),
		})
		for queryPaginator.HasMorePages() {
			queryResponse, err = queryPaginator.NextPage(context.TODO())
			if err != nil {
				log.Printf("Couldn't query for data: %v\n", err)
				break
			} else {
				var moviePage []T
				err = attributevalue.UnmarshalListOfMaps(queryResponse.Items, &moviePage)
				if err != nil {
					log.Printf("Couldn't unmarshal query response: %v\n", err)
					break
				} else {
					data = append(data, moviePage...)
				}
			}
		}
	}
	return data, err
}
