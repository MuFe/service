package football

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"mufe_service/camp/cache"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/jwt"
	"mufe_service/camp/server"
	pb "mufe_service/jsonRpc"
	"mufe_service/manager"
)

func init() {
	server.Get("/footballApi/teacher", check, checkTeacher)
	server.Get("/footballApi/grade", check, checkGrade)
	server.Get("/footballApi/class", check, checkClass)
	server.Get("/footballApi/app", check, checkList)

}

func check(c *gin.Context) {
	uid := int64(0)
	typeInt := int64(1)
	userHeader := c.GetHeader(jwt.AuthHeader)
	if userHeader != "" {
		claims, err := jwt.CheckUserJwt(userHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
			return
		}
		uid = claims.Uid
		typeInt = 1
	} else {
		userHeader = c.GetHeader(jwt.AdminAuthHeader)
		if userHeader != "" {
			claims, err := jwt.CheckAdminJwt(cache.AgentToken, userHeader)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
				return
			}
			uid = claims.Uid
			typeInt = 2
		} else {
			c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
			return
		}
	}
	scResult, err := manager.GetFootBallService().GetSchool(c, &pb.FootBallSchoolRequest{Uid: uid, Type: typeInt})
	if err == nil {
		c.Set("school", scResult.SchoolId)
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}

}

func checkTeacher(c *gin.Context) {
	type query struct {
		ClassId int64 `form:"class_id" json:"class_id" `
		GradeId int64 `form:"grade_id" json:"grade_id" `
	}

	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	schoolId, ok := c.MustGet("school").(int64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	result, err := manager.GetSchoolService().UserList(c, &pb.TeacherUserRequest{SchoolId: schoolId, ClassId: params.ClassId, GradeId: params.GradeId, Type: enum.ClassAdminType})
	if err == nil {
		type User struct {
			Id   int64  `json:"id" `
			Name string `json:"name" `
		}
		list := make([]*User, 0)
		resultMap := make(map[int64]*User)
		for _, v := range result.List {
			temp := &User{
				Id: v,
			}
			list = append(list, temp)
			resultMap[v] = temp
		}
		uResult, err := manager.GetUserService().GetUserList(c, &pb.GetUserListRequest{IdList: result.List})
		if err == nil {
			for _, v := range uResult.List {
				temp, ok := resultMap[v.Uid]
				if ok {
					temp.Name = v.Name
				}
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}

}

func checkGrade(c *gin.Context) {
	type query struct {
		ClassId int64 `form:"class_id" json:"class_id" `
		GradeId int64 `form:"grade_id" json:"grade_id" `
	}

	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	schoolId, ok := c.MustGet("school").(int64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	result, err := manager.GetSchoolService().GradeList(c, &pb.GradeRequest{SchoolId: schoolId})
	if err == nil {
		type Grade struct {
			Id   int64  `json:"id" `
			Name string `json:"name" `
		}
		list := make([]Grade, 0)
		for _, v := range result.List {
			for _, vv := range v.List {
				list = append(list, Grade{
					Id:   vv.Id,
					Name: vv.Name,
				})
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}

}

func checkClass(c *gin.Context) {
	type query struct {
		Teacher int64 `form:"teacher_id" json:"teacher_id" `
		GradeId int64 `form:"grade_id" json:"grade_id" `
	}

	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	schoolId, ok := c.MustGet("school").(int64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	if params.Teacher != 0 {
		schoolId = 0
	}
	result, err := manager.GetSchoolService().ClassList(c, &pb.ClassRequest{School: schoolId, Grade: params.GradeId, Uid: params.Teacher})
	if err == nil {
		type Grade struct {
			Id   int64  `json:"id" `
			Name string `json:"name" `
		}
		list := make([]Grade, 0)
		for _, v := range result.List {
			list = append(list, Grade{
				Id:   v.Id,
				Name: v.Name,
			})
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}

}

func checkList(c *gin.Context) {
	type query struct {
		Teacher int64 `form:"teacher_id" json:"teacher_id" `
		GradeId int64 `form:"grade_id" json:"grade_id" `
		Class   int64 `form:"class" json:"class" `
		Start   int64 `form:"start" json:"start" `
		End     int64 `form:"end" json:"end" `
	}

	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	schoolId, ok := c.MustGet("school").(int64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	if params.Teacher != 0 {
		schoolId = 0
	}
	_, err := manager.GetSchoolService().UserList(c, &pb.TeacherUserRequest{SchoolId: schoolId, ClassId: params.Class, GradeId: params.GradeId, Type: enum.ClassStudentType})
	if err == nil {
		type List struct {
			Time     int64  `json:"time" `
			Name     string `json:"name" `
			People   int64  `json:"people" `
			Duration int64  `json:"duration" `
			Number   int64  `json:"number" `
			Total    int64  `json:"total" `
			Rate     string `json:"rate" `
		}
		type Result struct {
			List [] List `json:"list" `
		}


		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(Result{}))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}

}
