package fleetdbapi

import (
	"database/sql"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

func (r *Router) bomsUpload(c *gin.Context) {
	var boms []Bom
	if err := c.ShouldBindJSON(&boms); err != nil {
		badRequestResponse(c, "invalid payload: []Bom{}", err)
		return
	}

	err := crdb.ExecuteTx(c.Request.Context(), r.DB.DB, nil, func(tx *sql.Tx) error {
		for _, bom := range boms {
			dbBomInfo, err := (bom).toDBModel()
			if err != nil {
				return err
			}

			if err := dbBomInfo.Insert(c.Request.Context(), r.DB, boil.Infer()); err != nil {
				return err
			}

			dbAocMacAddrsBoms, err := (bom).toAocMacAddressDBModels()
			if err != nil {
				return err
			}

			for _, dbAocMacAddrsBom := range dbAocMacAddrsBoms {
				if err := dbAocMacAddrsBom.Insert(c.Request.Context(), r.DB, boil.Infer()); err != nil {
					return err
				}
			}

			dbBmcMacAddrsBoms, err := (bom).toBmcMacAddressDBModels()
			if err != nil {
				return err
			}

			for _, dbBmcMacAddrsBom := range dbBmcMacAddrsBoms {
				if err := dbBmcMacAddrsBom.Insert(c.Request.Context(), r.DB, boil.Infer()); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	createdResponse(c, "")
}

func (r *Router) getBomFromAocMacAddress(c *gin.Context) {
	mods := []qm.QueryMod{
		qm.Where("aoc_mac_address=?", c.Param("aoc_mac_address")),
	}

	aocMacAddr, err := models.AocMacAddresses(mods...).One(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	mods = []qm.QueryMod{
		qm.Where("serial_num=?", aocMacAddr.SerialNum),
	}

	bomInfo, err := models.BomInfos(mods...).One(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	bom := Bom{}
	if err = bom.fromDBModel(bomInfo); err != nil {
		dbErrorResponse(c, err)
		return
	}

	itemResponse(c, bom)
}

func (r *Router) getBomFromBmcMacAddress(c *gin.Context) {
	mods := []qm.QueryMod{
		qm.Where("bmc_mac_address=?", c.Param("bmc_mac_address")),
	}

	bmcMacAddr, err := models.BMCMacAddresses(mods...).One(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	mods = []qm.QueryMod{
		qm.Where("serial_num=?", bmcMacAddr.SerialNum),
	}

	bomInfo, err := models.BomInfos(mods...).One(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	bom := Bom{}
	if err = bom.fromDBModel(bomInfo); err != nil {
		dbErrorResponse(c, err)
		return
	}

	itemResponse(c, bom)
}
