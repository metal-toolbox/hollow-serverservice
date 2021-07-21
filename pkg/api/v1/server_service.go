package hollow

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const (
	serversEndpoint                   = "servers"
	serverVersionedAttributesEndpoint = "versioned-attributes"
)

// ServerService provides the ability to interact with servers via Hollow
type ServerService interface {
	Create(context.Context, Server) (*uuid.UUID, error)
	Delete(context.Context, Server) error
	Get(context.Context, uuid.UUID) (*Server, error)
	List(context.Context, *ServerListParams) ([]Server, error)
	GetVersionedAttributes(context.Context, uuid.UUID) ([]VersionedAttributes, error)
	CreateVersionedAttributes(context.Context, uuid.UUID, VersionedAttributes) (*uuid.UUID, error)
}

// ServerServiceClient implements ServerService
type ServerServiceClient struct {
	client *Client
}

// Create will attempt to create a server in Hollow and return the new server's UUID
func (c *ServerServiceClient) Create(ctx context.Context, srv Server) (*uuid.UUID, error) {
	return c.client.post(ctx, serversEndpoint, srv)
}

// Delete will attempt to delete a server in Hollow and return an error on failure
func (c *ServerServiceClient) Delete(ctx context.Context, srv Server) error {
	return c.client.delete(ctx, fmt.Sprintf("%s/%s", serversEndpoint, srv.UUID))
}

// Get will return a server by it's UUID
func (c *ServerServiceClient) Get(ctx context.Context, srvUUID uuid.UUID) (*Server, error) {
	path := fmt.Sprintf("%s/%s", serversEndpoint, srvUUID)

	var srv Server
	if err := c.client.get(ctx, path, &srv); err != nil {
		return nil, err
	}

	return &srv, nil
}

// List will return all servers with optional params to filter the results
func (c *ServerServiceClient) List(ctx context.Context, params *ServerListParams) ([]Server, error) {
	var srv []Server
	if err := c.client.list(ctx, serversEndpoint, params, &srv); err != nil {
		return nil, err
	}

	return srv, nil
}

// GetVersionedAttributes will return all the versioned attributes for a given server
func (c *ServerServiceClient) GetVersionedAttributes(ctx context.Context, srvUUID uuid.UUID) ([]VersionedAttributes, error) {
	path := fmt.Sprintf("%s/%s/%s", serversEndpoint, srvUUID, serverVersionedAttributesEndpoint)

	var val []VersionedAttributes
	if err := c.client.list(ctx, path, nil, &val); err != nil {
		return nil, err
	}

	return val, nil
}

// CreateVersionedAttributes will create a new versioned attribute for a given server
func (c *ServerServiceClient) CreateVersionedAttributes(ctx context.Context, srvUUID uuid.UUID, va VersionedAttributes) (*uuid.UUID, error) {
	path := fmt.Sprintf("%s/%s/%s", serversEndpoint, srvUUID, serverVersionedAttributesEndpoint)
	return c.client.put(ctx, path, va)
}
