package liveModel

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"strconv"
	"strings"
	"time"
)

type LiveData struct {
	Id int64
	Type int64
	StartTime int64
	EndTime int64
	Status int64
	Home LiveTeamData
	Visiting LiveTeamData
	HomeScore string
	VisitingScore string
	Address string

}




func GetLiveList(idList []int64,uid int64) ([]*LiveData,error) {
	var buf strings.Builder
	args:=make([]interface{},0)
	buf.WriteString(`select tb.id,tb.start_time,tb.end_time,tb.status,tb.address,tb.type,tb1.home_score,tb1.home_team_id,tb1.visiting_score,tb1.visiting_team_id,tb2.name,tb3.name from qz_live tb
inner join qz_live_info tb1 on tb1.id=tb.info_id
inner join qz_live_team tb2 on tb2.id=tb1.home_team_id
inner join qz_live_team tb3 on tb3.id=tb1.visiting_team_id`)
	if len(idList)>0{
		utils.MysqlStringInUtils(&buf,idList," where tb.id")
	}else{
		if uid!=0{
			buf.WriteString(" where tb.create_by=?")
			args=append(args,uid)
		}
	}

	result,err:=db.GetLiveDb().Query(buf.String(),args...)
	if err!=nil{
		return nil,xlog.Error(err)
	}
	list:=make([]*LiveData,0)
	var id,startTime,endTime,status,typeInt,homeTeamId,visitingTeamId sql.NullInt64
	var address,homeScore,visitingScore,homeName,visitingName sql.NullString
	for result.Next(){
		err=result.Scan(&id,&startTime,&endTime,&status,&address,&typeInt,&homeScore,&homeTeamId,&visitingScore,&visitingTeamId,&homeName,&visitingName)
		if err==nil{
			list=append(list,&LiveData{
				Id:id.Int64,
				StartTime:startTime.Int64,
				EndTime:endTime.Int64,
				Status:status.Int64,
				Address:address.String,
				Home:LiveTeamData{
					Id:homeTeamId.Int64,
					Name:homeName.String,
				},
				Visiting:LiveTeamData{
					Id:visitingTeamId.Int64,
					Name:visitingName.String,
				},
				HomeScore:homeScore.String,
				VisitingScore:visitingScore.String,
				Type:typeInt.Int64,
			})
		}
	}
	return list,nil
}

func CreateLive(typeInt,classId,matchId,createBy,packageId int64,homeName,visitingName,address string,homeMemberList,visitingMemberList []*LiveTeamMemberData)(int64,string,error){
	id:=int64(0)
	pass:=""
	if typeInt==enum.Temporary{
		err:=db.GetLiveDb().WithTransaction(func(tx *db.Tx) error {
			homeTeamId,err:=CreateLiveTeam(homeName,classId,tx)
			if err!=nil{
				return err
			}
			for _,v:=range homeMemberList{
				v.TeamID=homeTeamId
			}
			err=CreateLiveTeamMember(homeMemberList,tx)
			if err!=nil{
				return err
			}

			visitingTeamId,err:=CreateLiveTeam(visitingName,classId,tx)
			if err!=nil{
				return err
			}
			for _,v:=range visitingMemberList{
				v.TeamID=visitingTeamId
			}
			err=CreateLiveTeamMember(visitingMemberList,tx)
			if err!=nil{
				return err
			}

			infoId,err:=CreateLiveInfo(homeTeamId,visitingTeamId,matchId,tx)
			if err!=nil{
				return err
			}
			pageInfo,err:=GetPackageInfo(packageId)
			max:=int64(0)
			duration:=int64(0)
			if err==nil{
				max=pageInfo.Max
				duration=pageInfo.Duration
			}
			inviteResult, err := tx.Query("select pass from qz_live where pass<>''")
			if err != nil {
				return xlog.Error("创建直播失败")
			}
			inviteMap := make(map[string]bool)
			var invite sql.NullString
			for inviteResult.Next() {
				err = inviteResult.Scan(&invite)
				if err == nil {
					inviteMap[invite.String] = true
				}
			}
			pass = utils.BaseCode4(inviteMap)
			result,err:=tx.Exec("insert into qz_live (info_id,start_time,address,create_time,create_by,type,max,duration,pass) values(?,?,?,?,?,?,?,?,?)",infoId,time.Now().Unix(),address,time.Now().Unix(),createBy,typeInt,max,duration,pass)
			if err!=nil{
				return xlog.Error(err)
			}
			id,err=result.LastInsertId()
			if err!=nil{
				return xlog.Error(err)
			}
			return nil
		})
		 if err!=nil{
		 	return 0,pass,err
		 }
	}
	return id,pass,nil

}

func UpdateLivePass(uid int64,pass string)(string,int64,error){
	var createBy,id,duration sql.NullInt64
	err:=db.GetLiveDb().QueryRow("select create_by,id,duration from qz_live where pass=?",pass).Scan(&createBy,&id,&duration)
	if err!=nil&&err==sql.ErrNoRows{
		return "",0,xlog.Error("比赛信息有误")
	}else if err!=nil{
		return "",0,xlog.Error(err)
	}else{
		if createBy.Int64!=uid{
			return "",0,xlog.Error("您不是比赛创建人，无法进行直播")
		}
		pushDomain,pullDomain, streamName, key := "173103.push.tlivecloud.com","play.sgsports-test.com", strconv.FormatInt(id.Int64,10), "428bc1e44bed86472c8ef69dc11ec526"
		end:=int64(0)
		 if duration.Int64>0{
			 end= time.Now().Unix()+duration.Int64
		 }else{
			 end= time.Now().Unix()+86400
		 }
		address:=GetPushUrl(pushDomain, streamName, key, end)
		address1:=GetPullUrl(pullDomain, streamName, key, 0)
		_,err=db.GetLiveDb().Exec("update  qz_live set live_address=?,pull_address=? where id=?",address,address1,id.Int64)
		if err!=nil{
			return "",0,xlog.Error(err)
		}
		return address,duration.Int64,nil
	}
}



func UpdateLiveScore(id int64,homeScore,visitingScore string,list []*ScoreInfo,)error{
	var infoId sql.NullInt64
	err:=db.GetLiveDb().QueryRow("select info_id from qz_live where id=?",id).Scan(&infoId)
	if err!=nil&&err==sql.ErrNoRows{
		return xlog.Error("比赛信息有误")
	}else if err!=nil{
		return xlog.Error(err)
	}else{
		return db.GetLiveDb().WithTransaction(func(tx *db.Tx) error {
			_,err:=tx.Exec("update  qz_live_info set home_score=?,visiting_score=? where id=?",homeScore,visitingScore,infoId)
			if err!=nil{
				return xlog.Error(err)
			}
			return AddScoreInfo(list,id,tx)
		})
	}
}

func EndLive(id,uid int64)error{
	var createBy sql.NullInt64
	err:=db.GetLiveDb().QueryRow("select create_by from qz_live where id=?",id).Scan(&createBy)
	if err!=nil&&err==sql.ErrNoRows{
		return xlog.Error("比赛信息有误")
	}else if err!=nil{
		return xlog.Error(err)
	}else{
		if createBy.Int64!=uid{
			return xlog.Error("您不是比赛创建人，无法进行操作")
		}
		_,err=db.GetLiveDb().Exec("update  qz_live set status=? where id=?",enum.LIVE_END,id)
		if err!=nil{
			return xlog.Error(err)
		}
	}
	return nil
}


func End(id,uid int64)error{
	var createBy sql.NullInt64
	err:=db.GetLiveDb().QueryRow("select create_by from qz_live where id=?",id).Scan(&createBy)
	if err!=nil&&err==sql.ErrNoRows{
		return xlog.Error("比赛信息有误")
	}else if err!=nil{
		return xlog.Error(err)
	}else{
		if createBy.Int64!=uid{
			return xlog.Error("您不是比赛创建人，无法进行操作")
		}
		_,err=db.GetLiveDb().Exec("update qz_live set end_time=? where id=?",time.Now().Unix(),id)
		if err!=nil{
			return xlog.Error(err)
		}
	}
	return nil
}


func StartWatch(uid int64,pass string)(string,int64,error){
	var pullAddress sql.NullString
	var max,id,createBy,status sql.NullInt64
	err:= db.GetLiveDb().WithTransaction(func(tx *db.Tx) error {
		err:=tx.QueryRow("select id,max,pull_address,create_by,status from qz_live where pass=?",pass).Scan(&id,&max,&pullAddress,&createBy,&status)
		if err!=nil&&err==sql.ErrNoRows{
			return xlog.Error("比赛信息有误")
		}else if err!=nil {
			return xlog.Error(err)
		}else if createBy.Int64==uid{
			if status.Int64==enum.LIVE_END{
				return xlog.Error("比赛已经结束")
			}
			return nil
		}else{
			if status.Int64==enum.LIVE_END{
				return xlog.Error("比赛已经结束")
			}
			return AddWatch(id.Int64,uid,max.Int64,tx)
		}
	})
	if err==nil{
		return pullAddress.String,id.Int64,nil
	}else{
		return "",0,err
	}
}



func GetPushUrl(domain, streamName, key string, time int64)(addrstr string){
	var ext_str string
	if key != "" && time != 0{
		txTime := strings.ToUpper(strconv.FormatInt(time, 16))
		txSecret := md5.Sum([]byte(key + streamName + txTime))
		txSecretStr := fmt.Sprintf("%x", txSecret)
		ext_str = "?txSecret=" + txSecretStr + "&txTime=" + txTime
	}
	addrstr = "rtmp://" + domain + "/live/" + streamName + ext_str
	return
}

func GetPullUrl(domain, streamName, key string, time int64)(addrstr string){
	var ext_str string
	if key != "" && time != 0{
		txTime := strings.ToUpper(strconv.FormatInt(time, 16))
		txSecret := md5.Sum([]byte(key + streamName + txTime))
		txSecretStr := fmt.Sprintf("%x", txSecret)
		ext_str = "?txSecret=" + txSecretStr + "&txTime=" + txTime
	}else{
		ext_str=""
	}
	addrstr = "http://" + domain + "/live/" + streamName+".flv" + ext_str
	return
}
