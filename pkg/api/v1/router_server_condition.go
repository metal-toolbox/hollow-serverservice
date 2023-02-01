package serverservice

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"

	"go.hollow.sh/serverservice/internal/models"
)

// serverConditionList lists all conditions associated with a server.
func (r *Router) serverConditionList(c *gin.Context) {
	pager := parsePagination(c)

	mods := []qm.QueryMod{
		models.ServerConditionWhere.ServerID.EQ(c.Param("uuid")),
		qm.InnerJoin(fmt.Sprintf("%s as t on t.%s = %s.%s",
			models.TableNames.ServerConditionTypes,
			models.ServerConditionTypeColumns.ID,
			models.TableNames.ServerConditions,
			models.ServerConditionColumns.ServerConditionTypeID,
		)),
		qm.InnerJoin(fmt.Sprintf("%s as st on st.%s = %s.%s",
			models.TableNames.ServerConditionStatusTypes,
			models.ServerConditionStatusTypeColumns.ID,
			models.TableNames.ServerConditions,
			models.ServerConditionColumns.ServerConditionStatusTypeID,
		)),
		qm.Load(models.ServerConditionRels.ServerConditionType),
		qm.Load(models.ServerConditionRels.ServerConditionStatusType),
	}

	dbC, err := models.ServerConditions(mods...).All(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	count, err := models.ServerConditions(mods...).Count(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	conditions := make([]ServerCondition, 0, len(dbC))

	for _, dbCElem := range dbC {
		conditions = append(conditions, ServerCondition{
			Slug:         dbCElem.R.ServerConditionType.Slug,
			Parameters:   json.RawMessage(dbCElem.Parameters),
			StatusOutput: []byte(dbCElem.StatusOutput),
			Status:       dbCElem.R.ServerConditionStatusType.Slug,
			CreatedAt:    dbCElem.CreatedAt.Time,
			UpdatedAt:    dbCElem.UpdatedAt.Time,
		})
	}

	pd := paginationData{
		pageCount:  len(conditions),
		totalCount: count,
		pager:      pager,
	}

	listResponse(c, conditions, pd)
}

// serverConditionGet lists the a condition associated with a server, matched by the condition slug parameter.
func (r *Router) serverConditionGet(c *gin.Context) {
	mods := []qm.QueryMod{
		models.ServerConditionWhere.ServerID.EQ(c.Param("uuid")),
		qm.InnerJoin(fmt.Sprintf("%s as t on t.%s = %s.%s",
			models.TableNames.ServerConditionTypes,
			models.ServerConditionTypeColumns.ID,
			models.TableNames.ServerConditions,
			models.ServerConditionColumns.ServerConditionTypeID,
		)),
		qm.InnerJoin(fmt.Sprintf("%s as st on st.%s = %s.%s",
			models.TableNames.ServerConditionStatusTypes,
			models.ServerConditionStatusTypeColumns.ID,
			models.TableNames.ServerConditions,
			models.ServerConditionColumns.ServerConditionStatusTypeID,
		)),
		qm.Where(fmt.Sprintf("t.%s=?", models.ServerConditionTypeColumns.Slug), c.Param("slug")),
		qm.Load(models.ServerConditionRels.ServerConditionType),
		qm.Load(models.ServerConditionRels.ServerConditionStatusType),
	}

	dbC, err := models.ServerConditions(mods...).One(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	condition := &ServerCondition{
		Slug:         dbC.R.ServerConditionType.Slug,
		Parameters:   json.RawMessage(dbC.Parameters),
		StatusOutput: []byte(dbC.StatusOutput),
		Status:       dbC.R.ServerConditionStatusType.Slug,
		CreatedAt:    dbC.CreatedAt.Time,
		UpdatedAt:    dbC.UpdatedAt.Time,
	}

	itemResponse(c, condition)
}

// serverConditionDelete removes a condition associated with server matched by the condition slug parameter.
func (r *Router) serverConditionDelete(c *gin.Context) {
	mods := []qm.QueryMod{
		models.ServerConditionWhere.ServerID.EQ(c.Param("uuid")),
		qm.InnerJoin(fmt.Sprintf("%s as t on t.%s = %s.%s",
			models.TableNames.ServerConditionTypes,
			models.ServerConditionTypeColumns.ID,
			models.TableNames.ServerConditions,
			models.ServerConditionColumns.ServerConditionTypeID,
		)),
		qm.Where(fmt.Sprintf("t.%s=?", models.ServerConditionTypeColumns.Slug), c.Param("slug")),
	}

	dbC, err := models.ServerConditions(mods...).One(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	if _, err = dbC.Delete(c.Request.Context(), r.DB); err != nil {
		dbErrorResponse(c, err)
		return
	}

	deletedResponse(c)
}

// serverConditionUpsert inserts or updates a server condition matched by the condition slug.
func (r *Router) serverConditionUpsert(c *gin.Context) {
	srvUUID, err := r.parseUUID(c)
	if err != nil {
		return
	}

	conditionSlug := c.Param("slug")

	exists, err := models.ServerExists(c.Request.Context(), r.DB, srvUUID.String())
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	if !exists {
		notFoundResponse(c, "server not found")
		return
	}

	conditionType, err := models.ServerConditionTypes(models.ServerConditionTypeWhere.Slug.EQ(conditionSlug)).One(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	var newValue ServerCondition
	if err := c.ShouldBindJSON(&newValue); err != nil {
		badRequestResponse(c, "invalid ServerCondition{} payload", err)
		return
	}

	conditionStatusType, err := models.ServerConditionStatusTypes(models.ServerConditionStatusTypeWhere.Slug.EQ(newValue.Status)).One(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	condition := models.ServerCondition{
		ServerID:                    srvUUID.String(),
		ServerConditionTypeID:       conditionType.ID,
		ServerConditionStatusTypeID: conditionStatusType.ID,
		Parameters:                  types.JSON(newValue.Parameters),
		StatusOutput:                types.JSON(newValue.StatusOutput),
	}

	// set a valid JSON default status output
	if len(newValue.StatusOutput) == 0 {
		condition.StatusOutput = []byte(`{}`)
	}

	err = condition.Upsert(
		c.Request.Context(),
		r.DB,
		true,
		// match records by server ID and server condition type ID
		[]string{
			models.ServerConditionColumns.ServerID,
			models.ServerConditionColumns.ServerConditionTypeID,
		},
		// update columns
		boil.Whitelist(
			models.ServerConditionColumns.ServerConditionStatusTypeID,
			models.ServerConditionColumns.Parameters,
			models.ServerConditionColumns.StatusOutput,
		),
		// insert columns
		boil.Whitelist(
			models.ServerConditionColumns.ServerID,
			models.ServerConditionColumns.ServerConditionTypeID,
			models.ServerConditionColumns.ServerConditionStatusTypeID,
			models.ServerConditionColumns.Parameters,
			models.ServerConditionColumns.StatusOutput,
			models.ServerConditionColumns.CreatedAt,
			models.ServerConditionColumns.UpdatedAt,
		),
	)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	updatedResponse(c, "")
}
