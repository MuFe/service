package coachModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"strings"
	"time"
)

type Institution struct {
	Icon    string
	Id      int64
	Name    string
	Address string
	Code    string
}

func GetInstitution(status, nowId int64,selectCode string) ([]Institution, error) {
	list := make([]Institution, 0)
	var buf strings.Builder
	args := make([]interface{}, 0)
	buf.WriteString("select id,name,address,icon,prefix,code,code_prefix from qz_institution where 1=1 ")
	if status != enum.StatusAll {
		buf.WriteString(" and status=?")
		args = append(args, status)
	} else {
		buf.WriteString(" and status<>?")
		args = append(args, enum.StatusDelete)
	}

	if selectCode!=""{
		buf.WriteString(" and CONCAT(code_prefix,`code`)=?")
		args = append(args, selectCode)
	}
	if nowId != 0 {
		buf.WriteString(" and id=?")
		args = append(args, nowId)
	}
	buf.WriteString(" order by id desc")
	rows, err := db.GetCoachDb().Query(buf.String(), args...)
	if err == nil {
		var id sql.NullInt64
		var name, icon, prefix,address,code,codePrefix sql.NullString
		for rows.Next() {
			err = rows.Scan(&id, &name,&address, &icon, &prefix,&code,&codePrefix)
			if err == nil {
				list = append(list, Institution{
					Id:id.Int64,
					Name:name.String,
					Address:address.String,
					Icon:prefix.String+icon.String,
					Code:codePrefix.String+code.String,
				})
			}
		}
	}
	return list, nil
}


func EditInstitution(createBy, id int64, name,address,icon,prefix string,del bool) (int64,error) {
	err:=db.GetCoachDb().WithTransaction(func(tx *db.Tx) error {
		now:=time.Now().Unix()
		if id==0{
			inviteResult,err:=tx.Query("select code from qz_institution ")
			if err!=nil{
				return xlog.Error("创建机构失败")
			}
			inviteMap:=make(map[string]bool)
			var invite sql.NullString
			for inviteResult.Next(){
				err=inviteResult.Scan(&invite)
				if err==nil{
					inviteMap[invite.String]=true
				}
			}
			code:=utils.BaseCode6(inviteMap)
			idResult,err:=tx.Exec("insert into qz_institution (name,address,create_time,modify_time,create_by,code,code_prefix) values(?,?,?,?,?,?,?)",name,address,now,now,createBy,code,enum.InstitutionCodePrefix)
			if err!=nil{
				return xlog.Error("创建机构失败")
			}
			id,_=idResult.LastInsertId()
		}else if del{
			_,err:=tx.Exec("update qz_institution set `status`=?,modify_time=? where id=?",enum.StatusDelete,now,id)
			if err!=nil{
				return xlog.Error("修改机构失败")
			}
			err=DelInstitutionCoach(id,tx)
			if err!=nil{
				return xlog.Error("修改机构失败")
			}
		}else if icon==""{
			_,err:=tx.Exec("update qz_institution set name=?,address=?,modify_time=? where id=?",name,address,now,id)
			if err!=nil{
				return xlog.Error("修改机构失败")
			}
		}else{
			_,err:=tx.Exec("update qz_institution set icon=?,prefix=?,modify_time=? where id=?",icon,prefix,now,id)
			if err!=nil{
				return xlog.Error("修改机构失败")
			}
		}
		return nil
	})
	return id,err
}


