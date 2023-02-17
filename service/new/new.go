package new

import (
	"context"
	"mufe_service/camp/enum"
	"mufe_service/camp/service"
	pb "mufe_service/jsonRpc"
	"mufe_service/model/new"
)

type rpcServer struct {
}



func (rpc *rpcServer) GetNews(ctx context.Context, request *pb.GetNewsRequest) ( *pb.GetNewsResponse,  error) {
	list,err:=newModel.GetNews(request.Id,request.Status,request.Type)
	resultList:=make([]*pb.NewsData,0)
	if err==nil{
		for _,v:=range list{
			resultList=append(resultList,&pb.NewsData{
				Id:v.Id,
				Time:v.Time,
				Title:v.Title,
				Source:v.Source,
				Cover:v.Cover,
				Content:v.Content,
				Type:v.Type,
			})
		}
	}

	return &pb.GetNewsResponse{List:resultList}, err
}

func (rpc *rpcServer) EditNewType(ctx context.Context, request *pb.EditNewTypeRequest) ( *pb.EmptyResponse,  error) {
	return &pb.EmptyResponse{}, newModel.EditNewType(request.Id,request.Type)
}

func (rpc *rpcServer) EditNew(ctx context.Context, request *pb.EditNewRequest) ( *pb.EditNewResponse,  error) {
	id,err:= newModel.EditNew(request.Id,request.Title,request.Source)
	return &pb.EditNewResponse{Id:id},err
}

func (rpc *rpcServer) EditContent(ctx context.Context, request *pb.EditNewRequest) ( *pb.EmptyResponse,  error) {
	return &pb.EmptyResponse{}, newModel.EditNewContent(request.Id,request.Content)
}

func (rpc *rpcServer) EditCover(ctx context.Context, request *pb.EditNewCoverRequest) ( *pb.EmptyResponse,  error) {
	return &pb.EmptyResponse{}, newModel.EditNewCover(request.Id,request.Cover,request.Prefix)
}

func (rpc *rpcServer) DelNews(ctx context.Context, request *pb.EditNewRequest) ( *pb.EmptyResponse,  error) {
	return &pb.EmptyResponse{}, newModel.DelNew(request.Id,enum.StatusDelete)
}


func init() {
	nSer := &rpcServer{}
	pb.RegisterNewsServiceServer(service.GetRegisterRpc(), nSer)
}
