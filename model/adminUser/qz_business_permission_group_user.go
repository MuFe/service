package adminUserModel

import (
	"mufe_service/camp/db"
)

func DeleteUserPermission(uid, businessId,permissionGroupId int64)error {
	_,err := db.GetAdminDb().Exec(`delete from qz_business_permission_group_user where business_id=? and uid=? and permission_group_id=?`,businessId,uid,permissionGroupId)
	return err
}
