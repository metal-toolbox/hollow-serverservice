package hollow

import (
	"context"
	"encoding/json"

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

// HardwareComponentTypeListParams allows you to filter the results
type HardwareComponentTypeListParams struct {
	Name string
}

func (f *HardwareComponentTypeListParams) queryMap() map[string]string {
	m := make(map[string]string)

	if f.Name != "" {
		m["name"] = f.Name
	}

	return m
}

// Create will attempt to create a hardware component type in Hollow
func (c *HardwareComponentTypeServiceClient) Create(ctx context.Context, t HardwareComponentType) (*uuid.UUID, error) {
	request, err := newPostRequest(ctx, c.client.url, hardwareComponentTypeEndpoint, t)
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

	var r serverResponse
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r.UUID, nil
}

// List will return the hardware component types with optional params
func (c *HardwareComponentTypeServiceClient) List(ctx context.Context, params *HardwareComponentTypeListParams) ([]HardwareComponentType, error) {
	request, err := newGetRequest(ctx, c.client.url, hardwareComponentTypeEndpoint)
	if err != nil {
		return nil, err
	}

	if params != nil {
		q := request.URL.Query()

		for k, v := range params.queryMap() {
			q.Add(k, v)
		}

		request.URL.RawQuery = q.Encode()
	}

	resp, err := c.client.do(request)
	if err != nil {
		return nil, err
	}

	if err := ensureValidServerResponse(resp); err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var ct []HardwareComponentType
	if err = json.NewDecoder(resp.Body).Decode(&ct); err != nil {
		return nil, err
	}

	return ct, nil
}
