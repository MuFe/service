package shop

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"net/http"
	"os"
	"mufe_service/camp/dataUtil"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/handler"
	"mufe_service/camp/jwt"
	"mufe_service/camp/server"
	"mufe_service/camp/wx/util"
	pb "mufe_service/jsonRpc"
	"mufe_service/manager"
	"time"
)

func init() {
	server.Post("/appApi/shop/detail", detail)
	server.Post("/appApi/order/orderList", handler.UserLogin, orderList)
	server.Post("/appApi/order/address", handler.UserLogin, address)
	server.Post("/appApi/order/editAddress", handler.UserLogin, editAddress)
	server.Post("/appApi/order/createOrder", handler.UserLogin, createOrder)
	server.Post("/appApi/order/pay", handler.UserLogin, pay)
	server.Post("/appApi/pay/notify", notify)
	server.Post("/appApi/order/su", handler.UserLogin, su)
	server.Post("/appApi/order/cancelOrder", handler.UserLogin, cancelOrder)
	server.Post("/appApi/order/receiveGood", handler.UserLogin, receiveGood)

}

type OrderData struct {
	Data       []*pb.OrderData
	Title      string
	Desc       string
	OrderPrice int64
}

type Address struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Area      string `json:"area"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Phone     string `json:"phone"`
	IsDefault bool   `json:"is_default"`
}

func detail(c *gin.Context) {
	type query struct {
		Id int64 `form:"id"`
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	resultData, err := manager.GetGoodService().GoodDetail(c, &pb.GoodDetailRequest{SkuId: params.Id, BusinessGroupId: enum.BRAND_ADMIN_GROUP, SkuStatus: enum.StatusNormal})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	info := dataUtil.ParseShopGoodDetail(resultData)
	infos := make([]*pb.DeliveryData, 0)
	if len(resultData.DeliveryInfo) > 0 {
		for _, tInfo := range resultData.DeliveryInfo {
			if tInfo.Type == enum.DeliveryTypeShop && tInfo.Id > 0 {
				tempResult, err := manager.GetGoodDeliveryService().GetDeliveryTemplate(c, &pb.DeliveryListRequest{BusinessId: resultData.BusinessId, Id: []int64{tInfo.Id}})
				if err == nil {
					infos = tempResult.List
				}
			}
		}
	}
	info.Info.DeliveryInfo = dataUtil.ParseSkuDelivery(infos)
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(info))
}

func address(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	resultData, err := manager.GetOrderService().GetUserAddress(c, &pb.AddressListServiceRequest{Uid: userData.Uid, IsFirst: false})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	list := make([]Address, 0)
	for _, info := range resultData.List {
		data := Address{}
		data.Address = info.Address
		data.Id = info.Id
		data.Area = info.Area
		data.Province = info.Province
		data.City = info.City
		data.Phone = info.Phone
		data.Name = info.Name
		data.IsDefault = info.IsDefault
		list = append(list, data)
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
}

func editAddress(c *gin.Context) {
	type query struct {
		Id        int64  `form:"id" json:"id"`
		Name      string `form:"name" json:"name"`
		Province  string `form:"province" json:"province"`
		City      string `form:"city" json:"city"`
		Area      string `form:"area" json:"area"`
		Address   string `form:"address" json:"address"`
		Phone     string `form:"phone" json:"phone"`
		IsDefault bool   `form:"is_default" json:"is_default"`
		IsDel     bool   `form:"del" json:"del"`
	}
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetOrderService().EditUserAddress(c, &pb.EditAddressServiceRequest{Uid: userData.Uid, Address: &pb.Address{
		Id:        params.Id,
		Address:   params.Address,
		Name:      params.Name,
		Phone:     params.Phone,
		Province:  params.Province,
		City:      params.City,
		Area:      params.Area,
		IsDefault: params.IsDefault,
	}, IsDel: params.IsDel})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	} else {
		if params.IsDel {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("删除成功"))
		} else {
			if params.Id != 0 {
				c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("修改成功"))
			} else {
				c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("添加成功"))
			}
		}
	}
}

func createOrder(c *gin.Context) {
	type OrderGoodRequestData struct {
		SkuId int64 `form:"sku_id" json:"sku_id"`
		Num   int64 `form:"num" json:"num"`
	}
	type Param struct {
		Goods     []OrderGoodRequestData `form:"goods" json:"goods"`
		Phone     string                 `form:"id" json:"phone"`
		Province  string                 `form:"province" json:"province"`
		City      string                 `form:"city" json:"city"`
		Area      string                 `form:"area" json:"area"`
		Address   string                 `form:"address" json:"address"`
		Consignee string                 `form:"consignee" json:"consignee"`
	}
	params := Param{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	goods := make(map[int64]int64)
	for _, info := range params.Goods {
		goods[info.SkuId] = info.Num
	}

	pOrderResult, err := CreateOrder(goods, enum.OrderTypeApp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if len(pOrderResult.Data) == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	orderResult, err := manager.GetOrderService().Create(c, &pb.OrderRequest{
		BuyerId:   userData.Uid,
		Response:  pOrderResult.Data,
		Time:      int64(30 * 60),
		Title:     pOrderResult.Title,
		Desc:      pOrderResult.Desc,
		Price:     pOrderResult.OrderPrice,
		OrderType: enum.OrderTypeApp,
		Phone:     params.Phone,
		Province:  params.Province,
		City:      params.City,
		Area:      params.Area,
		Address:   params.Address,
		Consignee: params.Consignee,
		ShowOrderList:enum.ShowOrderList,
	})

	if err == nil {
		list := make([]*pb.EditShopCarData, 0)
		for _, info := range pOrderResult.Data {
			for _, v := range info.List {
				list = append(list, &pb.EditShopCarData{
					SkuId: v.SkuId,
					Num:   0,
				})
			}
		}
		_, err = manager.GetOrderService().EditShopCar(c, &pb.EditShopCarRequest{Uid: userData.Uid, List: list, Clear: true})
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(orderResult.Id))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func CreateOrder(dataMap map[int64]int64, typeInt int64) (*OrderData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, err := manager.GetGoodService().GetOrderSkuList(ctx, &pb.GetOrderSkuListRequest{Data: dataMap})
	if err == nil {
		list := make([]*pb.OrderData, 0)
		allPrice := int64(0)
		for _, info := range result.List {
			dataList := make([]*pb.OrderGood, 0)
			listPrice := int64(0)
			deliverPrice := int64(0)
			for _, v := range info.List {
				price := v.Price
				allPrice += price * v.Num
				listPrice += price * v.Num
				dataList = append(dataList, &pb.OrderGood{
					SkuId:          v.SkuId,
					Num:            v.Num,
					Price:          price,
					SpuId:          v.SpuId,
					AgreementPrice: v.AgreementPrice,
				})
			}
			if typeInt == enum.OrderTypeApp {
				for _, v := range info.ShopDeliverPriceMap {
					if deliverPrice < v {
						deliverPrice = v
					}
				}
			}
			listPrice += deliverPrice
			allPrice += deliverPrice
			list = append(list, &pb.OrderData{
				SellerId:     info.SellerId,
				List:         dataList,
				Price:        listPrice,
				DeliverPrice: deliverPrice,
			})
		}
		return &OrderData{
			Title:      "爱球知",
			Desc:       "购买商品",
			Data:       list,
			OrderPrice: allPrice,
		}, nil
	} else {
		return nil, err
	}
}

func pay(c *gin.Context) {
	type Param struct {
		Id      int64 `form:"id" json:"id"`
		Channel int64 `form:"channel" json:"channel"`
	}
	params := Param{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}

	result, err := manager.GetOrderService().Pay(c, &pb.OrderPayRequest{
		OrderId:     params.Id,
		AppId:       os.Getenv("APP_WX_ID"),
		NotifyUrl:   os.Getenv("API") + "/pay/notify",
		ChannelType: params.Channel,
		PayType:     enum.AppPayType,
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if result.Pay {
		c.AbortWithStatusJSON(http.StatusNotModified, errcode.ParseOK("已经支付过"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(dataUtil.ParseWxPay(result)))
	}
}

func su(c *gin.Context) {
	type Param struct {
		Id      int64 `form:"id" json:"id"`
	}
	params := Param{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetOrderService().EditStatus(c,&pb.EditStatusRequest{
		OrderId: params.Id,
		Status:  enum.OrderStatusPay,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
}

func cancelOrder(c *gin.Context) {
	type Param struct {
		Id      int64 `form:"id" json:"id"`
	}
	params := Param{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	_, err := manager.GetOrderService().EditStatus(c,&pb.EditStatusRequest{
		OrderId: params.Id,
		BuyerId:userData.Uid,
		Status:  enum.OrderStatusUserCancel,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("取消订单成功"))
	}
}

func receiveGood(c *gin.Context) {
	type Param struct {
		Id      int64 `form:"id" json:"id"`
	}
	params := Param{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	_, err := manager.GetOrderService().EditStatus(c,&pb.EditStatusRequest{
		SubOrderId: params.Id,
		BuyerId:userData.Uid,
		Status:  enum.OrderStatusFinish,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("收货成功"))
	}
}


func orderList(c *gin.Context){
	type Param struct {
		Id      int64 `form:"id" json:"id"`
		Type      int64 `form:"type" json:"type"`
		Page      int64 `form:"page" json:"page"`
		Size      int64 `form:"size" json:"size"`
	}
	params := Param{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	orderResult, err := manager.GetOrderService().GetOrders(c,
		&pb.GetOrderRequest{
			BuyerId:    userData.Uid,
			Status:     enum.GetOrderStatusFromQuery(params.Type),
			OrderTypes: []int64{enum.OrderTypeApp},
			OrderId:   params.Id,
			Page:       params.Page,
			Size:      params.Size,
			ShowOrder:enum.ShowOrderList,
		})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	list := make([]dataUtil.UserOrderData, 0, len(orderResult.List))
	for _, info := range orderResult.List {
		list = append(list, dataUtil.ParseOrder(info))
	}
	type OrderResult struct {
		List  []dataUtil.UserOrderData `json:"list"`
		Total int64                    `json:"total"`
	}

	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(OrderResult{
		List:  list,
		Total: orderResult.Total,
	}))
}

func notify(c *gin.Context){
	notify:=util.GetNotify(c)
	content :=new(payments.Transaction)
	notifyReq, err := notify.ParseNotifyRequest(c, c.Request, content)
	// 如果验签未通过，或者解密失败
	if err != nil {
		fmt.Println(err)
		return
	}
	if notifyReq.EventType=="TRANSACTION.SUCCESS"{
		_, err := manager.GetOrderService().EditStatus(c,&pb.EditStatusRequest{
			OrderSn: *content.OutTradeNo,
			Status:  enum.OrderStatusPay,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
			return
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
	}

}
