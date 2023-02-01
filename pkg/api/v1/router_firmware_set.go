package serverservice

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.hollow.sh/serverservice/internal/models"
)

var (
	errComponentFirmwareSetRequest = errors.New("error in component firmware set request")
	errComponentFirmwareSetMap     = errors.New("error mapping firmware in set")
	errDBErr                       = errors.New("db error")
)

// Firmware sets group firmware versions
//
// - firmware sets can only reference to unique firmware versions based on the vendor, model, component attributes.

func (r *Router) serverComponentFirmwareSetList(c *gin.Context) {
	pager := parsePagination(c)

	// unmarshal query parameters
	var params ComponentFirmwareSetListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		badRequestResponse(c, "invalid filter payload: ComponentFirmwareSetListParams{}", err)
		return
	}

	// query parameters to query mods
	params.AttributeListParams = parseQueryAttributesListParams(c, "attr")
	mods := params.queryMods(models.TableNames.ComponentFirmwareSet)
	mods = append(mods, qm.Load(models.ComponentFirmwareSetRels.FirmwareSetAttributesFirmwareSets))

	// count rows
	count, err := models.ComponentFirmwareSets(mods...).Count(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	// add pagination
	pager.Preload = false
	pager.OrderBy = models.ComponentFirmwareSetColumns.CreatedAt + " DESC"

	// load firmware sets
	dbFirmwareSets, err := models.ComponentFirmwareSets(mods...).All(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	firmwareSets := make([]ComponentFirmwareSet, 0, count)

	// load firmware set mappings
	for _, dbFS := range dbFirmwareSets {
		f := ComponentFirmwareSet{}

		firmwares, err := r.queryFirmwareSetFirmware(c.Request.Context(), dbFS.ID)
		if err != nil {
			dbErrorResponse(c, err)
			return
		}

		if err := f.fromDBModel(dbFS, firmwares); err != nil {
			failedConvertingToVersioned(c, err)
			return
		}

		firmwareSets = append(firmwareSets, f)
	}

	pd := paginationData{
		pageCount:  len(firmwareSets),
		totalCount: count,
		pager:      pager,
	}

	listResponse(c, firmwareSets, pd)
}

func (r *Router) serverComponentFirmwareSetGet(c *gin.Context) {
	setID := c.Param("uuid")
	if setID == "" || setID == uuid.Nil.String() {
		badRequestResponse(c, "expected a firmware set UUID, got none", errComponentFirmwareSetRequest)
		return
	}

	setIDParsed, err := uuid.Parse(setID)
	if err != nil {
		badRequestResponse(c, "invalid firmware set UUID: "+setID, err)
	}

	// query firmware set
	mods := []qm.QueryMod{
		qm.Where("id=?", setIDParsed),
		qm.Load(models.ComponentFirmwareSetRels.FirmwareSetAttributesFirmwareSets),
	}

	dbFirmwareSet, err := models.ComponentFirmwareSets(mods...).One(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	firmwares, err := r.queryFirmwareSetFirmware(c.Request.Context(), dbFirmwareSet.ID)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	// convert from db model type
	var firmwareSet ComponentFirmwareSet
	if err = firmwareSet.fromDBModel(dbFirmwareSet, firmwares); err != nil {
		failedConvertingToVersioned(c, err)
		return
	}

	itemResponse(c, firmwareSet)
}

func (r *Router) queryFirmwareSetFirmware(ctx context.Context, firmwareSetID string) ([]*models.ComponentFirmwareVersion, error) {
	mapMods := []qm.QueryMod{
		qm.Where("firmware_set_id=?", firmwareSetID),
		qm.Load(models.ComponentFirmwareSetMapRels.Firmware),
	}

	// query firmware set references
	dbFirmwareSetMap, err := models.ComponentFirmwareSetMaps(mapMods...).All(ctx, r.DB)
	if err != nil {
		return nil, err
	}

	firmwares := []*models.ComponentFirmwareVersion{}

	for _, m := range dbFirmwareSetMap {
		if m.R != nil && m.R.Firmware != nil {
			firmwares = append(firmwares, m.R.Firmware)
		}
	}

	return firmwares, nil
}

func (r *Router) serverComponentFirmwareSetCreate(c *gin.Context) {
	var firmwareSetPayload ComponentFirmwareSetRequest

	if err := c.ShouldBindJSON(&firmwareSetPayload); err != nil {
		badRequestResponse(c, "invalid payload: ComponentFirmwareSetCreate{}", err)
		return
	}

	if firmwareSetPayload.Name == "" {
		badRequestResponse(
			c,
			"invalid payload: ComponentFirmwareSetCreate{}",
			errors.Wrap(errSrvComponentPayload, "required attribute not set: Name"),
		)

		return
	}

	// vet and parse firmware uuids
	if len(firmwareSetPayload.ComponentFirmwareUUIDs) == 0 {
		err := errors.Wrap(errComponentFirmwareSetRequest, "expected one or more firmware UUIDs, got none")
		badRequestResponse(
			c,
			"",
			err,
		)

		return
	}

	firmwareUUIDs, err := r.firmwareSetVetFirmwareUUIDsForCreate(c, firmwareSetPayload.ComponentFirmwareUUIDs)
	if err != nil {
		if errors.Is(errDBErr, err) {
			dbErrorResponse(c, err)
			return
		}

		badRequestResponse(c, "", err)

		return
	}

	dbFirmwareSet, err := firmwareSetPayload.toDBModelFirmwareSet()
	if err != nil {
		badRequestResponse(c, "invalid db model: ComponentFirmwareSet", err)
		return
	}

	err = r.firmwareSetCreateTx(c.Request.Context(), dbFirmwareSet, firmwareSetPayload.Attributes, firmwareUUIDs)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	createdResponse(c, dbFirmwareSet.ID)
}

func (r *Router) firmwareSetVetFirmwareUUIDsForCreate(c *gin.Context, firmwareUUIDs []string) ([]uuid.UUID, error) {
	// validate and collect firmware UUIDs
	vetted := []uuid.UUID{}

	// unique is a map of keys to limit firmware sets to firmwares with unique vendor, version, component attributes.
	unique := map[string]bool{}

	for _, firmwareUUID := range firmwareUUIDs {
		// parse uuid
		firmwareUUIDParsed, err := uuid.Parse(firmwareUUID)
		if err != nil {
			return nil, errors.Wrap(errComponentFirmwareSetRequest, err.Error()+" invalid firmware UUID: "+firmwareUUID)
		}

		// validate component firmware version exists
		firmwareVersion, err := models.FindComponentFirmwareVersion(c.Request.Context(), r.DB, firmwareUUID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errors.Wrap(err, "firmware object with given UUID does not exist: "+firmwareUUID)
			}

			return nil, errors.Wrap(errDBErr, err.Error())
		}

		// validate firmware is unique based on vendor, version, component attributes
		key := strings.ToLower(firmwareVersion.Vendor) + strings.ToLower(firmwareVersion.Version) + strings.ToLower(firmwareVersion.Component)

		_, exists := unique[key]
		if exists {
			return nil, errors.Wrap(
				errComponentFirmwareSetMap,
				"A firmware set can only reference unique firmware versions based on the vendor, version, component attributes",
			)
		}

		unique[key] = true

		vetted = append(vetted, firmwareUUIDParsed)
	}

	return vetted, nil
}

func (r *Router) firmwareSetCreateTx(ctx context.Context, dbFirmwareSet *models.ComponentFirmwareSet, attrs []Attributes, firmwareUUIDs []uuid.UUID) error {
	// being transaction to insert a new firmware set and its references
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// nolint:errcheck // TODO(joel): log error
	defer tx.Rollback()

	// insert set
	if err := dbFirmwareSet.Insert(ctx, tx, boil.Infer()); err != nil {
		return err
	}

	// insert attributes
	for _, attributes := range attrs {
		dbAttributes := attributes.toDBModelAttributesFirmwareSet()
		dbAttributes.FirmwareSetID = null.StringFrom(dbFirmwareSet.ID)

		err = dbFirmwareSet.AddFirmwareSetAttributesFirmwareSets(ctx, tx, true, dbAttributes)
		if err != nil {
			return err
		}
	}

	// add firmware references
	for _, id := range firmwareUUIDs {
		m := models.ComponentFirmwareSetMap{FirmwareSetID: dbFirmwareSet.ID, FirmwareID: id.String()}

		err := m.Insert(ctx, tx, boil.Infer())
		if err != nil {
			return err
		}
	}

	// commit
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *Router) serverComponentFirmwareSetUpdate(c *gin.Context) {
	dbFirmware, err := r.componentFirmwareSetFromParams(c)
	if err != nil {
		badRequestResponse(c, "invalid payload: ComponentFirmwareSet{}", err)
		return
	}

	var newValues ComponentFirmwareSetRequest
	if err := c.ShouldBindJSON(&newValues); err != nil {
		badRequestResponse(c, "invalid payload: ComponentFirmwareSet{}", err)
		return
	}

	// firmware set ID is expected for updates
	if newValues.ID == uuid.Nil {
		badRequestResponse(
			c,
			"",
			errors.Wrap(errComponentFirmwareSetRequest, "expected a valid firmware set ID, got none"),
		)

		return
	}

	dbFirmwareSet, err := newValues.toDBModelFirmwareSet()
	if err != nil {
		badRequestResponse(c, "invalid db model: ComponentFirmwareSet", err)
		return
	}

	dbAttributesFirmwareSet := make([]*models.AttributesFirmwareSet, 0, len(newValues.Attributes))

	for _, attributes := range newValues.Attributes {
		attr := attributes.toDBModelAttributesFirmwareSet()

		attr.FirmwareSetID = null.StringFrom(newValues.ID.String())
		dbAttributesFirmwareSet = append(dbAttributesFirmwareSet, attr)
	}

	// vet and parse firmware uuids
	var firmwareUUIDs []uuid.UUID

	if len(newValues.ComponentFirmwareUUIDs) > 0 {
		firmwareUUIDs, err = r.firmwareSetVetFirmwareUUIDsForUpdate(c.Request.Context(), dbFirmwareSet, newValues.ComponentFirmwareUUIDs)
		if err != nil {
			if errors.Is(errDBErr, err) {
				dbErrorResponse(c, err)
				return
			}

			badRequestResponse(c, "", err)

			return
		}
	}

	err = r.firmwareSetUpdateTx(c.Request.Context(), dbFirmwareSet, dbAttributesFirmwareSet, firmwareUUIDs)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	updatedResponse(c, dbFirmware.ID)
}

func (r *Router) firmwareSetVetFirmwareUUIDsForUpdate(ctx context.Context, firmwareSet *models.ComponentFirmwareSet, firmwareUUIDs []string) ([]uuid.UUID, error) {
	// firmware uuids are expected
	if len(firmwareUUIDs) == 0 {
		return nil, errors.Wrap(errComponentFirmwareSetRequest, "expected one or more firmware UUIDs, got none")
	}

	// validate and collect firmware UUIDs
	vetted := []uuid.UUID{}

	// unique is a map of keys to limit firmware sets to include only firmwares with,
	// unique vendor, version, component attributes.
	unique := map[string]bool{}

	if len(firmwareUUIDs) == 0 {
		return nil, errors.Wrap(errComponentFirmwareSetRequest, "expected one or more firmware UUIDs, got none")
	}

	for _, firmwareUUID := range firmwareUUIDs {
		// parse uuid
		firmwareUUIDParsed, err := uuid.Parse(firmwareUUID)
		if err != nil {
			return nil, errors.Wrap(errComponentFirmwareSetRequest, err.Error()+"invalid firmware UUID: "+firmwareUUID)
		}

		// validate firmware isn't part of set
		setMap, err := r.firmwareSetMap(ctx, firmwareSet, firmwareUUIDParsed)
		if err != nil {
			return nil, err
		}

		if len(setMap) > 0 {
			return nil, errors.Wrap(
				errComponentFirmwareSetRequest,
				fmt.Sprintf("firmware '%s' exists in firmware set '%s' ", firmwareUUID, firmwareSet.Name),
			)
		}

		// validate component firmware version exists
		firmwareVersion, err := models.FindComponentFirmwareVersion(ctx, r.DB, firmwareUUID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errors.Wrap(errComponentFirmwareSetRequest, "firmware object with given UUID does not exist: "+firmwareUUID)
			}

			return nil, errors.Wrap(errDBErr, err.Error())
		}

		// validate firmware is unique based on vendor, version, component attributes
		key := strings.ToLower(firmwareVersion.Vendor) + strings.ToLower(firmwareVersion.Version) + strings.ToLower(firmwareVersion.Component)

		_, duplicateVendorModelComponent := unique[key]
		if duplicateVendorModelComponent {
			return nil, errors.Wrap(
				errComponentFirmwareSetMap,
				"A firmware set can only reference unique firmware versions based on the vendor, version, component attributes",
			)
		}

		unique[key] = true

		vetted = append(vetted, firmwareUUIDParsed)
	}

	return vetted, nil
}

func (r *Router) firmwareSetMap(ctx context.Context, firmwareSet *models.ComponentFirmwareSet, firmwareUUID uuid.UUID) ([]*models.ComponentFirmwareSetMap, error) {
	var m []*models.ComponentFirmwareSetMap

	// validate component firmware version does not already exist in map
	query := firmwareSet.FirmwareSetComponentFirmwareSetMaps(qm.Where("firmware_id=?", firmwareUUID))

	m, err := query.All(ctx, r.DB)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return m, errors.Wrap(errDBErr, err.Error())
	}

	return m, nil
}

func (r *Router) firmwareSetUpdateTx(ctx context.Context, newValues *models.ComponentFirmwareSet, attributes models.AttributesFirmwareSetSlice, firmwareUUIDs []uuid.UUID) error {
	// being transaction to update a firmware set and its references
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// nolint:errcheck // TODO(joel): log error
	defer tx.Rollback()

	currentValues, err := models.FindComponentFirmwareSet(ctx, tx, newValues.ID)
	if err != nil {
		return err
	}

	// update name column
	if newValues.Name != "" && newValues.Name != currentValues.Name {
		currentValues.Name = newValues.Name
	}

	if _, err := currentValues.Update(ctx, tx, boil.Infer()); err != nil {
		return err
	}

	// retrieve referenced firmware set attributes
	attrs, err := currentValues.FirmwareSetAttributesFirmwareSets().All(ctx, tx)
	if err != nil {
		return err
	}

	// remove current referenced firmware set attributes
	_, err = attrs.DeleteAll(ctx, tx)
	if err != nil {
		return err
	}

	// add new firmware set attributes
	if err := newValues.AddFirmwareSetAttributesFirmwareSets(ctx, tx, true, attributes...); err != nil {
		return err
	}

	// add new firmware references
	for _, id := range firmwareUUIDs {
		m := models.ComponentFirmwareSetMap{FirmwareSetID: newValues.ID, FirmwareID: id.String()}

		err := m.Insert(ctx, tx, boil.Infer())
		if err != nil {
			return err
		}
	}

	// commit
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *Router) serverComponentFirmwareSetRemoveFirmware(c *gin.Context) {
	firmwareSet, err := r.componentFirmwareSetFromParams(c)
	if err != nil {
		badRequestResponse(c, "invalid payload: ComponentFirmwareSet{}", err)
		return
	}

	var payload ComponentFirmwareSetRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		badRequestResponse(c, "invalid payload: ComponentFirmwareSet{}", err)
		return
	}

	// firmware set ID is expected for
	if payload.ID == uuid.Nil {
		badRequestResponse(
			c,
			"",
			errors.Wrap(errComponentFirmwareSetRequest, "expected a valid firmware set ID in payload, got none"),
		)

		return
	}

	// firmware set ID URL param is expected to match payload firmware set ID
	if payload.ID.String() != firmwareSet.ID {
		badRequestResponse(
			c,
			"",
			errors.Wrap(errComponentFirmwareSetRequest, "firmware set ID does not match payload ID attribute"),
		)

		return
	}

	// identify firmware set - firmware mappings for removal
	removeMappings := []*models.ComponentFirmwareSetMap{}

	for _, firmwareUUID := range payload.ComponentFirmwareUUIDs {
		// parse uuid
		firmwareUUIDParsed, err := uuid.Parse(firmwareUUID)
		if err != nil {
			badRequestResponse(
				c,
				"invalid firmware UUID: "+firmwareUUID,
				errors.Wrap(errComponentFirmwareSetRequest, err.Error()),
			)

			return
		}

		// validate firmware is part of set
		setMap, err := r.firmwareSetMap(c.Request.Context(), firmwareSet, firmwareUUIDParsed)
		if err != nil {
			dbErrorResponse(c, err)

			return
		}

		if len(setMap) == 0 {
			badRequestResponse(
				c,
				"invalid firmware UUID: "+firmwareUUID,
				errors.Wrap(
					errComponentFirmwareSetRequest,

					fmt.Sprintf("firmware set '%s' does not contain firmware '%s'", firmwareSet.Name, firmwareUUID),
				),
			)

			return
		}

		removeMappings = append(removeMappings, setMap...)
	}

	err = r.firmwareSetDeleteMappingTx(c.Request.Context(), firmwareSet, removeMappings)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	deletedResponse(c)
}

func (r *Router) serverComponentFirmwareSetDelete(c *gin.Context) {
	dbFirmware, err := r.componentFirmwareSetFromParams(c)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	if _, err = dbFirmware.Delete(c.Request.Context(), r.DB); err != nil {
		dbErrorResponse(c, err)
		return
	}

	deletedResponse(c)
}

func (r *Router) componentFirmwareSetFromParams(c *gin.Context) (*models.ComponentFirmwareSet, error) {
	u, err := r.parseUUID(c)
	if err != nil {
		return nil, err
	}

	if u == uuid.Nil {
		return nil, errors.Wrap(errComponentFirmwareSetRequest, "expected a valid firmware set UUID")
	}

	firmwareSet, err := models.FindComponentFirmwareSet(c.Request.Context(), r.DB, u.String())
	if err != nil {
		return nil, err
	}

	return firmwareSet, nil
}

func (r *Router) firmwareSetDeleteMappingTx(ctx context.Context, firmwareSet *models.ComponentFirmwareSet, removeMappings []*models.ComponentFirmwareSetMap) error {
	// being transaction to insert a new firmware set and its mapping
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// nolint:errcheck // (joel) once a logger is made available, this tx rollback can be logged.
	defer tx.Rollback()

	for _, mapping := range removeMappings {
		if _, err := mapping.Delete(ctx, r.DB); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
