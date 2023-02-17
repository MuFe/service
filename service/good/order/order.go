package order

import (
	"context"
	"fmt"
	"os"
	"mufe_service/camp/enum"
	"mufe_service/camp/service"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	pb "mufe_service/jsonRpc"
	"mufe_service/manager"
	"mufe_service/model/order"
	"strings"
	"time"
)

func init() {
	pb.RegisterOrderServiceServer(service.GetRegisterRpc(), &rpcServer{})
}

type rpcServer struct{}

func (rpc *rpcServer) GetUserAddress(ctx context.Context, request *pb.AddressListServiceRequest) (*pb.AddressListServiceResponse,  error) {
	result,err:=orderModel.GetUserAddress(request.Uid,request.IsFirst,request.Id)
	list:=make([]*pb.Address,0)
	if err==nil{
		for _,v:=range result{
			list=append(list,&pb.Address{
				Id:v.Id,
				Address:v.Address,
				Area:v.Address,
				Province:v.Province,
				Phone:v.Phone,
				City:v.City,
				Name:v.Name,
				IsDefault:v.IsDefault,
			})
		}
	}
	return &pb.AddressListServiceResponse{List:list},err
}

func (rpc *rpcServer) EditUserAddress(ctx context.Context, request *pb.EditAddressServiceRequest) (result *pb.EmptyResponse, err error) {
	address:=request.Address
	return &pb.EmptyResponse{},orderModel.EditUserAddress(request.Address.Id,request.Uid,
		address.Name,address.Phone,address.Province,address.City,address.Area,address.Address,request.IsDel,address.IsDefault)
}

func (rpc *rpcServer) Create(ctx context.Context, request *pb.OrderRequest) (result *pb.CreateOrderResponse, err error) {
	orderSn := enum.CreateOrderSn(request.OrderType, utils.Get3Code(), false)
	var timeNow = time.Now()
	var endTime = timeNow.Unix() +  request.Time
	list:=make([]orderModel.OrderData,0)
	for _,v:=range request.Response{
		temp:=make([]*orderModel.OrderDetail,0)
		for _,vv:=range v.List{
			temp=append(temp,&orderModel.OrderDetail{
				SkuId:vv.SkuId,
				Price:vv.Price,
				Num:vv.Num,
			})
		}
		list=append(list,orderModel.OrderData{
			SellerId:v.SellerId,
			DeliverPrice:v.DeliverPrice,
			Price:v.Price,
			Detail:temp,
		})
	}
	orderId,err:=orderModel.CreateOrder(orderSn,request.Title,request.Desc,
		request.Phone,request.Province,request.City,request.Area,request.Address,request.Consignee,
		list,timeNow.Unix(),endTime,request.BuyerId,request.Price,request.OrderType,request.ShowOrderList)
	return &pb.CreateOrderResponse{
		Id:        orderId,
		PayAmount: request.PayPrice,
		EndTime:   endTime,
		OrderSn:   orderSn,
	},err
}


func (rpc *rpcServer) EditShopCar(ctx context.Context, request *pb.EditShopCarRequest) (result *pb.EmptyResponse, err error) {
	return &pb.EmptyResponse{},nil
}

func (rpc *rpcServer) CheckPay(ctx context.Context, request *pb.CheckPayRequest) (response *pb.CheckOrderStatusResponse, err error) {
	result, err := orderModel.CreateOrderFromId(request.OrderId, request.OrderSn)
	if err != nil {
		return nil, err
	}
	payStatusResult, err := manager.GetPayService().CheckOrderStatus(ctx, &pb.CheckOrderStatusRequest{
		AppId: request.AppId, ApiKey: os.Getenv("WX_PAY_API_KEY"), MchId: os.Getenv("WX_PAY_MACHINE_ID"), OrderSn: result.OutTradeNo,
	})
	return payStatusResult, err
}

func (rpc *rpcServer) EditStatus(ctx context.Context, request *pb.EditStatusRequest) (response *pb.EditStatusResponse, err error) {
	if request.Status == enum.OrderStatusUserCancel || request.Status == enum.OrderStatusOverTime {
		err:=orderModel.CancelOrder(request.OrderId, request.BuyerId, request.Status)
		return &pb.EditStatusResponse{},err
	} else if request.Status == enum.OrderStatusPay {
		orderSn,orderId,err:=orderModel.Pay(request.OrderId, request.OrderSn)
		return &pb.EditStatusResponse{OrderSn:orderSn,OrderId:orderId},err
	} else if request.Status == enum.OrderStatusSend {
		list:=make([]*orderModel.BatchExpressData,0)
		for _,v:=range request.List{
			list=append(list,&orderModel.BatchExpressData{
				Id:v.Id,
				ExpressCompany:v.ExpressCompany,
				ExpressNumber:v.ExpressNumber,
			})
		}
		err:=orderModel.SendGood(request.SendUid, request.BusinessId, list)
		return &pb.EditStatusResponse{},err
	} else if request.Status == enum.OrderStatusFinish {
		err:=orderModel.StartReceiveGood(request.SubOrderId, request.BuyerId, request.OrderId)
		return &pb.EditStatusResponse{},err
	} else if request.Status == enum.OrderStatusDelete {
		err := orderModel.DelOrder(request.SubOrderId, request.BuyerId)
		return &pb.EditStatusResponse{}, err
	}
	return nil, nil
}


func (rpc *rpcServer) Pay(ctx context.Context, request *pb.OrderPayRequest) (response *pb.PayServiceResponse, err error) {
	//获取支付需要的数据
	result, err := orderModel.CreateOrderFromId(request.OrderId, request.OrderSn)
	if err == nil {
		isNew := false
		if result.PrepayId == "" || (result.TotalFree != request.PayTotalFree && result.PrepayId != "") {
			if result.TotalFree != request.PayTotalFree && result.PrepayId != "" {
				oNo := enum.CreateOrderSn(result.OrderType, utils.Get3Code(),  false)
				err = orderModel.UpdateOrderOutTradeNo(result.Id, oNo)
				if err != nil {
					xlog.ErrorP(err)
					return nil, err
				}
				result.OutTradeNo = oNo
			}
			isNew = true
		}
		if request.OnlyGetPayInfo && !isNew {
			//不需要重新请求获取支付信息
			return &pb.PayServiceResponse{
				PayFree: result.TotalFree,
				Status:  result.Status,
			}, nil
		}
		if request.PayTotalFree!=0{
			result.TotalFree=request.PayTotalFree
		}


		payResult, err := manager.GetPayService().CreatePay(ctx, createPayInfo(result, request.OpenId, request.AppId, request.NotifyUrl, result.PrepayId, request.PayType,request.ChannelType))
		if err != nil {
			return nil, err
		} else {
			if payResult.Pay {
				return &pb.PayServiceResponse{
					Pay: payResult.Pay,
				}, nil
			}
			//更新预订单id
			err = orderModel.UpdateOrderPrepayId(result.Id,request.ChannelType, result.OrderSn, payResult.PrepayId)
			if err != nil {
				return nil, err
			}
			payResult.OrderSn = result.OrderSn
			payResult.OrderId = result.Id
			payResult.PayFree = result.TotalFree
			payResult.Status = result.Status
			return payResult, nil
		}
	} else {
		xlog.ErrorP(err)
		return nil, err
	}
}


func (rpc *rpcServer) GetOrderNum(ctx context.Context, request *pb.GetOrderNumRequest) (*pb.GetOrderNumResponse, error) {
	result, skuResult, err := orderModel.GetAllSpuSales(request.SpuSkuMap, request.Refresh)
	if err != nil {
		return nil, err
	}
	return &pb.GetOrderNumResponse{
		Result:    result,
		SkuResult: skuResult,
	}, nil

}


func (rpc *rpcServer) GetOrders(ctx context.Context, request *pb.GetOrderRequest) (response *pb.GetOrderResponse, err error) {
	return orders(ctx, request.SellerId, request.BuyerId, request.Page, request.Size, request.StartTime, request.EndTime, request.SubOrderId, request.OrderId,request.BuyerType,
		request.QueryType,request.ShowOrder,request.SpuName,request.BuyerName,request.Brand,request.Operator,request.OperatorNo,request.OrderSn,request.ExpressNumber,request.BuyerNumber,request.Phone, request.Status, request.OrderTypes)
}

func orders(ctx context.Context,  sellerId, buyerId, page, size, startTime, endTime, subOrderId, orderId,buyerType,queryType,showOrderType int64, spuName,buyerName,brand,operator,operatorNo,orderSn,expressNumber,buyerNumber,phone string, status, orderTypes []int64) (*pb.GetOrderResponse, error) {
	defaultName:="游客"
	skuIdList := make([]int64, 0)
	uidList := make([]int64, 0)
	sellerIdList := make([]int64, 0)
	brandIdList := make([]int64, 0)
	if buyerNumber!=""{
		uResult, err := manager.GetUserService().FindIdFromQuery(ctx, &pb.FindUserFromQueryRequest{Query: buyerNumber, No: true})
		if err == nil {
			uidList = append(uidList, uResult.List...)
		}
	}
	if buyerName!=""{
		if strings.Index(defaultName,buyerName)!=-1{
			uidList = append(uidList, 0)
		}
		uResult, err := manager.GetUserService().FindIdFromQuery(ctx, &pb.FindUserFromQueryRequest{Query: buyerName, Name: true})
		if err == nil {
			uidList = append(uidList, uResult.List...)
		}

	}
	if spuName != "" || len(brandIdList)>0{
		skuIdResult, err := manager.GetGoodService().GetSkuIdQuery(ctx, &pb.GetSkuIdQueryRequest{Query: spuName,BrandList:brandIdList})
		if err == nil {
			skuIdList = skuIdResult.List
		}
	}
	if sellerId!=0{
		sellerIdList=append(sellerIdList,sellerId)
	}
	orderList, total, err := orderModel.GetOrderList(buyerId, subOrderId, orderId, startTime, endTime, page, size,buyerType,queryType,showOrderType, status,orderTypes, skuIdList, sellerIdList, uidList, phone,expressNumber,orderSn)
	if err != nil {
		return nil, err
	}
	orderMap := make(map[int64]*pb.GetOrderData)
	userOrderMap := make(map[int64][]int64)
	orderSpuMap := make(map[string]*pb.GetOrderSpu)
	list := make([]*pb.GetOrderData, 0)
	//未支付订单map，用来做合并的
	unPayOrder := make(map[int64]*pb.GetOrderData)

	for _, info := range orderList {
		var temp *pb.GetOrderData
		ok := false
		if buyerId != 0 {
			if info.OrderStatus == enum.OrderStatusUnPay {
				//未支付的合并显示
				temp, ok = unPayOrder[info.OrderID]
			}
		}
		if !ok {
			temp = &pb.GetOrderData{
				OrderSn:        info.OrderSN,
				Status:         info.OrderStatus,
				Consignee:      info.Consignee,
				Phone:          info.Phone,
				Province:       info.Province,
				City:           info.City,
				Area:           info.Area,
				Address:        info.Address,
				OrderMount:     info.OrderAmount,
				PayMount:       info.PayAmount,
				ExpressNumber:  info.ExpressNumber,
				ExpressCompany: info.ExpressCompany,
				Message:        info.Message,
				AdminMessage:   info.Remark,
				SubOrderId:     info.SubOrderID,
				OrderId:        info.OrderID,
				OrderTime:      info.CreateTime,
				PayTime:        info.PayTime,
				EndTime:        info.EndTime,
				TransactionId:  info.TransactionId,
				IsHaveComment:  info.HaveCommented,
				DeliveryPrice:  info.DeliveryPrice,
				PayType:        info.PayType,
				PayChannel:     info.PayChannel,
				SellerId:       info.SellerID,
				CouponPrice:    info.CouponPrice,
				BuyerId:        info.BuyerId,
				OrderType:      info.OrderType,
				BuyType:        info.BuyType,
				ShowOrder:info.ShowOrder,
				List:           make([]*pb.GetOrderSpu, 0),
			}
			if info.BuyerId == 0 {
				temp.UserName = defaultName
			}
			list = append(list, temp)
			if info.OrderStatus == enum.OrderStatusUnPay {
				unPayOrder[info.OrderID] = temp
			}
		} else if temp != nil {
			//未支付的合并
			temp.DeliveryPrice += info.DeliveryPrice
			temp.PayMount += info.PayAmount
			temp.OrderMount += info.OrderAmount
		}
		//获取购买者信息
		userOrderList, ok := userOrderMap[info.BuyerId]
		if !ok {
			userOrderList = make([]int64, 0)
		}
		userOrderList = append(userOrderList, info.SubOrderID)
		userOrderMap[info.BuyerId] = userOrderList

		orderMap[info.SubOrderID] = temp
	}
	ids := make([]int64, 0)
	for k := range userOrderMap {
		ids = append(ids, k)
	}
	if len(ids) > 0 {
		//获取用户信息
		temp, err := manager.GetUserService().GetUserList(ctx, &pb.GetUserListRequest{IdList: ids})
		if err == nil {
			for _, info := range temp.List {
				temps, ok := userOrderMap[info.Uid]
				if ok {
					for _, id := range temps {
						orderResult, ok := orderMap[id]
						if ok {
							orderResult.UserSn = info.No
							orderResult.Head = info.Head
							orderResult.UserName = info.Name
						}
					}
				}
			}
		}
	}



	if len(orderMap) > 0 {
		subOrderIdList := make([]int64, 0)
		for k,v := range orderMap {
			if v.ShowOrder{
				subOrderIdList = append(subOrderIdList, k)
			}
		}
		if len(subOrderIdList)>0{
			skuResult, err := orderModel.GetOrderGoodDetailBySubOrderID(subOrderIdList)
			if err == nil {
				skuOptionsMap := make(map[int64]string)
				spuIdList := make([]int64, 0)
				skuIdList := make([]int64, 0)
				spuResultMap := make(map[int64]*pb.Spu)
				skuSpuMap := make(map[int64]int64)
				for _, info := range skuResult {
					skuIdList = append(skuIdList, info.SkuId)
				}
				if len(skuIdList)>0{
					optionsResult, err := manager.GetGoodService().GetSkuOptions(ctx, &pb.GetSkuOptionsRequest{SkuIdList: skuIdList})
					if err == nil {
						for k, info := range optionsResult.Result {
							tempList := make([]string, 0)
							for _, value := range info.List {
								tempList = append(tempList, value.OptionValue)
								spuIdList = append(spuIdList, value.SpuId)
								skuSpuMap[k] = value.SpuId
							}
							skuOptionsMap[k] = strings.Join(tempList, "-")
						}
					}
					spuResult, err := manager.GetGoodService().GetSpuList(ctx, &pb.GetSpuListRequest{List: spuIdList})
					if err == nil {
						for _, info := range spuResult.List {
							spuResultMap[info.Id] = info
						}
						for _, info := range skuResult {
							if temp, ok := orderMap[info.SubOrderID]; ok {
								if spuId, ok := skuSpuMap[info.SkuId]; ok {
									if spuInfo, ok := spuResultMap[spuId]; ok {
										spu, ok := orderSpuMap[fmt.Sprintf("%d+%d", info.SubOrderID, spuInfo.Id)]
										if !ok {
											spu = &pb.GetOrderSpu{
												Id:    spuInfo.Id,
												Name:  spuInfo.Name,
												Photo: spuInfo.Photo,
												List:  make([]*pb.GetOrderSku, 0),
											}
											orderSpuMap[fmt.Sprintf("%d+%d", info.SubOrderID, spuInfo.Id)] = spu
											temp.List = append(temp.List, spu)
										}
										spu.List = append(spu.List, &pb.GetOrderSku{
											SkuId: info.SkuId,
											Num:   info.Num,
											Price: info.Price,
											Spec:  skuOptionsMap[info.SkuId],
										})
									}
								}
							}
						}
					}
				}
			} else {
				xlog.ErrorP(err)
				return nil, err
			}
		}
	}
	return &pb.GetOrderResponse{List: list, Total: total}, nil
}
// PostOrderRefund 退款退货
func (rpc *rpcServer) PostOrderRefund(ctx context.Context, request *pb.PostRefundReq) (*pb.EmptyResponse, error) {
	var result = &pb.EmptyResponse{}
	subOrder, err := orderModel.GetSubOrderByID(request.SubOrderId,request.OrderId)
	if err != nil {
		return result, xlog.Error(err)
	}
	if subOrder.BuyerId != request.BuyerId && !request.Pass {
		return result, xlog.Error("非当前用户订单")
	}
	if !request.Cancel && subOrder.OrderStatus != enum.OrderStatusPay && subOrder.OrderStatus != enum.OrderStatusSend {
		return result, xlog.Error("已申请退款")
	}
	if request.Cancel {
		return result, nil
	} else {

		orderDetails, err := orderModel.GetOrderGoodDetailBySubOrderID([]int64{subOrder.SubOrderID})
		if err != nil {
			return result, err
		}
		if request.All {
			//全部退款退货，不传递List
			request.List = make([]*pb.PostRefundReqList, 0)
			for _, info := range orderDetails {
				request.List = append(request.List, &pb.PostRefundReqList{
					OgID: info.ID,
					Num:  info.Num,
				})
			}
		}
		var orderDetailMap = make(map[int64]orderModel.OrderDetail) // 购买的商品总数
		for _, o := range orderDetails {
			orderDetailMap[o.ID] = o
		}
		refunds, err := orderModel.GetRefundBySubOrderID(subOrder.SubOrderID)
		if err != nil {
			return result, err
		}
		var refundIDs []int64
		for _, v := range refunds {
			if v.Status == enum.RefundStatusFail {
				continue
			}
			refundIDs = append(refundIDs, v.ID)
		}
		refundDetails, err := orderModel.GetRefundDetailByRefundID(refundIDs)
		if err != nil {
			return result, err
		}
		// 已经退款的商品总数
		var refundDetailMap = make(map[int64]int64)
		for _, r := range refundDetails {
			if v, ok := refundDetailMap[r.OgID]; ok {
				refundDetailMap[r.OgID] = r.Num + v
			} else {
				refundDetailMap[r.OgID] = r.Num
			}
		}

		// 退全部商品
		var amount int64 = 0
		for _, r := range request.List {
			if v, ok := orderDetailMap[r.OgID]; ok {
				if r.Num+refundDetailMap[r.OgID] > v.Num {
					return result, xlog.Errorf("退货数量不正确")
				}
				amount += v.PayAmount
			} else {
				return result, xlog.Errorf("ogid：%d，不存在", r.OgID)
			}
		}

		// 新增
		refundOrderSn := enum.CreateRefundOrderSn(subOrder.OrderType, subOrder.BuyType, utils.Get3Code())
		err = orderModel.CreateRefund(subOrder.SubOrderID, subOrder.OrderStatus, amount,request.RefundMethod,request.RefundType, request.List,request.PhotoList, refundOrderSn,request.Reason,request.ExpressCompany,request.ExpressNumber,request.Explain)
		_, err1 := rpc.PutOrderRefundStatus(ctx, &pb.PutOrderRefundStatusReq{SubOrderId: subOrder.SubOrderID, Pass: true, AdminNo: ""})
		if err1 != nil {
			return result, err1
		}
	}
	return result, nil
}


// PutOrderRefundStatus 退款退货审核
func (rpc *rpcServer) PutOrderRefundStatus(ctx context.Context, request *pb.PutOrderRefundStatusReq) (*pb.EmptyResponse, error) {
	var result = &pb.EmptyResponse{}
	// 获取订单状态
	subOrder, err := orderModel.GetSubOrderByID(request.SubOrderId,request.OrderId)
	if err != nil {
		return result, err
	}
	if subOrder.OrderStatus != enum.OrderStatusUnRefund && subOrder.OrderStatus != enum.OrderStatusPayRefund {
		return result, xlog.Error("订单状态不正确")
	}

	// 更新订单状态
	orderStatus := enum.OrderStatusCloseRefund
	if subOrder.OrderStatus == enum.OrderStatusPayRefund {
		orderStatus = enum.OrderStatusCloseRefundReturn
	}
	var refundStatus = enum.RefundStatusFinish
	if !request.Pass {
		if subOrder.OrderStatus == enum.OrderStatusPayRefund {
			orderStatus = enum.OrderStatusPay
		} else {
			orderStatus = enum.OrderStatusSend
		}
		refundStatus = enum.RefundStatusFail
	} else {

		refunds, err := orderModel.GetRefundBySubOrderID(request.SubOrderId)
		total := int64(0)
		if err != nil {
			return nil, err
		}
		for _, v := range refunds {
			if v.Status == enum.RefundStatusVerify {
				total += v.Amount
				if total>subOrder.PayAmount{
					total=subOrder.PayAmount
				}
				_, err = manager.GetPayService().RefundOrder(ctx, &pb.RefundOrderRequest{
					OrderSn:     subOrder.OrderSN,
					TotalFree:   subOrder.PayAmount,
					RefundFree:  total,
					OutRefundNo: v.OutRefundNo,
					ChannelType:subOrder.PayChannel,
				})
				if err != nil {
					return nil, err
				}
			}
		}
	}

	err = orderModel.PutOrderRefundStatus(refundStatus, orderStatus, subOrder.SubOrderID, request.Remark, request.AdminNo)
	if err != nil {
		if request.Pass {
			//回滚

		}
		return result, err
	}
	if request.Pass {
		//成功后确定事务消息

	}

	return result, nil
}


func createPayInfo(orderResult *orderModel.OrderResponse, openId, appId, url, prepayId string, payType,channelType int64) *pb.PayServiceRequest {
	request := &pb.PayServiceRequest{AppId: appId, ApiKey: os.Getenv("WX_PAY_API_KEY")}
	request.MchId = os.Getenv("WX_PAY_MACHINE_ID")
	request.NotifyUrl = url
	request.OpenId = openId
	request.Title = orderResult.Title
	request.OrderSn = orderResult.OrderSn
	request.Total = orderResult.TotalFree
	request.PayType = payType
	request.PrepayId = prepayId
	request.OutTradeNo = orderResult.OutTradeNo
	request.ChannelType=channelType
	return request
}
