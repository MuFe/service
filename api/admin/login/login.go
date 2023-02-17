package login

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"mufe_service/camp/cache"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/handler"
	"mufe_service/camp/jwt"
	"mufe_service/camp/server"
	"mufe_service/camp/utils"
	pb "mufe_service/jsonRpc"
	servicemanager "mufe_service/manager"
	"strconv"
	"time"
)

func init() {
	server.Post("/adminAccount/login", login)
	server.Post("/adminAccount/getInfo", handler.AdminLogin, getUser)         // 获取用户信息
	//server.Post("/adminUser/getUserList", handler.AdminLogin, getUserList) // 获取用户信息
	//server.Post("/adminUser/permissionList", handler.AdminLogin, permissionList) // 获取用户信息
	server.Post("/adminAccount/logout", logout)
	server.Post("/adminAccount/getToken", getToken)

}

type UserInfo struct {
	Name       string `json:"name"`
	Head       string `json:"head"`
	Phone      string `json:"phone"`
	Sex        int64  `json:"sex"`
	Uid        int64  `json:"uid"`
	OpenId     string `json:"open_id"`
	VipStatus  int64  `json:"vip_status"`
	No         string `json:"no"`
	InviteCode string `json:"invite_code"`
	BusinessId int64  `json:"business_id"`
	AgentId    int64  `json:"agent_id"`
}

func login(c *gin.Context) {
	type query struct {
		Phone string `form:"phone" json:"phone" `
		Pass  string  `form:"pass"  json:"pass" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Phone==""||params.Pass==""{
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	userResult, err := servicemanager.GetAdminUserService().Login(c, &pb.AdminLoginRequest{Phone: params.Phone, Pass: params.Pass})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	token, err := jwt.GenerateAdminJwt(enum.SYSTEM_ADMIN_GROUP,cache.AgentToken,userResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(token))
}

func getUser(c *gin.Context) {
	type UserResult struct {
		Head       string `json:"head" `
		Name       string  `json:"name" `
	}

	result:=UserResult{}

	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
}

func getUserList(c *gin.Context) {
	type query struct {
		Page int64 `form:"page"`
		Size int64 `form:"size"`
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	type UserResult struct {
		Head       string `json:"head" `
		Name       string  `json:"name" `
		Role       string `json:"role" `
		Id         int64  `json:"id" `
		Department string `json:"department" `
	}
	list := make([]interface{}, 0)
	data:=UserResult{
		Head:"",
		Name:"张三",
		Role:"权限1",
		Id:1,
		Department:"运营补",
	}
	list=append(list,data)
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(utils.CreateListCurrentResultReturn(1, params.Page, list)))
}
func permissionList(c *gin.Context) {
	type RoleResult struct {
		Name       string  `json:"name" `
		Id         int64  `json:"id" `
	}
	type Result struct {
		Name       string  `json:"name" `
		Role       []RoleResult `json:"role" `
		Id         int64  `json:"id" `
	}
	list := make([]interface{}, 0)
	data:=Result{
		Name:"运营部",
		Id:1,
	}
	roleList:=make([]RoleResult,0)
	roleList=append(roleList,RoleResult{
		Name:"管理员",
		Id:1,
	})
	roleList=append(roleList,RoleResult{
		Name:"普通",
		Id:2,
	})
	data.Role=roleList
	list=append(list,data)


	data1:=Result{
		Name:"IT部",
		Id:2,
	}
	roleList1:=make([]RoleResult,0)
	roleList1=append(roleList1,RoleResult{
		Name:"超级管理员",
		Id:3,
	})
	roleList1=append(roleList1,RoleResult{
		Name:"普通员工",
		Id:4,
	})
	data1.Role=roleList1
	list=append(list,data1)
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
}


func logout(c *gin.Context) {
	var token = c.GetHeader(jwt.AdminAuthHeader)
	if token != "" {
		handler.AdminLogin(c)
		userData, ok := c.MustGet(handler.AdminData).(*jwt.AdminClaims)
		if ok {
			_ = cache.DeleteAdminToken(userData.Uid,cache.AgentToken)
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
}


type QiniuInfo struct {
	Token  string `json:"token"`
	Host   string `json:"host"`
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
}

func getToken(c *gin.Context) {
	type query struct {
		Name string `form:"name"`
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	filenameWithSuffix := path.Base(params.Name)
	fileSuffix := path.Ext(filenameWithSuffix) //获取文件后缀
	encodeString := base64.StdEncoding.EncodeToString([]byte(utils.MD5(params.Name+strconv.FormatInt(time.Now().Unix(), 10)) + fileSuffix))
	result, err := servicemanager.GetQiniuService().GetToken(c, &pb.QiniuServiceRequest{Bucket:os.Getenv("IMG_BUCKET"), IsBase64: false})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(&QiniuInfo{Token: result.Token, Host: result.UploadHost, Name: encodeString, Prefix: os.Getenv("IMG_PREFIX")}))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}
