package hollow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/uuid"

	"go.metalkube.net/hollow/internal/db"
)

const (
	hardwareEndpoint           = "hardware"
	hardwareBIOSConfigEndpoint = "bios-configs"
)

// HardwareService provides the ability to interact with hardware via Hollow
type HardwareService interface {
	Create(context.Context, Hardware) error
	GetBIOSConfigs(context.Context, uuid.UUID) ([]BIOSConfig, error)
	Get(context.Context, uuid.UUID) (*Hardware, error)
	List(context.Context, *HardwareListParams) ([]Hardware, error)
}

// HardwareServiceClient implements HardwareService
type HardwareServiceClient struct {
	client *Client
}

// HardwareListParams allows you to filter the results
type HardwareListParams struct {
	FacilityCode        string                `form:"facility-code" query:"facility-code"`
	AttributeListParams []AttributeListParams `form:"attributes" query:"attributes"`
}

func (p *HardwareListParams) setQuery(q url.Values) {
	if p.FacilityCode != "" {
		q.Set("facility-code", p.FacilityCode)
	}

	encodeAttributesListParams(p.AttributeListParams, q)
}

func (p *HardwareListParams) dbFilter() *db.HardwareFilter {
	dbF := &db.HardwareFilter{
		FacilityCode: p.FacilityCode,
	}

	for _, aF := range p.AttributeListParams {
		a := db.AttributesFilter{
			Namespace:        aF.Namespace,
			Keys:             aF.Keys,
			EqualValue:       aF.EqualValue,
			LessThanValue:    aF.LessThanValue,
			GreaterThanValue: aF.GreaterThanValue,
		}
		dbF.AttributesFilters = append(dbF.AttributesFilters, a)
	}

	return dbF
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

// Get will return a given piece of hardware by it's UUID
func (c *HardwareServiceClient) Get(ctx context.Context, hwUUID uuid.UUID) (*Hardware, error) {
	path := fmt.Sprintf("%s/%s", hardwareEndpoint, hwUUID)

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

	var hw Hardware
	if err = json.NewDecoder(resp.Body).Decode(&hw); err != nil {
		return nil, err
	}

	return &hw, nil
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

// List will return all hardware with optional params to filter the results
func (c *HardwareServiceClient) List(ctx context.Context, params *HardwareListParams) ([]Hardware, error) {
	request, err := newGetRequest(ctx, c.client.url, hardwareEndpoint)
	if err != nil {
		return nil, err
	}

	if params != nil {
		q := request.URL.Query()

		params.setQuery(q)

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

	var hw []Hardware
	if err = json.NewDecoder(resp.Body).Decode(&hw); err != nil {
		return nil, err
	}

	return hw, nil
}
