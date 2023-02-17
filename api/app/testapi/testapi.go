package testapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"mufe_service/camp/cache"
	"mufe_service/camp/dataUtil"
	"mufe_service/camp/errcode"
	"mufe_service/camp/server"
	"mufe_service/camp/xlog"
	pb "mufe_service/jsonRpc"
	"mufe_service/manager"
)

func init() {
	if os.Getenv("MODEL")=="test"{
		server.Post("/appApi/test/testPush", testPush)
		server.Post("/appApi/test/pay", testPay)
		server.Post("/appApi/test/re", re)
		server.Post("/appApi/test/delToken", del)
	}

}

func testPush(c *gin.Context) {
	type query struct {
		Ids []string `form:"ids"`
		Type int64 `form:"type"`
		Id int64 `form:"id"`
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	content:=fmt.Sprintf("{type:\"%d\",id:\"%d\"}", params.Type, params.Id)
	msg:=fmt.Sprintf("{code:\"%d\",data:\"%s\",msg:\"%s\"}", 200, content, "")
	_,err:=manager.GetPushService().PushMessage(c,&pb.PushRequest{
		DeviceList:params.Ids,
		Message:msg,
	})
	if err!=nil{
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("推送成功"))
}

func testPay(c *gin.Context){
	type Param struct {
		Id   int64                 `form:"id" json:"id"`
		Type   int64                 `form:"type" json:"type"`
	}
	params := Param{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}

	result, err := manager.GetOrderService().Pay(c, &pb.OrderPayRequest{
		OrderId: params.Id,
		AppId:       os.Getenv("APP_WX_ID"),
		NotifyUrl: os.Getenv("API") + "/appApi/pay/notify",
		PayTotalFree:1,
		ChannelType: 1,
		PayType:   params.Type,
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if result.Pay {
		c.AbortWithStatusJSON(http.StatusNotModified, errcode.ParseOK("已经支付过"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(dataUtil.ParseWxPay(result)))
	}
}


func re(c *gin.Context){
	bodyData, err := ioutil.ReadAll(c.Request.Body)
	if err==nil{
		xlog.Info(string(bodyData))
	}
}

func del(c *gin.Context){
	err:=cache.DeleteUserToken(1189)
	if err!=nil{
		xlog.ErrorP(err)
	}
}

