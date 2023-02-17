package football

import (
	"github.com/agiledragon/gomonkey"
	"github.com/smartystreets/goconvey/convey"
	"io"
	"io/ioutil"
	"net/http"
	app "mufe_service/jsonRpc"
	"reflect"
	"testing"
)

func TestSend(t *testing.T) {
	convey.Convey("TestSend", t, func() {
		client:=&http.Client{}
		//下面正常逻辑
		convey.Convey("return il", func() {
			patches := gomonkey.ApplyFunc(http.NewRequest, func(string,string,io.Reader) (*http.Request, error) {
				return nil, nil
			})
			patches.ApplyMethod(reflect.TypeOf(client),"Do", func(*http.Client,*http.Request) (*http.Response, error) {
				return &http.Response{}, nil
			})
			patches.ApplyFunc(ioutil.ReadAll, func(_ io.Reader) ([]byte, error) {
				return nil, nil
			})
			defer patches.Reset()
			flag, err := send("15816138010","11")
			convey.So(err, convey.ShouldBeNil)
			temp:=&app.GetDataResult{}
			convey.So(flag, convey.ShouldEqual,temp)
		})
		////返回错误
		//convey.Convey("return error when os.open error", func() {
		//	//open返回error
		//	patches := gomonkey.ApplyMethod(reflect.TypeOf(client),"SendSms", func(*dysmsapi.Client,*dysmsapi.SendSmsRequest) (*dysmsapi.SendSmsResponse, error) {
		//		return nil, fmt.Errorf("send error")
		//	})
		//	defer patches.Reset()
		//	flag, err := send("15816138010","","","")
		//	convey.So(err, convey.ShouldNotBeNil)
		//	convey.So(flag, convey.ShouldEqual, nil)
		//})
	})

}
