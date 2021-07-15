package hollow

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"

	"go.metalkube.net/hollow/internal/db"
)

// VersionedAttributes represents a set of attributes of an entity at a given time
type VersionedAttributes struct {
	EntityType string          `json:"entity_type" binding:"required"`
	EntityUUID uuid.UUID       `json:"entity_uuid" binding:"required"`
	Namespace  string          `json:"namespace" binding:"required"`
	Values     json.RawMessage `json:"values" binding:"required"`
	CreatedAt  time.Time       `json:"created_at"`
}

func hardwareVersionedAttributesList(c *gin.Context) {
	hwUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid hardware uuid", "error": err.Error()})
		return
	}

	dbVA, err := db.VersionedAttributesList(hwUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed fetching records from datastore", "error": err.Error()})
		return
	}

	va := []VersionedAttributes{}

	for _, dbA := range dbVA {
		a := VersionedAttributes{}
		if err := a.fromDBModel(dbA); err != nil {
			failedConvertingToVersioned(c, err)
			return
		}

		va = append(va, a)
	}

	c.JSON(http.StatusOK, va)
}

func biosConfigCreate(c *gin.Context) {
	var va VersionedAttributes
	if err := c.ShouldBindJSON(&va); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid versioned attributes",
			"error":   err.Error(),
		})

		return
	}

	dbVA, err := va.toDBModel()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid versioned attributes", "error": err.Error})
		return
	}

	// ensure the hardware for the UUID exist already
	if _, err := db.FindOrCreateHardwareByUUID(dbVA.EntityID); err != nil {
		c.JSON(http.StatusInternalServerError, newErrorResponse("failed ensuring hardware with uuid exists", err))
		return
	}

	if err := db.VersionedAttributesCreate(dbVA); err != nil {
		c.JSON(http.StatusInternalServerError, newErrorResponse("failed to create versioned attributes", err))
		return
	}

	c.JSON(http.StatusOK, createdResponse(dbVA.ID))
}

func (a *VersionedAttributes) toDBModel() (*db.VersionedAttributes, error) {
	dbc := &db.VersionedAttributes{
		EntityType: a.EntityType,
		EntityID:   a.EntityUUID,
		Namespace:  a.Namespace,
		Values:     datatypes.JSON(a.Values),
	}

	return dbc, nil
}

func (a *VersionedAttributes) fromDBModel(dba db.VersionedAttributes) error {
	a.EntityType = dba.EntityType
	a.EntityUUID = dba.EntityID
	a.CreatedAt = dba.CreatedAt
	a.Namespace = dba.Namespace
	a.Values = json.RawMessage(dba.Values)

	return nil
}
