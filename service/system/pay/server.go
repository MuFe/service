package wxPay

import (
	"context"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/shopspring/decimal"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"io/ioutil"
	"os"
	"mufe_service/camp/enum"
	"mufe_service/camp/service"
	utils3 "mufe_service/camp/utils"
	utils2 "mufe_service/camp/wx/util"
	pb "mufe_service/jsonRpc"
	"strconv"
	"time"

	"mufe_service/camp/xlog"
)

type Charset string

type rpcServer struct {
}

var wxClient *core.Client
var aliClient *alipay.Client
func init() {
	nSer := &rpcServer{}
	pb.RegisterPayServiceServer(service.GetRegisterRpc(), nSer)
	createWxClient()
	createAliClient()
}

func createWxClient(){
	var err error
	ctx := context.Background()
	wxClient, err= utils2.GetWxClient(ctx)
	xlog.ErrorP(err)
}

func createAliClient(){
	appID:=os.Getenv("ALI_PAY_APP_ID")
	privateKey, err := ioutil.ReadFile("./alipay_private.txt")
	if err != nil {
		xlog.ErrorP(err)
	}
	aliClient,err = alipay.NewClient(appID, string(privateKey), true)
	//publicKey, err := ioutil.ReadFile("./alipay_public.txt")
	if err != nil {
		xlog.ErrorP(err)
	}
	if os.Getenv("MODEL")=="test"{
		aliClient.DebugSwitch=gopay.DebugOn
		aliClient.NotifyUrl="https://www.sgsports-test.com:8443/appApi/test/re"
	}
	//temp:=xrsa.FormatAlipayPublicKey(string(publicKey))
	//aliClient.AutoVerifySign([]byte(temp))

}

func (rpc *rpcServer) CreatePay(ctx context.Context, request *pb.PayServiceRequest) (result *pb.PayServiceResponse, err error) {
	if request.ChannelType==enum.WEIXINPAY{
	 	return wxPay(ctx,request)
	 }else{
	 	return aliPay(ctx,request)
	 }

}

func (rpc *rpcServer) RefundOrder(ctx context.Context, request *pb.RefundOrderRequest) (result *pb.EmptyResponse, err error) {
	if request.ChannelType==enum.WEIXINPAY{
		svc :=refunddomestic.RefundsApiService{Client:wxClient}
		_, _, err := svc.Create(ctx,refunddomestic.CreateRequest{
			OutTradeNo:core.String(request.OrderSn),
			OutRefundNo:core.String(request.OutRefundNo),
			Amount:&refunddomestic.AmountReq{
				Refund:core.Int64(request.RefundFree),
				Total:core.Int64(request.TotalFree),
				Currency:core.String("CNY"),
			},
		})
		if err!=nil{
			if ne, ok := err.(*core.APIError); ok {
				return nil,xlog.Error(ne.Message)
			}
			return nil,xlog.Error(err)
		}else {
			return &pb.EmptyResponse{}, nil
		}
	}else{
		bm := make(gopay.BodyMap)
		decimalValue2 := decimal.NewFromInt(request.RefundFree)
		decimalValue2 = decimalValue2.Div(decimal.NewFromInt(100))
		num2,_ := decimalValue2.Float64()
		bm.
			Set("out_trade_no", request.OrderSn).
			Set("refund_amount", num2)
		_,err:=aliClient.TradeRefund(ctx,bm)
		if err!=nil{
			xlog.ErrorP(err)
			return nil,xlog.Error(err)
		}else {
			xlog.Info(1111)
			return &pb.EmptyResponse{}, nil
		}
	}
}

func (rpc *rpcServer) CheckOrderStatus(ctx context.Context, request *pb.CheckOrderStatusRequest) (result *pb.CheckOrderStatusResponse, err error) {
	svc := jsapi.JsapiApiService{Client: wxClient}
	resp, _, err :=svc.QueryOrderByOutTradeNo(ctx,jsapi.QueryOrderByOutTradeNoRequest{
		Mchid:core.String(request.MchId),
		OutTradeNo:core.String(request.OrderSn),
	})
	status := enum.PayStatus_No
	if err!=nil{
		if ne, ok := err.(*core.APIError); ok {
			return nil,xlog.Error(ne.Message)
		}
		return nil,xlog.Error(err)
	}
	payTime := int64(0)
	if *resp.TradeState == "NOTPAY" {
		status = enum.PayStatus_No
	} else if *resp.TradeState == "SUCCESS" {
		status = enum.PayStatus_Su
		datetime, _ := time.ParseInLocation(*resp.SuccessTime, "20060102150405", time.Local)
		payTime = datetime.Unix()
	} else if *resp.TradeState == "CLOSED" {
		status = enum.PayStatus_Close
	}
	return &pb.CheckOrderStatusResponse{OutTradeNo: *resp.OutTradeNo, TransactionId: *resp.TransactionId, Status: status, PayTime: payTime, TotalFee: int64(*resp.Amount.Total)}, nil
}

func (rpc *rpcServer) CheckPayData(ctx context.Context, request *pb.PayDataRequest) (result *pb.PayDataResponse, err error) {


	return &pb.PayDataResponse{}, nil
}



func wxPay(ctx context.Context, request *pb.PayServiceRequest) (result *pb.PayServiceResponse, err error) {
	if request.PrepayId != "" && (request.PayType == enum.MiniPayType || request.PayType == enum.JSAPIPayType) {
		resp:= new(jsapi.PrepayWithRequestPaymentResponse)
		resp.PrepayId = core.String(request.PrepayId)
		resp.SignType = core.String("RSA")
		resp.Appid = core.String(request.AppId)
		resp.TimeStamp = core.String(strconv.FormatInt(time.Now().Unix(), 10))
		nonce, err := utils.GenerateNonce()
		if err != nil {
			return nil,xlog.Error(err)
		}
		resp.NonceStr = core.String(nonce)
		resp.Package = core.String("prepay_id=" +request.PrepayId)
		message := fmt.Sprintf("%s\n%s\n%s\n%s\n", *resp.Appid, *resp.TimeStamp, *resp.NonceStr, *resp.Package)
		signatureResult, err := wxClient.Sign(ctx, message)
		if err != nil {
			return nil,xlog.Error(err)
		}
		resp.PaySign = core.String(signatureResult.Signature)
		return &pb.PayServiceResponse{
			AppId:request.AppId,
			PrepayId:*resp.PrepayId,
			Package:*resp.Package,
			NonceStr:*resp.NonceStr,
			TimeStamp:*resp.TimeStamp,
			PaySign:*resp.PaySign,
			SignType:*resp.SignType,
		}, err
	} else {
		if request.PayType == enum.WEBPayType {
			svc := native.NativeApiService{Client: wxClient}
			resp, _, err := svc.Prepay(ctx,native.PrepayRequest{
				Appid:core.String(request.AppId),
				Mchid:core.String(request.MchId),
				Description:core.String(request.Title),
				OutTradeNo:core.String(request.OutTradeNo),
				NotifyUrl:core.String(request.NotifyUrl),
				Attach:core.String(request.OrderSn),
				Amount:&native.Amount{
					Total:core.Int64(request.Total),
				},
			})
			if err!=nil{
				if ne, ok := err.(*core.APIError); ok {
					return nil,xlog.Error(ne.Message)
				}
				return nil,xlog.Error(err)
			}else {

				return &pb.PayServiceResponse{CodeUrl: *resp.CodeUrl}, err
			}
		} else if request.PayType==enum.AppPayType{
			svc := app.AppApiService{Client: wxClient}
			resp, _, err := svc.PrepayWithRequestPayment(ctx,app.PrepayRequest{
				Appid:core.String(request.AppId),
				Mchid:core.String(request.MchId),
				Description:core.String(request.Title),
				OutTradeNo:core.String(request.OutTradeNo),
				NotifyUrl:core.String(request.NotifyUrl),
				Attach:core.String(request.OrderSn),
				Amount:&app.Amount{
					Total:core.Int64(request.Total),
				},
			})
			if err!=nil{
				if ne, ok := err.(*core.APIError); ok {
					return nil,xlog.Error(ne.Message)
				}
				return nil,xlog.Error(err)
			}else {
				return &pb.PayServiceResponse{
					AppId:request.AppId,
					PrepayId:*resp.PrepayId,
					Package:*resp.Package,
					NonceStr:*resp.NonceStr,
					TimeStamp:*resp.TimeStamp,
					PaySign:*resp.Sign,
					PartnerId:*resp.PartnerId,
				}, err
			}
		}else if request.PayType == enum.H5PayType{
			svc := h5.H5ApiService{Client: wxClient}
			resp, _, err := svc.Prepay(ctx,h5.PrepayRequest{
				Appid:core.String(request.AppId),
				Mchid:core.String(request.MchId),
				Description:core.String(request.Title),
				OutTradeNo:core.String(request.OutTradeNo),
				NotifyUrl:core.String(request.NotifyUrl),
				Attach:core.String(request.OrderSn),
				Amount:&h5.Amount{
					Total:core.Int64(request.Total),
				},
				SceneInfo:&h5.SceneInfo{
					PayerClientIp:core.String(utils3.LocalIP()),
					H5Info: &h5.H5Info{
						AppName:     core.String(request.Title),
						AppUrl:      core.String(request.Title),
						BundleId:    core.String(request.Title),
						PackageName: core.String(request.Title),
						Type:        core.String(request.Title),
					},
				},
			})
			if err!=nil{
				if ne, ok := err.(*core.APIError); ok {
					return nil,xlog.Error(ne.Message)
				}
				return nil,xlog.Error(err)
			}else {

				return &pb.PayServiceResponse{CodeUrl: *resp.H5Url}, err
			}
		} else if request.PayType==enum.JSAPIPayType||request.PayType==enum.MiniPayType{
			svc := jsapi.JsapiApiService{Client: wxClient}
			resp, _, err := svc.PrepayWithRequestPayment(ctx,jsapi.PrepayRequest{
				Appid:core.String(request.AppId),
				Mchid:core.String(request.MchId),
				Description:core.String(request.Title),
				OutTradeNo:core.String(request.OutTradeNo),
				NotifyUrl:core.String(request.NotifyUrl),
				Attach:core.String(request.OrderSn),
				Amount:&jsapi.Amount{
					Total:core.Int64(request.Total),
				},
				Payer:&jsapi.Payer{
					Openid:core.String(request.OpenId),
				},
			})
			if err!=nil{
				if ne, ok := err.(*core.APIError); ok {
					return nil,xlog.Error(ne.Message)
				}
				return nil,xlog.Error(err)
			}else {
				return &pb.PayServiceResponse{
					AppId:request.AppId,
					PrepayId:*resp.PrepayId,
					Package:*resp.Package,
					NonceStr:*resp.NonceStr,
					TimeStamp:*resp.TimeStamp,
					PaySign:*resp.PaySign,
					SignType:*resp.SignType,
				}, err
			}
		}
	}
	return &pb.PayServiceResponse{

	},nil
}

func aliPay(ctx context.Context,  request *pb.PayServiceRequest) (result *pb.PayServiceResponse, err error) {
	if request.PrepayId != "" && (request.PayType == enum.MiniPayType || request.PayType == enum.JSAPIPayType) {

		return &pb.PayServiceResponse{

		}, err
	} else {
		bm := make(gopay.BodyMap)
		decimalValue2 := decimal.NewFromInt(request.Total)
		decimalValue2 = decimalValue2.Div(decimal.NewFromInt(100))
		num2,_ := decimalValue2.Float64()
			bm.Set("out_trade_no", request.OutTradeNo).
			Set("total_amount", num2).Set("subject", request.Title)

		if request.PayType == enum.AppPayType {
			result,err:=aliClient.TradeAppPay(ctx,bm)
			if err!=nil{
				return nil,xlog.Error(err)
			}else{
				xlog.Info(result)
				return &pb.PayServiceResponse{CodeUrl:result}, err
			}
		} else if request.PayType==enum.WEBPayType{
			bm.Set("product_code", "FAST_INSTANT_TRADE_PAY")
			result,err:=aliClient.TradePagePay(ctx,bm)
			if err!=nil{
				return nil,xlog.Error(err)
			}else{
				return &pb.PayServiceResponse{CodeUrl:result}, err
			}
		}
	}
	return &pb.PayServiceResponse{

	},nil
}
