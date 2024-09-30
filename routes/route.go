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
		{"POST", "/auth/v1/login"}: controllers.SignIn,

		// need bear auth header
		{"POST", "/user/v1/refresh"}: controllers.RefreshToken,
		{"POST", "/user/v1/revoke"}:  controllers.RevokeToken,
		{"POST", "/user/v1/logout"}:  controllers.GlobalSignOut,

		// s3
		{"GET", "/s3/fetch"}:     controllers.GetS3Obj,
		{"GET", "/s3/fetch-url"}: controllers.GetS3ObjRequest,
	}
}
