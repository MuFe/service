package cache

import (
	"github.com/go-redis/redis"
	"time"
	"mufe_service/camp/errcode"
)

var (
	phoneExpiration      = 30 * time.Minute // 验证码有效期
	phoneLimitExpiration = 99 * time.Second // 限制这个时间内不能重发
)

// checkPhoneCode 校验
func checkPhoneCode(key, code string) error {
	value, err := client.Get(key).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if value != code {
		return errcode.CommErrorVerifiedCodeWrong.RPCError()
	}
	// 校验成功，删除
	client.Del(key).Err()
	return nil
}

// getPhoneCodeLimitTime 获取手机验证码限制时间，返回秒数
func getPhoneCodeLimitTime(key string) (int64, error) {
	timeDuration, err := client.TTL(key).Result()
	if err != nil {
		return 0, err
	}
	t := phoneLimitExpiration - (phoneExpiration - timeDuration)
	if t > 0 {
		return int64(t.Seconds()), nil
	}
	return 0, nil
}

// setPhoneCode 设置手机验证码
func setPhoneCode(key, code string) (int64, error) {
	return int64(phoneLimitExpiration.Seconds()), client.Set(key, code, phoneExpiration).Err()
}
