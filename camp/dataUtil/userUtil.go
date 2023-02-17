package dataUtil

import (
	"mufe_service/camp/cache"
	"mufe_service/camp/jwt"
	pb "mufe_service/jsonRpc"
)

type User struct {
	Name           string `json:"name"`
	Head           string `json:"head"`
	Phone          string `json:"phone"`
	Sex            int64  `json:"sex"`
	UserNo         string `json:"user_no"`
	UserInviteCode string `json:"user_invite_code"`
	Sign           string `json:"sign" `
	Introduce           string `json:"introduce" `
	Address       string `json:"address" `
	Age           int64  `json:"age" `
	HaveWx         bool   `json:"wx" `
	HavePass         bool   `json:"have_pass" `
	Identity         int64   `json:"identity" `
}

type AdminUser struct {
	Head             string     `json:"head"`
	Name             string     `json:"name"`
	Phone            string     `json:"phone"`
	BusinessName     string     `json:"business_name"`
	BusinessPhoto    string     `json:"business_photo"`
	Group            string     `json:"group"`
	Roles            []string   `json:"roles"`
	Sex				 int64		`json:"sex"`
	Birthday		 int64		`json:"birthday"`
}

func ParseUser(result *pb.UserDataResponse) User {
	user := User{Phone: result.Phone, Name: result.Name, Head: result.Head, Sex: result.Sex}
	user.UserNo = result.No
	user.UserInviteCode = result.InviteCode
	user.Sign = result.Sign
	user.Address = result.Address
	user.Age = result.Age
	user.HaveWx=result.HaveWx
	user.HavePass=result.HavePass
	user.Identity=result.Identity
	return user
}

func ParseUserCache(result *pb.UserDataResponse) *cache.UserClaims {
	user := ParseUser(result)
	u := &cache.UserClaims{
		Name:           user.Name,
		Head:           user.Head,
		Phone:          user.Phone,
		Sex:            user.Sex,
		UserNo:         user.UserNo,
		UserInviteCode: user.UserInviteCode,
		Sign:           user.Sign,
		Address:        user.Address,
		Age:            user.Age,
		HaveWx:         user.HaveWx,
		Identity:       user.Identity,
		HavePass:       user.HavePass,
	}
	return u
}

func ParseUserFromCache(data *cache.UserClaims) User {
	user := User{Phone: data.Phone, Name: data.Name, Head: data.Head, Sex: data.Sex}
	user.UserNo = data.UserNo
	user.UserInviteCode = data.UserInviteCode
	user.Sign = data.Sign
	user.Address = data.Address
	user.Age = data.Age
	user.HaveWx=data.HaveWx
	user.Identity=data.Identity
	user.HavePass=data.HavePass
	return user
}

func ParseAdminUser(businessInfo *cache.BusinessClaims,userResult *jwt.AdminClaims) AdminUser {
	user := AdminUser{Head: userResult.Head, Name: userResult.Name}
	user.BusinessPhoto = businessInfo.BusinessPhoto
	user.BusinessName = businessInfo.BusinessName
	user.Phone = userResult.Phone
	return user
}
