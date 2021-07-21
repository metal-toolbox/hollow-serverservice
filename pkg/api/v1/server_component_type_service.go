package hollow

import (
	"context"

	"github.com/google/uuid"
)

const (
	serverComponentTypeEndpoint = "server-component-types"
)

// ServerComponentTypeService provides the ability to interact with server component types via Hollow
type ServerComponentTypeService interface {
	Create(context.Context, ServerComponentType) (*uuid.UUID, error)
	List(context.Context, *ServerComponentTypeListParams) ([]ServerComponentType, error)
}

// ServerComponentTypeServiceClient implements ServerComponentTypeService
type ServerComponentTypeServiceClient struct {
	client *Client
}

// Create will attempt to create a server component type in Hollow
func (c *ServerComponentTypeServiceClient) Create(ctx context.Context, t ServerComponentType) (*uuid.UUID, error) {
	return c.client.post(ctx, serverComponentTypeEndpoint, t)
}

// List will return the server component types with optional params
func (c *ServerComponentTypeServiceClient) List(ctx context.Context, params *ServerComponentTypeListParams) ([]ServerComponentType, error) {
	var ct []ServerComponentType
	if err := c.client.list(ctx, serverComponentTypeEndpoint, params, &ct); err != nil {
		return nil, err
	}

	return ct, nil
}
