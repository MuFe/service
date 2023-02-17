package school

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"mufe_service/camp/cache"
	"mufe_service/camp/dataUtil"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/handler"
	"mufe_service/camp/jwt"
	"mufe_service/camp/server"
	"mufe_service/camp/xlog"
	app "mufe_service/jsonRpc"
	"mufe_service/manager"
)

func init() {

	server.Post("/appApi/gradeInfo", gradeInfo)
	server.Post("/appApi/classInfo", classInfo)

	server.Post("/appApi/quitClass", handler.UserLogin, quitClass)
	server.Post("/appApi/disClass", handler.UserLogin, disClass)
	server.Post("/appApi/joinClass", handler.UserLogin, joinClass)
	server.Post("/appApi/createClass", handler.UserLogin, handler.TeacherCheck, createClass)

	server.Post("/appApi/mySchool", handler.UserLogin, mySchool)
	server.Post("/appApi/myClass", handler.UserLogin, myClass)
	server.Post("/appApi/classList", handler.UserLogin, classList)
	server.Post("/appApi/classDetail", handler.UserLogin, classDetail)

	server.Post("/appApi/joinSchool", handler.UserLogin, addSchool)
	server.Post("/appApi/quitSchool", handler.UserLogin, quitSchool)
	server.Post("/appApi/scanSchool", scanSchool)
	server.Post("/appApi/scan", scan)

	server.Post("/appApi/course/joinCourse", handler.UserLogin, handler.TeacherCheck, joinCourse)
	server.Post("/appApi/course/removeCourse", handler.UserLogin, handler.TeacherCheck, removeCourse)

}

func gradeInfo(c *gin.Context) {
	type Grade struct {
		Id   int64  `json:"id" `
		Name string `json:"name" `
	}
	type Result struct {
		Name string  `json:"name" `
		List []Grade `json:"list" `
	}
	type query struct {
		Id int64 `form:"id" json:"id" `
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
	result, err := manager.GetSchoolService().GradeList(c, &app.GradeRequest{SchoolId: params.Id})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	list := make([]Result, 0)
	for _, v := range result.List {
		tempList := make([]Grade, 0)
		for _, vv := range v.List {
			tempList = append(tempList, Grade{
				Id:   vv.Id,
				Name: vv.Name,
			})
		}
		list = append(list, Result{
			Name: v.Type,
			List: tempList,
		})
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
}

func mySchool(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	result, err := manager.GetSchoolService().MySchool(c, &app.SchoolRequest{Uid: userData.Uid})
	if err == nil {
		type Result struct {
			Id      int64  `json:"id" `
			Title   string `json:"title" `
			Address string `json:"address" `
			Icon    string `json:"icon" `
		}
		list := make([]Result, 0)
		for _, v := range result.List {
			list = append(list, Result{
				Id:      v.Id,
				Title:   v.Name,
				Address: v.Address,
				Icon:    v.Icon,
			})
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func myClass(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	result, err := manager.GetSchoolService().MySchool(c, &app.SchoolRequest{Uid: userData.Uid})
	if err == nil {
		type Result struct {
			Id      int64  `json:"id" `
			Title   string `json:"title" `
			Address string `json:"address" `
			Icon    string `json:"icon" `
		}
		list := make([]Result, 0)
		for _, v := range result.List {
			list = append(list, Result{
				Id:      v.Id,
				Title:   v.Name,
				Address: v.Address,
				Icon:    v.Icon,
			})
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func scanSchool(c *gin.Context) {
	type query struct {
		Content string `form:"content" json:"content" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	v, err := manager.GetSchoolService().Scan(c, &app.ScanRequest{Content: params.Content})
	if err == nil {
		if v.School == nil {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseMsg("无效的学校二维码"))
		} else {
			type Result struct {
				Id    int64  `json:"id" `
				Title string `json:"title" `
				Desc  string `json:"desc" `
				Icon  string `json:"icon" `
			}
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(Result{
				Id:    v.School.Id,
				Title: v.School.Name,
				Desc:  v.School.Desc,
				Icon:  v.School.Icon,
			}))
		}

	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func scan(c *gin.Context) {
	type query struct {
		Content string `form:"content" json:"content" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	v, err := manager.GetSchoolService().Scan(c, &app.ScanRequest{Content: params.Content})
	type Result struct {
		Id    int64  `json:"id" `
		Title string `json:"title" `
		Desc  string `json:"desc" `
		Icon  string `json:"icon" `
		Type  int64  `json:"type" `
	}
	if err == nil {

		result := Result{}
		if v.School != nil {
			result = Result{
				Id:    v.School.Id,
				Title: v.School.Name,
				Desc:  v.School.Address,
				Icon:  v.School.Icon,
				Type:  enum.SCAN_SCHOOL_TYPE,
			}
		} else {
			result = Result{
				Id:    v.Class.Id,
				Title: v.Class.Name,
				Desc:  v.Class.Grade,
				Type:  enum.SCAN_CLASS_TYPE,
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
	} else {
		tResult, err := manager.GetCoachService().GetInstitution(c, &app.InstitutionRequest{Code: params.Content, Status: enum.StatusNormal})
		if err == nil {
			if len(tResult.List) == 0 {
				c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(xlog.Error("无法识别二维码")))
				return
			}
			result := Result{
				Id:    tResult.List[0].Id,
				Title: tResult.List[0].Name,
				Desc:  tResult.List[0].Address,
				Icon:  tResult.List[0].Icon,
				Type:  enum.SCAN_INSTITUTION_TYPE,
			}
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
			return
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func quitSchool(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Id int64 `form:"id" json:"id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetSchoolService().QuitSchool(c, &app.AddSchoolRequest{Uid: userData.Uid, SchoolId: params.Id})
	if err == nil {
		result, err := manager.GetUserService().UpdateUser(c, &app.UpdateUserRequest{Uid: userData.Uid, IdentityType: enum.DEFAULT_TYPE, Type: enum.UpdateUserIdentity})
		if err == nil {
			token, err := jwt.GenerateUserJwt(userData.Uid, enum.DEFAULT_TYPE, userData.OpenId)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
				return
			}
			user := dataUtil.ParseUser(result)
			//将常用的用户信息存储到Redis
			u := &cache.UserClaims{
				Name:           user.Name,
				Head:           user.Head,
				Phone:          user.Phone,
				Sex:            user.Sex,
				UserNo:         user.UserNo,
				UserInviteCode: user.UserInviteCode,
				Sign:           user.Sign,
				Address:        user.Address,
				Identity:       user.Identity,
				Age:            user.Age,
				HaveWx:         user.HaveWx,
				HavePass:       user.HavePass,
			}
			_ = cache.SetUserInfo(result.Uid, u)
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(token))
		} else {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		}
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func addSchool(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Id int64 `form:"id" json:"id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetSchoolService().AddSchool(c, &app.AddSchoolRequest{Uid: userData.Uid, SchoolId: params.Id})
	if err == nil {
		result, err := manager.GetUserService().UpdateUser(c, &app.UpdateUserRequest{Uid: userData.Uid, IdentityType: enum.TEACHER_TYPE, Type: enum.UpdateUserIdentity})
		if err == nil {
			token, err := jwt.GenerateUserJwt(userData.Uid, enum.TEACHER_TYPE, userData.OpenId)
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

func createClass(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}

	type query struct {
		Id int64 `form:"id" json:"id" `
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
	_, err := manager.GetSchoolService().CreateClass(c, &app.CreateClassRequest{Uid: userData.Uid, ClassInfoId: params.Id})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("创建班级成功"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func joinClass(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Id int64 `form:"id" json:"id" `
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
	typeInt := enum.ClassStudentType
	if userData.Identity == enum.TEACHER_TYPE {
		typeInt = enum.ClassAdminType
	} else if userData.Identity == enum.STUDENT_TYPE {
		c.AbortWithStatusJSON(http.StatusOK, xlog.Error("您已经加入班级"))
		return
	}
	_, err := manager.GetSchoolService().JoinClass(c, &app.JoinClassRequest{Uid: userData.Uid, ClassId: params.Id, Type: typeInt})
	if err == nil {
		if typeInt == enum.ClassStudentType {
			result, err := manager.GetUserService().UpdateUser(c, &app.UpdateUserRequest{Uid: userData.Uid, IdentityType: enum.STUDENT_TYPE, Type: enum.UpdateUserIdentity})
			if err == nil {
				u := dataUtil.ParseUserCache(result)
				_ = cache.SetUserInfo(result.Uid, u)
			}
			v, err := manager.GetSchoolService().ClassDetail(c, &app.ClassDetailRequest{Id: params.Id, Uid: userData.Uid})
			if err == nil {
				ids := make([]int64, 0)
				ids = append(ids, v.AdminList...)
				ids = append(ids, userData.Uid)
				userResult, err := manager.GetUserService().GetUserList(c, &app.GetUserListRequest{IdList: ids})
				if err == nil {
					deviceList := make([]string, 0)
					userName := ""
					content := "%s已经加入了您的班级【%s】"
					for _, v := range userResult.List {
						if v.Uid != userData.Uid {
							deviceList = append(deviceList, v.RegistrationId)
						} else {
							userName = v.Name
						}

					}
					manager.GetPushService().PushMessage(c, &app.PushRequest{DeviceList: deviceList, Content: fmt.Sprintf(content, userName, v.Grade+v.Name)})
				}
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("加入成功"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func disClass(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Id int64 `form:"class_id" json:"class_id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	result, err := manager.GetSchoolService().DissolutionClass(c, &app.QuitClassRequest{Uid: userData.Uid, ClassId: params.Id})
	if err == nil {
		_, err = manager.GetUserService().UpdateUser(c, &app.UpdateUserRequest{UidList: result.StudentList, IdentityType: enum.DEFAULT_TYPE, Type: enum.UpdateUserIdentity})
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("解散成功"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func quitClass(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Id  int64 `form:"class_id" json:"class_id" `
		Uid int64 `form:"uid" json:"uid" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	uid := params.Uid
	if uid == 0 {
		uid = userData.Uid
	}
	_, err := manager.GetSchoolService().QuitClass(c, &app.QuitClassRequest{Uid: userData.Uid, ClassId: params.Id, QuitUid: uid})
	if err == nil {
		if userData.Uid != uid {
			result, err := manager.GetUserService().UpdateUser(c, &app.UpdateUserRequest{Uid: uid, IdentityType: enum.DEFAULT_TYPE, Type: enum.UpdateUserIdentity})
			if err == nil {
				u := dataUtil.ParseUserCache(result)
				_ = cache.SetUserInfo(result.Uid, u)
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("退出成功"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func classInfo(c *gin.Context) {
	type query struct {
		Grade  int64 `form:"id" json:"id" `
		School int64 `form:"school" json:"school" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Grade == 0 || params.School == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	result, err := manager.GetSchoolService().ClassInfo(c, &app.ClassInfoRequest{GradeId: params.Grade, SchoolId: params.School})
	if err == nil {
		type Result struct {
			Id         int64  `json:"id" `
			Grade      string `json:"grade" `
			SchoolType string `json:"school_type" `
			Class      string `json:"class" `
		}
		list := make([]Result, 0)
		for _, v := range result.List {
			list = append(list, Result{
				Id:         v.Id,
				Grade:      v.Grade,
				SchoolType: v.SchoolType,
				Class:      v.Name,
			})
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func joinCourse(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		ClassId  int64 `form:"class_id" json:"class_id" `
		CourseId int64 `form:"course_id" json:"course_id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.CourseId == 0 || params.ClassId == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	_, err := manager.GetSchoolService().AddCourse(c, &app.AddCourseRequest{Uid: userData.Uid, CourseId: params.CourseId, ClassId: params.ClassId})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("添加成功"))
	}
}

func removeCourse(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Id int64 `form:"id" json:"id" `
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
	_, err := manager.GetSchoolService().RemoveCourse(c, &app.AddCourseRequest{Uid: userData.Uid, ClassId: params.Id})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("移除成功"))
	}
}

func classList(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	result, err := manager.GetSchoolService().ClassList(c, &app.ClassRequest{Uid: userData.Uid})
	if err == nil {
		type Result struct {
			Id         int64  `json:"id" `
			Grade      string `json:"grade" `
			SchoolType string `json:"school_type" `
			Class      string `json:"class" `
			CreateTime int64  `json:"create_time" `
			Head       string `json:"head" `
			Name       string `json:"name" `
			School     string `json:"school" `
			Icon       string `json:"icon" `
			Tag        string `json:"tag" `
			Number     int64  `json:"number" `
		}
		list := make([]*Result, 0)
		uidList := make([]int64, 0)
		uidMap := make(map[int64]int64)
		for _, v := range result.List {
			temp := &Result{
				Id:         v.Id,
				Grade:      v.Grade,
				SchoolType: v.SchoolType,
				Class:      v.Name,
				CreateTime: v.CreateTime,
				Tag:        v.Tag,
				Number:     v.Number,
				Icon:       v.SchoolIcon,
				School:     v.SchoolName,
			}

			list = append(list, temp)
			if userData.Identity==enum.STUDENT_TYPE{
				uidList = append(uidList, v.CreateBy)
				uidMap[v.Id] = v.CreateBy
			}else{
				uidList = append(uidList, v.Uid)
				uidMap[v.Id] = v.Uid
			}
		}
		xlog.Info(uidList)
		if len(uidList) > 0 {
			userResult, err := manager.GetUserService().GetUserList(c, &app.GetUserListRequest{IdList: uidList, Status: enum.StatusNormal})
			if err == nil {
				userMap := make(map[int64]*app.UserDataResponse)
				for _, v := range userResult.List {
					userMap[v.Uid] = v
				}
				for _, v := range list {
					uidTemp, ok := uidMap[v.Id]
					if ok {
						temp, ok := userMap[uidTemp]
						if ok {
							v.Head = temp.Head
							v.Name = temp.Name
						}
					}
				}
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func classDetail(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		ClassId int64 `form:"id" json:"id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	v, err := manager.GetSchoolService().ClassDetail(c, &app.ClassDetailRequest{Id: params.ClassId, Uid: userData.Uid})
	if err == nil {
		type UserResult struct {
			Id   int64  `json:"id" `
			Name string `json:"name" `
			Head string `json:"head" `
			Type string `json:"type" `
		}
		type Result struct {
			Id         int64        `json:"id" `
			Grade      string       `json:"grade" `
			SchoolType string       `json:"school_type" `
			Class      string       `json:"class" `
			CreateTime int64        `json:"create_time" `
			Name       string       `json:"name" `
			Tag        string       `json:"tag" `
			Code       string       `json:"code" `
			Icon       string       `json:"icon" `
			Title      string       `json:"title" `
			CourseId   int64        `json:"course_id" `
			Number     int64        `json:"number" `
			Progress   float64      `json:"progress" `
			AdminList  []UserResult `json:"admin_list" `
			List       []UserResult `json:"list" `
		}

		uidList := make([]int64, 0)
		for _, value := range v.AdminList {
			uidList = append(uidList, value)
		}
		for _, value := range v.StudentList {
			uidList = append(uidList, value)
		}
		uidList = append(uidList, userData.Uid)
		userResult, err := manager.GetUserService().GetUserList(c, &app.GetUserListRequest{IdList: uidList, Status: enum.StatusNormal})
		name := ""
		adminList := make([]UserResult, 0)
		list := make([]UserResult, 0)

		if err == nil {
			userMap := make(map[int64]*app.UserDataResponse)
			for _, value := range userResult.List {
				userMap[value.Uid] = value
			}
			temp, ok := userMap[userData.Uid]
			if ok {
				name = temp.Name
			}

			for _, value := range v.AdminList {
				temp, ok := userMap[value]

				if ok {
					adminList = append(adminList, UserResult{
						Id:   value,
						Type: "老师",
						Head: temp.Head,
						Name: temp.Name,
					})
				}
			}
			for _, value := range v.StudentList {
				temp, ok := userMap[value]
				if ok {
					list = append(list, UserResult{
						Id:   value,
						Type: "学生",
						Head: temp.Head,
						Name: temp.Name,
					})
				}
			}

		}
		title := ""
		if v.CourseId != 0 {
			courseResult, err := manager.GetCourseService().GetCourse(c, &app.CourseServiceRequest{Ids: []int64{v.CourseId}, Status: enum.StatusNormal})
			if err == nil && len(courseResult.List) > 0 {
				title = fmt.Sprintf("%s%s%d%s%d%s", courseResult.List[0].Title, ":第", v.ChapterIndex, "章/第", v.VideoIndex, "节")
			}
		}
		if v.VideoId != 0 {
			videoResult, err := manager.GetVideoService().GetVideo(c, &app.VideoRequest{VideoId: v.VideoId})
			if err == nil && len(videoResult.VideoList) > 0 {
				title = title + videoResult.VideoList[0].Title
			}
		}

		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(Result{
			Id:         v.Id,
			Grade:      v.Grade,
			SchoolType: v.SchoolType,
			Class:      v.Name,
			CreateTime: v.CreateTime,
			Tag:        v.Tag,
			Name:       name,
			Code:       v.Code,
			AdminList:  adminList,
			CourseId:   v.CourseId,
			Title:      title,
			List:       list,
			Number:     int64(len(v.StudentList)),
			Progress:   v.Progress,
		}))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}
