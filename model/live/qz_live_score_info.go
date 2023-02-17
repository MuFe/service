package liveModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/xlog"
	"strings"
)


type ScoreInfo struct{
	MemberId int64
	Time string
}


func AddScoreInfo(list []*ScoreInfo,liveId int64,tx *db.Tx)error{
	_,err:=tx.Exec("delete from qz_live_score_info where live_id=?",liveId)
	if err != nil {
		return xlog.Error(err)
	}
	if len(list)>0{
		sqlStr := "insert into qz_live_score_info (live_id,member_id,time) values"
		placeHolder := "(?,?,?)"

		var values []string
		var args []interface{}
		for _, info := range list {
			values = append(values, placeHolder)
			args = append(args,liveId)
			args = append(args, info.MemberId)
			args = append(args, info.Time)
		}
		sqlStr += strings.Join(values, ",")
		_, err := tx.Exec(sqlStr, args...)
		if err != nil {
			return xlog.Error(err)
		}
	}
	return nil
}

func GetScoreInfo(liveId int64)([]ScoreInfo,error){
	result,err:=db.GetLiveDb().Query("select time,member_id from qz_live_score_info tb where tb.live_id=?",liveId)
	if err!=nil{
		return nil,xlog.Error(err)
	}
	list:=make([]ScoreInfo,0)
	var memberId sql.NullInt64
	var time sql.NullString
	for result.Next(){
		err=result.Scan(&time,&memberId)
		if err==nil{
			list=append(list,ScoreInfo{
				MemberId:memberId.Int64,
				Time:time.String,
			})
		}
	}
	return list,nil
}
