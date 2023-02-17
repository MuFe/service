package goodmodel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"mufe_service/jsonRpc"
	"strings"
)

func AddOption(infos []*app.Sku,spuId int64, tx *db.Tx)(error,map[string]*app.SkuOption){
	optionsMap := make(map[string]*app.SkuOption)
	delOptionsMap := make(map[int64]*app.SkuOption)
	for _, info := range infos {
		for _, temp := range info.Options {
			if temp.OptionValueId == 0 {
				_, ok := optionsMap[temp.Uuid]
				if !ok {
					optionsMap[temp.Uuid] = temp
				}
			} else if temp.Uuid == "" {
				_, ok := delOptionsMap[temp.OptionValueId]
				if !ok {
					delOptionsMap[temp.OptionValueId] = temp
				}
			} else {
				_,err := tx.Exec("update qz_good_specification_option set content =? where id=?", temp.OptionValue, temp.OptionValueId)
				if err != nil {
					return xlog.Error(err),optionsMap
				}
			}
		}
	}
	for _, info := range optionsMap {
		optionResult, err := tx.Exec("insert into qz_good_specification_option (content,spec_id,spu_id) values(?,?,?)", info.OptionValue, info.OptionId, spuId)
		if err != nil {
			return xlog.Error(err),optionsMap
		}
		optionInt, _ := optionResult.LastInsertId()
		info.OptionValueId = optionInt
	}

	for _, info := range delOptionsMap {
		_, err := tx.Exec("delete from qz_good_specification_option  where id=?", info.OptionValueId)
		if err!=nil{
			return xlog.Error(err),optionsMap
		}
	}
	return nil,optionsMap
}


type SkuOptions struct {
	Id      int64
	ValueId int64
	Value   string
	SkuId   int64
	SpuId   int64
}

//获取Sku规格
func GetSkuOptions(spuIdList, skuIdList []int64) (map[int64][]SkuOptions, error) {
	skuOptionMap := make(map[int64][]SkuOptions)
	var buf strings.Builder
	buf.WriteString(`select tb.sku_id,tb2.id,tb2.content,tb2.spec_id,tb.spu_id from qz_good_sku tb 
	left join qz_good_sku_specification tb1 on tb1.sku_id=tb.sku_id
	left join qz_good_specification_option tb2 on tb2.id=tb1.option_id`)
	if len(spuIdList) > 0 {
		utils.MysqlStringInUtils(&buf, spuIdList, " where tb.spu_id ")
	} else {
		utils.MysqlStringInUtils(&buf, skuIdList, " where tb.sku_id ")
	}
	spResult, err := db.GetGoodDb().Query(buf.String())
	if err == nil {
		var content sql.NullString
		var tId, skuId, specId, spuId sql.NullInt64
		for spResult.Next() {
			err = spResult.Scan(&skuId,&tId, &content, &specId,  &spuId)
			if err == nil {
				temp, ok := skuOptionMap[skuId.Int64]
				if !ok {
					temp = make([]SkuOptions, 0)
				}
				temp = append(temp, SkuOptions{
					Id:      specId.Int64,
					Value:   content.String,
					ValueId: tId.Int64,
					SkuId:   skuId.Int64,
					SpuId:   spuId.Int64,
				})
				skuOptionMap[skuId.Int64] = temp
			}
		}
		return skuOptionMap, nil
	} else {
		return skuOptionMap, xlog.Error(err)
	}
}
