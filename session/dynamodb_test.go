package session_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/borgoat/farmfa/session"
	"github.com/testcontainers/testcontainers-go"
)

const (
	DynamoTestTableName = "test-table"
)

// Locally launch a Dynalite container, create a test table in it, and return a compatible DynamoDB client for tests.
func setupDynalite(ctx context.Context) (*dynamodb.Client, error) {
	req := testcontainers.ContainerRequest{
		Image:        "docker.io/amazon/dynamodb-local",
		ExposedPorts: []string{"8000/tcp"},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "8000")
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	ddbClient := dynamodb.New(
		dynamodb.Options{
			EndpointResolver: dynamodb.EndpointResolverFunc(func(region string, options dynamodb.EndpointResolverOptions) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           fmt.Sprintf("http://%s:%d", hostIP, mappedPort.Int()),
					PartitionID:   "aws",
					SigningRegion: "us-east-1",
				}, nil
			}),
			Credentials: credentials.NewStaticCredentialsProvider("dummy", "dummy", ""),
		})

	err = createTestTable(ctx, ddbClient)
	if err != nil {
		return nil, err
	}

	return ddbClient, err
}

func createTestTable(ctx context.Context, client *dynamodb.Client) error {
	_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(DynamoTestTableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       types.KeyTypeRange,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	})
	if err != nil {
		return err
	}

	w := dynamodb.NewTableExistsWaiter(client)
	err = w.Wait(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(DynamoTestTableName)}, time.Minute)
	if err != nil {
		return err
	}

	return nil
}

func getDynamoDbStore() session.Store {
	ctx := context.TODO()

	ddbClient, err := setupDynalite(ctx)
	if err != nil {
		panic(err)
	}

	dynamoSession := session.NewDynamoDbStore(ctx, ddbClient, DynamoTestTableName)

	return dynamoSession
}

var testDynamoStore = getDynamoDbStore()

func TestDynamoDb_CreateSession(t *testing.T) {
	genericOracleCreateSesssion(t, testDynamoStore)
}

func TestDynamoDb_AddToc_valid(t *testing.T) {
	genericOracleAddToc_valid(t, testDynamoStore)
}

func TestDynamoDb_AddToc_empty(t *testing.T) {
	genericOracleAddToc_empty(t, testDynamoStore)
}

func TestDynamoDb_AddToc_notEncrypted(t *testing.T) {
	genericOracleAddToc_notEncrypted(t, testDynamoStore)
}

func TestDynamoDb_AddToc_alreadyExists(t *testing.T) {
	genericOracleAddToc_alreadyExists(t, testDynamoStore)
}

func TestDynamoDb_GenerateTOTP(t *testing.T) {
	genericOracleGenerateTOTP(t, testDynamoStore)
}
