package fleetdbapi

import (
	"net/url"

	"github.com/volatiletech/sqlboiler/types"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// ComponentFirmwareVersionListParams allows you to filter the results
type ComponentFirmwareVersionListParams struct {
	Vendor     string   `form:"vendor"`
	Model      []string `form:"model"`
	Version    string   `form:"version"`
	Filename   string   `form:"filename"`
	Checksum   string   `form:"checksum"`
	Pagination *PaginationParams
}

func (p *ComponentFirmwareVersionListParams) setQuery(q url.Values) {
	if p == nil {
		return
	}

	if p.Vendor != "" {
		q.Set("vendor", p.Vendor)
	}

	if p.Model != nil {
		for _, m := range p.Model {
			q.Add("model", m)
		}
	}

	if p.Version != "" {
		q.Set("version", p.Version)
	}

	if p.Filename != "" {
		q.Set("filename", p.Filename)
	}

	if p.Checksum != "" {
		q.Set("checksum", p.Checksum)
	}

	p.Pagination.setQuery(q)
}

// queryMods converts the list params into sql conditions that can be added to sql queries
func (p *ComponentFirmwareVersionListParams) queryMods() []qm.QueryMod {
	mods := []qm.QueryMod{}

	if p.Vendor != "" {
		m := models.ComponentFirmwareVersionWhere.Vendor.EQ(p.Vendor)
		mods = append(mods, m)
	}

	if p.Model != nil {
		m := qm.Where("model @> ?", types.StringArray(p.Model))
		mods = append(mods, m)
	}

	if p.Version != "" {
		m := models.ComponentFirmwareVersionWhere.Version.EQ(p.Version)
		mods = append(mods, m)
	}

	if p.Filename != "" {
		m := models.ComponentFirmwareVersionWhere.Filename.EQ(p.Filename)
		mods = append(mods, m)
	}

	if p.Checksum != "" {
		m := models.ComponentFirmwareVersionWhere.Checksum.EQ(p.Checksum)
		mods = append(mods, m)
	}

	return mods
}
