package sqs_worker

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func LoadAwsConfig(ctx context.Context) aws.Config {
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Could not load default configuration: %v", err)
	}
	return awsConfig
}

func GetSqsClient(ctx context.Context) *sqs.Client {
	awsConfig := LoadAwsConfig(ctx)
	sqsClient := sqs.NewFromConfig(awsConfig)
	return sqsClient
}

func GetSfnClient(ctx context.Context) *sfn.Client {
	awsConfig := LoadAwsConfig(ctx)
	sfnClient := sfn.NewFromConfig(awsConfig)
	return sfnClient
}

func GetDynamoClient(ctx context.Context) *dynamodb.Client {
	awsConfig := LoadAwsConfig(ctx)
	dynamoClient := dynamodb.NewFromConfig(awsConfig)
	return dynamoClient
}

func GetEcsClient(ctx context.Context) *ecs.Client {
	awsConfig := LoadAwsConfig(ctx)
	ecsClient := ecs.NewFromConfig(awsConfig)
	return ecsClient
}

func GetAwsContext(parentContext context.Context) context.Context {
	if parentContext != nil {
		return context.TODO()
	}
	return parentContext
}
