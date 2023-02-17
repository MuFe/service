package liveModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/xlog"
)

type Package struct {
	Id int64
	Name string
	Max int64
	Duration int64
	Price int64
}

func GetPackage()([]Package,error){
	result,err:=db.GetLiveDb().Query("select id,name,max,duration,price from qz_live_package where `status`=?",enum.StatusNormal)
	list:=make([]Package,0)
	if err==nil{
		var id,max,duration,price sql.NullInt64
		var name sql.NullString
		for result.Next(){
			err=result.Scan(&id,&name,&max,&duration,&price)
			if err==nil{
				list=append(list,Package{
					Id:id.Int64,
					Name:name.String,
					Max:max.Int64,
					Price:price.Int64,
					Duration:duration.Int64,
				})
			}
		}
		return list,nil
	}else{
		return nil,xlog.Error(err)
	}
}


func GetPackageInfo(id int64)(*Package,error){
	var max,duration sql.NullInt64
	err:=db.GetLiveDb().QueryRow("select max,duration from qz_live_package where `status`=? and id=?",enum.StatusNormal,id).Scan(&max,&duration)
	if err==nil{

		return &Package{
			Max:max.Int64,
			Duration:duration.Int64,
		},nil
	}else{
		return nil,xlog.Error(err)
	}
}
