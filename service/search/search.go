package search

import (
	"context"
	"mufe_service/camp/service"
	pb "mufe_service/jsonRpc"
	searchModel "mufe_service/model/search"
)

func init() {
	nSer := &rpcServer{}
	pb.RegisterSearchServiceServer(service.GetRegisterRpc(), nSer)

}

type rpcServer struct {
}

func (rpc *rpcServer) AddSearch(ctx context.Context, request *pb.SearchRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, searchModel.AddSearch(request.Uid, request.Content)
}

func (rpc *rpcServer) GetSearchHistory(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	list,hot:=searchModel.GetSearchHistory(request.Uid)
	resultList:=make([]*pb.SearchData,0)
	resultHot:=make([]*pb.SearchData,0)
	for _,v:=range list{
		resultList=append(resultList,&pb.SearchData{
			Content:v.Content,
			Number:v.Number,
			TodayNumber:v.TodayNumber,
		})
	}
	for _,v:=range hot{
		resultHot=append(resultHot,&pb.SearchData{
			Content:v.Content,
			Number:v.Number,
			TodayNumber:v.TodayNumber,
		})
	}
	return &pb.SearchResponse{
		List:resultList,
		Hot:resultHot,
	}, nil
}
func (rpc *rpcServer) SearchHint(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	list:=searchModel.GetSearchHint(request.Content)
	resultList:=make([]*pb.SearchData,0)
	for _,v:=range list{
		resultList=append(resultList,&pb.SearchData{
			Content:v.Content,
			Number:v.Number,
			TodayNumber:v.TodayNumber,
		})
	}
	return &pb.SearchResponse{
		List:resultList,
	}, nil
}
