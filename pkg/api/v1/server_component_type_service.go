package fleetdbapi

import (
	"context"
)

const (
	serverComponentTypeEndpoint = "server-component-types"
)

// CreateServerComponentType will attempt to create a server component type in Hollow
func (c *Client) CreateServerComponentType(ctx context.Context, t ServerComponentType) (*ServerResponse, error) {
	return c.post(ctx, serverComponentTypeEndpoint, t)
}

// ListServerComponentTypes will return the server component types with optional params
func (c *Client) ListServerComponentTypes(ctx context.Context, params *ServerComponentTypeListParams) (ServerComponentTypeSlice, *ServerResponse, error) {
	cts := &ServerComponentTypeSlice{}
	resp := ServerResponse{Records: cts}

	if err := c.list(ctx, serverComponentTypeEndpoint, params, &resp); err != nil {
		return nil, nil, err
	}

	return *cts, &resp, nil
}
