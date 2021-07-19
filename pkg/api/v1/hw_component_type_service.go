package hollow

import (
	"context"

	"github.com/google/uuid"
)

const (
	hardwareComponentTypeEndpoint = "hardware-component-types"
)

// HardwareComponentTypeService provides the ability to interact with hardware component types via Hollow
type HardwareComponentTypeService interface {
	Create(context.Context, HardwareComponentType) (*uuid.UUID, error)
	List(context.Context, *HardwareComponentTypeListParams) ([]HardwareComponentType, error)
}

// HardwareComponentTypeServiceClient implements HardwareService
type HardwareComponentTypeServiceClient struct {
	client *Client
}

// Create will attempt to create a hardware component type in Hollow
func (c *HardwareComponentTypeServiceClient) Create(ctx context.Context, t HardwareComponentType) (*uuid.UUID, error) {
	return c.client.post(ctx, hardwareComponentTypeEndpoint, t)
}

// List will return the hardware component types with optional params
func (c *HardwareComponentTypeServiceClient) List(ctx context.Context, params *HardwareComponentTypeListParams) ([]HardwareComponentType, error) {
	var ct []HardwareComponentType
	if err := c.client.list(ctx, hardwareComponentTypeEndpoint, params, &ct); err != nil {
		return nil, err
	}

	return ct, nil
}
