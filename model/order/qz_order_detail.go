package orderModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"strings"
)

type OrderDetail struct{
	ID             int64
	SkuId int64
	Num int64
	SubOrderID     int64
	Price          int64
	PayAmount      int64
	OrderID        int64
	AgreementPrice int64
}

//创建订单详细
func CreateOrderGoodDetail(tx *db.Tx, dataMap map[int64][]*OrderDetail, orderID int64) error {
	sqlStr := "insert into qz_order_good_detail (order_id,sku_id,sub_order_id,num,price,pay_price) values"
	placeHolder := "(?,?,?,?,?,?)"

	var values []string
	var args []interface{}
	for subOrderId, goods := range dataMap {
		for _, good := range goods {
			values = append(values, placeHolder)
			args = append(args, orderID)
			args = append(args, good.SkuId)
			args = append(args, subOrderId)
			args = append(args, good.Num)
			args = append(args, good.Price)
			args = append(args, good.Num*good.Price)
		}
	}
	sqlStr += strings.Join(values, ",")
	_, err := tx.Exec(sqlStr, args...)
	if err != nil {
		reErr := xlog.Error("生成订单详细出错")
		xlog.ErrorP(reErr, err)
		return reErr
	}
	return nil
}

// GetOrderGoodDetailBySubOrderID 获取订单详情
func GetOrderGoodDetailBySubOrderID(ids []int64) (result []OrderDetail, err error) {
	var buf strings.Builder
	buf.WriteString(`select ogd.id,ogd.sku_id,ogd.num,ogd.sub_order_id,ogd.price,ogd.order_id,ogd.pay_price
		from qz_order_good_detail ogd`)
	utils.MysqlStringInUtils(&buf, ids, " where ogd.sub_order_id ")
	rows, err := db.GetOrderDb().Query(buf.String())
	if err != nil {
		return result, xlog.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, skuID, num, subOrderID, price, orderID, payPrice sql.NullInt64
		err := rows.Scan(&id, &skuID, &num, &subOrderID, &price, &orderID, &payPrice)
		if err != nil {
			return result, xlog.Error(err)
		}
		result = append(result, OrderDetail{
			ID:             id.Int64,
			SkuId:          skuID.Int64,
			Num:            num.Int64,
			SubOrderID:     subOrderID.Int64,
			Price:          price.Int64,
			OrderID:        orderID.Int64,
			PayAmount:      payPrice.Int64,
		})
	}
	err = rows.Err()
	if err != nil {
		return result, xlog.Error(err)
	}
	return
}


// GetAllSpuSales 获取销量,spuID-num
func GetAllSpuSales(spuMap map[int64]int64, refresh bool) (result map[int64]int64, skuResult map[int64]int64, err error) {
	result = make(map[int64]int64, 0)
	result, skuResult, err = getSpuSales(spuMap)
	if err != nil {
		return nil, nil, err
	}
	return result, skuResult, nil
}


// getSpuSales 获取销量 // sku-spu
func getSpuSales(spuMap map[int64]int64) (map[int64]int64, map[int64]int64, error) {
	var spuResult = make(map[int64]int64)
	var skuResult = make(map[int64]int64)
	var buf strings.Builder
	buf.WriteString(`select de.sku_id,sum(de.num) from qz_order_good_detail de
left join  qz_sub_order sub on sub.sub_order_id = de.sub_order_id
where sub.order_status = ?`)
	if len(spuMap) > 0 {
		list := make([]int64, 0)
		for skuId := range spuMap {
			list = append(list, skuId)
		}
		utils.MysqlStringInUtils(&buf, list, " and de.sku_id ")
	}
	buf.WriteString(" group by de.sku_id")
	rows, err := db.GetOrderDb().Query(buf.String(), enum.OrderStatusFinish)
	if err != nil {
		return nil, nil, xlog.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		var skuID, num sql.NullInt64
		err := rows.Scan(&skuID, &num)
		if err != nil {
			return nil, nil, xlog.Error(err)
		}
		spuID := spuMap[skuID.Int64]
		if spuID > 0 {
			if v, ok := spuResult[spuID]; ok {
				spuResult[spuID] = v + num.Int64
			} else {
				spuResult[spuID] = num.Int64
			}
			skuResult[skuID.Int64] = num.Int64
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, nil, xlog.Error(err)
	}

	return spuResult, skuResult, nil
}
