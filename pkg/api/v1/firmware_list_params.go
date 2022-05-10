package serverservice

import (
	"net/url"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.hollow.sh/serverservice/internal/models"
)

// ComponentFirmwareVersionListParams allows you to filter the results
type ComponentFirmwareVersionListParams struct {
	Vendor  string `form:"vendor"`
	Model   string `form:"model"`
	Version string `form:"version"`
}

func (p *ComponentFirmwareVersionListParams) setQuery(q url.Values) {
	if p == nil {
		return
	}

	if p.Vendor != "" {
		q.Set("vendor", p.Vendor)
	}

	if p.Model != "" {
		q.Set("model", p.Model)
	}

	if p.Version != "" {
		q.Set("version", p.Version)
	}
}

// queryMods converts the list params into sql conditions that can be added to sql queries
func (p *ComponentFirmwareVersionListParams) queryMods() []qm.QueryMod {
	mods := []qm.QueryMod{}

	if p.Vendor != "" {
		m := models.ComponentFirmwareVersionWhere.Vendor.EQ(null.StringFrom(p.Vendor))
		mods = append(mods, m)
	}

	if p.Model != "" {
		m := models.ComponentFirmwareVersionWhere.Model.EQ(null.StringFrom(p.Model))
		mods = append(mods, m)
	}

	if p.Version != "" {
		m := models.ComponentFirmwareVersionWhere.Version.EQ(null.StringFrom(p.Version))
		mods = append(mods, m)
	}

	return mods
}
