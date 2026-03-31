package api

import (
	v1 "github.com/edgehook/ithings/common/types/v1"
	"github.com/edgehook/ithings/common/utils"
	"github.com/edgehook/ithings/webserver/api/jwt"
	responce "github.com/edgehook/ithings/webserver/types"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

func Login(c *gin.Context) {
	var auth v1.Auth
	if err := c.Bind(&auth); err != nil {
		responce.FailWithMessage("Parameter error", c)
		return
	}
	type resp struct {
		AccessToken string `form:"accessToken" json:"accessToken"`
	}
	err := utils.AuthUser(auth.Username, auth.Password)
	if err != nil {
		klog.Errorf("System user login error: %s", err.Error())
		responce.FailWithMessage("User authentication error", c)
		return

	}

	token, err := jwt.GenerateToken(auth.Username)
	if err != nil {
		responce.FailWithMessage("Grnerate token error", c)
		return
	}
	responce.OkWithData(&resp{
		AccessToken: token,
	}, c)
}
