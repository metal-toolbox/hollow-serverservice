package serverservice

import (
	"net/url"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.hollow.sh/serverservice/internal/models"
)

// ComponentFirmwareVersionListParams allows you to filter the results
type ComponentFirmwareVersionListParams struct {
	Vendor  string   `form:"vendor"`
	Model   []string `form:"model"`
	Version string   `form:"version"`
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
}

// queryMods converts the list params into sql conditions that can be added to sql queries
func (p *ComponentFirmwareVersionListParams) queryMods() []qm.QueryMod {
	mods := []qm.QueryMod{}

	if p.Vendor != "" {
		m := models.ComponentFirmwareVersionWhere.Vendor.EQ(p.Vendor)
		mods = append(mods, m)
	}

	if p.Model != nil {
		m := models.ComponentFirmwareVersionWhere.Model.EQ(p.Model)
		mods = append(mods, m)
	}

	if p.Version != "" {
		m := models.ComponentFirmwareVersionWhere.Version.EQ(p.Version)
		mods = append(mods, m)
	}

	return mods
}
