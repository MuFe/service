package enum

import (
	"fmt"
	"time"
)

// 订单类型
const (
	OrderTypeAll int64 = -1 //所有订单
	OrderTypeApp int64 = 1  //商城订单
)

const (
	BuyTypeAll    int64 = 0 //
	BuyTypeApp    int64 = 1 //app购买
)



const (
	OrderStatusUnPay                   int64 = 1  //待付款状态
	OrderStatusPay                     int64 = 2  //已支付状态
	OrderStatusSend                    int64 = 3  //已发货状态
	OrderStatusFinish                  int64 = 4  //交易成功-待评价状态
	OrderStatusUnRefund                int64 = 5  //已发货申请退款状态
	OrderStatusCloseRefund             int64 = 6  //交易关闭-仅退款
	OrderStatusCloseRefundReturn       int64 = 7  //交易关闭-退货退款(包含部分仅退款)
	OrderStatusRefund                  int64 = 8  //交易关闭-用户退款
	OrderStatusUserCancel              int64 = 9  //交易关闭-用户取消
	OrderStatusOverTime                int64 = 10 //交易关闭-超时
	OrderStatusDelete                  int64 = 11 //交易关闭-订单删除
	OrderStatusPayRefund               int64 = 14 //已支付用户申请退款
	OrderStatusCloseRefundDelete       int64 = 17 //交易关闭-仅退款-订单删除
	OrderStatusCloseRefundReturnDelete int64 = 18 //交易关闭-退货退款（包含部分仅退款）-订单删除
	OrderStatusRefundDelete            int64 = 19 //交易关闭-用户退款-订单删除
	OrderStatusUserCancelDelete        int64 = 20 //交易关闭-用户取消-订单删除
	OrderStatusOverTimeDelete          int64 = 21 //交易关闭-超时-订单删除
	OrderStatusFinishComment           int64 = 22 //交易成功-已评价
)

const (
	OrderActionStatusUnPay        int64 = 1 //待付款，
	OrderActionStatusPay          int64 = 2 //已付款，
	OrderActionStatusGiveUpPay    int64 = 3 //放弃付款，
	OrderActionPayStatusUnPay     int64 = 1 //未支付
	OrderActionPayStatusPay       int64 = 2 //支付成功
	OrderActionPayStatusPayFail   int64 = 3 //支付失败
	OrderActionPayStatusGiveUpPay int64 = 4 //放弃支付
)

const (
	MiniPayType  int64 = 1 //小程序支付
	JSAPIPayType int64 = 2 //JSAPI支付
	H5PayType    int64 = 3 //H5支付
	WEBPayType   int64 = 4 //网页扫码支付
	AppPayType   int64 = 5 //App支付
)

const(
	WEIXINPAY int64=1
	ALIPAY int64=2
)

const (
	QueryAllOrder         int64 = 0 //查看所有的订单
	QueryUnPay            int64 = 1 //查看未支付订单
	QueryPay              int64 = 2 //查看已支付订单
	QuerySend             int64 = 3 //查看已发货订单
	QueryFinish           int64 = 4 //查看已完成订单
	QueryRefund           int64 = 5 //查看申请退款订单
	QueryError            int64 = 6 //查看交易失败订单
	QueryFinishNotComment int64 = 7 //查看完成未评价
	QueryFinance 		  int64 = 8 //查看失败和未支付之外的订单
)

const (
	PayStatus_No    int64 = 1
	PayStatus_Su    int64 = 2
	PayStatus_Close int64 = 3
)

// CreateOrderSn 创建订单号
func CreateOrderSn(orderType int64, code string, isMini bool) (result string) {
	var start, end string
	switch orderType {
	case OrderTypeApp:
		start = "01"
		end = "X"
	}
	result = fmt.Sprintf("%s%d%s%s", start, time.Now().Unix(), code, end)
	return
}



// GetOrderStatusMeaning 获取订单状态的中文意义
// 下单未付款、待发货、待收货、交易完成、待退款、交易关闭-仅退款、交易关闭-退货退款（包含部分仅退款）、退款、交易关闭-用户取消、交易关闭-超时、订单删除
func GetOrderStatusMeaning(s int64) (result string) {
	switch s {
	case OrderStatusUnPay:
		result = "待付款"
	case OrderStatusPay:
		result = "已支付"
	case OrderStatusSend:
		result = "已发货"
	case OrderStatusFinish:
		result = "交易完成-待评价"
	case OrderStatusUnRefund:
		result = "待退款"
	case OrderStatusCloseRefund:
		result = "交易关闭-仅退款"
	case OrderStatusCloseRefundReturn:
		result = "交易关闭-退货退款（包含部分仅退款）"
	case OrderStatusRefund:
		result = "退款"
	case OrderStatusUserCancel:
		result = "交易关闭-用户取消"
	case OrderStatusOverTime:
		result = "交易关闭-超时"
	case OrderStatusDelete:
		result = "订单删除"
	case OrderStatusPayRefund:
		result = "已支付申请退款"
	case OrderStatusFinishComment:
		result = "交易完成-已评价"
	}
	return
}


func DeleteStatus(s int64) int64 {
	result := int64(0)
	switch s {
	case OrderStatusFinish:
		result = OrderStatusDelete
	case OrderStatusFinishComment:
		result = OrderStatusDelete
	case OrderStatusCloseRefund:
		result = OrderStatusCloseRefundDelete
	case OrderStatusCloseRefundReturn:
		result = OrderStatusCloseRefundReturnDelete
	case OrderStatusRefund:
		result = OrderStatusRefundDelete
	case OrderStatusUserCancel:
		result = OrderStatusUserCancelDelete
	case OrderStatusOverTime:
		result = OrderStatusOverTimeDelete
	}
	return result
}



//获取查询订单状态
func GetOrderStatusFromQuery(query int64) []int64 {
	if query == QueryUnPay {
		return []int64{OrderStatusUnPay}
	} else if query == QueryPay {
		return []int64{OrderStatusPay, OrderStatusPayRefund}
	} else if query == QuerySend {
		return []int64{OrderStatusSend, OrderStatusUnRefund}
	} else if query == QueryFinish {
		return []int64{OrderStatusFinish, OrderStatusFinishComment}
	} else if query == QueryFinishNotComment {
		return []int64{OrderStatusFinish}
	} else if query == QueryRefund {
		return []int64{OrderStatusCloseRefund,OrderStatusUnRefund,OrderStatusPayRefund, OrderStatusCloseRefundReturn, OrderStatusRefund}
	} else if query == QueryError {
		return []int64{OrderStatusUserCancel, OrderStatusOverTime}
	} else if query==  QueryFinance{
		return []int64{ OrderStatusPay, OrderStatusSend, OrderStatusFinish, OrderStatusUnRefund, OrderStatusPayRefund, OrderStatusFinishComment, OrderStatusCloseRefund, OrderStatusCloseRefundReturn, OrderStatusRefund}
	} else {
		return []int64{OrderStatusUnPay, OrderStatusPay, OrderStatusSend, OrderStatusFinish, OrderStatusUnRefund, OrderStatusPayRefund, OrderStatusFinishComment, OrderStatusUserCancel, OrderStatusOverTime, OrderStatusCloseRefund, OrderStatusCloseRefundReturn, OrderStatusRefund}
	}
}


func GetPayChannel(typeInt int64) string {
	result := ""
	switch typeInt {
	case WEIXINPAY:
		result = "微信支付"
	case ALIPAY:
		result = "支付宝支付"
	}
	return result
}

const (
	QueryOrderOrType = 1
	QueryOrderAndType =2
)

const (
	ShowOrderList = 1
	NoShowOrderList =0
)

// 订单类型
const (
	RefundStatusVerify int64 = 0 //申请中
	RefundStatusFinish int64 = 1 //完成
	RefundStatusFail   int64 = 2 //失败
)

// CreateRefundOrderSn 创建退款订单号
func CreateRefundOrderSn(orderType, buyType int64, code string) (result string) {
	var start, end string
	switch orderType {
	case OrderTypeApp:
		start = "01"
		end = "X"
	}
	result = fmt.Sprintf("%s%d%s%s", start, time.Now().Unix(), code, end)
	return
}
