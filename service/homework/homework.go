package homework

import (
	"context"
	"mufe_service/camp/service"
	"mufe_service/camp/xlog"
	pb "mufe_service/jsonRpc"
	"mufe_service/model/homework"
)

func init() {
	nSer := &rpcServer{}
	pb.RegisterHomeWorkServiceServer(service.GetRegisterRpc(), nSer)

}

type rpcServer struct {
}

func (rpc *rpcServer) AddHomeWork(ctx context.Context, request *pb.AddHomeWorkRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, homeworkModel.AddHomeWork(request.ClassIds, request.Ids, request.Uid, request.Time, request.Desc)
}

func (rpc *rpcServer) HomeWorkList(ctx context.Context, request *pb.HomeWorkListRequest) (*pb.HomeWorkListResponse, error) {
	result, err := homeworkModel.HomeWorkList(request.Time, request.ClassId, request.Uid)
	list := make([]*pb.HomeWorkListData, 0)
	if err == nil {
		index := int64(0)
		for _, v := range result {
			xlog.Info(v.Progress)
			list = append(list, &pb.HomeWorkListData{
				Id:       v.ContentId,
				Title:    v.Title,
				Progress: v.Progress,
				Index:    index,
				Cover:    v.Cover,
				Number:   1,
				Time:     60,
			})
		}
	}
	return &pb.HomeWorkListResponse{List: list}, err
}

func (rpc *rpcServer) HomeWorkGroup(ctx context.Context, request *pb.HomeWorkListRequest) (*pb.HomeWorkGroupResponse, error) {
	result, err := homeworkModel.StudentHomeWorkList(request.Time, request.Uid, request.Id)
	list := make([]*pb.HomeWorkListData, 0)
	returnResult := &pb.HomeWorkGroupResponse{

	}
	if err == nil {
		returnResult.Desc = result.Desc
		index := int64(0)
		for _, v := range result.List {
			list = append(list, &pb.HomeWorkListData{
				Id:       v.ID,
				Title:    v.Title,
				Progress: v.Progress,
				Index:    index,
				Cover:    v.Cover,
			})
		}
		returnResult.List = list
		returnResult.GroupId = result.Id
	}
	return returnResult, err
}

func (rpc *rpcServer) HomeWorkRecord(ctx context.Context, request *pb.HomeWorkRecordRequest) (*pb.HomeWorkRecordResponse, error) {
	finish, inComplete := homeworkModel.GetHomeWorkRecord(request.Id)
	finishResult := make([]*pb.HomeWorkRecordData, 0)
	inCompleteResult := make([]*pb.HomeWorkRecordData, 0)
	for _, v := range finish {
		finishResult = append(finishResult, &pb.HomeWorkRecordData{
			Uid: v,
		})
	}
	for _, v := range inComplete {
		inCompleteResult = append(inCompleteResult, &pb.HomeWorkRecordData{
			Uid: v,
		})
	}
	return &pb.HomeWorkRecordResponse{
		Finish:     finishResult,
		Incomplete: inCompleteResult,
	}, nil
}

func (rpc *rpcServer) HomeWorkInfo(ctx context.Context, request *pb.HomeWorkRequest) (*pb.HomeWorkResponse, error) {
	result := &pb.HomeWorkResponse{
		List: make([]*pb.HomeWorkData, 0),
	}
	list, err := homeworkModel.GetHomeWorkInfo(request.Page, request.Size, request.Status)
	if err == nil {
		for _, v := range list {
			tagList := make([]*pb.TagData, 0)
			for _, vv := range v.TagList {
				tagList = append(tagList, &pb.TagData{
					Id:    vv.ID,
					Title: vv.Name,
				})
			}
			result.List = append(result.List, &pb.HomeWorkData{
				Id:    v.ID,
				Title: v.Title,
				Cover: v.Cover,
				Tag:   tagList,
			})
		}
	}
	return result, err
}

func (rpc *rpcServer) HomeWork(ctx context.Context, request *pb.HomeWorkRequest) (*pb.HomeWorkResponse, error) {
	result := &pb.HomeWorkResponse{
		List: make([]*pb.HomeWorkData, 0),
	}
	list, err := homeworkModel.GetHomeWork(request.Status, request.ContentId)
	if err == nil {
		for _, v := range list {
			tagList := make([]*pb.TagData, 0)
			for _, vv := range v.TagList {
				tagList = append(tagList, &pb.TagData{
					Id:    vv.ID,
					Title: vv.Name,
				})
			}
			result.List = append(result.List, &pb.HomeWorkData{
				Id:        v.ID,
				Title:     v.Title,
				Cover:     v.Cover,
				ContentId: v.ContentId,
				Tag:       tagList,
				Level:     v.Level,
				InfoId:v.InfoId,
			})
		}
	}
	return result, err
}

func (rpc *rpcServer) EditHomeWork(ctx context.Context, request *pb.EditHomeWorkRequest) (*pb.EditHomeWorkResponse, error) {
	id, err := homeworkModel.EditHomeWork(request.Title, request.Id,request.Level, request.Tag)
	return &pb.EditHomeWorkResponse{Id: id}, err
}

func (rpc *rpcServer) HomeWorkDetail(ctx context.Context, request *pb.HomeWorkDetailRequest) (*pb.HomeWorkData, error) {
	if request.InfoId!=0{
		detail, err := homeworkModel.GetHomeWorkDetailInfo(request.InfoId)
		if err == nil {
			tagList := make([]*pb.TagData, 0)
			for _, v := range detail.TagList {
				tagList = append(tagList, &pb.TagData{
					Id:    v.ID,
					Title: v.Name,
				})
			}
			return &pb.HomeWorkData{
				Content: detail.Content,
				Cover:   detail.Cover,
				Video:   detail.Video,
				Tag:     tagList,
				Level:detail.Level,
				Title:   detail.Title,
			}, nil
		}
		return nil, err
	}else{

		detail, err := homeworkModel.GetHomeWorkDetail(request.Id)
		if err == nil {
			tagList := make([]*pb.TagData, 0)
			for _, v := range detail.TagList {
				tagList = append(tagList, &pb.TagData{
					Id:    v.ID,
					Title: v.Name,
				})
			}
			return &pb.HomeWorkData{
				Id:      request.Id,
				Content: detail.Content,
				Cover:   detail.Cover,
				Video:   detail.Video,
				Tag:     tagList,
				ContentId:detail.ContentId,
				Title:   detail.Title,
			}, nil
		}
		return nil, err
	}


}

func (rpc *rpcServer) EditHomeworkInfo(ctx context.Context, request *pb.EditHomeWorkInfoRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, homeworkModel.EditHomeWorkInfo(request.Type, request.Id, request.Content, request.Prefix)
}

func (rpc *rpcServer) UnbindHomeWork(ctx context.Context, request *pb.UnBindHomeWorkRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, homeworkModel.UnBindHomeWork(request.Id)
}

func (rpc *rpcServer) FinishHomeWork(ctx context.Context, request *pb.FinishHomeWorkRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, homeworkModel.FinishHomeWork(request.Id, request.Uid, request.Score)
}

func (rpc *rpcServer) CancelHomeWork(ctx context.Context, request *pb.CancelHomeWorkRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, homeworkModel.CancelHomeWorkRecord(request.List)
}
