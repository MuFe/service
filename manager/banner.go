package manager

import (
	"os"
	"mufe_service/camp/utils"
	pb "mufe_service/jsonRpc"
)
var (
	bannerService   pb.AdServiceClient
)

func GetBannerService() pb.AdServiceClient {
	if bannerService == nil {
		rpc, _ := utils.GetRPCService("banner_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		bannerService = pb.NewAdServiceClient(rpc)
	}
	return bannerService
}
