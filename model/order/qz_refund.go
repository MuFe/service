package orderModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"mufe_service/camp/enum"
	app "mufe_service/jsonRpc"
	"strconv"
	"strings"
	"time"
)

// Refund 退款
type Refund struct {
	ID          int64
	SubOrderID  int64
	Amount      int64
	CreateTime  int64
	Status      int64
	AdminNo     string
	Remark      string
	OutRefundNo string
	EndTime     int64
}

// RefundDetail 详情
type RefundDetail struct {
	ID       int64
	RefundID int64
	OgID     int64
	Num      int64
}

// GetRefundBySubOrderID 获取子订单的商品详情
func GetRefundBySubOrderID(id int64) (result []Refund, err error) {
	rows, err := db.GetOrderDb().Query(
		`select r.refund_id,r.sub_order_id,r.amount,r.create_time,r.status,r.admin_no,r.remark,r.end_time,r.out_refund_no
		from qz_refund r where r.sub_order_id = ? order by r.refund_id asc`, id)
	if err != nil {
		return result, xlog.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		var refundID, subOrderID, amount, createTime, status, endTime sql.NullInt64
		var adminNo, remark, outRefundNo sql.NullString
		err := rows.Scan(&refundID, &subOrderID, &amount, &createTime, &status, &adminNo, &remark, &endTime, &outRefundNo)
		if err != nil {
			return result, xlog.Error(err)
		}
		result = append(result, Refund{
			ID:          refundID.Int64,
			SubOrderID:  subOrderID.Int64,
			Amount:      amount.Int64,
			CreateTime:  createTime.Int64,
			Status:      status.Int64,
			AdminNo:     adminNo.String,
			Remark:      remark.String,
			EndTime:     endTime.Int64,
			OutRefundNo: outRefundNo.String,
		})
	}
	err = rows.Err()
	if err != nil {
		return result, xlog.Error(err)
	}

	return
}

// GetRefundDetailByRefundID 获取退款详情
func GetRefundDetailByRefundID(ids []int64) (result []RefundDetail, err error) {
	if len(ids) == 0 {
		return result, nil
	}
	ids = utils.RemoveDupsInt64(ids)
	args := []interface{}{}
	ph := []string{}
	for i, v := range ids {
		args = append(args, v)
		ph = append(ph, "?")
		if (i+1)%enum.SQLPlaceholderLimit == 0 || i == len(ids)-1 {
			rows, err := db.GetOrderDb().Query(
				`select rd.id,rd.refund_id,rd.og_id,rd.num from qz_refund_detail rd where rd.refund_id in (`+strings.Join(ph, ",")+`) order by rd.id asc`, args...)
			if err != nil {
				return result, xlog.Error(err)
			}
			defer rows.Close()
			for rows.Next() {
				var id, refundID, ogID, num sql.NullInt64
				err := rows.Scan(&id, &refundID, &ogID, &num)
				if err != nil {
					return result, xlog.Error(err)
				}
				result = append(result, RefundDetail{
					ID:       id.Int64,
					RefundID: refundID.Int64,
					OgID:     ogID.Int64,
					Num:      num.Int64,
				})
			}
			err = rows.Err()
			if err != nil {
				return result, xlog.Error(err)
			}
			args = []interface{}{}
			ph = []string{}
		}
	}

	return
}

//获取商家的退款信息
func GetBusinessRefund(businessId, lastTime int64, orderTypes []int64) ([]Refund, error) {
	// 退款记录
	var buf strings.Builder
	buf.WriteString(`select tr.end_time,tr.amount 
from qz_refund tr 
left join qz_sub_order tso on tr.sub_order_id = tso.sub_order_id 
left join qz_order tbo on tbo.order_id = tso.order_id 
where tso.seller_id = ? and tr.status = 1 and tr.end_time >= ? `)

	for k, typeInt := range orderTypes {
		if k == 0 {
			buf.WriteString(" and tbo.order_type in (")
		}
		buf.WriteString(strconv.FormatInt(typeInt, 10))
		if k == len(orderTypes)-1 {
			buf.WriteString(") ")
		} else {
			buf.WriteString(",")
		}
	}

	rows, err := db.GetOrderDb().Query(buf.String(), businessId, lastTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var memberRefund = make([]Refund, 0)
	for rows.Next() {
		var endTime, amount sql.NullInt64
		err := rows.Scan(&endTime, &amount)
		if err != nil {
			return nil, err
		}
		var data = Refund{
			EndTime: endTime.Int64,
			Amount:  amount.Int64,
		}

		memberRefund = append(memberRefund, data)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return memberRefund, nil
}

// PutOrderRefundStatus 退款退货审核
func PutOrderRefundStatus(status, orderStatus, subOrderId int64, remark, adminNo string) error {
	err := db.GetOrderDb().WithTransaction(func(tx *db.Tx) error {
		_, err := tx.Exec(`update qz_refund set status = ?,admin_no = ?,remark = ?,end_time = ? where sub_order_id = ? and status=?`,
			status, adminNo, remark, time.Now().Unix(), subOrderId, enum.RefundStatusVerify)
		if err != nil {
			return xlog.Error(err)
		}
		// 更新订单状态
		_, err = tx.Exec(`update qz_sub_order so set so.order_status = ? where so.sub_order_id = ?`,
			orderStatus, subOrderId)
		if err != nil {
			return xlog.Error(err)
		}
		return nil
	})
	return err
}

//取消退款申请
func CancelRefund(subOrderStatus, subOrderId int64) error {
	err := db.GetOrderDb().WithTransaction(func(tx *db.Tx) error {
		// 更新订单状态
		orderStatus := enum.OrderStatusSend
		if subOrderStatus == enum.OrderStatusPayRefund {
			orderStatus = enum.OrderStatusPay
		}
		_, err := tx.Exec(`update qz_refund set status = ?,admin_no = ?,remark = ?,end_time = ? where sub_order_id = ?`,
			enum.RefundStatusFail, "", "用户取消", time.Now().Unix(), subOrderId)
		if err != nil {
			return xlog.Error(err)
		}
		// 更新订单状态
		_, err = tx.Exec(`update qz_sub_order so set so.order_status = ? where so.sub_order_id = ?`,
			orderStatus, subOrderId)
		if err != nil {
			return xlog.Error(err)
		}
		return nil
	})
	return err
}

//新增退款申请
func CreateRefund(subOrderId, subOrderStatus, amount,refundMethod,refundType int64, list []*app.PostRefundReqList,photos []*app.PostRefundPhoto, refundSn,reason,expressCompany,expressNumber,explain string) error {
	// 新增
	var now = time.Now().Unix()
	var end = time.Now().AddDate(0, 0, 2).Unix()
	err := db.GetOrderDb().WithTransaction(func(tx *db.Tx) error {

		dbResult, err := tx.Exec("INSERT INTO qz_refund  (amount,create_time,end_time,`status`,sub_order_id,out_refund_no,reason,express_company,express_number,refund_method,refund_type,"+
		"`EXPLAIN`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)",
			amount,
			now,
			end,
			enum.RefundStatusVerify,
			subOrderId,
			refundSn,
			reason,
			expressCompany,
			expressNumber,
			refundMethod,
			refundType,
			explain,
		)
		if err != nil {
			return xlog.Error(err)
		}
		insertID, err := dbResult.LastInsertId()
		if err != nil {
			return xlog.Error(err)
		}
		var args = []interface{}{}
		var ph = []string{}
		for _, r := range list {
			args = append(args, r.Num)
			args = append(args, r.OgID)
			args = append(args, insertID)
			ph = append(ph, "(?,?,?)")
		}
		_, err = tx.Exec(`INSERT INTO qz_refund_detail (num, og_id, refund_id) VALUES `+strings.Join(ph, ","), args...)
		if err != nil {
			return xlog.Error(err)
		}
		// 更新订单状态
		orderStatus := enum.OrderStatusUnRefund
		if subOrderStatus == enum.OrderStatusPay {
			orderStatus = enum.OrderStatusPayRefund
		}
		_, err = tx.Exec(`update qz_sub_order so set so.order_status = ? where so.sub_order_id = ?`, orderStatus, subOrderId)
		if err != nil {
			return xlog.Error(err)
		}
		args = []interface{}{}
		ph = []string{}
		for _, r := range photos {
			for _,v:=range r.List{
				args = append(args, v)
				args = append(args, r.Type)
				args = append(args, insertID)
				ph = append(ph, "(?,?,?)")
			}
		}
		if len(ph)>0{
			_, err = tx.Exec(`INSERT INTO qz_refund_photo (url, type, refund_id) VALUES `+strings.Join(ph, ","), args...)
			if err != nil {
				return xlog.Error(err)
			}
		}
		return nil
	})
	return err
}

//获取退款信息
func GetRefundFromOrder(orderId int64, orderSn string, haveRefundVerify bool) map[int64]int64 {
	status := enum.RefundStatusFinish
	if haveRefundVerify {
		status = enum.RefundStatusVerify
	}
	refundMap := make(map[int64]int64)
	refundRow, err := db.GetOrderDb().Query(`SELECT
	tb1.og_id,
	tb1.num
FROM
	qz_refund tb
	inner join qz_refund_detail tb1 on tb1.refund_id=tb.refund_id
	inner join qz_sub_order tb2 on tb2.sub_order_id =tb.sub_order_id
	inner join qz_order tb3 on tb3.order_id=tb2.order_id
	where tb.status=? and (tb3.order_id =? or tb3.order_sn=?)`, status, orderId, orderSn)
	if err == nil {
		var ogId, num int64
		for refundRow.Next() {
			err = refundRow.Scan(&ogId, &num)
			if err == nil {
				refundMap[ogId] = num
			}
		}
	}
	return refundMap
}
