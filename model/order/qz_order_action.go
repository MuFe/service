package orderModel

import (
	"time"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/xlog"
)

func CreateAction(tx *db.Tx, orderID, buyerId int64) error {
	_, err := tx.Exec("insert into qz_order_action (order_id,buyer_id,order_status,pay_status,action_note,status_desc,create_time) values(?,?,?,?,'支付下单',?,?)", orderID, buyerId, enum.OrderStatusUnPay, enum.OrderActionStatusUnPay, enum.GetOrderStatusMeaning(enum.OrderStatusUnPay), time.Now().Unix())
	if err != nil {
		reErr := xlog.Error("生成流水出错")
		xlog.ErrorP(reErr, err)
		return reErr
	}
	return nil
}
