package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"os"
	"mufe_service/camp/cache"
	"mufe_service/camp/errcode"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	pb "mufe_service/jsonRpc"
	"strconv"
	"strings"
	"time"
)

// AdminAuthHeader header
const AdminAuthHeader = "AdminAuthorization"
const BrandAdminAuthHeader = "BrandAdminAuthorization"

// AdminClaims 代理商
type AdminClaims struct {
	jwt.StandardClaims
	Uid             int64  `json:"uid"`
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Head            string `json:"head"`
	BusinessGroupId int64  `json:"business_group_id"`
	BusinessId      int64  `json:"business_id"`
	UserGroupName   string `json:"user_group_name"`
	UserGroupId     int64  `json:"user_group_id"`
	WithdrawPass    string `json:"withdraw_pass"`
	Token           string `json:"token"`
}

// GenerateAdminJwt 生成令牌
func GenerateAdminJwt(groupId int64, key string, data *pb.AdminUserDataResponse) (string, error) {
	pass := os.Getenv("TOKEN_ADMIN_PASSWORD")
	if pass == "" {
		return "", errors.New("token加密密码为空")
	}

	expiresAt, _ := strconv.Atoi(os.Getenv("TOKEN_ADMIN_VALIDITY_TIME"))
	if expiresAt == 0 {
		return "", errors.New("token有效期为空")
	}
	token, err := utils.GenerateUUID()
	token = strings.ReplaceAll(token, "-", "")

	if err != nil {
		return "", err
	}
	now := time.Now()
	claims := &AdminClaims{
		Token:           token,
		Uid:             data.Uid,
		Name:            data.Name,
		Phone:           data.Phone,
		Head:            data.Head,
		BusinessId:      data.BusinessId,
		BusinessGroupId: groupId,
		UserGroupName:   data.UserGroupName,
		UserGroupId:     data.UserGroupId,
		WithdrawPass:    data.WithdrawPass,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),                                             // 签发时间
			ExpiresAt: now.Add(time.Duration(expiresAt) * time.Second).Unix(), // 过期时间，必须设置
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // 生成token
	accessToken, err := jwtToken.SignedString([]byte(pass))       // 加密
	// 添加到redis
	err = cache.SetAdminToken(data.Uid, key, token)
	if err != nil {
		return "", err
	}
	if data.BusinessId != 0 {
		info := &cache.BusinessClaims{
			BusinessID:    data.BusinessId,
			BusinessName:  data.BusinessName,
			BusinessPhoto: data.BusinessPhoto,
		}
		err = cache.SetBusinessInfo(data.BusinessId, info)
		if err != nil {
			return "", err
		}
	}
	return accessToken, err
}

// CheckAgentJwt 校验
func CheckAdminJwt(key, token string) (*AdminClaims, error) {
	var claims = &AdminClaims{}
	if token == "" {
		return claims, errcode.CommErrorNotLogin.RPCError()
	}
	pass := os.Getenv("TOKEN_ADMIN_PASSWORD")
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
	err = cache.CheckAdminToken(claims.Uid, key, claims.Token)
	if err != nil {
		return claims, err
	}
	return claims, nil
}
