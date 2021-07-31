package hollow

import (
	"context"
)

const (
	serverComponentTypeEndpoint = "server-component-types"
)

// ServerComponentTypeService provides the ability to interact with server component types via Hollow
type ServerComponentTypeService interface {
	Create(context.Context, ServerComponentType) (*ServerResponse, error)
	List(context.Context, *ServerComponentTypeListParams) ([]ServerComponentType, *ServerResponse, error)
}

// ServerComponentTypeServiceClient implements ServerComponentTypeService
type ServerComponentTypeServiceClient struct {
	client *Client
}

// Create will attempt to create a server component type in Hollow
func (c *ServerComponentTypeServiceClient) Create(ctx context.Context, t ServerComponentType) (*ServerResponse, error) {
	return c.client.post(ctx, serverComponentTypeEndpoint, t)
}

// List will return the server component types with optional params
func (c *ServerComponentTypeServiceClient) List(ctx context.Context, params *ServerComponentTypeListParams) ([]ServerComponentType, *ServerResponse, error) {
	cts := &[]ServerComponentType{}
	resp := ServerResponse{Records: cts}

	if err := c.client.list(ctx, serverComponentTypeEndpoint, params, &resp); err != nil {
		return nil, nil, err
	}

	return *cts, &resp, nil
}
