package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tinywell/baas/internal/model/response"
)

// OK ..
func OK(c *gin.Context, rsp *response.Response) {
	c.JSON(http.StatusOK, rsp)
}

// Fail ..
func Fail(c *gin.Context, rsp *response.Response) {
	c.JSON(http.StatusInternalServerError, rsp)
}

// FailWithStatus ..
func FailWithStatus(c *gin.Context, status int, rsp *response.Response) {
	c.JSON(status, rsp)
}
