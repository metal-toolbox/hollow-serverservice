package serverservice

import (
	"net/url"
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
		q.Set("Version", p.Version)
	}
}

// queryMods converts the list params into sql conditions that can be added to
// sql queries
// func (p *FirmwareListParams) queryMods() []qm.QueryMod {
// 	mods := []qm.QueryMod{}
//
// 	if p.Vendor != "" {
// 		m := models.FirmwareWhere.Vendor.EQ(null.StringFrom(p.Vendor))
// 		mods = append(mods, m)
// 	}
//
// 	if p.Model != "" {
// 		m := models.FirmwareWhere.Model.EQ(null.StringFrom(p.Model))
// 		mods = append(mods, m)
// 	}
//
// 	if p.Version != "" {
// 		m := models.FirmwareWhere.Version.EQ(null.StringFrom(p.Version))
// 		mods = append(mods, m)
// 	}
//
// 	return mods
// }
