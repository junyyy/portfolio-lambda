package routes

import (
	"context"

	"portfolio/v2/controllers"

	"github.com/aws/aws-lambda-go/events"
)

type RouteKey struct {
	Method   string
	Endpoint string
}

type Routes map[RouteKey]func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)

// Handler  func(ctx context.Context, request events.APIGatewayV2HTTPRequest)
func GetRoutes() Routes {
	return Routes{
		{"POST", "/auth/v1/login"}:   controllers.SignIn,
		{"POST", "/auth/v1/refresh"}: controllers.RefreshToken,
		{"POST", "/auth/v1/revoke"}:  controllers.RevokeToken,
	}
}
