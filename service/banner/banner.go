package qiniu

import (
	"context"
	"mufe_service/camp/service"
	pb "mufe_service/jsonRpc"
	"mufe_service/model/banner"
)

type rpcServer struct {
}



func (rpc *rpcServer) GetAds(ctx context.Context, request *pb.AdServiceRequest) (result *pb.AdServicesResponse, err error) {
	result = &pb.AdServicesResponse{}
	list,err:=bannerModel.GetBanner(request.Status,request.Id)
	for _,v := range list{
		result.Result = append(result.Result, &pb.AdServiceResponse{
			Photo:v.Prefix+v.Photo,
			Type:v.TypeInt,
			LinkId:v.ContentId,
			Url:v.Url,
			Sort:v.Sort,
			Id:v.Id,
		})
	}
	return result, err
}

func (rpc *rpcServer) EditAd(ctx context.Context, request *pb.EditAdRequest) (*pb.EditAdResponse, error) {
	var err error
	id:=int64(0)
	if request.Id!=0{
		err=bannerModel.EditAd(request.Type,request.LinkId,request.Id,request.Url)
		id=request.Id
	} else {
		id,err=bannerModel.AddAd(request.Type,request.LinkId,request.Url)
	}
	return &pb.EditAdResponse{Id:id}, err
}

func (rpc *rpcServer) EditAdSort(ctx context.Context, request *pb.EditAdRequest) (result *pb.EmptyResponse, err error) {
	return &pb.EmptyResponse{}, bannerModel.EditSort(request.Sort,request.Id)
}

func (rpc *rpcServer) DelAd(ctx context.Context, request *pb.EditAdRequest) (result *pb.EmptyResponse, err error) {
	return &pb.EmptyResponse{}, bannerModel.DelAd(request.Id)
}

func (rpc *rpcServer) EditAdPhoto(ctx context.Context, request *pb.EditAdRequest) (result *pb.EmptyResponse, err error) {
	return &pb.EmptyResponse{}, bannerModel.EditPhoto(request.Id,request.Key,request.Prefix)
}

func init() {
	nSer := &rpcServer{}
	pb.RegisterAdServiceServer(service.GetRegisterRpc(), nSer)
}
