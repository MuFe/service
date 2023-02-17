package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/jwt"
)

// UserData 账户信息
var UserData = "UserData"

// UserLogin 用户登录校验
func UserLogin(c *gin.Context) {
	claims, err := jwt.CheckUserJwt(c.GetHeader(jwt.AuthHeader))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.Set(UserData, claims)
}


// 老师权限校验
func TeacherCheck(c *gin.Context) {
	userData, err := jwt.CheckUserJwt(c.GetHeader(jwt.AuthHeader))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if userData.Identity!=enum.TEACHER_TYPE{
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorUnauthorized)
		return
	}
}

// 游客权限校验
func DefaultCheck(c *gin.Context) {
	userData, err := jwt.CheckUserJwt(c.GetHeader(jwt.AuthHeader))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if userData.Identity!=enum.DEFAULT_TYPE {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorUnauthorized)
		return
	}
}

// 教练权限校验
func CoachCheck(c *gin.Context) {
	userData, err := jwt.CheckUserJwt(c.GetHeader(jwt.AuthHeader))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if userData.Identity!=enum.INSTITUTION_TYPE{
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorUnauthorized)
		return
	}
}

