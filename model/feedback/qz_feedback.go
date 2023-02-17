package feedbackmodel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/xlog"
	"strings"
	"time"
)

type FeedbackData struct {
	ID      int64
	Uid     int64
	Name     string
	Content string
	Time    int64
}

func GetFeedback(status int64) ([]*FeedbackData, error) {
	resultList := make([]*FeedbackData, 0)
	var buf strings.Builder
	args:=make([]interface{},0)
	xlog.Info(status)
	buf.WriteString("select tb.id,tb.content,tb.create_by,tb.create_time,tb1.nick_name from qz_feedback tb left join qz_user tb1 on tb.create_by=tb1.uid")
	if status!=enum.StatusAll{
		buf.WriteString(" where tb.status=?")
		args=append(args,status)
	} else {
		buf.WriteString(" where tb.status<>?")
		args=append(args,enum.StatusDelete)
	}
	result, err := db.GetUserDb().Query(buf.String(),args...)
	if err == nil {
		var id,createTime,createBy sql.NullInt64
		var content,name sql.NullString
		for result.Next() {
			err = result.Scan(&id, &content, &createBy, &createTime, &name)
			if err == nil {
				nowName:="游客"
				if name.String!=""{
					nowName=name.String
				}
				resultList = append(resultList, &FeedbackData{
					ID:     id.Int64,
					Name:   nowName,
					Content:  content.String,
					Time: createTime.Int64,
					Uid:    createBy.Int64,
				})
			}
		}
	}
	return resultList, err
}

func AddFeedback(uid int64,content string) error {
	_,err:=db.GetUserDb().Exec("insert into qz_feedback (content,create_by,create_time) values (?,?,?)",content,uid,time.Now().Unix())
	return err
}

func EditFeedback(id,status int64) error {
	_,err:=db.GetUserDb().Exec("update  qz_feedback set status=? where id=?",status,id)
	return err
}
