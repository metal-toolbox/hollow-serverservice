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

func (r *Router) serverSecretGet(c *gin.Context) {
	mods := []qm.QueryMod{
		models.ServerSecretWhere.ServerID.EQ(c.Param("uuid")),
		qm.InnerJoin(fmt.Sprintf("%s as t on t.%s = %s.%s",
			models.TableNames.ServerSecretTypes,
			models.ServerSecretTypeColumns.ID,
			models.TableNames.ServerSecrets,
			models.ServerSecretColumns.ServerSecretTypeID,
		)),
		qm.Where(fmt.Sprintf("t.%s=?", models.ServerSecretTypeColumns.Slug), c.Param("slug")),
		qm.Load(models.ServerSecretRels.ServerSecretType),
	}

	dbS, err := models.ServerSecrets(mods...).One(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	decryptedValue, err := dbtools.Decrypt(c.Request.Context(), r.SecretsKeeper, dbS.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ServerResponse{Message: "error decrypting secret value", Error: err.Error()})
		return
	}

	sID, err := uuid.Parse(dbS.ServerID)
	if err != nil {
		failedConvertingToVersioned(c, err)
		return
	}

	secret := &ServerSecret{
		ServerID:   sID,
		SecretType: dbS.R.ServerSecretType.Slug,
		Value:      decryptedValue,
		CreatedAt:  dbS.CreatedAt,
		UpdatedAt:  dbS.UpdatedAt,
	}

	itemResponse(c, secret)
}

func (r *Router) serverSecretDelete(c *gin.Context) {
	mods := []qm.QueryMod{
		models.ServerSecretWhere.ServerID.EQ(c.Param("uuid")),
		qm.InnerJoin(fmt.Sprintf("%s as t on t.%s = %s.%s",
			models.TableNames.ServerSecretTypes,
			models.ServerSecretTypeColumns.ID,
			models.TableNames.ServerSecrets,
			models.ServerSecretColumns.ServerSecretTypeID,
		)),
		qm.Where(fmt.Sprintf("t.%s=?", models.ServerSecretTypeColumns.Slug), c.Param("slug")),
	}

	dbS, err := models.ServerSecrets(mods...).One(c.Request.Context(), r.DB)
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

func (r *Router) serverSecretUpsert(c *gin.Context) {
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

	secretType, err := models.ServerSecretTypes(models.ServerSecretTypeWhere.Slug.EQ(secretSlug)).One(c.Request.Context(), r.DB)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	var newValue serverSecretValue
	if err := c.ShouldBindJSON(&newValue); err != nil {
		badRequestResponse(c, "invalid server secret value", err)
		return
	}

	encryptedValue, err := dbtools.Encrypt(c.Request.Context(), r.SecretsKeeper, newValue.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ServerResponse{Message: "error encrypting secret value", Error: err.Error()})
		return
	}

	secret := models.ServerSecret{
		ServerSecretTypeID: secretType.ID,
		ServerID:           srvUUID.String(),
		Value:              encryptedValue,
	}

	err = secret.Upsert(
		c.Request.Context(),
		r.DB,
		true,
		// search for records by server id and type id to see if we need to update or insert
		[]string{models.ServerSecretColumns.ServerID, models.ServerSecretColumns.ServerSecretTypeID},
		// For updates only set the new value and updated at
		boil.Whitelist(models.ServerSecretColumns.Value, models.ServerSecretColumns.UpdatedAt),
		// For inserts set server id, type id and value
		boil.Whitelist(
			models.ServerSecretColumns.ServerID,
			models.ServerSecretColumns.ServerSecretTypeID,
			models.ServerSecretColumns.Value,
			models.ServerSecretColumns.CreatedAt,
			models.ServerSecretColumns.UpdatedAt,
		),
	)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	updatedResponse(c, secretSlug)
}
