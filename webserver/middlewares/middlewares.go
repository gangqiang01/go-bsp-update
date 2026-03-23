package middlewares

import (
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/edgehook/ithings/webserver/api/jwt"
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

		path := c.FullPath()
		if strings.Contains(path, "login") ||
			strings.Contains(path, "upload") {
			c.Next()
			return
		}
		token := GetToken(c)

		if strings.Contains(path, "v1") && !VerifyToken(token) {
			responce.FailWithCodeAndMessage(401, "illegal user", c)
			//stop context
			c.Abort()
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

func GetToken(c *gin.Context) string {
	token := c.Request.Header.Get("accesstoken")
	if token == "" {
		token = c.Request.Header.Get("Authorization")
		if token != "" {
			if strings.Contains(token, "bearer") || strings.Contains(token, "Bearer") {
				if len(token) > 7 {
					token = token[7:]
				}
			}
		}
	}

	return token
}

func VerifyToken(token string) bool {
	if _, err := jwt.ParseCliamsToken(token); err != nil {
		return false
	}
	return true
}
