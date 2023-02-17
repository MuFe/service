package cache

import (
	"os"
	"strconv"
	"time"
	"mufe_service/camp/errcode"

	"github.com/go-redis/redis"
	"mufe_service/camp/xlog"
)

const (
	AgentToken = "adminToken/"
	BrandAgentToken = "brandAdminToken/"
)

// CheckAdminToken 检查token
func CheckAdminToken(uid int64,tokenKey,token string) error {
	key:=strconv.FormatInt(uid,10)+tokenKey
	validityTime, err := strconv.Atoi(os.Getenv("TOKEN_ADMIN_VALIDITY_TIME"))
	if err != nil || validityTime <= 0 {
		return xlog.Error("token有效期设置错误", validityTime, err)
	}
	expiresAt := int(time.Now().Add(time.Duration(-validityTime) * time.Second).Unix())
	// 删除过期token
	_, err = client.ZRemRangeByScore(key, "-inf", strconv.Itoa(expiresAt)).Result()
	if err != nil {
		return err
	}

	// 判断是否有效
	score, err := client.ZScore(key, token).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if score == 0 {
		return errcode.CommErrorTokenExpired.RPCError()
	}
	return nil
}

// SetAdminToken 设置token
func SetAdminToken(uid int64,keyStr,token string) error {
	now := time.Now().Unix()
	key:=strconv.FormatInt(uid,10)+keyStr
	DeleteAdminToken(uid,keyStr)
	return client.ZAdd(key, redis.Z{
		Score:  float64(now),
		Member: token,
	}).Err()
}

// DeleteAdminToken 删除所有token
func DeleteAdminToken(uid int64,keyStr string) error {
	key:=strconv.FormatInt(uid,10)+keyStr
	return client.Del(key).Err()
}
