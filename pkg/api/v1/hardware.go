package hollow

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go.metalkube.net/hollow/internal/db"
)

// Hardware represents a piece of hardware in a facility. These are the
// details of the physical hardware
type Hardware struct {
	UUID               uuid.UUID           `json:"uuid"`
	FacilityCode       string              `json:"facility"`
	Attributes         []Attributes        `json:"attributes"`
	HardwareComponents []HardwareComponent `json:"hardware_components"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
}

// Attributes provide the ability to apply namespaced settings to an entity.
// For example hardware could have attributes in the `com.equinixmetal.api` namespace
// that represents equinix metal specific attributes that are stored in the API.
// The namespace is meant to define who owns the schema and values.
type Attributes struct {
	Namespace string          `json:"namespace"`
	Values    json.RawMessage `json:"values"`
}

func (a *Attributes) fromDBModel(dbA db.Attributes) error {
	a.Namespace = dbA.Namespace
	a.Values = json.RawMessage(dbA.Values)

	return nil
}

func convertDBAttributes(dbAttrs []db.Attributes) ([]Attributes, error) {
	attrs := []Attributes{}

	for _, dbA := range dbAttrs {
		a := Attributes{}
		if err := a.fromDBModel(dbA); err != nil {
			return nil, err
		}

		attrs = append(attrs, a)
	}

	return attrs, nil
}

func (h *Hardware) fromDBModel(dbH *db.Hardware) error {
	var err error

	h.UUID = dbH.ID
	h.FacilityCode = dbH.FacilityCode
	h.CreatedAt = dbH.CreatedAt
	h.UpdatedAt = dbH.UpdatedAt

	h.Attributes, err = convertDBAttributes(dbH.Attributes)
	if err != nil {
		return err
	}

	h.HardwareComponents, err = convertDBHardwareComponents(dbH.HardwareComponents)
	if err != nil {
		return err
	}

	return nil
}

func (h *Hardware) toDBModel() (*db.Hardware, error) {
	dbC := &db.Hardware{
		ID:           h.UUID,
		FacilityCode: h.FacilityCode,
	}

	for _, hc := range h.HardwareComponents {
		c, err := hc.toDBModel()
		if err != nil {
			return nil, err
		}

		dbC.HardwareComponents = append(dbC.HardwareComponents, *c)
	}

	// attrs, err := convertDBAttributes(dbC.Attributes)
	// if err != nil {
	// 	return nil, err
	// }

	// c.Attributes = attrs

	return dbC, nil
}

// func convertDBHardware(d db.Hardware) (Hardware, error) {
// 	var bc []BIOSConfig

// 	for _, dbc := range d {
// 		var b BIOSConfig
// 		if err := b.fromDBModel(&dbc); err != nil {
// 			return nil, err
// 		}

// 		bc = append(bc, b)
// 	}

// 	return bc, nil
// }

func hardwareGet(c *gin.Context) {
	hwUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse("Invalid hardware UUID", err))
		return
	}

	dbHW, err := db.FindHardwareByUUID(hwUUID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, notFoundResponse())
			return
		}

		c.JSON(http.StatusInternalServerError, newErrorResponse("failed fetching records from datastore", err))

		return
	}

	hw := &Hardware{}

	if err = hw.fromDBModel(dbHW); err != nil {
		c.JSON(http.StatusInternalServerError, newErrorResponse("failed parsing the results", err))
		return
	}

	c.JSON(http.StatusOK, hw)
}

func hardwareCreate(c *gin.Context) {
	var hw Hardware
	if err := c.ShouldBindJSON(&hw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid hardware",
			"error":   err.Error(),
		})

		return
	}

	dbHW, err := hw.toDBModel()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Hardware", "error": err.Error})
		return
	}

	if err := db.CreateHardware(*dbHW); err != nil {
		c.JSON(http.StatusInternalServerError, newErrorResponse("Failed to create Hardware", err))
		return
	}

	c.JSON(http.StatusOK, createdResponse())
}
