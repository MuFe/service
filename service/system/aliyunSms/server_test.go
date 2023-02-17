package aliyunSms

import (
	"fmt"
	"github.com/agiledragon/gomonkey"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/smartystreets/goconvey/convey"
	app "mufe_service/jsonRpc"
	"reflect"
	"testing"
)

func TestSend(t *testing.T) {
	convey.Convey("TestSend", t, func() {
		client:=&dysmsapi.Client{}
		//下面正常逻辑
		convey.Convey("return il", func() {
			//open返回nil
			patches := gomonkey.ApplyFunc(dysmsapi.NewClientWithAccessKey, func(_,_,_ string) (*dysmsapi.Client, error) {
				return client, nil
			})
			defer patches.Reset()
			//readAll返回nil
			patches.ApplyMethod(reflect.TypeOf(client),"SendSms", func(*dysmsapi.Client,*dysmsapi.SendSmsRequest) (*dysmsapi.SendSmsResponse, error) {
				return &dysmsapi.SendSmsResponse{
					Code:"111",
				}, nil
			})
			flag, err := send("15816138010","","","")
			convey.So(err, convey.ShouldBeNil)
			temp:=&app.SendResponse{Result:"111"}
			convey.So(flag.Result, convey.ShouldEqual,temp.Result)
		})
		////返回错误
		convey.Convey("return error when os.open error", func() {
			//open返回error
			patches := gomonkey.ApplyMethod(reflect.TypeOf(client),"SendSms", func(*dysmsapi.Client,*dysmsapi.SendSmsRequest) (*dysmsapi.SendSmsResponse, error) {
				return nil, fmt.Errorf("send error")
			})
			defer patches.Reset()
			flag, err := send("15816138010","","","")
			convey.So(err, convey.ShouldNotBeNil)
			convey.So(flag, convey.ShouldEqual, nil)
		})
	})

}
