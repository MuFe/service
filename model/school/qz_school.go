package schoolModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"strings"
)

type School struct {
	Id int64
	Name string
	Desc string
	Icon string
	Code string
	Type string
	Address string
	TypeId int64
}

type SchoolCode struct {
	Id int64
	Type int64
	Value int64
	School int64
	Status int64
	Code string
}

type SchoolType struct {
	Id int64
	Name string
	List []Grade
}

func AddSchool(uid,id int64)error{
	if id==0||uid==0{
		return xlog.Error("参数有误")
	}
	return db.GetSchool().WithTransaction(func(tx *db.Tx) error {
		var status,count sql.NullInt64
		err:=tx.QueryRow("select `status` from qz_school where id=?",id).Scan(&status)
		if err==sql.ErrNoRows{
			return errcode.HttpErrorWringParam.RPCError()
		}else if err!=nil{
			return err
		}else if status.Int64!=enum.StatusNormal{
			return errcode.HttpErrorWringParam.RPCError()
		}
		err=tx.QueryRow("select count(id) from qz_school_record where uid=? and school_id=?",uid,id).Scan(&count)
		if err==nil{
			if count.Int64>0{
				return xlog.Error("您已经加入该学校")
			}
		}	else if err!=sql.ErrNoRows{
			return err
		}
		_,err=tx.Exec("insert into qz_school_record (uid,school_id) values (?,?)",uid,id)
		return nil
	})
}

func QuitSchool(uid,id int64)error{
	if id==0||uid==0{
		return xlog.Error("参数有误")
	}
	return db.GetSchool().WithTransaction(func(tx *db.Tx) error {
		var status sql.NullInt64
		err:=tx.QueryRow("select `status` from qz_school where id=?",id).Scan(&status)
		if err==sql.ErrNoRows{
			return errcode.HttpErrorWringParam.RPCError()
		}else if err!=nil{
			return err
		}else if status.Int64!=enum.StatusNormal{
			return errcode.HttpErrorWringParam.RPCError()
		}
		_,err=tx.Exec("delete from qz_school_record where uid=? and school_id=?",uid,id)
		if err != nil {
			return xlog.Error(err)
		}
		err=DeleteClassRecord(0,uid,tx)
		return err
	})
}

func GetSchool(uid,id,status int64)[]School{
	var buf strings.Builder
	args:=make([]interface{},0)
	buf.WriteString(`select tb.id,tb.name,tb.desc,tb.icon,tb.code,tb.code_prefix,tb.school_type,tb2.name,tb.address from qz_school tb left join qz_school_record tb1 on tb1.school_id=tb.id left join qz_school_type tb2 on tb2.id=tb.school_type where 1=1`)
	if uid!=0{
		buf.WriteString(" and tb1.uid=?")
		args=append(args,uid)
	}
	if id!=0{
		buf.WriteString(" and tb.id=?")
		args=append(args,id)
	}
	if status!=0{
		buf.WriteString(" and tb.status=?")
		args=append(args,status)
	}
	buf.WriteString(" GROUP BY tb.id  order by tb.id desc")
	result,err:=db.GetSchool().Query(buf.String(),args...)
	list:=make([]School,0)
	if err==nil{
		var sid,schoolTypeId sql.NullInt64
		var name ,desc,icon,code,prefix,schoolType,address sql.NullString
		for result.Next(){
			err=result.Scan(&sid,&name,&desc,&icon,&code,&prefix,&schoolTypeId,&schoolType,&address)
			if err==nil{
				list=append(list,School{
					Id:sid.Int64,
					Name:name.String,
					Desc:desc.String,
					Icon:icon.String,
					Code:prefix.String+code.String,
					TypeId:schoolTypeId.Int64,
					Type:schoolType.String,
					Address:address.String,
				})
			}
		}
	}
	return list
}

func EditSchool(name,address,icon string,id,typeId int64)(int64,error){
	var err error
	var schoolId int64
	if id==0{
		err=db.GetSchool().WithTransaction(func(tx *db.Tx) error {
			inviteResult,err:=tx.Query("select code from qz_school ")
			if err!=nil{
				return xlog.Error("创建学校失败")
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
			idResult,err:=tx.Exec("insert into qz_school (`name`,`address`,`code`,code_prefix,school_type) values(?,?,?,?,?)",name,address,code,enum.SchoolCodePrefix,typeId)
			if err==nil{
				schoolId,_=idResult.LastInsertId()
			}
			return err
		})
	} else {
		if icon!=""{
			_,err=db.GetSchool().Exec("update qz_school set `icon`=? where id=?",icon,id)
		} else {
			_,err=db.GetSchool().Exec("update qz_school set `name`=?,`address`=?,school_type=? where id=?",name,address,typeId,id)
		}
		schoolId=id
	}

	return schoolId,err
}

func FindSchool(content string,status int64)[]School{
	var buf strings.Builder
	args:=make([]interface{},0)
	buf.WriteString("select tb.id,tb.name,tb.desc,tb.icon,tb.code,tb.code_prefix,tb.address from qz_school tb where CONCAT(code_prefix,`code`)=? and status=?")
	args=append(args,content,status)
	result,err:=db.GetSchool().Query(buf.String(),args...)
	list:=make([]School,0)
	if err==nil{
		var sid sql.NullInt64
		var name ,desc,icon,code,prefix,address sql.NullString
		for result.Next(){
			err=result.Scan(&sid,&name,&desc,&icon,&code,&prefix,&address)
			if err==nil{
				list=append(list,School{
					Id:sid.Int64,
					Name:name.String,
					Desc:desc.String,
					Icon:icon.String,
					Code:prefix.String+code.String,
					Address:address.String,
				})
			}
		}
	}
	return list
}

func GetSchoolType()[]*SchoolType{
	result,err:=db.GetSchool().Query("select tb.id,tb.name,tb1.id,tb1.name from qz_school_type tb inner join qz_grade tb1 on tb1.school_type=tb.id")
	list:=make([]*SchoolType,0)
	if err==nil{
		var id,gradeId sql.NullInt64
		var name,gradeName sql.NullString
		resultMap:=make(map[int64]*SchoolType)
		for result.Next(){
			err=result.Scan(&id,&name,&gradeId,&gradeName)
			if err==nil{
				temp,ok:=resultMap[id.Int64]
				if !ok{
					temp=&SchoolType{
						Id:id.Int64,
						Name:name.String,
					}
					list=append(list,temp)
					resultMap[id.Int64]=temp
				}
				temp.List=append(temp.List,Grade{
					Id:gradeId.Int64,
					Name:gradeName.String,
				})
			}
		}
	}
	return list
}




