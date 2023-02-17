package webModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/utils"
	"strings"
)

type WebData struct {
	Id int64
	Type int64
	Content string
}
func GetWebValue(typeIntList []int64)([]WebData,error){
	var content sql.NullString
	var typeInt sql.NullInt64
	list:=make([]WebData,0)
	var buf strings.Builder
	buf.WriteString("select `value`,`type` from qz_web where 1=1")
	utils.MysqlStringInUtils(&buf,typeIntList," and type")
	result,err:=db.GetAdminDb().Query(buf.String())
	if err!=nil{
		return nil,err
	}
	for result.Next(){
		err=result.Scan(&content,&typeInt)
		if err==nil{
			list=append(list,WebData{
				Type:typeInt.Int64,
				Content:content.String,
			})
		}
	}
	return list,nil
}

func AddWebContactUs(name,phone,email,content string)error{
	_,err:=db.GetAdminDb().Exec("insert into qz_contact_us (phone,email,name,content) values (?,?,?,?)",phone,email,name,content)
	return err
}

func EditWebContactUsStatus(id,status int64)error{
	_,err:=db.GetAdminDb().Exec("update qz_contact_us set `status`=? where id=?",status,id)
	return err
}
