package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"portfolio/v2/controllers"
	"portfolio/v2/routes"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func getHandler(request events.APIGatewayV2HTTPRequest, routesConfig routes.Routes) (func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error), error) {
	routeKeyStr := request.RouteKey
	routeKeys := strings.Split(routeKeyStr, " ")
	httpMethod := strings.ToUpper(routeKeys[0])
	endPoint := strings.ToLower(routeKeys[1])
	// check if wildcard is used in api gateway
	if strings.Contains(endPoint, "/{proxy+}") {
		endPoint = strings.ToLower(request.RawPath)
	}
	if handler, ok := routesConfig[routes.RouteKey{Method: httpMethod, Endpoint: endPoint}]; ok {
		return handler, nil
	} else {
		return nil, fmt.Errorf("not found: %s", routeKeyStr)
	}
}

func handleRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// auth are handled in api gateway
	routesConfig := routes.GetRoutes()
	handler, err := getHandler(request, routesConfig)
	if err != nil {
		log.Println(err.Error())
		return controllers.GetNotFoundResponse(), nil
	}
	var resp events.APIGatewayV2HTTPResponse

	resp, err = handler(ctx, request)
	if err != nil {
		log.Println(err.Error())
	}
	if resp.Headers == nil {
		resp.Headers = make(map[string]string)
	}
	return resp, nil
}

func main() {
	lambda.Start(handleRequest)
}
