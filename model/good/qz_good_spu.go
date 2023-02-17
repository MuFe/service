package goodmodel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"strconv"
	"strings"
	"time"
)

// Spu 商品
type Spu struct {
	SpuID     int64
	SpuStatus int64
	SpuName   string

	// 表外数据
	SkuID        int64
	BusinessID   int64
	BusinessName string
	Price        int64
	MemberPrice  int64
	Photo        string
	CategoryID   int64
	CategoryName string
	Status       int64
	SaleStatus   bool
	CreateTime   int64
	ModifyTime   int64
}

type GoodDetail struct {
	Id               int64
	Number           string
	DeliveryInfo     []GoodDeliveryInfo
	Location         string
	Name             string
	List             []Sku
	Detail           string
	CategoryId       int64
	ParentCategoryId int64
	Photos           []SpuPhoto
	SaleTime         int64
	Status           int64
	CommentNum       int64
	StatusMessage    string
	BusinessId       int64
}

type GoodDeliveryInfo struct {
	Type       int64
	TemplateId int64
}

const insertSpuSql = `
INSERT INTO qz_good_spu (
spu_name,
spu_number,
STATUS,
spu_delivery_location,
sale_time,
create_time,
modify_time,
spu_detail
)
VALUES
	(?,?,?,?,?,?,?,'')
`

func AddGood(tx *db.Tx,name,spuNo string,status int64,  spuDeliveryLocation string,)(int64,error){
	nowTime := time.Now().Unix()
	spuResult, err := tx.Exec(
		insertSpuSql,
		name,
		spuNo,
		status,
		spuDeliveryLocation,
		nowTime,
		nowTime,
		nowTime,
	)
	if err != nil {
		return 0,xlog.Error(err)
	}
	spuId, err := spuResult.LastInsertId()
	if err != nil {
		return 0,xlog.Error(err)
	}
	return spuId,nil
}


//获取详情
func GetDetail(spuId, skuStatus int64) (*GoodDetail, error) {
	sqlStr := `SELECT
		tb.spu_number,
		tb.spu_name,
		tb.spu_detail,
		tb.spu_delivery_location,
		tb.sale_time,
		tb.STATUS,
		tb.status_message,
		tb2.category_id,
		tb4.category_id,
		tb6.business_id,
		GROUP_CONCAT( tb1.type ) AS type,
		GROUP_CONCAT( tb1.template_id ) AS template_id
	FROM
		qz_good_spu tb
		LEFT JOIN qz_good_delivery tb1 ON tb1.spu_id = tb.spu_id
		LEFT JOIN qz_good_category_record tb2 ON tb2.spu_id = tb.spu_id 
		LEFT JOIN qz_good_category tb3 ON tb3.category_id = tb2.category_id 
		LEFT JOIN qz_good_category tb4 ON tb4.category_id = tb3.parent_id 
		LEFT JOIN qz_good_brand_record tb5 on tb5.spu_id=tb.spu_id
		LEFT JOIN qz_good_brand tb6 on tb6.brand_id=tb5.brand_id
	WHERE
		tb.spu_id = ? `
	var spuNumber, spuName, location, detail, deliveryType, deliveryId, StatusMessage sql.NullString
	var saleTime, categoryId, parentId, spuStatus, businessId sql.NullInt64
	err := db.GetGoodDb().QueryRow(sqlStr, spuId).Scan(
		&spuNumber,
		&spuName,
		&detail,
		&location,
		&saleTime,
		&spuStatus,
		&StatusMessage,
		&categoryId,
		&parentId,
		&businessId,
		&deliveryType,
		&deliveryId,
	)
	if err == nil {
		deliveryInfos := make([]GoodDeliveryInfo, 0)
		if deliveryType.String != "" {
			types := strings.Split(deliveryType.String, ",")
			ids := strings.Split(deliveryId.String, ",")
			for k := range types {
				typeInt, _ := strconv.Atoi(types[k])
				idInt, _ := strconv.Atoi(ids[k])
				deliveryInfos = append(deliveryInfos, GoodDeliveryInfo{
					Type:       int64(typeInt),
					TemplateId: int64(idInt),
				})
			}
		}
		spuPhoto, _ := GetSpuPhotoBySpuIDs([]int64{spuId})
		skuList, _ := GetSkuIDBySpuIDs([]int64{spuId}, skuStatus)
		skuIds := make([]int64, 0)
		for _, info := range skuList {
			skuIds = append(skuIds, info.SkuID)
		}

		return &GoodDetail{
			Id:               spuId,
			Number:           spuNumber.String,
			Name:             spuName.String,
			Status:           spuStatus.Int64,
			Location:         location.String,
			Detail:           detail.String,
			CategoryId:       categoryId.Int64,
			ParentCategoryId: parentId.Int64,
			DeliveryInfo:     deliveryInfos,
			SaleTime:         saleTime.Int64,
			Photos:           spuPhoto,
			List:             skuList,
			StatusMessage:    StatusMessage.String,
			BusinessId:       businessId.Int64,
		}, nil
	} else {
		return nil, err
	}
}

func GetSpuListBySpuIDS(spuIDS []int64) ([]Spu, error) {
	var buf strings.Builder
	var list []Spu

	buf.WriteString(`select spu.spu_id,spu.spu_name,spu.status,spu.sale_time,cr.category_id,ca.category_name,spu.create_time,spu.modify_time,tb2.brand_name,tb2.business_id from qz_good_spu spu
	left join qz_good_category_record cr on cr.spu_id = spu.spu_id
	left join qz_good_category ca on ca.category_id = cr.category_id
	left join qz_good_brand_record tb on tb.spu_id=spu.spu_id
	left join qz_good_brand tb2 on tb2.brand_id=tb.brand_id`)
	utils.MysqlStringInUtils(&buf, spuIDS, " where spu.spu_id")
	buf.WriteString(` order by spu.spu_id desc`)
	rows, err := db.GetGoodDb().Query(buf.String())
	if err != nil {
		return nil, xlog.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		var spuID, categoryID, saleTime, status, createTime, modifyTime, businessId sql.NullInt64
		var spuName, brandName, categoryName sql.NullString
		err := rows.Scan(&spuID, &spuName, &status, &saleTime, &categoryID, &categoryName, &createTime, &modifyTime, &brandName, &businessId)
		if err != nil {
			return nil, xlog.Error(err)
		}
		list = append(list, Spu{
			SpuID:        spuID.Int64,
			SpuName:      spuName.String,
			CategoryID:   categoryID.Int64,
			Status:       status.Int64,
			SaleStatus:   saleTime.Int64 > 0,
			CreateTime:   createTime.Int64,
			ModifyTime:   modifyTime.Int64,
			BusinessName: brandName.String,
			BusinessID:   businessId.Int64,
			CategoryName: categoryName.String,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, xlog.Error(err)
	}

	return list, nil
}
