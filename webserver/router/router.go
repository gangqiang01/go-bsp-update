package router

import (
	"net/http"
	"os"
	"path"
	"path/filepath"

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
		apiv1.POST("/system/reboot", v1.RebootHandler)
		apiv1.GET("/system/process", v1.UploadProcessHandler)
		apiv1.GET("/system/isUpdate", v1.IsUpdating)
	}

	//web

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	vueAssetsRoutePath := filepath.Join(dir, "frontend")
	r.StaticFile("/", path.Join(vueAssetsRoutePath, "index.html"))           // 指定资源文件 url.  127.0.0.1/ 这种
	r.StaticFile("/fav32.png", path.Join(vueAssetsRoutePath, "fav32.png"))   // 127.0.0.1/favicon.ico
	r.StaticFS("/assets", http.Dir(path.Join(vueAssetsRoutePath, "assets"))) // 以 assets 为前缀的 url
	r.StaticFS("/static", http.Dir(path.Join(vueAssetsRoutePath, "static")))
	return r

}
