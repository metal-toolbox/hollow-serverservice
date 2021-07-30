package hollow

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ServerResponse represents the data that the server will return on any given call
type ServerResponse struct {
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
	UUID    *uuid.UUID  `json:"uuid,omitempty"`
	Item    interface{} `json:"item,omitempty"`
	Items   interface{} `json:"items,omitempty"`
}

func newErrorResponse(m string, err error) *ServerResponse {
	return &ServerResponse{
		Message: m,
		Error:   err.Error(),
	}
}

func badRequestResponse(c *gin.Context, message string, err error) {
	c.JSON(http.StatusBadRequest, newErrorResponse(message, err))
}

func notFoundResponse(c *gin.Context, err error) {
	c.JSON(http.StatusNotFound, newErrorResponse("resource not found", err))
}

func createdResponse(c *gin.Context, u *uuid.UUID) {
	c.JSON(http.StatusOK, &ServerResponse{Message: "resource created", UUID: u})
}

func deletedResponse(c *gin.Context) {
	c.JSON(http.StatusOK, &ServerResponse{Message: "resource deleted"})
}

func dbFailureResponse(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, newErrorResponse("datastore error", err))
}

func failedConvertingToVersioned(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, newErrorResponse("failed parsing the datastore results", err))
}

func listResponse(c *gin.Context, i interface{}) {
	c.JSON(http.StatusOK, &ServerResponse{Items: i})
}

func itemResponse(c *gin.Context, i interface{}) {
	c.JSON(http.StatusOK, &ServerResponse{Item: i})
}
