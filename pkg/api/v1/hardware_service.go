package hollow

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

const (
	hardwareEndpoint           = "hardware"
	hardwareBIOSConfigEndpoint = "bios-configs"
)

// HardwareService provides the ability to interact with hardware via Hollow
type HardwareService interface {
	Create(context.Context, Hardware) error
	GetBIOSConfigs(context.Context, uuid.UUID) ([]BIOSConfig, error)
}

// HardwareServiceClient implements HardwareService
type HardwareServiceClient struct {
	client *Client
}

// GetBIOSConfigs will return all the BIOS Configs for a given piece of hardware
func (c *HardwareServiceClient) GetBIOSConfigs(ctx context.Context, hwUUID uuid.UUID) ([]BIOSConfig, error) {
	path := fmt.Sprintf("%s/%s/%s", hardwareEndpoint, hwUUID, hardwareBIOSConfigEndpoint)

	request, err := newGetRequest(ctx, c.client.url, path)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.do(request)
	if err != nil {
		return nil, err
	}

	if err := ensureValidServerResponse(resp); err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var bcl []BIOSConfig
	if err = json.NewDecoder(resp.Body).Decode(&bcl); err != nil {
		return nil, err
	}

	return bcl, nil
}

// Create will attempt to create hardware in Hollow
func (c *HardwareServiceClient) Create(ctx context.Context, hw Hardware) error {
	request, err := newPostRequest(ctx, c.client.url, hardwareEndpoint, hw)
	if err != nil {
		return err
	}

	resp, err := c.client.do(request)
	if err != nil {
		return err
	}

	return ensureValidServerResponse(resp)
}
