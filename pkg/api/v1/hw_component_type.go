package hollow

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go.metalkube.net/hollow/internal/db"
)

// HardwareComponentType provides a way to group hardware components by the type
type HardwareComponentType struct {
	UUID uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
}

func (t *HardwareComponentType) fromDBModel(dbT db.HardwareComponentType) error {
	t.UUID = dbT.ID
	t.Name = dbT.Name

	return nil
}

func (t *HardwareComponentType) toDBModel() (*db.HardwareComponentType, error) {
	dbT := &db.HardwareComponentType{
		ID:   t.UUID,
		Name: t.Name,
	}

	return dbT, nil
}

func hardwareComponentTypeCreate(c *gin.Context) {
	var t HardwareComponentType
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid hardware component type",
			"error":   err.Error(),
		})

		return
	}

	dbT, err := t.toDBModel()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid hardware component type", "error": err.Error()})
		return
	}

	if err := db.CreateHardwareComponentType(dbT); err != nil {
		c.JSON(http.StatusInternalServerError, newErrorResponse("failed to create hardware component type", err))
		return
	}

	c.JSON(http.StatusOK, createdResponse(dbT.ID))
}

func hardwareComponentTypeList(c *gin.Context) {
	dbFilter := &db.HardwareComponentTypeFilter{
		Name: c.Query("name"),
	}

	dbTypes, err := db.GetHardwareComponentTypes(dbFilter)
	if err != nil {
		dbQueryFailureResponse(c, err)
		return
	}

	types := []HardwareComponentType{}

	for _, dbT := range dbTypes {
		t := HardwareComponentType{}
		if err := t.fromDBModel(dbT); err != nil {
			failedConvertingToVersioned(c, err)
			return
		}

		types = append(types, t)
	}

	c.JSON(http.StatusOK, types)
}
