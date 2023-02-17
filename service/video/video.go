package video

import (
	"context"
	"mufe_service/camp/enum"
	"mufe_service/camp/service"
	pb "mufe_service/jsonRpc"
	"mufe_service/model/course"
)

func init() {
	nSer := &rpcServer{}
	pb.RegisterVideoServiceServer(service.GetRegisterRpc(), nSer)

}

type rpcServer struct {
}

func (rpc *rpcServer) AddVideo(ctx context.Context, request *pb.AddVideoRequest) (*pb.AddVideoResponse, error) {
	return &pb.AddVideoResponse{}, nil
}

func (rpc *rpcServer) EditItem(ctx context.Context, request *pb.EditItemRequest) (*pb.EditItemResponse, error) {
	result,err:=courseModel.EditItem(request.Id,request.ChapterId,request.Title,request.TagId)
	return &pb.EditItemResponse{Id:result}, err
}


func (rpc *rpcServer) EditItemInfo(ctx context.Context, request *pb.EditItemInfoRequest) (*pb.EmptyResponse, error) {
	if request.Type==enum.EDIT_VIDEO{
		list:=make([]courseModel.AddVideoData,0)
		for _,v:=range request.List{
			list=append(list,courseModel.AddVideoData{
				Cover:v.Cover,
				CoverPrefix:v.CoverPrefix,
				UrlPrefix:v.UrlPrefix,
				Url:v.Url,
				DownUrlPrefix:v.DownUrlPrefix,
				DownUrl:v.DownUrl,
				Duration:v.Duration,
				ID:request.Id,
			})
		}
		return &pb.EmptyResponse{}, courseModel.AddVideo(list)
	}else{
		return &pb.EmptyResponse{}, courseModel.EditItemInfo(request.Content,request.Prefix,request.Id,request.Type)
	}
}

func (rpc *rpcServer) GetAdminItem(ctx context.Context, request *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	result,err:=courseModel.GetAdminItem(request.Id)
	if err!=nil{
		return nil,err
	}

	list:=make([]*pb.VideoData,0)
	for _,v:=range result{
		tagList:=make([]*pb.TagData,0)
		for _,vv:=range v.TagList{
			tagList=append(tagList,&pb.TagData{
				Id:vv.ID,
				Title:vv.Name,
			})
		}
		list=append(list,&pb.VideoData{
			Id:v.ID,
			Title:v.Title,
			Content:v.Content,
			Cover:v.Cover,
			Url:v.Url,
			Tag:tagList,
		})
	}



	return &pb.GetItemResponse{Data:list},nil
}


func (rpc *rpcServer) VideoList(ctx context.Context, request *pb.VideoRequest) (*pb.VideoResponse, error) {
	result:=&pb.VideoResponse{}
	videoList,err:=courseModel.GetVideoList(request.Page,request.Pagesize)
	if err!=nil{
		return nil,err
	}
	for _,v:=range videoList{
		result.VideoList=append(result.VideoList,&pb.VideoData{
			Id:v.ID,
			Duration:v.Duration,
			Url:v.Url,
			DownUrl:v.DownUrl,
			Cover:v.Cover,
			Title:v.Title,
		})
	}
	return result, err
}

func (rpc *rpcServer) HistoryVideoList(ctx context.Context, request *pb.VideoRequest) (*pb.VideoResponse, error) {
	result:=&pb.VideoResponse{}
	videoList,err:=courseModel.GetHistoryVideoList(request.Uid)
	if err!=nil{
		return nil,err
	}
	for _,v:=range videoList{
		result.VideoList=append(result.VideoList,&pb.VideoData{
			Id:v.ID,
			Duration:v.Duration,
			Url:v.Url,
			DownUrl:v.DownUrl,
			Cover:v.Cover,
			Title:v.Title,
		})
	}
	return result, err
}

func (rpc *rpcServer) DelHistoryVideo(ctx context.Context, request *pb.DelVideoHistoryRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{},courseModel.DelHistoryVideo(request.Uid,request.Id)
}

func (rpc *rpcServer) AddHistoryVideo(ctx context.Context, request *pb.AddVideoHistoryRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, courseModel.AddHistoryVideo(request.Id,request.Uid)
}
func (rpc *rpcServer) DelChapterVideo(ctx context.Context, request *pb.VideoRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, courseModel.DelChapterVideo(request.VideoId)
}


func (rpc *rpcServer) GetVideo(ctx context.Context, request *pb.VideoRequest) (*pb.VideoResponse, error) {
	result:=&pb.VideoResponse{}
	videoList,err:=courseModel.GetVideo([]int64{request.VideoId},request.TagId,request.Key)
	if err!=nil{
		return nil,err
	}

	for _,v:=range videoList{
		tagList:=make([]*pb.TagData,0)
		for _,vv:=range v.TagList{
			tagList=append(tagList,&pb.TagData{
				Id:vv.ID,
				Title:vv.Name,
			})
		}
		result.VideoList=append(result.VideoList,&pb.VideoData{
			Id:v.ID,
			Duration:v.Duration,
			Url:v.Url,
			DownUrl:v.DownUrl,
			Cover:v.Cover,
			Title:v.Title,
			Content:v.Content,
			Tag:tagList,
			ChapterId:v.ChapterId,
		})
	}
	return result, err
}



func (rpc *rpcServer) EditVideoSort(ctx context.Context, request *pb.VideoRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, courseModel.EditVideoSort(request.VideoId,request.Sort)
}

func (rpc *rpcServer) EditVideoCover(ctx context.Context, request *pb.VideoRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, courseModel.EditVideo(request.Cover,request.Prefix,request.VideoId)
}
