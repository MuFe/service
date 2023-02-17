package wxPay_test

import (
	"context"
	"log"
	"time"
	utils2 "mufe_service/camp/wx/util"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"

)

func ExampleAppApiService_CloseOrder() {
	ctx := context.Background()
	client, err:= utils2.GetWxClient(ctx)
	if err != nil {
		log.Printf("new wechat pay client err:%s", err)
	}

	svc := app.AppApiService{Client: client}
	result, err := svc.CloseOrder(ctx,
		app.CloseOrderRequest{
			OutTradeNo: core.String("OutTradeNo_example"),
			Mchid:      core.String("1230000109"),
		},
	)

	if err != nil {
		// 处理错误
		log.Printf("call CloseOrder err:%s", err)
	} else {
		// 处理返回结果
		log.Printf("status=%d", result.Response.StatusCode)
	}
}

func ExampleAppApiService_Prepay() {
	ctx := context.Background()
	client, err:= utils2.GetWxClient(ctx)
	if err != nil {
		log.Printf("new wechat pay client err:%s", err)
	}

	svc := app.AppApiService{Client: client}
	resp, result, err := svc.Prepay(ctx,
		app.PrepayRequest{
			Appid:         core.String("wxd678efh567hg6787"),
			Mchid:         core.String("1230000109"),
			Description:   core.String("Image形象店-深圳腾大-QQ公仔"),
			OutTradeNo:    core.String("1217752501201407033233368018"),
			TimeExpire:    core.Time(time.Now()),
			Attach:        core.String("自定义数据说明"),
			NotifyUrl:     core.String("https://www.weixin.qq.com/wxpay/pay.php"),
			GoodsTag:      core.String("WXG"),
			LimitPay:      []string{"LimitPay_example"},
			SupportFapiao: core.Bool(false),
			Amount: &app.Amount{
				Currency: core.String("CNY"),
				Total:    core.Int64(100),
			},
			Detail: &app.Detail{
				CostPrice: core.Int64(608800),
				GoodsDetail: []app.GoodsDetail{app.GoodsDetail{
					GoodsName:        core.String("iPhoneX 256G"),
					MerchantGoodsId:  core.String("ABC"),
					Quantity:         core.Int64(1),
					UnitPrice:        core.Int64(828800),
					WechatpayGoodsId: core.String("1001"),
				}},
				InvoiceId: core.String("wx123"),
			},
			SceneInfo: &app.SceneInfo{
				DeviceId:      core.String("013467007045764"),
				PayerClientIp: core.String("14.23.150.211"),
				StoreInfo: &app.StoreInfo{
					Address:  core.String("广东省深圳市南山区科技中一道10000号"),
					AreaCode: core.String("440305"),
					Id:       core.String("0001"),
					Name:     core.String("腾讯大厦分店"),
				},
			},
			SettleInfo: &app.SettleInfo{
				ProfitSharing: core.Bool(false),
			},
		},
	)

	if err != nil {
		// 处理错误
		log.Printf("call Prepay err:%s", err)
	} else {
		// 处理返回结果
		log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
	}
}

func ExampleAppApiService_QueryOrderById() {
	ctx := context.Background()
	client, err:= utils2.GetWxClient(ctx)
	if err != nil {
		log.Printf("new wechat pay client err:%s", err)
	}

	svc := app.AppApiService{Client: client}
	resp, result, err := svc.QueryOrderById(ctx,
		app.QueryOrderByIdRequest{
			TransactionId: core.String("TransactionId_example"),
			Mchid:         core.String("Mchid_example"),
		},
	)

	if err != nil {
		// 处理错误
		log.Printf("call QueryOrderById err:%s", err)
	} else {
		// 处理返回结果
		log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
	}
}

func ExampleAppApiService_QueryOrderByOutTradeNo() {
	ctx := context.Background()
	client, err:= utils2.GetWxClient(ctx)
	if err != nil {
		log.Printf("new wechat pay client err:%s", err)
	}

	svc := app.AppApiService{Client: client}
	resp, result, err := svc.QueryOrderByOutTradeNo(ctx,
		app.QueryOrderByOutTradeNoRequest{
			OutTradeNo: core.String("OutTradeNo_example"),
			Mchid:      core.String("Mchid_example"),
		},
	)

	if err != nil {
		// 处理错误
		log.Printf("call QueryOrderByOutTradeNo err:%s", err)
	} else {
		// 处理返回结果
		log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
	}
}
