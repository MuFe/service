package brand

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"mufe_service/camp/cache"
	"mufe_service/camp/dataUtil"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/handler"
	"mufe_service/camp/jwt"
	"mufe_service/camp/server"
	pb "mufe_service/jsonRpc"
	servicemanager "mufe_service/manager"
)

func init() {
	server.Post("/adminAccount/brand/login", login)
	server.Post("/adminAccount/brand/getInfo", handler.BrandAdminLogin, getUser)         // 获取用户信息
	server.Post("/adminAccount/brand/logout", logout)
	server.Get("/brandIndex/index", getIndex)
	server.Get("/brandFinance/finance",  getFinance)
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
	userResult, err := servicemanager.GetAdminUserService().BrandLogin(c, &pb.AdminLoginRequest{Phone: params.Phone, Pass: params.Pass,BgID: enum.BRAND_ADMIN_GROUP,})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	token, err := jwt.GenerateAdminJwt(enum.BRAND_ADMIN_GROUP,cache.BrandAgentToken,userResult)
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

	adminData, ok := c.MustGet(handler.BrandAdminData).(*jwt.AdminClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	businessInfo, err := cache.GetBusinessInfo(adminData.BusinessId)
	list := make([]string, 0)
	permission, err := servicemanager.GetPermissionService().GetUserPagePermission(c, &pb.GetUserPagePermissionRequest{Uid: adminData.Uid, BusinessGroupId: enum.BRAND_ADMIN_GROUP, BusinessId: adminData.BusinessId})
	if err == nil {
		for _, info := range permission.PermissionName {
			list = append(list, info)
		}
	}
	user := dataUtil.ParseAdminUser(businessInfo,adminData)
	user.Roles = list
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(user))
}


func logout(c *gin.Context) {
	var token = c.GetHeader(jwt.AdminAuthHeader)
	if token != "" {
		handler.AdminLogin(c)
		userData, ok := c.MustGet(handler.AdminData).(*jwt.AdminClaims)
		if ok {
			_ = cache.DeleteAdminToken(userData.Uid,cache.BrandAgentToken)
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
}


type getBrandIndexGoodsResponse struct {
	InSale    int64 `json:"inSale"`
	FewStock  int64 `json:"fewStock"`
	Withdrawn int64 `json:"withdrawn"`
	Draft     int64 `json:"draft"`
}

type data struct {
	Title string `json:"title"`
	Data  string `json:"data"`
	Per   int64  `json:"per"`
	IsUp  bool   `json:"isUp"`
}

type sdr struct {
	Data      []data `json:"data"` //0成交额，1订单数，2浏览量，3客单价，4退款金额，5退款单数
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
}

type order struct {
	Unpaid int64 `json:"unpaid"`
	Unsend int64 `json:"unsend"`
	Refund int64 `json:"refund"`
	Send   int64 `json:"send"`
	All    int64 `json:"all"`
}


type getBrandIndexResponse struct {
	SDR   []sdr                      `json:"sdr"` //0日报，1周报，2月报
	Index getBrandIndexSDRResponse   `json:"index"`
	List  []order                    `json:"order"` //0会员，1批发
	Good  getBrandIndexGoodsResponse `json:"good"`
}

type getBrandIndexSDRResponse struct {
	Quality       float64 `json:"quality"`
	QualityBeyond int64   `json:"qualityBeyond"`
	Express       float64 `json:"express"`
	ExpressBeyond int64   `json:"expressBeyond"`
	Service       float64 `json:"service"`
	ServiceBeyond int64   `json:"serviceBeyond"`
	Stars         float64 `json:"stars"`
}

func getIndex(c *gin.Context) {


	var result getBrandIndexResponse

	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
	return
}

type getFinanceRsp struct {
	WithdrawPer        int64 `json:"withdrawPer"`
	Withdraw           int64 `json:"withdraw"`
	WithdrawNotArrival int64 `json:"withdrawNotArrival"`
	Withdrawn          int64 `json:"withdrawn"`
	NotArrival         int64 `json:"notArrival"`
}

func getFinance(c *gin.Context) {
	var result getFinanceRsp

	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
}
