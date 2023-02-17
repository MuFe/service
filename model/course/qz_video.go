package courseModel

import (
	"database/sql"
	"fmt"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"mufe_service/model/tag"
	"strconv"
	"strings"
	"time"
)

type Video struct {
	ID       int64
	VideoId  int64
	Cover    string
	Url      string
	DownUrl  string
	ChapterId int64
	Duration int64
	Title    string
	TagList  []tagMoel.Tag
	Content  string
}

type AddVideoData struct {
	ID            int64
	VideoId       int64
	Cover         string
	CoverPrefix   string
	Url           string
	UrlPrefix     string
	DownUrl       string
	DownUrlPrefix string
	Duration      int64

	Title   string
	Content string
	TagId   []int64
	Price   int64
}

type ProgressInfoDetail struct {
	ChapterId int64
	Id        int64
}
type ProgressInfo struct {
	ChapterIndex int64
	VideoIndex   int64
	Progress     float64
	List         []ProgressInfoDetail
}

func GetChapterVideo(chapterIdList []int64) (map[int64][]Video, error) {
	buf := strings.Builder{}
	args := make([]interface{}, 0)
	buf.WriteString("select qcv.chapter_id,qcv.id,qv.cover,qv.cover_prefix,qv.url,qv.url_prefix,qv.down_url,qv.down_prefix," +
		"qv.duration,qcv.title,qv.id from qz_video qv inner join qz_chapter_video qcv on qcv.video_id=qv.id where qcv.`status`=1")
	utils.MysqlStringInUtils(&buf, chapterIdList, " and qcv.chapter_id")
	buf.WriteString(" order by qcv.chapter_id  desc,qcv.sort desc,qcv.id asc ")

	result, err := db.GetCourse().Query(buf.String(), args...)
	if err != nil {
		return nil, err
	}
	resultMap := make(map[int64][]Video, 0)

	var chapterId, id, videoId, duration sql.NullInt64
	var cover, coverPrefix, url, urlPrefix, down, downPrefix, title sql.NullString
	for result.Next() {
		err := result.Scan(&chapterId, &id, &cover, &coverPrefix, &url, &urlPrefix, &down, &downPrefix, &duration, &title, &videoId)
		if err != nil {
			return nil, err
		}
		list, ok := resultMap[chapterId.Int64]
		if !ok {
			list = make([]Video, 0)
		}
		list = append(list, Video{
			ID:       id.Int64,
			VideoId:  videoId.Int64,
			Cover:    coverPrefix.String + cover.String,
			Url:      urlPrefix.String + url.String,
			DownUrl:  downPrefix.String + down.String,
			Title:    title.String,
			Duration: duration.Int64,
		})
		resultMap[chapterId.Int64] = list
	}
	return resultMap, nil
}

func GetAdminChapterVideo(chapterIdList []int64) (map[int64][]Video, error) {
	buf := strings.Builder{}
	args := make([]interface{}, 0)
	buf.WriteString("select qcv.chapter_id,qcv.id,qcv.cover,qcv.prefix,qv.url,qv.url_prefix,qv.down_url,qv.down_prefix," +
		"qv.duration,qcv.title,qv.id from qz_chapter_video qcv left join qz_video qv on qcv.video_id=qv.id where qcv.`status`=1")
	utils.MysqlStringInUtils(&buf, chapterIdList, " and qcv.chapter_id")
	buf.WriteString(" order by qcv.chapter_id  desc,qcv.sort desc,qcv.id asc ")

	result, err := db.GetCourse().Query(buf.String(), args...)
	if err != nil {
		return nil, err
	}
	resultMap := make(map[int64][]Video, 0)

	var chapterId, id, videoId, duration sql.NullInt64
	var cover, coverPrefix, url, urlPrefix, down, downPrefix, title sql.NullString
	for result.Next() {
		err := result.Scan(&chapterId, &id, &cover, &coverPrefix, &url, &urlPrefix, &down, &downPrefix, &duration, &title, &videoId)
		if err != nil {
			return nil, err
		}
		list, ok := resultMap[chapterId.Int64]
		if !ok {
			list = make([]Video, 0)
		}
		list = append(list, Video{
			ID:       id.Int64,
			VideoId:  videoId.Int64,
			Cover:    coverPrefix.String + cover.String,
			Url:      urlPrefix.String + url.String,
			DownUrl:  downPrefix.String + down.String,
			Title:    title.String,
			Duration: duration.Int64,
		})
		resultMap[chapterId.Int64] = list
	}
	return resultMap, nil
}

func EditItem(id, chapterId int64, title string, tagId []int64) (int64, error) {
	err := db.GetVideo().WithTransaction(func(tx *db.Tx) error {
		if id == 0 {
			chapterSql := "insert into qz_chapter_video (chapter_id,title,video_id,content,sort) values (?,?,?,?,?)"
			chapterVideoResult, err := tx.Exec(chapterSql, chapterId, title, 0, "", enum.SortMaxValue)
			if err != nil {
				return xlog.Error("生成视频数据出错")
			}
			id, err = chapterVideoResult.LastInsertId()
			if err != nil {
				return xlog.Error("生成视频数据出错")
			}

		} else {
			chapterSql := "update  qz_chapter_video set title=? where id=?"
			_, err := tx.Exec(chapterSql, title, id)
			if err != nil {
				return xlog.Error("生成视频数据出错")
			}
		}
		return tagMoel.EditTagRecord(tagId, id, enum.VideoTagType, tx)
	})
	return id, err
}

func EditItemInfo(cover, prefix string, id, typeInt int64) error {
	if typeInt == enum.EDIT_COVER {
		_, err := db.GetCourse().Exec("update qz_chapter_video set `cover`=?,`prefix`=? where id=?", cover, prefix, id)
		return err
	} else if typeInt == enum.EDIT_CONTENT {
		_, err := db.GetCourse().Exec("update qz_chapter_video set `content`=? where id=?", cover, id)
		return err
	}
	return nil
}

func GetAdminItem(idList []int64) ([]*Video, error) {
	buf := strings.Builder{}
	buf.WriteString("select qcv.id,qcv.cover,qcv.prefix,qv.url,qv.url_prefix,qcv.title,qcv.content " +
		"from qz_chapter_video qcv left join qz_video qv on qcv.video_id=qv.id where qcv.`status`=? ")
	utils.MysqlStringInUtils(&buf, idList, "and qcv.id")
	buf.WriteString(" group by qcv.id")
	var id sql.NullInt64
	var cover, coverPrefix, url, urlPrefix, title, content sql.NullString
	result, err := db.GetVideo().Query(buf.String(), enum.StatusNormal)
	list := make([]*Video, 0)
	resultMap:=make(map[int64]*Video)
	if err != nil {
		return nil, err
	}
	for result.Next() {
		err = result.Scan(&id, &cover, &coverPrefix, &url, &urlPrefix, &title, &content)
		if err == nil {
			temp:= &Video{
				ID:      id.Int64,
				Title:   title.String,
				Cover:   coverPrefix.String + cover.String,
				Url:     urlPrefix.String + url.String,
				Content: content.String,
				TagList: make([]tagMoel.Tag, 0),
			}
			resultMap[id.Int64]=temp
			list = append(list,temp)
		}

	}

	tagList, err := tagMoel.GetTagRecord(enum.VideoTagType, idList)
	if err != nil {
		return nil, err
	}
	for _,v:=range tagList{
		temp,ok:=resultMap[v.ContentId]
		if ok{
			temp.TagList=append(temp.TagList,v)
		}
	}
	return list, nil
}

func AddVideo(list []AddVideoData) error {
	sqlStr := "insert into qz_video (cover,cover_prefix,url,url_prefix,down_url,down_prefix,duration,status) values"
	return db.GetVideo().WithTransaction(func(tx *db.Tx) error {
		for _, v := range list {
			var id int64
			args := make([]interface{}, 0)
			args = append(args, v.Cover)
			args = append(args, v.CoverPrefix)
			args = append(args, v.Url)
			args = append(args, v.UrlPrefix)
			args = append(args, v.DownUrl)
			args = append(args, v.DownUrlPrefix)
			args = append(args, v.Duration)
			result, err := tx.Exec(sqlStr+"(?,?,?,?,?,?,?,1)", args...)
			if err != nil {
				return xlog.Error(" 生成视频数据出错")
			}
			id, err = result.LastInsertId()
			if err != nil {
				return xlog.Error("生成视频数据出错")
			}
			chapterSql := "update  qz_chapter_video set video_id=? where id=?"
			_, err = tx.Exec(chapterSql, id, v.ID)
			if err != nil {
				return xlog.Error("生成视频数据出错")
			}
			return nil
		}
		return nil
	})
}

func GetVideoList(page, size int64) ([]Video, error) {
	var start = (page - 1) * size
	result, err := db.GetVideo().Query("select qv.id,qv.cover,qv.cover_prefix,qv.url,qv.url_prefix,qv.down_url,qv.down_prefix,"+
		"qv.duration from qz_video qv  where qv.`status`=1 limit ?,?", start, size)
	if err != nil {
		return nil, err
	}
	list := make([]Video, 0)
	var id, duration sql.NullInt64
	var cover, coverPrefix, url, urlPrefix, down, downPrefix sql.NullString
	for result.Next() {
		err := result.Scan(&id, &cover, &coverPrefix, &url, &urlPrefix, &down, &downPrefix, &duration)
		if err != nil {
			return nil, err
		}
		list = append(list, Video{
			ID:       id.Int64,
			Cover:    coverPrefix.String + cover.String,
			Url:      urlPrefix.String + url.String,
			DownUrl:  downPrefix.String + down.String,
			Duration: duration.Int64,
		})
	}
	return list, nil
}

func GetHistoryVideoList(uid int64) ([]Video, error) {
	buf := strings.Builder{}
	args := make([]interface{}, 0)
	buf.WriteString("select qcv.id,qv.cover,qv.cover_prefix,qv.url,qv.url_prefix,qv.down_url,qv.down_prefix," +
		"qv.duration,qcv.title from qz_video qv inner join qz_chapter_video qcv on qcv.video_id=qv.id inner join qz_history qh on qh.chapter_video_id=qcv.id" +
		" where qv.`status`=1 and qh.uid=? group by qcv.id")
	args = append(args, uid)
	result, err := db.GetCourse().Query(buf.String(), args...)
	if err != nil {
		return nil, err
	}
	list := make([]Video, 0)
	var id, duration sql.NullInt64
	var cover, coverPrefix, url, urlPrefix, down, downPrefix, title sql.NullString
	for result.Next() {
		err := result.Scan(&id, &cover, &coverPrefix, &url, &urlPrefix, &down, &downPrefix, &duration, &title)
		if err != nil {
			return nil, err
		}
		list = append(list, Video{
			ID:       id.Int64,
			Cover:    coverPrefix.String + cover.String,
			Url:      urlPrefix.String + url.String,
			DownUrl:  downPrefix.String + down.String,
			Title:    title.String,
			Duration: duration.Int64,
		})
	}
	return list, nil
}

func AddHistoryVideo(id, uid int64) error {
	_, err := db.GetVideo().Exec("insert into qz_history (chapter_video_id,uid,create_time) values (?,?,?)", id, uid, time.Now().Unix())
	return err
}

func DelHistoryVideo(uid int64, idList []int64) error {
	return db.GetVideo().WithTransaction(func(tx *db.Tx) error {
		for _, v := range idList {
			_, err := tx.Exec("delete from qz_history where chapter_video_id=? and uid=?", v, uid)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func DelChapterVideo(id int64) error {
	_, err := db.GetVideo().Exec("update  qz_chapter_video set status=0 where id=?", id)
	return err
}

func GetVideo(videoId []int64, queryTagId int64,key string) ([]*Video, error) {
	buf := strings.Builder{}
	args := make([]interface{}, 0)
	buf.WriteString(`SELECT qcv.id,qcv.chapter_id,qcv.cover,qcv.prefix,qv.url,qv.url_prefix,qv.down_url,qv.down_prefix,qcv.content,qcv.title
	FROM  qz_chapter_video qcv
	left JOIN qz_video qv ON qcv.video_id=qv.id
	`)
	if queryTagId != 0 {
		buf.WriteString(" inner join qz_tag_record qtr on qtr.content_id=qcv.id and qtr.type=3 and qtr.tag_id=? ")
		args = append(args, queryTagId)
	}

	buf.WriteString(" WHERE qv.`status`=?")
	args = append(args, enum.StatusNormal)
	if key!=""{
		buf.WriteString(" and qcv.title like (?)")
		args = append(args, "%"+key+"%")
	}

	utils.MysqlStringInUtils(&buf, videoId, " and qcv.id")

	buf.WriteString(" order by qcv.sort desc,qcv.id desc ")

	result, err := db.GetCourse().Query(buf.String(), args...)
	if err != nil {
		return nil, err
	}
	list := make([]*Video, 0)
	idList:=make([]int64,0)
	mapList:=make(map[int64]*Video)
	var id, chapterId sql.NullInt64
	var cover, coverPrefix, url, urlPrefix, down, content, downPrefix, title sql.NullString
	for result.Next() {
		err := result.Scan(&id, &chapterId,&cover, &coverPrefix, &url, &urlPrefix, &down, &downPrefix, &content,  &title)
		if err != nil {
			return nil, err
		}
		idList=append(idList,id.Int64)
		v:=&Video{
			ID:       id.Int64,
			Cover:    coverPrefix.String + cover.String,
			Url:      urlPrefix.String + url.String,
			DownUrl:  downPrefix.String + down.String,
			Title:    title.String,
			Content:  content.String,
			ChapterId:chapterId.Int64,
		}
		mapList[id.Int64]=v
		list = append(list, v)
	}
	resultTag,err:=tagMoel.GetTagRecord(enum.VideoTagType,idList)
	if err==nil{
		for _,v:=range resultTag{
			temp,ok:=mapList[v.ContentId]
			if ok{
				temp.TagList=append(temp.TagList,v)
			}
		}
	}
	return list, nil
}


//修改排序
func EditVideoSort(videoId, sort int64) error {
	//列表页面修改
	err := db.GetCourse().WithTransaction(func(tx *db.Tx) error {
		var chapterId sql.NullInt64
		err := tx.QueryRow("select chapter_id from qz_chapter_video where id=?", videoId).Scan(&chapterId)
		if err != nil {
			return xlog.Error(err)
		}
		if sort == enum.SortMax {
			_, err = tx.Exec(`update qz_chapter_video tb set sort=sort-1 where tb.id<>? and tb.chapter_id=?`, videoId, chapterId)
			if err != nil {
				return xlog.Error(err)
			}
			_, err = tx.Exec(`update qz_chapter_video set sort=? where id=?`, enum.SortMaxValue, videoId)
			if err != nil {
				return xlog.Error(err)
			}
		} else {
			var sSort int64
			err = tx.QueryRow(`select sort from qz_chapter_video where id=?`, videoId).Scan(&sSort)
			if err != nil {
				return xlog.Error(err)
			}
			//获取要交换的排序
			var nextVideoId, nextSort int64
			if sort == enum.SortDown {
				err = tx.QueryRow(`select id,sort from qz_chapter_video where id<>? and chapter_id=? and sort<=?  order by sort desc,id asc`, videoId, chapterId, sSort).Scan(&nextVideoId, &nextSort)
			} else {
				err = tx.QueryRow(`select id,sort from qz_chapter_video where id<>? and chapter_id=? and sort>=?  order by sort desc,id asc`, videoId, chapterId, sSort).Scan(&nextVideoId, &nextSort)
			}

			if err != nil {
				return xlog.Error(err)
			}
			//交换
			_, err = tx.Exec(`update qz_chapter_video set sort=? where id=?`, nextSort, videoId)
			if err != nil {
				return xlog.Error(err)
			}
			_, err = tx.Exec(`update qz_chapter_video set sort=? where id=?`, sSort, nextVideoId)
			if err != nil {
				return xlog.Error(err)
			}
		}

		return nil
	})
	return err
}

func EditVideo(cover, prefix string, id int64) error {
	return db.GetCourse().WithTransaction(func(tx *db.Tx) error {
		_, err := tx.Exec("update qz_video set `cover`=?,`cover_prefix`=? where id=?", cover, prefix, id)
		return err

	})
}

func GetVideoIndexInfo(courseId, videoId int64) ProgressInfo {
	info := ProgressInfo{
		ChapterIndex: 1,
		VideoIndex:   1,
		Progress:     0,
	}
	if videoId == 0 || courseId == 0 {
		return info
	}
	result, err := db.GetCourse().Query(`select tb.id,tb1.id from qz_chapter tb
inner join qz_chapter_video tb1 on tb1.chapter_id=tb.id
 where tb.course_id=? and tb.status=1 and tb1.status=1 ORDER BY tb.sort desc,tb.id asc,tb1.sort desc,tb1.id asc`, courseId)

	if err == nil {
		var chapterId, id sql.NullInt64
		list := make([]ProgressInfoDetail, 0)
		for result.Next() {
			err = result.Scan(&chapterId, &id)
			if err == nil {
				list = append(list, ProgressInfoDetail{
					ChapterId: chapterId.Int64,
					Id:        id.Int64,
				})
			}
		}
		info.List = list
		index := 0
		videoIndex := int64(0)
		chapterIndex := int64(0)
		chapter := int64(0)
		for _, v := range list {
			if chapter != v.ChapterId {
				chapterIndex++
				chapter = v.ChapterId
				videoIndex = 1
			} else {
				videoIndex++
			}
			if v.Id == videoId {
				value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(index)/float64(len(list))), 64)
				info.Progress = value
				break
			} else {
				index++
			}
		}
		info.ChapterIndex = chapterIndex
		info.VideoIndex = videoIndex
	}
	return info
}
