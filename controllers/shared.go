package controllers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func GetBadResponse(err error) events.APIGatewayV2HTTPResponse {
	var response events.APIGatewayV2HTTPResponse
	response.Body = err.Error()
	response.StatusCode = http.StatusBadRequest
	return response
}

func GetOKResponse(body string) events.APIGatewayV2HTTPResponse {
	var response events.APIGatewayV2HTTPResponse
	response.Body = body
	response.StatusCode = http.StatusOK
	return response
}

func GetNotFoundResponse() events.APIGatewayV2HTTPResponse {
	var response events.APIGatewayV2HTTPResponse
	response.Body = "not found"
	response.StatusCode = http.StatusNotFound
	return response
}
