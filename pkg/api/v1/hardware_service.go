package hollow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
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

// AttributeListParams allow you to filter the results based on attributes
type AttributeListParams struct {
	Namespace        string   `form:"namespace" query:"namespace"`
	Keys             []string `form:"keys" query:"keys"`
	EqualValue       string   `form:"equals" query:"equals"`
	LessThanValue    int      `form:"less-than" query:"less-than"`
	GreaterThanValue int      `form:"greater-than" query:"greater-than"`
}

func (p *HardwareListParams) setQuery(q url.Values) {
	if p.FacilityCode != "" {
		q.Set("facility-code", p.FacilityCode)
	}

	for i, ap := range p.AttributeListParams {
		keyPrefix := fmt.Sprintf("attributes_%d_", i)

		q.Set(keyPrefix+"namespace", ap.Namespace)

		for _, k := range ap.Keys {
			q.Add(keyPrefix+"keys", k)
		}

		if ap.EqualValue != "" {
			q.Set(keyPrefix+"equals", ap.EqualValue)
		}

		if ap.LessThanValue != 0 {
			q.Set(keyPrefix+"less-than", fmt.Sprint(ap.LessThanValue))
		}

		if ap.GreaterThanValue != 0 {
			q.Set(keyPrefix+"greater-than", fmt.Sprint(ap.GreaterThanValue))
		}
	}
}

func parseQueryAttributesListParams(c *gin.Context) ([]AttributeListParams, error) {
	var err error

	alp := []AttributeListParams{}
	i := 0

	for {
		keyPrefix := fmt.Sprintf("attributes_%d_", i)

		ns := c.Query(keyPrefix + "namespace")
		if ns == "" {
			break
		}

		a := AttributeListParams{
			Namespace: ns,
			Keys:      c.QueryArray(keyPrefix + "keys"),
		}

		equals := c.Query(keyPrefix + "equals")
		if equals != "" {
			a.EqualValue = equals
		}

		lt := c.Query(keyPrefix + "less-than")
		if lt != "" {
			a.LessThanValue, err = strconv.Atoi(lt)
			if err != nil {
				return nil, err
			}
		}

		gt := c.Query(keyPrefix + "greater-than")
		if gt != "" {
			a.GreaterThanValue, err = strconv.Atoi(gt)
			if err != nil {
				return nil, err
			}
		}

		alp = append(alp, a)
		i++
	}

	return alp, nil
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
