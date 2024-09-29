package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BucketBasics struct {
	S3Client *s3.Client
}

func GetS3Obj(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	bucketName, ok := os.LookupEnv("PORTFOLIO_S3_BUCKET_NAME")
	var err error
	if !ok {
		err = errors.New("S3 config does not exist")
		return GetBadResponse(err), err
	}

	fileName, ok := request.QueryStringParameters["file"]

	if !ok {
		err = errors.New("missing file name")
		return GetBadResponse(err), err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return GetBadResponse(err), err
	}

	s3Client := s3.NewFromConfig(cfg)

	result, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})

	if err != nil {
		return GetBadResponse(err), err
	}

	defer result.Body.Close()

	buf := make([]byte, *result.ContentLength)
	log.Printf("content length: %d", *result.ContentLength)
	_, err = result.Body.Read(buf)
	if err != nil {
		return GetBadResponse(err), err
	}
	pdfBase64 := base64.StdEncoding.EncodeToString(buf)
	resp := GetOKResponse(pdfBase64)
	resp.Headers = map[string]string{
		"Content-Type":        "application/pdf",
		"Content-Disposition": fmt.Sprintf("attachment; filename=%s", fileName),
	}
	return resp, nil
}

func GetS3ObjRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	bucketName, ok := os.LookupEnv("PORTFOLIO_S3_BUCKET_NAME")
	var err error
	if !ok {
		err = errors.New("S3 config does not exist")
		return GetBadResponse(err), err
	}

	fileName, ok := request.QueryStringParameters["file"]

	if !ok {
		err = errors.New("missing file name")
		return GetBadResponse(err), err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return GetBadResponse(err), err
	}

	expires := time.Now().Add(time.Duration(time.Minute * 5))
	s3Client := s3.NewFromConfig(cfg)
	presigner := s3.NewPresignClient(s3Client)
	presignReq, err := presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket:          aws.String(bucketName),
		Key:             aws.String(fileName),
		ResponseExpires: &expires,
	})

	if err != nil {
		return GetBadResponse(err), err
	}

	// \u0026 in url needs to be replaced by &
	bytes, err := json.Marshal(presignReq)
	if err != nil {
		return GetBadResponse(err), err
	}

	return GetOKResponse(string(bytes)), nil
}
