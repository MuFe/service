package usermodel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/xlog"
)

// OpenID open_id
type OpenID struct {
	ID         int64
	OpenID     string
	SessionKey string
	AppID      string
	UID        int64
}

// GetOpenIDByOpenIDAndAppID 获取用户的openID
func GetOpenIDByOpenIDAndAppID(oID, aID string) (OpenID, error) {
	var result OpenID
	var id, uID sql.NullInt64
	var openID, sessionKey, appID sql.NullString
	err := db.GetUserDb().QueryRow(
		`select o.id,o.open_id,o.session_key,o.app_id,o.uid from qz_open_id o where o.open_id = ? and o.app_id = ?`,
		oID, aID).Scan(&id, &openID, &sessionKey, &appID, &uID)
	if err != nil && err != sql.ErrNoRows {
		return result, xlog.Error(err)
	}
	result.ID = id.Int64
	result.OpenID = openID.String
	result.SessionKey = sessionKey.String
	result.AppID = appID.String
	result.UID = uID.Int64
	return result, nil
}

//创建openId
func CreateOpenId(uid int64,appId,openId string,tx *db.Tx)error{
	var err error
	if tx==nil{
		err=db.GetUserDb().WithTransaction(func(tx *db.Tx) error {
			_,err:=tx.Exec("insert into qz_open_id(uid,app_id,open_id) values(?,?,?)", uid, appId, openId)
			return err
		})
	} else {
		_, err =tx.Exec("insert into qz_open_id(uid,app_id,open_id) values(?,?,?)", uid, appId, openId)
	}
	if err != nil {
		return xlog.Error("生成第三方数据出错")
	}
	return nil
}

func UnbindOpenId(uid int64,tx *db.Tx)error{
	var err error
	if tx==nil{
		err=db.GetUserDb().WithTransaction(func(tx *db.Tx) error {
			_,err:=tx.Exec("delete from qz_open_id where uid=?", uid)
			return err
		})
	} else {
		_,err=tx.Exec("delete from qz_open_id where uid=?", uid)
	}
	if err != nil {
		return xlog.Error("解绑出错")
	}
	return nil
}
