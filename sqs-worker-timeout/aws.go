package sqs_worker

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func GetSqsClient(ctx context.Context) *sqs.Client {
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Could not load default configuration: %v", err)
	}
	sqsClient := sqs.NewFromConfig(awsConfig)
	return sqsClient
}
