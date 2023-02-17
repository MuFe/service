package addressmodel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
)

type AddressData struct {
	ID     int64
	Pid    int64
	Name   string
	First  string
	Letter string
	List   []*AddressData
}


func GetAddress(id, typeInt64 int64) ([]*AddressData, error) {
	resultList:=make([]*AddressData,0)
	var result *sql.Rows
	var err error
	if typeInt64==enum.GetCity{
		if id==0{
			result,err=db.GetUserDb().Query("select i,p,n,f,l from qz_addr_city")
		} else {
			result,err=db.GetUserDb().Query("select i,p,n,f,l from qz_addr_city where p=?",id)
		}
	}else if typeInt64==enum.GetDistrict{
		if id==0{
			result,err=db.GetUserDb().Query("select i,p,n,f,l from qz_addr_district")
		} else {
			result,err=db.GetUserDb().Query("select i,p,n,f,l from qz_addr_district where p=?",id)
		}
	} else {
		result,err=db.GetUserDb().Query("SELECT p.i,p.p,p.n,p.f,p.l,c.i,c.n,c.f,c.l FROM qz_addr_province p INNER JOIN qz_addr_city c ON c.p=p.i WHERE p.p=1")
	}
	if err==nil{
		var id, pid,cid sql.NullInt64
		var name, first,letter,cname,cfirst,cletter sql.NullString
		resultMap:=make(map[int64]*AddressData)
		for result.Next() {
			err=result.Scan(&id,&pid,&name,&first,&letter,&cid,&cname,&cfirst,&cletter)
			if err==nil{
				temp,ok:=resultMap[id.Int64]
				if !ok{
					temp=&AddressData{
						ID:id.Int64,
						Name:name.String,
						First:first.String,
						Letter:letter.String,
						Pid:pid.Int64,
					}
					resultList=append(resultList,temp)
					resultMap[id.Int64]=temp
				}
				temp.List=append(temp.List,&AddressData{
					ID:cid.Int64,
					Name:cname.String,
					First:cfirst.String,
					Letter:cletter.String,
					Pid:id.Int64,
				})
			}
		}
	}
	return resultList, err
}


