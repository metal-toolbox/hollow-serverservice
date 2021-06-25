package hollow

import (
	"context"
)

const (
	biosConfigEndpoint = "bios-config"
)

// BIOSConfigService provides the ability to interact with bios configs via Hollow
type BIOSConfigService interface {
	CreateBIOSConfig(context.Context, BIOSConfig) error
}

// BIOSConfigServiceClient implements BIOSConfigService
type BIOSConfigServiceClient struct {
	client *Client
}

// CreateBIOSConfig will create a new BIOS Config
func (b *BIOSConfigServiceClient) CreateBIOSConfig(ctx context.Context, bc BIOSConfig) error {
	request, err := newPostRequest(ctx, b.client.url, biosConfigEndpoint, bc)
	if err != nil {
		return err
	}

	resp, err := b.client.do(request)
	if err != nil {
		return err
	}

	return ensureValidServerResponse(resp)
}
