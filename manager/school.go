package manager

import (
	"os"
	"mufe_service/camp/utils"
	pb "mufe_service/jsonRpc"
)
var (
	schoolService   pb.SchoolServiceClient
)

func GetSchoolService() pb.SchoolServiceClient {
	if schoolService == nil {
		rpc, _ := utils.GetRPCService("school_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		schoolService = pb.NewSchoolServiceClient(rpc)
	}
	return schoolService
}
