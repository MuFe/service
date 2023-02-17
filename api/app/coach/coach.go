package course

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"mufe_service/camp/cache"
	"mufe_service/camp/dataUtil"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/handler"
	"mufe_service/camp/jwt"
	"mufe_service/camp/server"
	app "mufe_service/jsonRpc"
	"mufe_service/manager"
	"time"
)

func init() {
	server.Post("/appApi/coach/joinInstitution", handler.UserLogin, handler.DefaultCheck, addInstitution)
	server.Post("/appApi/coach/editIntroduce", handler.UserLogin, handler.CoachCheck, editIntroduce)
	server.Post("/appApi/coach/quitInstitution", handler.UserLogin, handler.CoachCheck, quitInstitution)
	server.Get("/appApi/coach/institution", handler.UserLogin, handler.CoachCheck,institution)
	server.Post("/appApi/coach/addWork", handler.UserLogin, handler.CoachCheck, addWork)
	server.Get("/appApi/coach/work", handler.UserLogin, handler.CoachCheck, workList)
	server.Get("/appApi/coach/relist", handler.UserLogin, reList)
	server.Get("/appApi/coach/course", handler.UserLogin, course)
	server.Get("/appApi/coach/detail", handler.UserLogin, coachDetail)
	server.Get("/appApi/coach/list", handler.UserLogin, coachList)
	server.Get("/appApi/coach/place", handler.UserLogin, placeList)
	server.Get("/appApi/coach/order", handler.UserLogin, orderList)
	server.Get("/appApi/coach/orderDetail", handler.UserLogin, orderDetail)
	server.Post("/appApi/coach/pay", handler.UserLogin, pay)
	server.Post("/appApi/coach/paySu", handler.UserLogin, paySu)
	server.Post("/appApi/coach/receive", handler.UserLogin, receive)
	server.Post("/appApi/coach/refund", handler.UserLogin, refund)
	server.Get("/appApi/coach/myPlace", handler.UserLogin, myPlaceLis)
	server.Get("/appApi/coach/courseList", handler.UserLogin, courseList)
}

func addInstitution(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Id   int64  `form:"id" json:"id" `
		Info string `form:"info" json:"info" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetCoachService().JoinInstitution(c, &app.JoinInstitutionCourseRequest{Uid: userData.Uid, InstitutionId: params.Id, Info: params.Info})
	if err == nil {
		result, err := manager.GetUserService().UpdateUser(c, &app.UpdateUserRequest{Uid: userData.Uid, IdentityType: enum.INSTITUTION_TYPE, Type: enum.UpdateUserIdentity})
		if err == nil {
			token, err := jwt.GenerateUserJwt(userData.Uid, enum.INSTITUTION_TYPE, userData.OpenId)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
				return
			}
			u := dataUtil.ParseUserCache(result)
			_ = cache.SetUserInfo(result.Uid, u)
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(token))
		} else {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		}
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func editIntroduce(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Info string `form:"info" json:"info" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	tempResult,err:=manager.GetCoachService().CoachList(c,&app.CoachListRequest{Id:userData.Uid})
	if err!=nil{
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}

	_, err = manager.GetCoachService().JoinInstitution(c, &app.JoinInstitutionCourseRequest{Uid: userData.Uid, Id: tempResult.List[0].Id, Info: params.Info})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func quitInstitution(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Id   int64  `form:"id" json:"id" `
		Info string `form:"info" json:"info" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetCoachService().JoinInstitution(c, &app.JoinInstitutionCourseRequest{Uid: userData.Uid, Id: params.Id, Quit: true})
	if err == nil {
		result, err := manager.GetUserService().UpdateUser(c, &app.UpdateUserRequest{Uid: userData.Uid, IdentityType: enum.DEFAULT_TYPE, Type: enum.UpdateUserIdentity})
		if err == nil {
			token, err := jwt.GenerateUserJwt(userData.Uid, enum.DEFAULT_TYPE, userData.OpenId)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
				return
			}
			u := dataUtil.ParseUserCache(result)
			_ = cache.SetUserInfo(result.Uid, u)
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(token))
		} else {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		}
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func institution(c *gin.Context){
	type Result struct {
		Name string `json:"name" `
		Id   int64  `json:"id" `
		Address string `json:"info" `
		Icon string `json:"icon" `
	}
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	result,err:=manager.GetCoachService().GetInstitution(c,&app.InstitutionRequest{Uid:userData.Uid,Status:enum.StatusNormal})
	if err!=nil{
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}else{
		if len(result.List)==0{
			c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
			return
		}

		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(Result{
			Name:result.List[0].Name,
			Id:result.List[0].Id,
			Address:result.List[0].Address,
			Icon:result.List[0].Icon,
		}))
	}
}

func coachList(c *gin.Context) {
	type Result struct {
		Name string `json:"name" `
		Id   int64  `json:"id" `
		Info string `json:"info" `
		Head string `json:"head" `
	}

	type query struct {
		Page int64 `form:"page" json:"page" `
		Size int64 `form:"size" json:"size" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}

	chapterResult, err := manager.GetCoachService().CoachList(c, &app.CoachListRequest{Page: params.Page, Size: params.Size})
	if err == nil {
		list := make([]*Result, 0)
		listMap := make(map[int64]*Result)
		uidList := make([]int64, 0)
		for _, v := range chapterResult.List {
			temp := &Result{
				Info: v.Info,
				Id:   v.Uid,
			}
			listMap[v.Uid] = temp
			uidList = append(uidList, v.Uid)
			list = append(list, temp)
		}
		uResult, err := manager.GetUserService().GetUserList(c, &app.GetUserListRequest{IdList: uidList})
		if err == nil {
			for _, v := range uResult.List {
				temp, ok := listMap[v.Uid]
				if ok {
					temp.Name = v.Name
					temp.Head = v.Head
				}
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func placeList(c *gin.Context) {
	type Result struct {
		Name  string  `json:"name" `
		Id    int64   `json:"id" `
		Phone string  `json:"phone" `
		Place string  `json:"place" `
		Time  string  `json:"time" `
		Lat   float64 `json:"lat" `
		Lng   float64 `json:"lng" `
	}

	type query struct {
		Page int64 `form:"page" json:"page" `
		Size int64 `form:"size" json:"size" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}

	tResult, err := manager.GetCoachService().SchoolList(c, &app.InstitutionSchoolRequest{Page: params.Page, Size: params.Size, Status: enum.StatusNormal})
	if err == nil {
		tTemp := time.FixedZone("CST", 8*3600)
		list := make([]Result, 0)
		for _, v := range tResult.List {
			list = append(list, Result{
				Id:    v.Id,
				Name:  v.Name,
				Place: v.Address,
				Phone: v.Phone,
				Time:  time.Unix(v.Start, 0).In(tTemp).Format(enum.TimeFormatHour) + "点-" + time.Unix(v.End, 0).In(tTemp).Format(enum.TimeFormatHour) + "点",
				Lat:   0,
				Lng:   0,
			})
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func myPlaceLis(c *gin.Context) {
	type Result struct {
		Name  string  `json:"name" `
		Id    int64   `json:"id" `
		Phone string  `json:"phone" `
		Place string  `json:"place" `
		Time  string  `json:"time" `
		Lat   float64 `json:"lat" `
		Lng   float64 `json:"lng" `
	}

	type query struct {
		Page int64 `form:"page" json:"page" `
		Size int64 `form:"size" json:"size" `
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
	result, err := manager.GetCoachService().GetInstitution(c, &app.InstitutionRequest{Uid: userData.Uid, Status: enum.StatusNormal})
	if err == nil {
		if len(result.List) == 0 {
			c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
			return
		}
		list := make([]Result, 0)
		v := result.List[0]
		tResult, err := manager.GetCoachService().SchoolList(c, &app.InstitutionSchoolRequest{Iid: v.Id, Page: params.Page, Size: params.Size, Status: enum.StatusNormal})
		if err == nil {
			tTemp := time.FixedZone("CST", 8*3600)
			for _, v := range tResult.List {
				list = append(list, Result{
					Id:    v.Id,
					Name:  v.Name,
					Place: v.Address,
					Phone: v.Phone,
					Time:  time.Unix(v.Start, 0).In(tTemp).Format(enum.TimeFormatHour) + "点-" + time.Unix(v.End, 0).In(tTemp).Format(enum.TimeFormatHour) + "点",
					Lat:   0,
					Lng:   0,
				})
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func courseList(c *gin.Context) {
	type Result struct {
		Name  string `json:"name" `
		Id    int64  `json:"id" `
		Max   int64  `json:"max" `
		Price int64  `json:"price" `
		Level string `json:"level" `
	}

	type query struct {
		Page int64 `form:"page" json:"page" `
		Size int64 `form:"size" json:"size" `
		SchoolId int64 `form:"place" json:"place" `
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
	iId:=int64(0)
	if params.SchoolId==0{
		result, err := manager.GetCoachService().GetInstitution(c, &app.InstitutionRequest{Uid: userData.Uid, Status: enum.StatusNormal})
		if err!=nil{
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
			return
		}else{
			if len(result.List) == 0 {
				c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
				return
			}
			iId= result.List[0].Id
		}

	}
	list := make([]Result, 0)
	courseResult, err := manager.GetCoachService().CourseList(c, &app.InstitutionCourseRequest{Iid:iId,SchoolId:[]int64{params.SchoolId}, Page: params.Page, Size: params.Size, Status: enum.StatusAll})
	if err == nil {
		for _, v := range courseResult.List {
			list = append(list, Result{
				Name:  v.Name,
				Id:    v.Id,
				Price: v.Price,
				Max:   v.Max,
				Level: v.Level,
			})
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
}

func addWork(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Start  int64  `form:"start" json:"start" `
		End    int64  `form:"end" json:"end" `
		Place  int64  `form:"place" json:"place" `
		Course int64  `form:"course" json:"course" `
		Desc   string `form:"desc" json:"desc" `
		Id     int64  `form:"id" json:"id" `
		Del    bool   `form:"del" json:"del" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetCoachService().AddWork(c, &app.InstitutionWorkRequest{Uid: userData.Uid, Start: params.Start, End: params.End, Place: params.Place, Course: params.Course, Desc: params.Desc, Id: params.Id, Del: params.Del})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func workList(c *gin.Context) {
	type ResultData struct {
		Name    string `json:"name" `
		Id      int64  `json:"id" `
		Start   int64  `json:"start" `
		End     int64  `json:"end" `
		Now     int64  `json:"now" `
		Max     int64  `json:"max" `
		Place   int64  `json:"place" `
		Level   string `json:"level" `
		Reserve bool   `json:"reserve" `
	}

	type Result struct {
		Time int64        `json:"time" `
		List []ResultData `json:"list" `
	}

	type query struct {
		Place int64 `form:"place" json:"place" `
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
	start := (time.Now().Unix() / 86400) * 86400
	end := start + 8*86400
	result, err := manager.GetCoachService().WorkList(c, &app.InstitutionWorkListRequest{Uid: userData.Uid, Place: params.Place, Start: start, End: end})
	if err == nil {
		list := make([]*Result, 0)
		listMap := make(map[int64]*Result)
		for i := 0; i <= 6; i++ {
			tempTime := start + int64(i)*86400
			t := &Result{
				Time: tempTime,
				List: make([]ResultData, 0),
			}
			list = append(list, t)
			listMap[tempTime] = t
		}
		for _, v := range result.List {
			tempTime := (v.Start / 86400) * 86400
			t, ok := listMap[tempTime]
			if ok {
				t.List = append(t.List, ResultData{
					Place:   v.Place.Id,
					Id:      v.Id,
					Name:    v.Name,
					Start:   v.Start,
					End:     v.End,
					Now:     v.Now,
					Max:     v.Max,
					Level:   v.Level,
					Reserve: v.Reserve,
				})
			}

		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func reList(c *gin.Context) {
	type ResultData struct {
		Name    string `json:"name" `
		Head    string `json:"head" `
		Info    string `json:"info" `
		Id      int64  `json:"id" `
		Start   int64  `json:"start" `
		Level   string `json:"level" `
		End     int64  `json:"end" `
		Now     int64  `json:"now" `
		Max     int64  `json:"max" `
		Reserve bool   `json:"reserve" `
	}

	type Result struct {
		Time int64         `json:"time" `
		List []*ResultData `json:"list" `
	}

	type query struct {
		Place int64 `form:"place" json:"place" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	start := (time.Now().Unix() / 86400) * 86400
	end := start + 8*86400
	result, err := manager.GetCoachService().WorkList(c, &app.InstitutionWorkListRequest{Place: params.Place, Start: start, End: end})
	if err == nil {
		list := make([]*Result, 0)
		uidList := make([]int64, 0)
		listMap := make(map[int64]*Result)
		for i := 0; i <= 6; i++ {
			tempTime := start + int64(i)*86400
			t := &Result{
				Time: tempTime,
				List: make([]*ResultData, 0),
			}
			list = append(list, t)
			listMap[tempTime] = t
		}
		for _, v := range result.List {
			uidList = append(uidList, v.Uid)
			tempTime := (v.Start / 86400) * 86400
			t, ok := listMap[tempTime]
			if ok {
				tData := &ResultData{
					Id:      v.Uid,
					Start:   v.Start,
					End:     v.End,
					Now:     v.Now,
					Max:     v.Max,
					Info:    v.UserInfo,
					Level:   v.Level,
					Reserve: v.Reserve,
				}
				t.List = append(t.List, tData)
			}
		}
		uResult, err := manager.GetUserService().GetUserList(c, &app.GetUserListRequest{IdList: uidList})
		if err == nil {
			uMap:=make(map[int64]*app.UserDataResponse)
			for _, v := range uResult.List {
				uMap[v.Uid]=v

			}
			for _,v:=range list{
				for _,vv:=range v.List{
					tt, ok := uMap[vv.Id]
					if ok {
						vv.Head = tt.Head
						vv.Name = tt.Name
					}
				}

			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func coachDetail(c *gin.Context) {
	type ResultData struct {
		Name    string `json:"name" `
		Id      int64  `json:"id" `
		Start   int64  `json:"start" `
		End     int64  `json:"end" `
		Now     int64  `json:"now" `
		Max     int64  `json:"max" `
		Level   string `json:"level" `
		Reserve bool   `json:"reserve" `
	}

	type Result struct {
		Time int64        `json:"time" `
		List []ResultData `json:"list" `
	}

	type DetailResult struct {
		Id   int64  `json:"id" `
		Name string `json:"name" `
		Head string `json:"head" `
		Info string `json:"info" `
		Desc string `json:"desc" `
	}
	type ReturnResult struct {
		Detail DetailResult `json:"detail" `
		List   []*Result    `json:"list" `
	}

	type query struct {
		Id int64 `form:"id" json:"id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	start := (time.Now().Unix() / 86400) * 86400
	end := start + 8*86400
	result, err := manager.GetCoachService().WorkList(c, &app.InstitutionWorkListRequest{Uid: params.Id, Start: start, End: end})
	if err == nil {
		reResult := ReturnResult{
			List:   make([]*Result, 0),
			Detail: DetailResult{},
		}
		if len(result.List) == 0 {
			chapterResult, err := manager.GetCoachService().CoachList(c, &app.CoachListRequest{Id: params.Id, Page: 1, Size: 10,})
			if err == nil && len(chapterResult.List) > 0 {
				reResult.Detail = DetailResult{
					Id:   chapterResult.List[0].Uid,
					Info: chapterResult.List[0].Info,
				}
			} else {
				c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
				return
			}
		} else {
			reResult.Detail = DetailResult{
				Id:   result.List[0].Uid,
				Info: result.List[0].UserInfo,
			}
		}
		listMap := make(map[int64]*Result)
		for i := 0; i <= 6; i++ {
			tempTime := start + int64(i)*86400
			t := &Result{
				Time: tempTime,
				List: make([]ResultData, 0),
			}
			reResult.List = append(reResult.List, t)
			listMap[tempTime] = t
		}
		for _, v := range result.List {
			tempTime := (v.Start / 86400) * 86400
			t, ok := listMap[tempTime]
			if ok {
				t.List = append(t.List, ResultData{
					Id:      v.Id,
					Name:    v.Name,
					Start:   v.Start,
					End:     v.End,
					Now:     v.Now,
					Max:     v.Max,
					Level:   v.Level,
					Reserve: v.Reserve,
				})
			}
		}
		uResult, err := manager.GetUserService().GetUserList(c, &app.GetUserListRequest{IdList: []int64{params.Id}})
		if err == nil {
			for _, v := range uResult.List {
				reResult.Detail.Head = v.Head
				reResult.Detail.Name = v.Name
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(reResult))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func course(c *gin.Context) {
	type Place struct {
		Name   string  `json:"name" `
		Phone  string  `json:"phone" `
		Lat    float64 `json:"lat" `
		Lng    float64 `json:"lng" `
		Place  string  `json:"place" `
		Time   string  `json:"time" `
		Desc   string  `json:"desc" `
		Id     int64   `json:"id" `
		Select bool    `json:"select" `
	}
	type ResultData struct {
		Name    string `json:"name" `
		Desc    string `json:"desc" `
		Id      int64  `json:"id" `
		Start   int64  `json:"start" `
		End     int64  `json:"end" `
		Price   int64  `json:"price" `
		Time    int64  `json:"time" `
		Reserve bool   `json:"reserve" `
		Place   Place  `json:"place" `
	}

	type Detail struct {
		Id   int64  `json:"id" `
		Name string `json:"name" `
		Head string `json:"head" `
		Info string `json:"info" `
		Desc string `json:"desc" `
	}
	type Result struct {
		Detail Detail     `json:"detail" `
		Course ResultData `json:"course" `
	}

	type query struct {
		Id int64 `form:"id" json:"id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	result, err := manager.GetCoachService().WorkList(c, &app.InstitutionWorkListRequest{Id: []int64{params.Id}})
	if err == nil {
		if len(result.List) == 0 {
			c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
			return
		}
		v := result.List[0]
		reResult := Result{
			Detail: Detail{
				Id:   v.Uid,
				Info: v.UserInfo,
				Desc: "",
			},
			Course: ResultData{
				Name:    v.Name,
				Id:      v.Id,
				Start:   v.Start,
				End:     v.End,
				Desc:    v.Desc,
				Price:   v.Price,
				Time:    v.Duration,
				Reserve: v.Reserve,
			},
		}
		if v.Place != nil {
			tTemp := time.FixedZone("CST", 8*3600)
			reResult.Course.Place = Place{
				Id:    v.Place.Id,
				Name:  v.Place.Name,
				Place: v.Place.Address,
				Phone: v.Place.Phone,
				Time:  time.Unix(v.Place.Start, 0).In(tTemp).Format(enum.TimeFormatHour) + "点-" + time.Unix(v.Place.End, 0).In(tTemp).Format(enum.TimeFormatHour) + "点",
				Lat:   0,
				Lng:   0,
			}
		}
		uResult, err := manager.GetUserService().GetUserList(c, &app.GetUserListRequest{IdList: []int64{v.Uid}})
		if err == nil {
			for _, v := range uResult.List {
				reResult.Detail.Head = v.Head
				reResult.Detail.Name = v.Name
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(reResult))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func pay(c *gin.Context) {
	type query struct {
		Id      int64 `form:"id" json:"id" `
		Channel int64 `form:"channel" json:"channel" `
		Num     int64 `form:"num" json:"num" `
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
	list := make([]*app.OrderData, 0)
	data := &app.OrderData{
		SellerId: 0,
	}
	result, err := manager.GetCoachService().WorkList(c, &app.InstitutionWorkListRequest{Id: []int64{params.Id}})
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
			temp := make([]*app.OrderGood, 0)
			temp = append(temp, &app.OrderGood{
				Num:   params.Num,
				SkuId: v.Id,
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
	orderResult, err := manager.GetOrderService().Create(c, &app.OrderRequest{
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
	payResult, err := manager.GetOrderService().Pay(c, &app.OrderPayRequest{
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

func paySu(c *gin.Context) {
	type Param struct {
		Id      int64 `form:"id" json:"id"`
		OrderId int64 `form:"order_id" json:"order_id"`
	}
	params := Param{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	_, err := manager.GetOrderService().EditStatus(c, &app.EditStatusRequest{
		OrderId: params.OrderId,
		Status:  enum.OrderStatusPay,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err = manager.GetCoachService().WorkOrder(c, &app.InstitutionWorkOrderRequest{Id: params.Id, Uid: userData.Uid, Status: enum.OrderStatusPay,OrderId:params.OrderId})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
}

func receive(c *gin.Context) {
	type Param struct {
		Id int64 `form:"id" json:"id"`
	}
	params := Param{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	_, err := manager.GetOrderService().EditStatus(c, &app.EditStatusRequest{
		OrderId: params.Id,
		BuyerId: userData.Uid,
		Status:  enum.OrderStatusFinish,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err = manager.GetCoachService().WorkOrder(c, &app.InstitutionWorkOrderRequest{OrderId: params.Id, Uid: userData.Uid, Status: enum.OrderStatusFinish})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
}

func refund(c *gin.Context) {
	type Param struct {
		Id int64 `form:"id" json:"id"`
	}
	params := Param{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	_, err := manager.GetOrderService().PostOrderRefund(c, &app.PostRefundReq{
		OrderId: params.Id,
		BuyerId: userData.Uid,
		All:     true,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err = manager.GetCoachService().WorkOrder(c, &app.InstitutionWorkOrderRequest{OrderId: params.Id, Uid: userData.Uid, Status: enum.OrderStatusRefund})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
}

func orderList(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type Place struct {
		Name  string  `json:"name" `
		Phone string  `json:"phone" `
		Lat   float64 `json:"lat" `
		Lng   float64 `json:"lng" `
		Place string  `json:"place" `
		Time  string  `json:"time" `
		Desc  string  `json:"desc" `
		Id    int64   `json:"id" `
	}

	type Coach struct {
		Id   int64  `json:"id" `
		Name string `json:"name" `
		Head string `json:"head" `
		Info string `json:"info" `
		Desc string `json:"desc" `
	}
	type Result struct {
		Name    string `json:"name" `
		Level   string `json:"level" `
		Id      int64  `json:"id" `
		Start   int64  `json:"start" `
		End     int64  `json:"end" `
		Status  int64  `json:"status" `
		OrderId int64  `json:"order_id" `
		Num     int64  `json:"num" `
		Coach   Coach  `json:"coach" `
		Place   Place  `json:"place" `
	}
	type query struct {
		Page   int64 `form:"page" json:"page" `
		Size   int64 `form:"size" json:"size" `
		Status int64 `form:"status" json:"status" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	statusList := make([]int64, 0)
	if params.Status == 1 {
		statusList = append(statusList, enum.OrderStatusPay)
	} else if params.Status==2{
		statusList = append(statusList, enum.OrderStatusFinish, enum.OrderStatusRefund)
	}else{
		statusList = append(statusList,enum.OrderStatusPay, enum.OrderStatusFinish, enum.OrderStatusRefund)
	}
	result, err := manager.GetCoachService().WorkOrderList(c, &app.InstitutionWorkOrderRequest{
		Uid: userData.Uid, QueryStatus: statusList, Page: params.Page, Size: params.Size,
	})
	if err == nil {
		workIdList := make([]int64, 0)
		reList := make([]*Result, 0)
		for _, v := range result.List {
			workIdList = append(workIdList, v.WorkId)
			temp := &Result{
				Id:      v.WorkId,
				OrderId: v.OrderId,
				Num:     v.Num,
				Status:  v.Status,
			}
			reList = append(reList, temp)
		}

		workResult, err := manager.GetCoachService().WorkList(c, &app.InstitutionWorkListRequest{Id: workIdList})
		if err == nil {
			workMap := make(map[int64]*app.InstitutionWorkData)
			coachUidList := make([]int64, 0)
			for _, v := range workResult.List {
				workMap[v.Id] = v
			}
			tTemp := time.FixedZone("CST", 8*3600)
			for _, temp := range reList {
				v, ok := workMap[temp.Id]
				if ok {
					temp.Coach = Coach{
						Id:   v.Uid,
						Info: v.UserInfo,
						Desc: "",
					}
					temp.Name = v.Name
					temp.Start = v.Start
					temp.End = v.End
					temp.Level = v.Level
					if v.Place != nil {
						temp.Place = Place{
							Id:    v.Place.Id,
							Name:  v.Place.Name,
							Place: v.Place.Address,
							Phone: v.Place.Phone,
							Time:  time.Unix(v.Place.Start, 0).In(tTemp).Format(enum.TimeFormatHour) + "点-" + time.Unix(v.Place.End, 0).In(tTemp).Format(enum.TimeFormatHour) + "点",
							Lat:   0,
							Lng:   0,
						}
					}
					coachUidList = append(coachUidList, v.Uid)
				}
			}
			uResult, err := manager.GetUserService().GetUserList(c, &app.GetUserListRequest{IdList: coachUidList})
			if err == nil {
				uResultMap := make(map[int64]*app.UserDataResponse)
				for _, v := range uResult.List {
					uResultMap[v.Uid] = v
				}
				for _, v := range reList {
					temp, ok := uResultMap[v.Coach.Id]
					if ok {
						v.Coach.Name = temp.Name
						v.Coach.Head = temp.Head
					}
				}
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(reList))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func orderDetail(c *gin.Context) {
	type Place struct {
		Name  string  `json:"name" `
		Phone string  `json:"phone" `
		Lat   float64 `json:"lat" `
		Lng   float64 `json:"lng" `
		Place string  `json:"place" `
		Time  string  `json:"time" `
		Desc  string  `json:"desc" `
		Id    int64   `json:"id" `
	}

	type User struct {
		Id   int64  `json:"id" `
		Name string `json:"name" `
		Head string `json:"head" `
		Info string `json:"info" `
		Desc string `json:"desc" `
	}
	type Order struct {
		OrderSn string `json:"order_sn" `
		Create  int64  `json:"create" `
		Pay     int64  `json:"pay" `
	}
	type Detail struct {
		Name    string `json:"name" `
		Desc    string `json:"desc" `
		Id      int64  `json:"id" `
		Start   int64  `json:"start" `
		End     int64  `json:"end" `
		Price   int64  `json:"price" `
		Reserve bool   `json:"reserve" `
		Time    int64  `json:"time" `
		Place   Place  `json:"place" `
	}
	type Result struct {
		Detail Detail `json:"detail" `
		User   User   `json:"user" `
		Order  Order  `json:"order" `
	}
	type query struct {
		Id int64 `form:"id" json:"id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	result, err := manager.GetCoachService().WorkOrderList(c, &app.InstitutionWorkOrderRequest{
		OrderId: params.Id,
	})
	if err == nil {
		if len(result.List) == 0 {
			c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
			return
		}
		v := result.List[0]
		reResult := Result{
			Detail: Detail{

			},
			User:  User{},
			Order: Order{},
		}

		workResult, err := manager.GetCoachService().WorkList(c, &app.InstitutionWorkListRequest{Id: []int64{v.WorkId}})
		if err == nil {
			if len(workResult.List) == 0 {
				c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
				return
			}
			reResult.User = User{
				Id:   workResult.List[0].Uid,
				Info: workResult.List[0].UserInfo,
				Desc: "",
			}
			tTemp := time.FixedZone("CST", 8*3600)
			reResult.Detail = Detail{
				Name:    workResult.List[0].Name,
				Start:   workResult.List[0].Start,
				End:     workResult.List[0].End,
				Desc:    workResult.List[0].Desc,
				Id:      v.WorkId,
				Reserve: true,
				Price:   workResult.List[0].Price,
				Time:    workResult.List[0].Duration,
			}
			if workResult.List[0].Place != nil {
				reResult.Detail.Place = Place{
					Id:    workResult.List[0].Place.Id,
					Name:  workResult.List[0].Place.Name,
					Place: workResult.List[0].Place.Address,
					Phone: workResult.List[0].Place.Phone,
					Time:  time.Unix(workResult.List[0].Place.Start, 0).In(tTemp).Format(enum.TimeFormatHour) + "点-" + time.Unix(workResult.List[0].Place.End, 0).In(tTemp).Format(enum.TimeFormatHour) + "点",
					Lat:   0,
					Lng:   0,
				}
			}
			uResult, err := manager.GetUserService().GetUserList(c, &app.GetUserListRequest{IdList: []int64{workResult.List[0].Uid}})
			if err == nil {
				for _, v := range uResult.List {
					reResult.User.Name = v.Name
					reResult.User.Head = v.Head
				}
			}
		}
		orderResult, err := manager.GetOrderService().GetOrders(c, &app.GetOrderRequest{OrderId: v.OrderId, ShowOrder: enum.NoShowOrderList})
		if err == nil && len(orderResult.List) > 0 {
			reResult.Order.OrderSn = orderResult.List[0].OrderSn
			reResult.Order.Create = orderResult.List[0].OrderTime
			reResult.Order.Pay = orderResult.List[0].PayChannel
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(reResult))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}
