package good

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"mufe_service/camp/dataUtil"
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

const defaultStatus = enum.StatusNormal

type DeliveryData struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	DefaultNum    int64  `json:"default_num"`
	DefaultPrice  int64  `json:"default_price"`
	IncreaseNum   int64  `json:"increase_num"`
	IncreasePrice int64  `json:"increase_price"`
}

type Sku struct {
	Id             int64       `json:"id"`
	Price          int64       `json:"price"`
	AgreementPrice int64       `json:"agreement_price"`
	MemberPrice    int64       `json:"member_price"`
	Stock          int64       `json:"stock"`
	Del            bool        `json:"del"`
	Options        []SkuOption `form:"options" json:"options"`
}

type SkuOption struct {
	OptionId      int64  `json:"id"`
	OptionValue   string `json:"value"`
	OptionValueId int64  `json:"o_id"`
	Uuid          string `json:"uuid"`
}

type Photo struct {
	Key  string `json:"key"`
	Type int64  `json:"type"`
}

func init() {
	server.Get("/adminGood/getToken", token)
	server.Get("/adminGood/good", handler.BrandAdminLogin, getGoodList)   //分类列表
	server.Get("/adminGood/goodCategory", handler.BrandAdminLogin, getGoodCategory)
	server.Get("/adminGood/commitment", handler.BrandAdminLogin, GetCommitment) //获取承诺服务列表
	server.Get("/adminGood/deliveryTemplate", handler.BrandAdminLogin, getDeliverTemplate)
	server.Post("/adminGood/deliveryTemplate", handler.BrandAdminLogin, putDeliveryTemplate)
	server.Post("/adminGood/editGood", handler.BrandAdminLogin, putAddGood)
}

func getGoodCategory(c *gin.Context) {
	adminData, ok := c.MustGet(handler.BrandAdminData).(*jwt.AdminClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	categoryResult, err := servicemanager.GetAdminUserService().GetBrandGoodCategory(c, &pb.GetBrandGoodCategoryRequest{BusinessId: adminData.BusinessId})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	} else {
		resultData, err := servicemanager.GetGoodCategoryService().GoodCategory(c, &pb.GetGoodCategoryRequest{ParentCategoryIdList: categoryResult.List, Type: enum.GetCategoryInfoFromId})
		if err == nil {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(dataUtil.ParseCategoryInfoList(resultData.List)))
		} else {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		}
	}
}

//获取承诺服务列表
func GetCommitment(c *gin.Context) {
	type Result struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}
	result, err := servicemanager.GetAdminUserService().GetCommitment(c, &pb.EmptyRequest{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	list := make([]Result, 0)
	for _, v := range result.List {
		list = append(list, Result{
			Id:   v.Id,
			Name: v.Name,
		})
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
}

func getDeliverTemplate(c *gin.Context) {

	adminData, ok := c.MustGet(handler.BrandAdminData).(*jwt.AdminClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	result, err := servicemanager.GetGoodDeliveryService().GetDeliveryTemplate(c, &pb.DeliveryListRequest{BusinessId: adminData.BusinessId})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	} else {
		list := make([]DeliveryData, 0)
		for _, info := range result.List {
			result := DeliveryData{}
			result.Id = info.Id
			result.Name = info.Name
			result.DefaultNum = info.DefaultNum
			result.DefaultPrice = info.DefaultPrice
			result.IncreaseNum = info.IncreaseNum
			result.IncreasePrice = info.IncreasePrice
			list = append(list, result)
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	}
}

func putDeliveryTemplate(c *gin.Context) {
	type Query struct {
		Id            int64  `form:"id" json:"id" `
		Del           bool   `form:"del" json:"del" `
		Name          string `form:"name" json:"name" `
		DefaultNum    int64  `form:"default_num" json:"default_num"`
		DefaultPrice  int64  `form:"default_price" json:"default_price"`
		IncreaseNum   int64  `form:"increase_num" json:"increase_num"`
		IncreasePrice int64  `form:"increase_price" json:"increase_price"`
	}
	param := Query{}
	err := c.ShouldBind(&param)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	adminData, ok := c.MustGet(handler.BrandAdminData).(*jwt.AdminClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	if param.Del {
		_, err := servicemanager.GetGoodDeliveryService().DelDeliveryTemplate(c, &pb.DelDeliveryRequest{
			Id:         param.Id,
			BusinessId: adminData.BusinessId,
			Uid:        adminData.Uid,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		} else {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("删除成功"))
		}
	} else {
		data := &pb.DeliveryData{
			Id:            param.Id,
			Name:          param.Name,
			DefaultNum:    param.DefaultNum,
			DefaultPrice:  param.DefaultPrice,
			IncreasePrice: param.IncreasePrice,
			IncreaseNum:   param.IncreaseNum,
		}
		result, err := servicemanager.GetGoodDeliveryService().EditDeliveryTemplate(c,
			&pb.EditDeliveryListRequest{
				BusinessId: adminData.BusinessId,
				Uid:        adminData.Uid,
				Data:       data,
			})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
			return
		} else {
			list := make([]DeliveryData, 0)
			for _, info := range result.List {
				result := DeliveryData{}
				result.Id = info.Id
				result.Name = info.Name
				result.DefaultNum = info.DefaultNum
				result.DefaultPrice = info.DefaultPrice
				result.IncreaseNum = info.IncreaseNum
				result.IncreasePrice = info.IncreasePrice
				list = append(list, result)
			}
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
		}
	}

}

//发布商品
func putAddGood(c *gin.Context) {
	type Params struct {
		Name         string                   `form:"name" json:"name"`               //商品名称
		Detail       string                   `form:"detail" json:"detail"`           //详情
		Location     string                   `form:"location" json:"location"`       //发货地
		CategoryId   int64                    `form:"category_id" json:"category_id"` //分类id
		Draft        bool                     `form:"draft" json:"draft"`             //草稿
		Infos        []Sku                    `form:"sku" json:"sku"`                 //sku
		DeliveryInfo []*dataUtil.DeliveryInfo `form:"delivery" json:"delivery"`       //配送
		Type         int64                    `form:"type" json:"type"`               //类型
	}
	params := Params{}
	err := c.ShouldBind(&params)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	list := make([]*pb.Sku, 0)
	for _, v := range params.Infos {
		options := make([]*pb.SkuOption, 0)
		for _, k := range v.Options {
			if k.OptionValue != "" {
				options = append(options, &pb.SkuOption{OptionValue: k.OptionValue, OptionId: k.OptionId, Uuid: k.Uuid})
			}
		}
		list = append(list, &pb.Sku{
			Options:        options,
			Stock:          v.Stock,
			Price:          v.Price,
			AgreementPrice: v.AgreementPrice,
			CostPrice:      v.AgreementPrice,
		})
	}
	status := defaultStatus
	if params.Draft {
		status = enum.StatusDraft
	}
	adminData, ok := c.MustGet(handler.BrandAdminData).(*jwt.AdminClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}

	deliverInfoList := make([]*pb.DeliveryData, 0)
	for _, info := range params.DeliveryInfo {
		deliverInfoList = append(deliverInfoList, &pb.DeliveryData{
			Type: info.Type,
			Id:   info.TemplateId,
		})
	}
	result, err := servicemanager.GetGoodService().EditGood(c, &pb.EditGoodRequest{
		Infos:      list,
		SpuName:    params.Name,
		Uid:        adminData.Uid,
		CategoryId: params.CategoryId,
		Location:   params.Location,
		Detail:     params.Detail,
		IsDraft:    params.Draft,
		Delivery:   deliverInfoList,
		BusinessId: adminData.BusinessId,
		Type:       params.Type,
		Status:     status,
	})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result.SpuId))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

type QiniuInfo struct {
	Token    string   `json:"token"`
	Host     string   `json:"host"`
	BaseHost string   `json:"base_host"`
	Keys     []string `json:"keys"`
	Prefix   string   `json:"prefix"`
}

func token(c *gin.Context) {
	var param struct {
		Names string `form:"names"`
		Temp  string `form:"is_base64"`
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	var isBase64 bool
	if param.Temp == "true" {
		isBase64 = true
	} else {
		isBase64 = false
	}
	var list []string
	err := json.Unmarshal([]byte(param.Names), &list)

	for i := 0; i < len(list); i++ {
		fullFilename := list[i]
		filenameWithSuffix := path.Base(fullFilename)
		fileSuffix := path.Ext(filenameWithSuffix) //获取文件后缀
		list[i] = utils.MD5(fullFilename+strconv.FormatInt(time.Now().Unix(), 10)) + fileSuffix
	}
	osStr := os.Getenv("IMG_BUCKET")
	prefix := os.Getenv("IMG_PREFIX")
	result, err := servicemanager.GetQiniuService().GetToken(c, &pb.QiniuServiceRequest{Bucket: osStr, IsBase64: isBase64})
	if err == nil {
		if isBase64 {
			for i := 0; i < len(list); i++ {
				encodeString := base64.StdEncoding.EncodeToString([]byte(list[i]))
				list[i] = encodeString
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(QiniuInfo{Token: result.Token, Keys: list, Host: result.UploadHost, Prefix: prefix}))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}


func getGoodList(c *gin.Context) {
	//page := utils.GetInt64ValueFromReq(c, "page")
	//pageSize := utils.GetInt64ValueFromReq(c, "size")
	//status := utils.GetInt64ValueFromReq(c, "status")
	//saleStatus := enum.StatusAllSale
	//if status == enum.StatusNormal {
	//	saleStatus = enum.StatusSale
	//} else if status == enum.StatusExceptional {
	//	saleStatus = enum.StatusUnSale
	//	status = enum.StatusNormal
	//}
	//isFew := utils.GetInt64ValueFromReq(c, "is_few")
	//name := c.Query("name")
	//keys := make([]string, 0)
	//if name != "" {
	//	keys = append(keys, name)
	//}
	//adminData, ok := c.MustGet(handler.BrandAdminData).(*jwt.AdminClaims)
	//if !ok {
	//	c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
	//	return
	//}
	//resultData, err := servicemanager.GetGoodService().GetSpuList(c, &pb.GetSpuListRequest{
	//	BusinessIdList: []int64{adminData.BusinessId},
	//	Page:           page,
	//	Size:           pageSize,
	//	Status:         status,
	//	SaleStatus:     saleStatus,
	//	Keywords:       keys,
	//	Channel:        enum.ChannelALl,
	//	Few:            isFew == 1,
	//})
	//if err == nil {
	//	list := make([]interface{}, 0)
	//	idName := make(map[int64]string)
	//	spuIdMap := make(map[int64]int64)
	//	shopCategoryResult, err := servicemanager.GetShopCategoryService().GetGoodRefCategory(c, &pb.GetGoodRefCategoryRequest{BusinessId: []int64{adminData.BusinessId}})
	//	if err == nil {
	//		for _, info := range shopCategoryResult.List {
	//			idName[info.CategoryId] = info.CategoryName
	//			for _, id := range info.SpuId {
	//				spuIdMap[id] = info.CategoryId
	//			}
	//		}
	//	}
	//	goodsResult, err := spuInfo.GetSpu(resultData.List, true)
	//	if err == nil {
	//		deliverInfoResult, err := servicemanager.GetGoodSpuService().GetSpuDeliverInfo(c, &pb.GetSpuDeliverInfoRequest{List: resultData.List})
	//		temps := dataUtil.ParseBrandSpu(goodsResult)
	//		for _, info := range temps {
	//			categoryId, ok := spuIdMap[info.Id]
	//			if ok {
	//				name, ok := idName[categoryId]
	//				if ok {
	//					info.Category = append(info.Category, dataUtil.BrandCategory{
	//						Id:   categoryId,
	//						Name: name,
	//						Type: 1,
	//					})
	//				}
	//			}
	//
	//			if err == nil {
	//				if v, ok := deliverInfoResult.Result[info.Id]; ok {
	//					info.Wholes = v.Wholes
	//					info.Shop = v.Shop
	//					info.Local = v.Tm
	//					info.DeliverInfoName = v.Name
	//				}
	//			}
	//			list = append(list, info)
	//		}
	//	}
	//
	//	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(utils.CreateListResultReturn(resultData.Total, list)))
	//} else {
	//	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	//}
}
