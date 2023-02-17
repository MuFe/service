package usermodel

import (
	"database/sql"
	"fmt"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"strings"
	"time"
)

// User 用户表
type User struct {
	UID             int64
	No              string
	InviteCode      string
	Phone           string
	NickName        string
	Sign            string
	Status          int64
	LoginTime       int64
	Head            string
	Pass            string
	Sex             int64
	Age             int64
	Identity        int64
	Device          string
	Province        string
	City            string
	HaveWx          bool
	CanEditHead     int64
	CanEditName     int64
	CancelStartTime int64
	Address         string
	RegistrationId  string
	LastType        int64
}

// GetUserByUnionid 通过Unionid获取用户信息
func GetUserByUnionid(unionid string) (User, error) {
	u := User{}
	list, _, err := GetUserQuery([]int64{}, []int64{}, "", "", unionid, "", "", 0, 0, enum.CANCEL_DEFUALT_TYPE)
	if err == nil {
		if len(list) > 0 {
			u = list[0]
		}
	}
	return u, nil
}

// GetUserByPhone 通过phone获取用户信息
func GetUserFromPhone(phone string) ([]User, error) {
	list, _, err := GetUserQuery([]int64{}, []int64{enum.StatusNormal}, "", "", "", phone, "", 0, 0, enum.CANCEL_DEFUALT_TYPE)
	return list, err
}

// GetUserByPhone 通过phone获取用户信息
func GetUserByPhone(phone, pass string) (User, error) {
	u := User{}
	list, _, err := GetUserQuery([]int64{}, []int64{enum.StatusNormal}, "", "", "", phone, pass, 0, 0, enum.CANCEL_DEFUALT_TYPE)
	if err == nil {
		if len(list) > 0 {
			u = list[0]
		}
	}else{
		xlog.ErrorP(err)
	}
	return u, nil
}

// GetUserByID 通过id获取用户信息
func GetUserByID(id int64) (User, error) {
	u := User{}
	list, _, err := GetUserQuery([]int64{id}, []int64{}, "", "", "", "", "", 0, 0, enum.CANCEL_DEFUALT_TYPE)
	if err == nil {
		if len(list) > 0 {
			u = list[0]
		}
	}
	return u, nil
}

// GetUserByNo 通过no获取用户信息
func GetUserByNo(nob string) ([]User, error) {
	list, _, err := GetUserQuery([]int64{}, []int64{enum.StatusNormal}, nob, "", "", "", "", 0, 0, enum.CANCEL_DEFUALT_TYPE)
	return list, err
}

// GetUserByName 通过昵称获取用户信息
func GetUserByName(name string) ([]User, error) {
	list, _, err := GetUserQuery([]int64{}, []int64{enum.StatusNormal}, "", name, "", "", "", 0, 0, enum.CANCEL_DEFUALT_TYPE)
	return list, err
}

func GetUserQuery(uidList, statusList []int64, no, name, unionId, phone, pass string, page, size, cancelType int64) ([]User, int64, error) {
	list := make([]User, 0)
	var buf strings.Builder
	args := make([]interface{}, 0)
	buf.WriteString(`select %s from qz_user u `)
	buf.WriteString(` left join qz_outh o on o.uid = u.uid and o.type = ?`)
	args = append(args, enum.OuthTypeWx)
	buf.WriteString(" where 1=1 ")
	utils.MysqlStringInUtils(&buf, uidList, " AND u.uid ")
	if no != "" {
		buf.WriteString(" AND u.no like (?)")
		args = append(args, "%"+no+"%")
	}
	if name != "" {
		buf.WriteString(" AND u.nick_name like (?)")
		args = append(args, "%"+name+"%")
	}
	if phone != "" {
		buf.WriteString(" AND u.phone like (?)")
		args = append(args, phone+"%")
	}
	if unionId != "" {
		buf.WriteString(" AND o.outh_id = ?")
		args = append(args, unionId)
	}
	if pass != "" {
		buf.WriteString(" AND u.pwd = ?")
		args = append(args, pass)
	}
	if cancelType == enum.CANCEL_START_TYPE {
		buf.WriteString(" AND u.cancel_start_time <>0 and u.cancel_end_time=0")
	} else if cancelType == enum.CANCEL_END_TYPE {
		buf.WriteString(" AND u.cancel_start_time <>0 and u.cancel_end_time<>0")
	}
	utils.MysqlStringInUtils(&buf, statusList, " AND u.status ")
	buf.WriteString(" order by u.uid desc")
	var count sql.NullInt64
	err := db.GetUserDb().QueryRow(fmt.Sprintf(buf.String(), " count(u.uid)"), args...).Scan(&count)
	if err != nil {
		xlog.ErrorP(err)
		return nil, 0, err
	}

	if page != 0 {
		buf.WriteString(" limit ?,?")
		start := (page - 1) * size
		args = append(args, start, size)
	}
	rows, err := db.GetUserDb().Query(fmt.Sprintf(buf.String(), "u.uid,u.no,u.phone,u.nick_name,u.status,u.login_time,u.sex,u.head,u.invite_code,u.last_device,u.user_type,u.pwd,u.sign,o.id,u.province,u.city,u.head_edit,u.name_edit,u.cancel_start_time,u.address,u.age,u.registration_id,u.last_type"), args...)
	if err == nil {
		var no, phone, nickName, head, inviteCode, device, pwd, sign, province, city, address, registrationId sql.NullString
		var uid, status, sex, loginTime, identity, outhId, canEditHead, canEditName, cancelStartTime, age, lastType sql.NullInt64
		for rows.Next() {
			err = rows.Scan(&uid, &no, &phone, &nickName, &status, &loginTime, &sex, &head, &inviteCode, &device, &identity, &pwd, &sign, &outhId, &province, &city, &canEditHead, &canEditName, &cancelStartTime, &address, &age, &registrationId, &lastType)
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
				result.InviteCode = inviteCode.String
				result.Device = device.String
				result.Identity = identity.Int64
				result.Pass = pwd.String
				result.Sign = sign.String
				result.HaveWx = outhId.Int64 != 0
				result.Province = province.String
				result.City = city.String
				result.CanEditHead = canEditHead.Int64
				result.CanEditName = canEditName.Int64
				result.CancelStartTime = cancelStartTime.Int64
				result.Age = age.Int64
				result.Address = address.String
				result.LastType=lastType.Int64
				result.RegistrationId = registrationId.String
				list = append(list, result)
			}
		}
		return list, count.Int64, nil
	} else {
		xlog.ErrorP(err)
		return nil, count.Int64, err
	}
}

//创建用户
func CreateUserOuth(uName, uHead, no, inviteCode, unionId, appId, openId, sign string, uSex, outhType int64) (int64, error) {
	t := time.Now()
	uid := int64(0)
	err := db.GetUserDb().WithTransaction(func(tx *db.Tx) error {
		result, err := tx.Exec("insert into qz_user(nick_name,sex,head,created,no,invite_code,sign) values(?,?,?,?,?,?,?)",
			uName, uSex, uHead, t.Unix(), no, inviteCode, sign)
		if err != nil {
			return xlog.Error("注册失败")
		}
		uid, _ = result.LastInsertId()
		err = CreateOuth(uid, outhType, unionId, tx)
		if err != nil {
			return err
		}
		err = CreateOpenId(uid, appId, openId, tx)
		return err
	})
	return uid, err
}

//创建用户
func CreateUser(uName, phone, no, inviteCode, pass, sign string, uSex int64) (int64, error) {
	t := time.Now()
	uid := int64(0)
	err := db.GetUserDb().WithTransaction(func(tx *db.Tx) error {
		result, err := tx.Exec("insert into qz_user(nick_name,sex,phone,created,no,invite_code,pwd,sign) values(?,?,?,?,?,?,?,?)",
			uName, uSex, phone, t.Unix(), no, inviteCode, pass, sign)
		if err != nil {
			return xlog.Error("注册失败")
		}
		uid, _ = result.LastInsertId()
		return err
	})
	return uid, err
}

func UpdateUser(userInfo User, name, head string, sex int64) {
	var buf strings.Builder
	buf.WriteString("update qz_user set ")
	querys := make([]string, 0)
	args := make([]interface{}, 0)
	if name != "" && userInfo.NickName != name && userInfo.CanEditName == 0 {
		querys = append(querys, "nick_name=?")
		args = append(args, name)
	}
	if head != "" && head != userInfo.Head && userInfo.CanEditHead == 0 {
		querys = append(querys, "head=?")
		args = append(args, head)
	}
	if sex != 0 && sex != userInfo.Sex {
		querys = append(querys, "sex=?")
		args = append(args, sex)
	}

	if len(querys) > 0 {
		buf.WriteString(strings.Join(querys, ","))
		buf.WriteString(" where uid=?")
		args = append(args, userInfo.UID)
		_, _ = db.GetUserDb().Exec(buf.String(), args...)
	}
}

func UpdateLoginInfo(device string, uid, lastType int64) {
	db.GetUserDb().Exec("update qz_user set last_device=?,login_time=?,cancel_start_time=0,last_type=? where uid=?", device, time.Now().Unix(), lastType, uid)
}

func UpdateUserPhone(phone string, uid int64) error {
	_, err := db.GetUserDb().Exec("update qz_user set phone=? where uid=?", phone, uid)
	return err
}

func UpdateUserName(name string, uid int64) error {
	_, err := db.GetUserDb().Exec("update qz_user set nick_name=?,name_edit=1 where uid=?", name, uid)
	return err
}

func UpdateUserHead(head string, uid int64) error {
	_, err := db.GetUserDb().Exec("update qz_user set head=?,head_edit=1 where uid=?", head, uid)
	return err
}

func UpdateUserSign(sign string, uid int64) error {
	_, err := db.GetUserDb().Exec("update qz_user set sign=? where uid=?", sign, uid)
	return err
}

func UpdateUserAddress(uid int64, address string) error {
	return db.GetUserDb().WithTransaction(func(tx *db.Tx) error {
		_, err := tx.Exec("update qz_user set address=? where uid=? ", address, uid)
		return err
	})
}

func UpdatePushInfo(id string, uid int64) error {
	return db.GetUserDb().WithTransaction(func(tx *db.Tx) error {
		_, err := tx.Exec("update qz_user set registration_id=? where uid=? ", id, uid)
		return err
	})
}

func UpdateUserSex(sex, uid int64) error {
	_, err := db.GetUserDb().Exec("update qz_user set sex=? where uid=?", sex, uid)
	return err
}

func UpdateUserAge(age, uid int64) error {
	_, err := db.GetUserDb().Exec("update qz_user set age=? where uid=?", age, uid)
	return err
}

func UpdateUserPass(phone, pass string) error {
	_, err := db.GetUserDb().Exec("update qz_user set pwd=? where phone=?", pass, phone)
	return err
}

func UpdateUserModifyPass(pass, newPass string, uid int64) error {
	var oldPass sql.NullString
	err := db.GetUserDb().QueryRow("select pwd from qz_user where uid=?", uid).Scan(&oldPass)
	if err == nil {
		if oldPass.String != pass {
			return xlog.Error("旧密码不正确")
		} else {
			_, err = db.GetUserDb().Exec("update qz_user set pwd=? where uid=?", newPass, uid)
		}
	} else if err != sql.ErrNoRows {

	} else {
		err = xlog.Error("参数有误")
	}
	return err
}

func UpdateUserIdentity(identity int64,uidList []int64) error {
	if len(uidList)==0{
		return nil
	}
	var buf strings.Builder
	buf.WriteString("update qz_user set user_type=? ")
	utils.MysqlStringInUtils(&buf,uidList," where uid")
	_, err := db.GetUserDb().Exec(buf.String(), identity)
	return err
}

func Unbind(uid, outhType int64) error {
	return db.GetUserDb().WithTransaction(func(tx *db.Tx) error {
		err := UnbindOuth(uid, outhType, tx)
		if err != nil {
			return err
		}
		err = UnbindOpenId(uid, tx)
		return err
	})
}

func CancelUser(uid int64) error {
	return db.GetUserDb().WithTransaction(func(tx *db.Tx) error {
		var status, cancelStart sql.NullInt64
		err := tx.QueryRow("select `status`,`cancel_start_time` from qz_user where uid=?", uid).Scan(&status, &cancelStart)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if err == sql.ErrNoRows {
			return xlog.Error("用户ID不存在")
		}
		if status.Int64 != enum.StatusNormal {
			return xlog.Error("用户状态有误，请联系客服")
		}
		if cancelStart.Int64 != 0 {
			return xlog.Error("该用户已经注销，请勿重复操作")
		}
		_, err = tx.Exec("update qz_user set cancel_start_time=? where uid=?", time.Now().Unix(), uid)
		return err
	})
}

func EnterCancel(uidList []int64) error {
	var buf strings.Builder
	args := make([]interface{}, 0)
	buf.WriteString("update qz_user set cancel_end_time=?,status=? where")
	utils.MysqlStringInUtils(&buf, uidList, " uid")
	args = append(args, time.Now().Unix(), enum.StatusDelete)
	return db.GetUserDb().WithTransaction(func(tx *db.Tx) error {
		_,err:=tx.Exec(buf.String(),args...)
		if err!=nil{
			return xlog.Error(err)
		}
		var buf strings.Builder
		buf.WriteString("delete from qz_outh where ")
		utils.MysqlStringInUtils(&buf, uidList, " uid")
		_,err=tx.Exec(buf.String())
		if err!=nil{
			return xlog.Error(err)
		}
		var buf1 strings.Builder
		buf1.WriteString("delete from qz_open_id where ")
		utils.MysqlStringInUtils(&buf1, uidList, " uid")
		_,err=tx.Exec(buf1.String())
		if err!=nil{
			return xlog.Error(err)
		}
		return nil
	})
}

func BatchModifyType(uidList []int64, userType int64) error {
	if len(uidList) == 0 {
		return nil
	}
	buf := strings.Builder{}
	args := make([]interface{}, 0)
	buf.WriteString("update qz_user set user_type=? where 1=1")
	args = append(args, userType)
	utils.MysqlStringInUtils(&buf, uidList, " and uid")
	_, err := db.GetUserDb().Exec(buf.String(), args...)
	if err != nil {
		return xlog.Error(err)
	}
	return nil
}
