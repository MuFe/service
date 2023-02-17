package courseModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/utils"
	"strings"
)

type Recommend struct {
	ID          int64
	InfoId      int64
	ContentType int64
	ContentId   int64
}

type RecommendInfo struct {
	ID   int64
	Type int64
	Icon string
	Name string
}

func GetAdminRecommendList(infoId, page, size int64) ([]Recommend, error) {
	start := (page - 1) * size
	result, err := db.GetCourse().Query("select id,content_type,content_id from qz_recommend where `info_id`=? limit ?,? ", infoId, start, size)
	if err != nil {
		return nil, err
	}
	list := make([]Recommend, 0)
	var contentId, id, contentType sql.NullInt64
	for result.Next() {
		err := result.Scan(&id, &contentType, &contentId)
		if err != nil {
			return nil, err
		}
		list = append(list, Recommend{
			ID:          id.Int64,
			ContentId:   contentId.Int64,
			ContentType: contentType.Int64,
			InfoId:infoId,
		})
	}
	return list, nil
}

func GetRecommendList(infoList []int64) ([]Recommend, error) {
	var buf strings.Builder
	buf.WriteString("select content_type,content_id,info_id from qz_recommend where 1=1 ")
	utils.MysqlStringInUtils(&buf, infoList, " and info_id")
	result, err := db.GetCourse().Query(buf.String())
	if err != nil {
		return nil, err
	}
	list := make([]Recommend, 0)
	var contentId, infoId, contentType sql.NullInt64
	for result.Next() {
		err := result.Scan(&contentType, &contentId, &infoId)
		if err != nil {
			return nil, err
		}
		list = append(list, Recommend{
			InfoId:      infoId.Int64,
			ContentId:   contentId.Int64,
			ContentType: contentType.Int64,
		})
	}
	return list, nil
}

func GetRecommendInfoList(typeInt int64)([]RecommendInfo,error) {
	result, err := db.GetCourse().Query("select id,name,icon,type from qz_recommend_info where `identity`=?", typeInt)
	list:=make([]RecommendInfo,0)
	if err == nil {
		var id,idType sql.NullInt64
		var name,icon sql.NullString
		for result.Next(){
			err=result.Scan(&id,&name,&icon,&idType)
			if err==nil{
				list=append(list,RecommendInfo{
					ID:id.Int64,
					Type:idType.Int64,
					Name:name.String,
					Icon:icon.String,
				})
			}
		}
	}
	return list,nil
}

func EditRecommend(id, typeInt, contentId, contentType int64, del bool) error {
	if del {
		_, err := db.GetCourse().Exec("delete from qz_recommend where id=?", id)
		return err
	} else {
		return db.GetCourse().WithTransaction(func(tx *db.Tx) error {
			_,err:=tx.Exec("delete from qz_recommend where `info_id`=? and content_type=? and content_id=?",typeInt, contentType, contentId)
			if err!=nil{
				return err
			}
			_, err = tx.Exec("insert into qz_recommend (`info_id`,`content_type`,`content_id`) values(?,?,?)", typeInt, contentType, contentId)
			return err
		})
	}
}
