package sqs_worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/mauri-codes/go-modules/aws/dynamo"
)

type EcsSfnConfig struct {
	TableName   string
	ServiceName string
	ClusterName string
}

type SqsWorkerForEcsSfnInput[SqsMessage any] struct {
	Config         *Configuration
	EcsSfnConfig   *EcsSfnConfig
	AwsContext     context.Context
	MessageProcess func(message SqsMessage) bool
}

func CreateSqsWorkerForEcsSfn[SqsMessage any](input SqsWorkerForEcsSfnInput[SqsMessage]) {
	awsContext := GetAwsContext(input.AwsContext)
	sfnClient := GetSfnClient(awsContext)
	ecsClient := GetEcsClient(awsContext)
	dynamoClient := GetDynamoClient(awsContext)
	CreateSqsWorker(SqsWorkerInput{
		Config: input.Config,
		Process: &Process{
			ShouldKeepAliveOnTimeOut: func() bool {
				return ShouldKeepAliveOnTimeOut(ShouldKeepAliveOnTimeOutInput{
					Config:       input.EcsSfnConfig,
					EcsClient:    ecsClient,
					DynamoClient: dynamoClient,
				}, awsContext)
			},
			ShutDownAction: func() {
				ShutDownAction(ShutDownActionInput{
					Config:    input.EcsSfnConfig,
					EcsClient: ecsClient,
				}, awsContext)
			},
			ProcessMessage: func(message types.Message) {
				ProcessMessage(&ProcessMessageInput[SqsMessage]{
					StepFunctionsClient: sfnClient,
					Message:             message,
				}, awsContext)
			},
		},
	})
}

type ShouldKeepAliveOnTimeOutInput struct {
	Config       *EcsSfnConfig
	EcsClient    *ecs.Client
	DynamoClient *dynamodb.Client
}

type TesterStatus struct {
	ServiceStop int
}

func ShouldKeepAliveOnTimeOut(input ShouldKeepAliveOnTimeOutInput, ctx context.Context) bool {
	table := dynamo.NewTable(input.Config.TableName, "pk", "sk", input.DynamoClient)
	action := dynamo.NewGetItem[TesterStatus]("tester", "activity")
	testerStatus, err := dynamo.GetItem(table, action)
	if err != nil {
		log.Println("Could not get tester status")
		return false
	}
	log.Println(testerStatus)
	ecsService, ecsError := input.EcsClient.DescribeServices(ctx, &ecs.DescribeServicesInput{
		Services: []string{input.Config.ServiceName},
		Cluster:  &input.Config.ClusterName,
	})
	if ecsError != nil {
		log.Println("Could not get ECS Service Data")
		return false
	}
	if len(ecsService.Services) == 0 {
		log.Println("No ECS Service Found")
		return false
	}
	now := time.Now().UnixMilli()
	serviceStopInTimeout := now < int64(testerStatus.ServiceStop)
	log.Println("serviceStopInTimeout")
	log.Println(serviceStopInTimeout)
	multipleServicesRunning := ecsService.Services[0].RunningCount > 1
	log.Println(multipleServicesRunning)
	if serviceStopInTimeout || multipleServicesRunning {
		return true
	}
	return false
}

type ShutDownActionInput struct {
	Config    *EcsSfnConfig
	EcsClient *ecs.Client
}

func ShutDownAction(input ShutDownActionInput, ctx context.Context) {
	var desiredCount int32 = 0
	input.EcsClient.UpdateService(ctx, &ecs.UpdateServiceInput{
		Service:      &input.Config.ServiceName,
		Cluster:      &input.Config.ClusterName,
		DesiredCount: &desiredCount,
	})
}

type ProcessMessageInput[SqsMessage any] struct {
	StepFunctionsClient *sfn.Client
	Message             types.Message
	ProcessMessage      func(message SqsMessage) bool
}

type SqsPayload[SqsMessage any] struct {
	StepFunctionsToken string
	Payload            SqsMessage
}

func ProcessMessage[SqsMessage any](input *ProcessMessageInput[SqsMessage], ctx context.Context) {
	var sqsPayload SqsPayload[SqsMessage]
	err := json.Unmarshal([]byte(*input.Message.Body), &sqsPayload)
	log.Println(sqsPayload)
	if err != nil {
		log.Printf("Could not transform sqs message %v", err)
		return
	}
	isSuccessful := input.ProcessMessage(sqsPayload.Payload)
	if isSuccessful {
		NotifyStepFunctions(NotifyStepFunctionsInput{
			Success: true, StepFunctionsClient: input.StepFunctionsClient, Token: sqsPayload.StepFunctionsToken, Output: "{\"success\": true}",
		}, ctx)
	} else {
		NotifyStepFunctions(NotifyStepFunctionsInput{
			Success: false, StepFunctionsClient: input.StepFunctionsClient, Token: sqsPayload.StepFunctionsToken,
		}, ctx)
	}
}

type NotifyStepFunctionsInput struct {
	Success             bool
	StepFunctionsClient *sfn.Client
	Token               string
	Output              string
}

func NotifyStepFunctions(input NotifyStepFunctionsInput, ctx context.Context) {
	log.Println(input.Token)
	if input.Success {
		_, err := input.StepFunctionsClient.SendTaskSuccess(ctx, &sfn.SendTaskSuccessInput{
			TaskToken: &input.Token,
			Output:    &input.Output,
		})
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		_, err := input.StepFunctionsClient.SendTaskFailure(ctx, &sfn.SendTaskFailureInput{
			TaskToken: &input.Token,
		})
		if err != nil {
			log.Println(err.Error())
		}
	}
}
