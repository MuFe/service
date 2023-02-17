package schoolModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"strings"
)

type Grade struct {
	Id int64
	Name string
}

type GradeType struct {
	Id int64
	Name string
	List []Grade
}

func GradeInfo(schoolId int64)([]*GradeType,error){
	var buf strings.Builder
	args:=make([]interface{},0)
	buf.WriteString(`select tb.id,tb.name,tb.school_type,tb1.name from qz_grade tb 
inner join qz_school_type tb1 on tb1.id=tb.school_type `)
	if schoolId!=0{
		buf.WriteString(" inner join qz_school tb2 on tb2.school_type=tb1.id where tb2.id=?")
		args=append(args,schoolId)
	}
	buf.WriteString(" order by tb.id asc ")
	result,err:=db.GetSchool().Query(buf.String(),args...)
	if err!=nil{
		return nil,err
	}
	list:=make([]*GradeType,0)
	listMap:=make(map[int64]*GradeType)
	var id,typeId sql.NullInt64
	var name,typeName sql.NullString
	for result.Next(){
		err=result.Scan(&id,&name,&typeId,&typeName)
		if err==nil{
			info,ok:=listMap[typeId.Int64]
			if !ok{
				info=&GradeType{
					Id:typeId.Int64,
					Name:typeName.String,
					List:make([]Grade,0),
				}
				listMap[typeId.Int64]=info
				list=append(list,info)
			}
			info.List=append(info.List,Grade{
				Id:id.Int64,
				Name:name.String,
			})
		}
	}
	return list,nil
}
