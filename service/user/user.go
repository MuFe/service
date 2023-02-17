package user

import (
	"context"
	"fmt"
	"github.com/Timothylock/go-signin-with-apple/apple"
	"mufe_service/camp/enum"
	"mufe_service/camp/sequence"
	"mufe_service/camp/service"
	"mufe_service/camp/wx/v2"
	"mufe_service/camp/xlog"
	pb "mufe_service/jsonRpc"
	addressmodel "mufe_service/model/address"
	"mufe_service/model/user"
)

var (
	teamId   = "Y799SFNZ48"
	keyId    = "2MFKN5ZQYV"
	clientId = "com.dzss.qiuzhi"
	secret   = `
-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgaoO8zZGULbAauDd9
kCLlHa2anusDYRcEBUZ8HbGVphKgCgYIKoZIzj0DAQehRANCAAQU4wT9+Lhglwq0
U3kkTTzdT3DGb6VQ8n2Ly1xM+0NCIh426sqQ5sCPX35cwLV5CNrmCMnhIb7z7X8w
XsXHD+Vp
-----END PRIVATE KEY-----`
)

func init() {
	nSer := &rpcServer{}
	pb.RegisterUserServiceServer(service.GetRegisterRpc(), nSer)

}

type rpcServer struct {
}

func (rpc *rpcServer) Login(ctx context.Context, request *pb.LoginRequest) (*pb.UserDataResponse, error) {
	userResult := &pb.UserDataResponse{}
	uid := int64(0)
	var uName, uHead string
	var uSex int64
	var nickname, avatar, no, inviteCode, phone, sign string
	var gender int64
	isNew := false
	userResultInfo := usermodel.User{}
	lastType := int64(0)
	if request.Type == enum.OuthType {
		userInfo := &wx.UserInfo{}
		var err error
		if request.IsMiniPrograms {
			lastType = enum.LoginTypeWxMini
			//userInfo, err := wx.Login(request.AppId, request.Secret, request.Code)
			//if err != nil {
			//	xlog.ErrorP(err)
			//	return nil, err
			//}
			//openId = userInfo.OpenID
			//UserInfoResult, err := wx.DecryptUserInfo(request.RawData, request.EncryptedData, request.Signature, request.Iv, userInfo.SessionKey)
			//if err != nil {
			//	return nil, err
			//}
			//uSex = int64(UserInfoResult.Gender)
			//uName = UserInfoResult.Nickname
			//uHead = UserInfoResult.Avatar
			//unionId = UserInfoResult.UnionID
		} else if request.OuthType == enum.AppleType {
			lastType = enum.LoginTypeApp
			secret, err := apple.GenerateClientSecret(secret, teamId, clientId, keyId)

			client := apple.New()

			vReq := apple.AppValidationTokenRequest{
				ClientID:     clientId,
				ClientSecret: secret,
				Code:         request.Code,
			}

			var resp apple.ValidationResponse

			err = client.VerifyAppToken(ctx, vReq, &resp)
			if err != nil {
				return nil, err

			}

			if resp.Error != "" {
				return nil, err

			}

			// Get the unique user ID
			unique, err := apple.GetUniqueID(resp.IDToken)
			if err != nil {
				return nil, err

			}

			// Get the email
			claim, err := apple.GetClaims(resp.IDToken)
			if err != nil {
				return nil, err

			}

			email := fmt.Sprintf("%v", (*claim)["email"])
			// Voila!
			fmt.Println(unique)
			fmt.Println(email)
			userInfo.Unionid = unique
			userInfo.OpenID = unique
			userInfo.Email = email
			userInfo.Nickname = request.Name
			request.AppId = clientId
		} else {
			lastType = enum.LoginTypeWxApp
			userInfo, err = wx.GetUserInfo(request.AppId, request.Secret, request.Code)
			if err != nil {
				return nil, err
			}
			uSex = int64(userInfo.Sex)
			uName = userInfo.Nickname
			uHead = userInfo.HeadImgURL
			openID, err := usermodel.GetOpenIDByOpenIDAndAppID(userInfo.OpenID, request.AppId)
			if err != nil {
				return nil, err
			}
			uid = openID.UID
		}

		if uid == 0 {
			//该appID和openID第一次登录平台
			if userInfo.Unionid == "" {
				return nil, xlog.Error("UnionID为空")
			}
			//判断用户是否建立
			uInfo, err := usermodel.GetUserByUnionid(userInfo.Unionid)
			if err != nil {
				return nil, err
			}
			uid = uInfo.UID
			if uInfo.UID == 0 {
				//建立新用户
				no, err := sequence.UserNo.NewNo()
				if err != nil {
					return nil, xlog.Error("生成用户编号失败")
				}
				inviteCode, err := sequence.UserInviteCode.NewNo()
				if err != nil {
					xlog.ErrorP("生成邀请码失败")
					return nil, xlog.Error("生成邀请码失败")
				}
				uid, err = usermodel.CreateUserOuth(userInfo.Nickname, userInfo.HeadImgURL, no, inviteCode, userInfo.Unionid, request.AppId, userInfo.OpenID, request.Sign, userInfo.Sex, request.OuthType)
				if err != nil {
					return nil, xlog.Error(err)
				}
				isNew = true
			} else {
				//只需要建立新的openID记录
				err = usermodel.CreateOpenId(uInfo.UID, request.AppId, userInfo.OpenID, nil)
				if err != nil {
					return nil, xlog.Error(err)
				}
			}
		}
		userResultInfo, err = usermodel.GetUserByID(uid)
		if err != nil {
			return nil, xlog.Error(err)
		}
	} else if request.Type == enum.JpushType {
		lastType = enum.LoginTypeJpush
		var err error
		userResultInfo, err = usermodel.GetUserByPhone(request.Phone, "")
		if err != nil {
			return nil, xlog.Error(err)
		}
		if userResultInfo.UID == 0 {
			no, err := sequence.UserNo.NewNo()
			if err != nil {
				return nil, xlog.Error("生成用户编号失败")
			}
			inviteCode, err := sequence.UserInviteCode.NewNo()
			if err != nil {
				xlog.ErrorP("生成邀请码失败")
				return nil, xlog.Error("生成邀请码失败")
			}
			uid, err = usermodel.CreateUser(request.Phone, request.Phone, no, inviteCode, "", request.Sign, 0)
			if err != nil {
				return nil, xlog.Error(err)
			}
			isNew = true
			userResultInfo, err = usermodel.GetUserByID(uid)
			if err != nil {
				return nil, xlog.Error(err)
			}
		}
	} else if request.Type == enum.PhoneType {
		lastType = enum.LoginTypePhone
		var err error
		userResultInfo, err = usermodel.GetUserByPhone(request.Phone, "")
		if err != nil {
			return nil, xlog.Error(err)
		}
		if userResultInfo.UID == 0 {
			no, err := sequence.UserNo.NewNo()
			if err != nil {
				return nil, xlog.Error("生成用户编号失败")
			}
			inviteCode, err := sequence.UserInviteCode.NewNo()
			if err != nil {
				xlog.ErrorP("生成邀请码失败")
				return nil, xlog.Error("生成邀请码失败")
			}
			uid, err = usermodel.CreateUser(request.Phone, request.Phone, no, inviteCode, "", request.Sign, 0)
			if err != nil {
				return nil, xlog.Error(err)
			}
			isNew = true
			userResultInfo, err = usermodel.GetUserByID(uid)
			if err != nil {
				return nil, xlog.Error(err)
			}
		}
	} else if request.Type == enum.PhonePassType {
		lastType = enum.LoginTypePhonePass
		var err error
		userResultInfo, err = usermodel.GetUserByPhone(request.Phone, "")
		if err != nil {
			return nil, xlog.Error(err)
		}
		if userResultInfo.UID == 0 {
			return nil, xlog.Error("该手机没有注册")
		} else if userResultInfo.Pass != request.Pass {
			return nil, xlog.Error("密码错误")
		}
	}

	if request.Type == enum.OuthType {
		//第三方登录校验信息是否是最新的
		usermodel.UpdateUser(userResultInfo, uName, uHead, uSex)
	}

	//更新登录时间和登录设备
	usermodel.UpdateLoginInfo(request.Device, userResultInfo.UID, lastType)
	nickname = userResultInfo.NickName
	avatar = userResultInfo.Head
	no = userResultInfo.No
	sign = userResultInfo.Sign
	inviteCode = userResultInfo.InviteCode
	phone = userResultInfo.Phone
	gender = userResultInfo.Sex
	userResult.Name = nickname
	userResult.CancelStatus = userResultInfo.CancelStartTime != 0
	userResult.Sex = gender
	userResult.Head = avatar
	userResult.Uid = userResultInfo.UID
	userResult.InviteCode = inviteCode
	userResult.No = no
	userResult.Sign = sign
	userResult.Phone = phone
	userResult.Identity = userResultInfo.Identity
	userResult.HaveWx = userResultInfo.HaveWx
	userResult.HavePass = userResultInfo.Pass != ""
	userResult.IsNew = isNew
	userResult.City = userResultInfo.City
	userResult.Province = userResultInfo.Province
	userResult.Address = userResultInfo.Address
	userResult.Age = userResultInfo.Age
	return userResult, nil
}

func (rpc *rpcServer) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.UserDataResponse, error) {
	var err error
	if request.Type == enum.UpdateUserPhone {
		err = usermodel.UpdateUserPhone(request.Phone, request.Uid)
	} else if request.Type == enum.UpdateUserPass {
		err = usermodel.UpdateUserPass(request.Phone, request.Pass)
	} else if request.Type == enum.UpdateUserIdentity {
		list:=make([]int64,0)
		if len(request.UidList)!=0{
			list=append(list,request.UidList...)
		}else{
			list=append(list,request.Uid)
		}
		err = usermodel.UpdateUserIdentity(request.IdentityType,list)
	} else if request.Type == enum.UpdateUserModifyPass {
		err = usermodel.UpdateUserModifyPass(request.Pass, request.NewPass, request.Uid)
	} else if request.Type == enum.UpdateUserName {
		err = usermodel.UpdateUserName(request.Name, request.Uid)
	} else if request.Type == enum.UpdateUserHead {
		err = usermodel.UpdateUserHead(request.Head, request.Uid)
	} else if request.Type == enum.UpdateUserSex {
		err = usermodel.UpdateUserSex(request.Sex, request.Uid)
	} else if request.Type == enum.UpdateUserSign {
		err = usermodel.UpdateUserSign(request.Sign, request.Uid)
	} else if request.Type == enum.UpdateUserAddress {
		//err = usermodel.UpdateUserAddress(request.District, request.Uid)
	} else if request.Type == enum.UpdatePushInfo {
		err = usermodel.UpdatePushInfo(request.RegistrationId, request.Uid)
	} else if request.Type == enum.UpdateUserInfo {
		if request.Phone != "" {
			err = usermodel.UpdateUserPhone(request.Phone, request.Uid)
		} else {
			err = usermodel.UpdateUserSign(request.Sign, request.Uid)
		}
		err = usermodel.UpdateUserName(request.Name, request.Uid)
		err = usermodel.UpdateUserSex(request.Sex, request.Uid)
		err = usermodel.UpdateUserAge(request.Age, request.Uid)
		err = usermodel.UpdateUserAddress(request.Uid, request.District)
	}
	if err != nil {
		return nil, err
	}
	userResultInfo, err := usermodel.GetUserByID(request.Uid)
	if err != nil {
		return nil, err
	}
	userResult := &pb.UserDataResponse{}
	userResult.Name = userResultInfo.NickName
	userResult.Sex = userResultInfo.Sex
	userResult.Head = userResultInfo.Head
	userResult.Uid = userResultInfo.UID
	userResult.InviteCode = userResultInfo.InviteCode
	userResult.No = userResultInfo.No
	userResult.Sign = userResultInfo.Sign
	userResult.Phone = userResultInfo.Phone
	userResult.Identity = userResultInfo.Identity
	userResult.HaveWx = userResultInfo.HaveWx
	userResult.HavePass = userResultInfo.Pass != ""
	userResult.City = userResultInfo.City
	userResult.Province = userResultInfo.Province
	userResult.Address = userResultInfo.Address
	userResult.Age = userResultInfo.Age
	return userResult, err
}

func (rpc *rpcServer) Register(ctx context.Context, request *pb.LoginRequest) (*pb.UserDataResponse, error) {
	//判断用户是否建立
	uInfo, err := usermodel.GetUserByPhone(request.Phone, "")
	if err != nil {
		return nil, err
	}
	uid := uInfo.UID
	if uInfo.UID == 0 {
		//建立新用户
		no, err := sequence.UserNo.NewNo()
		if err != nil {
			return nil, xlog.Error("生成用户编号失败")
		}
		inviteCode, err := sequence.UserInviteCode.NewNo()
		if err != nil {
			xlog.ErrorP("生成邀请码失败")
			return nil, xlog.Error("生成邀请码失败")
		}

		uid, err = usermodel.CreateUser(request.Name, request.Phone, no, inviteCode, request.Pass, request.Sign, 0)
	} else {
		return nil, xlog.Error("该手机已经被使用")
	}
	userResultInfo, err := usermodel.GetUserByID(uid)
	if err != nil {
		return nil, err
	}
	userResult := &pb.UserDataResponse{}
	userResult.Name = userResultInfo.NickName
	userResult.Sex = userResultInfo.Sex
	userResult.Head = userResultInfo.Head
	userResult.Uid = uid
	userResult.InviteCode = userResultInfo.InviteCode
	userResult.No = userResultInfo.No
	userResult.Phone = userResultInfo.Phone
	userResult.HaveWx = userResultInfo.HaveWx
	userResult.HavePass = userResultInfo.Pass != ""
	userResult.IsNew = true
	userResult.Address = userResultInfo.Address
	userResult.Age = userResultInfo.Age
	return userResult, nil
}

func (rpc *rpcServer) PhoneCheck(ctx context.Context, request *pb.PhoneCheckRequest) (*pb.PhoneCheckResponse, error) {
	//判断用户是否建立
	uInfo, err := usermodel.GetUserByPhone(request.Phone, "")
	if err != nil {
		return nil, err
	}
	uid := uInfo.UID
	if uid == 0 {
		return &pb.PhoneCheckResponse{Result: 0}, nil
	} else if uInfo.Device == request.Device {
		return &pb.PhoneCheckResponse{Result: 1}, nil
	} else {
		return &pb.PhoneCheckResponse{Result: 2}, nil
	}
}

func (rpc *rpcServer) Address(ctx context.Context, request *pb.AddressRequest) (*pb.AddressResponse, error) {
	result := make([]*pb.AddressData, 0)
	resultList, err := addressmodel.GetAddress(request.Id, request.Type)
	if err != nil {
		return nil, err
	}
	for _, v := range resultList {
		list := make([]*pb.AddressData, 0)
		for _, vv := range v.List {
			list = append(list, &pb.AddressData{
				Id:     vv.ID,
				Name:   vv.Name,
				First:  vv.First,
				Letter: vv.Letter,
				Pid:    vv.Pid,
			})
		}
		result = append(result, &pb.AddressData{
			Id:     v.ID,
			Name:   v.Name,
			First:  v.First,
			Letter: v.Letter,
			Pid:    v.Pid,
			List:   list,
		})
	}
	return &pb.AddressResponse{List: result}, nil
}

func (rpc *rpcServer) ModifyOuth(ctx context.Context, request *pb.ModifyOuthRequest) (*pb.EmptyResponse, error) {
	if request.IsUnbind {
		return &pb.EmptyResponse{}, usermodel.Unbind(request.Uid, request.OuthType)
	} else {
		userInfo, err := wx.GetUserInfo(request.AppId, request.Secret, request.Code)
		if err != nil {
			return nil, err
		}
		if userInfo.Unionid == "" {
			return nil, xlog.Error("UnionID为空")
		}
		openID, err := usermodel.GetOpenIDByOpenIDAndAppID(userInfo.OpenID, request.AppId)
		if err != nil {
			return nil, err
		}
		if openID.UID == 0 {
			err = usermodel.CreateOpenId(request.Uid, request.AppId, userInfo.OpenID, nil)
			if err != nil {
				return nil, xlog.Error(err)
			}
		}
		//判断用户是否建立
		uInfo, err := usermodel.GetUserByUnionid(userInfo.Unionid)
		if err != nil {
			return nil, err
		}
		if uInfo.UID == 0 {
			err = usermodel.CreateOuth(request.Uid, request.OuthType, userInfo.Unionid, nil)
			if err != nil {
				return nil, xlog.Error(err)
			}
		}else{
			return nil,xlog.Error("该微信已经被绑定")
		}
	}

	return &pb.EmptyResponse{}, nil
}

func (rpc *rpcServer) FindIdFromQuery(ctx context.Context, request *pb.FindUserFromQueryRequest) (*pb.FindUserFromQueryResponse, error) {
	var list []usermodel.User
	var err error
	if request.No {
		list, err = usermodel.GetUserByNo(request.Query)
	} else if request.Phone {
		list, err = usermodel.GetUserFromPhone(request.Query)
	} else if request.Name {
		list, err = usermodel.GetUserByName(request.Query)
	}

	if err != nil {
		return nil, err
	}
	resultList := make([]int64, 0)
	for _, info := range list {
		resultList = append(resultList, info.UID)
	}
	return &pb.FindUserFromQueryResponse{List: resultList}, err
}

func (rpc *rpcServer) GetUserList(ctx context.Context, request *pb.GetUserListRequest) (*pb.GetUserListResponse, error) {
	result := &pb.GetUserListResponse{
		List: make([]*pb.UserDataResponse, 0),
	}
	list, total, err := usermodel.GetUserQuery(request.IdList, []int64{request.Status}, "", "", "", "", "", request.Page, request.Size, request.CancelType)
	if err == nil {
		for _, v := range list {
			result.List = append(result.List, &pb.UserDataResponse{
				Identity:       v.Identity,
				Uid:            v.UID,
				Head:           v.Head,
				Name:           v.NickName,
				Phone:          v.Phone,
				CancelTime:     v.CancelStartTime,
				Address:        v.Address,
				Age:            v.Age,
				RegistrationId: v.RegistrationId,
				LoginType:      v.LastType,
				LoginTime:      v.LoginTime,
			})
		}
		result.Total = total
	}

	return result, err
}

func (rpc *rpcServer) Cancel(ctx context.Context, request *pb.CancelRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, usermodel.CancelUser(request.Uid)
}

func (rpc *rpcServer) EnterCancel(ctx context.Context, request *pb.EnterCancelRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, usermodel.EnterCancel(request.List)
}

func (rpc *rpcServer) BatchModifyType(ctx context.Context, request *pb.BatchModifyRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, usermodel.BatchModifyType(request.List, request.Type)
}
