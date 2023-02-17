package courseModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"mufe_service/model/homework"
	tagMoel "mufe_service/model/tag"
	"strings"
	"time"
)

type Chapter struct {
	ID       int64
	Price    int64
	Cover    string
	Video    string
	Plan     string
	Desc     string
	TagList  []tagMoel.Tag
	Title    string
	Level    string
	LevelId  int64
	User     int64
	Section  int64
	Duration int64
	HomeWork []int64
}

func GetChapterFromID(ids []int64, page, size, courseId, status int64) ([]*Chapter, error) {
	buf := strings.Builder{}
	args := make([]interface{}, 0)
	buf.WriteString("select qc.id,qc.cover,qc.title,qc.`desc`,qc.`prefix`,qc.price,qn.total,qvn.total,qh.total,qc.plan from qz_chapter qc " +
		"left join (select count(*) as total,chapter_id from qz_chapter_video GROUP BY chapter_id) qn on qn.chapter_id=qc.id " +
		"left join (SELECT count(qh.id) AS total,qcv.chapter_id FROM  qz_history qh  INNER JOIN qz_chapter_video qcv ON qh.chapter_video_id=qcv.id GROUP BY qcv.chapter_id) qh on qh.chapter_id=qc.id  " +
		"left join (select sum(qv.duration) as total,qcv.chapter_id from  qz_chapter_video qcv  inner join qz_video qv on qv.id=qcv.video_id GROUP BY qcv.chapter_id) qvn on qvn.chapter_id=qc.id  " +
		"where 1=1 ")
	if status != enum.StatusAll {
		buf.WriteString(" and qc.`status`=?")
		args = append(args, status)
	} else {
		buf.WriteString(" and qc.`status`<>?")
		args = append(args, enum.StatusDelete)
	}

	utils.MysqlStringInUtils(&buf, ids, " and qc.id")
	if courseId != 0 {
		buf.WriteString(" and qc.course_id=?")
		args = append(args, courseId)
	}
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}
	buf.WriteString(" order by qc.sort desc,qc.id asc ")
	start := (page - 1) * size
	buf.WriteString(" limit ?,?")
	args = append(args, start, size)
	result, err := db.GetCourse().Query(buf.String(), args...)
	if err != nil {
		return nil, err
	}
	list := make([]*Chapter, 0)
	mapList := make(map[int64]*Chapter)
	idList := make([]int64, 0)
	var id, price, total, durtotal, userTotal sql.NullInt64
	var cover, title, desc, prefix,plan sql.NullString
	for result.Next() {
		err := result.Scan(&id, &cover, &title, &desc, &prefix, &price, &total, &durtotal, &userTotal,&plan)
		if err != nil {
			return nil, err
		}
		idList = append(idList, id.Int64)
		v := &Chapter{
			ID:       id.Int64,
			Cover:    prefix.String + cover.String,
			Title:    title.String,
			Desc:     desc.String,
			Price:    price.Int64,
			Section:  total.Int64,
			Plan:plan.String,
			Duration: durtotal.Int64,
			User:     userTotal.Int64,
		}
		list = append(list, v)
		mapList[id.Int64] = v
	}
	resultTag, err := tagMoel.GetTagRecord(enum.ChapterTagType, idList)
	if err == nil {
		for _, v := range resultTag {
			temp, ok := mapList[v.ContentId]
			if ok {
				temp.TagList = append(temp.TagList, v)
			}
		}
	}
	return list, nil
}

func GetAdminChapter(ids []int64, page, size, courseId, status int64) ([]Chapter, error) {
	buf := strings.Builder{}
	args := make([]interface{}, 0)
	buf.WriteString("select qc.id,qc.cover,qc.title,qc.`prefix` from qz_chapter qc " +
		"where 1=1 ")
	if status != enum.StatusAll {
		buf.WriteString(" and qc.`status`=?")
		args = append(args, status)
	} else {
		buf.WriteString(" and qc.`status`<>?")
		args = append(args, enum.StatusDelete)
	}

	utils.MysqlStringInUtils(&buf, ids, " and qc.id")
	if courseId != 0 {
		buf.WriteString(" and qc.course_id=?")
		args = append(args, courseId)
	}
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}
	buf.WriteString(" order by qc.sort desc,qc.id asc ")
	start := (page - 1) * size
	buf.WriteString(" limit ?,?")
	args = append(args, start, size)
	result, err := db.GetCourse().Query(buf.String(), args...)
	if err != nil {
		return nil, err
	}
	list := make([]Chapter, 0)
	var id sql.NullInt64
	var cover, title, prefix sql.NullString
	for result.Next() {
		err := result.Scan(&id, &cover, &title, &prefix)
		if err != nil {
			return nil, err
		}
		list = append(list, Chapter{
			ID:    id.Int64,
			Cover: prefix.String + cover.String,
			Title: title.String,
		})
	}
	return list, nil
}

func GetAdminChapterDetail(id int64) ([]Chapter, error) {
	buf := strings.Builder{}
	args := make([]interface{}, 0)
	buf.WriteString("select qc.cover,qc.title,qc.`prefix`,qc.video,qc.desc,qc.video_prefix,qc.plan from qz_chapter qc " +
		"where 1=1 ")

	buf.WriteString(" and qc.id=?")
	args = append(args, id)

	result, err := db.GetCourse().Query(buf.String(), args...)
	if err != nil {
		return nil, err
	}
	list := make([]Chapter, 0)

	var cover, title, prefix, video, desc, videoPrefix, plan sql.NullString
	for result.Next() {
		err := result.Scan(&cover, &title, &prefix, &video, &desc, &videoPrefix, &plan)
		if err != nil {
			return nil, err
		}
		list = append(list, Chapter{
			ID:    id,
			Cover: prefix.String + cover.String,
			Video: videoPrefix.String + video.String,
			Plan:  plan.String,
			Title: title.String,
			Desc:  desc.String,
		})
	}
	return list, nil
}

func EditChapter(title, desc string, id, sourceId, createBy int64, homeWorkList  []int64) (int64, error) {
	err := db.GetCourse().WithTransaction(func(tx *db.Tx) error {
		if id == 0 {
			sqlStr := "insert into qz_chapter (`title`,`desc`,`course_id`,`create_by`,`create_time`,`sort`,`plan`) values"
			args := make([]interface{}, 0)
			placeHolder := "(?,?,?,?,?,?,?)"
			values := make([]string, 0)
			values = append(values, placeHolder)
			args = append(args, title)
			args = append(args, desc)
			args = append(args, sourceId)
			args = append(args, createBy)
			args = append(args, time.Now().Unix())
			args = append(args, enum.SortMaxValue)
			args = append(args, "")
			sqlStr += strings.Join(values, ",")
			newIdResult, err := tx.Exec(sqlStr, args...)
			if err != nil {
				return xlog.Error("生成章节数据出错")
			}
			id, err = newIdResult.LastInsertId()
			if err != nil {
				return xlog.Error("生成章节数据出错")
			}
			return homeworkModel.BindHomeWork(id, homeWorkList, tx)
		} else {
			_, err := tx.Exec("update qz_chapter set `title`=?,`desc`=? where id=?", title, desc, id)
			if err != nil {
				return err
			}
			return homeworkModel.BindHomeWork(id, homeWorkList, tx)
		}

	})
	return id, err
}

func EditChapterCover(cover, prefix string, id, typeInt int64) error {
	if typeInt == enum.EDIT_COVER {
		_, err := db.GetCourse().Exec("update qz_chapter set `cover`=?,`prefix`=? where id=?", cover, prefix, id)
		return err
	} else if typeInt == enum.EDIT_VIDEO {
		_, err := db.GetCourse().Exec("update qz_chapter set `video`=?,`video_prefix`=? where id=?", cover, prefix, id)
		return err
	} else if typeInt == enum.EDIT_PLAN {
		_, err := db.GetCourse().Exec("update qz_chapter set `plan`=? where id=?", cover, id)
		return err
	}
	return nil

}

//修改排序
func EditChapterSort(chapterId, sort int64) error {
	//列表页面修改
	err := db.GetCourse().WithTransaction(func(tx *db.Tx) error {
		var courserId sql.NullInt64
		err := tx.QueryRow("select course_id from qz_chapter where id=?", chapterId).Scan(&courserId)
		if err != nil {
			return xlog.Error(err)
		}
		if sort == enum.SortMax {
			_, err = tx.Exec(`update qz_chapter tb set sort=sort-1 where tb.id<>? and tb.course_id=?`, chapterId, courserId)
			if err != nil {
				return xlog.Error(err)
			}
			_, err = tx.Exec(`update qz_chapter set sort=? where id=?`, enum.SortMaxValue, chapterId)
			if err != nil {
				return xlog.Error(err)
			}
		} else {
			var sSort int64
			err = tx.QueryRow(`select sort from qz_chapter where id=?`, chapterId).Scan(&sSort)
			if err != nil {
				return xlog.Error(err)
			}
			//获取要交换的排序
			var nextChapterId, nextSort int64
			if sort == enum.SortDown {
				err = tx.QueryRow(`select id,sort from qz_chapter where id<>? and course_id=? and sort<=?  order by sort desc,id asc`, chapterId, courserId, sSort).Scan(&nextChapterId, &nextSort)
			} else {
				err = tx.QueryRow(`select id,sort from qz_chapter where id<>? and course_id=? and sort>=?  order by sort desc,id asc`, chapterId, courserId, sSort).Scan(&nextChapterId, &nextSort)
			}

			if err != nil {
				return xlog.Error(err)
			}
			//交换
			_, err = tx.Exec(`update qz_chapter set sort=? where id=?`, nextSort, chapterId)
			if err != nil {
				return xlog.Error(err)
			}
			_, err = tx.Exec(`update qz_chapter set sort=? where id=?`, sSort, nextChapterId)
			if err != nil {
				return xlog.Error(err)
			}
		}

		return nil
	})
	return err
}
