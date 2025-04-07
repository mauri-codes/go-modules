package utils

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type LambdaResponse struct {
	Success bool
	Message *string
	Data    any
}

type ResponseInput struct {
	Message string
	Code    int
	Data    any
}

func Error400(input ResponseInput) (events.APIGatewayProxyResponse, error) {
	input.Code = 400
	return ErrorResponse(input)
}

func Error500(input ResponseInput) (events.APIGatewayProxyResponse, error) {
	input.Code = 500
	return ErrorResponse(input)
}

func ErrorResponse(input ResponseInput) (events.APIGatewayProxyResponse, error) {
	out, _ := json.Marshal(LambdaResponse{
		Message: &input.Message,
		Success: false,
		Data:    input.Message,
	})
	return events.APIGatewayProxyResponse{Body: string(out), StatusCode: input.Code}, nil
}

func SuccessResponse(event ResponseInput) (events.APIGatewayProxyResponse, error) {
	var response = LambdaResponse{
		Success: true,
		Message: &event.Message,
		Data:    event.Data,
	}
	out, _ := json.Marshal(response)
	return events.APIGatewayProxyResponse{Body: string(out), StatusCode: 200}, nil
}
