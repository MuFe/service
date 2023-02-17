package manager

import (
	"os"
	"mufe_service/camp/utils"
	pb "mufe_service/jsonRpc"
)

var (
	coachService pb.CoachServiceClient
)

func GetCoachService() pb.CoachServiceClient {
	if coachService == nil {
		rpc, _ := utils.GetRPCService("coach_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		coachService = pb.NewCoachServiceClient(rpc)
	}
	return coachService
}
