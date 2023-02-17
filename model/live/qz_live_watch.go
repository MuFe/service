package liveModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/xlog"
	"time"
)



func AddWatch(id,uid,max int64,tx *db.Tx)error{
	var count sql.NullInt64
	err:=tx.QueryRow("select count(id) from qz_live_watch where uid=? and live_id=?",uid,id).Scan(&count)
	if err!=nil{
		return xlog.Error(err)
	}
	if count.Int64==0{
		err:=tx.QueryRow("select count(id) from qz_live_watch where  live_id=?",id).Scan(&count)
		if err!=nil{
			return xlog.Error(err)
		}
		if count.Int64>=max&&max>0{
			return xlog.Error("对不起，直播间人数过多暂时无法进入")
		}
		_,err=tx.Exec("insert into qz_live_watch (uid,live_id,create_time) values(?,?,?)",uid,id,time.Now().Unix())
		if err!=nil{
			return xlog.Error(err)
		}
	} else {
		_,err:=db.GetLiveDb().Exec("update qz_live_watch set status=1 where uid=? and id=?",uid,id)
		if err!=nil{
			return xlog.Error(err)
		}
	}

	return nil
}

func OutWatch(uid int64)error{
	_,err:=db.GetLiveDb().Exec("update qz_live_watch set status=0 where uid=?",uid)
	if err!=nil{
		return xlog.Error(err)
	}
	return nil
}


func GetUserWatch(uid int64)([]int64,error){
	result,err:=db.GetLiveDb().Query("select live_id from qz_live_watch where uid=?",uid)
	if err==nil{
		var id sql.NullInt64
		list:=make([]int64,0)
		for result.Next(){
			err=result.Scan(&id)
			if err==nil{
				list=append(list,id.Int64)
			}
		}
		return list,nil
	}else{
		return nil,xlog.Error(err)
	}
}

func GetWatchNumber(id int64)([]int64,error){
	result,err:=db.GetLiveDb().Query("select uid from qz_live_watch where live_id=? and status=?",id,enum.StatusNormal)
	if err==nil{
		var uid sql.NullInt64
		re:=make([]int64,0)
		for result.Next(){
			err=result.Scan(&uid)
			if err==nil{
				re=append(re,uid.Int64)
			}
		}
		return re,nil
	}else{
		return nil,xlog.Error(err)
	}
}
