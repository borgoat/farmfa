package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	ddb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"

	"github.com/borgoat/farmfa/api"
	"github.com/borgoat/farmfa/server"
	"github.com/borgoat/farmfa/session"
)

var echoLambda = func() *echoadapter.EchoLambdaV2 {
	e := echo.New()
	apiObj, err := api.GetSwagger()
	if err != nil {
		panic(fmt.Errorf("error loading OpenAPI spec: %w", err))
	}
	e.Use(middleware.OapiRequestValidator(apiObj))

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Errorf("failed to get AWS SDK config: %w", err))
	}
	ddbClient := ddb.NewFromConfig(cfg)

	store := session.NewDynamoDbStore(context.TODO(), ddbClient, os.Getenv("FARMFA_DYNAMODB_TABLE"))
	oracle := session.NewOracle(store)
	s := server.New(oracle)
	api.RegisterHandlers(e, s)

	return echoadapter.NewV2(e)
}()

func HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return echoLambda.ProxyWithContext(ctx, event)
}

func main() {
	lambda.Start(HandleRequest)
}
