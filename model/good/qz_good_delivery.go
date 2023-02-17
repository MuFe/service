package goodmodel

import (
	"mufe_service/jsonRpc"
	"strings"
	"time"
	"mufe_service/camp/db"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
)

// SpuDelivery 配送方式
type SpuDelivery struct {
	SpuID  int64
	Name   string
}

type DeliveryData struct {
	Id            int64
	DefaultNum    int64
	DefaultPrice  int64
	IncreaseNum   int64
	IncreasePrice int64

	Name              string
}


//获取配送信息
func GetDeliveryTemplate(businessId int64, spuId int64, idsList []int64) ([]DeliveryData, error) {
	var name string
	var id, defaultNum, increaseNum,defaultPrice, increasePrice int64
	var buf strings.Builder

	args := make([]interface{}, 0)
	buf.WriteString(`SELECT
	tb.id,
	tb.template_name,
	tb.template_default_num,
	tb.template_default_price,
	tb.template_increase_num,
	tb.template_increase_price
FROM
	qz_good_delivery_template  tb`)
	if spuId != 0 {
		buf.WriteString(` inner join qz_good_delivery tb1 on tb1.template_id=tb.id`)
	}
	buf.WriteString(` WHERE tb.status = 1`)
	if businessId != 0 {
		buf.WriteString(" AND tb.business_id= ?")
		args = append(args, businessId)
	}
	if len(idsList) > 0 {
		utils.MysqlStringInUtils(&buf, idsList, " AND tb.id ")
	}
	if spuId != 0 {
		buf.WriteString(" AND tb1.spu_id= ?")
		args = append(args, spuId)
	}
	result, err := db.GetGoodDb().Query(buf.String(), args...)
	if err != nil {
		xlog.ErrorP(err)
		return nil, err
	}
	list := make([]DeliveryData, 0)
	for result.Next() {
		err = result.Scan(
			&id,
			&name,
			&defaultNum,
			&defaultPrice,
			&increaseNum,
			&increasePrice,
		)
		if err == nil {
			list = append(list, DeliveryData{
				Id:                id,
				Name:              name,
				DefaultNum:        defaultNum,
				DefaultPrice:      defaultPrice,
				IncreaseNum:       increaseNum,
				IncreasePrice:     increasePrice,
			})
		}
	}
	return list, nil
}

//修改配送信息
func EditDeliveryTemplate(id, dNum, dPrice, iNum, iPrice, businessId, uid int64, name string) (int64, error) {
	var resultId int64
	if id != 0 {
		resultId = id
		_, err := db.GetGoodDb().Exec(
			`UPDATE qz_good_delivery_template 
SET template_name =?,
template_default_num =?,
template_default_price =?,
template_increase_num =?,
template_increase_price =?
WHERE
	id =?`,
			name,
			dNum,
			dPrice,
			iNum,
			iPrice,
			id,
		)
		if err != nil {
			return 0, err
		}
	} else {
		result, err := db.GetGoodDb().Exec(
			`INSERT INTO qz_good_delivery_template (
template_name,
template_default_num,
template_default_price,
template_increase_num,
template_increase_price,
business_id,
create_time,
create_by
)
VALUES
	( ?,?,?,?,?,?,?,?)`,
			name,
			dNum,
			dPrice,
			iNum,
			iPrice,
			businessId,
			time.Now().Unix(),
			uid,
		)
		if err != nil {
			return 0, err
		}
		resultId, err = result.LastInsertId()
		if err != nil {
			return 0, err
		}
	}
	return resultId, nil
}

//删除配送信息
func DelDeliveryTemplate(id int64, businessId int64, uid int64) error {
	_, err := db.GetGoodDb().Exec("update qz_good_delivery_template set status =0  where id=?", id)
	if err != nil {
		return err
	}
	return nil
}

func AddDelivery(tx *db.Tx, deliverInfo []*app.DeliveryData,spuId,uid int64)error{
	if len(deliverInfo) > 0 {
		sqlStr := "insert into qz_good_delivery (type,spu_id,template_id,create_time,create_by) values"
		placeHolder := "(?,?,?,?,?)"

		nowTime := time.Now().Unix()
		var values []string
		var args []interface{}
		for _, v := range deliverInfo {
			values = append(values, placeHolder)
			args = append(args, v.Type)
			args = append(args, spuId)
			args = append(args, v.Id)
			args = append(args, nowTime)
			args = append(args, uid)
		}
		sqlStr += strings.Join(values, ",")
		_, err := tx.Exec(sqlStr, args...)
		if err != nil {
			return xlog.Error(err)
		}
	}
	return nil
}
