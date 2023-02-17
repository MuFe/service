package goodmodel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"mufe_service/jsonRpc"
	"strings"
	"time"
)

type CreateStockData struct {
	SkuId           int64
	Num             int64
	BusinessId      int64
	BusinessGroupId int64
}

func UpdateStock(tx *db.Tx,info *app.Sku,businessId int64)(error,[][]int64,[]CreateStockData){
	stockMsgList := make([][]int64, 0)
	stockList := make([]CreateStockData, 0)
	var num sql.NullInt64
	err:= tx.QueryRow("select num from qz_good_sku_stock where sku_id=? and business_id=? and business_group_id=?", info.SkuId, businessId, enum.BRAND_ADMIN_GROUP).Scan(&num)
	if err == nil {
		dif := info.Stock - num.Int64
		_, err = tx.Exec("update qz_good_sku_stock set num=? where sku_id=? and business_id=? and business_group_id=?", info.Stock, info.SkuId, businessId, enum.BRAND_ADMIN_GROUP)
		if err != nil {
			return xlog.Error(err),nil,nil
		}
		if dif != 0 {
			stockMsgList = append(stockMsgList, []int64{info.SkuId, dif})
		}
	} else if err == sql.ErrNoRows {
		err = nil
		stockMsgList = append(stockMsgList, []int64{info.SkuId, info.Stock})
		stockList = append(stockList, CreateStockData{
			SkuId:           info.SkuId,
			Num:             info.Stock,
			BusinessGroupId: enum.BRAND_ADMIN_GROUP,
			BusinessId:      businessId,
		})
	} else {
		xlog.ErrorP(err)
	}
	return nil,stockMsgList,stockList
}


//创建库存
func CreateStock(stockList []CreateStockData, tx *db.Tx) error {
	//大于0添加库存记录消息
	sqlStr := "insert into qz_good_sku_stock (sku_id,num,modify_time,business_id,business_group_id) values"
	placeHolder := "(?,?,?,?,?)"

	var values []string
	var args []interface{}
	for _, info := range stockList {
		values = append(values, placeHolder)
		args = append(args, info.SkuId)
		args = append(args, info.Num)
		args = append(args, time.Now().Unix())
		args = append(args, info.BusinessId)
		args = append(args, info.BusinessGroupId)
	}
	sqlStr += strings.Join(values, ",")
	_, err := tx.Exec(sqlStr, args...)
	if err != nil {
		return xlog.Error(err)
	}
	return nil
}


// GetGoodStockBySkuIDs 获取sku商品库存
func GetGoodStockBySkuIDs(skuIDs []int64, businessId, businessGroupId int64) (map[int64]int64, error) {
	var result = make(map[int64]int64, 0)
	if len(skuIDs) == 0 {
		return result, nil
	}
	var args []interface{}
	var buf strings.Builder
	buf.WriteString(`select sku_id,num from qz_good_sku_stock where 1=1`)
	if businessId != 0 {
		buf.WriteString(" and business_id=? ")
		args = append(args, businessId)
	}
	if businessGroupId != 0 {
		buf.WriteString(" and business_group_id=? ")
		args = append(args, businessGroupId)
	}
	utils.MysqlStringInUtils(&buf, skuIDs, " and sku_id ")
	rows, err := db.GetGoodDb().Query(
		buf.String(), args...)
	if err != nil {
		return result, xlog.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		var skuID, num sql.NullInt64
		err := rows.Scan(&skuID, &num)
		if err != nil {
			return result, xlog.Error(err)
		}
		result[skuID.Int64] = num.Int64
	}
	err = rows.Err()
	if err != nil {
		return result, xlog.Error(err)
	}
	return result, nil
}


// GetAllGoodStockOnSpu 获取商品库存 spuIds有值就是获取指定的商品库存
func GetAllGoodStockOnSpu(spuIDs []int64) (result map[int64]int64, err error) {
	result = make(map[int64]int64, 0)
	var buf strings.Builder
	buf.WriteString(`select sum(s.num),sku.spu_id from qz_good_sku sku
		left join qz_good_sku_stock s on s.sku_id = sku.sku_id `)
	utils.MysqlStringInUtils(&buf, spuIDs, " where sku.spu_id")
	buf.WriteString(" and s.business_group_id=? group by sku.spu_id")
	rows, err := db.GetGoodDb().Query(buf.String(), enum.BRAND_ADMIN_GROUP)
	if err != nil {
		return result, xlog.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		var sum, spuID sql.NullInt64
		err := rows.Scan(&sum, &spuID)
		if err != nil {
			return result, xlog.Error(err)
		}
		result[spuID.Int64] = sum.Int64
	}
	err = rows.Err()
	if err != nil {
		return result, xlog.Error(err)
	}
	return result, nil
}
