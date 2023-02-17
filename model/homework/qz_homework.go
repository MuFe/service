package homeworkModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"mufe_service/model/tag"
	"strings"
	"time"
)

type HomeWorkData struct {
	ID        int64
	Title     string
	Index     int64
	Number    int64
	Progress  float64
	Time      int64
	Duration  int64
	Total     int64
	ContentId int64
	Content   string
	Cover     string
	Info      int64
}

type HomeWorkGroup struct {
	Desc string
	Id   int64
	List []*HomeWorkData
}

type HomeWorkInfo struct {
	ID        int64
	InfoId    int64
	GroupId   int64
	Title     string
	Cover     string
	Video     string
	Content   string
	ContentId int64
	Level     int64
	TagList   []tagMoel.Tag
}

func GetHomeWorkInfo(page, size, status int64) ([]*HomeWorkInfo, error) {
	resultList := make([]*HomeWorkInfo, 0)
	limit := (page - 1) * size
	result, err := db.GetCourse().Query("select tb.id,tb.title,tb.content,tb.cover,tb.prefix,tb.video,tb.video_prefix from  qz_homework_info tb  where tb.status=? order by id desc limit ?,? ", status, limit, size)
	if err == nil {
		var id sql.NullInt64
		var cover, coverPrefix, video, videoPrefix, title, content sql.NullString
		idList := make([]int64, 0)
		resultMap := make(map[int64]*HomeWorkInfo)
		for result.Next() {
			err = result.Scan(&id, &title, &content, &cover, &coverPrefix, &video, &videoPrefix)
			if err == nil {
				temp := &HomeWorkInfo{
					ID:      id.Int64,
					Title:   title.String,
					Cover:   coverPrefix.String + cover.String,
					Video:   videoPrefix.String + video.String,
					Content: content.String,
					TagList: make([]tagMoel.Tag, 0),
				}
				resultList = append(resultList, temp)
				resultMap[id.Int64] = temp
				idList = append(idList, id.Int64)
			}
		}
		tagList, err := tagMoel.GetTagRecord(enum.HomeWorkTagType, idList)
		if err == nil {
			for _, v := range tagList {
				temp, ok := resultMap[v.ContentId]
				if ok {
					temp.TagList = append(temp.TagList, tagMoel.Tag{
						ID:   v.ID,
						Name: v.Name,
					})
				}
			}
		}
	}
	return resultList, err
}

func GetHomeWork(status int64, contentId []int64) ([]*HomeWorkInfo, error) {
	resultList := make([]*HomeWorkInfo, 0)
	var buf strings.Builder
	buf.WriteString("select tb1.id,tb.id,tb.title,tb.content,tb.cover,tb.prefix,tb.video,tb.video_prefix,tb1.content_id,tb.level from  qz_homework_info tb inner join qz_homework tb1 on tb1.home_info_id=tb.id where tb.status=?")
	utils.MysqlStringInUtils(&buf, contentId, " and tb1.content_id")
	result, err := db.GetCourse().Query(buf.String(), status)
	if err == nil {
		idList := make([]int64, 0)
		resultMap := make(map[int64]*HomeWorkInfo)
		var id, infoId, contentId, level sql.NullInt64
		var cover, coverPrefix, video, videoPrefix, title, content sql.NullString
		for result.Next() {
			err = result.Scan(&id, &infoId, &title, &content, &cover, &coverPrefix, &video, &videoPrefix, &contentId, &level)
			if err == nil {
				temp := &HomeWorkInfo{
					ID:        id.Int64,
					Title:     title.String,
					InfoId:    infoId.Int64,
					Cover:     coverPrefix.String + cover.String,
					Video:     videoPrefix.String + video.String,
					Content:   content.String,
					ContentId: contentId.Int64,
					Level:     level.Int64,
				}
				resultList = append(resultList, temp)
				resultMap[id.Int64] = temp
				idList = append(idList, id.Int64)
			}
		}
		tagList, err := tagMoel.GetTagRecord(enum.HomeWorkTagType, idList)
		if err == nil {
			for _, v := range tagList {
				temp, ok := resultMap[v.ContentId]
				if ok {
					temp.TagList = append(temp.TagList, tagMoel.Tag{
						ID:   v.ID,
						Name: v.Name,
					})
				}
			}
		}
	}
	return resultList, err
}

func GetHomeWorkDetail(id int64) (*HomeWorkInfo, error) {
	var cover, coverPrefix, video, videoPrefix, title, content sql.NullString
	var contentId,tagId sql.NullInt64
	buffer:=strings.Builder{}
	args:=make([]interface{},0)
	buffer.WriteString("select tb.id,tb.title,tb.content,tb.cover,tb.prefix,tb.video,tb.video_prefix,tb1.content_id from  qz_homework_info tb inner join qz_homework tb1 on tb1.home_info_id=tb.id where tb1.id=?")
	args=append(args,id)
	err := db.GetCourse().QueryRow(buffer.String(),args...).Scan(&tagId,&title, &content, &cover, &coverPrefix, &video, &videoPrefix,&contentId)
	if err == nil {
		temp := &HomeWorkInfo{
			Title:   title.String,
			Cover:   coverPrefix.String + cover.String,
			Video:   videoPrefix.String + video.String,
			Content: content.String,
			ContentId:contentId.Int64,
		}
		tagList, err := tagMoel.GetTagRecord(enum.HomeWorkTagType, []int64{tagId.Int64})
		if err == nil {
			temp.TagList = tagList
		}
		return temp, nil
	}
	return nil, err
}

func GetHomeWorkDetailInfo(infoId int64) (*HomeWorkInfo, error) {
	var cover, coverPrefix, video, videoPrefix, title, content sql.NullString
	var level sql.NullInt64
	buffer:=strings.Builder{}
	args:=make([]interface{},0)
	buffer.WriteString("select tb.title,tb.content,tb.cover,tb.prefix,tb.video,tb.video_prefix,tb.level from  qz_homework_info tb where tb.id=?")
	args=append(args,infoId)
	err := db.GetCourse().QueryRow(buffer.String(),args...).Scan(&title, &content, &cover, &coverPrefix, &video, &videoPrefix,&level)
	if err == nil {
		temp := &HomeWorkInfo{
			Title:   title.String,
			Cover:   coverPrefix.String + cover.String,
			Video:   videoPrefix.String + video.String,
			Content: content.String,
			Level:level.Int64,
		}
		tagList, err := tagMoel.GetTagRecord(enum.HomeWorkTagType, []int64{infoId})
		if err == nil {
			temp.TagList = tagList
		}
		return temp, nil
	}
	return nil, err
}

func EditHomeWork(title string, id,level int64, tagList []int64) (int64, error) {
	err := db.GetSchool().WithTransaction(func(tx *db.Tx) error {
		if id == 0 {
			idResult, err := tx.Exec("insert into qz_homework_info (title,sort,content,`level`) values(?,?,'',?)", title, enum.SortMaxValue,level)
			if err == nil {
				id, err = idResult.LastInsertId()
				if err == nil {
					return tagMoel.EditTagRecord(tagList, id, enum.HomeWorkTagType, tx)
				} else {
					return err
				}
			} else {
				return err
			}
		} else {
			_, err := tx.Exec("update qz_homework_info set  title=?,`level`=? where id=?", title,level, id)
			if err == nil {
				return tagMoel.EditTagRecord(tagList, id, enum.HomeWorkTagType, tx)
			} else {
				return err
			}
		}
	})
	return id, err
}

func EditHomeWorkInfo(typeInt, id int64, content, prefix string) error {
	if typeInt == enum.HomeWork_Edit_Cover {
		_, err := db.GetCourse().Exec("update qz_homework_info set `cover`=?,`prefix`=? where id=?", content, prefix, id)
		return err
	} else if typeInt == enum.HomeWork_Edit_Video {
		_, err := db.GetCourse().Exec("update qz_homework_info set `video`=?,`video_prefix`=? where id=?", content, prefix, id)
		return err
	} else if typeInt == enum.HomeWork_Edit_Content {
		_, err := db.GetCourse().Exec("update qz_homework_info set `content`=? where id=?", content, id)
		return err
	}
	return nil
}

func UnBindHomeWork(id int64) error {
	_, err := db.GetCourse().Exec("delete from qz_homework where id=?", id)
	return err
}

func BindHomeWork(contentId int64, homeInfoList []int64, tx *db.Tx) error {

	sqlStr := "insert into qz_homework (home_info_id,content_id) values"
	placeHolder := "(?,?)"

	var values []string
	args := make([]interface{}, 0)
	for _, info := range homeInfoList {
		values = append(values, placeHolder)
		args = append(args, info)
		args = append(args, contentId)
	}
	sqlStr += strings.Join(values, ",")
	if tx == nil {
		return db.GetCourse().WithTransaction(func(tx *db.Tx) error {
			return startBind(contentId, sqlStr, args, tx)
		})
	} else {
		return startBind(contentId, sqlStr, args, tx)
	}
}

func startBind(contentId int64, sqlStr string, args []interface{}, tx *db.Tx) error {
	_, err := tx.Exec("delete from qz_homework where content_id=?", contentId)
	if err != nil {
		return xlog.Error(err)
	} else {
		if len(args) > 0 {
			_, err = tx.Exec(sqlStr, args...)
			if err != nil {
				return xlog.Error(err)
			} else {
				return nil
			}
		} else {
			return nil
		}
	}
}

func AddHomeWork(classIds, ids []int64, uid, timeInt int64, desc string) error {

	sqlStr := "insert into qz_homework_record (homework_id,class_id,uid,create_time,group_id) values"
	placeHolder := "(?,?,?,?,?)"

	var values []string
	var args []interface{}

	err := db.GetCourse().WithTransaction(func(tx *db.Tx) error {
		_, err := tx.Exec("update  qz_homework_group set `status`=? where homework_time=?", enum.StatusDelete, timeInt)
		if err != nil {
			return xlog.Error(err)
		}
		result, err := tx.Exec("insert into qz_homework_group (create_time,homework_time,`desc`,`status`) values (?,?,?,?)", time.Now().Unix(), timeInt, desc, enum.StatusNormal)
		if err == nil {
			id, err := result.LastInsertId()
			if err == nil {
				for _, info := range classIds {
					for _, v := range ids {
						if info == 0 || v == 0 {
							continue
						}
						values = append(values, placeHolder)
						args = append(args, v)
						args = append(args, info)
						args = append(args, uid)
						args = append(args, time.Now().Unix())
						args = append(args, id)
					}
				}
				sqlStr += strings.Join(values, ",")
				_, err = tx.Exec(sqlStr, args...)
				return err
			}
		}
		return err
	})
	if err != nil {
		return xlog.Error(err)
	} else {
		return nil
	}
}

func HomeWorkList(time, classId, uid int64) ([]*HomeWorkData, error) {
	var buf strings.Builder
	args := make([]interface{}, 0)
	buf.WriteString(`select tb.id,tb.homework_id,tb12.content,tb12.title,tb12.cover,tb12.prefix,tb3.number,tb2.class_id from qz_homework_record tb 
inner join qz_homework tb1 on tb1.id=tb.homework_id
inner join qz_homework_info tb12 on tb12.id=tb1.home_info_id
inner join qz_class_record tb2 on tb2.class_id=tb.class_id
inner join qz_class tb3 on tb3.id=tb.class_id
inner join qz_homework_group tb4 on tb4.id=tb.group_id
where tb2.uid=? and tb4.status=?`)
	args = append(args, uid, enum.StatusNormal)
	if classId != 0 {
		buf.WriteString(" and tb2.class_id=?")
		args = append(args, classId)
	}
	if time != 0 {
		buf.WriteString(" and tb4.homework_time=?")
		args = append(args, time)
	}
	result, err := db.GetCourse().Query(buf.String(), args...)
	idList := make([]int64, 0)
	homeWorkList := make([]int64, 0)
	list := make([]*HomeWorkData, 0)
	resultMap := make(map[int64]*HomeWorkData)
	resultHomeMap := make(map[int64]int64)
	if err == nil {
		var id, homeWorkId, number, class sql.NullInt64
		var title, content, coverPrefix, cover sql.NullString
		for result.Next() {
			err = result.Scan(&id, &homeWorkId, &content, &title, &cover, &coverPrefix, &number, &class)
			if err == nil {
				idList = append(idList, id.Int64)
				homeWorkList = append(homeWorkList, homeWorkId.Int64)
				temp := &HomeWorkData{
					ID:      id.Int64,
					Title:   title.String,
					Total:   number.Int64,
					Content: content.String,
					ContentId:homeWorkId.Int64,
					Cover:   coverPrefix.String + cover.String,
				}
				resultMap[id.Int64] = temp
				resultHomeMap[homeWorkId.Int64]=id.Int64
				list = append(list, temp)
			}
		}
	}
	if len(homeWorkList) != 0 {
		var buf1 strings.Builder
		buf1.WriteString("select count(uid),homework_id from qz_homework_submit_record ")
		utils.MysqlStringInUtils(&buf1, homeWorkList, " where homework_id ")
		buf1.WriteString(" GROUP BY homework_id")
		result1, err := db.GetCourse().Query(buf1.String())
		if err == nil {
			var total, homeWorkId sql.NullInt64
			for result1.Next() {
				err = result1.Scan(&total, &homeWorkId)
				if err == nil {
					idTemp, ok := resultHomeMap[homeWorkId.Int64]
					if ok {
						temp,ok:=resultMap[idTemp]
						if ok{
							temp.Total=total.Int64
						}
					}
				}
			}
		}
	}
	if len(idList) != 0 {
		var buf1 strings.Builder
		buf1.WriteString(`select count(tb2.uid),tb.id from qz_homework_record tb 
inner join qz_class_record tb2 on tb2.class_id=tb.class_id ` )
		utils.MysqlStringInUtils(&buf1, idList, " where tb.id  ")
		buf1.WriteString(" and tb2.type=? GROUP BY tb.id")
		result1, err := db.GetCourse().Query(buf1.String(),enum.ClassStudentType)
		if err == nil {
			var total, id sql.NullInt64
			for result1.Next() {
				err = result1.Scan(&total, &id)
				if err == nil {
					temp, ok := resultMap[id.Int64]
					if ok {
						if temp.Total==total.Int64{
							temp.Progress=100
						}else{
							temp.Progress = 100*float64(temp.Total) / float64(total.Int64)
						}

					}
				}
			}
		}
	}
	return list, nil

}

func StudentHomeWorkList(time, uid, nowId int64) (HomeWorkGroup, error) {
	var buf strings.Builder
	args := make([]interface{}, 0)
	buf.WriteString(`select tb.id,tb.homework_id,tb12.content,tb12.title,tb12.cover,tb12.prefix,tb4.desc,tb4.id,tb12.id from qz_homework_record tb 
inner join qz_homework tb1 on tb1.id=tb.homework_id
inner join qz_homework_info tb12 on tb12.id=tb1.home_info_id
inner join qz_homework_group tb4 on tb4.id=tb.group_id
inner join qz_class_record tb2 on tb2.class_id=tb.class_id
inner join qz_class tb5 on tb5.id=tb.class_id
where tb4.status=? and tb5.status=?`)
	args = append(args, enum.StatusNormal, enum.StatusNormal)
	if nowId != 0 {
		buf.WriteString(" and tb4.id=?")
		args = append(args, nowId)
		if uid != 0 {
			buf.WriteString(" and tb2.uid=?")
			args = append(args, uid)
		}

	} else {
		buf.WriteString(" and tb4.homework_time=? and tb2.uid=?")
		args = append(args, time)
		args = append(args, uid)
	}
	result, err := db.GetCourse().Query(buf.String(), args...)
	idList := make([]int64, 0)
	resultRe := HomeWorkGroup{
		List: make([]*HomeWorkData, 0),
	}
	resultMap := make(map[int64]*HomeWorkData)
	if err == nil {
		var id, homeWorkId, groupId, infoId sql.NullInt64
		var title, desc, content, cover, coverPrefix sql.NullString
		for result.Next() {
			err = result.Scan(&id, &homeWorkId, &content, &title, &cover, &coverPrefix, &desc, &groupId, &infoId)
			if err == nil {
				idList = append(idList, homeWorkId.Int64)
				resultRe.Desc = desc.String
				temp := &HomeWorkData{
					ID:      homeWorkId.Int64,
					Title:   title.String,
					Content: content.String,
					Info:    infoId.Int64,
					Cover:   coverPrefix.String + cover.String,
				}
				resultMap[homeWorkId.Int64] = temp
				resultRe.List = append(resultRe.List, temp)
			}
		}
	}
	if len(idList) != 0 {
		var buf1 strings.Builder
		buf1.WriteString("select count(uid),homework_id from qz_homework_submit_record ")
		utils.MysqlStringInUtils(&buf1, idList, " where homework_id ")
		buf1.WriteString(" and uid=? ")
		buf1.WriteString(" GROUP BY homework_id")
		result1, err := db.GetCourse().Query(buf1.String(), uid)
		if err == nil {
			var total, homeWorkId sql.NullInt64
			for result1.Next() {
				err = result1.Scan(&total, &homeWorkId)
				if err == nil {
					temp, ok := resultMap[homeWorkId.Int64]
					if ok {
						temp.Progress = 100

					}
				}
			}
		}
	}
	return resultRe, nil

}

func GetHomeWorkRecord(id int64) ([]int64, []int64) {
	var uid sql.NullInt64
	finish := make([]int64, 0)
	finishMap := make(map[int64]bool)
	result, err := db.GetCourse().Query("select uid from qz_homework_submit_record where homework_id=?", id)
	if err == nil {
		for result.Next() {
			err = result.Scan(&uid)
			if err == nil {
				finish = append(finish, uid.Int64)
				finishMap[uid.Int64] = true
			}
		}
	}
	inComplete := make([]int64, 0)
	result, err = db.GetCourse().Query(`select tb2.uid from qz_homework_record tb 
inner join qz_class_record tb2 on tb2.class_id=tb.class_id where tb.id=? and tb2.type=?`, id, enum.ClassStudentType)
	if err == nil {
		for result.Next() {
			err = result.Scan(&uid)
			if err == nil {
				_, ok := finishMap[uid.Int64]
				if !ok {
					inComplete = append(inComplete, uid.Int64)
				}

			}
		}
	}
	return finish, inComplete
}

func FinishHomeWork(id, uid, score int64) error {
	return db.GetCourse().WithTransaction(func(tx *db.Tx) error {
		_, err := tx.Exec("delete from qz_homework_submit_record where uid=? and homework_id=?", uid, id)
		if err == nil {
			_, err = tx.Exec("insert into qz_homework_submit_record (uid,homework_id,`score`) values (?,?,?)", uid, id, score)
			return err
		}
		return err
	})
}


func CancelHomeWorkRecord(quitUid []int64)error{
	var err error
	if len(quitUid)==0{
		return nil
	}
	var buf strings.Builder
	buf.WriteString("delete from qz_homework_submit_record ")
	utils.MysqlStringInUtils(&buf,quitUid," where uid")
	_, err=db.GetSchool().Exec(buf.String())

	return err
}
