package hollow

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go.metalkube.net/hollow/internal/db"
)

// Hardware provides a versioned representation of hardware from the datastore
type Hardware struct {
	UUID         uuid.UUID `json:"uuid"`
	FacilityCode string    `json:"facility"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func hardwareList(c *gin.Context) {
	hw, err := db.HardwareList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
	}

	c.JSON(http.StatusOK, hw)
}
