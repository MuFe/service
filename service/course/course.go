package course

import (
	"context"
	"mufe_service/camp/enum"
	"mufe_service/camp/service"
	"mufe_service/camp/xlog"
	pb "mufe_service/jsonRpc"
	"mufe_service/model/course"
	"mufe_service/model/tag"
)

func init() {
	nSer := &rpcServer{}
	pb.RegisterCourseServiceServer(service.GetRegisterRpc(), nSer)

}

type rpcServer struct {
}

func (rpc *rpcServer) GetCourse(ctx context.Context, request *pb.CourseServiceRequest) (*pb.CourseServiceResponse, error) {
	result := &pb.CourseServiceResponse{}
	list, err := courseModel.GetCourseFromID(request.Ids, request.LevelId, request.Status)
	if err == nil {
		for _, v := range list {
			tag := make([]*pb.TagData, 0)
			for _, va := range v.Tag {
				tag = append(tag, &pb.TagData{
					Id:    va.ID,
					Title: va.Name,
				})
			}
			result.List = append(result.List, &pb.CourseData{
				Id:      v.ID,
				Title:   v.Title,
				Desc:    v.Desc,
				Cover:   v.Cover,
				Section: v.Section,
				Bg:      v.Bg,
				Tag:     tag,
			})
		}
	}
	return result, err
}

func (rpc *rpcServer) GetAdminCourse(ctx context.Context, request *pb.CourseServiceRequest) (*pb.CourseServiceResponse, error) {
	result := &pb.CourseServiceResponse{}
	list, err := courseModel.GetAdminCourse(request.Ids, request.Status)
	if err == nil {
		for _, v := range list {
			result.List = append(result.List, &pb.CourseData{
				Id:    v.ID,
				Title: v.Title,
				Cover: v.Cover,
			})
		}
	}
	return result, err
}

func (rpc *rpcServer) GetAdminCourseDetail(ctx context.Context, request *pb.CourseServiceRequest) (*pb.CourseServiceResponse, error) {
	result := &pb.CourseServiceResponse{}
	v, err := courseModel.GetAdminCourseDetail(request.Ids[0])
	if err == nil {
		data, err := courseModel.GetOrigin(request.Ids[0])
		temp := &pb.GetOriginResponse{

		}
		if err == nil {
			temp.Photo = data.Photo
			temp.Info = data.Info
			temp.Title = data.Title
			temp.Name = data.Name
			temp.Desc = data.Desc
			temp.Auth = data.Auth
			temp.InfoTitle=data.InfoTitle
			temp.Certificate=data.Certificate
		}
		tag := make([]*pb.TagData, 0)
		for _, va := range v.Tag {
			tag = append(tag, &pb.TagData{
				Id:    va.ID,
				Title: va.Name,
			})
		}
		result.List = append(result.List, &pb.CourseData{
			Id:    v.ID,
			Title: v.Title,
			Cover: v.Cover,
			Bg:    v.Bg,
			Desc:  v.Desc,
			Tag:   tag,
			Data:  temp,
		})

	}
	return result, err
}

func (rpc *rpcServer) EditCourse(ctx context.Context, request *pb.EditCourseRequest) (*pb.EditCourseResponse, error) {
	id, err := courseModel.EditCourse(request.Title, request.OriginName, request.OriginTitle, request.OriginDesc, request.OriginInfo, request.OriginInfoTitle, request.Certificate,
		request.Id, request.Level, request.CreateBy)
	return &pb.EditCourseResponse{Id: id}, err

}

func (rpc *rpcServer) EditCourseCover(ctx context.Context, request *pb.EditCourseRequest) (*pb.EmptyResponse, error) {
	if request.Type == 1 {
		return &pb.EmptyResponse{}, courseModel.EditCourseCover(request.Id, request.Cover, request.Prefix)
	} else if request.Type == 2 {
		return &pb.EmptyResponse{}, courseModel.EditCourseBgCover(request.Id, request.Cover, request.Prefix)
	} else if request.Type == 3 {
		data, err := courseModel.GetOrigin(request.Id)
		if err != nil {
			return nil, err
		}
		return &pb.EmptyResponse{}, courseModel.EditOriginCourseCover(data.Id, request.Cover, request.Prefix)
	} else if request.Type == 4 {
		data, err := courseModel.GetOrigin(request.Id)
		if err != nil {
			return nil, err
		}
		list := make([]courseModel.Auth, 0)
		for _, v := range request.AuthList {
			list = append(list, courseModel.Auth{
				Cover:  v.Cover,
				Prefix: v.Prefix,
			})
		}
		xlog.Info(list)
		return &pb.EmptyResponse{}, courseModel.EditOriginAuth(data.Id, list)
	} else {
		return nil, nil
	}

}

func (rpc *rpcServer) GetCourseLevel(ctx context.Context, request *pb.CourseLevelRequest) (*pb.CourseLevelResponse, error) {
	result := &pb.CourseLevelResponse{}
	list, err := courseModel.GetLevel(request.Type)
	if err == nil {
		for _, v := range list {
			result.List = append(result.List, &pb.CourseLevel{
				Id:   v.ID,
				Name: v.Name,
			})
		}
	}
	return result, err
}

func (rpc *rpcServer) TagList(ctx context.Context, request *pb.TagRequest) (*pb.TagResponse, error) {
	result := &pb.TagResponse{}
	list, err := tagMoel.GetTagList(request.Status, request.Id)
	if err == nil {
		for _, v := range list {
			result.List = append(result.List, &pb.TagData{
				Id:      v.ID,
				Title:   v.Name,
				Content: v.Content,
				Cover:   v.Cover,
			})
		}
	}
	return result, err
}

func (rpc *rpcServer) EditTag(ctx context.Context, request *pb.EditTagRequest) (*pb.TagData, error) {
	if request.Type == enum.AddTag {
		if request.Id == 0 {
			id, err := tagMoel.AddTag(request.Title, request.Content)
			return &pb.TagData{Id: id}, err
		} else {
			return &pb.TagData{Id: request.Id}, tagMoel.EditTagDetail(request.Title, request.Content, request.Id)
		}
	} else if request.Type == enum.EditTagCover {
		return &pb.TagData{}, tagMoel.EditTagCover(request.Cover, request.Prefix, request.Id)
	} else {
		return &pb.TagData{}, nil
	}
}

func (rpc *rpcServer) Notice(ctx context.Context, request *pb.NoticeRequest) (*pb.NoticeResponse, error) {
	result := courseModel.GetNotice(request.Id)
	list := make([]*pb.NoticeData, 0)
	for _, v := range result {
		list = append(list, &pb.NoticeData{
			Id:    v.Id,
			Time:  v.Time,
			Title: v.Title,
		})
	}
	return &pb.NoticeResponse{List: list}, nil
}

func (rpc *rpcServer) NoticeDetail(ctx context.Context, request *pb.NoticeRequest) (*pb.NoticeData, error) {
	result, err := courseModel.NoticeDeatil(request.Id)
	if err == nil {
		return &pb.NoticeData{
			Id:       result.Id,
			Title:    result.Title,
			Content:  result.Content,
			Time:     result.Time,
			CreateBy: result.CreateBy,
		}, nil
	} else {
		return nil, err
	}
}

func (rpc *rpcServer) AddNotice(ctx context.Context, request *pb.NoticeRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, courseModel.AddNotice(request.Id, request.CreateBy, request.Title, request.Content)
}

func (rpc *rpcServer) EditCourseLevel(ctx context.Context, request *pb.EditCourseLevelRequest) (*pb.EditCourseLevelResponse, error) {
	id, err := courseModel.EditLevel(request.Name, request.Id, request.Type)
	return &pb.EditCourseLevelResponse{Id: id}, err
}

func (rpc *rpcServer) GetOrigin(ctx context.Context, request *pb.GetOriginRequest) (*pb.GetOriginResponse, error) {
	data, err := courseModel.GetOrigin(request.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetOriginResponse{
		Photo: data.Photo,
		Info:  data.Info,
		Title: data.Title,
		Name:  data.Name,
		Desc:  data.Desc,
		Auth:  data.Auth,
		InfoTitle:data.InfoTitle,
		Certificate:data.Certificate,
	}, nil
}
