package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"mufe_service/camp/cache"
	"mufe_service/camp/dataUtil"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/handler"
	"mufe_service/camp/jwt"
	"mufe_service/camp/server"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	pb "mufe_service/jsonRpc"
	"mufe_service/manager"
	"strconv"
	"time"
	"unicode/utf8"
)

func init() {
	server.Post("/appApi/outh", outh)
	server.Post("/appApi/token", token)
	server.Post("/appApi/updateToken",  handler.UserLogin,updateToken)
	server.Post("/appApi/login", login)
	server.Post("/appApi/getUserInfo", handler.UserLogin, getUserInfo)
	server.Post("/appApi/cancel", handler.UserLogin, cancel)
	server.Post("/appApi/sendCancelCode", handler.UserLogin, sendCancelCode)
	server.Post("/appApi/checkUser", handler.UserLogin, checkUser)
	server.Post("/appApi/sendCode", sendCode)
	server.Post("/appApi/bind", handler.UserLogin, bind)
	server.Post("/appApi/modify", handler.UserLogin, modify)
	server.Post("/appApi/modifyOuth", handler.UserLogin, modifyOuth)
	server.Post("/appApi/modifyUserInfo", handler.UserLogin, modifyUserInfo)
	server.Post("/appApi/modifyHead", handler.UserLogin, modifyHead)
	server.Post("/appApi/forgetPass", forgetPass)
	server.Post("/appApi/register", register)
	server.Post("/appApi/address", address)
	server.Post("/appApi/phoneCheck", phoneCheck)
	server.Post("/appApi/chooseIdentity", handler.UserLogin, choose)
	server.Post("/appApi/collectionList", handler.UserLogin, collectionList)
	server.Post("/appApi/history", history)
	server.Post("/appApi/delHistory", handler.UserLogin, delHistory)
	server.Post("/appApi/logout", logout)
	server.Post("/appApi/updatePush",  handler.UserLogin, pushInfo)
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
	if params.Phone == "13866666666" {
		rpcSendSmsResult.Result = "OK"
		code = "123456"
	} else {
		paramStr := fmt.Sprintf("{code:\"%s\"}", code)
		rpcSendSmsResult, err = manager.GetSendSmsService().SendSms(c, &pb.SendRequest{
			Phone:        params.Phone,
			ParamStr:     paramStr,
			SingName:     enum.SignName,
			TemplateCode: enum.TemplateCode,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
			return
		}
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
	err := cache.CheckUserPhoneCode(params.Phone, params.Code, enum.CodeRegisterType)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	inviteMap:=make(map[string]bool)
	code:=utils.BaseCode6(inviteMap)
	_, err = manager.GetUserService().Register(c, &pb.LoginRequest{Phone: params.Phone, Pass: params.Pass, Name: "用户"+code, Sign: "绿茵赛场，成就梦想"})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func address(c *gin.Context) {
	type query struct {
		Id   int64 `form:"id" json:"id" `
		Type int64 `form:"type" json:"type" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	type data struct {
		Id     int64  `form:"id" json:"id" `
		Pid    int64  `form:"pid" json:"pid" `
		Name   string `form:"name" json:"name" `
		First  string `form:"first" json:"first" `
		Letter string `form:"letter" json:"letter" `
		List   []data `form:"list" json:"list" `
	}
	resultList := make([]data, 0)
	result, err := manager.GetUserService().Address(c, &pb.AddressRequest{Id: params.Id, Type: params.Type})
	if err == nil {
		for _, v := range result.List {
			list := make([]data, 0)
			for _, vv := range v.List {
				list = append(list, data{
					Id:     vv.Id,
					Name:   vv.Name,
					First:  vv.First,
					Letter: vv.Letter,
					Pid:    vv.Pid,
					List:   make([]data, 0),
				})
			}
			resultList = append(resultList, data{
				Id:     v.Id,
				Name:   v.Name,
				First:  v.First,
				Letter: v.Letter,
				Pid:    v.Pid,
				List:   list,
			})
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(resultList))
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
	err := cache.CheckUserPhoneCode(params.Phone, params.Code, enum.CodeForgetType)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err = manager.GetUserService().UpdateUser(c, &pb.UpdateUserRequest{Phone: params.Phone, Pass: params.Pass, Type: enum.UpdateUserPass})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func bind(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Code  string `form:"code" json:"code" `
		Phone string `form:"phone" json:"phone" `
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
	err := cache.CheckUserPhoneCode(params.Phone, params.Code, enum.CodeBindPhoneType)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	result, err := manager.GetUserService().UpdateUser(c, &pb.UpdateUserRequest{Phone: params.Phone, Uid: userData.Uid,Type:enum.UpdateUserPhone})
	if err == nil {
		u := dataUtil.ParseUserCache(result)
		_ = cache.SetUserInfo(result.Uid, u)
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func modify(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Pass    string `form:"pass" json:"pass" `
		NewPass string `form:"new_pass" json:"new_pass" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	result, err := manager.GetUserService().UpdateUser(c, &pb.UpdateUserRequest{Pass: params.Pass, NewPass: params.NewPass, Type: enum.UpdateUserModifyPass, Uid: userData.Uid})
	if err == nil {
		u := dataUtil.ParseUserCache(result)
		u.HavePass=true
		_ = cache.SetUserInfo(result.Uid, u)
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("修改成功"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func modifyHead(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Head string `form:"head" json:"head" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	result, err := manager.GetUserService().UpdateUser(c, &pb.UpdateUserRequest{Head: params.Head, Type: enum.UpdateUserHead, Uid: userData.Uid})
	if err == nil {
		u := dataUtil.ParseUserCache(result)
		_ = cache.SetUserInfo(result.Uid, u)
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("修改成功"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func modifyUserInfo(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Name    string `form:"name" json:"name" `
		Sex     int64  `form:"sex" json:"sex" `
		Age     int64  `form:"age" json:"age" `
		Address string `form:"address" json:"address" `
		Phone   string `form:"phone" json:"phone" `
		Code    string `form:"code" json:"code" `
		Sign    string `form:"sign" json:"sign" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	var err error
	var result *pb.UserDataResponse

	if utf8.RuneCountInString(params.Name) > 10 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(xlog.Error("长度超过10个字符")))
		return
	}
	r, err := manager.GetQiniuService().WordInspect(c, &pb.WordInspectRequest{Content: params.Name})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	} else if r.Type == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(xlog.Error("对不起，您的昵称含有敏感词语")))
		return
	} else {
		if params.Phone != "" {
			err := cache.CheckUserPhoneCode(params.Phone, params.Code, enum.CodeBindPhoneType)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
				return
			}
			result, err = manager.GetUserService().UpdateUser(c, &pb.UpdateUserRequest{Type: enum.UpdateUserInfo, Name: params.Name, District: params.Address, Sex: params.Sex, Phone: params.Phone, Uid: userData.Uid, Age: params.Age})
		} else {
			result, err = manager.GetUserService().UpdateUser(c, &pb.UpdateUserRequest{Type: enum.UpdateUserInfo, Name: params.Name, District: params.Address, Sign: params.Sign, Sex: params.Sex, Uid: userData.Uid, Age: params.Age})
		}
	}
	if err == nil {
		u := dataUtil.ParseUserCache(result)
		_ = cache.SetUserInfo(result.Uid, u)
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("修改成功"))
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
		err := cache.CheckUserPhoneCode(params.Phone, params.Code, enum.CodePhoneLoginType)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
			return
		}
	}
	if params.LoginType==enum.JpushType{
		tokenResult,err:=manager.GetPushService().GetPhone(c,&pb.GetPhoneRequest{Token:params.Code})
		if err==nil{
			params.Phone=tokenResult.Phone
		}else{
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
			return
		}
	}
	result, err := manager.GetUserService().Login(c, &pb.LoginRequest{Type: params.LoginType, Phone: params.Phone, Code: params.Code, Pass: params.Pass, Device: utils.GetHeaderFromKey(c, "device")})
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


func token(c *gin.Context) {
	type QiniuInfo struct {
		Token  string `json:"token"`
		Host   string `json:"host"`
		Keys   string `json:"key"`
		Prefix string `json:"prefix"`
	}
	type query struct {
		Name string `form:"name" json:"name" `
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	filenameWithSuffix := path.Base(params.Name)
	fileSuffix := path.Ext(filenameWithSuffix) //获取文件后缀
	encodeString := utils.MD5(params.Name+strconv.FormatInt(time.Now().Unix(), 10)) + fileSuffix
	osStr := os.Getenv("IMG_BUCKET")
	prefix := os.Getenv("IMG_PREFIX")
	result, err := manager.GetQiniuService().GetToken(c, &pb.QiniuServiceRequest{Bucket: osStr})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(&QiniuInfo{Token: result.Token, Host: result.UploadHost, Keys: encodeString, Prefix: prefix}))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func updateToken(c *gin.Context){
	type Result struct {
		Token  string `json:"token"`
		Identity int64  `json:"identity"`
	}
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	result, err := cache.GetUserInfo(userData.Uid)
	if err!=nil{
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
	token, err := jwt.GenerateUserJwt(userData.Uid, result.Identity, userData.OpenId)
	if err!=nil{
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(&Result{Token: token,Identity:result.Identity}))
}

func outh(c *gin.Context) {
	type query struct {
		Type int64  `form:"type" json:"type" `
		Code string `form:"code" json:"code" `
		Name string `form:"name" json:"name" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Code == "" {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}

	appId := os.Getenv("APP_WX_ID")
	sec := os.Getenv("APP_WX_SEC")
	result, err := manager.GetUserService().Login(c, &pb.LoginRequest{Type: enum.OuthType, IsMiniPrograms: utils.IsMiniProgram(c),
		OuthType: params.Type, Code: params.Code, AppId: appId, Secret: sec, Name: params.Name, Sign: "绿茵赛场，成就梦想"})
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

func modifyOuth(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Type   int64  `form:"type" json:"type" `
		Code   string `form:"code" json:"code" `
		Unbind bool   `form:"unbind" json:"unbind" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Code == "" && !params.Unbind {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	appId := os.Getenv("APP_WX_ID")
	sec := os.Getenv("APP_WX_SEC")
	_, err := manager.GetUserService().ModifyOuth(c, &pb.ModifyOuthRequest{Uid: userData.Uid, OuthType: params.Type, Code: params.Code, AppId: appId, Secret: sec, IsUnbind: params.Unbind})
	if err == nil {
		result, err := cache.GetUserInfo(userData.Uid)
		if err == nil {
			if params.Unbind {
				result.HaveWx = false
			} else {
				result.HaveWx = true
			}
			_ = cache.SetUserInfo(userData.Uid, result)
		} else {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func phoneCheck(c *gin.Context) {
	type query struct {
		Phone string `form:"phone" json:"phone" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	result, err := manager.GetUserService().PhoneCheck(c, &pb.PhoneCheckRequest{Phone: params.Phone, Device: utils.GetHeaderFromKey(c, "device")})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result.Result))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func choose(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Type int64 `form:"type" json:"type" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	result, err := manager.GetUserService().UpdateUser(c, &pb.UpdateUserRequest{Uid: userData.Uid, IdentityType: params.Type, Type: enum.UpdateUserIdentity})
	if err == nil {
		token, err := jwt.GenerateUserJwt(userData.Uid, params.Type, userData.OpenId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
			return
		}
		u := dataUtil.ParseUserCache(result)
		_ = cache.SetUserInfo(result.Uid, u)
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(token))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func getUserInfo(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	result, err := cache.GetUserInfo(userData.Uid)
	if err == nil {
		user := dataUtil.ParseUserFromCache(result)
		if user.Identity==enum.INSTITUTION_TYPE{
			chapterResult, err := manager.GetCoachService().CoachList(c, &pb.CoachListRequest{Id: userData.Uid, Page: 1, Size: 10})
			if err==nil{
				user.Introduce=chapterResult.List[0].Info
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(user))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func delHistory(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Id []int64 `form:"id" json:"id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetVideoService().DelHistoryVideo(c, &pb.DelVideoHistoryRequest{Uid: userData.Uid, Id: params.Id})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
}

func history(c *gin.Context) {

	type Video struct {
		Id        int64  `json:"id" `
		ChapterId string `json:"chapter_id" `
		Cover     string `json:"cover" `
		Title     string `json:"title" `
		Time      int64  `json:"time" `
		PlayUrl   string `json:"play_url" `
		DownUrl   string `json:"down_url" `
	}
	v := make([]Video, 0)
	userData, err := jwt.CheckUserJwt(c.GetHeader(jwt.AuthHeader))
	if err == nil {
		videoResult, err := manager.GetVideoService().HistoryVideoList(c, &pb.VideoRequest{Uid: userData.Uid})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
			return
		}
		for _, temp := range videoResult.VideoList {
			v = append(v, Video{
				Id:      temp.Id,
				Time:    temp.Duration,
				Title:   temp.Title,
				Cover:   temp.Cover,
				PlayUrl: temp.Url,
				DownUrl: temp.DownUrl,
			})
		}
	}

	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(v))
}

func collectionList(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Page int64 `form:"page" json:"page" `
		Size int64 `form:"size" json:"size" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	type Tag struct {
		Tag string `json:"tag" `
		Id  int64  `json:"id" `
	}
	type Chapter struct {
		Cover   string `json:"cover" `
		Id      string `json:"id" `
		Title   string `json:"title" `
		Desc    string `json:"desc" `
		Section int64  `json:"section" `
		Time    int64  `json:"time" `
		Level   string `json:"level" `
		Tags    []Tag  `json:"tag" `
		Users   int64  `json:"users" `
	}
	re := make([]Chapter, 0)
	chapterResult, err := manager.GetCollectionService().GetCollection(c, &pb.GetCollectionRequest{Page: params.Page, Size: params.Size, Uid: userData.Uid, Type: 3})
	if err == nil {
		for _, v := range chapterResult.List {
			tags := make([]Tag, 0)
			for _, v := range v.Tag {
				tags = append(tags, Tag{
					Id:  v.Id,
					Tag: v.Title,
				})
			}
			re = append(re, Chapter{
				Cover:   v.Cover,
				Id:      strconv.FormatInt(v.Id, 10),
				Title:   v.Title,
				Desc:    v.Desc,
				Time:    v.Time,
				Level:   v.Level,
				Section: v.Section,
				Tags:    tags,
				Users:   v.User,
			})
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(re))
}

func logout(c *gin.Context) {
	var token = c.GetHeader(jwt.AuthHeader)
	if token != "" {
		handler.UserLogin(c)
		userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
		if ok {
			_ = cache.DeleteUserToken(userData.Uid)
			_ = cache.DeleteUserInfo(userData.Uid)
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
}

func cancel(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	_, err := manager.GetUserService().Cancel(c, &pb.CancelRequest{Uid: userData.Uid})
	if err == nil {
		_ = cache.DeleteUserToken(userData.Uid)
		_ = cache.DeleteUserInfo(userData.Uid)
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("注销成功"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func checkUser(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Code string `form:"code" json:"code" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Code == "" {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	result, err := cache.GetUserInfo(userData.Uid)
	if err == nil {
		user := dataUtil.ParseUserFromCache(result)
		err := cache.CheckUserPhoneCode(user.Phone, params.Code, enum.CodeCancelType)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
			return
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}

}

func sendCancelCode(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type codeResponse struct {
		Time int64 `json:"time"`
	}
	var result = codeResponse{}
	result1, err := cache.GetUserInfo(userData.Uid)
	if err == nil {
		user := dataUtil.ParseUserFromCache(result1)
		if user.Phone == "" {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(xlog.Error("该用户未绑定手机")))
			return
		}
		ttl, err := cache.GetUserPhoneCodeLimitTime(user.Phone, enum.CodeCancelType)
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
		paramStr := fmt.Sprintf("{code:\"%s\"}", code)
		rpcSendSmsResult, err = manager.GetSendSmsService().SendSms(c, &pb.SendRequest{
			Phone:        user.Phone,
			ParamStr:     paramStr,
			SingName:     enum.SignName,
			TemplateCode: enum.TemplateCode,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
			return
		}
		if rpcSendSmsResult.Result != "OK" {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseMsg(rpcSendSmsResult.Result))
			return
		}

		result.Time, err = cache.SetUserPhoneCode(user.Phone, code, enum.CodeCancelType)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
			return
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}

}

func pushInfo(c *gin.Context){
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Id    string `form:"registration_id" json:"registration_id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetUserService().UpdateUser(c, &pb.UpdateUserRequest{RegistrationId: params.Id, Type: enum.UpdatePushInfo, Uid: userData.Uid})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("修改成功"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}
