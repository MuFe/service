package chapter

import (
	"context"
	"mufe_service/camp/enum"
	"mufe_service/camp/service"
	"mufe_service/camp/xlog"
	pb "mufe_service/jsonRpc"
	"mufe_service/model/course"
	homeworkModel "mufe_service/model/homework"
)

func init() {
	nSer := &rpcServer{}
	pb.RegisterChapterServiceServer(service.GetRegisterRpc(), nSer)

}

type rpcServer struct {
}

func (rpc *rpcServer) GetChapter(ctx context.Context, request *pb.ChapterServiceRequest) (*pb.ChapterServiceResponse, error) {
	result := &pb.ChapterServiceResponse{}
	list, err := courseModel.GetChapterFromID(request.Ids, request.Page, request.Size, request.CourseId, request.Status)
	if err == nil {
		for _, v := range list {
			tagList := make([]*pb.TagData, 0)
			for _, vv := range v.TagList {
				tagList = append(tagList, &pb.TagData{
					Id:    vv.ID,
					Title: vv.Name,
				})
			}
			result.List = append(result.List, &pb.ChapterData{
				Id:      v.ID,
				Title:   v.Title,
				Desc:    v.Desc,
				Level:   v.Level,
				Cover:   v.Cover,
				Tag:     tagList,
				Price:   v.Price,
				LevelId: v.LevelId,
				User:    v.User,
				Section: v.Section,
				Time:    v.Duration,
			})
		}
	}
	return result, err
}

func (rpc *rpcServer) GetAdminChapter(ctx context.Context, request *pb.ChapterServiceRequest) (*pb.ChapterServiceResponse, error) {
	result := &pb.ChapterServiceResponse{}
	list, err := courseModel.GetAdminChapter(request.Ids, request.Page, request.Size, request.CourseId, request.Status)
	if err == nil {
		for _, v := range list {
			result.List = append(result.List, &pb.ChapterData{
				Id:    v.ID,
				Title: v.Title,
				Cover: v.Cover,
			})
		}
	}
	return result, err
}

func (rpc *rpcServer) GetAdminChapterDetail(ctx context.Context, request *pb.ChapterServiceRequest) (*pb.ChapterServiceResponse, error) {
	result := &pb.ChapterServiceResponse{}
	list, err := courseModel.GetAdminChapterDetail(request.Ids[0])
	if err == nil {
		homeWorkList := make([]int64, 0)
		homeWorkInfo, err := homeworkModel.GetHomeWork(enum.StatusNormal, request.Ids)
		if err == nil {
			for _, v := range homeWorkInfo {
				homeWorkList = append(homeWorkList, v.InfoId)
			}
		}
		for _, v := range list {
			result.List = append(result.List, &pb.ChapterData{
				Id:       v.ID,
				Title:    v.Title,
				Cover:    v.Cover,
				Desc:     v.Desc,
				Plan:     v.Plan,
				Video:    v.Video,
				Homework: homeWorkList,
			})
		}
	}
	return result, err
}

func (rpc *rpcServer) GetChapterWithVideo(ctx context.Context, request *pb.ChapterVideoServiceRequest) (*pb.ChapterVideoServiceResponse, error) {
	result := &pb.ChapterVideoServiceResponse{
		List: make([]*pb.ChapterData, 0),
	}
	m := make(map[int64]*pb.ChapterData)
	list, err := courseModel.GetChapterFromID(request.ChapterId, request.Page, request.Size, request.CourseId, request.Status)
	if err == nil {
		if len(list) == 0 {
			return nil, xlog.Error("参数有误")
		}
		idList := make([]int64, 0)
		for _, v := range list {
			tagList := make([]*pb.TagData, 0)
			for _, vv := range v.TagList {
				tagList = append(tagList, &pb.TagData{
					Id:    vv.ID,
					Title: vv.Name,
				})
			}
			temp := &pb.ChapterData{
				Id:      v.ID,
				Title:   v.Title,
				Desc:    v.Desc,
				Level:   v.Level,
				Cover:   v.Cover,
				Tag:     tagList,
				Price:   v.Price,
				Plan:v.Plan,
				LevelId: v.LevelId,
				User:    v.User,
				Section: v.Section,
				Time:    v.Duration,
			}
			result.List = append(result.List, temp)
			m[v.ID] = temp
			idList = append(idList, v.ID)
		}
		videoListMap, err := courseModel.GetChapterVideo(idList)
		if err != nil {
			return nil, err
		}
		for k, vv := range videoListMap {
			temp, ok := m[k]
			if ok {
				for _, v := range vv {
					temp.VideoList = append(temp.VideoList, &pb.VideoData{
						Id:       v.ID,
						VideoId:  v.VideoId,
						Duration: v.Duration,
						Url:      v.Url,
						DownUrl:  v.DownUrl,
						Cover:    v.Cover,
						Title:    v.Title,
					})
				}
			}
		}

	}
	return result, err
}

func (rpc *rpcServer) GetAdminChapterWithVideo(ctx context.Context, request *pb.ChapterVideoServiceRequest) (*pb.ChapterVideoServiceResponse, error) {
	result := &pb.ChapterVideoServiceResponse{
		List: make([]*pb.ChapterData, 0),
	}
	m := make(map[int64]*pb.ChapterData)
	list, err := courseModel.GetChapterFromID(request.ChapterId, request.Page, request.Size, request.CourseId, request.Status)
	if err == nil {
		if len(list) == 0 {
			return nil, xlog.Error("参数有误")
		}
		idList := make([]int64, 0)
		for _, v := range list {
			tagList := make([]*pb.TagData, 0)
			for _, vv := range v.TagList {
				tagList = append(tagList, &pb.TagData{
					Id:    vv.ID,
					Title: vv.Name,
				})
			}
			temp := &pb.ChapterData{
				Id:      v.ID,
				Title:   v.Title,
				Desc:    v.Desc,
				Level:   v.Level,
				Cover:   v.Cover,
				Tag:     tagList,
				Price:   v.Price,
				LevelId: v.LevelId,
				User:    v.User,
				Section: v.Section,
				Time:    v.Duration,
			}
			result.List = append(result.List, temp)
			m[v.ID] = temp
			idList = append(idList, v.ID)
		}
		videoListMap, err := courseModel.GetAdminChapterVideo(idList)
		if err != nil {
			return nil, err
		}
		for k, vv := range videoListMap {
			temp, ok := m[k]
			if ok {
				for _, v := range vv {
					temp.VideoList = append(temp.VideoList, &pb.VideoData{
						Id:       v.ID,
						VideoId:  v.VideoId,
						Duration: v.Duration,
						Url:      v.Url,
						DownUrl:  v.DownUrl,
						Cover:    v.Cover,
						Title:    v.Title,
					})
				}
			}
		}

	}
	return result, err
}

func (rpc *rpcServer) EditChapter(ctx context.Context, request *pb.EditChapterRequest) (*pb.EditChapterResponse, error) {
	id, err := courseModel.EditChapter(request.Title, request.Desc, request.Id, request.Source, request.CreateBy, request.HomeWork)
	return &pb.EditChapterResponse{Id: id}, err
}

func (rpc *rpcServer) EditChapterCoverInfo(ctx context.Context, request *pb.EditChapterRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, courseModel.EditChapterCover(request.InfoContent, request.Prefix, request.Id, request.Type)
}

func (rpc *rpcServer) EditChapterSort(ctx context.Context, request *pb.ChapterRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, courseModel.EditChapterSort(request.Id, request.Sort)
}
