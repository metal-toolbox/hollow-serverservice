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

// BIOSConfig represents the BIOS config settings of a server at a given time
type BIOSConfig struct {
	HardwareUUID uuid.UUID       `json:"hardware_uuid" binding:"required"`
	ConfigValues json.RawMessage `json:"values" binding:"required"`
	CreatedAt    time.Time       `json:"created_at"`
}

func hardwareBIOSConfigList(c *gin.Context) {
	hwUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid hardware UUID", "error": err.Error()})
		return
	}

	bcl, err := db.BIOSConfigList(hwUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed fetching records from datastore", "error": err.Error()})
		return
	}

	l, err := dbSliceToBIOSConfig(bcl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed parsing results", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, l)
}

func biosConfigCreate(c *gin.Context) {
	var b BIOSConfig
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid BIOS Config",
			"error":   err.Error(),
		})

		return
	}

	dbc, err := b.toDBModel()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid BIOS Config", "error": err.Error})
		return
	}

	if err := db.CreateBIOSConfig(*dbc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create BIOS Config", "error": err.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "created"})
}

func (b *BIOSConfig) toDBModel() (*db.BIOSConfig, error) {
	dbc := &db.BIOSConfig{
		HardwareUUID: b.HardwareUUID,
		ConfigValues: datatypes.JSON(b.ConfigValues),
	}

	return dbc, nil
}

func (b *BIOSConfig) fromDBModel(dbc *db.BIOSConfig) error {
	b.HardwareUUID = dbc.HardwareUUID
	b.CreatedAt = dbc.CreatedAt
	b.ConfigValues = json.RawMessage(dbc.ConfigValues)

	return nil
}

func dbSliceToBIOSConfig(d []db.BIOSConfig) ([]BIOSConfig, error) {
	var bc []BIOSConfig

	for _, dbc := range d {
		var b BIOSConfig
		if err := b.fromDBModel(&dbc); err != nil {
			return nil, err
		}

		bc = append(bc, b)
	}

	return bc, nil
}
