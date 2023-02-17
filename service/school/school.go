package school

import (
	"context"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/service"
	"mufe_service/camp/xlog"
	pb "mufe_service/jsonRpc"
	"mufe_service/model/course"
	schoolModel "mufe_service/model/school"
)

func init() {
	nSer := &rpcServer{}
	pb.RegisterSchoolServiceServer(service.GetRegisterRpc(), nSer)

}

type rpcServer struct {
}

func (rpc *rpcServer) AddSchool(ctx context.Context, request *pb.AddSchoolRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, schoolModel.AddSchool(request.Uid, request.SchoolId)
}

func (rpc *rpcServer) QuitSchool(ctx context.Context, request *pb.AddSchoolRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, schoolModel.QuitSchool(request.Uid, request.SchoolId)
}

func (rpc *rpcServer) MySchool(ctx context.Context, request *pb.SchoolRequest) (*pb.SchoolResponse, error) {
	list := schoolModel.GetSchool(request.Uid, request.Id, enum.StatusNormal)
	result := &pb.SchoolResponse{
		List: make([]*pb.SchoolData, 0),
	}
	for _, v := range list {
		result.List = append(result.List, &pb.SchoolData{
			Id:      v.Id,
			Icon:    v.Icon,
			Desc:    v.Desc,
			Name:    v.Name,
			Address: v.Address,
			Code:    v.Code,
		})
	}
	return result, nil
}
func (rpc *rpcServer) EditSchool(ctx context.Context, request *pb.SchoolData) (*pb.AddSchoolResponse, error) {
	id, err := schoolModel.EditSchool(request.Name, request.Address, request.Icon, request.Id, request.TypeId)
	return &pb.AddSchoolResponse{Id: id}, err
}

func (rpc *rpcServer) SchoolList(ctx context.Context, request *pb.SchoolRequest) (*pb.SchoolResponse, error) {
	list := schoolModel.GetSchool(request.Uid, request.Id, enum.StatusNormal)
	result := &pb.SchoolResponse{
		List: make([]*pb.SchoolData, 0),
	}
	for _, v := range list {
		result.List = append(result.List, &pb.SchoolData{
			Id:       v.Id,
			Icon:     v.Icon,
			Desc:     v.Desc,
			Name:     v.Name,
			Code:     v.Code,
			Address:  v.Address,
			TypeId:   v.TypeId,
			TypeName: v.Type,
		})
	}
	return result, nil
}

func (rpc *rpcServer) SchoolTypeList(ctx context.Context, request *pb.EmptyRequest) (*pb.SchoolTypeResponse, error) {
	list := schoolModel.GetSchoolType()
	result := &pb.SchoolTypeResponse{
		List: make([]*pb.SchoolTypeData, 0),
	}
	for _, v := range list {
		data := &pb.SchoolTypeData{
			Id:    v.Id,
			Name:  v.Name,
			Grade: make([]*pb.GradeInfo, 0),
		}
		for _, vv := range v.List {
			data.Grade = append(data.Grade, &pb.GradeInfo{
				Id:   vv.Id,
				Name: vv.Name,
			})
		}
		result.List = append(result.List, data)
	}
	return result, nil
}

func (rpc *rpcServer) Scan(ctx context.Context, request *pb.ScanRequest) (*pb.ScanResponse, error) {
	result := &pb.ScanResponse{}
	list := schoolModel.FindSchool(request.Content, enum.StatusNormal)
	if len(list) == 0 {
		classResult := schoolModel.FindClass(request.Content, enum.StatusNormal)
		if len(classResult) == 0 {
			return nil, xlog.Error("无法识别二维码")
		}
		v := classResult[0]
		result.Class = &pb.Class{
			Id:         v.Id,
			Name:       v.Name,
			SchoolType: v.SchoolType,
			Grade:      v.Grade,
			GradeId:    v.GradeId,
			CreateTime: v.CreateTime,
			Number:     v.Number,
		}
	} else {
		v := list[0]
		result.School = &pb.SchoolData{
			Id:      v.Id,
			Icon:    v.Icon,
			Desc:    v.Desc,
			Name:    v.Name,
			Code:    v.Code,
			Address: v.Address,
		}
	}
	return result, nil
}

func (rpc *rpcServer) GradeList(ctx context.Context, request *pb.GradeRequest) (*pb.GradeTypeResponse, error) {
	list, err := schoolModel.GradeInfo(request.SchoolId)
	if err != nil {
		return nil, xlog.Error(err)
	}
	result := &pb.GradeTypeResponse{
		List: make([]*pb.GradeTypeInfo, 0),
	}
	for _, v := range list {
		tempList := make([]*pb.GradeInfo, 0)
		for _, vv := range v.List {
			tempList = append(tempList, &pb.GradeInfo{
				Id:   vv.Id,
				Name: vv.Name,
			})
		}
		result.List = append(result.List, &pb.GradeTypeInfo{
			Type:   v.Name,
			TypeId: v.Id,
			List:   tempList,
		})
	}
	return result, nil
}

func (rpc *rpcServer) CreateClassInfo(ctx context.Context, request *pb.CreateClassInfoRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, schoolModel.CreateClassInfo(request.Uid, request.SchoolId, request.GradeId, request.Name)
}

func (rpc *rpcServer) CreateClass(ctx context.Context, request *pb.CreateClassRequest) (*pb.EmptyResponse, error) {
	list := schoolModel.GetSchool(request.Uid, 0, enum.StatusNormal)
	if len(list) == 0 {
		return nil, xlog.Error("您还没有加入任何的学校")
	}
	idList := make([]int64, 0)
	for _, v := range list {
		idList = append(idList, v.Id)
	}
	return &pb.EmptyResponse{}, schoolModel.CreateClass(request.Uid, request.ClassInfoId, idList)
}

func (rpc *rpcServer) JoinClass(ctx context.Context, request *pb.JoinClassRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, schoolModel.JoinClass(request.Uid, request.ClassId, request.Type, nil)
}

func (rpc *rpcServer) QuitClass(ctx context.Context, request *pb.QuitClassRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, schoolModel.QuitClass(request.Uid, request.ClassId, request.QuitUid)
}

func (rpc *rpcServer) DissolutionClass(ctx context.Context, request *pb.QuitClassRequest) (*pb.QuitClassResponse, error) {
	list,err:=schoolModel.DisClass(request.Uid, request.ClassId)
	return &pb.QuitClassResponse{
		StudentList:list,
	}, err
}

func (rpc *rpcServer) ClassList(ctx context.Context, request *pb.ClassRequest) (*pb.ClassResponse, error) {
	var err error
	var result []*schoolModel.Class
	if request.School != 0 {
		result, err = schoolModel.AdminClassList(request.School, request.Grade)
	} else {
		result, err = schoolModel.ClassList(request.Uid, request.Id)
	}

	if err != nil {
		return nil, xlog.Error(err)
	}
	list := make([]*pb.Class, 0)
	for _, v := range result {
		list = append(list, &pb.Class{
			Id:         v.Id,
			Name:       v.Name,
			SchoolType: v.SchoolType,
			Grade:      v.Grade,
			GradeId:    v.GradeId,
			CreateTime: v.CreateTime,
			Number:     v.Number,
			Uid:        request.Uid,
			Code:       v.Code,
			Tag:        v.Tag,
			SchoolIcon: v.SchoolIcon,
			SchoolName: v.SchoolName,
			CreateBy:v.CreateBy,
		})
	}
	return &pb.ClassResponse{List: list}, nil
}

func (rpc *rpcServer) ClassDetail(ctx context.Context, request *pb.ClassDetailRequest) (*pb.Class, error) {
	v, err := schoolModel.ClassDetail(request.Uid, request.Id)
	if err != nil {
		return nil, xlog.Error(err)
	}
	info := courseModel.GetVideoIndexInfo(v.CourseId, v.VideoId)
	return &pb.Class{
		Id:           v.Id,
		Name:         v.Name,
		SchoolType:   v.SchoolType,
		Grade:        v.Grade,
		GradeId:      v.GradeId,
		CreateTime:   v.CreateTime,
		Number:       v.Number,
		Uid:          request.Uid,
		Code:         v.Code,
		Tag:          v.Tag,
		ChapterIndex: info.ChapterIndex,
		VideoIndex:   info.VideoIndex,
		VideoId:      v.VideoId,
		Progress:     info.Progress,
		CourseId:     v.CourseId,
		AdminList:    v.AdminList,
		StudentList:  v.List,
	}, nil
}

func (rpc *rpcServer) ClassInfo(ctx context.Context, request *pb.ClassInfoRequest) (*pb.ClassResponse, error) {
	result, err := schoolModel.ClassInfo(request.SchoolId, request.GradeId)
	if err != nil {
		return nil, xlog.Error(err)
	}
	list := make([]*pb.Class, 0)
	for _, v := range result {
		list = append(list, &pb.Class{
			Id:         v.Id,
			Name:       v.Name,
			SchoolType: v.SchoolType,
			Grade:      v.Grade,
			GradeId:    v.GradeId,
			CreateTime: v.CreateTime,
			Number:     v.Number,
		})
	}
	return &pb.ClassResponse{List: list}, nil
}

func (rpc *rpcServer) AddCourse(ctx context.Context, request *pb.AddCourseRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, schoolModel.AddClassCourse(request.ClassId, request.CourseId, request.Uid)
}

func (rpc *rpcServer) RemoveCourse(ctx context.Context, request *pb.AddCourseRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, schoolModel.RemoveClassCourse(request.ClassId, request.Uid)
}

func (rpc *rpcServer) EditCourseProgress(ctx context.Context, request *pb.EditCourseProgressRequest) (*pb.EmptyResponse, error) {
	v, err := schoolModel.ClassDetail(request.Uid, request.ClassId)
	if err != nil {
		return nil, xlog.Error(err)
	}
	info := courseModel.GetVideoIndexInfo(v.CourseId, v.VideoId)
	isHave := false
	for _, v := range info.List {
		if request.VideoId == v.Id {
			isHave = true
		}
	}
	if !isHave {
		return nil, xlog.Error(errcode.HttpErrorWringParam.Msg)
	}
	return &pb.EmptyResponse{}, schoolModel.EditCourseProgress(request.Uid, request.ClassId, request.VideoId)
}

func (rpc *rpcServer) UserList(ctx context.Context, request *pb.TeacherUserRequest) (*pb.TeacherUserResponse, error) {
	return &pb.TeacherUserResponse{List: schoolModel.GetUserList(request.SchoolId, request.GradeId, request.ClassId, request.Type)}, nil
}

func (rpc *rpcServer) CancelSchool(ctx context.Context, request *pb.CancelSchoolRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, schoolModel.CancelClassRecord(request.List)
}
