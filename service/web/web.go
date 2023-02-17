package web

import (
	"context"
	"mufe_service/camp/service"
	pb "mufe_service/jsonRpc"
	"mufe_service/model/web"
)

type rpcServer struct {
}



func (rpc *rpcServer) GetWebInfo(ctx context.Context, request *pb.GetWebInfoRequest) ( *pb.GetWebInfoResponse,  error) {
	result:= &pb.GetWebInfoResponse{
		Content:make(map[int64]string),
	}
	list,err:=webModel.GetWebValue(request.List)
	if err!=nil{
		return nil,err
	}
	for _,v:=range list{
		result.Content[v.Type]=v.Content
	}
	return result, nil
}

func (rpc *rpcServer) ContactUs(ctx context.Context, request *pb.ContactUsRequest) ( *pb.EmptyResponse,  error) {
	var err error
	if request.Id==0{
		err=webModel.AddWebContactUs(request.Name,request.Phone,request.Email,request.Content)
	} else {
		err=webModel.EditWebContactUsStatus(request.Id,request.Status)
	}
	if err!=nil{
		return nil,err
	}
	return &pb.EmptyResponse{}, nil
}



func init() {
	nSer := &rpcServer{}
	pb.RegisterWebServiceServer(service.GetRegisterRpc(), nSer)
}
