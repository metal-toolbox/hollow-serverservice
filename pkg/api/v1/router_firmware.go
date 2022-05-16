package serverservice

import (
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.hollow.sh/serverservice/internal/models"
)

func (r *Router) serverComponentFirmwareList(c *gin.Context) {
	pager := parsePagination(c)

	var params ComponentFirmwareVersionListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		badRequestResponse(c, "invalid filter", err)
		return
	}

	mods := params.queryMods()

	count, err := models.ComponentFirmwareVersions(mods...).Count(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	// add pagination
	pager.Preload = false
	pager.OrderBy = models.ComponentFirmwareVersionTableColumns.Vendor + " DESC"
	mods = append(mods, pager.queryMods()...)

	dbFirmwares, err := models.ComponentFirmwareVersions(mods...).All(c.Request.Context(), r.DB)
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

func (r *Router) serverComponentFirmwareGet(c *gin.Context) {
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

func (r *Router) serverComponentFirmwareCreate(c *gin.Context) {
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

func (r *Router) serverComponentFirmwareDelete(c *gin.Context) {
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

func (r *Router) serverComponentFirmwareUpdate(c *gin.Context) {
	dbFirmware, err := r.loadComponentFirmwareVersionFromParams(c)
	if err != nil {
		return
	}

	var newValues ComponentFirmwareVersion
	if err := c.ShouldBindJSON(&newValues); err != nil {
		badRequestResponse(c, "invalid dbFirmware", err)
		return
	}

	dbFirmware.Vendor = newValues.Vendor
	dbFirmware.Model = newValues.Model
	dbFirmware.Filename = newValues.Filename
	dbFirmware.Version = newValues.Version
	dbFirmware.Component = newValues.Component
	dbFirmware.Checksum = newValues.Checksum
	dbFirmware.UpstreamURL = newValues.UpstreamURL
	dbFirmware.S3URL = newValues.S3URL

	cols := boil.Infer()

	if _, err := dbFirmware.Update(c.Request.Context(), r.DB, cols); err != nil {
		dbErrorResponse(c, err)
		return
	}

	updatedResponse(c, dbFirmware.ID)
}
