package aliyunSms

import (
	"context"
	"os"
	"mufe_service/camp/service"
	pb "mufe_service/jsonRpc"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type rpcServer struct {
}



func init() {
	AccessKeyID = os.Getenv("SMSACCESSKEYID")
	AccessKeySecret = os.Getenv("SMSACCESSKEYSECRET")
	pb.RegisterSendSmsServiceServer(service.GetRegisterRpc(), &rpcServer{})
}

func (rpc *rpcServer) SendSms(ctx context.Context, request *pb.SendRequest) (result *pb.SendResponse, err error) {
	if request.Key!=""{
		AccessKeyID=request.KeyId
		AccessKeySecret=request.Key
	}
	return send(request.Phone, request.ParamStr, request.SingName, request.TemplateCode)
}

var (
	AccessKeyID     string
	AccessKeySecret string
)

func send(phone, paramStr, singName, tempCode string) (*pb.SendResponse, error) {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", AccessKeyID, AccessKeySecret)
	if err != nil {
		return nil, err
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phone
	request.SignName = singName
	request.TemplateCode = tempCode

	//request.TemplateParam = fmt.Sprintf("{code:\"%s\"}", code)
	request.TemplateParam = paramStr

	response, err := client.SendSms(request)
	if err != nil {
		return nil, err
	} else {
		return &pb.SendResponse{Result: response.Code}, nil
	}
}

