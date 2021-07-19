package hollow

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const (
	hardwareEndpoint                    = "hardware"
	hardwareVersionedAttributesEndpoint = "versioned-attributes"
)

// HardwareService provides the ability to interact with hardware via Hollow
type HardwareService interface {
	Create(context.Context, Hardware) (*uuid.UUID, error)
	Delete(context.Context, Hardware) error
	Get(context.Context, uuid.UUID) (*Hardware, error)
	List(context.Context, *HardwareListParams) ([]Hardware, error)
	GetVersionedAttributes(context.Context, uuid.UUID) ([]VersionedAttributes, error)
	CreateVersionedAttributes(context.Context, uuid.UUID, VersionedAttributes) (*uuid.UUID, error)
}

// HardwareServiceClient implements HardwareService
type HardwareServiceClient struct {
	client *Client
}

// Create will attempt to create hardware in Hollow and return the new hardware UUID
func (c *HardwareServiceClient) Create(ctx context.Context, hw Hardware) (*uuid.UUID, error) {
	return c.client.post(ctx, hardwareEndpoint, hw)
}

// Delete will attempt to delete hardware in Hollow and return an error on failure
func (c *HardwareServiceClient) Delete(ctx context.Context, hw Hardware) error {
	return c.client.delete(ctx, fmt.Sprintf("%s/%s", hardwareEndpoint, hw.UUID))
}

// Get will return a given piece of hardware by it's UUID
func (c *HardwareServiceClient) Get(ctx context.Context, hwUUID uuid.UUID) (*Hardware, error) {
	path := fmt.Sprintf("%s/%s", hardwareEndpoint, hwUUID)

	var hw Hardware
	if err := c.client.get(ctx, path, &hw); err != nil {
		return nil, err
	}

	return &hw, nil
}

// List will return all hardware with optional params to filter the results
func (c *HardwareServiceClient) List(ctx context.Context, params *HardwareListParams) ([]Hardware, error) {
	var hw []Hardware
	if err := c.client.list(ctx, hardwareEndpoint, params, &hw); err != nil {
		return nil, err
	}

	return hw, nil
}

// GetVersionedAttributes will return all the versioned attributes for a given piece of hardware
func (c *HardwareServiceClient) GetVersionedAttributes(ctx context.Context, hwUUID uuid.UUID) ([]VersionedAttributes, error) {
	path := fmt.Sprintf("%s/%s/%s", hardwareEndpoint, hwUUID, hardwareVersionedAttributesEndpoint)

	var val []VersionedAttributes
	if err := c.client.list(ctx, path, nil, &val); err != nil {
		return nil, err
	}

	return val, nil
}

// CreateVersionedAttributes will create a new versioned attribute for a given piece of hardware
func (c *HardwareServiceClient) CreateVersionedAttributes(ctx context.Context, hwUUID uuid.UUID, va VersionedAttributes) (*uuid.UUID, error) {
	path := fmt.Sprintf("%s/%s/%s", hardwareEndpoint, hwUUID, hardwareVersionedAttributesEndpoint)
	return c.client.put(ctx, path, va)
}
