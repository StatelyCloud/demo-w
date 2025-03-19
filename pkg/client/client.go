package client

import (
	"context"
	"time"

	"github.com/StatelyCloud/demo-w/pkg/schema"
	"github.com/StatelyCloud/go-sdk/stately"
	"github.com/google/uuid"
)

type Client struct {
	client stately.Client
}

func NewClient(ctx context.Context, storeID uint64) (*Client, error) {
	statelyClient, err := schema.NewClient(ctx, storeID, &stately.Options{
		NoAuth:   true,
		Endpoint: "http://localhost:3030",
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		statelyClient,
	}, nil
}

func (c *Client) CreateUser(ctx context.Context, displayName, email string) (*schema.User, error) {
	item, err := c.client.Put(ctx, &schema.User{
		DisplayName: displayName,
		Email:       email,
	})
	if err != nil {
		return nil, err
	}
	return item.(*schema.User), nil
}

func (c *Client) CreateResource(ctx context.Context, name string) (*schema.Resource, error) {
	item, err := c.client.Put(ctx, &schema.Resource{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return item.(*schema.Resource), nil
}

func (c *Client) CreateLease(ctx context.Context, userID, resourceID uuid.UUID, reason string, duration time.Duration) (*schema.Lease, error) {
	item, err := c.client.Put(ctx, &schema.Lease{
		UserId:   userID,
		ResId:    resourceID,
		Reason:   reason,
		Duration: duration,
	})
	if err != nil {
		return nil, err
	}
	return item.(*schema.Lease), nil
}

func (c *Client) DeleteLease(ctx context.Context, leaseID uuid.UUID) error {
	return c.client.Delete(ctx, "/lease-"+stately.ToKeyID(leaseID[:]))
}

func (c *Client) GetLeasesForUser(ctx context.Context, userID uuid.UUID) ([]*schema.Lease, error) {
	resp, err := c.client.BeginList(ctx, "/user-"+stately.ToKeyID(userID[:])+"/res")
	if err != nil {
		return nil, err
	}
	var leases []*schema.Lease
	for resp.Next() {
		if lease, ok := resp.Value().(*schema.Lease); ok {
			leases = append(leases, lease)
		}
	}
	return leases, nil
}

func (c *Client) GetLeasesForResource(ctx context.Context, resourceID uuid.UUID) ([]*schema.Lease, error) {
	resp, err := c.client.BeginList(ctx, "/res-"+stately.ToKeyID(resourceID[:])+"/lease")
	if err != nil {
		return nil, err
	}
	var leases []*schema.Lease
	for resp.Next() {
		if lease, ok := resp.Value().(*schema.Lease); ok {
			leases = append(leases, lease)
		}
	}
	return leases, nil
}

func (c *Client) GetUserByEmail(ctx context.Context, email string) (*schema.User, error) {
	user, err := c.client.Get(ctx, "/user_email-"+stately.ToKeyID(email))
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user.(*schema.User), nil
	}
	return nil, nil
}
