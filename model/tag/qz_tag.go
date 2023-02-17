package tagMoel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"strings"
)

type Tag struct {
	ID       int64
	Cover    string
	Name      string
	Content  string
	ContentId int64
}


func GetTagList(status,queryId int64) ([]Tag, error) {
	buf := strings.Builder{}
	args := make([]interface{}, 0)
	buf.WriteString("select id,content,title,cover,prefix from qz_tag where 1=1")
	if status==enum.StatusAll{
		buf.WriteString(" and status<>?")
		args=append(args,enum.StatusDelete)
	} else {
		buf.WriteString(" and status=?")
		args=append(args,status)
	}

	if queryId!=0{
		buf.WriteString(" and id=?")
		args=append(args,queryId)
	}
	result, err := db.GetCourse().Query(buf.String(), args...)
	if err != nil {
		return nil, err
	}
	list := make([]Tag, 0)
	var id sql.NullInt64
	var cover, coverPrefix, content, title sql.NullString
	for result.Next() {
		err := result.Scan(&id, &content, &title,&cover, &coverPrefix)
		if err != nil {
			return nil, err
		}
		list = append(list, Tag{
			ID:       id.Int64,
			Cover:    coverPrefix.String + cover.String,
			Name:    title.String,
			Content: content.String,
		})
	}
	return list, nil
}

func AddTag(title,content string)(int64,error){
	result,err:=db.GetCourse().Exec("insert into qz_tag (title,content,status) values (?,?,1)",title,content)
	if err==nil{
		id,err:=result.LastInsertId()
		if err==nil{
			return id,nil
		}
	}
	return 0,err
}

func EditTagCover(cover,prefix string,id int64)error{
	_,err:=db.GetCourse().Exec("update qz_tag set cover=?,prefix=? where id=?",cover,prefix,id)
	return err
}
func EditTagDetail(title,content string,id int64)error{
	_,err:=db.GetCourse().Exec("update qz_tag set title=?,content=? where id=?",title,content,id)
	return err
}



func EditTagRecord(idList []int64,contentId,typeInt int64,tx *db.Tx)error{
	args:=make([]interface{},0)
	_,err:=tx.Exec("delete from qz_tag_record where type=? and content_id=?",typeInt,contentId)
	if err != nil {
		return xlog.Error("修改标签出错")
	}
	if len(idList)==0{
		return nil
	}
	sqlStr := "insert into qz_tag_record (`tag_id`,`type`,`content_id`) values"
	placeHolder := "(?,?,?)"
	values := make([]string, 0)
	for _,value:=range idList{
		values = append(values, placeHolder)
		args=append(args,value)
		args=append(args,typeInt)
		args=append(args,contentId)
	}
	sqlStr += strings.Join(values, ",")
	_, err = tx.Exec(sqlStr, args...)
	if err != nil {
		return xlog.Error("修改标签出错")
	}else{
		return nil
	}
}

func GetTagRecord(typeInt int64,contentIdList []int64)([]Tag,error){
	args:=make([]interface{},0)
	var buf strings.Builder
	buf.WriteString("select tb1.id,tb1.title,tb.content_id from qz_tag_record tb inner join qz_tag tb1 on tb1.id=tb.tag_id where tb.type=? ")
	args=append(args,typeInt)
	utils.MysqlStringInUtils(&buf,contentIdList," and tb.content_id")
	result, err := db.GetCourse().Query(buf.String(),args...)
	if err != nil {
		return nil, err
	}
	list := make([]Tag, 0)
	var id,contentId sql.NullInt64
	var   title sql.NullString
	for result.Next() {
		err := result.Scan(&id,  &title,&contentId)
		if err != nil {
			return nil, err
		}
		list = append(list, Tag{
			ID:       id.Int64,
			Name:    title.String,
			ContentId:contentId.Int64,

		})
	}
	return list, nil
}

