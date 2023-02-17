package cache

import (
	"github.com/mitchellh/mapstructure"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"strconv"
)

const (
	userLoginCode = "userCode/"
)

//
// CheckUserPhoneCode 检查验证码
func CheckUserPhoneCode(phone, code string, typeInt int64) error {
	return checkPhoneCode(userLoginCode+strconv.FormatInt(typeInt, 10)+phone, code)
}

// GetUserPhoneCodeLimitTime 获取验证码
func GetUserPhoneCodeLimitTime(phone string, typeInt int64) (int64, error) {
	return getPhoneCodeLimitTime(userLoginCode + strconv.FormatInt(typeInt, 10) + phone)
}

// SetUserPhoneCode 设置验证码
func SetUserPhoneCode(phone, code string, typeInt int64) (int64, error) {
	return setPhoneCode(userLoginCode+strconv.FormatInt(typeInt, 10)+phone, code)
}

const (
	//用户信息
	userInfo = "userInfo/"
)

// UserClaims
type UserClaims struct {
	Name           string
	Head           string
	Phone          string
	Sex            int64
	Identity       int64
	UserNo         string
	UserInviteCode string
	Sign           string
	Address        string
	Age            int64
	HaveWx         bool
	HavePass       bool
}

// GetUserInfo 获取用户信息
func GetUserInfo(uid int64) (*UserClaims, error) {
	return getUserInfo(userInfo + strconv.FormatInt(uid, 10))
}

// SetUserInfo 设置用户信息
func SetUserInfo(uid int64, value *UserClaims) error {
	return setUserInfo(userInfo+strconv.FormatInt(uid, 10), value)
}

// DeleteUserToken 删除所有token
func DeleteUserInfo(uid int64) error {
	key := userInfo + strconv.FormatInt(uid, 10)
	return client.Del(key).Err()
}

// getUserInfo 获取用户信息
func getUserInfo(key string) (*UserClaims, error) {
	user := &UserClaims{}
	res := make(map[string]interface{})
	if err := mapstructure.Decode(user, &res); err != nil {
		xlog.ErrorP(err)
		return nil, err
	}
	j := 0
	keys := make([]string, len(res))
	for k := range res {
		keys[j] = k
		j++
	}
	info, err := client.HGetAll(key).Result()
	if err != nil {
		xlog.ErrorP(err)
		return nil, err
	}
	for _, v := range keys {
		temp, err := utils.Scan([]byte(info[v]), res[v])
		if err != nil {
			xlog.ErrorP(err, info[v], v)
		} else {
			res[v] = temp
		}
	}
	if err := mapstructure.Decode(res, user); err != nil {
		xlog.ErrorP(err)
		return nil, err
	}
	return user, nil
}

// setUserInfo 设置用户信息
func setUserInfo(key string, value *UserClaims) error {
	res := make(map[string]interface{})
	if err := mapstructure.Decode(value, &res); err != nil {
		xlog.ErrorP(err)
		return err
	}
	info := client.HMSet(key, res)
	if info.Err() != nil {
		xlog.ErrorP(info.Err())
		return info.Err()
	} else {
		return nil
	}
}
