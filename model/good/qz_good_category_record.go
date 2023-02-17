package goodmodel

import (
	"mufe_service/camp/db"
	"mufe_service/camp/xlog"
)

func AddGoodCategoryRecord(tx *db.Tx,spuId,categoryId int64)error{
	_, err := tx.Exec("insert into qz_good_category_record (spu_id,category_id) values(?,?)", spuId, categoryId)
	if err != nil {
		return xlog.Error(err)
	}else{
		return nil
	}
}

