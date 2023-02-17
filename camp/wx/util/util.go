package util

import (
	"context"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"os"
	"mufe_service/camp/xlog"
)

// TokenAPI 获取带 token 的 API 地址
func GetWxClient(ctx context.Context ) (*core.Client, error) {
	mchID:=os.Getenv("WX_PAY_MACHINE_ID")
	mchAPIv3Key:=os.Getenv("WX_PAY_API_KEY")
	mchCertificateSerialNumber:=os.Getenv("NUMBER")

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath("./apiclient_key.pem")
	if err != nil {
		xlog.ErrorP(err)
	}

	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	wxClient, err := core.NewClient(ctx, opts...)
	return wxClient, nil
}


func GetNotify(ctx context.Context )*notify.Handler{
	mchID:=os.Getenv("WX_PAY_MACHINE_ID")
	mchAPIv3Key:=os.Getenv("WX_PAY_API_KEY")
	mchCertificateSerialNumber:=os.Getenv("NUMBER")

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath("./apiclient_key.pem")
	if err != nil {
		xlog.ErrorP(err)
	}
	err = downloader.MgrInstance().RegisterDownloaderWithPrivateKey(ctx, mchPrivateKey, mchCertificateSerialNumber, mchID, mchAPIv3Key)
	if err != nil {
		xlog.ErrorP(err)
	}
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(mchID)
	handler:=notify.NewNotifyHandler(mchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	return handler
}
