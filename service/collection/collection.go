package course

import (
	"context"
	"mufe_service/camp/service"
	"mufe_service/camp/xlog"
	pb "mufe_service/jsonRpc"
	"mufe_service/model/collection"
	courseModel "mufe_service/model/course"
)

func init() {
	nSer := &rpcServer{}
	pb.RegisterCollectionServiceServer(service.GetRegisterRpc(), nSer)

}

type rpcServer struct {
}

func (rpc *rpcServer) IsCollection(ctx context.Context, request *pb.CollectionServiceRequest) (*pb.CollectionServiceResponse, error) {
	result:=&pb.CollectionServiceResponse{Collection:collectionModel.IsCollection(request.ContentId,request.Type,request.Uid)}
	return result, nil
}


func (rpc *rpcServer) EditCollection(ctx context.Context, request *pb.CollectionServiceRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, collectionModel.EditCollection(request.ContentId,request.Type,request.Uid,request.Del)
}

func (rpc *rpcServer) GetCollection(ctx context.Context, request *pb.GetCollectionRequest) (*pb.GetCollectionResponse, error) {
	result:=&pb.GetCollectionResponse{}
	xlog.Info(request)
	ids,err:=collectionModel.GetCollection(request.Uid,request.Type)
	if err==nil&&len(ids)>0{
		list,err:=courseModel.GetVideo(ids,0,"")
		if err==nil{
			for _,v := range list{
				tagList:=make([]*pb.TagData,0)
				for _,vv:=range v.TagList{
					tagList=append(tagList,&pb.TagData{
						Id:vv.ID,
						Title:vv.Name,
					})
				}
				result.List=append(result.List,&pb.ChapterData{
					Id:v.ID,
					Title:v.Title,
					Cover:v.Cover,
					Tag:tagList,
				})
			}
		} else {
			xlog.ErrorP(err)
		}
	}
	return result,err
}
