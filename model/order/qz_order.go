package orderModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"strconv"
	"strings"
	"time"
)

type OrderData struct{
	SellerId int64
	DeliverPrice int64
	Price int64
	Detail []*OrderDetail
}

type OrderResponse struct {
	Id          int64
	TotalFree   int64
	Title       string
	Desc        string
	OrderSn     string
	Status      int64
	OrderType   int64
	BuyerId     int64
	PrepayId    string
	OutTradeNo  string
}

type ReceiveGoodData struct {
	SubOrderId  int64
	OrderStatus int64
	OrderId     int64
	OrderType   int64
	OrderSn     string
	Handler     string
}


type BatchExpressData struct{
	Id int64
	ExpressNumber string
	ExpressCompany string
}


func CreateOrder(orderSn,title,desc,phone,province,city,area,address,name string,list []OrderData,startTime,endTime ,buyerId,price,orderType,showList int64) (int64,error) {
	var id int64
	err := db.GetOrderDb().WithTransaction(func(tx *db.Tx) error {
		//生成订单
		orderID, err := createOrder(tx,  orderSn,title,desc, startTime, endTime,buyerId,price,orderType,showList)
		if err != nil {
			reErr := xlog.Error("生成订单出错")
			xlog.ErrorP(reErr, err)
			return reErr
		}
		subOrderMap := make(map[int64][]*OrderDetail)
		//分批生成子订单

		// 避免除数为0
		if price == 0 {
			price = 1
		}
		for _, temp := range list{
			subOrderId, err := CreateSubOrder(tx, phone,province,city,area,address,name, orderID, temp.SellerId, temp.Price, temp.DeliverPrice)
			if err != nil {
				reErr := xlog.Error("生成子订单出错")
				xlog.ErrorP(reErr, err)
				return reErr
			}
			subOrderMap[subOrderId] = temp.Detail
		}
		//生成订单详细
		err = CreateOrderGoodDetail(tx, subOrderMap, orderID)
		if err != nil {
			return err
		}
		//生成流水信息
		err = CreateAction(tx, orderID, buyerId)
		if err != nil {
			return err
		}
		id=orderID
		return nil
	})
	return id, err

}



//创建订单
func createOrder(tx *db.Tx, orderSn,title,desc string, startTime, endTime,buyerId,price,orderType,showList int64) (int64, error) {
	bType := enum.BuyTypeApp

	result,err:= tx.Exec(`INSERT INTO qz_order (buyer_id,order_sn,create_time,end_time,order_amount,pay_amount,pay_title,pay_desc,order_type,buy_type,out_trade_no,show_order_list) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`,
		buyerId,
		orderSn,
		startTime,
		endTime,
		price,
		price,
		title,
		desc,
		orderType,
		bType,
		orderSn,
		showList,
	)
	if err!=nil{
		reErr := xlog.Error("生成订单出错")
		xlog.ErrorP(reErr, err)
		return 0,reErr
	}
	id,err:=result.LastInsertId()
	if err!=nil{
		return 0,xlog.Error(err)
	}
	xlog.Info(id)
	return id,nil
}

//更新支付系统订单号
func UpdateOrderOutTradeNo(id int64, no string) error {
	_, err := db.GetOrderDb().Exec(`update qz_order set out_trade_no=? where order_id=?`, no, id)
	return err
}


func CreateOrderFromId(id int64, oSn string) (response *OrderResponse, err error) {
	var desc, title, orderSn, prepayId, oNo string
	var payFree, orderType, buyerId, status, orderId, bType int64
	err = db.GetOrderDb().QueryRow("select buy_type,order_id,prepay_id,order_sn,pay_title,pay_desc,pay_amount,order_type,buyer_id,order_status,out_trade_no from qz_order where order_id=? or order_sn=?", id, oSn).
		Scan(&bType, &orderId, &prepayId, &orderSn, &title, &desc, &payFree,  &orderType, &buyerId, &status, &oNo)
	if err != nil {
		return nil, err
	}
	if payFree <= 0 {
		payFree = 1
	}
	createOrderResponse := &OrderResponse{Id: orderId}
	createOrderResponse.Desc = desc
	createOrderResponse.Title = title
	createOrderResponse.TotalFree = payFree
	createOrderResponse.OrderSn = orderSn
	createOrderResponse.OrderType = orderType
	createOrderResponse.BuyerId = buyerId
	createOrderResponse.Status = status
	createOrderResponse.PrepayId = prepayId
	createOrderResponse.OutTradeNo = oNo
	return createOrderResponse, nil
}

func UpdateOrderPrepayId(id,payChannel int64, oSn, prepayId string) error {
	_, err := db.GetOrderDb().Exec(`update qz_order set prepay_id=?,pay_channel=? where order_id=? or order_sn=?`, prepayId,payChannel, id, oSn)
	return err
}

func CancelOrder(orderId, buyerId, selectStatus int64)error{
	//用户取消
	var count, status, orderType sql.NullInt64
	var orderSn sql.NullString
	var err error
	err = db.GetOrderDb().QueryRow("select count(order_id),order_status,order_type,order_sn from qz_order where order_id=? and buyer_id=?", orderId, buyerId).Scan(&count, &status, &orderType, &orderSn)
	if err != nil {
		return  err
	}
	if count.Int64 == 0 {
		return  xlog.Error("不是你的订单")
	}
	//取消订单流程
	 return Cancel(status.Int64, selectStatus, orderId)
}

func Pay(queryOrderId int64, queryOrderSn string) (string,int64,error) {
	var status, orderType, uid, orderId sql.NullInt64
	var orderSn sql.NullString
	err := db.GetOrderDb().QueryRow("select order_id,order_status,order_type,order_sn,buyer_id from qz_order where order_id=? or order_sn=?", queryOrderId, queryOrderSn).Scan(&orderId, &status, &orderType, &orderSn, &uid)
	if err != nil && err != sql.ErrNoRows {
		return "",0, nil
	}

	if err == sql.ErrNoRows {
		//已经消费过了
		return "",0, nil
	}


	if status.Int64 != enum.OrderStatusUnPay {
		//订单状态有误
		return "",0, xlog.Error("订单状态有误")
	}

	selectStatus := enum.OrderStatusPay
	err = db.GetOrderDb().WithTransaction(func(tx *db.Tx) error {
		list := make([]string, 0)
		list = append(list, strconv.FormatInt(orderId.Int64, 10))
		err = EditOrderStatus(list, tx, selectStatus, enum.OrderStatusUnPay)
		if err != nil {
			xlog.ErrorP(err)
			return err
		}
		return nil
	})
	if err != nil {
		return "",0, err
	}
	return orderSn.String,orderId.Int64, nil
}


func Cancel(status, selectStatus, orderId int64) error {
	if status == enum.OrderStatusUserCancel || status == enum.OrderStatusOverTime {
		//已经消费过了
		return nil
	}
	if status != enum.OrderStatusUnPay {
		//订单状态有误
		return nil
	}
	err:= db.GetOrderDb().WithTransaction(func(tx *db.Tx) error {
		list := make([]string, 0)
		list = append(list, strconv.FormatInt(orderId, 10))
		err := EditOrderStatus(list, tx, selectStatus, enum.OrderStatusUnPay)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func EditOrderStatus(list []string, tx *db.Tx, status, oldStatus int64) error {
	var err error
	if (status == enum.OrderStatusPay || status == enum.OrderStatusSend) && oldStatus == enum.OrderStatusUnPay {
		//修改为支付状态
		_, err = tx.Exec("update qz_order set order_status=?,pay_time=? where order_status=? and order_id in ("+strings.Join(list, ",")+")", status, time.Now().Unix(), oldStatus)
		if err != nil {
			return err
		}
	} else {
		_, err = tx.Exec("update qz_order set order_status=? where order_status=? and order_id in ("+strings.Join(list, ",")+")", status, oldStatus)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec("update qz_sub_order set order_status=? where order_status=? and order_id in ("+strings.Join(list, ",")+")", status, oldStatus)
	if err != nil {
		return err
	}
	return nil
}



//发货流程
func SendGood(sendUid, businessId int64, list []*BatchExpressData)  error {
	var status, orderId, sellerId, subOrderId int64
	var orderSn, expressNumber, expressCompany string
	subOrderIds := make([]string, 0)
	batchExpressMap := make(map[int64]*BatchExpressData)
	for _, info := range list {
		subOrderIds = append(subOrderIds, strconv.FormatInt(info.Id, 10))
		batchExpressMap[info.Id] = info
	}
	result, err := db.GetOrderDb().Query(`select tb.sub_order_id,tb.order_id,tb.order_status,tb1.order_sn,tb.seller_id,tb.express_number,tb.express_company from qz_sub_order  tb inner join qz_order tb1 on tb1.order_id=tb.order_id where tb.sub_order_id in (` + strings.Join(subOrderIds, ",") + ")")
	if err == nil {
		for result.Next() {
			err = result.Scan(&subOrderId, &orderId, &status, &orderSn, &sellerId, &expressNumber, &expressCompany)
			if err == nil {
				if sellerId != businessId {
					return  xlog.Error("对不起，您没有操作该订单的权限")
				} else {
					isEditExpress := false
					data, _ := batchExpressMap[subOrderId]
					if len(subOrderIds) == 1 {
						if status == enum.OrderStatusSend {
							if expressNumber == data.ExpressNumber && expressCompany == data.ExpressCompany {
								return nil
							} else {
								isEditExpress = true
							}
						}
						if status != enum.OrderStatusPay {
							//订单状态有误
							return  xlog.Error("订单状态有误")
						}
					} else {
						if status == enum.OrderStatusSend {
							if expressNumber == data.ExpressNumber && expressCompany == data.ExpressCompany {
								continue
							} else {
								isEditExpress = true
							}
						}
						if status != enum.OrderStatusPay {
							//订单状态有误
							continue
						}
					}
					if isEditExpress {
						//只修改物流信息
						err=EditExpress(data,subOrderId)
						if err != nil {
							return err
						}
					} else {
						err=EditExpressInfo(data,subOrderId,sendUid)
						if err != nil {
							return err
						}
					}

				}
			}
		}
	} else {
		return  err
	}
	return nil
}


func StartReceiveGood(subOrderId, buyerId, orderId int64)  error {
	var total, status sql.NullInt64
	var err error
	if orderId != 0 {
		err = db.GetOrderDb().QueryRow("select count(tb.sub_order_id),tb.order_status,tb.sub_order_id from qz_sub_order tb inner join qz_order tb1 on tb1.order_id =tb.order_id where tb.order_id=? and tb1.buyer_id=?", orderId, buyerId).Scan(&total, &status, &subOrderId)
		if err != nil {
			return  err
		}
	} else if buyerId != 0 {
		err = db.GetOrderDb().QueryRow("select count(tb.sub_order_id),tb.order_status from qz_sub_order tb inner join qz_order tb1 on tb1.order_id =tb.order_id where tb.sub_order_id=? and tb1.buyer_id=?", subOrderId, buyerId).Scan(&total, &status)
		if err != nil {
			return  err
		}
	}
	if total.Int64 == 0 {
		return xlog.Error("对不起，您没有操作该订单的权限")
	}
	return  ReceiveGoods( []int64{subOrderId})
}



//确认收货
func ReceiveGoods(subOrderId []int64) error {
	list := make([]ReceiveGoodData, 0)
	//获取相关订单信息
	var queryBuf strings.Builder
	queryBuf.WriteString(`select tb1.order_status,tb.sub_order_id,tb1.order_type,tb1.order_sn,tb.order_id,tb1.buyer_id from qz_sub_order tb inner join qz_order tb1 on tb1.order_id=tb.order_id`)
	utils.MysqlStringInUtils(&queryBuf, subOrderId, " where tb.sub_order_id")
	result, err := db.GetOrderDb().Query(queryBuf.String())
	if err == nil {
		var orderStatus, subOrderId, orderId, orderType, buyerId int64
		var orderSn string
		for result.Next() {
			err = result.Scan(&orderStatus, &subOrderId, &orderType, &orderSn, &orderId, &buyerId)
			if err == nil {
				list = append(list, ReceiveGoodData{
					SubOrderId:  subOrderId,
					OrderSn:     orderSn,
					OrderStatus: orderStatus,
					OrderId:     orderId,
					OrderType:   orderType,
				})
			} else {
				xlog.ErrorP(err)
			}
		}
	}

	for _, v := range list {
		err := receiveGood(v.SubOrderId, v.OrderStatus, v.OrderId)
		if err == nil {

		} else {

		}
	}
	return  nil
}

func receiveGood(subOrderId, status, orderId int64) error {
	if status == enum.OrderStatusFinish {
		//已经消费过了
		return nil
	}
	if status != enum.OrderStatusSend {
		//订单状态有误
		return xlog.Error("订单状态有误")
	}
	err := db.GetOrderDb().WithTransaction(func(tx *db.Tx) error {
		_, err := tx.Exec("update qz_sub_order set order_status=? where sub_order_id=?", enum.OrderStatusFinish, subOrderId)
		if err != nil {
			return err
		}
		countResult, err := tx.Query(`select tb1.order_status from qz_sub_order tb 
   						  inner join qz_sub_order tb1 on tb1.order_id=tb.order_id
						  where tb1.sub_order_id=?`, subOrderId)
		isFinish := true
		if err == nil {
			var status int64
			for countResult.Next() {
				err = countResult.Scan(&status)
				if err == nil {
					if status == enum.OrderStatusPay || status == enum.OrderStatusSend {
						isFinish = false
					}
				} else {
					return xlog.Error(err)
				}
			}
			if isFinish {
				_, err = tx.Exec("update qz_order set order_status=? where order_id=?", enum.OrderStatusFinish, orderId)
			}
		} else {
			xlog.ErrorP(err)
		}
		return nil
	})
	return err
}

func DelOrder(subOrderId, buyerId int64) error {
	//用户取消
	var count int64
	var orderStatus sql.NullInt64
	err := db.GetOrderDb().QueryRow("select count(tb.order_id),tb1.order_status from qz_sub_order tb inner join qz_order tb1 on tb1.order_id=tb.order_id where tb.sub_order_id=? and tb1.buyer_id=?", subOrderId, buyerId).Scan(&count, &orderStatus)
	if err != nil {
		xlog.ErrorP(err)
		return err
	}
	if count == 0 {
		return xlog.Error("不是你的订单")
	}
	selectStatus := enum.DeleteStatus(orderStatus.Int64)
	//订单超时
	err = db.GetOrderDb().WithTransaction(func(tx *db.Tx) error {
		_, err := tx.Exec("update qz_sub_order set order_status=? where sub_order_id=?", selectStatus, subOrderId)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
