package ddb

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

/*
Package client provides a DynamoDB implementation of the resource leasing system.

Table Design:
This implementation uses a single DynamoDB table with the following structure:
- Primary Key: PK (partition key) and SK (sort key)
- GSI1: Index for querying leases by user
- GSI2: Index for querying leases by resource
- GSI3: Index for querying users by email

Key Patterns:
- User records:     PK=USER#{id}, SK=METADATA, GSI3PK=EMAIL#{email}
- Resource records: PK=RESOURCE#{id}, SK=METADATA
- Lease records:    PK=LEASE#{id}, SK=METADATA, GSI1PK=USER#{id}, GSI2PK=RESOURCE#{id}

Table creation command:

	aws dynamodb create-table \
		--table-name YourTableName \
		--attribute-definitions \
				AttributeName=PK,AttributeType=S \
				AttributeName=SK,AttributeType=S \
				AttributeName=GSI1PK,AttributeType=S \
				AttributeName=GSI1SK,AttributeType=S \
				AttributeName=GSI2PK,AttributeType=S \
				AttributeName=GSI2SK,AttributeType=S \
        AttributeName=GSI3PK,AttributeType=S \
		--key-schema \
				AttributeName=PK,KeyType=HASH \
				AttributeName=SK,KeyType=RANGE \
		--global-secondary-indexes \
				"[
						{
								\"IndexName\": \"GSI1\",
								\"KeySchema\": [
										{\"AttributeName\":\"GSI1PK\",\"KeyType\":\"HASH\"},
										{\"AttributeName\":\"GSI1SK\",\"KeyType\":\"RANGE\"}
								],
								\"Projection\": {
										\"ProjectionType\":\"ALL\"
								},
								\"ProvisionedThroughput\": {
										\"ReadCapacityUnits\": 5,
										\"WriteCapacityUnits\": 5
								}
						},
						{
								\"IndexName\": \"GSI2\",
								\"KeySchema\": [
										{\"AttributeName\":\"GSI2PK\",\"KeyType\":\"HASH\"},
										{\"AttributeName\":\"GSI2SK\",\"KeyType\":\"RANGE\"}
								],
								\"Projection\": {
										\"ProjectionType\":\"ALL\"
								},
								\"ProvisionedThroughput\": {
										\"ReadCapacityUnits\": 5,
										\"WriteCapacityUnits\": 5
								}
						},
            {
                \"IndexName\": \"GSI3\",
                \"KeySchema\": [
                    {\"AttributeName\":\"GSI3PK\",\"KeyType\":\"HASH\"}
                ],
                \"Projection\": {
                    \"ProjectionType\":\"ALL\"
                },
                \"ProvisionedThroughput\": {
                    \"ReadCapacityUnits\": 5,
                    \"WriteCapacityUnits\": 5
                }
            }
				]" \
		--provisioned-throughput \
				ReadCapacityUnits=5,WriteCapacityUnits=5

	aws dynamodb update-time-to-live \
    --table-name YourTableName \
    --time-to-live-specification "Enabled=true, AttributeName=ttl"
*/

// User represents a user in DynamoDB
type User struct {
	ID          uuid.UUID `dynamodbav:"id"`
	DisplayName string    `dynamodbav:"display_name"`
	Email       string    `dynamodbav:"email"`
}

// Resource represents a resource in DynamoDB
type Resource struct {
	ID   uuid.UUID `dynamodbav:"id"`
	Name string    `dynamodbav:"name"`
}

// Lease represents a lease in DynamoDB
type Lease struct {
	ID       uuid.UUID     `dynamodbav:"id"`
	UserId   uuid.UUID     `dynamodbav:"user_id"`
	ResId    uuid.UUID     `dynamodbav:"resource_id"`
	Reason   string        `dynamodbav:"reason"`
	Duration time.Duration `dynamodbav:"duration"`
	TTL      int64         `dynamodbav:"ttl"` // DynamoDB TTL field
}

type DynamoDBClient struct {
	client *dynamodb.Client
	table  string
}

func NewDynamoDBClient(ctx context.Context, tableName string) (*DynamoDBClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return &DynamoDBClient{
		client: client,
		table:  tableName,
	}, nil
}

var emailRegex = regexp.MustCompile(`[^@]+@[^@]+`)

func (c *DynamoDBClient) CreateUser(ctx context.Context, displayName, email string) (*User, error) {
	if displayName == "" {
		return nil, fmt.Errorf("display name cannot be empty")
	}
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	if !emailRegex.MatchString(email) {
		return nil, fmt.Errorf("invalid email format")
	}

	user := &User{
		ID:          uuid.New(),
		DisplayName: displayName,
		Email:       email,
	}

	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user: %w", err)
	}

	av["PK"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", user.ID.String())}
	av["SK"] = &types.AttributeValueMemberS{Value: "METADATA"}
	av["GSI3PK"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("EMAIL#%s", email)}

	_, err = c.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(c.table),
		Item:      av,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (c *DynamoDBClient) CreateResource(ctx context.Context, name string) (*Resource, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	resource := &Resource{
		ID:   uuid.New(),
		Name: name,
	}

	av, err := attributevalue.MarshalMap(resource)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource: %w", err)
	}

	av["PK"] = &types.AttributeValueMemberS{Value: "RESOURCE#" + resource.ID.String()}
	av["SK"] = &types.AttributeValueMemberS{Value: "METADATA"}

	_, err = c.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(c.table),
		Item:      av,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	return resource, nil
}

func (c *DynamoDBClient) CreateLease(ctx context.Context, userID, resourceID uuid.UUID, reason string, duration time.Duration) (*Lease, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user ID cannot be empty")
	}
	if resourceID == uuid.Nil {
		return nil, fmt.Errorf("resource ID cannot be empty")
	}
	if reason == "" {
		return nil, fmt.Errorf("reason cannot be empty")
	}
	if duration <= 0 {
		return nil, fmt.Errorf("duration must be positive")
	}

	now := time.Now()
	lease := &Lease{
		ID:       uuid.New(),
		UserId:   userID,
		ResId:    resourceID,
		Reason:   reason,
		Duration: duration,
		TTL:      now.Add(duration).Unix(), // Set TTL to creation time + duration
	}

	av, err := attributevalue.MarshalMap(lease)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal lease: %w", err)
	}

	av["PK"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("LEASE#%s", lease.ID.String())}
	av["SK"] = &types.AttributeValueMemberS{Value: "METADATA"}
	av["GSI1PK"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID.String())}
	av["GSI1SK"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("LEASE#%s", lease.ID.String())}
	av["GSI2PK"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("RESOURCE#%s", resourceID.String())}
	av["GSI2SK"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("LEASE#%s", lease.ID.String())}

	_, err = c.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(c.table),
		Item:      av,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create lease: %w", err)
	}

	return lease, nil
}

func (c *DynamoDBClient) DeleteLease(ctx context.Context, leaseID uuid.UUID) error {
	if leaseID == uuid.Nil {
		return fmt.Errorf("lease ID cannot be empty")
	}

	_, err := c.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(c.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("LEASE#%s", leaseID.String())},
			"SK": &types.AttributeValueMemberS{Value: "METADATA"},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to delete lease: %w", err)
	}

	return nil
}

func (c *DynamoDBClient) GetLeasesForUser(ctx context.Context, userID uuid.UUID) ([]*Lease, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	result, err := c.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(c.table),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("GSI1PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID.String())},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query leases: %w", err)
	}

	leases := make([]*Lease, 0)
	for _, item := range result.Items {
		var leaseData map[string]types.AttributeValue
		if v, ok := item["LeaseData"]; ok {
			if m, ok := v.(*types.AttributeValueMemberM); ok {
				leaseData = m.Value
			}
		}

		var lease Lease
		err = attributevalue.UnmarshalMap(leaseData, &lease)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal lease: %w", err)
		}
		leases = append(leases, &lease)
	}

	return leases, nil
}

func (c *DynamoDBClient) GetLeasesForResource(ctx context.Context, resourceID uuid.UUID) ([]*Lease, error) {
	if resourceID == uuid.Nil {
		return nil, fmt.Errorf("resource ID cannot be empty")
	}

	result, err := c.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(c.table),
		IndexName:              aws.String("GSI2"),
		KeyConditionExpression: aws.String("GSI2PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("RESOURCE#%s", resourceID.String())},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query leases: %w", err)
	}

	leases := make([]*Lease, 0)
	for _, item := range result.Items {
		var leaseData map[string]types.AttributeValue
		if v, ok := item["LeaseData"]; ok {
			if m, ok := v.(*types.AttributeValueMemberM); ok {
				leaseData = m.Value
			}
		}

		var lease Lease
		err = attributevalue.UnmarshalMap(leaseData, &lease)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal lease: %w", err)
		}
		leases = append(leases, &lease)
	}

	return leases, nil
}

func (c *DynamoDBClient) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	if !emailRegex.MatchString(email) {
		return nil, fmt.Errorf("invalid email format")
	}

	// In DDB you can't do a GetItem on a GSI, so we have to use Query
	result, err := c.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(c.table),
		IndexName:              aws.String("GSI3"),
		KeyConditionExpression: aws.String("GSI3PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("EMAIL#%s", email)},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	var user User
	err = attributevalue.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}
