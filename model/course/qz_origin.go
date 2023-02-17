package courseModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/xlog"
	"strings"
)

type Origin struct {
	Id          int64
	Photo       string
	Title       string
	Desc        string
	Name        string
	Info        string
	InfoTitle   string
	Certificate string
	Auth        []string
}

type Auth struct {
	Cover  string
	Prefix string
}

func GetOrigin(id int64) (*Origin, error) {
	result, err := db.GetCourse().Query(`select tb.id,tb.name,tb.title,tb.desc,tb.info,tb.prefix,tb.cover,tb1.prefix,tb1.cover,tb.info_title,tb.certificate from qz_origin tb 
inner join qz_course tb2 on tb2.origin_id=tb.id
left join qz_origin_auth_photo tb1 on tb1.oid=tb.id  and tb1.status=?
where tb2.id=?`, enum.StatusNormal, id)
	if err != nil {
		return nil, xlog.Error(err)
	}
	temp := &Origin{
		Auth: make([]string, 0),
	}
	var name, info, title, desc, oPrefix, oCover, prefix, cover,infoTitle,certificate sql.NullString
	var tId sql.NullInt64
	for result.Next() {
		err := result.Scan(&tId, &name, &title, &desc, &info, &oPrefix, &oCover, &prefix, &cover,&infoTitle,&certificate)
		if err != nil {
			return nil, err
		}
		temp.Title = title.String
		temp.Id = tId.Int64
		temp.Name = name.String
		temp.Desc = desc.String
		temp.Info = info.String
		temp.InfoTitle=infoTitle.String
		temp.Certificate=certificate.String
		temp.Photo = oPrefix.String + oCover.String
		if prefix.String != "" {
			temp.Auth = append(temp.Auth, prefix.String+cover.String)
		}
	}
	return temp, nil
}

func EditOrigin(name, title, desc, info, infoTitle, certificate string, id int64, tx *db.Tx) (int64, error) {
	if id == 0 {
		idResult, err := tx.Exec("insert into qz_origin (`name`,title,`desc`,`info`,`info_title`,`certificate`) value(?,?,?,?,?,?)", name, title, desc, info, infoTitle, certificate)
		if err == nil {
			id, err = idResult.LastInsertId()
			return id, err
		} else {
			return 0, xlog.Error(err)
		}
	} else {
		_, err := tx.Exec("update qz_origin set `name`=?,`title`=?,`desc`=?,`info`=?,`info_title`=?,`certificate`=? where id=?", name, title, desc, info, infoTitle, certificate, id)
		if err == nil {
			return id, err
		} else {
			return 0, xlog.Error(err)
		}
	}
}

func EditOriginCourseCover(id int64, cover, prefix string) error {
	_, err := db.GetCourse().Exec("update qz_origin set `cover`=?,`prefix`=? where id=?", cover, prefix, id)
	return err
}

func EditOriginAuth(originId int64, list []Auth) error {
	return db.GetCourse().WithTransaction(func(tx *db.Tx) error {
		_, err := tx.Exec("update qz_origin_auth_photo set `status`=? where oid=?", enum.StatusDelete, originId)

		if err == nil {
			args := make([]interface{}, 0)
			sqlStr := "insert into qz_origin_auth_photo (`cover`,`prefix`,`oid`) values"
			placeHolder := "(?,?,?)"
			values := make([]string, 0)
			for _, v := range list {
				values = append(values, placeHolder)
				args = append(args, v.Cover)
				args = append(args, v.Prefix)
				args = append(args, originId)
			}
			sqlStr += strings.Join(values, ",")
			_, err := tx.Exec(sqlStr, args...)
			if err != nil {
				return xlog.Error("生成课程数据出错")
			}
			return nil
		} else {
			return xlog.Error(err)
		}
	})
}
