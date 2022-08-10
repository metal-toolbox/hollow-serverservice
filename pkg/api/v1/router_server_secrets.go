package serverservice

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.hollow.sh/serverservice/internal/dbtools"
	"go.hollow.sh/serverservice/internal/models"
)

func (r *Router) serverCredentialGet(c *gin.Context) {
	mods := []qm.QueryMod{
		models.ServerCredentialWhere.ServerID.EQ(c.Param("uuid")),
		qm.InnerJoin(fmt.Sprintf("%s as t on t.%s = %s.%s",
			models.TableNames.ServerCredentialTypes,
			models.ServerCredentialTypeColumns.ID,
			models.TableNames.ServerCredentials,
			models.ServerCredentialColumns.ServerCredentialTypeID,
		)),
		qm.Where(fmt.Sprintf("t.%s=?", models.ServerCredentialTypeColumns.Slug), c.Param("slug")),
		qm.Load(models.ServerCredentialRels.ServerCredentialType),
	}

	dbS, err := models.ServerCredentials(mods...).One(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	decryptedValue, err := dbtools.Decrypt(c.Request.Context(), r.SecretsKeeper, dbS.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ServerResponse{Message: "error decrypting value", Error: err.Error()})
		return
	}

	sID, err := uuid.Parse(dbS.ServerID)
	if err != nil {
		failedConvertingToVersioned(c, err)
		return
	}

	secret := &ServerCredential{
		ServerID:   sID,
		SecretType: dbS.R.ServerCredentialType.Slug,
		Username:   dbS.Username,
		Password:   decryptedValue,
		CreatedAt:  dbS.CreatedAt,
		UpdatedAt:  dbS.UpdatedAt,
	}

	itemResponse(c, secret)
}

func (r *Router) serverCredentialDelete(c *gin.Context) {
	mods := []qm.QueryMod{
		models.ServerCredentialWhere.ServerID.EQ(c.Param("uuid")),
		qm.InnerJoin(fmt.Sprintf("%s as t on t.%s = %s.%s",
			models.TableNames.ServerCredentialTypes,
			models.ServerCredentialTypeColumns.ID,
			models.TableNames.ServerCredentials,
			models.ServerCredentialColumns.ServerCredentialTypeID,
		)),
		qm.Where(fmt.Sprintf("t.%s=?", models.ServerCredentialTypeColumns.Slug), c.Param("slug")),
	}

	dbS, err := models.ServerCredentials(mods...).One(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	if _, err = dbS.Delete(c.Request.Context(), r.DB); err != nil {
		dbErrorResponse(c, err)
		return
	}

	deletedResponse(c)
}

func (r *Router) serverCredentialUpsert(c *gin.Context) {
	srvUUID, err := r.parseUUID(c)
	if err != nil {
		return
	}

	secretSlug := c.Param("slug")

	exists, err := models.ServerExists(c.Request.Context(), r.DB, srvUUID.String())
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	if !exists {
		notFoundResponse(c, "server not found")
		return
	}

	secretType, err := models.ServerCredentialTypes(models.ServerCredentialTypeWhere.Slug.EQ(secretSlug)).One(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	var newValue serverCredentialValues
	if err := c.ShouldBindJSON(&newValue); err != nil {
		badRequestResponse(c, "invalid server secret value", err)
		return
	}

	encryptedValue, err := dbtools.Encrypt(c.Request.Context(), r.SecretsKeeper, newValue.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ServerResponse{Message: "error encrypting secret value", Error: err.Error()})
		return
	}

	secret := models.ServerCredential{
		ServerCredentialTypeID: secretType.ID,
		ServerID:               srvUUID.String(),
		Password:               encryptedValue,
		Username:               newValue.Username,
	}

	err = secret.Upsert(
		c.Request.Context(),
		r.DB,
		true,
		// search for records by server id and type id to see if we need to update or insert
		[]string{models.ServerCredentialColumns.ServerID, models.ServerCredentialColumns.ServerCredentialTypeID},
		// For updates only set the new value and updated at
		boil.Whitelist(models.ServerCredentialColumns.Password, models.ServerCredentialColumns.UpdatedAt),
		// For inserts set server id, type id and value
		boil.Whitelist(
			models.ServerCredentialColumns.ServerID,
			models.ServerCredentialColumns.ServerCredentialTypeID,
			models.ServerCredentialColumns.Username,
			models.ServerCredentialColumns.Password,
			models.ServerCredentialColumns.CreatedAt,
			models.ServerCredentialColumns.UpdatedAt,
		),
	)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	updatedResponse(c, secretSlug)
}
