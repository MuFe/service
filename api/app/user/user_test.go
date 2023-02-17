package user

import (
	"bou.ke/monkey"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"mufe_service/camp/cache"
	"mufe_service/jsonRpc"
	"testing"
)

func TestSendCode(t *testing.T) {
	type query struct {
		Phone string `form:"phone"`
		Type  int64  `form:"type"`
	}
	params := query{
		Phone:"15816138010",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest(
		"POST",
		"/hello",
		bytes.NewBuffer(b),
	)

	c := &gin.Context{Request:req}


	//// 对 varys.GetInfoByUID 进行打桩
	//// 无论传入的uid是多少，都返回 &varys.UserInfo{Name: "liwenzhou"}, nil
	//monkey.PatchInstanceMethod(tSend, func(string,string,*gin.Context)(*app.SendResponse, error) {
	//	return &app.SendResponse{Result:"aaaaa"}, nil
	//})
	monkey.Patch(cache.GetUserPhoneCodeLimitTime, func(string,int64)(int64, error) {
		return 300, nil
	})
	monkey.Patch(cache.SetUserPhoneCode, func(string,string,int64)(int64, error) {
		return 300, nil
	})
	SendCode(c)

}
