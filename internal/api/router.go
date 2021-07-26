package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	v1 "github.com/tinywell/baas/internal/api/v1"
)

// AddRouter ...
func AddRouter(r *gin.Engine) {

	r.Use(cors.Default())

	apiRoot := r.Group("/api")

	apiv1 := apiRoot.Group("/v1")
	{
		apiv1.POST("/network/init", v1.Network.Init)
		apiv1.GET("/network", v1.Network.Info)
	}
}
