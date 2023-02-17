package jwt

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
	"mufe_service/camp/cache"
	"mufe_service/camp/errcode"
	"mufe_service/camp/utils"

	"github.com/golang-jwt/jwt"
	"mufe_service/camp/xlog"
)

// AuthHeader header
const AuthHeader = "Authorization"

// UserClaims
type UserClaims struct {
	jwt.StandardClaims
	Uid      int64  `json:"uid"`
	OpenId   string `json:"open_id"`
	Token    string `json:"token"`
	Identity int64  `json:"identity"`
}

// GenerateUserJwt 生成令牌
func GenerateUserJwt(uid, identity int64, openId string) (string, error) {
	pass := os.Getenv("TOKEN_USER_PASSWORD")
	if pass == "" {
		return "", errors.New("token加密密码为空")
	}

	expiresAt, _ := strconv.Atoi(os.Getenv("TOKEN_USER_VALIDITY_TIME"))
	if expiresAt == 0 {
		return "", errors.New("token有效期为空")
	}
	token, err := utils.GenerateUUID()
	token = strings.ReplaceAll(token, "-", "")

	if err != nil {
		return "", err
	}
	now := time.Now()
	claims := &UserClaims{
		Token:  token,
		Uid:    uid,
		OpenId: openId,
		Identity: identity,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),                                             // 签发时间
			ExpiresAt: now.Add(time.Duration(expiresAt) * time.Second).Unix(), // 过期时间，必须设置
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // 生成token
	accessToken, err := jwtToken.SignedString([]byte(pass))       // 加密
	// 添加到redis
	err = cache.SetUserToken(uid,token)
	if err != nil {
		return "", err
	}
	return accessToken, err
}

// CheckUserJwt 校验
func CheckUserJwt(token string) (*UserClaims, error) {
	var claims = &UserClaims{}
	if token == "" {
		return claims, errcode.CommErrorNotLogin.RPCError()
	}
	pass := os.Getenv("TOKEN_USER_PASSWORD")
	if pass == "" {
		return claims, errors.New("token加密密码为空")
	}
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(pass), nil
	})
	if err != nil {
		if strings.Index(err.Error(), "token is expired by") != -1 {
			return claims, errcode.CommErrorTokenExpired.RPCError()
		}
		return claims, xlog.Error(err)
	}
	if !jwtToken.Valid {
		return claims, errcode.CommErrorTokenExpired.RPCError()
	}
	// 判断redis

	err = cache.CheckUserToken(claims.Uid,claims.Token)
	if err != nil {
		return claims, err
	}
	return claims, nil
}
