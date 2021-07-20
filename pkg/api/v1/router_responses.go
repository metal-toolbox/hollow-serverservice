package hollow

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type serverResponse struct {
	Message string     `json:"message"`
	UUID    *uuid.UUID `json:"uuid,omitempty"`
	Error   string     `json:"error,omitempty"`
}

func newErrorResponse(m string, err error) *serverResponse {
	return &serverResponse{
		Message: m,
		Error:   err.Error(),
	}
}

func notFoundResponse(c *gin.Context, err error) {
	c.JSON(http.StatusNotFound, newErrorResponse("resource not found", err))
}

func createdResponse(c *gin.Context, u *uuid.UUID) {
	c.JSON(http.StatusOK, &serverResponse{Message: "resource created", UUID: u})
}

func deletedResponse(c *gin.Context) {
	c.JSON(http.StatusOK, &serverResponse{Message: "resource deleted"})
}

func dbFailureResponse(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, newErrorResponse("datastore error", err))
}

func failedConvertingToVersioned(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, newErrorResponse("failed parsing the datastore results", err))
}
