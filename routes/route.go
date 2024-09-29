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

func GetRoutes() Routes {
	return Routes{
		{"POST", "/auth/v1/login"}:   controllers.SignIn,
		{"POST", "/user/v1/refresh"}: controllers.RefreshToken,
		{"POST", "/user/v1/revoke"}:  controllers.RevokeToken,
		{"GET", "/user/v1/logout"}:   controllers.GlobalSignOut,
	}
}
