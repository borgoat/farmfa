package session

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	ddb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"golang.org/x/crypto/blake2b"

	"github.com/borgoat/farmfa/api"
)

type dynamoItem struct {
	PK         string
	SK         string
	RecordType string
}

type dynamoSession struct {
	dynamoItem
	*api.Session
}

type dynamoEncryptedString struct {
	dynamoItem
	EncryptedValue string
}

type dynamoEncryptedBytes struct {
	dynamoItem
	EncryptedValue []byte
}

type DDBGetItemApi interface {
	GetItem(ctx context.Context, params *ddb.GetItemInput, optFns ...func(*ddb.Options)) (*ddb.GetItemOutput, error)
}

type DDBQueryApi interface {
	Query(ctx context.Context, params *ddb.QueryInput, optFns ...func(*ddb.Options)) (*ddb.QueryOutput, error)
}

type DDBPutItemApi interface {
	PutItem(ctx context.Context, params *ddb.PutItemInput, optFns ...func(*ddb.Options)) (*ddb.PutItemOutput, error)
}

type DDBBatchWriteItemApi interface {
	BatchWriteItem(ctx context.Context, params *ddb.BatchWriteItemInput, optFns ...func(*ddb.Options)) (*ddb.BatchWriteItemOutput, error)
}

type DDBTransactWriteItemsApi interface {
	TransactWriteItems(ctx context.Context, params *ddb.TransactWriteItemsInput, optFns ...func(*ddb.Options)) (*ddb.TransactWriteItemsOutput, error)
}

type DDBClient interface {
	DDBGetItemApi
	DDBQueryApi
	DDBPutItemApi
	DDBBatchWriteItemApi
	DDBTransactWriteItemsApi
}

type DynamoDbStore struct {
	ctx    context.Context
	client DDBClient
	table  string
}

func NewDynamoDbStore(ctx context.Context, client DDBClient, table string) Store {
	return &DynamoDbStore{
		ctx:    ctx,
		client: client,
		table:  table,
	}
}

func (d *DynamoDbStore) encryptedTocItem(sessionId, encryptedToc string) (map[string]types.AttributeValue, error) {
	h, err := blake2b.New256(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create hash function: %w", err)
	}
	h.Write([]byte(encryptedToc))

	tocItem, err := attributevalue.MarshalMap(&dynamoEncryptedString{
		dynamoItem{
			PK:         fmt.Sprintf("SESSION#%s", sessionId),
			SK:         fmt.Sprintf("TOC#%s", base64.StdEncoding.EncodeToString(h.Sum(nil))),
			RecordType: "Toc",
		},
		encryptedToc,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal toc: %w", err)
	}
	return tocItem, nil
}

func (d *DynamoDbStore) getTablePtr() *string {
	return aws.String(d.table)
}

func (d *DynamoDbStore) CreateSession(session *api.Session, encryptedTEK []byte, encryptedTocZero string) error {
	sessionItem, err := attributevalue.MarshalMap(&dynamoSession{
		dynamoItem{
			PK:         fmt.Sprintf("SESSION#%s", session.Id),
			SK:         "SESSION",
			RecordType: "Session",
		},
		session,
	})
	if err != nil {
		return fmt.Errorf("failed to create session object: %w", err)
	}

	tocZeroItem, err := d.encryptedTocItem(session.Id, encryptedTocZero)
	if err != nil {
		return fmt.Errorf("failed to prepare toc zero: %w", err)
	}

	tekItem, err := attributevalue.MarshalMap(&dynamoEncryptedBytes{
		dynamoItem{
			PK:         fmt.Sprintf("SESSION#%s", session.Id),
			SK:         "TEK",
			RecordType: "TEK",
		},
		encryptedTEK,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal TEK: %w", err)
	}

	_, err = d.client.TransactWriteItems(d.ctx, &ddb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					Item:                     sessionItem,
					TableName:                d.getTablePtr(),
					ExpressionAttributeNames: map[string]string{"#pk": "PK", "#sk": "SK"},
					ConditionExpression:      aws.String("attribute_not_exists(#pk) AND attribute_not_exists(#sk)"),
				},
			},
			{
				Put: &types.Put{
					Item:      tocZeroItem,
					TableName: d.getTablePtr(),
				},
			},
			{
				Put: &types.Put{
					Item:      tekItem,
					TableName: d.getTablePtr(),
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to write to dynamo: %w", err)
	}

	return nil
}

func (d *DynamoDbStore) GetSession(id string) (*api.Session, error) {
	// TODO Check ttl
	resp, err := d.client.GetItem(d.ctx, &ddb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("SESSION#%s", id)},
			"SK": &types.AttributeValueMemberS{Value: "SESSION"},
		},
		TableName:      d.getTablePtr(),
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve session: %w", err)
	}

	var sessionItem dynamoSession
	err = attributevalue.UnmarshalMap(resp.Item, &sessionItem)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return sessionItem.Session, nil
}

func (d *DynamoDbStore) AddEncryptedToc(id string, encryptedToc string) error {
	encryptedTocItem, err := d.encryptedTocItem(id, encryptedToc)
	if err != nil {
		return fmt.Errorf("failed to prepare toc: %w", err)
	}

	_, err = d.client.PutItem(d.ctx, &ddb.PutItemInput{
		Item:                     encryptedTocItem,
		TableName:                d.getTablePtr(),
		ConditionExpression:      aws.String("attribute_not_exists(#pk) AND attribute_not_exists(#sk)"),
		ExpressionAttributeNames: map[string]string{"#pk": "PK", "#sk": "SK"},
	})
	if err != nil {
		return ErrTocAlreadyExists
	}

	return nil
}

func (d *DynamoDbStore) GetEncryptedTocs(id string) ([]string, error) {
	resp, err := d.client.Query(d.ctx, &ddb.QueryInput{
		TableName:                d.getTablePtr(),
		KeyConditionExpression:   aws.String("#pk = :pk AND begins_with(#sk, :prefix)"),
		ConsistentRead:           aws.Bool(true),
		ExpressionAttributeNames: map[string]string{"#pk": "PK", "#sk": "SK"},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: fmt.Sprintf("SESSION#%s", id)},
			":prefix": &types.AttributeValueMemberS{Value: "TOC#"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve encrypted tocs: %w", err)
	}

	tocs := make([]string, len(resp.Items))

	for i, item := range resp.Items {
		var tocItem dynamoEncryptedString
		err = attributevalue.UnmarshalMap(item, &tocItem)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal toc: %w", err)
		}
		tocs[i] = tocItem.EncryptedValue
	}

	return tocs, nil
}

func (d *DynamoDbStore) GetTEK(id string) ([]byte, error) {
	resp, err := d.client.GetItem(d.ctx, &ddb.GetItemInput{
		TableName:      d.getTablePtr(),
		ConsistentRead: aws.Bool(true),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("SESSION#%s", id)},
			"SK": &types.AttributeValueMemberS{Value: "TEK"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tek: %w", err)
	}

	var tekItem dynamoEncryptedBytes
	err = attributevalue.UnmarshalMap(resp.Item, &tekItem)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarhsal tek: %w", err)
	}

	return tekItem.EncryptedValue, nil
}

func (d *DynamoDbStore) GarbageCollect(shouldDelete func(session *api.Session) bool) {
	panic("implement me")
}
