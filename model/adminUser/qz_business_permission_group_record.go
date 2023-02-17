package adminUserModel

import (
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/xlog"
)

func CheckGroupPermission(uid, businessGroupId,businessId int64, roleName string)bool {
	var count int64
	err:=db.GetAdminDb().QueryRow(`select count(tb1.id) from qz_business_permission_group_user tb 
inner join qz_business_permission_group tb1 on tb1.id=tb.permission_group_id
where tb.uid=? and tb1.business_group_id=? and tb1.super=1 and tb.business_id=?`,uid,businessGroupId,businessId).Scan(&count)
	if err!=nil{
		xlog.ErrorP(err)
	} else if count==1{
		//超级管理员有所有的权限
		return true
	}
	err = db.GetAdminDb().QueryRow(`select count(tb.id) from qz_business_permission_group_user tb 
		inner join qz_business_permission_group tb1 on tb1.id=tb.permission_group_id
		inner join qz_business_permission_group_record tb2 on tb2.permission_group_id=tb1.id
		inner join qz_business_permission tb3 on tb3.id=tb2.permission_id
		where tb1.business_group_id=? and tb3.role_name=? and tb.uid=? and tb3.status=? and tb2.status=? and tb.business_id=?`, businessGroupId, roleName, uid,enum.StatusNormal,enum.StatusNormal,businessId).Scan(&count)
	if err!=nil{
		xlog.ErrorP(err)
	}
	return count!=0
}

func GetUserPagePermission(uid, businessGroupId,businessId int64)[]string{
	roles:=make([]string,0)
	var count int64
	err:=db.GetAdminDb().QueryRow(`select count(tb1.id) from qz_business_permission_group_user tb 
inner join qz_business_permission_group tb1 on tb1.id=tb.permission_group_id
where tb.uid=? and tb1.business_group_id=? and tb1.super=1 and tb.business_id=?`,uid,businessGroupId,businessId).Scan(&count)
	if err!=nil{
		xlog.ErrorP(err)
	} else if count==1{
		//超级管理员有所有的权限
		roles=append(roles,"system_super")
		return roles
	}
	rows,err := db.GetAdminDb().Query(`select tb3.role_name from qz_business_permission_group_user tb 
		inner join qz_business_permission_group tb1 on tb1.id=tb.permission_group_id
		inner join qz_business_permission_group_record tb2 on tb2.permission_group_id=tb1.id
		inner join qz_business_group_permission tb3 on tb3.id=tb2.permission_id
		where tb1.business_group_id=? and tb3.permission_type=1 and tb.uid=? and tb3.status=? and tb2.status=? and tb.business_id=?`, businessGroupId, uid,enum.StatusNormal,enum.StatusNormal,businessId)
	if err==nil{
		var name string
		for rows.Next(){
			err=rows.Scan(&name)
			if err==nil{
				roles=append(roles,name)
			} else {
				xlog.ErrorP(err)
			}
		}
	}
	return roles
}
