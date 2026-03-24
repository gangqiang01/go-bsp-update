package router

import (
	v1 "github.com/edgehook/ithings/webserver/api/v1"
	"github.com/edgehook/ithings/webserver/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors())
	apiv1 := r.Group("/v1")
	{
		apiv1.POST("/upload/chunk", v1.UploadChunkHandler)
		apiv1.PUT("/system/reboot", v1.RebootHandler)
	}
	return r

}
