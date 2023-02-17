package courseModel

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

// Outh tb_outh
type Course struct {
	ID       int64
	Price    int64
	Cover    string
	Bg       string
	Desc     string
	Title    string
	Section  int64
	Duration int64
	User     int64
	Study    []string
	Tag      []tagMoel.Tag
}

func GetCourseFromID(ids, levelIds []int64, status int64) ([]*Course, error) {
	buf := strings.Builder{}
	args := make([]interface{}, 0)
	buf.WriteString("select qc.id,qc.cover,qc.title,qc.`desc`,qc.`prefix`,qc.price,qn.total,qvn.total,qh.total,qc.bg,qc.bg_prefix from qz_course qc " +
		"left join (select count(*) as total,course_id from qz_chapter GROUP BY course_id) qn on qn.course_id=qc.id " +
		"left join (SELECT count(qh.id) AS total,qzc.course_id FROM qz_chapter qzc INNER JOIN qz_chapter_video qcv ON qcv.chapter_id=qzc.id INNER JOIN qz_history qh ON qh.chapter_video_id=qcv.id GROUP BY qzc.course_id) qh on qh.course_id=qc.id " +
		"left join (select sum(qv.duration) as total,qzc.course_id from qz_chapter qzc inner join qz_chapter_video qcv on qcv.chapter_id=qzc.id inner join qz_video qv on qv.id=qcv.video_id GROUP BY qzc.course_id) qvn on qvn.course_id=qc.id " +
		"where 1=1 ")
	if status != enum.StatusAll {
		buf.WriteString(" and qc.`status`=?")
		args = append(args, status)
	} else {
		buf.WriteString(" and qc.`status`<>?")
		args = append(args, enum.StatusDelete)
	}

	utils.MysqlStringInUtils(&buf, ids, " and qc.id")
	utils.MysqlStringInUtils(&buf, levelIds, " and qc.level")
	result, err := db.GetCourse().Query(buf.String(), args...)
	if err != nil {
		return nil, err
	}
	list := make([]*Course, 0)
	mapList := make(map[int64]*Course)
	idList := make([]int64, 0)
	var id, price, total, durtotal, userTotal sql.NullInt64
	var cover, title, desc, prefix, bg, bgprefix sql.NullString
	for result.Next() {
		err := result.Scan(&id, &cover, &title, &desc, &prefix, &price, &total, &durtotal, &userTotal, &bg, &bgprefix)
		if err != nil {
			return nil, err
		}
		temp := &Course{
			ID:       id.Int64,
			Cover:    prefix.String + cover.String,
			Title:    title.String,
			Tag:      make([]tagMoel.Tag, 0),
			Desc:     desc.String,
			Price:    price.Int64,
			Section:  total.Int64,
			Duration: durtotal.Int64,
			User:     userTotal.Int64,
			Bg:       bgprefix.String + bg.String,
		}
		idList = append(idList, id.Int64)
		mapList[id.Int64] = temp
		list = append(list, temp)
	}
	tagList, err := tagMoel.GetTagRecord(enum.CourseTagType, idList)
	if err == nil {
		for _, v := range tagList {
			temp, ok := mapList[v.ContentId]
			if ok {
				temp.Tag = append(temp.Tag, tagMoel.Tag{
					ID:   v.ID,
					Name: v.Name,
				})
			}
		}
	}
	return list, nil
}

func GetAdminCourse(levelIds []int64, status int64) ([]Course, error) {
	buf := strings.Builder{}
	args := make([]interface{}, 0)
	buf.WriteString("select qc.id,qc.cover,qc.title,qc.`prefix` from qz_course qc " +
		"where 1=1")
	if status != enum.StatusAll {
		buf.WriteString(" and qc.`status`=?")
		args = append(args, status)
	} else {
		buf.WriteString(" and qc.`status`<>?")
		args = append(args, enum.StatusDelete)
	}
	utils.MysqlStringInUtils(&buf, levelIds, " and qc.id")
	result, err := db.GetCourse().Query(buf.String(), args...)
	if err != nil {
		return nil, err
	}
	list := make([]Course, 0)
	var id sql.NullInt64
	var cover, title, prefix sql.NullString
	for result.Next() {
		err := result.Scan(&id, &cover, &title, &prefix)
		if err != nil {
			return nil, err
		}
		list = append(list, Course{
			ID:    id.Int64,
			Cover: prefix.String + cover.String,
			Title: title.String,
		})
	}
	return list, nil
}

func GetAdminCourseDetail(id int64) (*Course, error) {
	result, err := db.GetCourse().Query(`select qc.cover,qc.bg,qc.bg_prefix,qc.title,qc.prefix,qc.desc
from qz_course qc
where qc.id=?
`, id)
	if err != nil {
		return nil, err
	}

	resultCourse := &Course{
		Tag:   make([]tagMoel.Tag, 0),
		Study: make([]string, 0),
	}
	var cover, title, prefix, desc, bg,bgPrefix sql.NullString
	for result.Next() {
		err := result.Scan(&cover,&bg,&bgPrefix, &title, &prefix, &desc)
		if err != nil {
			return nil, err
		}
		resultCourse.ID = id
		resultCourse.Cover = prefix.String + cover.String
		resultCourse.Bg = bgPrefix.String + bg.String
		resultCourse.Title = title.String
		resultCourse.Desc = desc.String
	}
	tag,err:=tagMoel.GetTagRecord(enum.CourseTagType,[]int64{id})
	if err==nil{
		resultCourse.Tag =tag
	}
	return resultCourse, nil
}

func EditCourseCover(id int64, cover, prefix string) error {
	_, err := db.GetCourse().Exec("update qz_course set `cover`=?,`prefix`=? where id=?", cover, prefix, id)
	return err
}

func EditCourseBgCover(id int64, cover, prefix string) error {
	_, err := db.GetCourse().Exec("update qz_course set `bg`=?,`bg_prefix`=? where id=?", cover, prefix, id)
	return err
}

func EditCourse(title, originName, originTitle, originDesc, originInfo, originInfoTitle, certificate string, id, level, createBy int64) (int64, error) {
	err := db.GetCourse().WithTransaction(func(tx *db.Tx) error {
		if id == 0 {
			args := make([]interface{}, 0)
			sqlStr := "insert into qz_course (`title`,`level`,`create_by`,`create_time`,`status`) values"
			err := db.GetCourse().WithTransaction(func(tx *db.Tx) error {
				placeHolder := "(?,?,?,?,?)"
				values := make([]string, 0)
				values = append(values, placeHolder)
				args = append(args, title)
				args = append(args, level)
				args = append(args, createBy)
				args = append(args, time.Now().Unix())
				args = append(args, enum.StatusNormal)
				sqlStr += strings.Join(values, ",")
				idResult, err := tx.Exec(sqlStr, args...)
				if err != nil {
					return xlog.Error("生成课程数据出错")
				}
				id, _ = idResult.LastInsertId()
				originId, err := EditOrigin(originName, originTitle, originDesc, originInfo, originInfoTitle, certificate, 0, tx)
				if err == nil {
					_, err = tx.Exec("update qz_course set  origin_id=? where id=?", originId, id)
					if err != nil {
						return xlog.Error("生成课程数据出错")
					}
				}
				return nil
			})
			return err
		} else {
			_, err := tx.Exec("update qz_course set `title`=? where id=?", title, id)
			if err != nil {
				return err
			}
			var originId sql.NullInt64
			err = tx.QueryRow("select origin_id from qz_course where id=?", id).Scan(&originId)
			if err != nil {
				return xlog.Error("生成课程数据出错")
			}
			oId, err := EditOrigin(originName, originTitle, originDesc, originInfo, originInfoTitle, certificate, originId.Int64, tx)
			if err != nil {
				return xlog.Error("生成课程数据出错")
			}
			if originId.Int64 == 0 {
				_, err = tx.Exec("update qz_course set origin_id=? where id=?", oId, id)
				if err != nil {
					return xlog.Error("生成课程数据出错")
				}
			}
			return nil
			//return tagMoel.EditTagRecord(tagId, id, enum.CourseTagType, tx)
		}

	})
	return id, err
}
