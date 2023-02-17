package recommend

import (
	"context"
	"mufe_service/camp/service"
	pb "mufe_service/jsonRpc"
	"mufe_service/model/course"
)

func init() {
	nSer := &rpcServer{}
	pb.RegisterRecommendServiceServer(service.GetRegisterRpc(), nSer)

}

type rpcServer struct {
}

func (rpc *rpcServer) GetRecommendList(ctx context.Context, request *pb.RecommendRequest) (*pb.RecommendResponse, error) {
	result:=&pb.RecommendResponse{
		List:make([]*pb.RecommendData,0),
	}
	list,err:=courseModel.GetRecommendList(request.Ids)
	if err==nil{
		for _,v := range list{
			result.List=append(result.List,&pb.RecommendData{
				Id:v.ID,
				ContentId:v.ContentId,
				ContentType:v.ContentType,
				InfoId:v.InfoId,
			})
		}
	}
	return result, err
}

func (rpc *rpcServer) GetAdminRecommendList(ctx context.Context, request *pb.RecommendRequest) (*pb.RecommendResponse, error) {
	result:=&pb.RecommendResponse{
		List:make([]*pb.RecommendData,0),
	}
	list,err:=courseModel.GetAdminRecommendList(request.Ids[0],request.Page,request.Size)
	if err==nil{
		for _,v := range list{
			result.List=append(result.List,&pb.RecommendData{
				Id:v.ID,
				ContentId:v.ContentId,
				ContentType:v.ContentType,
				InfoId:v.InfoId,
			})
		}
	}
	return result, err
}

func (rpc *rpcServer) GetRecommendInfoList(ctx context.Context, request *pb.RecommendInfoRequest) (*pb.RecommendInfoResponse, error) {
	result:=&pb.RecommendInfoResponse{
		List:make([]*pb.RecommendInfoData,0),
	}
	list,err:=courseModel.GetRecommendInfoList(request.Type)
	if err==nil{
		for _,v := range list{
			result.List=append(result.List,&pb.RecommendInfoData{
				Id:v.ID,
				Icon:v.Icon,
				Title:v.Name,
				Type:v.Type,
			})
		}
	}
	return result, err
}

func (rpc *rpcServer) EditRecommend(ctx context.Context, request *pb.EditRecommendRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, courseModel.EditRecommend(request.Id,request.Type,request.ContentId,request.ContentType,request.IsDel)
}
