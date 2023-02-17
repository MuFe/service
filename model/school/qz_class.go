package schoolModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"strings"
	"time"
)

type Class struct {
	Id         int64
	Name       string
	Number     int64
	CreateTime int64
	Tag        string
	GradeId    int64
	Grade      string
	SchoolId   int64
	SchoolType string
	SchoolName string
	SchoolIcon string
	Code       string
	CourseId    int64
	VideoId    int64
	CreateBy    int64
	AdminList  []int64
	List       []int64
}

func CreateClassInfo(uid, schoolId, gradeId int64, name []string) error {
	sqlStr:="insert into qz_class_info (name,grade_id,school_id,create_time,tag) values "
	placeHolder := "(?,?,?,?,?)"
	values := make([]string, 0)
	args := make([]interface{}, 0)
	for _,v:=range name{
		values = append(values, placeHolder)
		args = append(args, v)
		args = append(args, gradeId)
		args = append(args, schoolId)
		args = append(args, time.Now().Unix(), enum.FootBallType)
	}

	sqlStr += strings.Join(values, ",")
	_, err := db.GetSchool().Exec(sqlStr, args...)
	return err
}

func CreateClass(uid, classInfoId int64,schoolIdList []int64) error {
	return db.GetSchool().WithTransaction(func(tx *db.Tx) error {
		var schoolId sql.NullInt64
		err:=tx.QueryRow("select school_id from qz_class_info where id=?",classInfoId).Scan(&schoolId)
		if err==sql.ErrNoRows{
			return errcode.HttpErrorWringParam.RPCError()
		}else
		if err != nil {
			return xlog.Error(err)
		}
		isHave:=false
		for _,v:=range schoolIdList{
			if v==schoolId.Int64{
				isHave=true
			}
		}
		if !isHave{
			return errcode.CommErrorUnauthorized.RPCError()
		}
		var count sql.NullInt64
		err=tx.QueryRow("select count(id) from qz_class where class_info_id=? and `status`=?",classInfoId,enum.StatusNormal).Scan(&count)
		if err != nil {
			return xlog.Error(err)
		}
		if count.Int64>0{
			return xlog.Error("您已创建过该班级了")
		}
		idResult, err := tx.Exec("insert into qz_class (class_info_id,create_time,create_by) values(?,?,?)", classInfoId, time.Now().Unix(),uid)
		if err != nil {
			return xlog.Error(err)
		}
		id, err := idResult.LastInsertId()
		if err != nil {
			return xlog.Error(err)
		}
		return JoinClass(uid, id, enum.ClassAdminType, tx)
	})
}

func JoinClass(uid, classId, typeInt int64, tx *db.Tx) error {
	if tx == nil {
		return db.GetSchool().WithTransaction(func(tx *db.Tx) error {
			return joinClass(uid, classId, typeInt, tx)
		})
	} else {
		return joinClass(uid, classId, typeInt, tx)
	}
}

func joinClass(uid, classId, typeInt int64, tx *db.Tx) error {
	var count sql.NullInt64
	err := tx.QueryRow("select count(id) from qz_class_record where uid=? and class_id=? and `type`=?", uid, classId, typeInt).Scan(&count)
	if err != nil {
		return xlog.Error("加入班级失败")
	}
	if count.Int64 > 0 {
		return xlog.Error("已经加入该班级了")
	}
	inviteResult, err := tx.Query("select code from qz_class_record where code<>''")
	if err != nil {
		return xlog.Error("加入班级失败")
	}
	inviteMap := make(map[string]bool)
	var invite sql.NullString
	for inviteResult.Next() {
		err = inviteResult.Scan(&invite)
		if err == nil {
			inviteMap[invite.String] = true
		}
	}
	code := utils.BaseCode6(inviteMap)
	_, err = tx.Exec("insert into qz_class_record (uid,class_id,`type`,`code`,`prefix`) values (?,?,?,?,?)", uid, classId, typeInt, code, enum.ClassCodePrefix)
	return err
}

func DisClass(uid, classId int64) ([]int64,error) {
	list:=make([]int64,0)
	err:= db.GetSchool().WithTransaction(func(tx *db.Tx) error {
		var classUid sql.NullInt64

		err := tx.QueryRow("select create_by from qz_class where id=?", classId).Scan(&classUid)
		if err != nil {
			return xlog.Error(err)
		}
		if classUid.Int64!=uid{
			return xlog.Error(errcode.CommErrorUnauthorized.Msg)
		}
		_, err = tx.Exec("update qz_class set `status`=? where id=?", enum.StatusDelete, classId)
		if err != nil {
			return xlog.Error(err)
		}
		result,err:=tx.Query("select uid from qz_class_record where class_id=? and `type`=?",classId,enum.ClassStudentType)
		if err==nil{
			var tUid sql.NullInt64
			for result.Next(){
				err=result.Scan(&tUid)
				if err==nil{
					list=append(list,tUid.Int64)
				}
			}
		}else{
			return xlog.Error(err)
		}
		return nil
	})
	return list,err
}

func QuitClass(uid, classId,quitUid int64) error {
	return db.GetSchool().WithTransaction(func(tx *db.Tx) error {
		if uid!=quitUid{
			var count sql.NullInt64
			err := tx.QueryRow("select count(id) from qz_class_record where uid=? and class_id=? and type=?", uid, classId, enum.ClassAdminType).Scan(&count)
			if err != nil {
				return xlog.Error(err)
			}
			if count.Int64 == 0 {
				return errcode.CommErrorUnauthorized.RPCError()
			}
		}
		err:=DeleteClassRecord(classId,quitUid,tx)
		if err != nil {
			return xlog.Error(err)
		}
		return nil
	})
}

func DeleteClassRecord(classId,quitUid int64,tx *db.Tx)error{
	var err error
	if classId==0{
		_, err=tx.Exec("delete from qz_class_record where uid=?",quitUid)
	} else{
		_, err=tx.Exec("delete from qz_class_record where uid=? and class_id=?",quitUid,classId)
	}

	return err
}

func CancelClassRecord(quitUid []int64)error{
	var err error
	if len(quitUid)==0{
		return nil
	}
	var buf strings.Builder
	buf.WriteString("delete from qz_class_record ")
	utils.MysqlStringInUtils(&buf,quitUid," where uid")
	_, err=db.GetSchool().Exec(buf.String())

	return err
}

func AdminClassList(schoolId, gradeId int64) ([]*Class, error) {
	var startBuf strings.Builder
	args:=make([]interface{},0)
	startBuf.WriteString(`select tb.id,tb1.name from qz_class_info tb1
left join qz_class tb on tb1.id=tb.class_info_id and tb.status=?
 where tb1.status=? `)
	args=append(args,enum.StatusNormal,enum.StatusNormal,)
	if schoolId!=0{
		startBuf.WriteString(" and tb1.school_id=? ")
		args=append(args,schoolId)
	}
	if gradeId!=0{
		startBuf.WriteString(" and tb1.grade_id=? ")
		args=append(args,gradeId)
	}
	result, err := db.GetSchool().Query(startBuf.String(),args...)

	if err != nil {
		return nil, xlog.Error(err)
	}
	var classId sql.NullInt64
	var  className sql.NullString
	list := make([]*Class, 0)
	listMap := make(map[int64]*Class)
	classIdList := make([]int64, 0)

	for result.Next() {
		err = result.Scan(&classId, &className)
		if err == nil {
			temp := &Class{
				Id:         classId.Int64,
				Name:       className.String,
			}
			classIdList = append(classIdList, classId.Int64)
			list = append(list, temp)
			listMap[classId.Int64] = temp
		}
	}
	var buf strings.Builder
	buf.WriteString("select count(id),class_id from qz_class_record where `type`=?")
	utils.MysqlStringInUtils(&buf, classIdList, " and class_id")
	buf.WriteString(" group by class_id")
	countResult, err := db.GetSchool().Query(buf.String(), enum.ClassStudentType)
	if err == nil {
		var count, classId sql.NullInt64
		for countResult.Next() {
			err = countResult.Scan(&count, &classId)
			if err == nil {
				info, ok := listMap[classId.Int64]
				if ok {
					info.Number = count.Int64
				}
			}
		}
	}
	return list, nil
}

func ClassList(uid, id int64) ([]*Class, error) {
	var bufStr strings.Builder
	args := make([]interface{}, 0)
	bufStr.WriteString(`select tb.id,tb1.name,tb2.name,tb4.name,tb.create_time,tb1.tag,CONCAT(tb5.prefix,tb5.code),tb3.name,tb3.icon,tb.create_by from qz_class tb
inner join qz_class_info tb1 on tb1.id=tb.class_info_id
inner join qz_grade tb2 on tb2.id=tb1.grade_id
inner join qz_school tb3 on tb3.id=tb1.school_id
inner join qz_school_type tb4 on tb4.id=tb3.school_type
inner join qz_class_record tb5 on tb5.class_id=tb.id`)
	bufStr.WriteString(" where tb.status=? and tb5.uid=? ")
	args = append(args, enum.StatusNormal, uid)

	if id != 0 {
		bufStr.WriteString(" and tb.id=?")
		args = append(args, id)
	}

	result, err := db.GetSchool().Query(bufStr.String(), args...)

	if err != nil {
		return nil, xlog.Error(err)
	}
	var classId, createTime,createBy sql.NullInt64
	var gradeName, schoolType, className,schoolName,schoolIcon, tag, code sql.NullString
	list := make([]*Class, 0)
	listMap := make(map[int64]*Class)
	classIdList := make([]int64, 0)

	for result.Next() {
		err = result.Scan(&classId, &className, &gradeName, &schoolType, &createTime, &tag, &code,&schoolName,&schoolIcon,&createBy)
		if err == nil {
			temp := &Class{
				Id:         classId.Int64,
				Name:       className.String,
				Grade:      gradeName.String,
				SchoolType: schoolType.String,
				Tag:        tag.String,
				CreateTime: createTime.Int64,
				Code:       code.String,
				SchoolName:schoolName.String,
				SchoolIcon:schoolIcon.String,
				CreateBy:createBy.Int64,
			}
			classIdList = append(classIdList, classId.Int64)
			list = append(list, temp)
			listMap[classId.Int64] = temp
		}
	}
	var buf strings.Builder
	buf.WriteString("select count(id),class_id from qz_class_record where `type`=?")
	utils.MysqlStringInUtils(&buf, classIdList, " and class_id")
	buf.WriteString(" group by class_id")
	countResult, err := db.GetSchool().Query(buf.String(), enum.ClassStudentType)
	if err == nil {
		var count, classId sql.NullInt64
		for countResult.Next() {
			err = countResult.Scan(&count, &classId)
			if err == nil {
				info, ok := listMap[classId.Int64]
				if ok {
					info.Number = count.Int64
				}
			}
		}
	}
	return list, nil
}

func ClassDetail(uid, id int64) (*Class, error) {
	result := &Class{
		AdminList:  make([]int64, 0),
		List:       make([]int64, 0),
	}
	var buf strings.Builder
	var count sql.NullInt64
	isHave:=false
	var gradeName, schoolType, className, tag, code sql.NullString
	buf.WriteString("select uid,`type`,CONCAT(prefix,`code`) from qz_class_record where  class_id=?")
	listResult, err := db.GetSchool().Query(buf.String(), id)
	if err == nil {
		var uidResult, typeResult sql.NullInt64
		for listResult.Next() {
			err = listResult.Scan(&uidResult, &typeResult, &code)
			if err == nil {
				if uidResult.Int64 == uid {
					result.Code = code.String
					isHave=true
				}
				if typeResult.Int64 == enum.ClassAdminType {
					result.AdminList = append(result.AdminList, uidResult.Int64)
				}
				if typeResult.Int64 == enum.ClassStudentType {
					result.List = append(result.List, uidResult.Int64)
				}
			}
		}
	}
	if !isHave {
		return nil, xlog.Error(errcode.HttpErrorWringParam.Msg)
	}
	if err == nil {
		result.Number = count.Int64
	}

	var bufStr strings.Builder
	var classId, createTime,courseId,videoId sql.NullInt64

	bufStr.WriteString(`select tb.id,tb1.name,tb2.name,tb4.name,tb.create_time,tb1.tag,tb.course_id,tb.video_id from qz_class tb
inner join qz_class_info tb1 on tb1.id=tb.class_info_id
inner join qz_grade tb2 on tb2.id=tb1.grade_id
inner join qz_school tb3 on tb3.id=tb1.school_id
inner join qz_school_type tb4 on tb4.id=tb3.school_type
where tb.id=?  and tb.status=?`)
	err = db.GetSchool().QueryRow(bufStr.String(), id,enum.StatusNormal).Scan(&classId, &className, &gradeName, &schoolType, &createTime, &tag,&courseId,&videoId)

	if err != nil && err != sql.ErrNoRows {
		return nil, xlog.Error(err)
	} else if err == sql.ErrNoRows {
		return nil, xlog.Error("参数不正确")
	}
	result.Id=  classId.Int64
	result.Name=  className.String
	result.Grade=  gradeName.String
	result.SchoolType=  schoolType.String
	result.Tag= tag.String
	result.CreateTime=  createTime.Int64
	result.CourseId=courseId.Int64
	result.VideoId=videoId.Int64
	return result, nil
}

func ClassInfo(schoolId, gradeId int64) ([]*Class, error) {
	var bufStr strings.Builder
	args := make([]interface{}, 0)
	bufStr.WriteString(`select tb.id,tb2.name,tb4.name,tb.name,tb.create_time,tb.tag from qz_class_info tb
inner join qz_grade tb2 on tb2.id=tb.grade_id
inner join qz_school tb3 on tb3.id=tb.school_id
inner join qz_school_type tb4 on tb4.id=tb3.school_type`)

	bufStr.WriteString(" where tb.status=?")
	args = append(args, enum.StatusNormal)
	bufStr.WriteString(" and tb2.id=? and tb3.id=?")
	args = append(args, gradeId, schoolId)

	result, err := db.GetSchool().Query(bufStr.String(), args...)

	if err != nil {
		return nil, xlog.Error(err)
	}
	var classId, createTime sql.NullInt64
	var gradeName, schoolType, className, tag sql.NullString
	list := make([]*Class, 0)
	listMap := make(map[int64]*Class)
	for result.Next() {
		err = result.Scan(&classId, &gradeName, &schoolType, &className, &createTime, &tag)
		if err == nil {
			temp := &Class{
				Id:         classId.Int64,
				Name:       className.String,
				Grade:      gradeName.String,
				SchoolType: schoolType.String,
				Tag:        tag.String,
				CreateTime: createTime.Int64,
			}
			list = append(list, temp)
			listMap[classId.Int64] = temp
		}
	}
	return list, nil
}

func AddClassCourse(id, courseId,uid int64) error {
	return db.GetSchool().WithTransaction(func(tx *db.Tx) error {
		var tempId,count sql.NullInt64
		err:=tx.QueryRow("select count(id) from qz_class_record where class_id=? and uid=? and `type`=?",id,uid,enum.ClassAdminType).Scan(&count)
		if err!=nil&&err!=sql.ErrNoRows{
			return xlog.Error(err)
		}
		if count.Int64==0{
			return xlog.Error(errcode.CommErrorUnauthorized.Msg)
		}

		err = tx.QueryRow("select course_id from qz_class where id=?", id).Scan(&tempId)
		if err != nil && err == sql.ErrNoRows {
			return xlog.Error(errcode.HttpErrorWringParam.Msg)

		} else if err != nil {
			return xlog.Error(err)
		}
		if tempId.Int64 == 0 {
			_,err:=tx.Exec("update qz_class set course_id=? where id=?", courseId, id)
			if err != nil {
				return xlog.Error(err)
			}
		} else {
			return xlog.Error(errcode.HttpErrorWringParam.Msg)
		}
		return nil
	})
}

func RemoveClassCourse(id,uid int64) error {
	return db.GetSchool().WithTransaction(func(tx *db.Tx) error {
		var tempId,count sql.NullInt64
		err:=tx.QueryRow("select count(id) from qz_class_record where class_id=? and uid=? and `type`=?",id,uid,enum.ClassAdminType).Scan(&count)
		if err!=nil&&err!=sql.ErrNoRows{
			return xlog.Error(err)
		}
		if count.Int64==0{
			return xlog.Error(errcode.CommErrorUnauthorized.Msg)
		}

		err = tx.QueryRow("select course_id from qz_class where id=?", id).Scan(&tempId)
		if err != nil && err == sql.ErrNoRows {
			return xlog.Error(errcode.HttpErrorWringParam.Msg)

		} else if err != nil {
			return xlog.Error(err)
		}
		if tempId.Int64 != 0 {
			_,err:=tx.Exec("update qz_class set course_id=0 where id=?", id)
			if err != nil {
				return xlog.Error(err)
			}
		} else {
			return xlog.Error(errcode.HttpErrorWringParam.Msg)
		}
		return nil
	})
}




func FindClass(content string, status int64) []Class {
	list := make([]Class, 0)
	var bufStr strings.Builder
	args := make([]interface{}, 0)
	bufStr.WriteString(`select tb.id,tb2.name,tb4.name,tb1.name,tb.create_time,tb1.tag from qz_class tb
inner join qz_class_info tb1 on tb1.id=tb.class_info_id
inner join qz_grade tb2 on tb2.id=tb1.grade_id
inner join qz_school tb3 on tb3.id=tb1.school_id
inner join qz_school_type tb4 on tb4.id=tb3.school_type
inner join qz_class_record tb5 on tb5.class_id=tb.id`)
	bufStr.WriteString(" where tb.status=? and CONCAT(tb5.prefix,tb5.`code`)=? ")
	args = append(args, status, content)

	result, err := db.GetSchool().Query(bufStr.String(), args...)

	if err != nil {
		return list
	}
	var classId, createTime sql.NullInt64
	var gradeName, schoolType, className, tag sql.NullString

	for result.Next() {
		err = result.Scan(&classId, &gradeName, &schoolType, &className, &createTime, &tag)
		if err == nil {
			temp := Class{
				Id:         classId.Int64,
				Name:       className.String,
				Grade:      gradeName.String,
				SchoolType: schoolType.String,
				Tag:        tag.String,
				CreateTime: createTime.Int64,
			}
			list = append(list, temp)
		}
	}
	return list
}

func EditCourseProgress(uid,classId,videoId int64)error{
	return db.GetSchool().WithTransaction(func(tx *db.Tx) error {
		var count sql.NullInt64
		err:=tx.QueryRow("select count(id) from qz_class_record where uid=? and class_id=? and `type`=?",uid,classId,enum.ClassAdminType).Scan(&count)
		if err==nil{
			if count.Int64==0{
				return xlog.Error(errcode.CommErrorUnauthorized.Msg)
			}
			_,err=tx.Exec("update qz_class set video_id=? where id=?",videoId,classId)
			return err
		}else{
			return err
		}
	})
}

func GetUserList(schoolId,gradeId,classId,typeInt int64)[]int64{
	resultList:=make([]int64,0)
	var buf strings.Builder
	args:=make([]interface{},0)
	buf.WriteString(`select tb4.uid from qz_school_record tb 
inner join qz_school tb1 on tb1.id=tb.school_id
inner join qz_class_info tb2 on tb2.school_id=tb.school_id
inner join qz_class tb3 on tb3.class_info_id=tb2.id
inner join qz_class_record tb4 on tb4.class_id=tb3.id
where tb1.status=? and tb2.status=? and tb3.status=? and tb4.type=?
`)
	args=append(args,enum.StatusNormal,enum.StatusNormal,enum.StatusNormal)
	if typeInt==enum.ClassAdminType{
		args=append(args,enum.ClassAdminType)

	} else{
		args=append(args,enum.ClassStudentType)
		buf.WriteString(" and tb3.id=?")
	}
	if schoolId!=0{
		buf.WriteString(" and tb1.id=?")
		args=append(args,schoolId)
	}
	if gradeId!=0{
		buf.WriteString(" and tb2.grade_id=?")
		args=append(args,gradeId)
	}

	if classId!=0{
		buf.WriteString(" and tb3.id=?")
		args=append(args,classId)
	}

	buf.WriteString(" group by tb4.uid")
	result,err:=db.GetSchool().Query(buf.String(),args...)
	if err==nil{
		var uid sql.NullInt64
		for result.Next(){
			err=result.Scan(&uid)
			if err==nil{
				resultList=append(resultList,uid.Int64)
			}
		}
	}
	return resultList
}
