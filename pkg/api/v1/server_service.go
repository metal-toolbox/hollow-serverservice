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
	Create(context.Context, Server) (*uuid.UUID, *ServerResponse, error)
	Delete(context.Context, Server) (*ServerResponse, error)
	Get(context.Context, uuid.UUID) (*Server, *ServerResponse, error)
	List(context.Context, *ServerListParams) ([]Server, *ServerResponse, error)
	GetVersionedAttributes(context.Context, uuid.UUID) ([]VersionedAttributes, *ServerResponse, error)
	CreateVersionedAttributes(context.Context, uuid.UUID, VersionedAttributes) (*uuid.UUID, *ServerResponse, error)
}

// ServerServiceClient implements ServerService
type ServerServiceClient struct {
	client *Client
}

// Create will attempt to create a server in Hollow and return the new server's UUID
func (c *ServerServiceClient) Create(ctx context.Context, srv Server) (*uuid.UUID, *ServerResponse, error) {
	return c.client.post(ctx, serversEndpoint, srv)
}

// Delete will attempt to delete a server in Hollow and return an error on failure
func (c *ServerServiceClient) Delete(ctx context.Context, srv Server) (*ServerResponse, error) {
	return c.client.delete(ctx, fmt.Sprintf("%s/%s", serversEndpoint, srv.UUID))
}

// Get will return a server by it's UUID
func (c *ServerServiceClient) Get(ctx context.Context, srvUUID uuid.UUID) (*Server, *ServerResponse, error) {
	path := fmt.Sprintf("%s/%s", serversEndpoint, srvUUID)
	srv := &Server{}
	r := ServerResponse{Record: srv}

	if err := c.client.get(ctx, path, &r); err != nil {
		return nil, nil, err
	}

	return srv, &r, nil
}

// List will return all servers with optional params to filter the results
func (c *ServerServiceClient) List(ctx context.Context, params *ServerListParams) ([]Server, *ServerResponse, error) {
	servers := &[]Server{}
	r := ServerResponse{Records: servers}

	if err := c.client.list(ctx, serversEndpoint, params, &r); err != nil {
		return nil, nil, err
	}

	return *servers, &r, nil
}

// GetVersionedAttributes will return all the versioned attributes for a given server
func (c *ServerServiceClient) GetVersionedAttributes(ctx context.Context, srvUUID uuid.UUID) ([]VersionedAttributes, *ServerResponse, error) {
	path := fmt.Sprintf("%s/%s/%s", serversEndpoint, srvUUID, serverVersionedAttributesEndpoint)
	val := &[]VersionedAttributes{}
	r := ServerResponse{Records: val}

	if err := c.client.list(ctx, path, nil, &r); err != nil {
		return nil, nil, err
	}

	return *val, &r, nil
}

// CreateVersionedAttributes will create a new versioned attribute for a given server
func (c *ServerServiceClient) CreateVersionedAttributes(ctx context.Context, srvUUID uuid.UUID, va VersionedAttributes) (*uuid.UUID, *ServerResponse, error) {
	path := fmt.Sprintf("%s/%s/%s", serversEndpoint, srvUUID, serverVersionedAttributesEndpoint)
	return c.client.put(ctx, path, va)
}
