package coachModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/xlog"
	"strings"
	"time"
)

type InstitutionCourse struct {
	Id     int64
	School int64
	Name   string
	Level  string
	Price  int64
	Max    int64
	Duration int64
}

func GetInstitutionCourse(status, nowId, parentId int64) ([]InstitutionCourse, error) {
	list := make([]InstitutionCourse, 0)
	var buf strings.Builder
	args := make([]interface{}, 0)
	buf.WriteString("select id,name,`level`,price,max,`duration` from qz_institution_course where 1=1 ")
	if status != enum.StatusAll {
		buf.WriteString(" and status=?")
		args = append(args, status)
	} else {
		buf.WriteString(" and status<>?")
		args = append(args, enum.StatusDelete)
	}

	if nowId != 0 {
		buf.WriteString(" and id=?")
		args = append(args, nowId)
	}

	if parentId != 0 {
		buf.WriteString(" and parent_id=?")
		args = append(args, parentId)
	}
	buf.WriteString(" order by id desc")
	rows, err := db.GetCoachDb().Query(buf.String(), args...)
	if err == nil {
		var id, price, max,duration sql.NullInt64
		var name, level sql.NullString
		for rows.Next() {
			err = rows.Scan(&id, &name, &level, &price, &max,&duration)
			if err == nil {
				list = append(list, InstitutionCourse{
					Id:    id.Int64,
					Name:  name.String,
					Price: price.Int64,
					Level: level.String,
					Max:   max.Int64,
					Duration:duration.Int64,
				})
			}
		}
	}
	return list, nil
}

func EditInstitutionCourse(createBy, id, max, price, parentId,duration int64, name, level string, del bool) error {
	return db.GetCoachDb().WithTransaction(func(tx *db.Tx) error {
		now := time.Now().Unix()
		if id == 0 {
			idResult, err := tx.Exec("insert into qz_institution_course (name,`level`,create_time,modify_time,create_by,max,price,parent_id,`duration`) values(?,?,?,?,?,?,?,?,?)", name, level, now, now, createBy, max, price, parentId,duration)
			if err != nil {
				return xlog.Error("创建课程失败")
			}
			id, _ = idResult.LastInsertId()
		} else if del {
			_, err := tx.Exec("update qz_institution_course set `status`=?,modify_time=? where id=?", enum.StatusDelete, now, id)
			if err != nil {
				return xlog.Error("修改课程失败")
			}
		} else {
			_, err := tx.Exec("update qz_institution_course set `level`=?,name=?,max=?,price=?,`duration`=?,modify_time=? where id=?", level, name, max, price,duration, now, id)
			if err != nil {
				return xlog.Error("修改课程失败")
			}
		}
		return nil
	})
}
