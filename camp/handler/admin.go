package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"mufe_service/camp/cache"
	"mufe_service/camp/errcode"
	"mufe_service/camp/jwt"
)

// AdminData 账户信息
var AdminData = "AdminData"

// AdminLogin 后台登录校验
func AdminLogin(c *gin.Context) {
	claims, err := jwt.CheckAdminJwt(cache.AgentToken,c.GetHeader(jwt.AdminAuthHeader))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.Set(AdminData, claims)
}


// AdminData 账户信息
var BrandAdminData = "BrandAdminData"

// AdminLogin 后台登录校验
func BrandAdminLogin(c *gin.Context) {
	claims, err := jwt.CheckAdminJwt(cache.BrandAgentToken,c.GetHeader(jwt.BrandAdminAuthHeader))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.Set(BrandAdminData, claims)
}
