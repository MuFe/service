package coach

import (
	"context"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/service"
	"mufe_service/camp/xlog"
	pb "mufe_service/jsonRpc"
	coachModel "mufe_service/model/coach"
)

type rpcServer struct {
}

func (rpc *rpcServer) GetInstitution(ctx context.Context, request *pb.InstitutionRequest) (result *pb.InstitutionResponse, err error) {
	result = &pb.InstitutionResponse{}
	id := int64(0)
	if request.Uid != 0 {
		list, err := coachModel.GetInstitutionCoach(0, request.Uid)
		if err == nil && len(list) > 0 {
			id = list[0].InstitutionId
		} else {
			return nil, xlog.Error(errcode.HttpErrorWringParam.Msg)
		}
	} else {
		id = request.Id
	}
	list, err := coachModel.GetInstitution(request.Status, id, request.Code)
	for _, v := range list {
		result.List = append(result.List, &pb.InstitutionData{
			Icon:    v.Icon,
			Name:    v.Name,
			Address: v.Address,
			Id:      v.Id,
			Code:    v.Code,
		})
	}
	return result, err
}

func (rpc *rpcServer) CancelCoach(ctx context.Context, request *pb.CancelCoachRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, coachModel.DelInstitutionCoachFromUid(request.List)
}

func (rpc *rpcServer) EditInstitution(ctx context.Context, request *pb.EditInstitutionRequest) (*pb.EditInstitutionResponse, error) {
	var err error
	id := int64(0)
	id, err = coachModel.EditInstitution(request.Create, request.Id, request.Name, request.Address, request.Icon, request.Prefix, request.Del)
	return &pb.EditInstitutionResponse{Id: id}, err
}

func (rpc *rpcServer) EditInstitutionSchool(ctx context.Context, request *pb.EditInstitutionSchoolRequest) (*pb.EditInstitutionResponse, error) {
	var err error
	id := int64(0)
	id, err = coachModel.EditInstitutionSchool(request.Create, request.Id, request.Start, request.End, request.ParentId, request.Name, request.Address, request.Phone, request.Icon, request.Prefix, request.Del, request.Course)
	return &pb.EditInstitutionResponse{Id: id}, err
}

func (rpc *rpcServer) EditInstitutionCourse(ctx context.Context, request *pb.EditInstitutionCourseRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, coachModel.EditInstitutionCourse(request.Create, request.Id, request.Max, request.Price, request.ParentId, request.Duration, request.Name, request.Level, request.Del)
}

func (rpc *rpcServer) JoinInstitution(ctx context.Context, request *pb.JoinInstitutionCourseRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, coachModel.EditInstitutionCoach(request.Id,request.Uid, request.InstitutionId, request.Info, request.Quit)
}

func (rpc *rpcServer) WorkOrder(ctx context.Context, request *pb.InstitutionWorkOrderRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, coachModel.EditInstitutionCoachOrder(request.Id, request.Uid, request.Status,request.OrderId)
}

func (rpc *rpcServer) AddWork(ctx context.Context, request *pb.InstitutionWorkRequest) (*pb.EmptyResponse, error) {
	courseResult, err := coachModel.GetInstitutionCourse(enum.StatusNormal, request.Course, 0)
	if err != nil {
		return nil, xlog.Error(err)
	} else if len(courseResult) == 0 {
		return nil, xlog.Error(errcode.HttpErrorWringParam.Msg)
	}
	return &pb.EmptyResponse{}, coachModel.EditInstitutionCoachWork(request.Id, request.Start, request.End, request.Place, request.Course, request.Uid,
		courseResult[0].Price, courseResult[0].Max, courseResult[0].Duration, courseResult[0].Level, request.Desc, courseResult[0].Name, request.Del)
}

func (rpc *rpcServer) SchoolList(ctx context.Context, request *pb.InstitutionSchoolRequest) (result *pb.InstitutionSchoolResponse, err error) {
	result = &pb.InstitutionSchoolResponse{}
	list, err := coachModel.GetInstitutionSchool(request.Status, request.Iid, []int64{request.Id})
	resultMap := make(map[int64]*pb.InstitutionSchool)
	schoolIdList := make([]int64, 0)
	for _, v := range list {
		temp := &pb.InstitutionSchool{
			Icon:    v.Icon,
			Name:    v.Name,
			Address: v.Address,
			Id:      v.Id,
			Start:   v.Start,
			End:     v.End,
			Phone:   v.Phone,
			Course:  make([]*pb.InstitutionCourse, 0),
		}
		resultMap[v.Id] = temp
		result.List = append(result.List, temp)
		schoolIdList = append(schoolIdList, v.Id)
	}
	courseList, err := coachModel.GetInstitutionSchoolCourse(schoolIdList)
	if err == nil {
		for _, v := range courseList {
			temp, ok := resultMap[v.School]
			if ok {
				temp.Course = append(temp.Course, &pb.InstitutionCourse{
					Name:     v.Name,
					Id:       v.Id,
					SchoolId: v.School,
				})
			}
		}
	}
	return result, err
}

func (rpc *rpcServer) CourseList(ctx context.Context, request *pb.InstitutionCourseRequest) (result *pb.InstitutionCourseResponse, err error) {
	result = &pb.InstitutionCourseResponse{}
	if len(request.SchoolId) > 0 {
		list, err := coachModel.GetInstitutionSchoolCourse(request.SchoolId)
		for _, v := range list {
			result.List = append(result.List, &pb.InstitutionCourse{
				Price: v.Price,
				Name:  v.Name,
				Level: v.Level,
				Id:    v.Id,
				Max:   v.Max,
				Duration:v.Duration,
			})
		}
		return result, err
	} else {
		list, err := coachModel.GetInstitutionCourse(request.Status, request.Id, request.Iid)
		for _, v := range list {
			result.List = append(result.List, &pb.InstitutionCourse{
				Price: v.Price,
				Name:  v.Name,
				Level: v.Level,
				Id:    v.Id,
				Max:   v.Max,
				Duration:v.Duration,
			})
		}
		return result, err
	}
}

func (rpc *rpcServer) CoachList(ctx context.Context, request *pb.CoachListRequest) (result *pb.CoachListResponse, err error) {
	result = &pb.CoachListResponse{}
	list, err := coachModel.GetInstitutionCoach(request.Iid, request.Id)
	for _, v := range list {
		result.List = append(result.List, &pb.CoachData{
			Uid:           v.Uid,
			InstitutionId: v.InstitutionId,
			Info:          v.Info,
			Id:v.Id,
		})
	}
	return result, err
}

func (rpc *rpcServer) WorkOrderList(ctx context.Context, request *pb.InstitutionWorkOrderRequest) (result *pb.InstitutionWorkOrderResponse, err error) {
	result = &pb.InstitutionWorkOrderResponse{}
	list, err := coachModel.GetInstitutionCoachOrder(request.Uid, request.Page, request.Size, request.OrderId,request.QueryStatus)
	for _, v := range list {
		result.List = append(result.List, &pb.InstitutionWorkOrderData{
			Uid:     v.Uid,
			WorkId:  v.WorkId,
			OrderId: v.OrderId,
			Num:     v.Num,
			Status:  v.Status,
			Id:      v.Id,
		})
	}
	return result, err
}

func (rpc *rpcServer) WorkList(ctx context.Context, request *pb.InstitutionWorkListRequest) (result *pb.InstitutionWorkListResponse, err error) {
	result = &pb.InstitutionWorkListResponse{}
	list, err := coachModel.GetInstitutionCoachWork(request.Place, request.Uid, request.Start, request.End, request.Id)
	placeIdMap := make(map[int64]int64)
	listMap := make(map[int64]*pb.InstitutionWorkData)
	placeIdList := make([]int64, 0)
	idList := make([]int64, 0)
	for _, v := range list {
		placeIdMap[v.Id] =  v.PlaceId
		placeIdList = append(placeIdList, v.PlaceId)
		idList = append(idList, v.Id)
		temp := &pb.InstitutionWorkData{
			Uid:      v.Uid,
			Id:       v.Id,
			Start:    v.Start,
			End:      v.End,
			Price:    v.Price,
			Max:      v.Max,
			Now:      v.Now,
			Level:    v.Level,
			Name:     v.Name,
			UserInfo: v.Info,
			Desc:     v.Desc,
			Duration: v.Duration,
			Reserve:  v.Reserve,
		}
		result.List = append(result.List, temp)
		listMap[v.Id] = temp
	}
	placeResult, err := coachModel.GetInstitutionSchool(enum.StatusNormal, 0, placeIdList)
	if err == nil {
		placeResultMap:=make(map[int64]coachModel.InstitutionSchool)
		for _, v := range placeResult {
			placeResultMap[v.Id]=v
		}
		for _,vv:=range result.List{
			idTemp, ok := placeIdMap[vv.Id]
			if ok {
				v, ok := placeResultMap[idTemp]
				if ok {
					vv.Place = &pb.InstitutionSchool{
						Icon:    v.Icon,
						Name:    v.Name,
						Address: v.Address,
						Id:      v.Id,
						Start:   v.Start,
						End:     v.End,
						Phone:   v.Phone,
					}
				}
			}
		}
	}
	re, err := coachModel.GetInstitutionCoachOrderNum(idList)
	if err == nil {
		for _, v := range re {
			dataTemp, ok := listMap[v.Id]
			if ok {
				dataTemp.Now = v.Num
			}
		}
	}
	return result, err
}

func init() {
	nSer := &rpcServer{}
	pb.RegisterCoachServiceServer(service.GetRegisterRpc(), nSer)
}
