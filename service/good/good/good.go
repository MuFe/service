package delivery

import (
	"context"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/sequence"
	"mufe_service/camp/service"
	"mufe_service/camp/xlog"
	pb "mufe_service/jsonRpc"
	"mufe_service/manager"
	goodmodel "mufe_service/model/good"
)

func init() {
	pb.RegisterGoodServiceServer(service.GetRegisterRpc(), &rpcServer{})
}

type rpcServer struct{}

func (rpc *rpcServer) EditGood(ctx context.Context, request *pb.EditGoodRequest) (result *pb.EditGoodResponse, err error) {
	return editGood(request)
}


func editGood(request *pb.EditGoodRequest) (*pb.EditGoodResponse, error) {
	spuId := request.SpuId
	if request.Type == enum.GoodAdd {
		result, err := addGood(
			request.SpuName,
			request.CategoryId,
			request.Delivery,
			request.Location,
			request.Infos,
			request.BusinessId,
			request.IsDraft,
			request.Uid,
			request.Status)
		return result, err
	}
	return &pb.EditGoodResponse{SpuId: spuId}, nil
}


func addGood(name string, categoryId int64, deliverInfo []*pb.DeliveryData, spuDeliveryLocation string, infos []*pb.Sku, businessId int64, draft bool, uid int64, status int64) (*pb.EditGoodResponse, error) {
	spuNo := ""
	categoryName := ""
	if categoryId != 0 {
		categoryName,err:=goodmodel.GetGoodCategoryName(categoryId)
		if err != nil {
			return nil, xlog.Error(err)
		}
		spuNo, err = sequence.SpuNo.NewNo(categoryName)
		if err != nil {
			spuNo = ""
		}
	}
	var spuId int64
	err := db.GetGoodDb().WithTransaction(func(tx *db.Tx) error {
		spuId,err:=goodmodel.AddGood(tx,name,spuNo,status,spuDeliveryLocation)
		if err != nil {
			return xlog.Error(err)
		}
		if categoryId != 0 {
			err=goodmodel.AddGoodCategoryRecord(tx,spuId,categoryId)
		}
		if err != nil {
			return xlog.Error(err)
		}
		err=goodmodel.AddBrandRecord(businessId,spuId,tx)
		if err != nil {
			return xlog.Error(err)
		}
		err = goodmodel.EditSku(infos, int64(spuId), name, categoryName, int64(status), tx, businessId)
		if err != nil {
			return xlog.Error(err)
		}
		return goodmodel.AddDelivery(tx,deliverInfo,spuId,uid)
	})
	if err != nil {
		return nil, err
	}
	return &pb.EditGoodResponse{SpuId: int64(spuId)}, err

}


func (rpc *rpcServer) GoodDetail(context context.Context, request *pb.GoodDetailRequest) (*pb.GoodDetailResponse, error) {
	if request.SpuId == 0 {
		request.SpuId = goodmodel.GetSpuIdFromSkuId(request.SkuId)
	}
	if request.SpuId == 0 {
		return nil, xlog.Error("参数有误")
	}
	result, err := goodmodel.GetDetail(request.SpuId, request.SkuStatus)
	if err != nil {
		return nil, err
	}
	tempMap, _ := goodmodel.GetSkuOptions([]int64{request.SpuId}, []int64{})
	photos := make([]*pb.PhotoInfo, 0)
	for _, info := range result.Photos {
		photos = append(photos, &pb.PhotoInfo{
			Key: info.Key,
			Url: info.Prefix + info.Key,
		})
	}
	skus := make([]*pb.Sku, 0)
	deliveryInfoList := make([]*pb.DeliveryData, 0)
	for _, info := range result.DeliveryInfo {
		deliveryInfoList = append(deliveryInfoList, &pb.DeliveryData{
			Id: info.TemplateId,
			Type:       info.Type,
		})
	}
	skuIdList := make([]int64, 0)
	for _, info := range result.List {
		skuIdList = append(skuIdList, info.SkuID)
	}
	// 库存
	stockMap, err := goodmodel.GetGoodStockBySkuIDs(skuIdList, result.BusinessId, request.BusinessGroupId)
	if err != nil {
		return nil, xlog.Error(err)
	}
	for _, info := range result.List {
		options := make([]*pb.SkuOption, 0)
		temp, ok := tempMap[info.SkuID]
		if ok {
			for _, optionsInfo := range temp {
				options = append(options, &pb.SkuOption{
					OptionValue:   optionsInfo.Value,
					OptionId:      optionsInfo.Id,
					OptionValueId: optionsInfo.ValueId,
				})
			}
		}
		skus = append(skus, &pb.Sku{
			SkuId:          info.SkuID,
			SkuName:        info.SkuName,
			Price:          info.Price,
			Stock:          stockMap[info.SkuID],
			Status:         info.Status,
			Options:        options,
		})
	}
	return &pb.GoodDetailResponse{
		SpuId:            result.Id,
		SpuName:          result.Name,
		Status:           result.Status,
		StatusMessage:    result.StatusMessage,
		SpuNumber:        result.Number,
		Location:         result.Location,
		Detail:           result.Detail,
		CategoryId:       result.CategoryId,
		ParentCategoryId: result.ParentCategoryId,
		CommentNum:       result.CommentNum,
		Photos:           photos,
		List:             skus,
		DeliveryInfo: deliveryInfoList,
		BusinessId:       result.BusinessId,
	}, nil
}

func (rpc *rpcServer) GetOrderSkuList(ctx context.Context, request *pb.GetOrderSkuListRequest) (*pb.GetOrderSkuListResponse, error) {
	result, err := goodmodel.GetSkuOrderInfo(request.Data)
	if err != nil {
		return nil, err
	}
	xlog.Info(result)
	detailMap := make(map[int64]*pb.GetSubOrderData)
	for _, info := range result.List {
		temp, ok := detailMap[info.BusinessId]
		if !ok {
			temp = &pb.GetSubOrderData{
				SellerId:              info.BusinessId,
				List:                  make([]*pb.GetSubOrderDetailData, 0),
				ShopDeliverPriceMap:   make(map[int64]int64),
			}
			detailMap[info.BusinessId] = temp
		}

		if _, ok = temp.ShopDeliverPriceMap[info.SpuID]; !ok {
			temp.ShopDeliverPriceMap[info.SpuID] = result.ShopDeliverPrice[info.SpuID]
		}
		temp.List = append(temp.List, &pb.GetSubOrderDetailData{
			SkuId:          info.SkuID,
			Num:            info.Num,
			Price:          info.Price,
			SpuId:          info.SpuID,
		})
	}
	list := make([]*pb.GetSubOrderData, 0)
	for _, info := range detailMap {
		list = append(list, info)
	}
	xlog.Info(list)
	return &pb.GetOrderSkuListResponse{
		List: list,
	},nil
}

func (rpc *rpcServer) GetSpuList(context context.Context, request *pb.GetSpuListRequest) (*pb.GetSpuListResponse, error) {
	returnLis := make([]*pb.Spu, 0)
	requestList := make([]int64, 0)

	requestList = request.List
	spuList, err := goodmodel.GetSpuListBySpuIDS(requestList)
	if err != nil {
		return nil, err
	}
	if len(spuList) > 0 {
		spuIDs := make([]int64, 0)
		for _, info := range spuList {
			spuIDs = append(spuIDs, info.SpuID)
		}
		// 第一张图片
		var photoMap = make(map[int64]string, 0)
		photos, err := goodmodel.GetSpuPhotoBySpuIDs(spuIDs)
		if err != nil {
			return nil, xlog.Error(err)
		}
		for _, v := range photos {
			if _, ok := photoMap[v.SpuID]; !ok {
				photoMap[v.SpuID] = v.Prefix + v.Key
			}
		}
		// 最低价的sku
		var skuMap = make(map[int64]goodmodel.Sku, 0)
		var memberSkuMap = make(map[int64]goodmodel.Sku, 0)
		spuMap := make(map[int64]int64) //sku-spu
		skus, err := goodmodel.GetSkuIDBySpuIDs(spuIDs, enum.StatusNormal)
		if err != nil {
			return nil, xlog.Error(err)
		}
		for _, info := range skus {
			spuMap[info.SkuID] = info.SpuID
			if v, ok := skuMap[info.SpuID]; !ok || v.Price > info.Price {
				skuMap[info.SpuID] = info
			}
		}
		for _, info := range skus {
			if _, ok := memberSkuMap[info.SpuID]; !ok  {
				memberSkuMap[info.SpuID] = info
			}
		}
		// 库存
		stockMap, err := goodmodel.GetAllGoodStockOnSpu(spuIDs)
		if err != nil {
			return nil, xlog.Error(err)
		}
		// 销量
		saleResult, err := manager.GetOrderService().GetOrderNum(context, &pb.GetOrderNumRequest{
			SpuSkuMap: spuMap,
		})
		if err != nil {
			return nil, xlog.Error(err)
		}
		saleMap := saleResult.Result
		for _, g := range spuList {
			s := &pb.Spu{
				Id:           g.SpuID,
				Name:         g.SpuName,
				Photo:        photoMap[g.SpuID],
				Sales:        saleMap[g.SpuID],
				Stock:        stockMap[g.SpuID],
				Status:       g.Status,
				SaleStatus:   g.SaleStatus,
				CreateTime:   g.CreateTime,
				ModifyTime:   g.ModifyTime,
				BusinessId:   g.BusinessID,
				BrandName:    g.BusinessName,
				CategoryId:   g.CategoryID,
				CategoryName: g.CategoryName,
				Normal:       &pb.SpuSku{},
				Member:       &pb.SpuSku{},
			}
			if k, ok := skuMap[g.SpuID]; ok {
				if s.Normal.SkuId < 1 {
					s.Normal.SkuId = k.SkuID
					s.Normal.Price = k.Price
				}
			}
			if k, ok := memberSkuMap[g.SpuID]; ok {
				if s.Member.SkuId < 1 {
					s.Member.SkuId = k.SkuID
					s.Member.Price = k.Price
				}
			}
			returnLis = append(returnLis, s)
		}
	}

	return &pb.GetSpuListResponse{List: returnLis}, nil
}

func (rpc *rpcServer) GetSkuOptions(context context.Context, request *pb.GetSkuOptionsRequest) (*pb.GetSkuOptionsResponse, error) {
	tempMap, err := goodmodel.GetSkuOptions(request.SpuIdList, request.SkuIdList)
	if err != nil {
		return nil, err
	}
	resultMap := make(map[int64]*pb.GetSkuOptions)
	for k, v := range tempMap {
		list := make([]*pb.SkuOption, 0)
		for _, info := range v {
			list = append(list, &pb.SkuOption{
				OptionValue:   info.Value,
				OptionId:      info.Id,
				OptionValueId: info.ValueId,
				SpuId:         info.SpuId,
			})
		}
		resultMap[k] = &pb.GetSkuOptions{
			List: list,
		}
	}
	return &pb.GetSkuOptionsResponse{Result: resultMap}, nil
}


func (rpc *rpcServer) GetSkuIdQuery(context context.Context, request *pb.GetSkuIdQueryRequest) (*pb.GetSkuIdQueryResponse, error) {
	list, err := goodmodel.GetSkuIdFromQuery(request.Query,request.BrandList)
	if err != nil {
		return nil, err
	}
	return &pb.GetSkuIdQueryResponse{List: list}, nil
}
