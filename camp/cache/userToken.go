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
	userToken = "userToken/"
)

// CheckUserToken 检查token
func CheckUserToken( uid int64,token string) error {
	key := strconv.FormatInt(uid,10)+userToken
	validityTime, err := strconv.Atoi(os.Getenv("TOKEN_USER_VALIDITY_TIME"))
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

// SetUserToken 设置token
func SetUserToken(uid int64,token string) error {
	now := time.Now().Unix()
	key:=strconv.FormatInt(uid,10)+userToken
	return client.ZAdd(key, redis.Z{
		Score:  float64(now),
		Member: token,
	}).Err()
}

// DeleteUserToken 删除所有token
func DeleteUserToken(uid int64) error {
	key:=strconv.FormatInt(uid,10)+userToken
	return client.Del(key).Err()
}
