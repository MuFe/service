package newModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"strings"
	"time"
)

type New struct {
	Id int64
	Time int64
	Type int64
	Source string
	Cover string
	Title string
	Content string
}
func GetNews(newId,status,typeInt int64)([]New,error){
	var content,title,source,cover,prefix sql.NullString
	var id,time,newType sql.NullInt64
	list:=make([]New,0)
	var buf strings.Builder
	args:=make([]interface{},0)
	buf.WriteString("select id,cover,`prefix`,title,create_time,content,source,`type` from qz_news where 1=1 ")
	if status==0{
		buf.WriteString(" and `status`<>?")
		args=append(args,enum.StatusDelete)
	}else{
		buf.WriteString(" and `status`=?")
		args=append(args,enum.StatusNormal)
	}
	if typeInt!=0{
		buf.WriteString(" and `type`=?")
		args=append(args,enum.NewsAppHome)
	}
	if newId!=0{
		buf.WriteString(" and id=?")
		args=append(args,newId)
	}
	buf.WriteString(" order by `type` desc,id desc")
	result,err:=db.GetSchool().Query(buf.String(),args...)
	if err!=nil{
		return nil,err
	}
	for result.Next(){
		err=result.Scan(&id,&cover,&prefix,&title,&time,&content,&source,&newType)
		if err==nil{
			list=append(list,New{
				Id:id.Int64,
				Title:title.String,
				Time:time.Int64,
				Source:source.String,
				Content:content.String,
				Cover:prefix.String+cover.String,
				Type:newType.Int64,
			})
		}
	}
	return list,nil
}

func EditNewType(id,typeInt int64)error{
	_,err:=db.GetSchool().Exec("update qz_news set `type`=? where id=?",typeInt,id)
	return err
}

func EditNewContent(id int64,content string)error{
	_,err:=db.GetSchool().Exec("update qz_news set `content`=? where id=?",content,id)
	return err
}

func DelNew(id,status int64)error{
	_,err:=db.GetSchool().Exec("update qz_news set `status`=? where id=?",status,id)
	return err
}

func EditNewCover(id int64,cover,prefix string)error{
	_,err:=db.GetSchool().Exec("update qz_news set `cover`=?,prefix=? where id=?",cover,prefix,id)
	return err
}

func EditNew(id int64,title,source string)(int64,error){
	var err error
	if id==0{
		result,err:=db.GetSchool().Exec("insert into qz_news (title,source,content,create_time) values(?,?,'',?)",title,source,time.Now().Unix())
		if err==nil{
			id,err=result.LastInsertId()
		}
	}else{
		_,err=db.GetSchool().Exec("update qz_news set `title`=?,source=? where id=?",title,source,id)
	}

	return id,err
}
