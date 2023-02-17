package bannerModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"strings"
)

type Banner struct {
	Photo     string
	TypeInt   int64
	Url       string
	Prefix    string
	ContentId int64
	Id int64
	Sort int64
}

func GetBanner(status,nowId int64) ([]Banner, error) {
	list:=make([]Banner,0)
	var buf strings.Builder
	args:=make([]interface{},0)
	buf.WriteString("select id,url,type,photo,content_id,prefix,sort from qz_banner where 1=1 ")
	if status!=enum.StatusAll{
		buf.WriteString(" and status=?")
		args=append(args,status)
	} else {
		buf.WriteString(" and status<>?")
		args=append(args,enum.StatusDelete)
	}

	if nowId!=0{
		buf.WriteString(" and id=?")
		args=append(args,nowId)
	}
	buf.WriteString(" order by sort desc,id desc")
	rows,err:=db.GetBannerDb().Query(buf.String(),args...)
	if err == nil {
		var contentId,id,sort, typeInt sql.NullInt64
		var bannerUrl, photo, prefix sql.NullString
		for rows.Next() {
			err = rows.Scan(&id,&bannerUrl, &typeInt, &photo, &contentId, &prefix,&sort)
			if err == nil {
				var result Banner
				result.Photo = photo.String
				result.Id=id.Int64
				result.TypeInt = typeInt.Int64
				result.Url = bannerUrl.String
				result.Prefix = prefix.String
				result.ContentId = contentId.Int64
				result.Sort=sort.Int64
				list = append(list, result)
			}
		}
	}
	return list, nil
}

func AddAd(typeInt,contentId int64,url string)(int64,error){
	result,err:=db.GetBannerDb().Exec("insert into  qz_banner (`type`,`status`,content_id,url) values (?,?,?,?)",typeInt,enum.StatusVerify,contentId,url)
	if err!=nil{
		return 0,err
	}
	id,err:=result.LastInsertId()
	return id,err
}
func DelAd(id int64)error{
	_,err:=db.GetBannerDb().Exec("update qz_banner set `status`=?  where id=?",enum.StatusDelete,id)
	return err
}
func EditAd(typeInt,contentId,id int64,url string)error{
	_,err:=db.GetBannerDb().Exec("update qz_banner set `type`=?,url=?,content_id=? where id=?",typeInt,url,contentId,id)
	return err
}

func EditSort(sort,id int64)error{
	_,err:=db.GetBannerDb().Exec("update qz_banner set `sort`=? where id=?",sort,id)
	return err
}

func EditPhoto(id int64,photo,prefix string)error{
	_,err:=db.GetBannerDb().Exec("update qz_banner set status=?,photo=?,`prefix`=? where id=?",enum.StatusNormal,photo,prefix,id)
	return err
}
