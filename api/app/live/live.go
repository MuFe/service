package football

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"mufe_service/camp/dataUtil"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/handler"
	"mufe_service/camp/jwt"
	"mufe_service/camp/server"
	"mufe_service/camp/xlog"
	pb "mufe_service/jsonRpc"
	"mufe_service/manager"
)

func init() {
	server.Post("/appApi/live/list", handler.UserLogin, list)
	server.Post("/appApi/live/create", handler.UserLogin, create)
	server.Post("/appApi/live/start", handler.UserLogin, start)
	server.Post("/appApi/live/pay", handler.UserLogin, pay)
	server.Post("/appApi/live/end", handler.UserLogin, end)
	server.Post("/appApi/live/state", handler.UserLogin, state)
	server.Post("/appApi/live/updateScore", handler.UserLogin, updateScore)
	server.Post("/appApi/live/endLive", handler.UserLogin, endLive)
	server.Post("/appApi/live/watch", handler.UserLogin, watch)
	server.Post("/appApi/live/endWatch", handler.UserLogin, endWatch)
	server.Post("/appApi/live/package", handler.UserLogin, packageList)
	server.Post("/appApi/live/teamMember", handler.UserLogin, teamMember)
	server.Post("/appApi/live/liveData", handler.UserLogin, liveData)
	server.Post("/appApi/live/watchInfo", handler.UserLogin, watchData)

}

func list(c *gin.Context) {
	type Team struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
		Head string `json:"head"`
	}
	type Data struct {
		Score     string `json:"score"`
		Address   string `json:"address"`
		Id        int64  `json:"id"`
		StartTime int64  `json:"start_time"`
		EndTime   int64  `json:"end_time"`
		Status    int64  `json:"status"`
		Home      Team   `json:"home_team"`
		Visiting  Team   `json:"visiting_team"`
	}
	type Result struct {
		Type int64  `json:"type"`
		List []Data `json:"list"`
	}

	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	result, err := manager.GetLiveService().GetList(c, &pb.GetLiveDataRequest{Uid: userData.Uid})
	list := make([]Result, 0)
	if err == nil {
		for _, v := range result.List {
			tempList := make([]Data, 0)
			for _, vv := range v.List {
				home := Team{
					Id:   vv.Home.Id,
					Name: vv.Home.Name,
					Head: vv.Home.Head,
				}
				visiting := Team{
					Id:   vv.Visiting.Id,
					Name: vv.Visiting.Name,
					Head: vv.Visiting.Head,
				}
				tempList = append(tempList, Data{
					Home:      home,
					Visiting:  visiting,
					Score:     vv.HomeScore + "-" + vv.VisitingScore,
					Address:   vv.Address,
					Id:        vv.Id,
					StartTime: vv.Start,
					EndTime:   vv.End,
					Status:    vv.Status,
				})
			}

			list = append(list, Result{
				Type: v.Type,
				List: tempList,
			})
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
}

func create(c *gin.Context) {
	type team struct {
		Name   string `form:"name" json:"name" `
		Number string `form:"number" json:"number" `
		Uid    int64  `form:"uid" json:"uid" `
	}
	type query struct {
		Type           int64  `form:"type" json:"type" `
		Address        string `form:"address" json:"address" `
		Home           string `form:"home_team" json:"home_team" `
		HomeNumber     []team `form:"home_member" json:"home_member" `
		Visiting       string `form:"visiting_team" json:"visiting_team" `
		VisitingNumber []team `form:"visiting_member" json:"visiting_member" `
		Time           int64  `form:"time" json:"time" `
		Package        int64  `form:"package" json:"package" `
		OrderId        int64  `form:"order_id" json:"order_id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	xlog.Info(params)
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	if params.Package != 0 && params.OrderId == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	if params.Package != 0 && params.OrderId != 0 {
		_, err := manager.GetOrderService().EditStatus(c, &pb.EditStatusRequest{
			OrderId: params.OrderId,
			Status:  enum.OrderStatusPay,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
			return
		}
	}
	hNumber := make([]*pb.TeamMemberData, 0)
	for _, v := range params.HomeNumber {
		hNumber = append(hNumber, &pb.TeamMemberData{
			Number: v.Number,
			Name:   v.Name,
			Uid:    v.Uid,
		})
	}
	vNumber := make([]*pb.TeamMemberData, 0)
	for _, v := range params.VisitingNumber {
		vNumber = append(vNumber, &pb.TeamMemberData{
			Number: v.Number,
			Name:   v.Name,
			Uid:    v.Uid,
		})
	}
	result, err := manager.GetLiveService().Create(c, &pb.CreateLiveRequest{Type: params.Type, Uid: userData.Uid, Home: params.Home, Visiting: params.Visiting, HomeMember: hNumber, VisitingMember: vNumber, PackageId: params.Package, Address: params.Address})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	} else {
		type Result struct {
			Id   int64  `json:"id"`
			Pass string `json:"pass"`
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(Result{
			Id:   result.Id,
			Pass: result.Pass,
		}))
	}
}

func start(c *gin.Context) {
	type Result struct {
		Address string `form:"address" json:"address" `
		Time    int64  `form:"time" json:"time" `
		Key     string `form:"key" json:"key" `
		Url     string `form:"url" json:"url" `
	}
	type query struct {
		Pass string `form:"pass" json:"pass" `
	}

	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	result, err := manager.GetLiveService().Start(c, &pb.StartLivRequest{Uid: userData.Uid, Pass: params.Pass})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(Result{
			Address: result.Address,
			Time:    result.End,
			Key:     "85704abc3be90dfb746405979835c78a",
			Url:     "https://license.vod2.myqcloud.com/license/v2/1312887211_1/v_cube.license",
		}))
	}
}

func updateScore(c *gin.Context) {
	type Score struct {
		Id   int64  `form:"id" json:"id" `
		Time string `form:"time" json:"time" `
	}
	type query struct {
		Id            int64   `form:"id" json:"id" `
		HomeScore     string  `form:"home_score" json:"home_score" `
		VisitingScore string  `form:"visiting_score" json:"visiting_score" `
		Score         []Score `form:"score_info" json:"score_info" `
	}

	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	temp := make([]*pb.ScoreData, 0)
	for _, v := range params.Score {
		temp = append(temp, &pb.ScoreData{
			Id:   v.Id,
			Time: v.Time,
		})
	}

	_, err := manager.GetLiveService().UpdateScore(c, &pb.UpdateScoreRequest{Id: params.Id, HomeScore: params.HomeScore, VisitingScore: params.VisitingScore, List: temp})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	} else {
		result, err := manager.GetLiveService().GetWatchInfo(c, &pb.LiveRequest{Id: params.Id})
		if err == nil {
			go func() {
				userResult, err := manager.GetUserService().GetUserList(c, &pb.GetUserListRequest{IdList: result.User})
				if err == nil {
					deviceList := make([]string, 0)
					content := fmt.Sprintf("{type:\"%d\",id:\"%d\"}", 1, params.Id)
					msg := fmt.Sprintf("{code:\"%d\",data:\"%s\",msg:\"%s\"}", 200, content, "")
					for _, v := range userResult.List {
						deviceList = append(deviceList, v.RegistrationId)
					}
					manager.GetPushService().PushMessage(c, &pb.PushRequest{DeviceList: deviceList, Message: msg})
				}
			}()
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
	}
}

func end(c *gin.Context) {
	type query struct {
		Id int64 `form:"id" json:"id" `
	}
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}

	_, err := manager.GetLiveService().End(c, &pb.EndLivRequest{Id: params.Id, Uid: userData.Uid})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
	}
}

func endLive(c *gin.Context) {

	type query struct {
		Id int64 `form:"id" json:"id" `
	}

	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	_, err := manager.GetLiveService().EndLive(c, &pb.StartLivRequest{Id: params.Id, Uid: userData.Uid})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	} else {
		result, err := manager.GetLiveService().GetWatchInfo(c, &pb.LiveRequest{Id: params.Id})
		if err == nil {
			go func() {
				userResult, err := manager.GetUserService().GetUserList(c, &pb.GetUserListRequest{IdList: result.User})
				if err == nil {
					deviceList := make([]string, 0)
					content := fmt.Sprintf("{type:\"%d\",id:\"%d\"}", 3, params.Id)
					msg := fmt.Sprintf("{code:\"%d\",data:\"%s\",msg:\"%s\"}", 200, content, "")
					for _, v := range userResult.List {
						deviceList = append(deviceList, v.RegistrationId)
					}
					manager.GetPushService().PushMessage(c, &pb.PushRequest{DeviceList: deviceList, Message: msg})
				}
			}()
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
	}
}

func state(c *gin.Context) {
	type query struct {
		Id      int64 `form:"id" json:"id" `
		Suspend bool  `form:"suspend" json:"suspend" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	result, err := manager.GetLiveService().GetWatchInfo(c, &pb.LiveRequest{Id: params.Id})
	if err == nil {
		go func() {
			userResult, err := manager.GetUserService().GetUserList(c, &pb.GetUserListRequest{IdList: result.User})
			if err == nil {
				deviceList := make([]string, 0)
				state := 1
				if params.Suspend {
					state = 0
				}
				content := fmt.Sprintf("{type:\"%d\",id:\"%d\",state:\"%d\"}", 2, params.Id, state)
				msg := fmt.Sprintf("{code:\"%d\",data:\"%s\",msg:\"%s\"}", 200, content, "")
				for _, v := range userResult.List {
					deviceList = append(deviceList, v.RegistrationId)
				}
				manager.GetPushService().PushMessage(c, &pb.PushRequest{DeviceList: deviceList, Message: msg})
			}
		}()
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
}

func watch(c *gin.Context) {

	type Result struct {
		Address string `form:"address" json:"address" `
		Id      int64  `form:"id" json:"id" `
	}
	type query struct {
		Pass string `form:"pass" json:"pass" `
	}

	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	result, err := manager.GetLiveService().Watch(c, &pb.WatchLiveRequest{Uid: userData.Uid, Pass: params.Pass})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(Result{
			Address: result.Address,
			Id:      result.Id,
		}))
	}
}

func endWatch(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	_, err := manager.GetLiveService().EndWatch(c, &pb.WatchLiveRequest{Uid: userData.Uid})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
	}
}

func packageList(c *gin.Context) {
	type Result struct {
		Id       int64  `json:"id"`
		Name     string `json:"name"`
		Max      int64  `json:"max"`
		Duration int64  `json:"duration"`
		Price    int64  `json:"price"`
	}
	result, err := manager.GetLiveService().Package(c, &pb.EmptyRequest{})
	list := make([]Result, 0)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	} else {
		for _, v := range result.List {
			list = append(list, Result{
				Id:       v.Id,
				Max:      v.Max,
				Name:     v.Name,
				Duration: v.Duration,
				Price:    v.Price,
			})
		}

		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	}
}

func pay(c *gin.Context) {
	type query struct {
		Id      int64 `form:"id" json:"id" `
		Channel int64 `form:"channel" json:"channel" `
	}

	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	list := make([]*pb.OrderData, 0)
	data := &pb.OrderData{
		SellerId: 0,
	}
	result, err := manager.GetLiveService().Package(c, &pb.EmptyRequest{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	title := ""
	desc := ""
	for _, v := range result.List {
		if v.Id == params.Id {
			data.Price = v.Price
			title = v.Name
			desc = v.Name
			temp := make([]*pb.OrderGood, 0)
			temp = append(temp, &pb.OrderGood{
				Num:   1,
				SkuId: v.SkuId,
				Price: v.Price,
			})
			data.List = temp

			list = append(list, data)
		}
	}
	if len(list) == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	orderResult, err := manager.GetOrderService().Create(c, &pb.OrderRequest{
		BuyerId:       userData.Uid,
		Response:      list,
		Time:          int64(30 * 60),
		Title:         title,
		Desc:          desc,
		Price:         data.Price,
		OrderType:     enum.OrderTypeApp,
		Phone:         "",
		Province:      "",
		City:          "",
		Area:          "",
		Address:       "",
		Consignee:     "",
		ShowOrderList: enum.NoShowOrderList,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	payResult, err := manager.GetOrderService().Pay(c, &pb.OrderPayRequest{
		OrderId:     orderResult.Id,
		AppId:       os.Getenv("APP_WX_ID"),
		NotifyUrl:   os.Getenv("API") + "/pay/notify",
		ChannelType: params.Channel,
		PayType:     enum.AppPayType,
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(dataUtil.ParseWxPay(payResult)))
}

func teamMember(c *gin.Context) {
	type query struct {
		Id int64 `form:"id" json:"id" `
	}
	type team struct {
		Name   string `json:"name" `
		Number string ` json:"number" `
		Head   string ` json:"head" `
		Id     int64  `json:"id" `
	}
	type Result struct {
		HomeNumber     []team `form:"home_member" json:"home_member" `
		VisitingNumber []team `form:"visiting_member" json:"visiting_member" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}

	result, err := manager.GetLiveService().TeamMember(c, &pb.TeamMemberRequest{Id: params.Id})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	} else {
		re := Result{
			VisitingNumber: make([]team, 0),
			HomeNumber:     make([]team, 0),
		}
		for _, v := range result.VisitingMember {
			re.VisitingNumber = append(re.VisitingNumber, team{
				Id:     v.Id,
				Name:   v.Name,
				Number: v.Number,
				Head:   v.Head,
			})
		}
		for _, v := range result.HomeMember {
			re.HomeNumber = append(re.HomeNumber, team{
				Id:     v.Id,
				Name:   v.Name,
				Number: v.Number,
				Head:   v.Head,
			})
		}

		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(re))
	}
}

func liveData(c *gin.Context) {
	type query struct {
		Id int64 `form:"id" json:"id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	type Team struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
		Head string `json:"head"`
	}
	type ScoreInfo struct {
		Name string `json:"name"`
		Head string `json:"head"`
		Time string `json:"time"`
		Home bool   `json:"home"`
	}
	type Data struct {
		HomeScore     string      `json:"home_score"`
		VisitingScore string      `json:"visiting_score"`
		Address       string      `json:"address"`
		Id            int64       `json:"id"`
		Home          Team        `json:"home_team"`
		Visiting      Team        `json:"visiting_team"`
		Info          []ScoreInfo `json:"info"`
	}
	result, err := manager.GetLiveService().LiveInfo(c, &pb.LiveRequest{Id: params.Id})
	if err == nil {
		home := Team{
			Id:   result.Home.Id,
			Name: result.Home.Name,
			Head: result.Home.Head,
		}
		visiting := Team{
			Id:   result.Visiting.Id,
			Name: result.Visiting.Name,
			Head: result.Visiting.Head,
		}
		info := make([]ScoreInfo, 0)
		for _, v := range result.Info {
			info = append(info, ScoreInfo{
				Name: v.Data.Name,
				Head: v.Data.Head,
				Time: v.Time,
				Home: v.IsHome,
			})
		}
		re := Data{
			Home:          home,
			Visiting:      visiting,
			HomeScore:     result.HomeScore,
			VisitingScore: result.VisitingScore,
			Address:       result.Address,
			Id:            result.Id,
			Info:          info,
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(re))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}

}

func watchData(c *gin.Context) {
	type query struct {
		Id int64 `form:"id" json:"id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	type Data struct {
		Watch int64 `json:"watch"`
	}
	result, err := manager.GetLiveService().GetWatchInfo(c, &pb.LiveRequest{Id: params.Id})
	if err == nil {
		re := Data{
			Watch: int64(len(result.User)),
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(re))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}

}
