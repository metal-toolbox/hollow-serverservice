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
	ListBIOSConfigs(context.Context, uuid.UUID) ([]BIOSConfig, error)
	CreateBIOSConfig(context.Context, uuid.UUID, BIOSConfig) error
}

// HardwareServiceClient implements HardwareService
type HardwareServiceClient struct {
	client *Client
}

// ListBIOSConfigs will return all the BIOS Configs for a given piece of hardware
func (h *HardwareServiceClient) ListBIOSConfigs(ctx context.Context, hwUUID uuid.UUID) ([]BIOSConfig, error) {
	path := fmt.Sprintf("%s/%s/%s", hardwareEndpoint, hwUUID, hardwareBIOSConfigEndpoint)

	request, err := newGetRequest(ctx, h.client.url, path)
	if err != nil {
		return nil, err
	}

	resp, err := h.client.do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var bcl []BIOSConfig
	if err = json.NewDecoder(resp.Body).Decode(&bcl); err != nil {
		return nil, err
	}

	return bcl, nil
}

// CreateBIOSConfig will create a new BIOS Config for a given piece of hardware
func (h *HardwareServiceClient) CreateBIOSConfig(ctx context.Context, hwUUID uuid.UUID, bc BIOSConfig) error {
	return nil
}
