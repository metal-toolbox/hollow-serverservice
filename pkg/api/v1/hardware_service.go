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
	hardwareEndpoint                    = "hardware"
	hardwareVersionedAttributesEndpoint = "versioned-attributes"
)

// HardwareService provides the ability to interact with hardware via Hollow
type HardwareService interface {
	Create(context.Context, Hardware) (*uuid.UUID, error)
	Get(context.Context, uuid.UUID) (*Hardware, error)
	List(context.Context, *HardwareListParams) ([]Hardware, error)
	VersionedAttributesGet(context.Context, uuid.UUID) ([]VersionedAttributes, error)
}

// HardwareServiceClient implements HardwareService
type HardwareServiceClient struct {
	client *Client
}

// HardwareListParams allows you to filter the results
type HardwareListParams struct {
	FacilityCode                 string `form:"facility-code"`
	AttributeListParams          []AttributeListParams
	VersionedAttributeListParams []AttributeListParams
}

func (p *HardwareListParams) setQuery(q url.Values) {
	if p.FacilityCode != "" {
		q.Set("facility-code", p.FacilityCode)
	}

	encodeAttributesListParams(p.AttributeListParams, "attr", q)
	encodeAttributesListParams(p.VersionedAttributeListParams, "ver_attr", q)
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

	for _, aF := range p.VersionedAttributeListParams {
		a := db.AttributesFilter{
			Namespace:        aF.Namespace,
			Keys:             aF.Keys,
			EqualValue:       aF.EqualValue,
			LessThanValue:    aF.LessThanValue,
			GreaterThanValue: aF.GreaterThanValue,
		}
		dbF.VersionedAttributesFilters = append(dbF.VersionedAttributesFilters, a)
	}

	return dbF
}

// VersionedAttributesGet will return all the versioned attributes for a given piece of hardware
func (c *HardwareServiceClient) VersionedAttributesGet(ctx context.Context, hwUUID uuid.UUID) ([]VersionedAttributes, error) {
	path := fmt.Sprintf("%s/%s/%s", hardwareEndpoint, hwUUID, hardwareVersionedAttributesEndpoint)

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

	var val []VersionedAttributes
	if err = json.NewDecoder(resp.Body).Decode(&val); err != nil {
		return nil, err
	}

	return val, nil
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

// Create will attempt to create hardware in Hollow and return the new hardware UUID
func (c *HardwareServiceClient) Create(ctx context.Context, hw Hardware) (*uuid.UUID, error) {
	request, err := newPostRequest(ctx, c.client.url, hardwareEndpoint, hw)
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
