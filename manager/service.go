package manager

import (
	"os"
	"mufe_service/camp/utils"
	pb "mufe_service/jsonRpc"
)

var (
	sendService         pb.SendSmsServiceClient
	footBallService         pb.FootballServiceClient
	liveService         pb.LiveServiceClient
	pushService         pb.PushServiceClient
)


func GetSendSmsService() pb.SendSmsServiceClient {
	if sendService == nil {
		// rpc, _ := core.GetRPCServiceWithTarget("127.0.0.1:9101")
		rpc, _ := utils.GetRPCService("system_aliyun_sms_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		sendService = pb.NewSendSmsServiceClient(rpc)
	}
	return sendService
}

func GetFootBallService() pb.FootballServiceClient {
	if footBallService == nil {
		// rpc, _ := core.GetRPCServiceWithTarget("127.0.0.1:9101")
		rpc, _ := utils.GetRPCService("system_aliyun_sms_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		footBallService = pb.NewFootballServiceClient(rpc)
	}
	return footBallService
}


func GetPushService() pb.PushServiceClient {
	if pushService == nil {
		// rpc, _ := core.GetRPCServiceWithTarget("127.0.0.1:9101")
		rpc, _ := utils.GetRPCService("system_push_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		pushService = pb.NewPushServiceClient(rpc)
	}
	return pushService
}


func GetLiveService() pb.LiveServiceClient {
	if liveService == nil {
		rpc, _ := utils.GetRPCService("live_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		liveService = pb.NewLiveServiceClient(rpc)
	}
	return liveService
}
