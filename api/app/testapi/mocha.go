package testapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"mufe_service/camp/cache"
	"mufe_service/camp/dataUtil"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/handler"
	"mufe_service/camp/jwt"
	"mufe_service/camp/server"
	"mufe_service/camp/utils"
	pb "mufe_service/jsonRpc"
	"mufe_service/manager"
)

func init() {
	if os.Getenv("MODEL")=="test"{
		server.Post("/appApi/test/mocha/sendCode", sendMochaCode)
		server.Post("/appApi/test/mocha/register", register)
		server.Post("/appApi/test/mocha/forgetPass", forgetPass)
		server.Post("/appApi/test/mocha/login", login)
		server.Post("/appApi/test/cancel",handler.UserLogin,mochaCancel)
		server.Post("/appApi/test/sendCode", sendCode)
	}

}



func sendMochaCode(c *gin.Context) {
	type codeResponse struct {
		Time int64 `json:"time"`
	}
	var result = codeResponse{}
	type query struct {
		Phone string `form:"phone"`
		Type  int64  `form:"type"`
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Phone == "" {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	if params.Type==3{
		phoneResult,err:=manager.GetUserService().PhoneCheck(c,&pb.PhoneCheckRequest{Phone:"m"+params.Phone})
		if err==nil{
			if phoneResult.Result!=0{
				c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorPhoneRepeat)
				return
			}
		}
	}
	ttl, err := cache.GetUserPhoneCodeLimitTime("m"+params.Phone, params.Type)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if ttl > 0 {
		result.Time = ttl
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorVerifiedCodeFrequent.SetData(result))
		return
	}

	code := utils.Get6Code()
	rpcSendSmsResult := &pb.SendResponse{}
	tempCode:="SMS_198585116"
	if params.Type==4{
		tempCode="SMS_198585115"
	}
	paramStr := fmt.Sprintf("{code:\"%s\"}", code)
	rpcSendSmsResult, err = manager.GetSendSmsService().SendSms(c, &pb.SendRequest{
		Key:"tT4IXgPraEzf1EUGV178kPmkBX75Cd",
		KeyId:"LTAI4GCZwryGBz5qnyWY1vpb",
		Phone:        params.Phone,
		ParamStr:     paramStr,
		SingName:     "万潮网络",
		TemplateCode: tempCode,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}

	if rpcSendSmsResult.Result != "OK" {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseMsg(rpcSendSmsResult.Result))
		return
	}

	result.Time, err = cache.SetUserPhoneCode("m"+params.Phone, code, params.Type)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
}



func mochaCancel(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	_, err := manager.GetUserService().EnterCancel(c, &pb.EnterCancelRequest{List:[]int64{ userData.Uid}})
	if err == nil {
		_ = cache.DeleteUserToken(userData.Uid)
		_ = cache.DeleteUserInfo(userData.Uid)
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("注销成功"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func register(c *gin.Context) {
	type query struct {
		Code  string `form:"code" json:"code" `
		Phone string `form:"phone" json:"phone" `
		Pass  string `form:"pass" json:"pass" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Phone == "" {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	if params.Code == "" {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	if params.Pass == "" {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	err := cache.CheckUserPhoneCode("m"+params.Phone, params.Code, enum.CodeRegisterType)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	inviteMap:=make(map[string]bool)
	code:=utils.BaseCode6(inviteMap)
	_, err = manager.GetUserService().Register(c, &pb.LoginRequest{Phone: "m"+params.Phone, Pass: params.Pass, Name: "用户"+code, Sign: ""})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func login(c *gin.Context) {
	type query struct {
		Code      string `form:"code" json:"code" `
		Phone     string `form:"phone" json:"phone" `
		Pass      string `form:"pass" json:"pass" `
		LoginType int64  `form:"login_type" json:"login_type" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Phone == "" {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	if params.LoginType == enum.PhoneType && params.Code == "" {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	if params.LoginType == enum.PhonePassType && params.Pass == "" {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	if params.LoginType == enum.PhoneType {
		err := cache.CheckUserPhoneCode("m"+params.Phone, params.Code, enum.CodePhoneLoginType)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
			return
		}
	}
	result, err := manager.GetUserService().Login(c, &pb.LoginRequest{Type: params.LoginType, Phone: "m"+params.Phone, Code: params.Code, Pass: params.Pass, Device: utils.GetHeaderFromKey(c, "device")})
	if err == nil {
		token, err := jwt.GenerateUserJwt(result.Uid, result.Identity, result.OpenId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
			return
		}
		u := dataUtil.ParseUserCache(result)
		_ = cache.SetUserInfo(result.Uid, u)
		type Result struct {
			Token    string `json:"token"`
			IsNew    bool   `json:"new"`
			Cancel   bool   `json:"cancel"`
			Identity int64  `json:"identity"`
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(Result{
			Token:    token,
			IsNew:    result.IsNew,
			Identity: result.Identity,
			Cancel:   result.CancelStatus,
		}))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func forgetPass(c *gin.Context) {
	type query struct {
		Code  string `form:"code" json:"code" `
		Phone string `form:"phone" json:"phone" `
		Pass  string `form:"pass" json:"pass" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Phone == "" {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	if params.Code == "" {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	if params.Pass == "" {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	err := cache.CheckUserPhoneCode("m"+params.Phone, params.Code, enum.CodeForgetType)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err = manager.GetUserService().UpdateUser(c, &pb.UpdateUserRequest{Phone:"m"+ params.Phone, Pass: params.Pass, Type: enum.UpdateUserPass})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func sendCode(c *gin.Context) {
	type codeResponse struct {
		Time int64 `json:"time"`
	}
	var result = codeResponse{}
	type query struct {
		Phone string `form:"phone"`
		Type  int64  `form:"type"`
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Phone == "" {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	if params.Type==3{
		phoneResult,err:=manager.GetUserService().PhoneCheck(c,&pb.PhoneCheckRequest{Phone:params.Phone})
		if err==nil{
			if phoneResult.Result!=0{
				c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorPhoneRepeat)
				return
			}
		}
	}
	ttl, err := cache.GetUserPhoneCodeLimitTime(params.Phone, params.Type)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if ttl > 0 {
		result.Time = ttl
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorVerifiedCodeFrequent.SetData(result))
		return
	}

	code := utils.Get6Code()
	rpcSendSmsResult := &pb.SendResponse{}
	tempCode:="SMS_198585116"
	if params.Type==4{
		tempCode="SMS_198585115"
	}
	paramStr := fmt.Sprintf("{code:\"%s\"}", code)
	rpcSendSmsResult, err = manager.GetSendSmsService().SendSms(c, &pb.SendRequest{
		Key:"",
		KeyId:"",
		Phone:        params.Phone,
		ParamStr:     paramStr,
		SingName:     "",
		TemplateCode: tempCode,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}

	if rpcSendSmsResult.Result != "OK" {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseMsg(rpcSendSmsResult.Result))
		return
	}

	result.Time, err = cache.SetUserPhoneCode(params.Phone, code, params.Type)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
}
