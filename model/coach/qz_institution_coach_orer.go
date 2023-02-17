package coachModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"strings"
)

type CoachOrder struct {
	Id      int64
	Num     int64
	WorkId  int64
	OrderId int64
	Status  int64
	Uid     int64
}

func EditInstitutionCoachOrder(workId, uid, status,orderId int64) error {
	return db.GetCoachDb().WithTransaction(func(tx *db.Tx) error {
		var iid sql.NullInt64
		err:=tx.QueryRow("select id from qz_institution_coach_order where uid=? and order_id=?",uid,orderId).Scan(&iid)
		if err!=nil&&err!=sql.ErrNoRows{
			return xlog.Error(err)
		}else if err==sql.ErrNoRows{
			_,err=tx.Exec("insert into qz_institution_coach_order (work_id,uid,`status`,order_id) values (?,?,?,?)",
				workId, uid, status,orderId)
			if err!=nil{
				return xlog.Error(err)
			}
		}else{
			_,err=tx.Exec("update qz_institution_coach_order set `status`=? where id=?",status,iid.Int64)
			if err!=nil{
				return xlog.Error(err)
			}
		}
		return nil
	})
}

func GetInstitutionCoachOrderNum(idList []int64) ([]*CoachOrder, error) {
	list := make([]*CoachOrder, 0)
	var buf strings.Builder
	args := make([]interface{}, 0)
	buf.WriteString("select sum(num),work_id from qz_institution_coach_order where 1=1 and `status`<>?")
	args=append(args,enum.OrderStatusRefund)
	utils.MysqlStringInUtils(&buf, idList, " and work_id")
	buf.WriteString(" group by work_id")
	rows, err := db.GetCoachDb().Query(buf.String(), args...)
	if err == nil {
		var num, id sql.NullInt64
		for rows.Next() {
			err = rows.Scan(&num, &id)
			if err == nil {
				temp := &CoachOrder{
					Id:  id.Int64,
					Num: num.Int64,
				}
				list = append(list, temp)
			}
		}
	}
	return list, nil
}

func GetInstitutionCoachOrder(nowUid, page, size,nowOrderId int64, selectStatus []int64) ([]*CoachOrder, error) {
	list := make([]*CoachOrder, 0)
	var buf strings.Builder
	args := make([]interface{}, 0)
	buf.WriteString("select id,uid,work_id,order_id,`status`,`num` from qz_institution_coach_order where 1=1 ")
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 20
	}
	if nowUid!=0{
		buf.WriteString(" and uid=?")
		args = append(args, nowUid)
	}

	if nowOrderId!=0{
		buf.WriteString(" and order_id=?")
		args = append(args, nowOrderId)
	}

	utils.MysqlStringInUtils(&buf, selectStatus, " and `status`")
	buf.WriteString(" order by id desc")
	rows, err := db.GetCoachDb().Query(buf.String(), args...)
	if err == nil {
		var id, uid, workId, orderId, status, num sql.NullInt64
		for rows.Next() {
			err = rows.Scan(&id,&uid, &workId, &orderId, &status, &num)
			if err == nil {
				temp := &CoachOrder{
					Uid:     uid.Int64,
					Id:      id.Int64,
					WorkId:  workId.Int64,
					OrderId: orderId.Int64,
					Status:  status.Int64,
					Num:     num.Int64,
				}
				list = append(list, temp)
			}
		}
	}
	return list, nil
}
