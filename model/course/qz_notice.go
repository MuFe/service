package courseModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/xlog"
	"time"
)

type Notice struct {
	Id int64
	Title string
	Content string
	Time int64
	CreateBy int64
}

func GetNotice(courseId int64)[]Notice{
	list:=make([]Notice,0)
	result,err:=db.GetCourse().Query("select id,title,create_time from qz_notice where course_id=?",courseId)
	if err==nil{
		var id,cTime sql.NullInt64
		var title sql.NullString
		for result.Next(){
			err=result.Scan(&id,&title,&cTime)
			if err==nil{
				list=append(list,Notice{
					Id:id.Int64,
					Time:cTime.Int64,
					Title:title.String,
				})
			}
		}
	}
	return list
}

func AddNotice(courseId,createBy int64,title,content string)error{
	_,err:=db.GetCourse().Exec("insert qz_notice (course_id,create_by,create_time,title,content) value(?,?,?,?,?)",courseId,createBy,time.Now().Unix(),title,content)
	return err
}

func NoticeDeatil(id int64)(*Notice,error){
	var cTime,createBy sql.NullInt64
	var title,content sql.NullString
	err:=db.GetCourse().QueryRow("select title,create_time,content,create_by from qz_notice where id=?",id).Scan(&title,&cTime,&content,&createBy)
	if err==nil{
		return &Notice{
			Id:id,
			Title:title.String,
			Time:cTime.Int64,
			Content:content.String,
			CreateBy:createBy.Int64,
		},nil
	}else if err!=sql.ErrNoRows{
		return nil,xlog.Error(err)
	}else {
		return nil,xlog.Error("参数不正确")
	}
}
