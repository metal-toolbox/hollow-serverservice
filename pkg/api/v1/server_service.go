package serverservice

import (
	"context"
	"encoding/json"
	"fmt"
	"path"

	"github.com/google/uuid"
)

const (
	serversEndpoint                     = "servers"
	serverAttributesEndpoint            = "attributes"
	serverComponentsEndpoint            = "components"
	serverVersionedAttributesEndpoint   = "versioned-attributes"
	serverComponentFirmwaresEndpoint    = "server-component-firmwares"
	serverCredentialsEndpoint           = "credentials"
	serverCredentialTypeEndpoint        = "server-credential-types"
	serverComponentFirmwareSetsEndpoint = "server-component-firmware-sets"
	bomInfoEndpoint                     = "bill-of-materials"
	uploadFileEndpoint                  = "batch-upload"
	bomByMacAOCAddressEndpoint          = "aoc-mac-address"
	bomByMacBMCAddressEndpoint          = "bmc-mac-address"
)

// ClientInterface provides an interface for the expected calls to interact with a server service api
type ClientInterface interface {
	Create(context.Context, Server) (*uuid.UUID, *ServerResponse, error)
	Delete(context.Context, Server) (*ServerResponse, error)
	DeleteNow(context.Context, Server) (*ServerResponse, error)
	Get(context.Context, uuid.UUID) (*Server, *ServerResponse, error)
	List(context.Context, *ServerListParams) ([]Server, *ServerResponse, error)
	Update(context.Context, uuid.UUID, Server) (*ServerResponse, error)
	CreateAttributes(context.Context, uuid.UUID, Attributes) (*ServerResponse, error)
	DeleteAttributes(ctx context.Context, u uuid.UUID, ns string) (*ServerResponse, error)
	GetAttributes(context.Context, uuid.UUID, string) (*Attributes, *ServerResponse, error)
	ListAttributes(context.Context, uuid.UUID, *PaginationParams) ([]Attributes, *ServerResponse, error)
	UpdateAttributes(ctx context.Context, u uuid.UUID, ns string, data json.RawMessage) (*ServerResponse, error)
	GetComponents(context.Context, uuid.UUID, *PaginationParams) ([]ServerComponent, *ServerResponse, error)
	ListComponents(context.Context, *ServerComponentListParams) ([]ServerComponent, *ServerResponse, error)
	CreateComponents(context.Context, uuid.UUID, ServerComponentSlice) (*ServerResponse, error)
	UpdateComponents(context.Context, uuid.UUID, ServerComponentSlice) (*ServerResponse, error)
	DeleteServerComponents(context.Context, uuid.UUID) (*ServerResponse, error)
	CreateVersionedAttributes(context.Context, uuid.UUID, VersionedAttributes) (*ServerResponse, error)
	GetVersionedAttributes(context.Context, uuid.UUID, string) ([]VersionedAttributes, *ServerResponse, error)
	ListVersionedAttributes(context.Context, uuid.UUID) ([]VersionedAttributes, *ServerResponse, error)
	CreateServerComponentFirmware(context.Context, ComponentFirmwareVersion) (*uuid.UUID, *ServerResponse, error)
	DeleteServerComponentFirmware(context.Context, ComponentFirmwareVersion) (*ServerResponse, error)
	GetServerComponentFirmware(context.Context, uuid.UUID) (*ComponentFirmwareVersion, *ServerResponse, error)
	ListServerComponentFirmware(context.Context, *ComponentFirmwareVersionListParams) ([]ComponentFirmwareVersion, *ServerResponse, error)
	UpdateServerComponentFirmware(context.Context, uuid.UUID, ComponentFirmwareVersion) (*ServerResponse, error)
	CreateServerComponentFirmwareSet(context.Context, ComponentFirmwareSetRequest) (*uuid.UUID, *ServerResponse, error)
	UpdateComponentFirmwareSetRequest(context.Context, ComponentFirmwareSetRequest) (*uuid.UUID, *ServerResponse, error)
	GetServerComponentFirmwareSet(context.Context, uuid.UUID) (*ComponentFirmwareSet, *ServerResponse, error)
	ListServerComponentFirmwareSet(context.Context, *ComponentFirmwareSetListParams) ([]ComponentFirmwareSet, *ServerResponse, error)
	DeleteServerComponentFirmwareSet(context.Context, uuid.UUID) (*ServerResponse, error)
	GetCredential(context.Context, uuid.UUID, string) (*ServerCredential, *ServerResponse, error)
	SetCredential(context.Context, uuid.UUID, string, string) (*ServerResponse, error)
	DeleteCredential(context.Context, uuid.UUID, string) (*ServerResponse, error)
	ListServerCredentialTypes(context.Context) (*ServerResponse, error)
	BillOfMaterialsBatchUpload(context.Context, []Bom) (*ServerResponse, error)
	GetBomInfoByAOCMacAddr(context.Context, string) (*Bom, *ServerResponse, error)
	GetBomInfoByBMCMacAddr(context.Context, string) (*Bom, *ServerResponse, error)
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

// Delete will attempt to soft delete a server in Hollow and return an error on failure
func (c *Client) Delete(ctx context.Context, srv Server) (*ServerResponse, error) {
	return c.delete(ctx, fmt.Sprintf("%s/%s", serversEndpoint, srv.UUID))
}

// DeleteNow will attempt to hard delete a server in Hollow and return an error on failure
func (c *Client) DeleteNow(ctx context.Context, srv Server) (*ServerResponse, error) {
	return c.delete(ctx, fmt.Sprintf("%s/%s/true", serversEndpoint, srv.UUID))
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

// GetComponents will get all the components for a given server
func (c *Client) GetComponents(ctx context.Context, srvUUID uuid.UUID, params *PaginationParams) (ServerComponentSlice, *ServerResponse, error) {
	sc := &ServerComponentSlice{}
	r := ServerResponse{Records: sc}

	path := fmt.Sprintf("%s/%s/%s", serversEndpoint, srvUUID, serverComponentsEndpoint)
	if err := c.list(ctx, path, params, &r); err != nil {
		return nil, nil, err
	}

	return *sc, &r, nil
}

// ListComponents will get all the components matching the given parameters
func (c *Client) ListComponents(ctx context.Context, params *ServerComponentListParams) (ServerComponentSlice, *ServerResponse, error) {
	sc := &ServerComponentSlice{}
	r := ServerResponse{Records: sc}

	path := fmt.Sprintf("%s/%s", serversEndpoint, serverComponentsEndpoint)
	if err := c.list(ctx, path, params, &r); err != nil {
		return nil, nil, err
	}

	return *sc, &r, nil
}

// CreateComponents will create given components for a given server
func (c *Client) CreateComponents(ctx context.Context, srvUUID uuid.UUID, components ServerComponentSlice) (*ServerResponse, error) {
	path := fmt.Sprintf("%s/%s/%s", serversEndpoint, srvUUID, serverComponentsEndpoint)
	return c.post(ctx, path, components)
}

// UpdateComponents will update given components for a given server
func (c *Client) UpdateComponents(ctx context.Context, srvUUID uuid.UUID, components ServerComponentSlice) (*ServerResponse, error) {
	path := fmt.Sprintf("%s/%s/%s", serversEndpoint, srvUUID, serverComponentsEndpoint)
	return c.put(ctx, path, components)
}

// DeleteServerComponents will delete components for the given server identifier.
func (c *Client) DeleteServerComponents(ctx context.Context, srvUUID uuid.UUID) (*ServerResponse, error) {
	return c.delete(ctx, fmt.Sprintf("%s/%s/%s", serversEndpoint, srvUUID, serverComponentsEndpoint))
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

// CreateServerComponentFirmwareSet will attempt to create a firmware set in Hollow and return the firmware UUID
func (c *Client) CreateServerComponentFirmwareSet(ctx context.Context, set ComponentFirmwareSetRequest) (*uuid.UUID, *ServerResponse, error) {
	resp, err := c.post(ctx, serverComponentFirmwareSetsEndpoint, set)
	if err != nil {
		return nil, nil, err
	}

	u, err := uuid.Parse(resp.Slug)
	if err != nil {
		return nil, resp, nil
	}

	return &u, resp, nil
}

// DeleteServerComponentFirmwareSet will attempt to delete a firmware set and return an error on failure
func (c *Client) DeleteServerComponentFirmwareSet(ctx context.Context, firmwareSetID uuid.UUID) (*ServerResponse, error) {
	return c.delete(ctx, fmt.Sprintf("%s/%s", serverComponentFirmwareSetsEndpoint, firmwareSetID))
}

// GetServerComponentFirmwareSet will return a firmware by its UUID
func (c *Client) GetServerComponentFirmwareSet(ctx context.Context, fwSetUUID uuid.UUID) (*ComponentFirmwareSet, *ServerResponse, error) {
	path := fmt.Sprintf("%s/%s", serverComponentFirmwareSetsEndpoint, fwSetUUID)
	fws := &ComponentFirmwareSet{}
	r := ServerResponse{Record: fws}

	if err := c.get(ctx, path, &r); err != nil {
		return nil, nil, err
	}

	return fws, &r, nil
}

// ListServerComponentFirmwareSet will return all firmwares with optional params to filter the results
func (c *Client) ListServerComponentFirmwareSet(ctx context.Context, params *ComponentFirmwareSetListParams) ([]ComponentFirmwareSet, *ServerResponse, error) {
	firmwareSets := &[]ComponentFirmwareSet{}
	r := ServerResponse{Records: firmwareSets}

	if err := c.list(ctx, serverComponentFirmwareSetsEndpoint, params, &r); err != nil {
		return nil, nil, err
	}

	return *firmwareSets, &r, nil
}

// UpdateComponentFirmwareSetRequest will add a firmware set with the new firmware id(s) passed in the firmwareSet parameter
func (c *Client) UpdateComponentFirmwareSetRequest(ctx context.Context, fwSetUUID uuid.UUID, firmwareSet ComponentFirmwareSetRequest) (*ServerResponse, error) {
	path := fmt.Sprintf("%s/%s", serverComponentFirmwareSetsEndpoint, fwSetUUID)
	return c.put(ctx, path, firmwareSet)
}

// RemoveServerComponentFirmwareSetFirmware will update a firmware set by removing the mapping for the firmware id(s) passed in the firmwareSet parameter
func (c *Client) RemoveServerComponentFirmwareSetFirmware(ctx context.Context, fwSetUUID uuid.UUID, firmwareSet ComponentFirmwareSetRequest) (*ServerResponse, error) {
	path := fmt.Sprintf("%s/%s/remove-firmware", serverComponentFirmwareSetsEndpoint, fwSetUUID)
	return c.post(ctx, path, firmwareSet)
}

// GetCredential will return the secret for the secret type for the given server UUID
func (c *Client) GetCredential(ctx context.Context, srvUUID uuid.UUID, secretSlug string) (*ServerCredential, *ServerResponse, error) {
	p := path.Join(serversEndpoint, srvUUID.String(), serverCredentialsEndpoint, secretSlug)
	secret := &ServerCredential{}
	r := ServerResponse{Record: secret}

	if err := c.get(ctx, p, &r); err != nil {
		return nil, nil, err
	}

	return secret, &r, nil
}

// SetCredential will set the secret for a given server UUID and secret type.
func (c *Client) SetCredential(ctx context.Context, srvUUID uuid.UUID, secretSlug, username, password string) (*ServerResponse, error) {
	p := path.Join(serversEndpoint, srvUUID.String(), serverCredentialsEndpoint, secretSlug)
	secret := &serverCredentialValues{
		Password: password,
		Username: username,
	}

	return c.put(ctx, p, secret)
}

// DeleteCredential will remove the secret for a given server UUID and secret type.
func (c *Client) DeleteCredential(ctx context.Context, srvUUID uuid.UUID, secretSlug string) (*ServerResponse, error) {
	p := path.Join(serversEndpoint, srvUUID.String(), serverCredentialsEndpoint, secretSlug)

	return c.delete(ctx, p)
}

// ListServerCredentialTypes will return all server secret types
func (c *Client) ListServerCredentialTypes(ctx context.Context, params *PaginationParams) ([]ServerCredentialType, *ServerResponse, error) {
	types := &[]ServerCredentialType{}
	r := ServerResponse{Records: types}

	if err := c.list(ctx, serverCredentialTypeEndpoint, params, &r); err != nil {
		return nil, nil, err
	}

	return *types, &r, nil
}

// CreateServerCredentialType will create a new server secret type
func (c *Client) CreateServerCredentialType(ctx context.Context, sType *ServerCredentialType) (*ServerResponse, error) {
	return c.post(ctx, serverCredentialTypeEndpoint, sType)
}

// BillOfMaterialsBatchUpload will attempt to write multiple boms to database.
func (c *Client) BillOfMaterialsBatchUpload(ctx context.Context, boms []Bom) (*ServerResponse, error) {
	path := fmt.Sprintf("%s/%s", bomInfoEndpoint, uploadFileEndpoint)

	resp, err := c.post(ctx, path, boms)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetBomInfoByAOCMacAddr will return the bom info object by the aoc mac address.
func (c *Client) GetBomInfoByAOCMacAddr(ctx context.Context, aocMacAddr string) (*Bom, *ServerResponse, error) {
	path := fmt.Sprintf("%s/%s/%s", bomInfoEndpoint, bomByMacAOCAddressEndpoint, aocMacAddr)
	bom := &Bom{}
	r := ServerResponse{Record: bom}

	if err := c.get(ctx, path, &r); err != nil {
		return nil, nil, err
	}

	return bom, &r, nil
}

// GetBomInfoByBMCMacAddr will return the bom info object by the bmc mac address.
func (c *Client) GetBomInfoByBMCMacAddr(ctx context.Context, bmcMacAddr string) (*Bom, *ServerResponse, error) {
	path := fmt.Sprintf("%s/%s/%s", bomInfoEndpoint, bomByMacBMCAddressEndpoint, bmcMacAddr)
	bom := &Bom{}
	r := ServerResponse{Record: bom}

	if err := c.get(ctx, path, &r); err != nil {
		return nil, nil, err
	}

	return bom, &r, nil
}
