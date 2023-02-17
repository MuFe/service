package courseModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/xlog"
	"strings"
)

type Level struct {
	ID       int64
	Name    string
}


func EditLevel(name string,id,typeInt int64) (int64,error) {
	if id!=0{
		_,err:=db.GetCourse().Exec("update  qz_level set name=? where id=?",name,id)
		if err!=nil{
			return 0,xlog.Error(err)
		}
		return id,nil
	}else{
		result,err:=db.GetCourse().Exec("insert into qz_level (name,`type`) values (?,?)",name,typeInt)
		if err!=nil{
			return 0,xlog.Error(err)
		}
		lastId,err:=result.LastInsertId()
		if err!=nil{
			return 0,xlog.Error(err)
		}
		return lastId,nil
	}
}

func GetLevel(typeInt int64)([]Level,error) {
	buf := strings.Builder{}
	args := make([]interface{}, 0)
	buf.WriteString("select id,name from qz_level where status=1")
	if typeInt!=0{
		buf.WriteString(" and `type`=?")
		args=append(args,typeInt)
	}
	result, err := db.GetCourse().Query(buf.String(), args...)
	if err != nil {
		return nil, err
	}
	list := make([]Level, 0)
	var id sql.NullInt64
	var name sql.NullString
	for result.Next() {
		err := result.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		list = append(list, Level{
			ID:       id.Int64,
			Name:    name.String,
		})
	}
	return list, nil
}


