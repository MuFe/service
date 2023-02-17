package adminUserModel

import (
	"database/sql"
	"mufe_service/camp/db"
)

type PermissionGroup struct {
	Id   int64
	Name string
	User []PermissionUser
	Role []Permission
	UserMap map[int64]int64
	PermissionMap map[int64]int64
}

type PermissionUser struct {
	Uid   int64
	Name  string
	Phone string
}

type Permission struct {
	Id     int64
	Name   string
	Status int64
}

//获取对应平台权限角色组
func GetPermissionGroup( businessGroupId,businessId int64) []*PermissionGroup {
	list:=make([]*PermissionGroup,0)
	listMap:=make(map[int64]*PermissionGroup)
	rows, err := db.GetAdminDb().Query(`SELECT tb.name,tb.id,tb1.uid,tb4.user_name,tb4.user_phone,tb3.name,tb2.id,tb2.status from qz_business_permission_group tb 
left join qz_business_permission_group_user tb1 on tb1.permission_group_id=tb.id and tb1.business_id=?
left join qz_business_permission_group_record tb2 on tb2.permission_group_id=tb.id
left join qz_business_permission tb3 on tb3.id=tb2.permission_id
left join qz_account tb4 on tb4.id=tb1.uid
where tb.super=0 and tb.business_group_id=?`,businessId, businessGroupId)
	if err == nil {
		var pName,uName,uPhone,gName sql.NullString
		var gId,uid,pId,pStatus sql.NullInt64
		for rows.Next() {
			err=rows.Scan(&gName,&gId,&uid,&uName,&uPhone,&pName,&pId,&pStatus)
			if err==nil{
				temp,ok:=listMap[gId.Int64]
				if !ok{
					temp=&PermissionGroup{
						Id:gId.Int64,
						Name:gName.String,
						User:make([]PermissionUser,0),
						Role:make([]Permission,0),
						UserMap:make(map[int64]int64),
						PermissionMap:make(map[int64]int64),
					}
					listMap[gId.Int64]=temp
					list=append(list,temp)
				}
				if uid.Int64!=0{
					if _,ok:=temp.UserMap[uid.Int64];!ok{
						temp.User=append(temp.User,PermissionUser{
							Uid:uid.Int64,
							Name:uName.String,
							Phone:uPhone.String,
						})
						temp.UserMap[uid.Int64]=1
					}
				}
				if pId.Int64!=0{
					if _,ok:=temp.PermissionMap[pId.Int64];!ok{
						temp.Role=append(temp.Role,Permission{
							Id:pId.Int64,
							Name:pName.String,
							Status:pStatus.Int64,
						})
						temp.PermissionMap[pId.Int64]=1
					}
				}
			}
		}
	}
	return list
}
