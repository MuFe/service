package course

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"mufe_service/camp/cache"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/handler"
	"mufe_service/camp/jwt"
	"mufe_service/camp/server"
	"mufe_service/camp/utils"
	pb "mufe_service/jsonRpc"
	"mufe_service/manager"
	"strconv"
	"time"
)

func init() {
	server.Post("/adminCoach/token", handler.AdminLogin, getToken)
	server.Get("/adminCoach/institution", handler.AdminLogin, institutionList)
	server.Get("/adminCoach/course", handler.AdminLogin, courseList)
	server.Get("/adminCoach/institutionDetail", handler.AdminLogin, institutionDetail)
	server.Post("/adminCoach/editInstitution", handler.AdminLogin, editInstitution)
	server.Post("/adminCoach/editSchool", handler.AdminLogin, editSchool)
	server.Post("/adminCoach/editCourse", handler.AdminLogin, editCourse)
	server.Post("/adminCoach/editInstitutionIcon", handler.AdminLogin, editInstitutionIcon)
	server.Post("/adminCoach/editSchoolIcon", handler.AdminLogin, editInstitutionSchoolIcon)

}

type QiniuInfo struct {
	Token    string   `json:"token"`
	Host     string   `json:"host"`
	BaseHost string   `json:"base_host"`
	Keys     []string `json:"keys"`
	Prefix   string   `json:"prefix"`
}

type Tag struct {
	Tag string `json:"tag" `
	Id  int64  `json:"id" `
}

func getToken(c *gin.Context) {
	type query struct {
		Name   []string `form:"names" json:"names" `
		IsBase bool     `form:"base64" json:"base64" `
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

	result, err := manager.GetQiniuService().GetToken(c, &pb.QiniuServiceRequest{Bucket: osStr})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(&QiniuInfo{Token: result.Token, Host: result.UploadHost, BaseHost: result.Base64UploadHost, Keys: []string{encodeString}, Prefix: prefix}))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func institutionList(c *gin.Context) {
	type Result struct {
		Name    string `json:"name" `
		Id      int64  `json:"id" `
		Address string `json:"address" `
		Icon    string `json:"icon" `
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

	chapterResult, err := manager.GetCoachService().GetInstitution(c, &pb.InstitutionRequest{Page: params.Page, Size: params.Size, Status: enum.StatusAll})
	if err == nil {
		list := make([]Result, 0)
		for _, v := range chapterResult.List {
			list = append(list, Result{
				Name:    v.Name,
				Id:      v.Id,
				Address: v.Address,
				Icon:    v.Icon,
			})
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func institutionDetail(c *gin.Context) {
	type Course struct {
		Name string `json:"name" `
		Id   int64  `json:"id" `
	}
	type School struct {
		Name    string   `json:"name" `
		Id      int64    `json:"id" `
		Phone   string   `json:"phone" `
		Icon    string   `json:"icon" `
		Address string   `json:"address" `
		Time    string   `json:"time" `
		Course  []Course `json:"course" `
	}
	type Result struct {
		Name       string   `json:"name" `
		Id         int64    `json:"id" `
		Address    string   `json:"address" `
		Icon       string   `json:"icon" `
		Code       string   `json:"code" `
		SchoolList []School `json:"school" `
	}

	type query struct {
		Id int64 `form:"id" json:"id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}

	result, err := manager.GetCoachService().GetInstitution(c, &pb.InstitutionRequest{Id: params.Id, Status: enum.StatusNormal})
	if err == nil {
		if len(result.List) == 0 {
			c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
			return
		}
		v := result.List[0]
		reResult := Result{
			Name:       v.Name,
			Id:         v.Id,
			Address:    v.Address,
			Icon:       v.Icon,
			Code:       v.Code,
			SchoolList: make([]School, 0),
		}
		tResult, err := manager.GetCoachService().SchoolList(c, &pb.InstitutionSchoolRequest{Iid: v.Id, Page: 1, Size: 20, Status: enum.StatusNormal})
		if err == nil {
			tTemp:=time.FixedZone("CST", 8*3600)
			for _, v := range tResult.List {
				cList := make([]Course, 0)
				for _, vv := range v.Course {
					cList = append(cList, Course{
						Id:   vv.Id,
						Name: vv.Name,
					})
				}
				reResult.SchoolList = append(reResult.SchoolList, School{
					Id:      v.Id,
					Name:    v.Name,
					Icon:    v.Icon,
					Address: v.Address,
					Phone:   v.Phone,
					Course:  cList,
					Time:    time.Unix(v.Start, 0).In(tTemp).Format(enum.TimeFormatHour) + "点-" + time.Unix(v.End, 0).In(tTemp).Format(enum.TimeFormatHour) + "点",
				})
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(reResult))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func editInstitution(c *gin.Context) {
	type query struct {
		Name    string `form:"name" json:"name" `
		Address string `form:"address" json:"address" `
		Id      int64  `form:"id" json:"id" `
		Del     bool   `form:"del" json:"del" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Del == true && params.Id == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	adminData, ok := c.MustGet(handler.AdminData).(*jwt.AdminClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	uidList:=make([]int64,0)
	if params.Del{
		uResult, err :=manager.GetCoachService().CoachList(c, &pb.CoachListRequest{Iid:params.Id})
		if err == nil {
			for _,v:=range uResult.List{
				uidList=append(uidList,v.Uid)
			}

		}
	}

	result, err := manager.GetCoachService().EditInstitution(c, &pb.EditInstitutionRequest{Id: params.Id, Del: params.Del, Name: params.Name, Address: params.Address, Create: adminData.Uid})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	} else {
		if params.Del{
			manager.GetUserService().BatchModifyType(c,&pb.BatchModifyRequest{List:uidList,Type:enum.DEFAULT_TYPE})
			for _,v:=range uidList{
				cache.DeleteUserToken(v)
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result.Id))
	}
}

func editInstitutionIcon(c *gin.Context) {
	type query struct {
		Prefix string `form:"prefix" json:"prefix" `
		Icon   string `form:"icon" json:"icon" `
		Id     int64  `form:"id" json:"id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	result, err := manager.GetCoachService().EditInstitution(c, &pb.EditInstitutionRequest{Id: params.Id, Icon: params.Icon, Prefix: params.Prefix})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result.Id))
	}
}

func editInstitutionSchoolIcon(c *gin.Context) {
	type query struct {
		Prefix string `form:"prefix" json:"prefix" `
		Icon   string `form:"icon" json:"icon" `
		Id     int64  `form:"id" json:"id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	result, err := manager.GetCoachService().EditInstitutionSchool(c, &pb.EditInstitutionSchoolRequest{Id: params.Id, Icon: params.Icon, Prefix: params.Prefix})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result.Id))
	}
}

func editSchool(c *gin.Context) {
	type query struct {
		Name     string  `form:"name" json:"name" `
		Address  string  `form:"address" json:"address" `
		Phone    string  `form:"phone" json:"phone" `
		Start    int64   `form:"start" json:"start" `
		End      int64   `form:"end" json:"end" `
		Id       int64   `form:"id" json:"id" `
		ParentId int64   `form:"parent_id" json:"parent_id" `
		Course   []int64 `form:"course" json:"course" `
		Del      bool    `form:"del" json:"del" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Del == true && params.Id == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	adminData, ok := c.MustGet(handler.AdminData).(*jwt.AdminClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	result, err := manager.GetCoachService().EditInstitutionSchool(c, &pb.EditInstitutionSchoolRequest{Id: params.Id, Del: params.Del, Name: params.Name, Address: params.Address, Phone: params.Phone, Start: params.Start, End: params.End, Create: adminData.Uid, ParentId: params.ParentId,Course:params.Course})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result.Id))
	}
}

func courseList(c *gin.Context) {
	type Result struct {
		Name  string `json:"name" `
		Id    int64  `json:"id" `
		Max   int64  `json:"max" `
		Price int64  `json:"price" `
		Level string `json:"level" `
		Duration int64 `json:"duration" `
	}

	type query struct {
		Page     int64 `form:"page" json:"page" `
		Size     int64 `form:"size" json:"size" `
		Id       int64 `form:"id" json:"id" `
		ParentId int64 `form:"parent_id" json:"parent_id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}

	if params.Id == 0 && params.ParentId == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	courseResult, err := manager.GetCoachService().CourseList(c, &pb.InstitutionCourseRequest{Iid: params.ParentId, Id: params.Id, Page: params.Page, Size: params.Size, Status: enum.StatusAll})
	if err == nil {
		list := make([]Result, 0)
		for _, v := range courseResult.List {
			list = append(list, Result{
				Name:  v.Name,
				Id:    v.Id,
				Price: v.Price,
				Max:   v.Max,
				Level: v.Level,
				Duration:v.Duration,
			})
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func editCourse(c *gin.Context) {
	type query struct {
		Name     string `form:"name" json:"name" `
		Level    string `form:"level" json:"level" `
		Max      int64  `form:"max" json:"max" `
		Price    int64  `form:"price" json:"price" `
		Id       int64  `form:"id" json:"id" `
		ParentId int64  `form:"parent_id" json:"parent_id" `
		Duration int64  `form:"duration" json:"duration" `
		Del      bool   `form:"del" json:"del" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Del == true && params.Id == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	adminData, ok := c.MustGet(handler.AdminData).(*jwt.AdminClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	_, err := manager.GetCoachService().EditInstitutionCourse(c, &pb.EditInstitutionCourseRequest{Id: params.Id, Del: params.Del, Name: params.Name, Level: params.Level, Max: params.Max, Price: params.Price, Create: adminData.Uid, ParentId: params.ParentId,Duration:params.Duration})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
	}
}
