package user

import (
	"context"
	"mufe_service/camp/errcode"
	"mufe_service/camp/service"
	pb "mufe_service/jsonRpc"
	"mufe_service/model/adminUser"
)

func init() {
	nSer := &rpcServer{}
	pb.RegisterAdminUserServiceServer(service.GetRegisterRpc(), nSer)
	pb.RegisterAdminPermissionServiceServer(service.GetRegisterRpc(), nSer)
}

type rpcServer struct {
}

func (*rpcServer) Login(ctx context.Context, request *pb.AdminLoginRequest) (*pb.AdminUserDataResponse, error) {
	result, err := adminUserModel.Login(request.Phone, request.Pass)
	if err != nil {
		return nil, err
	}

	return &pb.AdminUserDataResponse{
		Uid:   result.UID,
		Head:  result.Head,
		Phone: result.Phone,
		Name:  result.NickName,
	}, nil
}

func (*rpcServer) BrandLogin(ctx context.Context, request *pb.AdminLoginRequest) (*pb.AdminUserDataResponse, error) {
	result, err := adminUserModel.Login(request.Phone, request.Pass)
	if err != nil {
		return nil, err
	}

	//获取登录使用的用户信息
	userInfo, err := adminUserModel.GetLoginUserInfo(result.UID, request.BgID)
	if err != nil {
		return nil, err
	}
	return &pb.AdminUserDataResponse{
		Uid:           result.UID,
		Head:          result.Head,
		Phone:         result.Phone,
		Name:          result.NickName,
		BusinessId:    userInfo.BusinessId,
		No:            userInfo.No,
		UserGroupName: userInfo.UserGroupName,
		UserGroupId:   userInfo.UserGroupId,
		WithdrawPass:  userInfo.WithdrawPass,
		BusinessPhoto:userInfo.BusinessPhoto,
		BusinessName:userInfo.BusinessName,
	}, nil
}



func (*rpcServer) DeleteUserPermission(ctx context.Context, request *pb.DeleteUserPermissionRequest) (result *pb.EmptyResponse, err error) {
	return &pb.EmptyResponse{},adminUserModel.DeleteUserPermission(request.Uid,request.BusinessId,request.GroupId)
}

func (*rpcServer) GetPermissionGroup(ctx context.Context, request *pb.GetPermissionGroupRequest) (result *pb.GetPermissionGroupResponse, err error) {
	roleList:=adminUserModel.GetPermissionGroup(request.BusinessGroupId,request.BusinessId)
	list:=make([]*pb.GetPermissionGroup,0)
	for _,info:=range roleList{
		user:=make([]*pb.GetPermissionUser,0)
		for _,value:=range info.User{
			user=append(user,&pb.GetPermissionUser{
				Uid:value.Uid,
				Name:value.Name,
				Phone:value.Phone,
			})
		}
		permissionList:=make([]*pb.GetPermission,0)
		for _,value:=range info.Role{
			permissionList=append(permissionList,&pb.GetPermission{
				Id:value.Id,
				Name:value.Name,
				Status:value.Status,
			})
		}
		list=append(list,&pb.GetPermissionGroup{
			Id:info.Id,
			Name:info.Name,
			User:user,
			Role:permissionList,
		})
	}
	return &pb.GetPermissionGroupResponse{List:list},nil
}

func (*rpcServer) GetUserPagePermission(ctx context.Context, request *pb.GetUserPagePermissionRequest) (result *pb.PagePermissionResponse, err error) {
	roleList:=adminUserModel.GetUserPagePermission(request.Uid,request.BusinessGroupId,request.BusinessId)
	return &pb.PagePermissionResponse{PermissionName:roleList},nil
}

func (*rpcServer) CheckUserPermission(ctx context.Context, request *pb.CheckGroupPermissionRequest) (result *pb.EmptyResponse, err error) {
	havePermission:=adminUserModel.CheckGroupPermission(request.Uid,request.UserGroupID,request.BusinessId,request.RoleName)
	if !havePermission{
		return &pb.EmptyResponse{},errcode.CommErrorUnauthorized.RPCError()
	}else {
		return &pb.EmptyResponse{},nil
	}
}

func (*rpcServer) CheckUserGroupPermission(ctx context.Context, request *pb.CheckUserGroupPermissionRequest) (result *pb.EmptyResponse, err error) {
	result = &pb.EmptyResponse{}
	return result, nil
}
func (rpc *rpcServer) GetBrandGoodCategory(ctx context.Context, request *pb.GetBrandGoodCategoryRequest) (*pb.GetBrandGoodCategoryResponse, error) {
	result, err := adminUserModel.GetBusinessCategoryByBusinessID([]int64{request.BusinessId})
	if err != nil {
		return nil, err
	}
	list := make([]int64, 0)
	for _, info := range result {
		list = append(list, info.CategoryID)
	}
	return &pb.GetBrandGoodCategoryResponse{List: list}, nil
}

//获取承诺服务列表
func (rpc *rpcServer) GetCommitment(context.Context, *pb.EmptyRequest) (*pb.GetCommitmentResponse, error) {
	result := &pb.GetCommitmentResponse{}
	result.List = make([]*pb.Commitment, 0)
	dataList, err := adminUserModel.GetCommitment()
	if err != nil {
		return nil, err
	}
	for _, v := range dataList {
		result.List = append(result.List, &pb.Commitment{
			Id:     v.Id,
			Name:   v.Name,
			Status: v.Status,
		})
	}
	return result, nil
}
