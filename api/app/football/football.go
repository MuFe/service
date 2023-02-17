package football

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/handler"
	"mufe_service/camp/jwt"
	"mufe_service/camp/server"
	pb "mufe_service/jsonRpc"
	"mufe_service/manager"
	"time"
)

func init() {
	server.Post("/appApi/football/data", getFootBallData)
	server.Post("/appApi/football/getFootball", handler.UserLogin, getFoot)
	server.Post("/appApi/football/bind", handler.UserLogin, bind)
	server.Post("/appApi/football/unbind", handler.UserLogin, unbind)
	server.Post("/appApi/football/teacher", handler.UserLogin, teacher)


}

func getFootBallData(c *gin.Context) {
	type Result struct {
		Score    int64   `json:"score"`
		Duration float64 `json:"duration"`
		Time     int64   `json:"time"`
	}
	type query struct {
		Time int64  `form:"time" json:"time" `
		Mac  string `form:"mac" json:"mac" `
	}

	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	requestList := make([]*pb.GetFootData, 0)
	requestList = append(requestList, &pb.GetFootData{Mac: params.Mac, Date: time.Unix(params.Time, 0).Format(enum.TimeFormatDate)})
	result, err := manager.GetFootBallService().GetData(c, &pb.GetFootDataRequest{List: requestList})
	list := make([]Result, 0)
	if err == nil {
		for _, v := range result.List {
			list = append(list, Result{
				Score:    v.Score,
				Duration: v.Duration,
				Time:     params.Time,
			})
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
}

func getFoot(c *gin.Context) {
	type Result struct {
		Mac string `json:"mac"`
	}
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	result, err := manager.GetFootBallService().GetFoot(c, &pb.GetFootRequest{Uid: []int64{userData.Uid}})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	list := make([]Result, 0)
	for _, v := range result.List {
		if v.Uid == userData.Uid {
			list = append(list, Result{
				Mac: v.Mac,
			})
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
}

func bind(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Mac string `form:"mac" json:"mac" `
	}

	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetFootBallService().Bind(c, &pb.BindFootRequest{Mac: params.Mac, Uid: userData.Uid})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("绑定成功"))
}

func unbind(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Mac string `form:"mac" json:"mac" `
	}

	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetFootBallService().Unbind(c, &pb.BindFootRequest{Mac: params.Mac, Uid: userData.Uid})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("解除绑定成功"))
}

func teacher(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Time    int64 `form:"time" json:"time" `
		ClassId int64 `form:"class" json:"class" `
	}
	type Result struct {
		Score    int64   `json:"score"`
		Duration float64 `json:"duration"`
		Time     int64   `json:"time"`
		Uid      int64   `json:"uid"`
		Name     string  `json:"name"`
		Head     string  `json:"head"`
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	list := make([]Result, 0)
	v, err := manager.GetSchoolService().ClassDetail(c, &pb.ClassDetailRequest{Id: params.ClassId, Uid: userData.Uid})
	if err == nil {
		uidList := make([]int64, 0)
		for _, value := range v.StudentList {
			uidList = append(uidList, value)
		}
		reUidList:=make([]int64,0)
		userResult, err := manager.GetUserService().GetUserList(c, &pb.GetUserListRequest{IdList: uidList, Status: enum.StatusNormal})
		userMap := make(map[int64]*pb.UserDataResponse)
		if err == nil {
			for _, value := range userResult.List {
				userMap[value.Uid] = value
				reUidList=append(reUidList,value.Uid)
			}

		}
		footResult, err := manager.GetFootBallService().GetFoot(c, &pb.GetFootRequest{Uid: reUidList})
		requestList := make([]*pb.GetFootData, 0)
		macMap:=make(map[string]bool,0)
		if err == nil {
			for _, v := range footResult.List {
				_,ok:=macMap[v.Mac]
				if !ok{
					requestList = append(requestList, &pb.GetFootData{Mac: v.Mac, Date: time.Unix(params.Time, 0).Format(enum.TimeFormatDate)})
					macMap[v.Mac]=true
				}
			}
		}

		footDataResult, err := manager.GetFootBallService().GetData(c, &pb.GetFootDataRequest{List: requestList})
		if err == nil {
			for _, v := range footDataResult.List {
				for _, vv := range footResult.List {
					if vv.Mac==v.Mac{
						uResult, ok := userMap[vv.Uid]
						if ok {
							list = append(list, Result{
								Score:    v.Score,
								Duration: v.Duration,
								Time:     params.Time,
								Uid:      uResult.Uid,
								Name:     uResult.Name,
								Head:     uResult.Head,
							})
						}
					}
				}
			}
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))

}
