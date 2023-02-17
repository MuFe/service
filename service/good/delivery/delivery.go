package delivery

import (
	"context"
	"mufe_service/camp/service"
	pb "mufe_service/jsonRpc"
	goodmodel "mufe_service/model/good"
)

func init() {
	pb.RegisterGoodDeliveryServiceServer(service.GetRegisterRpc(), &rpcServer{})
}

type rpcServer struct{}

func (rpc *rpcServer) GetDeliveryTemplate(ctx context.Context, request *pb.DeliveryListRequest) (*pb.DeliveryListResponse, error) {
	list, err := goodmodel.GetDeliveryTemplate(request.BusinessId, request.SpuId, request.Id)
	if err != nil {
		return nil, err
	}
	response := &pb.DeliveryListResponse{}
	for _, info := range list {
		response.List = append(response.List, &pb.DeliveryData{
			Id:                info.Id,
			Name:              info.Name,
			DefaultPrice:      info.DefaultPrice,
			DefaultNum:        info.DefaultNum,
			IncreaseNum:       info.IncreaseNum,
			IncreasePrice:     info.IncreasePrice,
		})
	}
	return response, nil
}

func (rpc *rpcServer) EditDeliveryTemplate(ctx context.Context, request *pb.EditDeliveryListRequest) (response *pb.DeliveryListResponse, err error) {
	id, err := goodmodel.EditDeliveryTemplate(request.Data.Id,
		request.Data.DefaultNum,
		request.Data.DefaultPrice,
		request.Data.IncreaseNum,
		request.Data.IncreasePrice,
		request.BusinessId,
		request.Uid,
		request.Data.Name)
	if err != nil {
		return nil, err
	}
	list := make([]*pb.DeliveryData, 1)
	list[0] = &pb.DeliveryData{
		Id:                id,
		Name:              request.Data.Name,
		DefaultNum:        request.Data.DefaultNum,
		DefaultPrice:      request.Data.DefaultPrice,
		IncreaseNum:       request.Data.IncreaseNum,
		IncreasePrice:     request.Data.IncreasePrice,
	}
	return &pb.DeliveryListResponse{List: list}, nil
}

func (rpc *rpcServer) DelDeliveryTemplate(ctx context.Context, request *pb.DelDeliveryRequest) (response *pb.EmptyResponse, err error) {
	return &pb.EmptyResponse{}, goodmodel.DelDeliveryTemplate(request.Id, request.BusinessId, request.Uid)
}
