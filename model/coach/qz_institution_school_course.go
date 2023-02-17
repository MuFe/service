package coachModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"strings"
)



func GetInstitutionSchoolCourse(schoolId []int64) ([]InstitutionCourse, error) {
	list := make([]InstitutionCourse, 0)
	var buf strings.Builder
	args := make([]interface{}, 0)
	buf.WriteString("select tb.course_id,tb1.name,tb.school_id,tb1.`level`,tb1.price,tb1.max,tb1.`duration` from qz_institution_school_course tb inner join qz_institution_course tb1 on tb1.id=tb.course_id where 1=1 ")
	utils.MysqlStringInUtils(&buf,schoolId," and tb.school_id")
	rows, err := db.GetCoachDb().Query(buf.String(), args...)
	if err == nil {
		var id,school,price, max,duration sql.NullInt64
		var name,level sql.NullString
		for rows.Next() {
			err = rows.Scan(&id, &name,&school,&level,&price,&max,&duration)
			if err == nil {
				list = append(list, InstitutionCourse{
					Id:id.Int64,
					Name:name.String,
					School:school.Int64,
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


func EditInstitutionSchoolCourse(tx *db.Tx,schoolId int64,courseId []int64) error {
	_,err:=tx.Exec("delete from qz_institution_school_course where school_id=?",schoolId)
	if err!=nil{
		return xlog.Error(err)
	}
	if len(courseId)==0{
		return nil
	}
	sqlStr := "insert into qz_institution_school_course (`school_id`,`course_id`) values"
	args := make([]interface{}, 0)
	placeHolder := "(?,?)"
	values := make([]string, 0)
	for _,v:=range courseId{
		values = append(values, placeHolder)
		args = append(args, schoolId)
		args = append(args, v)
	}

	sqlStr += strings.Join(values, ",")
	_, err = tx.Exec(sqlStr, args...)
	if err != nil {
		return xlog.Error("生成数据出错")
	}
	return nil
}

