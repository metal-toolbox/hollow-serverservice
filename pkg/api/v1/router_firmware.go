package serverservice

import (
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.hollow.sh/serverservice/internal/models"
)

func (r *Router) firmwareList(c *gin.Context) {
	pager := parsePagination(c)

	var params ComponentFirmwareVersionListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		badRequestResponse(c, "invalid filter", err)
		return
	}

	mods := params.queryMods()

	dbFirmwares, err := models.ComponentFirmwareVersions(mods...).All(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	count, err := models.ComponentFirmwareVersions(mods...).Count(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	firmwares := []ComponentFirmwareVersion{}

	for _, dbF := range dbFirmwares {
		f := ComponentFirmwareVersion{}
		if err := f.fromDBModel(dbF); err != nil {
			failedConvertingToVersioned(c, err)
			return
		}

		firmwares = append(firmwares, f)
	}

	pd := paginationData{
		pageCount:  len(firmwares),
		totalCount: count,
		pager:      pager,
	}

	listResponse(c, firmwares, pd)
}

func (r *Router) firmwareGet(c *gin.Context) {
	mods := []qm.QueryMod{
		qm.Where("id=?", c.Param("uuid")),
	}

	dbFirmware, err := models.ComponentFirmwareVersions(mods...).One(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	var firmware ComponentFirmwareVersion
	if err = firmware.fromDBModel(dbFirmware); err != nil {
		failedConvertingToVersioned(c, err)
		return
	}

	itemResponse(c, firmware)
}

func (r *Router) firmwareCreate(c *gin.Context) {
	var firmware ComponentFirmwareVersion
	if err := c.ShouldBindJSON(&firmware); err != nil {
		badRequestResponse(c, "invalid firmware", err)
		return
	}

	dbFirmware, err := firmware.toDBModel()
	if err != nil {
		badRequestResponse(c, "invalid firmware", err)
		return
	}

	if err := dbFirmware.Insert(c.Request.Context(), r.DB, boil.Infer()); err != nil {
		dbErrorResponse(c, err)
		return
	}

	createdResponse(c, dbFirmware.ID)
}

func (r *Router) firmwareDelete(c *gin.Context) {
	dbFirmware, err := r.loadComponentFirmwareVersionFromParams(c)
	if err != nil {
		return
	}

	if _, err = dbFirmware.Delete(c.Request.Context(), r.DB); err != nil {
		dbErrorResponse(c, err)
		return
	}

	deletedResponse(c)
}

func (r *Router) firmwareUpdate(c *gin.Context) {
	dbFirmware, err := r.loadComponentFirmwareVersionFromParams(c)
	if err != nil {
		return
	}

	var newValues ComponentFirmwareVersion
	if err := c.ShouldBindJSON(&newValues); err != nil {
		badRequestResponse(c, "invalid dbFirmware", err)
		return
	}

	dbFirmware.Vendor = null.StringFrom(newValues.Vendor)
	dbFirmware.Model = null.StringFrom(newValues.Model)
	dbFirmware.Filename = null.StringFrom(newValues.Filename)
	dbFirmware.Version = null.StringFrom(newValues.Version)
	dbFirmware.Component = null.StringFrom(newValues.Component)
	dbFirmware.Utility = null.StringFrom(newValues.Utility)
	dbFirmware.Sha = null.StringFrom(newValues.Sha)
	dbFirmware.UpstreamURL = null.StringFrom(newValues.UpstreamURL)
	dbFirmware.S3URL = null.StringFrom(newValues.S3URL)

	cols := boil.Infer()

	if _, err := dbFirmware.Update(c.Request.Context(), r.DB, cols); err != nil {
		dbErrorResponse(c, err)
		return
	}

	updatedResponse(c, dbFirmware.ID)
}
