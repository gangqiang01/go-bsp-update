package router

import (
	"github.com/edgehook/ithings/webserver/api"
	v1 "github.com/edgehook/ithings/webserver/api/v1"
	"github.com/edgehook/ithings/webserver/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors())
	r.POST("/login", api.Login)
	apiv1 := r.Group("/v1")
	{
		apiv1.POST("/upload/chunk", v1.UploadChunkHandler)
		apiv1.GET("/bsp/byPage", v1.GetBspByPage)
		apiv1.DELETE("/bsp/:id", v1.DeleteBsp)
	}
	return r

}
