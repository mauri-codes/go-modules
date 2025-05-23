package sqs_worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type PollSqsInput struct {
	Config         *Configuration
	TimeoutCtx     context.Context
	AwsCtx         context.Context
	SqsClient      *sqs.Client
	ResetChan      chan struct{}
	ProcessMessage func(message types.Message)
}

func PollSqs(input *PollSqsInput) *sync.WaitGroup {
	log.Println("PollSqs")
	workerPool := make(chan struct{}, 2*input.Config.MaxConcurrency)
	var wg sync.WaitGroup
	for {
		select {
		case <-input.TimeoutCtx.Done():
			log.Println("Context cancelled, exiting processor loop")
			return &wg
		default:
			log.Println("poll-01")
			if len(workerPool) >= input.Config.MaxConcurrency {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			log.Println("poll-02")

			resp, err := input.SqsClient.ReceiveMessage(input.AwsCtx, &sqs.ReceiveMessageInput{
				QueueUrl:            &input.Config.QueueUrl,
				MaxNumberOfMessages: 5,
				WaitTimeSeconds:     20,
			})
			log.Println("poll-03")
			if err != nil {
				fmt.Println("Error receiving message:", err)
				time.Sleep(2 * time.Second)
				continue
			}

			if len(resp.Messages) == 0 {
				log.Println("poll-no-messages")
				continue
			}
			log.Println("poll-reset-chan")
			input.ResetChan <- struct{}{}

			for _, msg := range resp.Messages {
				fmt.Println("Got message:", *msg.Body)
				workerPool <- struct{}{}
				wg.Add(1)
				log.Println("poll-04")
				go func(m types.Message) {
					log.Println("m")
					log.Println(*m.Body)
					log.Println(*msg.Body)
					input.ProcessMessage(m)
					log.Println("poll-05")
					_, err := input.SqsClient.DeleteMessage(input.AwsCtx, &sqs.DeleteMessageInput{
						QueueUrl:      &input.Config.QueueUrl,
						ReceiptHandle: msg.ReceiptHandle,
					})
					log.Println("poll-06")
					if err != nil {
						fmt.Println("Error deleting message:", err)
					}
					defer func() {
						<-workerPool
						wg.Done()
					}()
				}(msg)
			}
		}
	}
}
