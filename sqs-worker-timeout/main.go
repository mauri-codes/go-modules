package sqs_worker

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SqsWorkerInput struct {
	Config     *Configuration
	Process    *Process
	AwsContext context.Context
}

func CreateSqsWorker(input SqsWorkerInput) {
	awsContext := GetAwsContext(input.AwsContext)
	sqsClient := GetSqsClient(awsContext)
	idleTimer, resetChan, timeoutCtx, cancel := SetShutDownConditions(SetShutDownConditionsInput{
		Configuration:            input.Config,
		ShouldKeepAliveOnTimeOut: input.Process.ShouldKeepAliveOnTimeOut,
		ShutDownAction:           input.Process.ShutDownAction,
	})

	wg := PollSqs(&PollSqsInput{
		Config:         input.Config,
		TimeoutCtx:     timeoutCtx,
		AwsCtx:         awsContext,
		SqsClient:      sqsClient,
		ResetChan:      resetChan,
		ProcessMessage: input.Process.ProcessMessage,
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
