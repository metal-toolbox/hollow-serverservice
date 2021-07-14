package hollow

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type serverResponse struct {
	Message      string `json:"message"`
	ErrorDetails string `json:"error_details,omitempty"`
}

func newErrorResponse(m string, err error) *serverResponse {
	return &serverResponse{
		Message:      m,
		ErrorDetails: err.Error(),
	}
}

func notFoundResponse() *serverResponse {
	return &serverResponse{
		Message: "resource not found",
	}
}

func createdResponse() *serverResponse {
	return &serverResponse{
		Message: "created",
	}
}

func dbQueryFailureResponse(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, newErrorResponse("failed fetching records from datastore", err))
}

func failedConvertingToVersioned(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, newErrorResponse("failed parsing the datastore results", err))
}
