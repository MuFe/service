package manager

import (
	"os"
	"mufe_service/camp/utils"
	pb "mufe_service/jsonRpc"
)
var (
	qiniuService   pb.QiniuServiceClient
	videoService   pb.VideoServiceClient
	adminUerService   pb.AdminUserServiceClient
	webService   pb.WebServiceClient
	newsService   pb.NewsServiceClient
	permissionService pb.AdminPermissionServiceClient
	recommendService   pb.RecommendServiceClient
)

func GetQiniuService() pb.QiniuServiceClient {
	if qiniuService == nil {
		rpc, _ := utils.GetRPCService("qiniu_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		qiniuService = pb.NewQiniuServiceClient(rpc)
	}
	return qiniuService
}

func GetVideoService() pb.VideoServiceClient {
	if videoService == nil {
		rpc, _ := utils.GetRPCService("course_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		videoService = pb.NewVideoServiceClient(rpc)
	}
	return videoService
}

func GetAdminUserService() pb.AdminUserServiceClient {
	if adminUerService == nil {
		rpc, _ := utils.GetRPCService("admin_user_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		adminUerService = pb.NewAdminUserServiceClient(rpc)
	}
	return adminUerService
}

func GetWebService() pb.WebServiceClient {
	if webService == nil {
		rpc, _ := utils.GetRPCService("banner_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		webService = pb.NewWebServiceClient(rpc)
	}
	return webService
}

func GetNewsService() pb.NewsServiceClient {
	if newsService == nil {
		rpc, _ := utils.GetRPCService("new_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		newsService = pb.NewNewsServiceClient(rpc)
	}
	return newsService
}

func GetRecommendService() pb.RecommendServiceClient {
	if recommendService == nil {
		rpc, _ := utils.GetRPCService("course_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		recommendService = pb.NewRecommendServiceClient(rpc)
	}
	return recommendService
}

func GetPermissionService() pb.AdminPermissionServiceClient {
	if permissionService == nil {
		rpc, _ := utils.GetRPCService("admin_user_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		permissionService = pb.NewAdminPermissionServiceClient(rpc)
	}
	return permissionService
}
