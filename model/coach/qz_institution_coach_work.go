package coachModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"strings"
	"time"
)

type CoachWork struct {
	Id    int64
	Uid   int64
	Start int64
	End   int64
	Max   int64
	Now   int64
	Price int64
	Desc  string
	Name  string
	Level string
	Info string
	Duration int64
	PlaceId  int64
	Reserve bool
	CourseId int64
}

func EditInstitutionCoachWork(id, start, end, place, course, uid, price, max,duration int64, level, desc, name string, del bool) error {
	if del {
		_, err := db.GetCoachDb().Exec("update   qz_institution_coach_work set `status`=? where id=?", enum.StatusDelete,id)
		if err != nil {
			return xlog.Error(err)
		}
		return nil
	} else {
		_, err := db.GetCoachDb().Exec("insert into qz_institution_coach_work (start,end,uid,place_id,course_id,price,create_time,max,`desc`,`level`,name,`duration`,`status`) values (?,?,?,?,?,?,?,?,?,?,?,?,?)",
			start, end, uid, place, course, price, time.Now().Unix(), max, desc, level, name,duration,enum.StatusNormal)
		if err != nil {
			return xlog.Error(err)
		}
		return nil
	}

}

func GetInstitutionCoachWork(placeId, nowUid, start, end int64,nowID []int64) ([]*CoachWork, error) {
	list := make([]*CoachWork, 0)
	listMap := make(map[int64]*CoachWork, 0)
	var buf strings.Builder
	args := make([]interface{}, 0)
	buf.WriteString("select tb.id,tb.start,tb.end,tb.`desc`,tb.price,tb.max,tb.`level`,tb.name,tb.uid,tb1.info,tb.duration,tb.place_id,tb.course_id from qz_institution_coach_work tb inner join qz_institution_coach " +
		" tb1 on tb1.uid=tb.uid where 1=1 ")


	if placeId != 0 {
		buf.WriteString(" and tb.place_id=? and tb.status=?")
		args = append(args, placeId,enum.StatusNormal)
	}
	if nowUid != 0 {
		buf.WriteString(" and tb.uid=? and tb.status=?")
		args = append(args, nowUid,enum.StatusNormal)
	}
	if len(nowID) > 0 {
		utils.MysqlStringInUtils(&buf,nowID," and tb.id")
	}else{
		buf.WriteString(" and tb.start>? and tb.end<?")
		args = append(args, start, end)
	}

	buf.WriteString(" order by tb.start asc")
	rows, err := db.GetCoachDb().Query(buf.String(), args...)
	if err == nil {
		var uid, id, start, end, price, max,duration,placeId,courseId sql.NullInt64
		var desc, level, name,info sql.NullString
		for rows.Next() {
			err = rows.Scan(&id, &start, &end, &desc, &price, &max, &level, &name, &uid,&info,&duration,&placeId,&courseId)
			if err == nil {
				temp := &CoachWork{
					Uid:   uid.Int64,
					Id:    id.Int64,
					Start: start.Int64,
					End:   end.Int64,
					Price: price.Int64,
					Max:   max.Int64,
					Desc:  desc.String,
					Name:  name.String,
					Level: level.String,
					Info:info.String,
					Duration:duration.Int64,
					PlaceId:placeId.Int64,
					CourseId:courseId.Int64,
				}
				list = append(list, temp)
				listMap[id.Int64] = temp
			}
		}
	}
	return list, nil
}
