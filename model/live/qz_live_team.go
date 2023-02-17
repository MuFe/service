package liveModel

import (
	"mufe_service/camp/db"
	"mufe_service/camp/xlog"
)

type LiveTeamData struct {
	Id int64
	Name string
	Head string
	ClassId int64
}

func CreateLiveTeam(name string,classId int64,tx *db.Tx)(int64,error){
	result,err:=tx.Exec("insert into qz_live_team (name,class_id) values(?,?)",name,classId)
	if err!=nil{
		return 0,xlog.Error(err)
	}
	id,err:=result.LastInsertId()
	if err!=nil{
		return 0,xlog.Error(err)
	}
	return id,nil
}

