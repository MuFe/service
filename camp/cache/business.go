package cache

import (
	"github.com/mitchellh/mapstructure"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	"strconv"
)

const (
	businessInfo = "businessInfo/"
)


// UserClaims
type BusinessClaims struct {
	BusinessID   int64
	BusinessName  string
	BusinessPhoto string
}

// GetBusinessInfo 获取商家信息
func GetBusinessInfo(businessId int64) (*BusinessClaims, error) {
	return getBusinessInfo(businessInfo + strconv.FormatInt(businessId, 10))
}

// SetBusinessInfo 设置公司信息
func SetBusinessInfo(businessId int64, value *BusinessClaims) error {
	return setBusinessInfo(businessInfo+strconv.FormatInt(businessId, 10), value)
}

// DeleteBusinessInfo 删除
func DeleteBusinessInfo(businessId int64) error {
	key := businessInfo + strconv.FormatInt(businessId, 10)
	return client.Del(key).Err()
}

// getBusinessInfo 获取商家信息
func getBusinessInfo(key string) (*BusinessClaims, error) {
	user := &BusinessClaims{}
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

// setBusinessInfo 设置商家信息
func setBusinessInfo(key string, value *BusinessClaims) error {
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
