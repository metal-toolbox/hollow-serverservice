package hollow

import (
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

func (h *Hardware) fromDBModel(dbH db.Hardware) error {
	var err error

	h.UUID = dbH.ID
	h.FacilityCode = dbH.FacilityCode
	h.CreatedAt = dbH.CreatedAt
	h.UpdatedAt = dbH.UpdatedAt

	h.Attributes, err = convertFromDBAttributes(dbH.Attributes)
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

	attrs, err := convertToDBAttributes(h.Attributes)
	if err != nil {
		return nil, err
	}

	dbC.Attributes = attrs

	return dbC, nil
}

func hardwareGet(c *gin.Context) {
	hwUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse("invalid hardware UUID", err))
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

	if err = hw.fromDBModel(*dbHW); err != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid hardware", "error": err.Error})
		return
	}

	if err := db.CreateHardware(dbHW); err != nil {
		c.JSON(http.StatusInternalServerError, newErrorResponse("failed to create hardware", err))
		return
	}

	c.JSON(http.StatusOK, createdResponse(dbHW.ID))
}

func hardwareList(c *gin.Context) {
	var params HardwareListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid filter",
			"error":   err.Error(),
		})
	}

	alp, err := parseQueryAttributesListParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid attributes list params",
			"error":   err.Error(),
		})
	}

	params.AttributeListParams = alp
	dbFilter := params.dbFilter()

	dbHW, err := db.GetHardware(dbFilter)
	if err != nil {
		dbQueryFailureResponse(c, err)
		return
	}

	hw := []Hardware{}

	for _, dbH := range dbHW {
		h := Hardware{}
		if err := h.fromDBModel(dbH); err != nil {
			failedConvertingToVersioned(c, err)
			return
		}

		hw = append(hw, h)
	}

	c.JSON(http.StatusOK, hw)
}
