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
		apiv1.GET("/system/process", v1.UploadProcessHandler)
	}
	return r

}
