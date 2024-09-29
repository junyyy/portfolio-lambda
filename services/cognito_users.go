package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	cp "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type CognitoConfig struct {
	CognitoClient   *cp.Client
	ClientId        string
	ClientSecretKey string
}

func (c CognitoConfig) SignIn(ctx context.Context, username string, password string, secretHash string) (*types.AuthenticationResultType, error) {
	var authResult *types.AuthenticationResultType

	output, err := c.CognitoClient.InitiateAuth(ctx, &cp.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		ClientId: aws.String(c.ClientId),
		AuthParameters: map[string]string{
			"USERNAME": username,
			"PASSWORD": password,
			// "SECRET_HASH": secretHash,
		},
	})
	if err != nil {
		var resetRequired *types.PasswordResetRequiredException
		if errors.As(err, &resetRequired) {
			log.Println(*resetRequired.Message)
		} else {
			log.Printf("Couldn't sign in user %v. Here's why: %v\n", username, err)
		}
	} else {
		authResult = output.AuthenticationResult
	}
	return authResult, err
}

func (c CognitoConfig) RefreshToken(ctx context.Context, refreshToken string, secretHash string) (*types.AuthenticationResultType, error) {
	output, err := c.CognitoClient.InitiateAuth(ctx, &cp.InitiateAuthInput{
		AuthFlow: "REFRESH_TOKEN_AUTH",
		ClientId: aws.String(c.ClientId),
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshToken,
			"SECRET_HASH":   secretHash,
		},
	})

	if err != nil {
		return nil, err
	} else {
		return output.AuthenticationResult, nil
	}
}

func (c CognitoConfig) GlobalSignOut(ctx context.Context, token string) (*cp.GlobalSignOutOutput, error) {
	return c.CognitoClient.GlobalSignOut(ctx, &cp.GlobalSignOutInput{AccessToken: aws.String(token)})
}

func (c CognitoConfig) ConfirmForgotPassword(ctx context.Context, username string, password string, code string, secretHash string) (*cp.ConfirmForgotPasswordOutput, error) {
	output, err := c.CognitoClient.ConfirmForgotPassword(ctx, &cp.ConfirmForgotPasswordInput{
		ClientId:         &c.ClientId,
		ConfirmationCode: &code,
		Username:         &username,
		Password:         &password,
		SecretHash:       &secretHash,
	})

	return output, err
}

func (c CognitoConfig) RevokeToken(ctx context.Context, token string) (*cp.RevokeTokenOutput, error) {
	output, err := c.CognitoClient.RevokeToken(ctx, &cp.RevokeTokenInput{
		ClientId:     &c.ClientId,
		ClientSecret: &c.ClientSecretKey,
		Token:        &token,
	})

	return output, err
}

func (c CognitoConfig) CalcSecretHash(username string) (string, error) {
	input := fmt.Sprintf("%s%s", username, c.ClientId)
	key := []byte(c.ClientSecretKey)

	h := hmac.New(sha256.New, key)

	_, err := h.Write([]byte(input))
	if err != nil {
		return "", err
	}

	sum := h.Sum(nil)

	encoded := base64.StdEncoding.EncodeToString(sum)

	return encoded, nil
}

func GetCognitoConfig() (*CognitoConfig, error) {
	region, regionOk := os.LookupEnv("REGION")
	clientID, clientIDOk := os.LookupEnv("COGNITO_CLIENT_ID")

	if regionOk && clientIDOk {
		cognito := CognitoConfig{
			CognitoClient: cp.NewFromConfig(aws.Config{
				Region: region,
			}),
			ClientId: clientID,
		}

		return &cognito, nil
	}

	return nil, errors.New("failed to read configs")
}
