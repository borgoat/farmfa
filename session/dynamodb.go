package session

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	ddb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"golang.org/x/crypto/blake2b"

	"github.com/borgoat/farmfa/api"
)

var (
	exprKeyPk = expression.Key("PK")
	exprKeySk = expression.Key("SK")
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

func prepareTocRecord(sessionId, encryptedToc string) (map[string]types.AttributeValue, error) {
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

func checkUniquePrimaryKey() (expression.Expression, error) {
	noPk := expression.Name("PK").AttributeNotExists()
	noSk := expression.Name("SK").AttributeNotExists()
	return expression.NewBuilder().WithCondition(noPk.And(noSk)).Build()
}

func sessionPk(id string) string {
	return fmt.Sprintf("SESSION#%s", id)
}

func (d *DynamoDbStore) getTableName() *string {
	return aws.String(d.table)
}

func (d *DynamoDbStore) CreateSession(session *api.Session, encryptedTEK []byte, encryptedTocZero string) error {
	sessionItem, err := attributevalue.MarshalMap(&dynamoSession{
		dynamoItem{
			PK:         sessionPk(session.Id),
			SK:         "SESSION",
			RecordType: "Session",
		},
		session,
	})
	if err != nil {
		return fmt.Errorf("failed to create session object: %w", err)
	}

	tocZeroItem, err := prepareTocRecord(session.Id, encryptedTocZero)
	if err != nil {
		return fmt.Errorf("failed to prepare toc zero: %w", err)
	}

	tekItem, err := attributevalue.MarshalMap(&dynamoEncryptedBytes{
		dynamoItem{
			PK:         sessionPk(session.Id),
			SK:         "TEK",
			RecordType: "TEK",
		},
		encryptedTEK,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal TEK: %w", err)
	}

	expr, err := checkUniquePrimaryKey()
	if err != nil {
		return fmt.Errorf("failed to build condition: %w", err)
	}

	_, err = d.client.TransactWriteItems(d.ctx, &ddb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					Item:                     sessionItem,
					TableName:                d.getTableName(),
					ExpressionAttributeNames: expr.Names(),
					ConditionExpression:      expr.Condition(),
				},
			},
			{
				Put: &types.Put{
					Item:      tocZeroItem,
					TableName: d.getTableName(),
				},
			},
			{
				Put: &types.Put{
					Item:      tekItem,
					TableName: d.getTableName(),
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
			"PK": &types.AttributeValueMemberS{Value: sessionPk(id)},
			"SK": &types.AttributeValueMemberS{Value: "SESSION"},
		},
		TableName:      d.getTableName(),
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
	encryptedTocItem, err := prepareTocRecord(id, encryptedToc)
	if err != nil {
		return fmt.Errorf("failed to prepare toc: %w", err)
	}

	expr, err := checkUniquePrimaryKey()
	if err != nil {
		return fmt.Errorf("failed to build condition: %w", err)
	}

	_, err = d.client.PutItem(d.ctx, &ddb.PutItemInput{
		Item:                     encryptedTocItem,
		TableName:                d.getTableName(),
		ConditionExpression:      expr.Condition(),
		ExpressionAttributeNames: expr.Names(),
	})
	if err != nil {
		return ErrTocAlreadyExists
	}

	return nil
}

func (d *DynamoDbStore) GetEncryptedTocs(id string) ([]string, error) {
	pkBySession := expression.KeyEqual(exprKeyPk, expression.Value(sessionPk(id)))
	skTypeToc := expression.KeyBeginsWith(exprKeySk, "TOC#")
	expr, err := expression.NewBuilder().WithKeyCondition(expression.KeyAnd(pkBySession, skTypeToc)).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build condition: %w", err)
	}

	resp, err := d.client.Query(d.ctx, &ddb.QueryInput{
		TableName:                 d.getTableName(),
		ConsistentRead:            aws.Bool(true),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
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
		TableName:      d.getTableName(),
		ConsistentRead: aws.Bool(true),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: sessionPk(id)},
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

func (d *DynamoDbStore) GarbageCollect(func(session *api.Session) bool) {
	// Probably shouldn't implement considering dynamo has TTL already
}
