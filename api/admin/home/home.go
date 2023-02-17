package login

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/handler"
	"mufe_service/camp/server"
	"mufe_service/camp/utils"
	"mufe_service/jsonRpc"
	"mufe_service/manager"
	"strconv"
	"time"
)

func init() {
	server.Post("/adminHome/home", handler.AdminLogin, home)
	server.Post("/adminBanner/list", handler.AdminLogin, bannerList)
	server.Post("/adminBanner/edit", bannerEdit)
	server.Post("/adminBanner/editSort", bannerSortEdit)
	server.Post("/adminBanner/editPhoto", bannerPhotoEdit)
	server.Post("/adminBanner/token", getToken)
	server.Post("/adminWeb/icp", icp)
	server.Post("/adminWeb/contact", contact)
	server.Post("/adminWeb/app", appInfo)
	server.Post("/adminWeb/editContact", editContact)

}

type Result struct {
	OrderNo   string `json:"order_no"`
	TimeStamp int64  `json:"timestamp"`
	UserName  string `json:"username"`
	Price     string `json:"price"`
	Status    string `json:"status"`
}

func home(c *gin.Context) {
	list := make([]interface{}, 0)
	data := Result{
		OrderNo:   "1111",
		TimeStamp: time.Now().Unix(),
		UserName:  "test",
		Price:     "100",
		Status:    "success",
	}
	list = append(list, data)
	list = append(list, data)
	list = append(list, data)
	list = append(list, data)
	list = append(list, data)
	list = append(list, data)
	list = append(list, data)
	list = append(list, data)
	list = append(list, data)
	list = append(list, data)
	list = append(list, data)
	list = append(list, data)
	list = append(list, data)

	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(utils.CreateListResultReturn(20, list)))
}

type Ad struct {
	Url    string `json:"url"`
	Photo  string `json:"photo"`
	Type   int64  `json:"type"`
	Id     int64  `json:"id"`
	LinkId int64  `json:"link_id"`
	Sort int64  `json:"sort"`
}

func bannerList(c *gin.Context) {
	type query struct {
		Id     int64  `form:"id" json:"id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	resultData, err := manager.GetBannerService().GetAds(c, &app.AdServiceRequest{Id:params.Id,Status: enum.StatusAll})
	var list []Ad
	if err == nil {
		list = make([]Ad, 0)
		for _, v := range resultData.Result {
			list = append(list, parseAd(v))
		}
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
}

func bannerEdit(c *gin.Context) {
	type query struct {
		Id     int64  `form:"id" json:"id" `
		Type   int64  `form:"type" json:"type" `
		LinkId int64  `form:"link_id" json:"link_id" `
		Url    string `form:"url" json:"url" `
		IsDel  bool   `form:"del" json:"del" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.IsDel {
		_, err := manager.GetBannerService().DelAd(c, &app.EditAdRequest{Id: params.Id})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
			return
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("操作成功"))
	} else {
		idResult, err := manager.GetBannerService().EditAd(c, &app.EditAdRequest{Id: params.Id, Type: params.Type, LinkId: params.LinkId, Url: params.Url})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
			return
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(idResult.Id))
	}

}

func bannerSortEdit(c *gin.Context) {
	type query struct {
		Id   int64 `form:"id" json:"id" `
		Sort int64 `form:"sort" json:"sort" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetBannerService().EditAdSort(c, &app.EditAdRequest{Id: params.Id, Sort: params.Sort})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("操作成功"))
}

func bannerPhotoEdit(c *gin.Context) {
	type query struct {
		Id     int64  `form:"id" json:"id" `
		Photo  string `form:"photo" json:"photo" `
		Prefix string `form:"prefix" json:"prefix" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetBannerService().EditAdPhoto(c, &app.EditAdRequest{Id: params.Id, Key: params.Photo, Prefix: params.Prefix})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("操作成功"))
}

func getToken(c *gin.Context) {
	type query struct {
		Name []string `form:"names" json:"names" `
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	list := make([]string, 0)
	for _, v := range params.Name {
		filenameWithSuffix := path.Base(v)
		fileSuffix := path.Ext(filenameWithSuffix) //获取文件后缀
		encodeString := utils.MD5(v+strconv.FormatInt(time.Now().Unix(), 10)) + fileSuffix
		list = append(list, encodeString)
	}
	osStr := os.Getenv("IMG_BUCKET")
	prefix := os.Getenv("IMG_PREFIX")
	type QiniuInfo struct {
		Token    string   `json:"token"`
		Host     string   `json:"host"`
		BaseHost string   `json:"base_host"`
		Keys     []string `json:"keys"`
		Prefix   string   `json:"prefix"`
	}
	result, err := manager.GetQiniuService().GetToken(c, &app.QiniuServiceRequest{Bucket: osStr})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(&QiniuInfo{Token: result.Token, Host: result.UploadHost, BaseHost: result.Base64UploadHost, Keys: list, Prefix: prefix}))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func parseAd(info *app.AdServiceResponse) Ad {
	ad := Ad{}
	ad.Url = info.Url
	ad.Type = info.Type
	ad.LinkId = info.LinkId
	ad.Id = info.Id
	ad.Photo = info.Photo
	ad.Sort=info.Sort
	return ad
}

func icp(c *gin.Context) {
	result, err := manager.GetWebService().GetWebInfo(c, &app.GetWebInfoRequest{List: []int64{enum.WebIcp}})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result.Content[enum.WebIcp]))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func contact(c *gin.Context) {
	type query struct {
		Content string `form:"content" json:"content" `
		Name    string `form:"name" json:"name" `
		Email   string `form:"email" json:"email" `
		Phone   string `form:"phone" json:"phone" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetWebService().ContactUs(c, &app.ContactUsRequest{
		Name:    params.Name,
		Content: params.Content,
		Email:   params.Email,
		Phone:   params.Phone,
	})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("提交成功"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func editContact(c *gin.Context) {
	type query struct {
		Id int64 `form:"id" json:"id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetWebService().ContactUs(c, &app.ContactUsRequest{
		Id:     params.Id,
		Status: enum.StatusNormal,
	})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("提交成功"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func appInfo(c *gin.Context) {
	type Result struct {
		Android string `json:"android"`
		Ios     string `json:"ios"`
	}
	returnResult:=Result{}
	result, err := manager.GetWebService().GetWebInfo(c, &app.GetWebInfoRequest{List: []int64{enum.WebAndroid,enum.WebIos}})
	if err == nil {
		returnResult.Android=result.Content[enum.WebAndroid]
		returnResult.Ios=result.Content[enum.WebIos]
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(returnResult))
}
