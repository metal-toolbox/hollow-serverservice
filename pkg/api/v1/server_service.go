package serverservice

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
	serverComponentFirmwaresEndpoint  = "server-component-firmwares"
)

// ClientInterface provides an interface for the expected calls to interact with a server service api
type ClientInterface interface {
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
	CreateServerComponentFirmware(context.Context, ComponentFirmwareVersion) (*uuid.UUID, *ServerResponse, error)
	DeleteServerComponentFirmware(context.Context, ComponentFirmwareVersion) (*ServerResponse, error)
	GetServerComponentFirmware(context.Context, uuid.UUID) (*ComponentFirmwareVersion, *ServerResponse, error)
	ListServerComponentFirmware(context.Context, *ComponentFirmwareVersionListParams) ([]ComponentFirmwareVersion, *ServerResponse, error)
	UpdateServerComponentFirmware(context.Context, uuid.UUID, ComponentFirmwareVersion) (*ServerResponse, error)
}

// Create will attempt to create a server in Hollow and return the new server's UUID
func (c *Client) Create(ctx context.Context, srv Server) (*uuid.UUID, *ServerResponse, error) {
	resp, err := c.post(ctx, serversEndpoint, srv)
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
func (c *Client) Delete(ctx context.Context, srv Server) (*ServerResponse, error) {
	return c.delete(ctx, fmt.Sprintf("%s/%s", serversEndpoint, srv.UUID))
}

// Get will return a server by it's UUID
func (c *Client) Get(ctx context.Context, srvUUID uuid.UUID) (*Server, *ServerResponse, error) {
	path := fmt.Sprintf("%s/%s", serversEndpoint, srvUUID)
	srv := &Server{}
	r := ServerResponse{Record: srv}

	if err := c.get(ctx, path, &r); err != nil {
		return nil, nil, err
	}

	return srv, &r, nil
}

// List will return all servers with optional params to filter the results
func (c *Client) List(ctx context.Context, params *ServerListParams) ([]Server, *ServerResponse, error) {
	servers := &[]Server{}
	r := ServerResponse{Records: servers}

	if err := c.list(ctx, serversEndpoint, params, &r); err != nil {
		return nil, nil, err
	}

	return *servers, &r, nil
}

// Update will to update a server with the new values passed in
func (c *Client) Update(ctx context.Context, srvUUID uuid.UUID, srv Server) (*ServerResponse, error) {
	path := fmt.Sprintf("%s/%s", serversEndpoint, srvUUID)
	return c.put(ctx, path, srv)
}

// CreateAttributes will to create the given attributes for a given server
func (c *Client) CreateAttributes(ctx context.Context, srvUUID uuid.UUID, attr Attributes) (*ServerResponse, error) {
	path := fmt.Sprintf("%s/%s/%s", serversEndpoint, srvUUID, serverAttributesEndpoint)
	return c.post(ctx, path, attr)
}

// GetAttributes will get all the attributes in a namespace for a given server
func (c *Client) GetAttributes(ctx context.Context, srvUUID uuid.UUID, ns string) (*Attributes, *ServerResponse, error) {
	attrs := &Attributes{}
	r := ServerResponse{Record: attrs}

	path := fmt.Sprintf("%s/%s/%s/%s", serversEndpoint, srvUUID, serverAttributesEndpoint, ns)
	if err := c.get(ctx, path, &r); err != nil {
		return nil, nil, err
	}

	return attrs, &r, nil
}

// DeleteAttributes will attempt to delete attributes by server uuid and namespace return an error on failure
func (c *Client) DeleteAttributes(ctx context.Context, srvUUID uuid.UUID, ns string) (*ServerResponse, error) {
	path := fmt.Sprintf("%s/%s/%s/%s", serversEndpoint, srvUUID, serverAttributesEndpoint, ns)
	return c.delete(ctx, path)
}

// ListAttributes will get all the attributes for a given server
func (c *Client) ListAttributes(ctx context.Context, srvUUID uuid.UUID, params *PaginationParams) ([]Attributes, *ServerResponse, error) {
	attrs := &[]Attributes{}
	r := ServerResponse{Records: attrs}

	path := fmt.Sprintf("%s/%s/%s", serversEndpoint, srvUUID, serverAttributesEndpoint)
	if err := c.list(ctx, path, params, &r); err != nil {
		return nil, nil, err
	}

	return *attrs, &r, nil
}

// UpdateAttributes will to update the data stored in a given namespace for a given server
func (c *Client) UpdateAttributes(ctx context.Context, srvUUID uuid.UUID, ns string, data json.RawMessage) (*ServerResponse, error) {
	path := fmt.Sprintf("%s/%s/%s/%s", serversEndpoint, srvUUID, serverAttributesEndpoint, ns)
	return c.put(ctx, path, Attributes{Data: data})
}

// ListComponents will get all the components for a given server
func (c *Client) ListComponents(ctx context.Context, srvUUID uuid.UUID, params *PaginationParams) ([]ServerComponent, *ServerResponse, error) {
	sc := &[]ServerComponent{}
	r := ServerResponse{Records: sc}

	path := fmt.Sprintf("%s/%s/%s", serversEndpoint, srvUUID, serverComponentsEndpoint)
	if err := c.list(ctx, path, params, &r); err != nil {
		return nil, nil, err
	}

	return *sc, &r, nil
}

// CreateVersionedAttributes will create a new versioned attribute for a given server
func (c *Client) CreateVersionedAttributes(ctx context.Context, srvUUID uuid.UUID, va VersionedAttributes) (*ServerResponse, error) {
	path := fmt.Sprintf("%s/%s/%s", serversEndpoint, srvUUID, serverVersionedAttributesEndpoint)

	return c.post(ctx, path, va)
}

// GetVersionedAttributes will return all the versioned attributes for a given server
func (c *Client) GetVersionedAttributes(ctx context.Context, srvUUID uuid.UUID, ns string) ([]VersionedAttributes, *ServerResponse, error) {
	path := fmt.Sprintf("%s/%s/%s/%s", serversEndpoint, srvUUID, serverVersionedAttributesEndpoint, ns)
	val := &[]VersionedAttributes{}
	r := ServerResponse{Records: val}

	if err := c.list(ctx, path, nil, &r); err != nil {
		return nil, nil, err
	}

	return *val, &r, nil
}

// ListVersionedAttributes will return all the versioned attributes for a given server
func (c *Client) ListVersionedAttributes(ctx context.Context, srvUUID uuid.UUID) ([]VersionedAttributes, *ServerResponse, error) {
	path := fmt.Sprintf("%s/%s/%s", serversEndpoint, srvUUID, serverVersionedAttributesEndpoint)
	val := &[]VersionedAttributes{}
	r := ServerResponse{Records: val}

	if err := c.list(ctx, path, nil, &r); err != nil {
		return nil, nil, err
	}

	return *val, &r, nil
}

// CreateServerComponentFirmware will attempt to create a firmware in Hollow and return the firmware UUID
func (c *Client) CreateServerComponentFirmware(ctx context.Context, firmware ComponentFirmwareVersion) (*uuid.UUID, *ServerResponse, error) {
	resp, err := c.post(ctx, serverComponentFirmwaresEndpoint, firmware)
	if err != nil {
		return nil, nil, err
	}

	u, err := uuid.Parse(resp.Slug)
	if err != nil {
		return nil, resp, nil
	}

	return &u, resp, nil
}

// DeleteServerComponentFirmware will attempt to delete firmware and return an error on failure
func (c *Client) DeleteServerComponentFirmware(ctx context.Context, firmware ComponentFirmwareVersion) (*ServerResponse, error) {
	return c.delete(ctx, fmt.Sprintf("%s/%s", serverComponentFirmwaresEndpoint, firmware.UUID))
}

// GetServerComponentFirmware will return a firmware by its UUID
func (c *Client) GetServerComponentFirmware(ctx context.Context, fwUUID uuid.UUID) (*ComponentFirmwareVersion, *ServerResponse, error) {
	path := fmt.Sprintf("%s/%s", serverComponentFirmwaresEndpoint, fwUUID)
	fw := &ComponentFirmwareVersion{}
	r := ServerResponse{Record: fw}

	if err := c.get(ctx, path, &r); err != nil {
		return nil, nil, err
	}

	return fw, &r, nil
}

// ListServerComponentFirmware will return all firmwares with optional params to filter the results
func (c *Client) ListServerComponentFirmware(ctx context.Context, params *ComponentFirmwareVersionListParams) ([]ComponentFirmwareVersion, *ServerResponse, error) {
	firmwares := &[]ComponentFirmwareVersion{}
	r := ServerResponse{Records: firmwares}

	if err := c.list(ctx, serverComponentFirmwaresEndpoint, params, &r); err != nil {
		return nil, nil, err
	}

	return *firmwares, &r, nil
}

// UpdateServerComponentFirmware will to update a firmware with the new values passed in
func (c *Client) UpdateServerComponentFirmware(ctx context.Context, fwUUID uuid.UUID, firmware ComponentFirmwareVersion) (*ServerResponse, error) {
	path := fmt.Sprintf("%s/%s", serverComponentFirmwaresEndpoint, fwUUID)
	return c.put(ctx, path, firmware)
}
