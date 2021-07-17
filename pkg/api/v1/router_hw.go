package hollow

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Router) hardwareList(c *gin.Context) {
	var params HardwareListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid filter",
			"error":   err.Error(),
		})
	}

	alp, err := parseQueryAttributesListParams(c, "attr")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid attributes list params",
			"error":   err.Error(),
		})
	}

	params.AttributeListParams = alp

	valp, err := parseQueryAttributesListParams(c, "ver_attr")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid versioned attributes list params",
			"error":   err.Error(),
		})
	}

	params.VersionedAttributeListParams = valp

	dbFilter := params.dbFilter()

	dbHW, err := r.Store.GetHardware(dbFilter)
	if err != nil {
		dbFailureResponse(c, err)
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

func (r *Router) hardwareGet(c *gin.Context) {
	dbHW, err := r.loadHardwareFromParams(c)
	if err != nil {
		return
	}

	var hw Hardware
	if err = hw.fromDBModel(*dbHW); err != nil {
		failedConvertingToVersioned(c, err)
		return
	}

	c.JSON(http.StatusOK, hw)
}

func (r *Router) hardwareCreate(c *gin.Context) {
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

	if err := r.Store.CreateHardware(dbHW); err != nil {
		c.JSON(http.StatusInternalServerError, newErrorResponse("failed to create hardware", err))
		return
	}

	createdResponse(c, &dbHW.ID)
}

func (r *Router) hardwareDelete(c *gin.Context) {
	dbHW, err := r.loadHardwareFromParams(c)
	if err != nil {
		return
	}

	if err = r.Store.DeleteHardware(dbHW); err != nil {
		c.JSON(http.StatusInternalServerError, newErrorResponse("failed deleting resource", err))
		return
	}

	deletedResponse(c)
}

func (r *Router) hardwareVersionedAttributesList(c *gin.Context) {
	hwUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid hardware uuid", "error": err.Error()})
		return
	}

	dbVA, err := r.Store.GetVersionedAttributes(hwUUID)
	if err != nil {
		dbFailureResponse(c, err)
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

func (r *Router) hardwareVersionedAttributesCreate(c *gin.Context) {
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
		failedConvertingToVersioned(c, err)
		return
	}

	hw, err := r.loadHardwareFromParams(c)
	if err != nil {
		return
	}

	err = r.Store.CreateVersionedAttributes(hw, dbVA)
	if err != nil {
		dbFailureResponse(c, err)
		return
	}

	createdResponse(c, &dbVA.ID)
}
