package goodmodel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/sequence"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"mufe_service/jsonRpc"
	"strings"
)

// Sku 具体规格的商品
type Sku struct {
	SkuID          int64
	SkuName        string
	SpuName        string
	Status         int64
	SpuID          int64
	Price          int64
	BusinessId     int64
	Num            int64
	CreateTime     int64
}

type OrderSku struct {
	List               []Sku
	ShopDeliverPrice   map[int64]int64
}

const insertSkuSql = `
INSERT INTO qz_good_sku (
sku_name,
sku_number,
spu_id,
price,
agreement_price,
cost_price,
member_price,
STATUS 
	)
VALUES
	( ?,?,?,?,?,?,?,? )
`

func EditSku(infos []*app.Sku, spuId int64, name string, categoryName string, status int64, tx *db.Tx, businessId int64) error {
	var err error
	var skuId int64
	stockMsgList := make([][]int64, 0)
	stockList := make([]CreateStockData, 0)
	err,optionsMap:=AddOption(infos,spuId,tx)
	if err!=nil{
		return xlog.Error(err)
	}
	for _, info := range infos {
		if info.SkuId != 0 {
			if info.IsDel {
				err =delSku(info.SkuId,tx)
				err=DelSkuOption(info.SkuId,tx)
			} else {
				err=updateSku(info,tx)
				_,tempStock,tmpList:=UpdateStock(tx,info,businessId)
				stockMsgList = append(stockMsgList, tempStock...)
				stockList = append(stockList,tmpList...)
			}
			if err != nil {
				return xlog.Error(err)
			}
			skuId = info.SkuId
		} else {
			skuNo := ""
			if categoryName != "" {
				skuNo, err = sequence.SkuNo.NewNo(categoryName)
				if err != nil {
					skuNo = ""
				}
			}
			err,skuId=addSku(info,name,skuNo,spuId,status,tx)
			if err != nil {
				return xlog.Error(err)
			}
			stockMsgList = append(stockMsgList, []int64{skuId, info.Stock})
			stockList = append(stockList, CreateStockData{
				SkuId:           skuId,
				Num:             info.Stock,
				BusinessGroupId: enum.BRAND_ADMIN_GROUP,
				BusinessId:      businessId,
			})
		}
		list:=make([][]int64,0)
		for _, option := range info.Options {
			temp, ok := optionsMap[option.Uuid]
			if ok {
				tempData:=make([]int64,0)
				tempData=append(tempData, option.OptionId,temp.OptionValueId)
				list=append(list,tempData)
			}
		}

		err=AddSkuOption(skuId,list,tx)
		if err!=nil{
			return nil
		}
	}
	if len(stockList) > 0 {
		err := CreateStock(stockList, tx)
		if err != nil {
			return xlog.Error(err)
		}
	}

	err=AddStockMessage(stockMsgList,businessId,tx)
	if err != nil {
		return err
	}
	return nil
}

func delSku(skuId int64,tx *db.Tx)error{
	_, err := tx.Exec("update qz_good_sku set status =? where sku_id=?", enum.StatusDelete,skuId)
	return err
}

func updateSku(info *app.Sku,tx *db.Tx)error{
	_, err:= tx.Exec("update qz_good_sku set price=?,agreement_price=?,cost_price=?,member_price=? where sku_id=?", info.Price, info.AgreementPrice, info.CostPrice, info.Price,  info.SkuId)
	return err
}

func addSku(info *app.Sku,name,skuNo string,spuId,status int64,tx *db.Tx)(error,int64){
	skuResult, err := tx.Exec(
		insertSkuSql,
		name,
		skuNo,
		spuId,
		info.Price,
		info.AgreementPrice,
		info.CostPrice,
		info.Price,
		status,
	)
	if err != nil {
		return xlog.Error(err),0
	}
	skuIdInt, _ := skuResult.LastInsertId()
	skuId:= int64(skuIdInt)
	return nil,skuId
}

func GetSpuIdFromSkuId(skuId int64) int64 {
	spuId := int64(0)
	err := db.GetGoodDb().QueryRow(`select spu_id from qz_good_sku where sku_id=?`, skuId).Scan(&spuId)
	if err != nil {
		xlog.ErrorP(err)
	}
	return spuId
}

func GetSkuFromQuery(spuIDs, skuIDs []int64, status int64) ([]Sku, error) {
	var buf strings.Builder
	arg := make([]interface{}, 0)
	buf.WriteString(`select sku.sku_id,sku.sku_name,sku.status,sku.spu_id,sku.price,sku.member_price,sku.agreement_price,sku.cost_price,spu.create_time,spu.spu_name
						from qz_good_sku sku inner join qz_good_spu spu on spu.spu_id=sku.spu_id `)
	buf.WriteString(" where 1=1 ")
	utils.MysqlStringInUtils(&buf, spuIDs, " AND sku.spu_id")
	utils.MysqlStringInUtils(&buf, skuIDs, " AND sku.sku_id")
	if status != enum.StatusAll {
		buf.WriteString(" AND sku.status=?")
		arg = append(arg, status)
	}
	buf.WriteString(" order by sku.sku_id asc")
	rows, err := db.GetGoodDb().Query(buf.String(), arg...)
	if err != nil {
		return nil, xlog.Error(err)
	}
	result := make([]Sku, 0)
	defer rows.Close()
	for rows.Next() {
		var skuID, status, spuID, price, memberPrice, agreementPrice, costPrice, createTime sql.NullInt64
		var name, spuName sql.NullString
		err := rows.Scan(&skuID, &name, &status, &spuID, &price, &memberPrice, &agreementPrice, &costPrice, &createTime, &spuName)
		if err != nil {
			return nil, xlog.Error(err)
		}
		result = append(result, Sku{
			SkuID:          skuID.Int64,
			SkuName:        name.String,
			SpuName:        spuName.String,
			Status:         status.Int64,
			SpuID:          spuID.Int64,
			Price:          price.Int64,
			CreateTime:     createTime.Int64,
		})
	}
	return result, nil
}

// GetSkuIDBySpuIDs  通过spu下的sku商品
func GetSkuIDBySpuIDs(spuIDs []int64, status int64) (result []Sku, err error) {
	return GetSkuFromQuery(spuIDs, []int64{}, status)
}

// GetAllSku 获取全部sku
func GetAllSku() (result []Sku, err error) {
	return GetSkuFromQuery([]int64{}, []int64{}, enum.StatusAll)
}

// GetSkuBySkuIDs 通过id获取sku信息
func GetSkuBySkuIDs(ids []int64) ([]Sku, error) {
	return GetSkuFromQuery([]int64{}, ids, enum.StatusAll)
}


func GetSkuOrderInfo(idMap map[int64]int64) (*OrderSku, error) {
	var buf, buf1 strings.Builder
	ids := make([]int64, 0)
	for k := range idMap {
		ids = append(ids, k)
	}
	buf.WriteString(`SELECT tb.sku_id,tb.price,tb3.business_id,tb.spu_id from qz_good_sku tb
	INNER JOIN qz_good_brand_record tb2 on tb2.spu_id=tb.spu_id
	INNER JOIN qz_good_brand tb3 on tb2.brand_id=tb3.brand_id`)
	utils.MysqlStringInUtils(&buf, ids, " where tb.sku_id")
	result, err := db.GetGoodDb().Query(buf.String())
	if err == nil {
		var skuId, price,businessId, spuID int64
		list := make([]Sku, 0)
		for result.Next() {
			err = result.Scan(&skuId, &price, &businessId, &spuID)
			if err == nil {
				list = append(list, Sku{
					SpuID:          spuID,
					BusinessId:     businessId,
					Price:          price,
					SkuID:          skuId,
					Num:            idMap[skuId],
				})
			}
		}

		buf1.WriteString(`select tb.template_default_num,tb.template_default_price,tb.template_increase_num,tb.template_increase_price,t.spu_id,tb4.type from qz_good_sku t
INNER JOIN qz_good_delivery tb4 on tb4.spu_id=t.spu_id 
inner JOIN qz_good_delivery_template tb on tb.id=tb4.template_id `)
		utils.MysqlStringInUtils(&buf1, ids, " where t.sku_id")
		buf1.WriteString(" GROUP BY t.spu_id")
		deliveryResult, err := db.GetGoodDb().Query(buf1.String())
		if err != nil {
			return nil, err
		}
		type Deliver struct {
			DefaultNum    int64
			DefaultPrice  int64
			IncreaseNum   int64
			IncreasePrice int64
		}
		var dNum, dPrice, iNum, iPrice, dType, spuId int64
		deliverMap := make(map[int64]map[int64]Deliver)
		for deliveryResult.Next() {
			err = deliveryResult.Scan(&dNum, &dPrice, &iNum, &iPrice, &spuId, &dType)
			if err == nil {
				temp, ok := deliverMap[spuId]
				if !ok {
					temp = make(map[int64]Deliver)
				}
				temp[dType] = Deliver{
					IncreasePrice: iPrice,
					DefaultPrice:  dPrice,
					IncreaseNum:   iNum,
					DefaultNum:    dNum,
				}
				deliverMap[spuId] = temp
			}
		}
		shopDeliverMap := make(map[int64]int64)
		spuNum := make(map[int64]int64)
		for _, info := range list {
			//统计spu购买的数量
			spuNum[info.SpuID] += idMap[info.SkuID]
		}
		for k, v := range spuNum {
			temp, ok := deliverMap[k]
			if ok {
				//spu存在运费
				for t, info := range temp {
					tempDeliveryPrice := info.DefaultPrice
					if v > info.DefaultNum {
						difNum := (v - info.DefaultNum) / info.IncreaseNum
						tempDeliveryPrice += info.IncreasePrice * difNum
					}
					if t == enum.DeliveryTypeShop {
						shopDeliverMap[k] = tempDeliveryPrice
					}
				}
			}
		}

		return &OrderSku{
			List:               list,
			ShopDeliverPrice:   shopDeliverMap,
		}, err
	} else {
		return nil, err
	}
}

func GetSkuIdFromQuery(query string,brandIdList []int64 ) ([]int64, error) {
	var buf strings.Builder
	args:=make([]interface{},0)
	buf.WriteString(`select tb.sku_id from tb_good_sku tb inner join qz_good_spu tb1 on tb1.spu_id=tb.spu_id`)
	if len(brandIdList)>0{
		buf.WriteString(` inner join qz_good_brand_record tb2 on tb2.spu_id=tb.spu_id `)
		buf.WriteString(` inner join qz_good_brand tb3 on tb3.brand_id=tb2.brand_id `)
	}
	buf.WriteString(" where 1=1 ")
	if query!=""{
		buf.WriteString(`tb1.spu_name like (?) or tb.sku_name like (?)`)
		args=append(args,"%"+query+"%","%"+query+"%")
	}
	if len(brandIdList)>0{
		utils.MysqlStringInUtils(&buf,brandIdList," AND tb3.business_id")
	}
	spuResult, err := db.GetGoodDb().Query(buf.String(),args...)
	if err == nil {
		var spuId int64
		list := make([]int64, 0)
		for spuResult.Next() {
			err = spuResult.Scan(&spuId)
			if err == nil {
				list = append(list, spuId)
			} else {
				xlog.ErrorP(err)
			}
		}
		return list, nil
	} else {
		return nil, err
	}
}
