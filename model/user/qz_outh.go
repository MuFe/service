package usermodel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/xlog"
)

// Outh tb_outh
type Outh struct {
	ID     int64
	OuthID string
	Type   int64
	UID    int64
}

// GetOuthByOuthID 获取Outh
func GetOuthByOuthID(t int64, o string) (Outh, error) {
	var result Outh
	var id, tp, uID sql.NullInt64
	var outhID sql.NullString
	err := db.GetUserDb().QueryRow(
		`select o.id,o.type,o.outh_id,o.uid from qz_outh o where o.outh_id = ? and o.type = ?`,
		o, t).Scan(&id, &tp, &outhID, &uID)
	if err != nil && err != sql.ErrNoRows {
		return result, xlog.Error(err)
	}
	result.ID = id.Int64
	result.OuthID = outhID.String
	result.Type = tp.Int64
	result.UID = uID.Int64
	return result, nil
}


//创建第三方登录数据
func CreateOuth(uid, typeInt int64, unionId string, tx *db.Tx) error {
	var err error
	if tx == nil {
		err = db.GetUserDb().WithTransaction(func(tx *db.Tx) error {
			_, err := tx.Exec("insert into qz_outh(uid,type,outh_id) values(?,?,?)", uid, typeInt, unionId)
			return err
		})
	} else {
		_, err = tx.Exec("insert into qz_outh(uid,type,outh_id) values(?,?,?)", uid, typeInt, unionId)
	}

	if err != nil {
		return xlog.Error("生成第三方数据出错")
	}
	return nil
}

func UnbindOuth(uid,typeInt int64,tx *db.Tx)error{
	var err error
	if tx==nil{
		err=db.GetUserDb().WithTransaction(func(tx *db.Tx) error {
			_,err:=tx.Exec("delete from qz_outh where uid=? and type=?", uid,typeInt)
			return err
		})
	} else {
		_,err=tx.Exec("delete from qz_outh where uid=? and type=?", uid,typeInt)
	}
	if err != nil {
		return xlog.Error("解绑出错")
	}
	return nil
}
