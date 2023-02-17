package goodmodel

import (
	"mufe_service/camp/db"
	"mufe_service/camp/xlog"
	"strings"
)

func DelSkuOption(skuId int64, tx *db.Tx)error{
	_, err := tx.Exec("delete from qz_good_sku_specification where sku_id=?", skuId)
	return err
}

func AddSkuOption(skuId int64,list [][]int64, tx *db.Tx)error{

	sqlStr := "insert into qz_good_sku_specification (sku_id,spec_id,option_Id) values "
	placeHolder := "(?,?,?)"

	var values []string
	var args []interface{}
	for _, info := range list {
		values = append(values, placeHolder)
		args = append(args, skuId)
		args = append(args, info[0])
		args = append(args, info[1])
	}
	sqlStr += strings.Join(values, ",")
	_, err := tx.Exec(sqlStr, args...)
	if err != nil {
		return xlog.Error(err)
	}else{
		return nil
	}
}
