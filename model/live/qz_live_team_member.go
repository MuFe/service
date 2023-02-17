package liveModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/xlog"
	"strings"
)

type LiveTeamMemberData struct {
	Id int64
	Name string
	Number string
	TeamID int64
	Uid int64
}

func CreateLiveTeamMember(list []*LiveTeamMemberData,tx *db.Tx)error{
	if len(list)>0{
		sqlStr := "insert into qz_live_team_member (number,name,team_id,uid) values"
		placeHolder := "(?,?,?,?)"

		var values []string
		var args []interface{}
		for _, info := range list {
			values = append(values, placeHolder)
			args = append(args, info.Number)
			args = append(args, info.Name)
			args = append(args, info.TeamID)
			args = append(args, info.Uid)
		}
		sqlStr += strings.Join(values, ",")
		_, err := tx.Exec(sqlStr, args...)
		if err != nil {
			return xlog.Error(err)
		}
	}
	return nil
}


func GetLiveTeamMember(liveId int64)([]LiveTeamMemberData,[]LiveTeamMemberData,error){
	result,err:=db.GetLiveDb().Query(`select tb.number,tb.id,tb.name,tb.uid,tb.team_id from qz_live_team_member tb 
inner join qz_live_info tb1 on tb1.home_team_id=tb.team_id
inner join qz_live tb2 on tb2.info_id=tb1.id
where tb2.id=?`,liveId)
	if err!=nil{
		return nil,nil,xlog.Error(err)
	}
	home:=make([]LiveTeamMemberData,0)
	visiting:=make([]LiveTeamMemberData,0)
	var id,teamId,uid sql.NullInt64
	var name,number sql.NullString
	for result.Next(){
		err=result.Scan(&number,&id,&name,&uid,&teamId)
		if err==nil{
			home=append(home,LiveTeamMemberData{
				Id:id.Int64,
				Name:name.String,
				Number:number.String,
				Uid:uid.Int64,
				TeamID:teamId.Int64,
			})
		}

	}

	result,err=db.GetLiveDb().Query(`select tb.number,tb.id,tb.name,tb.uid,tb.team_id from qz_live_team_member tb 
inner join qz_live_info tb1 on tb1.visiting_team_id=tb.team_id
inner join qz_live tb2 on tb2.info_id=tb1.id
where tb2.id=?`,liveId)
	if err!=nil{
		return nil,nil,xlog.Error(err)
	}

	for result.Next(){
		err=result.Scan(&number,&id,&name,&uid,&teamId)
		if err==nil{
			visiting=append(visiting,LiveTeamMemberData{
				Id:id.Int64,
				Name:name.String,
				Number:number.String,
				Uid:uid.Int64,
				TeamID:teamId.Int64,
			})
		}

	}

	return home,visiting,nil

}

