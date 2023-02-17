package goodmodel

import (
	"mufe_service/camp/db"
	"mufe_service/jsonRpc"
)

func AddBrand(businessId int64, brandName string) (*app.EmptyResponse, error) {
	_, err := db.GetGoodDb().Exec("insert into qz_good_brand (business_id,brand_name,status) values(?,?,?)", businessId, brandName, 1)
	return &app.EmptyResponse{}, err
}
func AddBrandRecord(businessId,spuId int64,tx *db.Tx) error {
	_, err:= tx.Exec("insert into qz_good_brand_record (spu_id,brand_id) select ?,brand_id from tb_good_brand where business_id=?", spuId, businessId)
	return err
}
