package middlewares

import (
	"net/http"
	"runtime/debug"

	responce "github.com/edgehook/ithings/webserver/types"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token, session, Content-Type, accesstoken, timeout, Srptoken, Origin, Accept")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			c.Header("Access-Control-Max-Age", "3600")
			// c.Header("Access-Control-Allow-Credentials", "true")
		}

		// 预检请求
		if method == "OPTIONS" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token, session, Content-Type, accesstoken, timeout, Srptoken, Origin, Accept")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		defer func() {
			if err := recover(); err != nil {
				klog.Errorf("WebServer error occurred at:%s", string(debug.Stack()))
				responce.FailWithMessage("Server error", c)
			}
		}()

		c.Next()
	}
}
