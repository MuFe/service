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

type InstitutionSchool struct {
	Icon    string
	Id      int64
	Name    string
	Address string
	Start   int64
	End     int64
	Phone   string
}

func GetInstitutionSchool(status,parentId int64,nowId []int64) ([]InstitutionSchool, error) {
	list := make([]InstitutionSchool, 0)
	var buf strings.Builder
	args := make([]interface{}, 0)
	buf.WriteString("select id,name,address,icon,prefix,phone,start,end from qz_institution_school where 1=1  ")
	if status != enum.StatusAll {
		buf.WriteString(" and status=?")
		args = append(args, status)
	} else {
		buf.WriteString(" and status<>?")
		args = append(args, enum.StatusDelete)
	}


	if len(nowId) != 0 {
		utils.MysqlStringInUtils(&buf,nowId," and id")
	}

	if parentId != 0 {
		buf.WriteString(" and parent_id=?")
		args = append(args, parentId)
	}else {
		buf.WriteString(" and parent_id<>0")
	}
	buf.WriteString(" order by id desc")
	rows, err := db.GetCoachDb().Query(buf.String(), args...)
	if err == nil {
		var id,start,end sql.NullInt64
		var name, icon, prefix,address,phone sql.NullString
		for rows.Next() {
			err = rows.Scan(&id, &name,&address, &icon, &prefix,&phone,&start,&end)
			if err == nil {
				list = append(list, InstitutionSchool{
					Id:id.Int64,
					Name:name.String,
					Address:address.String,
					Icon:prefix.String+icon.String,
					Phone:phone.String,
					Start:start.Int64,
					End:end.Int64,
				})
			}
		}
	}
	return list, nil
}


func EditInstitutionSchool(createBy, id,start,end,parentId int64, name,address,phone,icon,prefix string,del bool,courseId []int64) (int64,error) {
	err:=db.GetCoachDb().WithTransaction(func(tx *db.Tx) error {
		now:=time.Now().Unix()
		if id==0{
			idResult,err:=tx.Exec("insert into qz_institution_school (name,address,create_time,modify_time,create_by,phone,start,end,parent_id) values(?,?,?,?,?,?,?,?,?)",name,address,now,now,createBy,phone,start,end,parentId)
			if err!=nil{
				return xlog.Error("创建机构校区失败")
			}
			id,_=idResult.LastInsertId()
			err=EditInstitutionSchoolCourse(tx,id,courseId)
			return err
		}else if del{
			_,err:=tx.Exec("update qz_institution_school set `status`=?,modify_time=? where id=?",enum.StatusDelete,now,id)
			if err!=nil{
				return xlog.Error("修改机构校区失败")
			}
		}else if icon==""{
			_,err:=tx.Exec("update qz_institution_school set name=?,address=?,phone=?,start=?,end=?,modify_time=? where id=?",name,address,phone,start,end,now,id)
			if err!=nil{
				return xlog.Error("修改机构校区失败")
			}
			err=EditInstitutionSchoolCourse(tx,id,courseId)
			return err
		}else{
			_,err:=tx.Exec("update qz_institution_school set icon=?,prefix=?,modify_time=? where id=?",icon,prefix,now,id)
			if err!=nil{
				return xlog.Error("修改机构校区失败")
			}
		}
		return nil
	})
	return id,err
}

