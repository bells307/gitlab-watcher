package err_resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, err error) {
	resp := ErrorResponse{
		Message: err.Error(),
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": resp})
}
