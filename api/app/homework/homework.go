package homework

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"mufe_service/camp/dataUtil"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/handler"
	"mufe_service/camp/jwt"
	"mufe_service/camp/server"
	pb "mufe_service/jsonRpc"
	"mufe_service/manager"
)

func init() {
	server.Post("/appApi/homework/finish", handler.UserLogin, finish)
	server.Post("/appApi/homework/record", handler.UserLogin, record)

	server.Post("/appApi/homework/homework", handler.UserLogin,homework)
	server.Post("/appApi/homework/detail", handler.UserLogin, detail)
	server.Post("/appApi/homework/detail/v2", handler.UserLogin, detailV2)
	server.Post("/appApi/my/homework", handler.UserLogin, myHomework)
	server.Post("/appApi/homework/addHomeWork", handler.UserLogin, addHomework)
}

type HomeWork struct {
	Id       int64   `json:"id" `
	Time     int64   `json:"time" `
	Title    string  `json:"title" `
	Index    int64   `json:"index" `
	Number   int64   `json:"number" `
	Info     int64   `json:"info" `
	Cover    string  `json:"cover" `
	Progress float64 `json:"progress" `
	Finish   bool    `json:"finish" `
}

func homework(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		ClassId int64 `form:"class_id" json:"class_id" `
		Time    int64 `form:"time" json:"time" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Time == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	temp := int64(params.Time / 86400)
	params.Time = 86400 * temp
	rpcResult, err := manager.GetHomeWorkService().HomeWorkList(c, &pb.HomeWorkListRequest{Time: params.Time, ClassId: params.ClassId, Uid: userData.Uid})
	type Result struct {
		Id       int64   `json:"id"`
		Title    string  `json:"title"`
		Cover    string  `json:"cover"`
		Index    int64   `json:"index"`
		Number   int64   `json:"number"`
		Progress float64 `json:"progress"`
		Time     int64   `json:"time"`
		Finish   bool    `json:"finish" `
	}
	result := make([]Result, 0)
	if err == nil {
		for _, v := range rpcResult.List {
			result = append(result, Result{
				Id:       v.Id,
				Time:     v.Time,
				Title:    v.Title,
				Index:    v.Index,
				Progress: v.Progress,
				Number:   v.Number,
				Cover:    v.Cover,
				Finish:   v.Progress == 100,
			})
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func myHomework(c *gin.Context) {
	type query struct {
		Time int64 `form:"time" json:"time" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Time == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	rpcResult, err := manager.GetHomeWorkService().HomeWorkGroup(c, &pb.HomeWorkListRequest{Uid: userData.Uid, Time: params.Time})
	type Result struct {
		Desc string     `json:"desc"`
		Id   int64      `json:"id"`
		List []HomeWork `json:"list"`
	}
	list := make([]HomeWork, 0)
	result := Result{

	}
	if err == nil {
		for _, v := range rpcResult.List {
			list = append(list, HomeWork{
				Id:       v.Id,
				Time:     v.Time,
				Title:    v.Title,
				Index:    v.Index,
				Progress: v.Progress,
				Number:   v.Number,
				Cover:    v.Cover,
				Info:     v.Info,
				Finish:   v.Progress == 100,
			})
		}
		result.List = list
		result.Desc = rpcResult.Desc
		result.Id = rpcResult.GroupId
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func record(c *gin.Context) {
	type query struct {
		Id int64 `form:"id" json:"id" `
	}
	type HomeWorkRecord struct {
		Id   int64  `json:"id" `
		Name string `json:"name" `
		Head string `json:"head" `
		Type string `json:"type" `
	}
	type Result struct {
		Finish     []HomeWorkRecord `json:"finish" `
		InComplete []HomeWorkRecord `json:"incomplete" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Id == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	result, err := manager.GetHomeWorkService().HomeWorkRecord(c, &pb.HomeWorkRecordRequest{Id: params.Id})
	if err == nil {
		returnResult := Result{
			Finish:     make([]HomeWorkRecord, 0),
			InComplete: make([]HomeWorkRecord, 0),
		}
		uidList := make([]int64, 0)
		for _, v := range result.Incomplete {
			uidList = append(uidList, v.Uid)
		}
		for _, v := range result.Finish {
			uidList = append(uidList, v.Uid)
		}
		uResultMap := make(map[int64]HomeWorkRecord)
		uResult, err := manager.GetUserService().GetUserList(c, &pb.GetUserListRequest{IdList: uidList,Status:enum.StatusNormal})
		if err == nil {
			for _, v := range uResult.List {
				uResultMap[v.Uid] = HomeWorkRecord{
					Id:   v.Uid,
					Name: v.Name,
					Head: v.Head,
					Type: "学生",
				}
			}
			for _, v := range result.Incomplete {
				temp, ok := uResultMap[v.Uid]
				if ok {
					returnResult.InComplete = append(returnResult.InComplete, temp)
				}
			}
			for _, v := range result.Finish {
				temp, ok := uResultMap[v.Uid]
				if ok {
					returnResult.Finish = append(returnResult.Finish, temp)
				}
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(returnResult))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func finish(c *gin.Context) {
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
	_, err := manager.GetHomeWorkService().FinishHomeWork(c, &pb.FinishHomeWorkRequest{Id: params.Id, Uid: userData.Uid, Score: 100})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("提交成功"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func addHomework(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Ids      []int64 `form:"ids" json:"ids" `
		ClassIds []int64 `form:"class_ids" json:"class_ids" `
		Time     int64   `form:"time" json:"time" `
		Desc     string  `form:"desc" json:"desc" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if len(params.Ids) == 0 || len(params.ClassIds) == 0 || params.Time == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	temp := int64(params.Time / 86400)
	params.Time = 86400 * temp
	_, err := manager.GetHomeWorkService().AddHomeWork(c, &pb.AddHomeWorkRequest{Ids: params.Ids, ClassIds: params.ClassIds, Time: params.Time, Desc: params.Desc, Uid: userData.Uid})
	if err == nil {
		go func() {
			ids := make([]int64, 0)
			for _, id := range params.Ids {
				v, err := manager.GetSchoolService().ClassDetail(c, &pb.ClassDetailRequest{Id: id, Uid: userData.Uid})
				if err == nil {
					ids = append(ids, v.AdminList...)
				}
			}
			userResult, err := manager.GetUserService().GetUserList(c, &pb.GetUserListRequest{IdList: ids})
			if err == nil {
				deviceList := make([]string, 0)
				content := "作业已发布，记得及时完成哦"
				for _, v := range userResult.List {
					if v.Uid != userData.Uid {
						deviceList = append(deviceList, v.RegistrationId)
					}

				}
				manager.GetPushService().PushMessage(c, &pb.PushRequest{DeviceList: deviceList, Content: content})
			}
		}()

		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("发布成功"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func detail(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Id      int64 `form:"id" json:"id" `
		Detail  int64 `form:"detail_id" json:"detail_id" `
		GroupId int64 `form:"group" json:"group" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	detailResult, err := manager.GetHomeWorkService().HomeWorkDetail(c, &pb.HomeWorkDetailRequest{InfoId: params.Id})
	if err == nil {
		type HomeWork struct {
			Title string `json:"title"`
			Id    int64  `json:"id"`
			Info  int64  `json:"info"`
			Cover string `json:"cover"`
		}
		type Result struct {
			Title   string         `json:"title"`
			Cover   string         `json:"cover"`
			Video   string         `json:"video"`
			Content string         `json:"content"`
			Id      int64          `json:"id"`
			Tag     []dataUtil.Tag `json:"tag"`
			List    []HomeWork     `json:"list"`
		}
		result := Result{
			Id:      detailResult.Id,
			Title:   detailResult.Title,
			Cover:   detailResult.Cover,
			Video:   detailResult.Video,
			Content: detailResult.Content,
			Tag:     dataUtil.ParseTag(detailResult.Tag),
			List:    make([]HomeWork, 0),
		}
		rpcResult, err := manager.GetHomeWorkService().HomeWorkGroup(c, &pb.HomeWorkListRequest{Id: params.GroupId, Uid: userData.Uid})
		if err == nil {
			for _, v := range rpcResult.List {
				if v.Id != params.Detail {
					result.List = append(result.List, HomeWork{
						Id:    v.Id,
						Title: v.Title,
						Info:  v.Info,
						Cover: v.Cover,
					})
				}
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))

	}
}

func detailV2(c *gin.Context) {
	type query struct {
		Id     int64 `form:"id" json:"id" `
		Detail int64 `form:"detail_id" json:"detail_id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	detailResult, err := manager.GetHomeWorkService().HomeWorkDetail(c, &pb.HomeWorkDetailRequest{Id: params.Id})
	if err == nil {
		type HomeWork struct {
			Title string `json:"title"`
			Id    int64  `json:"id"`
			Cover string `json:"cover"`
		}
		type Result struct {
			Title   string         `json:"title"`
			Cover   string         `json:"cover"`
			Video   string         `json:"video"`
			Content string         `json:"content"`
			Id      int64          `json:"id"`
			Tag     []dataUtil.Tag `json:"tag"`
			List    []HomeWork     `json:"list"`
		}
		result := Result{
			Id:      detailResult.Id,
			Title:   detailResult.Title,
			Cover:   detailResult.Cover,
			Video:   detailResult.Video,
			Content: detailResult.Content,
			Tag:     dataUtil.ParseTag(detailResult.Tag),
			List:    make([]HomeWork, 0),
		}

		rpcResult, err := manager.GetHomeWorkService().HomeWork(c, &pb.HomeWorkRequest{ContentId: []int64{detailResult.ContentId}, Status: enum.StatusNormal})
		if err == nil {
			for _, v := range rpcResult.List {
				if v.Id != params.Id {
					result.List = append(result.List, HomeWork{
						Id:    v.Id,
						Title: v.Title,
						Cover: v.Cover,
					})
				}
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))

	}
}
