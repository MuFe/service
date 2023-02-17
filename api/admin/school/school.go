package school

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/handler"
	"mufe_service/camp/server"
	"mufe_service/camp/utils"
	"mufe_service/jsonRpc"
	"mufe_service/manager"
	"strconv"
	"time"
)

func init() {
	server.Post("/adminSchool/list", handler.AdminLogin, list)
	server.Post("/adminSchool/schoolType", schoolType)
	server.Post("/adminSchool/token", handler.AdminLogin, getToken)
	server.Post("/adminSchool/editSchool", handler.AdminLogin, editSchool)
	server.Post("/adminSchool/editSchoolIcon", handler.AdminLogin, editSchoolIcon)
	server.Post("/adminSchool/createClassInfo", handler.AdminLogin, createClassInfo)
	server.Post("/adminSchool/classStudentList", handler.AdminLogin, classStudentList)
	server.Post("/adminSchool/schoolGradeList", schoolGradeList)
	server.Post("/adminSchool/schoolGradeClassList", schoolGradeClassList)
	server.Post("/adminSchool/schoolClassList", schoolClassList)
}

func list(c *gin.Context) {
	type query struct {
		Id int64 `form:"id" json:"id" `
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	result, err := manager.GetSchoolService().SchoolList(c, &app.SchoolRequest{Id: params.Id})
	if err == nil {
		type School struct {
			Id      int64  `json:"id" `
			Title   string `json:"title" `
			Address string `json:"school_address" `
			Icon    string `json:"icon" `
			Code    string `json:"code" `
			TypeId  int64  `json:"type_id" `
			Type    string `json:"type_name" `
		}
		list := make([]interface{}, 0)
		for _, v := range result.List {
			list = append(list, School{
				Id:      v.Id,
				Title:   v.Name,
				Icon:    v.Icon,
				Address: v.Address,
				Code:    v.Code,
				TypeId:  v.TypeId,
				Type:    v.TypeName,
			})
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(utils.CreateListResultReturn(result.Total, list)))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func editSchool(c *gin.Context) {
	type query struct {
		Name    string `form:"name" json:"name" `
		Address string `form:"school_address" json:"school_address" `
		Id      int64  `form:"id" json:"id" `
		Type    int64  `form:"type" json:"type" `
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	result, err := manager.GetSchoolService().EditSchool(c, &app.SchoolData{
		Name:    params.Name,
		Address: params.Address,
		Id:      params.Id,
		TypeId:  params.Type,
	})
	if err == nil {
		if params.Id == 0 {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result.Id))
		} else {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result.Id))
		}

	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func editSchoolIcon(c *gin.Context) {
	type query struct {
		Photo  string `form:"photo" json:"photo" `
		Prefix string `form:"prefix" json:"prefix" `
		Id     int64  `form:"id" json:"id" `
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetSchoolService().EditSchool(c, &app.SchoolData{
		Id:   params.Id,
		Icon: params.Prefix + params.Photo,
	})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("修改成功"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func getToken(c *gin.Context) {
	type query struct {
		Name []string `form:"names" json:"names" `
		IsBase  bool     `form:"base64" json:"base64" `
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if len(params.Name) == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	filenameWithSuffix := path.Base(params.Name[0])
	fileSuffix := path.Ext(filenameWithSuffix) //获取文件后缀
	encodeString := utils.MD5(params.Name[0]+strconv.FormatInt(time.Now().Unix(), 10)) + fileSuffix
	if params.IsBase {
		encodeString = base64.StdEncoding.EncodeToString([]byte(encodeString))
	}
	osStr := os.Getenv("IMG_BUCKET")
	prefix := os.Getenv("IMG_PREFIX")
	type QiniuInfo struct {
		Token    string   `json:"token"`
		Host     string   `json:"host"`
		BaseHost string   `json:"base_host"`
		Keys     []string `json:"keys"`
		Prefix   string   `json:"prefix"`
	}

	result, err := manager.GetQiniuService().GetToken(c, &app.QiniuServiceRequest{Bucket: osStr})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(&QiniuInfo{Token: result.Token, Host: result.UploadHost, BaseHost: result.Base64UploadHost, Keys: []string{encodeString}, Prefix: prefix}))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func createClassInfo(c *gin.Context) {
	type query struct {
		Name  []string `form:"name" json:"name" `
		Id    int64  `form:"id" json:"id" `
		Grade int64  `form:"grade" json:"grade" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetSchoolService().CreateClassInfo(c, &app.CreateClassInfoRequest{SchoolId: params.Id, GradeId: params.Grade, Name: params.Name})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("添加成功"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func schoolType(c *gin.Context) {
	type Grade struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}
	type SchoolType struct {
		Id   int64   `json:"id"`
		Name string  `json:"name"`
		List []Grade `json:"list"`
	}
	list := make([]SchoolType, 0)
	result, err := manager.GetSchoolService().SchoolTypeList(c, &app.EmptyRequest{})
	if err == nil {
		for _, v := range result.List {
			data := SchoolType{
				Id:   v.Id,
				Name: v.Name,
				List: make([]Grade, 0),
			}

			for _, vv := range v.Grade {
				data.List = append(data.List, Grade{
					Id:   vv.Id,
					Name: vv.Name,
				})
			}
			list = append(list, data)
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func schoolGradeList(c *gin.Context) {
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
	if err == nil {
		type ClassInfo struct {
			Id   int64  `json:"id" `
			Name string `json:"name" `
		}
		type Grade struct {
			Id   int64  `json:"id" `
			Name string `json:"name" `
			Info []ClassInfo `json:"list" `
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

func schoolGradeClassList(c *gin.Context) {
	type query struct {
		Id    int64 `form:"id" json:"id" `
		Grade int64 `form:"grade" json:"grade" `
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
	result, err := manager.GetSchoolService().ClassInfo(c, &app.ClassInfoRequest{GradeId: params.Grade, SchoolId: params.Id})
	if err == nil {
		type ClassInfo struct {
			Id   int64  `json:"id" `
			Name string `json:"name" `
		}
		list := make([]ClassInfo, 0)
		for _, v := range result.List {
			list = append(list, ClassInfo{
				Id:   v.Id,
				Name: v.Name,
			})
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func schoolClassList(c *gin.Context) {
	type query struct {
		Id    int64 `form:"id" json:"id" `
		Grade int64 `form:"grade" json:"grade" `
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
	result, err := manager.GetSchoolService().ClassList(c, &app.ClassRequest{School: params.Id, Grade: params.Grade})
	if err == nil {
		type Grade struct {
			Id     int64  `json:"id" `
			Name   string `json:"name" `
			Number int64  `json:"number" `
		}
		list := make([]Grade, 0)
		for _, v := range result.List {
			list = append(list, Grade{
				Id:     v.Id,
				Name:   v.Name,
				Number: v.Number,
			})
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func classStudentList(c *gin.Context) {
	type query struct {
		School   int64 `form:"school_id" json:"school_id" `
		ClassId   int64 `form:"class_id" json:"class_id" `
		GradeId   int64 `form:"grade_id" json:"grade_id" `
		Type int64 `form:"type" json:"type" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	typeInt:=enum.ClassAdminType
	if params.Type==1{
		typeInt=enum.ClassStudentType
	}
	result, err := manager.GetSchoolService().UserList(c, &app.TeacherUserRequest{SchoolId: params.School,ClassId:params.ClassId,GradeId:params.GradeId, Type: typeInt})
	if err == nil {
		type User struct {
			Id      int64  `json:"id" `
			Name    string `json:"name" `
			Phone   string `json:"phone" `
			Address string `json:"address" `
			Sex     int64  `json:"sex" `
			Age     int64  `json:"age" `
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
		uResult, err := manager.GetUserService().GetUserList(c, &app.GetUserListRequest{IdList: result.List })
		if err == nil {
			for _, v := range uResult.List {
				temp,ok:=resultMap[v.Uid]
				if ok{
					temp.Name=v.Name
					temp.Phone=v.Phone
					temp.Address=v.Address
					temp.Age=v.Age
					temp.Sex=v.Sex
				}
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}
