package goodmodel

import (
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/xlog"
	"strings"
	"time"
)


func AddStockMessage(stockMsgList [][]int64,businessId int64, tx *db.Tx) error {
	if len(stockMsgList) > 0 {
		//大于0添加库存记录消息
		sqlStr := "insert into qz_good_sku_stock_message (sku_id,num,create_time,business_id,type) values"
		placeHolder := "(?,?,?,?,?)"

		var values []string
		var args []interface{}
		for _, info := range stockMsgList {
			values = append(values, placeHolder)
			args = append(args, info[0])
			args = append(args, info[1])
			args = append(args, time.Now().Unix())
			args = append(args, businessId)
			args = append(args, enum.SkuStockTypeBrandAdd)
		}
		sqlStr += strings.Join(values, ",")
		_, err := tx.Exec(sqlStr, args...)
		if err != nil {
			return xlog.Error(err)
		}
	}
	return nil
}
