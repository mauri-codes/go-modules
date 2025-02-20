package definitions

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Table struct {
	TableName string
	HashKey   string
	SortKey   string
	Client    *dynamodb.Client
}

func NewTable(name string, hashKey string, sortKey string, client *dynamodb.Client) *Table {
	return &Table{
		TableName: name,
		HashKey:   hashKey,
		SortKey:   sortKey,
		Client:    client,
	}
}

type ItemAction[T any] struct {
	HashKeyValue  string
	SortKeyValue  string
	Data          T
	Expression    expression.Expression
	SortKeyAction string
}

func (tableAction *ItemAction[T]) GetData() T {
	return tableAction.Data
}

func (tableAction *ItemAction[T]) GetHashKeyValue() string {
	return tableAction.HashKeyValue
}

func (tableAction *ItemAction[T]) GetSortKeyValue() string {
	return tableAction.SortKeyValue
}

func (tableAction *ItemAction[T]) GetSortKeyAction() string {
	return tableAction.SortKeyAction
}

func (tableAction *ItemAction[T]) GetExpression() expression.Expression {
	return tableAction.Expression
}

type IItemAction[T any] interface {
	GetHashKeyValue() string
	GetSortKeyValue() string
	GetSortKeyAction() string
	GetData() T
	GetExpression() expression.Expression
}

func NewPutItem[T any](hashKey string, sortKey string, data T) IItemAction[T] {
	return &ItemAction[T]{
		HashKeyValue: hashKey,
		SortKeyValue: sortKey,
		Data:         data,
	}
}

func NewQuery[T any](hashKey string, sortKey string) IItemAction[T] {
	return &ItemAction[T]{
		HashKeyValue: hashKey,
		SortKeyValue: sortKey,
	}
}

const (
	EQUALS      = "EQUALS"
	BEGINS_WITH = "BEGINS_WITH"
)
