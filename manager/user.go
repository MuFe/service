package manager

import (
	"os"
	"mufe_service/camp/utils"
	pb "mufe_service/jsonRpc"
)

var (
	userService     pb.UserServiceClient
	feedbackService pb.FeedbackServiceClient
)

func GetUserService() pb.UserServiceClient {
	if userService == nil {
		rpc, _ := utils.GetRPCService("user_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		userService = pb.NewUserServiceClient(rpc)
	}
	return userService
}

func GetFeedbackService() pb.FeedbackServiceClient {
	if feedbackService == nil {
		rpc, _ := utils.GetRPCService("user_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		feedbackService = pb.NewFeedbackServiceClient(rpc)
	}
	return feedbackService
}
