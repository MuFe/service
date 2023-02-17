package category

import (
	"context"
	"sort"
	"mufe_service/camp/dataUtil"
	"mufe_service/camp/enum"
	"mufe_service/camp/service"
	pb "mufe_service/jsonRpc"
	goodmodel "mufe_service/model/good"
)

func init() {
	pb.RegisterGoodCategoryServiceServer(service.GetRegisterRpc(), &rpcServer{})
}

type rpcServer struct{}

func (rpc *rpcServer) EditGoodCategory(ctx context.Context, request *pb.EditGoodCategoryRequest) (result *pb.EmptyResponse, err error) {
	return &pb.EmptyResponse{}, goodmodel.EditGoodCategory(request.PId, request.Name, request.Id, request.Level, request.List)
}
func (rpc *rpcServer) EditGoodCategoryStatus(ctx context.Context, request *pb.EditGoodCategoryStatusRequest) (result *pb.EmptyResponse, err error) {
	return &pb.EmptyResponse{}, goodmodel.EditGoodCategoryStatus(request.Id, request.Status)
}

func (rpc *rpcServer) GoodCategory(context context.Context, request *pb.GetGoodCategoryRequest) (*pb.GetGoodCategoryResponse, error) {
	list := make([]*pb.GoodCategory, 0)
	if request.Type == enum.GetCategoryInfoFromId {
		result, err := goodmodel.GetCategoryFromId(request.CategoryIdList, request.SecondCategoryIdList, request.ParentCategoryIdList)
		if err != nil {
			return nil, err
		}
		keys := make([]int, 0)
		for k := range result {
			keys = append(keys, int(k))
		}
		sort.Ints(keys)
		for _, key := range keys {
			value := result[int64(key)]
			children := make([]*pb.GoodCategory, 0)
			for _, child := range value.Children {
				temp := &pb.GoodCategory{
					Id:   child.ID,
					Name: child.Name,
				}
				children = append(children, temp)
				for _, v := range child.Spec {
					temp.Specs = append(temp.Specs, &pb.GoodSpecificationsInfo{
						Id:   v.Id,
						Name: v.Name,
					})
				}
			}
			list = append(list, &pb.GoodCategory{
				Id:       value.ID,
				Name:     value.Name,
				Children: children,
			})
		}
	} else if request.Type == enum.GetCategoryFromSpu {
		result, err := goodmodel.GetCategoryInfoFromSpu(request.Channel, request.BusinessId, request.BusinessGroupId)
		if err != nil {
			return nil, err
		}
		for _, info := range result {
			list = append(list, &pb.GoodCategory{
				Id:   info.ID,
				Name: info.Name,
			})
		}
	} else if request.Type == enum.GetCategoryFromQuery {
		result, err := goodmodel.GetCategoryInfoList(request.Level, request.CategoryIdList)
		if err != nil {
			return nil, err
		}
		for _, info := range result {
			list = append(list, &pb.GoodCategory{
				Id:   info.ID,
				Name: info.Name,
			})
		}
	} else if request.Type == enum.GetAllCategory {
		result, err := goodmodel.GetAllCategory()
		if err != nil {
			return nil, err
		}
		list = dataUtil.MakeCategory(result)
	}
	return &pb.GetGoodCategoryResponse{List: list}, nil
}
