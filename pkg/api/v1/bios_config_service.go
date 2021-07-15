package hollow

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

const (
	biosConfigEndpoint = "bios-config"
)

// BIOSConfigService provides the ability to interact with bios configs via Hollow
type BIOSConfigService interface {
	Create(context.Context, VersionedAttributes) (*uuid.UUID, error)
}

// BIOSConfigServiceClient implements BIOSConfigService
type BIOSConfigServiceClient struct {
	client *Client
}

// Create will create a new BIOS Config
func (b *BIOSConfigServiceClient) Create(ctx context.Context, va VersionedAttributes) (*uuid.UUID, error) {
	request, err := newPostRequest(ctx, b.client.url, biosConfigEndpoint, va)
	if err != nil {
		return nil, err
	}

	resp, err := b.client.do(request)
	if err != nil {
		return nil, err
	}

	if err := ensureValidServerResponse(resp); err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var r serverResponse
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r.UUID, nil
}
