package hollow

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

const (
	serversEndpoint                   = "servers"
	serverAttributesEndpoint          = "attributes"
	serverComponentsEndpoint          = "components"
	serverVersionedAttributesEndpoint = "versioned-attributes"
)

// ServerService provides the ability to interact with servers via Hollow
type ServerService interface {
	Create(context.Context, Server) (*uuid.UUID, *ServerResponse, error)
	Delete(context.Context, Server) (*ServerResponse, error)
	Get(context.Context, uuid.UUID) (*Server, *ServerResponse, error)
	List(context.Context, *ServerListParams) ([]Server, *ServerResponse, error)
	Update(context.Context, uuid.UUID, Server) (*ServerResponse, error)
	CreateAttributes(context.Context, uuid.UUID, Attributes) (*ServerResponse, error)
	DeleteAttributes(ctx context.Context, u uuid.UUID, ns string) (*ServerResponse, error)
	GetAttributes(context.Context, uuid.UUID, string) (*Attributes, *ServerResponse, error)
	ListAttributes(context.Context, uuid.UUID, *PaginationParams) ([]Attributes, *ServerResponse, error)
	UpdateAttributes(ctx context.Context, u uuid.UUID, ns string, data json.RawMessage) (*ServerResponse, error)
	ListComponents(context.Context, uuid.UUID, *PaginationParams) ([]ServerComponent, *ServerResponse, error)
	CreateVersionedAttributes(context.Context, uuid.UUID, VersionedAttributes) (*ServerResponse, error)
	GetVersionedAttributes(context.Context, uuid.UUID, string) ([]VersionedAttributes, *ServerResponse, error)
	ListVersionedAttributes(context.Context, uuid.UUID) ([]VersionedAttributes, *ServerResponse, error)
}

// ServerServiceClient implements ServerService
type ServerServiceClient struct {
	client *Client
}

// Create will attempt to create a server in Hollow and return the new server's UUID
func (c *ServerServiceClient) Create(ctx context.Context, srv Server) (*uuid.UUID, *ServerResponse, error) {
	resp, err := c.client.post(ctx, serversEndpoint, srv)
	if err != nil {
		return nil, nil, err
	}

	u, err := uuid.Parse(resp.Slug)
	if err != nil {
		return nil, resp, nil
	}

	return &u, resp, nil
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

// Update will to update a server with the new values passed in
func (c *ServerServiceClient) Update(ctx context.Context, srvUUID uuid.UUID, srv Server) (*ServerResponse, error) {
	path := fmt.Sprintf("%s/%s", serversEndpoint, srvUUID)
	return c.client.put(ctx, path, srv)
}

// CreateAttributes will to create the given attributes for a given server
func (c *ServerServiceClient) CreateAttributes(ctx context.Context, srvUUID uuid.UUID, attr Attributes) (*ServerResponse, error) {
	path := fmt.Sprintf("%s/%s/%s", serversEndpoint, srvUUID, serverAttributesEndpoint)
	return c.client.post(ctx, path, attr)
}

// GetAttributes will get all the attributes in a namespace for a given server
func (c *ServerServiceClient) GetAttributes(ctx context.Context, srvUUID uuid.UUID, ns string) (*Attributes, *ServerResponse, error) {
	attrs := &Attributes{}
	r := ServerResponse{Record: attrs}

	path := fmt.Sprintf("%s/%s/%s/%s", serversEndpoint, srvUUID, serverAttributesEndpoint, ns)
	if err := c.client.get(ctx, path, &r); err != nil {
		return nil, nil, err
	}

	return attrs, &r, nil
}

// DeleteAttributes will attempt to delete attributes by server uuid and namespace return an error on failure
func (c *ServerServiceClient) DeleteAttributes(ctx context.Context, srvUUID uuid.UUID, ns string) (*ServerResponse, error) {
	path := fmt.Sprintf("%s/%s/%s/%s", serversEndpoint, srvUUID, serverAttributesEndpoint, ns)
	return c.client.delete(ctx, path)
}

// ListAttributes will get all the attributes for a given server
func (c *ServerServiceClient) ListAttributes(ctx context.Context, srvUUID uuid.UUID, params *PaginationParams) ([]Attributes, *ServerResponse, error) {
	attrs := &[]Attributes{}
	r := ServerResponse{Records: attrs}

	path := fmt.Sprintf("%s/%s/%s", serversEndpoint, srvUUID, serverAttributesEndpoint)
	if err := c.client.list(ctx, path, params, &r); err != nil {
		return nil, nil, err
	}

	return *attrs, &r, nil
}

// UpdateAttributes will to update the data stored in a given namespace for a given server
func (c *ServerServiceClient) UpdateAttributes(ctx context.Context, srvUUID uuid.UUID, ns string, data json.RawMessage) (*ServerResponse, error) {
	path := fmt.Sprintf("%s/%s/%s/%s", serversEndpoint, srvUUID, serverAttributesEndpoint, ns)
	return c.client.put(ctx, path, Attributes{Data: data})
}

// ListComponents will get all the components for a given server
func (c *ServerServiceClient) ListComponents(ctx context.Context, srvUUID uuid.UUID, params *PaginationParams) ([]ServerComponent, *ServerResponse, error) {
	sc := &[]ServerComponent{}
	r := ServerResponse{Records: sc}

	path := fmt.Sprintf("%s/%s/%s", serversEndpoint, srvUUID, serverComponentsEndpoint)
	if err := c.client.list(ctx, path, params, &r); err != nil {
		return nil, nil, err
	}

	return *sc, &r, nil
}

// CreateVersionedAttributes will create a new versioned attribute for a given server
func (c *ServerServiceClient) CreateVersionedAttributes(ctx context.Context, srvUUID uuid.UUID, va VersionedAttributes) (*ServerResponse, error) {
	path := fmt.Sprintf("%s/%s/%s", serversEndpoint, srvUUID, serverVersionedAttributesEndpoint)

	return c.client.post(ctx, path, va)
}

// GetVersionedAttributes will return all the versioned attributes for a given server
func (c *ServerServiceClient) GetVersionedAttributes(ctx context.Context, srvUUID uuid.UUID, ns string) ([]VersionedAttributes, *ServerResponse, error) {
	path := fmt.Sprintf("%s/%s/%s/%s", serversEndpoint, srvUUID, serverVersionedAttributesEndpoint, ns)
	val := &[]VersionedAttributes{}
	r := ServerResponse{Records: val}

	if err := c.client.list(ctx, path, nil, &r); err != nil {
		return nil, nil, err
	}

	return *val, &r, nil
}

// ListVersionedAttributes will return all the versioned attributes for a given server
func (c *ServerServiceClient) ListVersionedAttributes(ctx context.Context, srvUUID uuid.UUID) ([]VersionedAttributes, *ServerResponse, error) {
	path := fmt.Sprintf("%s/%s/%s", serversEndpoint, srvUUID, serverVersionedAttributesEndpoint)
	val := &[]VersionedAttributes{}
	r := ServerResponse{Records: val}

	if err := c.client.list(ctx, path, nil, &r); err != nil {
		return nil, nil, err
	}

	return *val, &r, nil
}
