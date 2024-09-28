package controllers

import (
	"context"
	"encoding/json"
	"portfolio/v2/services"

	"github.com/aws/aws-lambda-go/events"
)

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	Username     string `json:"username"`
	RefreshToken string `json:"refresh_token"`
}

type RevokeTokenRequest struct {
	Token string `json:"token"`
}

type ForgotPasswordRequest struct {
	Username string `json:"username"`
}

type ConfirmForgotPasswordRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

func SignIn(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var reqBody SignInRequest
	var err error
	if parseErr := json.Unmarshal([]byte(request.Body), &reqBody); parseErr != nil {
		return GetBadResponse(parseErr), parseErr
	}
	cognito, err := services.GetCognitoActions()
	if err != nil {
		return GetBadResponse(err), err
	}

	secretHash, err := cognito.CalcSecretHash(reqBody.Username)
	if err != nil {
		return GetBadResponse(err), err
	}

	result, err := cognito.SignIn(ctx, reqBody.Username, reqBody.Password, secretHash)
	if err != nil {
		return GetBadResponse(err), err
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		return GetBadResponse(err), err
	}

	return GetOKResponse(string(bytes)), nil
}

func RefreshToken(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var reqBody RefreshTokenRequest
	var err error
	if parseErr := json.Unmarshal([]byte(request.Body), &reqBody); parseErr != nil {
		return GetBadResponse(parseErr), parseErr
	}
	cognito, err := services.GetCognitoActions()
	if err != nil {
		return GetBadResponse(err), err
	}

	secretHash, err := cognito.CalcSecretHash(reqBody.Username)
	if err != nil {
		return GetBadResponse(err), err
	}

	result, err := cognito.RefreshToken(ctx, reqBody.RefreshToken, secretHash)
	if err != nil {
		return GetBadResponse(err), err
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		return GetBadResponse(err), err
	}

	return GetOKResponse(string(bytes)), nil
}

func RevokeToken(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var reqBody RevokeTokenRequest
	var err error
	if paresErr := json.Unmarshal([]byte(request.Body), &reqBody); paresErr != nil {
		return GetBadResponse(paresErr), paresErr
	}
	cognito, err := services.GetCognitoActions()
	if err != nil {
		return GetBadResponse(err), err
	}

	result, err := cognito.RevokeToken(ctx, reqBody.Token)
	if err != nil {
		return GetBadResponse(err), err
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		return GetBadResponse(err), err
	}

	return GetOKResponse(string(bytes)), nil
}

func ConfirmForgotPassword(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var reqBody ConfirmForgotPasswordRequest
	var err error
	if paresErr := json.Unmarshal([]byte(request.Body), &reqBody); paresErr != nil {
		return GetBadResponse(paresErr), paresErr
	}
	cognito, err := services.GetCognitoActions()
	if err != nil {
		return GetBadResponse(err), err
	}

	secretHash, err := cognito.CalcSecretHash(reqBody.Username)
	if err != nil {
		return GetBadResponse(err), err
	}

	result, err := cognito.ConfirmForgotPassword(ctx, reqBody.Username, reqBody.Password, reqBody.Code, secretHash)
	if err != nil {
		return GetBadResponse(err), err
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		return GetBadResponse(err), err
	}

	return GetOKResponse(string(bytes)), nil
}
