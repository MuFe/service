package liveModel

import (
	"mufe_service/camp/db"
	"mufe_service/camp/xlog"
)



func CreateLiveInfo(homeTeamId,visitingTeamId ,matchId int64,tx *db.Tx)(int64,error){
	result,err:=tx.Exec("insert into qz_live_info (home_team_id,visiting_team_id,match_id) values(?,?,?)",homeTeamId,visitingTeamId,matchId)
	if err!=nil{
		return 0,xlog.Error(err)
	}
	id,err:=result.LastInsertId()
	if err!=nil{
		return 0,xlog.Error(err)
	}
	return id,nil
}

