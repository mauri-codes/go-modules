package sqs_worker

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SqsWorker struct {
}

func CreateSqsWorker(config *Configuration, process *Process) {
	awsContext := context.TODO()
	sqsClient := GetSqsClient(awsContext)
	idleTimer, resetChan, timeoutCtx, cancel := SetShutDownConditions(SetShutDownConditionsInput{
		Configuration:            config,
		ShouldKeepAliveOnTimeOut: process.ShouldKeepAliveOnTimeOut,
		ShutDownAction:           process.ShutDownAction,
	})

	wg := PollSqs(&PollSqsInput{
		Config:         config,
		TimeoutCtx:     timeoutCtx,
		AwsCtx:         awsContext,
		SqsClient:      sqsClient,
		ResetChan:      resetChan,
		ProcessMessage: process.ProcessMessage,
	})
	log.Println("Waiting for in-progress workers to complete...")
	(*wg).Wait()
	log.Println("All workers finished. Exiting.")
	defer cancel()
	defer idleTimer.Stop()
}

type Process struct {
	ShouldKeepAliveOnTimeOut func() bool
	ProcessMessage           func(message types.Message)
	ShutDownAction           func()
}

type Configuration struct {
	IdleTimeout    int
	MaxConcurrency int
	QueueUrl       string
}
