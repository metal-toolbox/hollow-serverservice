package hollow

import (
	"encoding/json"
	"time"

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

// func biosConfigCreate(c *gin.Context) {
// 	var va VersionedAttributes
// 	if err := c.ShouldBindJSON(&va); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": "invalid versioned attributes",
// 			"error":   err.Error(),
// 		})

// 		return
// 	}

// 	dbVA, err := va.toDBModel()
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid versioned attributes", "error": err.Error})
// 		return
// 	}

// 	// ensure the hardware for the UUID exist already
// 	if _, err := db.FindOrCreateHardwareByUUID(dbVA.EntityID); err != nil {
// 		c.JSON(http.StatusInternalServerError, newErrorResponse("failed ensuring hardware with uuid exists", err))
// 		return
// 	}

// 	if err := db.VersionedAttributesCreate(dbVA); err != nil {
// 		c.JSON(http.StatusInternalServerError, newErrorResponse("failed to create versioned attributes", err))
// 		return
// 	}

// 	c.JSON(http.StatusOK, createdResponse(dbVA.ID))
// }

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

func convertToDBVersionedAttributes(attrs []VersionedAttributes) ([]db.VersionedAttributes, error) {
	dbVerAttrs := []db.VersionedAttributes{}

	for _, a := range attrs {
		dbVA, err := a.toDBModel()
		if err != nil {
			return nil, err
		}

		dbVerAttrs = append(dbVerAttrs, *dbVA)
	}

	return dbVerAttrs, nil
}

func convertFromDBVersionedAttributes(dbAttrs []db.VersionedAttributes) ([]VersionedAttributes, error) {
	attrs := []VersionedAttributes{}

	for _, dbA := range dbAttrs {
		a := VersionedAttributes{}
		if err := a.fromDBModel(dbA); err != nil {
			return nil, err
		}

		attrs = append(attrs, a)
	}

	return attrs, nil
}
