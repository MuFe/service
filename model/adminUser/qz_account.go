package adminUserModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"strings"
	"time"
)

// User 用户表
type User struct {
	UID           int64
	BusinessId    int64
	No            string
	BusinessPhoto string
	BusinessName string
	UserGroupName string
	UserGroupId   int64
	Phone         string
	NickName      string
	Status        int64
	LoginTime     int64
	Head          string
	WithdrawPass     string
	Sex           int64
}

// GetUserByPhone 通过phone获取用户信息
func GetUserByPhone(phone, pass string) (User, error) {
	u := User{}
	list, err := GetUserQuery([]int64{}, []int64{}, "", "", phone, pass)
	if err == nil {
		if len(list) > 0 {
			u = list[0]
		} else {
			err = xlog.Error("没有该用户")
		}
	}
	return u, err
}

func Login(phone, pass string) (User, error) {
	u := User{}
	list, err := GetUserQuery([]int64{}, []int64{}, "", "", phone, pass)
	if err == nil {
		if len(list) > 0 {
			u = list[0]
			_, err = db.GetAdminDb().Exec("update qz_account  set login_time=? where uid=?", time.Now().Unix(), u.UID)
		} else {
			err = xlog.Error("没有该用户")
		}
	}
	return u, err
}

// GetUserByID 通过id获取用户信息
func GetUserByID(id int64) (User, error) {
	u := User{}
	list, err := GetUserQuery([]int64{id}, []int64{}, "", "", "", "")
	if err == nil {
		u = list[0]
	}
	return u, nil
}

func GetUserQuery(uidList, statusList []int64, no, name, phone, pass string) ([]User, error) {
	list := make([]User, 0)
	var buf strings.Builder
	args := make([]interface{}, 0)
	buf.WriteString(`select u.uid,u.no,u.user_phone,u.user_name,u.status,u.login_time,u.sex,u.user_head from qz_account u `)
	buf.WriteString(" where 1=1 ")
	utils.MysqlStringInUtils(&buf, uidList, " AND u.uid ")
	if no != "" {
		buf.WriteString(" AND u.no like (?)")
		args = append(args, "%"+no+"%")
	}
	if name != "" {
		buf.WriteString(" AND u.user_name like (?)")
		args = append(args, "%"+name+"%")
	}
	if phone != "" {
		buf.WriteString(" AND u.user_phone like (?)")
		args = append(args, "%"+phone+"%")
	}
	if pass != "" {
		buf.WriteString(" AND u.user_pass = ?")
		args = append(args, pass)
	}
	utils.MysqlStringInUtils(&buf, statusList, " AND u.status ")

	buf.WriteString(" order by u.uid desc")
	rows, err := db.GetAdminDb().Query(buf.String(), args...)
	if err == nil {
		var no, phone, nickName, head sql.NullString
		var uid, status, sex, loginTime sql.NullInt64
		for rows.Next() {
			err = rows.Scan(&uid, &no, &phone, &nickName, &status, &loginTime, &sex, &head)
			if err == nil {
				var result User
				result.UID = uid.Int64
				result.No = no.String
				result.Phone = phone.String
				result.NickName = nickName.String
				result.Status = status.Int64
				result.LoginTime = loginTime.Int64
				result.Sex = sex.Int64
				result.Head = head.String
				list = append(list, result)
			}
		}
		return list, nil
	} else {
		return nil, err
	}
}

//获取登录后的用户信息
func GetLoginUserInfo(uid, businessGroupId int64) (*User, error) {
	// 获取用户-商家分组关联信息
	var businessId, userGroupID, ubID, ubStatus, usbStatus, bgrStatus sql.NullInt64
	var adminNo, userGroupName, payPass,name,photo sql.NullString
	err := db.GetAdminDb().QueryRow(`
		select b.no,ub.id,ub.business_id,tb1.id,tb1.name,ub.status,usb.status,bgr.status,usb.withdraw_pass,bgr.photo,bgr.name 
		from qz_account usb 
		left join qz_business_permission_group_user ub on ub.uid=usb.uid  
		left join qz_business b on b.id=ub.business_id 
		left join qz_business_permission_group tb1 on tb1.id = ub.permission_group_id
		left join qz_business_group_record bgr on tb1.business_group_id = bgr.group_id  and bgr.business_id=ub.business_id
		where usb.uid = ? and tb1.business_group_id = ?`,
		uid, businessGroupId).Scan(&adminNo, &ubID, &businessId, &userGroupID, &userGroupName, &ubStatus, &usbStatus, &bgrStatus, &payPass,&photo,&name)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if ubStatus.Int64 != enum.StatusNormal || usbStatus.Int64 != enum.StatusNormal || bgrStatus.Int64 != enum.StatusNormal {
		return nil, errcode.CommErrorAccountUnavailale.RPCError()
	}
	//更新登录时间
	err = UpdateLoginTime(ubID.Int64)
	if err != nil {
		return nil, err
	}
	return &User{
		UID:           uid,
		BusinessId:    businessId.Int64,
		No:            adminNo.String,
		UserGroupName: userGroupName.String,
		UserGroupId:   userGroupID.Int64,
		WithdrawPass:  payPass.String,
		BusinessName:name.String,
		BusinessPhoto:photo.String,
	}, nil
}

//更新登录时间
func UpdateLoginTime(ubId int64) error {
	_, err := db.GetAdminDb().Exec("update qz_business_permission_group_user set login_time = ? where id = ?", time.Now().Unix(), ubId)
	if err != nil {
		return xlog.Error(err)
	}
	return nil
}
