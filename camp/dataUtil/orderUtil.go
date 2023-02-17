package dataUtil

import (
	"mufe_service/camp/enum"
	pb "mufe_service/jsonRpc"
)


type UserOrderData struct {
	ShopPhoto      string     `json:"shopPhoto"`
	ShopName       string     `json:"shop_name"`
	ShopPhone      string     `json:"shopPhone"`
	Id             int64      `json:"id"`
	SubOrderId     int64      `json:"sub_order_id"`
	OrderSn        string     `json:"order_sn"`
	PayChannel     string     `json:"pay_channel"`
	UserName       string     `json:"user_name"`
	Message        string     `json:"message"`
	ExpressNumber  string     `json:"express_number"`
	ExpressCompany string     `json:"express_company"`
	Phone          string     `json:"phone"`
	Address        string     `json:"address"`
	Time           int64      `json:"time"`
	EndTime        int64      `json:"end_time"`
	Status         int64      `json:"status"`
	Price          int64      `json:"price"`
	CouponPrice    int64      `json:"coupon_price"`
	DeliveryPrice  int64      `json:"delivery_price"`
	GoodPrice      int64      `json:"good_price"`
	SellerID       int64      `json:"sellerID"`
	Goods          []GoodInfo `json:"goods"`
}

type GoodInfo struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Photo    string `json:"photo"`
	Price    int64  `json:"price"`
	Num      int64  `json:"num"`
	Spec     string `json:"spec"`
	NowPrice int64  `json:"nowPrice"`
}

type WxPayServiceResponse struct {
	OrderId   int64  `json:"order_id"`
	OrderSn   string `json:"order_sn"`
	PrepayId   string `json:"prepay_id"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	PaySign   string `json:"paySign"`
	PartnerId   string `json:"partnerId"`
	SignType  string `json:"signType"`
	AppId     string `json:"appId"`
	CodeUrl     string `json:"url"`
}

func ParseWxPay(orderResult *pb.PayServiceResponse) WxPayServiceResponse {
	return WxPayServiceResponse{
		OrderSn:   orderResult.OrderSn,
		OrderId:   orderResult.OrderId,
		TimeStamp: orderResult.TimeStamp,
		PrepayId:orderResult.PrepayId,
		NonceStr:  orderResult.NonceStr,
		Package:   orderResult.Package,
		PaySign:   orderResult.PaySign,
		SignType:  orderResult.SignType,
		AppId:     orderResult.AppId,
		CodeUrl:orderResult.CodeUrl,
		PartnerId:orderResult.PartnerId,
	}
}


func ParseOrder(info *pb.GetOrderData) UserOrderData {
	result := UserOrderData{}
	result.Time = info.OrderTime
	result.Price = info.PayMount
	result.Id = info.OrderId
	result.SubOrderId = info.SubOrderId
	result.EndTime = info.EndTime
	result.UserName = info.Consignee
	result.Phone = info.Phone
	result.Address = info.Address
	result.OrderSn = info.OrderSn
	result.PayChannel = enum.GetPayChannel(info.PayChannel)
	result.Message = info.Message
	result.CouponPrice = info.CouponPrice
	result.DeliveryPrice = info.DeliveryPrice
	result.ExpressNumber = info.ExpressNumber
	result.ExpressCompany = info.ExpressCompany
	result.SellerID = info.SellerId
	result.ShopName = info.ShopName
	result.ShopPhoto = info.ShopPhoto
	result.ShopPhone = info.ShopPhone
	var goodPrice int64
	list := make([]GoodInfo, 0, len(info.List))
	for _, tempInfo := range info.List {
		for _, t := range tempInfo.List {
			temp := parseOrderGood(t)
			temp.Photo = tempInfo.Photo
			temp.Name = tempInfo.Name
			list = append(list, temp)
			goodPrice += temp.Price * temp.Num
		}
	}
	result.Status = info.Status
	result.GoodPrice = goodPrice - info.CouponPrice
	result.Goods = list
	return result
}

func parseOrderGood(t *pb.GetOrderSku) GoodInfo {
	t1 := GoodInfo{}
	t1.Id = t.SkuId
	t1.Price = t.Price
	t1.Num = t.Num
	t1.Spec = t.Spec
	t1.NowPrice = t.NowPrice
	return t1
}
