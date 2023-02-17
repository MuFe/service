package manager

import (
	"os"
	"mufe_service/camp/utils"
	pb "mufe_service/jsonRpc"
)

var (
	homeworkService     pb.HomeWorkServiceClient
)

func GetHomeWorkService() pb.HomeWorkServiceClient {
	if homeworkService == nil {
		rpc, _ := utils.GetRPCService("homework_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		homeworkService = pb.NewHomeWorkServiceClient(rpc)
	}
	return homeworkService
}


