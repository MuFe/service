package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/handler"
	"mufe_service/camp/server"
	"mufe_service/camp/utils"
	"mufe_service/jsonRpc"
	"mufe_service/manager"
)

func init() {
	server.Post("/adminUser/getUserList", handler.AdminLogin, list)
	server.Post("/adminUser/enterCancel", handler.AdminLogin, cancel)

}

func list(c *gin.Context) {
	type query struct {
		Page   int64 `form:"page" json:"page" `
		Size   int64 `form:"size" json:"size" `
		Status int64 `form:"status" json:"status" `
		Cancel int64 `form:"cancel" json:"cancel" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	typeInt:=enum.CANCEL_DEFUALT_TYPE
	if params.Cancel==1{
		typeInt=enum.CANCEL_START_TYPE
	}
	if params.Status==3{
		typeInt=enum.CANCEL_END_TYPE
		params.Status=enum.StatusDelete
	}
	result, err := manager.GetUserService().GetUserList(c, &app.GetUserListRequest{Size: params.Size, Page: params.Page, Status: params.Status,CancelType:typeInt})
	if err == nil {
		type School struct {
			Id         int64  `json:"id" `
			Head       string `json:"head" `
			Phone      string `json:"phone" `
			Name       string `json:"name" `
			Identity   int64  `json:"identity" `
			CancelTime int64  `json:"cancel_time" `
			LoginType  int64  `json:"login_type" `
			LoginTime  int64  `json:"login_time" `
		}
		list := make([]interface{}, 0)
		for _, v := range result.List {
			list = append(list, School{
				Id:         v.Uid,
				Name:       v.Name,
				Phone:      v.Phone,
				Head:       v.Head,
				Identity:   v.Identity,
				LoginType:  v.LoginType,
				CancelTime: v.CancelTime,
				LoginTime:v.LoginTime,
			})
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(utils.CreateListResultReturn(result.Total, list)))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func cancel(c *gin.Context) {
	type query struct {
		Id []int64 `form:"ids" json:"ids" `
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if len(params.Id) == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	_, err := manager.GetUserService().EnterCancel(c, &app.EnterCancelRequest{List: params.Id})
	if err == nil {
		manager.GetSchoolService().CancelSchool(c,&app.CancelSchoolRequest{List:params.Id})
		manager.GetHomeWorkService().CancelHomeWork(c,&app.CancelHomeWorkRequest{List:params.Id})
		manager.GetCoachService().CancelCoach(c,&app.CancelCoachRequest{List:params.Id})
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("操作成功"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}
