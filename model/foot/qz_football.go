package footModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/errcode"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"strings"
)

type Football struct {
	UID      int64
	Mac     string
}

func Bind(uid int64,mac string) error {
	_,err:=db.GetCourse().Exec("insert into qz_football (uid,mac) values(?,?)",uid,mac)
	return err
}
func UnBind(uid int64,mac string) error {
	var err error
	if mac!=""{
		_,err=db.GetCourse().Exec("delete from qz_football where uid=? and mac=?",uid,mac)
	} else{
		_,err=db.GetCourse().Exec("delete from qz_football where uid=? ",uid)
	}
	return err
}

func GetSchool(uid,typeInt int64) (int64,error) {
	var schoolId sql.NullInt64
	err:=db.GetCourse().QueryRow("select school_id from qz_area_use where uid=? and type=?",uid,typeInt).Scan(&schoolId)
	if err!=nil&&err!=sql.ErrNoRows{
		return 0,xlog.Error(err)
	}else if err!=nil{
		return 0,xlog.Error(errcode.HttpErrorWringParam.Msg)
	}
	return schoolId.Int64,err
}

func GetFoot(uidList []int64) []Football {
	list:=make([]Football,0)
	var build strings.Builder
	build.WriteString("select mac,uid from qz_football  ")
	utils.MysqlStringInUtils(&build,uidList," where uid")
	result,err:=db.GetCourse().Query(build.String())
	if err==nil{
		var mac sql.NullString
		var uid sql.NullInt64
		for result.Next(){
			err=result.Scan(&mac,&uid)
			if err==nil{
				list=append(list,Football{
					Mac:mac.String,
					UID:uid.Int64,
				})
			}
		}
	}
	return list
}


