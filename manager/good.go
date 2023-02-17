package manager

import (
"os"
"mufe_service/camp/utils"
pb "mufe_service/jsonRpc"
)
var (
	goodCategoryService   pb.GoodCategoryServiceClient
	goodDeliveryService   pb.GoodDeliveryServiceClient
	goodService   pb.GoodServiceClient
	orderService   pb.OrderServiceClient
	wxPayService   pb.PayServiceClient
)

func GetGoodCategoryService() pb.GoodCategoryServiceClient {
	if goodCategoryService == nil {
		rpc, _ := utils.GetRPCService("good_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		goodCategoryService = pb.NewGoodCategoryServiceClient(rpc)
	}
	return goodCategoryService
}


func GetGoodDeliveryService() pb.GoodDeliveryServiceClient {
	if goodDeliveryService == nil {
		rpc, _ := utils.GetRPCService("good_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		goodDeliveryService = pb.NewGoodDeliveryServiceClient(rpc)
	}
	return goodDeliveryService
}

func GetGoodService() pb.GoodServiceClient {
	if goodService == nil {
		rpc, _ := utils.GetRPCService("good_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		goodService = pb.NewGoodServiceClient(rpc)
	}
	return goodService
}



func GetOrderService() pb.OrderServiceClient {
	if orderService == nil {
		rpc, _ := utils.GetRPCService("order_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		orderService = pb.NewOrderServiceClient(rpc)
	}
	return orderService
}




func GetPayService() pb.PayServiceClient {
	if wxPayService == nil {
		rpc, _ := utils.GetRPCService("pay_service", os.Getenv("CONSUL_TAG"), os.Getenv("CONSUL_IP"))
		wxPayService = pb.NewPayServiceClient(rpc)
	}
	return wxPayService
}


