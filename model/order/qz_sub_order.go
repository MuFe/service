package orderModel

import (
	"database/sql"
	"fmt"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"strings"
	"time"
)

// SubOrder 子订单表
type SubOrder struct {
	OrderID       int64
	SubOrderID    int64
	OrderAmount   int64
	PayAmount     int64
	DeliveryPrice int64
	OrderType     int64
	OrderStatus   int64
	MachineID     int64
	BrandID       int64
	RcTime        int64
	RcStatus      int64
	RcAdminNo     string
	Remark        string
	Consignee     string
	Phone         string
	Province      string
	City          string
	Area          string
	Address       string
	TransactionId string
	HaveCommented bool
	ShowOrder     bool

	OrderSN        string
	ExpressNumber  string
	ExpressCompany string
	Message        string
	OutTradeNo     string
	PayType        int64
	PayChannel     int64
	SellerID       int64
	BuyerId        int64
	BuyType        int64
	CreateTime     int64
	EndTime        int64
	PayTime        int64
	CouponPrice    int64


}

//创建子订单
func CreateSubOrder(tx *db.Tx, phone,province,city,area,address,name string, orderID, sellerId, price, deliverPrice int64) (int64, error) {
	result,err:=  tx.Exec(`INSERT INTO qz_sub_order (order_id,seller_id,order_amount,pay_amount,delivery_price,phone,province,city,area,address,consignee) VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
		orderID,
		sellerId,
		price,
		price,
		deliverPrice,
		phone,
		province,
		city,
		area,
		address,
		name,
	)
	if err!=nil{
		reErr := xlog.Error("生成子订单出错")
		xlog.ErrorP(reErr, err)
		return 0,reErr
	}
	id,err:=result.LastInsertId()
	if err!=nil{
		return 0,xlog.Error(err)
	}
	return id,nil
}


func EditExpress(data *BatchExpressData,subOrderId int64)error{
	_, err := db.GetOrderDb().Exec(`
UPDATE tb_sub_order 
SET express_number =?,express_company =?
WHERE
	sub_order_id =?`, data.ExpressNumber, data.ExpressCompany, subOrderId)
	if err != nil {
		return xlog.Error(err)
	}
	return nil
}

func EditExpressInfo(data *BatchExpressData,subOrderId,sendUid int64)error{
	sendTime := time.Now().Unix()
	_, err:= db.GetOrderDb().Exec(`
UPDATE qz_sub_order,
tb_order 
SET qz_sub_order.order_status = ?,
qz_order.order_status = ?,
qz_sub_order.express_number =?,
qz_sub_order.express_company =?,
qz_sub_order.send_uid =?,
qz_sub_order.send_time =? 
WHERE
	qz_sub_order.order_id = qz_order.order_id 
	AND qz_sub_order.sub_order_id =?`, enum.OrderStatusSend, enum.OrderStatusSend, data.ExpressNumber, data.ExpressCompany, sendUid, sendTime, subOrderId)
	if err != nil {
		return xlog.Error(err)
	}
	return nil
}

//获取订单列表
func GetOrderList( buyerId, subOrderId, orderID, startTime, queryEndTime, page, size,mBuyerType,queryType,showOrderType int64, statusList, orderTypes, skuIdList,sellerIdList, queryUidList []int64, queryPhone,queryExpressNumber,queryOrderSn string) ([]SubOrder, int64, error) {
	var buf strings.Builder
	total := int64(0)
	query := "" +
		"tb.sub_order_id," 		+ // 子订单id
		"tb.order_id," 			+ // 订单id
		"tbo.order_sn," 		+ // 订单编号
		"tb.order_status," 		+ // 订单状态 1下单未付款，2待发货状态，3待收货状态，4交易完成状态，5待退款状态，6交易关闭-仅退款，7交易关闭-退货退款（包含部分仅退款）8（退款）9交易关闭-用户取消，10交易关闭-交易超时，11订单删除
		"tb.order_amount," 		+ // 订单金额
		"tb.pay_amount," 		+ // 支付金额
		"tb.express_number," 	+ // 快递单号
		"tb.express_company," 	+ // 快递公司
		"tb.message," 			+ // 信息
		"tb.remark," 			+ // 备注
		"tb.consignee," 		+ // 收货人
		"tb.phone," 			+ // 手机
		"tb.province," 			+ // 省
		"tb.city," 				+ // 市
		"tb.area," 				+ // 区
		"tb.address," 			+ // 地址
		"tb.delivery_price," 	+ // 交易金额
		"tbo.create_time," 		+ // 创建时间
		"tbo.pay_time," 		+ // 支付时间
		"tbo.end_time," 		+ // 订单支付结束时间
		"tbo.transaction_id," 	+ // 交易id
		"tb.seller_id," 		+ // 卖家id
		"tbo.buyer_id," 		+ // 买家id，普通订单为用户id，批发订单为商家id
		"tbo.order_type," 		+ // 订单类型
		"tbo.buy_type," 		+ // 购买方式，1为小程序购买，2为终端购买，3为网站购买
		"tbo.pay_type," 		+ // 支付类型
		"tbo.show_order_list," 		+ // 是否显示商品订单
		"tbo.pay_channel" 		 // 支付方式，1为微信支付，2为支付宝支付



	buf.WriteString(`SELECT %s FROM
	qz_sub_order tb
	INNER JOIN qz_order tbo ON tbo.order_id = tb.order_id`)
	if len(skuIdList) > 0 {
		buf.WriteString(" INNER JOIN qz_order_good_detail tb2 ON tb.sub_order_id = tb2.sub_order_id ")
	}
	buf.WriteString(" Where 1=1 ")
	args := make([]interface{}, 0)
	utils.MysqlStringInUtils(&buf, statusList, " AND tb.order_status ")
	if subOrderId != 0 {
		buf.WriteString(" AND tb.sub_order_id=?")
		args = append(args, subOrderId)
	} else if orderID != 0 {
		buf.WriteString(" AND tb.order_id=?")
		args = append(args, orderID)
	} else if len(orderTypes) > 0 {
		utils.MysqlStringInUtilsWithZero(&buf, orderTypes, " AND tbo.order_type")
	}
	if len(sellerIdList) > 0 {
		utils.MysqlStringInUtils(&buf, sellerIdList, " AND tb.seller_id ")
	}
	if buyerId != 0 {
		buf.WriteString(" AND tbo.buyer_id=?")
		args = append(args, buyerId)
	}
	buf.WriteString(" AND tbo.show_order_list=?")
	args = append(args, showOrderType)
	if len(skuIdList) > 0 {
		utils.MysqlStringInUtils(&buf, skuIdList, " AND tb2.sku_id ")
	}

	if mBuyerType!=0{
		buf.WriteString(" AND tbo.buy_type=?")
		args = append(args, mBuyerType)
	}

	if queryType==enum.QueryOrderOrType{
		args = append(args, queryPhone)
		args = append(args, queryExpressNumber)
		args = append(args, queryOrderSn)
		if len(queryUidList) > 0 {
			utils.MysqlStringInUtilsWithZero(&buf, queryUidList, " AND (tb.phone=? or tb.express_number=? or tbo.order_sn=? or tbo.buyer_id ")
			buf.WriteString(" )")
		} else {
			buf.WriteString(" AND (tb.phone=? or tb.express_number=? or tbo.order_sn=?)")
		}
	} else if queryType==enum.QueryOrderAndType{
		utils.MysqlStringInUtilsWithZero(&buf,queryUidList," AND tbo.buyer_id")
		if queryPhone!=""{
			buf.WriteString(" AND tb.phone=?")
			args = append(args, queryPhone)
		}
		if queryExpressNumber!=""{
			buf.WriteString(" AND tb.express_number=?")
			args = append(args, queryExpressNumber)
		}
		if queryOrderSn!=""{
			buf.WriteString(" AND tbo.order_sn=?")
			args = append(args, queryOrderSn)
		}
	}
	if startTime != 0 {
		buf.WriteString(" AND tbo.create_time>?")
		args = append(args, startTime)
	}

	if queryEndTime > 0 {
		buf.WriteString(" AND tbo.create_time<?")
		args = append(args, queryEndTime)
	}

	if page == 1 {
		err := db.GetOrderDb().QueryRow(fmt.Sprintf(buf.String(), " count(tb.sub_order_id)"), args...).Scan(&total)
		if err != nil {
			total = 0
		}
	}
	buf.WriteString(" ORDER BY tbo.create_time DESC ")
	if size != 0 {
		buf.WriteString(" LIMIT ?,?")
		start := (page - 1) * size
		args = append(args, start)
		args = append(args, size)
	}
	result, err := db.GetOrderDb().Query(fmt.Sprintf(buf.String(), query), args...)
	if err != nil {
		return nil, total, xlog.Error(err)
	}
	var orderSn, expressNumber, expressCompany, message, adminMessage, consignee, phone, province, city, area, address, transactionId string
	var amount, pMount, orderTime, payTime, endTime, status, oId, orderId, bUid, deliveryPrice, orderType, payType, payChannel, oSellerId,buyType,showOrder int64
	list := make([]SubOrder, 0)
	for result.Next() {
		err = result.Scan(
			&oId,
			&orderId,
			&orderSn,
			&status,
			&amount,
			&pMount,
			&expressNumber,
			&expressCompany,
			&message,
			&adminMessage,
			&consignee,
			&phone,
			&province,
			&city,
			&area,
			&address,
			&deliveryPrice,
			&orderTime,
			&payTime,
			&endTime,
			&transactionId,
			&oSellerId,
			&bUid,
			&orderType,
			&buyType,
			&payType,
			&showOrder,
			&payChannel,
		)
		if err == nil {
			isShow:=true
			if showOrder==0{
				isShow=false
			}
			temp := SubOrder{
				OrderSN:        orderSn,
				OrderID:        orderId,
				OrderStatus:    status,
				Consignee:      consignee,
				Phone:          phone,
				Province:       province,
				City:           city,
				Area:           area,
				Address:        address,
				OrderAmount:    amount,
				PayAmount:      pMount,
				ExpressNumber:  expressNumber,
				ExpressCompany: expressCompany,
				Message:        adminMessage,
				Remark:         message,
				SubOrderID:     oId,
				CreateTime:     orderTime,
				PayTime:        payTime,
				EndTime:        endTime,
				TransactionId:  transactionId,
				HaveCommented:  status == enum.OrderStatusFinishComment,
				DeliveryPrice:  deliveryPrice,
				BuyerId:        bUid,
				PayType:        payType,
				PayChannel:     payChannel,
				SellerID:       oSellerId,
				OrderType:      orderType,
				BuyType:buyType,
				ShowOrder:isShow,
			}
			list = append(list, temp)
		} else {
			xlog.ErrorP(err)
			return nil, total, err
		}
	}
	return list, total, nil
}

// GetSubOrderByID 获取子订单信息
func GetSubOrderByID(subId,id int64) (*SubOrder, error) {
	result:=&SubOrder{}
	var orderID, subOrderID, orderStatus, uid, orderType, sellerID, brandID,
	endTime, createTime, orderAmount, payAmount, payType, payChannel, rcTime, rcStatus, bType sql.NullInt64
	var orderSN, rcAdminNo, remark, outNo sql.NullString
	buf:=strings.Builder{}
	args:=make([]interface{},0)
	buf.WriteString(`select o.order_id,so.sub_order_id,so.order_amount,so.order_status,o.create_time,
		o.buyer_id,so.pay_amount,o.order_type,so.seller_id,so.brand_id,o.end_time,
		o.pay_type,o.pay_channel,o.order_sn,so.rc_time,so.rc_status,so.rc_admin_no,so.remark,o.out_trade_no,o.buy_type from qz_sub_order so
		left join qz_order o on o.order_id = so.order_id `)
	if subId==0&&id==0{
		return nil,xlog.Error(errcode.HttpErrorWringParam.Msg)
	}
	buf.WriteString(" where 1=1 ")
	if subId!=0{
		buf.WriteString(" and so.sub_order_id = ?")
		args=append(args,subId)
	}
	if id!=0{
		buf.WriteString(" and o.order_id = ?")
		args=append(args,id)
	}
	err := db.GetOrderDb().QueryRow(
		buf.String(),args...).
		Scan(&orderID, &subOrderID, &orderAmount, &orderStatus, &createTime, &uid,
			&payAmount, &orderType, &sellerID, &brandID, &endTime, &payType, &payChannel, &orderSN,
			&rcTime, &rcStatus, &rcAdminNo, &remark, &outNo, &bType)
	if err != nil && err != sql.ErrNoRows {
		return nil, xlog.Error(err)
	}

	result.OrderSN = orderSN.String
	result.OrderID = orderID.Int64
	result.SubOrderID = subOrderID.Int64
	result.OrderAmount = orderAmount.Int64
	result.OrderStatus = orderStatus.Int64
	result.CreateTime = createTime.Int64
	result.BuyerId = uid.Int64
	result.PayAmount = payAmount.Int64
	result.OrderType = orderType.Int64
	result.SellerID = sellerID.Int64
	result.BrandID = brandID.Int64
	result.EndTime = endTime.Int64
	result.PayType = payType.Int64
	result.PayChannel = payChannel.Int64
	result.RcAdminNo = rcAdminNo.String
	result.RcStatus = rcStatus.Int64
	result.RcTime = rcTime.Int64
	result.Remark = remark.String
	result.HaveCommented = orderStatus.Int64 == enum.OrderStatusFinishComment
	result.OutTradeNo = outNo.String
	result.BuyType = bType.Int64

	return result, nil
}
