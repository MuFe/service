package collectionModel

import (
	"database/sql"
	"mufe_service/camp/db"
)




func IsCollection(cId,typeInt,uid int64) bool {
	var id sql.NullInt64
	err:=db.GetCourse().QueryRow("select id from qz_collection where content_id=? and type=? and uid=?",cId,typeInt,uid).Scan(&id)
	if err==nil{
		return true
	}
	return false
}

func EditCollection(cId,typeInt,uid int64,del bool)error{
	var err error
	if del{
		_,err=db.GetCourse().Exec("delete from qz_collection where content_id=? and type=? and uid=?",cId,typeInt,uid)
	}else {
		_,err=db.GetCourse().Exec("insert into qz_collection (content_id,type,uid) values (?,?,?)",cId,typeInt,uid)
	}
	return err
}

func GetCollection(uid,typeInt int64)([]int64,error){
	result,err:=db.GetCourse().Query("select content_id from qz_collection where type=? and uid=?",typeInt,uid)
	idList:=make([]int64,0)
	if err==nil{
		var id sql.NullInt64
		for result.Next() {
			err := result.Scan(&id)
			if err != nil {
				return idList,err
			}
			idList = append(idList, id.Int64)
		}
		return idList,nil
	} else {
		return idList,err
	}
}
