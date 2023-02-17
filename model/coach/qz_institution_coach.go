package coachModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"strings"
)

type InstitutionCoach struct {
	Uid     int64
	InstitutionId int64
	Info  string
	Id int64
}

func GetInstitutionCoach(nowId,nowUid int64) ([]InstitutionCoach, error) {
	list := make([]InstitutionCoach, 0)
	var buf strings.Builder
	args := make([]interface{}, 0)
	buf.WriteString("select id,uid,iid,info from qz_institution_coach where 1=1 ")

	if nowId!=0{
		buf.WriteString(" and iid=?")
		args=append(args,nowId)
	}
	if nowUid!=0{
		buf.WriteString(" and uid=?")
		args=append(args,nowUid)
	}

	buf.WriteString(" order by id desc")
	rows, err := db.GetCoachDb().Query(buf.String(), args...)
	if err == nil {
		var uid, iid,id sql.NullInt64
		var info sql.NullString
		for rows.Next() {
			err = rows.Scan(&id,&uid, &iid,&info)
			if err == nil {
				list = append(list, InstitutionCoach{
					Uid:    uid.Int64,
					InstitutionId:  iid.Int64,
					Info:info.String,
					Id:id.Int64,
				})
			}
		}
	}
	return list, nil
}

func EditInstitutionCoach(id,uid,institutionId int64,info string,del bool) error {
	if del {
		_, err := db.GetCoachDb().Exec("delete from  qz_institution_coach where uid=? and iid=?", uid, institutionId)
		if err != nil {
			return xlog.Error(err)
		}
		return nil
	} else if id!=0{
		_,err:=db.GetCoachDb().Exec("update  qz_institution_coach set info=? where id=?",info,id)
		if err!=nil{
			return xlog.Error(err)
		}
		return nil
	}else {
		_,err:=db.GetCoachDb().Exec("insert into qz_institution_coach (uid,iid,info) values (?,?,?)",uid,institutionId,info)
		if err!=nil{
			return xlog.Error(err)
		}
		return nil
	}
}

func DelInstitutionCoach(iid int64,tx *db.Tx)error{
	_,err:=tx.Exec("delete from qz_institution_coach where iid=?",iid)
	return err
}

func DelInstitutionCoachFromUid(quitUid []int64)error{
	var err error
	if len(quitUid)==0{
		return nil
	}
	var buf strings.Builder
	buf.WriteString("delete from qz_institution_coach ")
	utils.MysqlStringInUtils(&buf,quitUid," where uid")
	_, err=db.GetSchool().Exec(buf.String())
	return err
}

