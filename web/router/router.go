package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iahfdoa/crawlsForBeauty/web/api"
)

func Router() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	v1.GET("/get", api.ImageController)
	return r
}
