package course

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/handler"
	"mufe_service/camp/jwt"
	"mufe_service/camp/server"
	"mufe_service/camp/utils"
	"mufe_service/camp/xlog"
	pb "mufe_service/jsonRpc"
	"mufe_service/manager"
	"strconv"
	"time"
)

func init() {
	server.Post("/adminCourse/token", handler.AdminLogin, getToken)

	server.Post("/adminCourse/levelList", handler.AdminLogin, courseLevel)
	server.Post("/adminCourse/sourceList", handler.AdminLogin, sourceList)
	server.Post("/adminCourse/courseList", handler.AdminLogin, courseList)
	server.Post("/adminCourse/itemList", handler.AdminLogin, itemList)

	server.Post("/adminCourse/sourceDetail", handler.AdminLogin, sourceDetail)
	server.Post("/adminCourse/courseDetail", handler.AdminLogin, courseDetail)
	server.Post("/adminCourse/itemDetail", handler.AdminLogin, itemDetail)

	server.Post("/adminCourse/tagList", handler.AdminLogin, tagList)
	server.Post("/adminCourse/tagDetail", handler.AdminLogin, tagDetail)
	server.Post("/adminCourse/editTag", handler.AdminLogin, editTag)
	server.Post("/adminCourse/editTagCover", handler.AdminLogin, editTagCover)

	server.Post("/adminCourse/editLevel", handler.AdminLogin, editLevel)
	server.Post("/adminCourse/editSource", handler.AdminLogin, editSource)
	server.Post("/adminCourse/editCourse", handler.AdminLogin, editCourse)
	server.Post("/adminCourse/editItem", handler.AdminLogin, editItem)

	server.Post("/adminCourse/editSourceCover", handler.AdminLogin, editCourseCover)
	server.Post("/adminCourse/editOriginAuth", handler.AdminLogin, editCourseOriginAuth)
	server.Post("/adminCourse/editCourseInfo", handler.AdminLogin, editCourseInfo)
	server.Post("/adminCourse/editItemInfo", handler.AdminLogin, editItemInfo)

	server.Post("/adminCourse/addAppContent", handler.AdminLogin, addAppContent)
	server.Post("/adminCourse/delAppContent", handler.AdminLogin, delAppContent)
	server.Post("/adminCourse/appContentList", handler.AdminLogin, appContent)

	server.Post("/adminCourse/homeworkList", handler.AdminLogin, homeworkList)
	server.Post("/adminCourse/unBindHomeWork", handler.AdminLogin, unBind)

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
		Name    []string `form:"names" json:"names" `
		IsVideo bool     `form:"video" json:"video" `
		IsBase  bool     `form:"base64" json:"base64" `
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	list := make([]string, 0)
	for _, v := range params.Name {
		filenameWithSuffix := path.Base(v)
		fileSuffix := path.Ext(filenameWithSuffix) //获取文件后缀
		encodeString := utils.MD5(v+strconv.FormatInt(time.Now().Unix(), 10)) + fileSuffix
		if params.IsBase {
			encodeString = base64.StdEncoding.EncodeToString([]byte(encodeString))
		}
		list = append(list, encodeString)
	}
	osStr := os.Getenv("IMG_BUCKET")
	prefix := os.Getenv("IMG_PREFIX")
	if params.IsVideo {
		osStr = os.Getenv("VIDEO_BUCKET")
		prefix = os.Getenv("VIDEO_PREFIX")
	}
	result, err := manager.GetQiniuService().GetToken(c, &pb.QiniuServiceRequest{Bucket: osStr})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(&QiniuInfo{Token: result.Token, Host: result.UploadHost, BaseHost: result.Base64UploadHost, Keys: list, Prefix: prefix}))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func sourceList(c *gin.Context) {

	type Course struct {
		Cover string `json:"cover" `
		Id    int64  `json:"id" `
		Title string `json:"title" `
	}

	type query struct {
		Page  int64 `form:"page" json:"page" `
		Size  int64 `form:"size" json:"size" `
		Level int64 `form:"level" json:"level" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Level == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseMsg("参数有误"))
		return
	}
	courseResult, err := manager.GetCourseService().GetCourse(c, &pb.CourseServiceRequest{Status: enum.StatusNormal, LevelId: []int64{params.Level}})
	if err == nil {
		data := make([]interface{}, 0)
		for _, v := range courseResult.List {
			data = append(data, Course{
				Cover: v.Cover,
				Id:    v.Id,
				Title: v.Title,
			})
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(data))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func courseList(c *gin.Context) {
	type Chapter struct {
		Cover string `json:"cover" `
		Id    int64  `json:"id" `
		Title string `json:"title" `
	}

	type query struct {
		Page   int64 `form:"page" json:"page" `
		Size   int64 `form:"size" json:"size" `
		Source int64 `form:"source" json:"source" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Source == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseMsg("参数有误"))
		return
	}

	chapterResult, err := manager.GetChapterService().GetAdminChapter(c, &pb.ChapterServiceRequest{Page: 1, Size: 100, CourseId: params.Source, Status: enum.StatusAll})
	if err == nil {
		list := make([]Chapter, 0)
		for _, v := range chapterResult.List {
			list = append(list, Chapter{
				Cover: v.Cover,
				Id:    v.Id,
				Title: v.Title,
			})
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func itemList(c *gin.Context) {
	type Video struct {
		Cover string `json:"cover" `
		Id    int64  `json:"id" `
		Title string `json:"title" `
	}
	type query struct {
		Page   int64 `form:"page" json:"page" `
		Size   int64 `form:"size" json:"size" `
		Course int64 `form:"course" json:"course" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Course == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseMsg("参数有误"))
		return
	}

	courseResult, err := manager.GetChapterService().GetAdminChapterWithVideo(c, &pb.ChapterVideoServiceRequest{ChapterId: []int64{params.Course}, Page: params.Page, Size: params.Size, Status: enum.StatusAll})
	if err == nil {
		if len(courseResult.List) == 0 {
			c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		}
		v := courseResult.List[0]
		list := make([]Video, 0)
		for _, v := range v.VideoList {
			list = append(list, Video{
				Cover: v.Cover,
				Id:    v.Id,
				Title: v.Title,
			})
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func sourceDetail(c *gin.Context) {
	type Origin struct {
		Photo     string   `json:"photo" `
		Name      string   `json:"name" `
		Title     string   `json:"title" `
		Desc      string   `json:"desc" `
		Info      string   `json:"info" `
		InfoTitle string   `json:"info_title" `
		AuthTitle string   `json:"auth_title" `
		Auth      []string `json:"auth" `
	}

	type Course struct {
		Cover  string `json:"cover" `
		Id     int64  `json:"id" `
		Title  string `json:"title" `
		Bg     string `json:"bg" `
		Origin Origin `json:"origin" `
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
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseMsg("参数有误"))
		return
	}

	courseResult, err := manager.GetCourseService().GetAdminCourseDetail(c, &pb.CourseServiceRequest{Ids: []int64{params.Id}, Status: enum.StatusAll})
	if err == nil {
		if len(courseResult.List) == 0 {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseMsg("参数有误"))
			return
		}
		v := courseResult.List[0]
		tags := make([]int64, 0)
		for _, v := range v.Tag {
			tags = append(tags, v.Id)
		}
		data := Course{
			Cover: v.Cover,
			Id:    v.Id,
			Title: v.Title,
			Bg:    v.Bg,
			Origin: Origin{
				Photo:     v.Data.Photo,
				Name:      v.Data.Name,
				Info:      v.Data.Info,
				InfoTitle: v.Data.InfoTitle,
				Title:     v.Data.Title,
				Desc:      v.Data.Desc,
				AuthTitle: v.Data.Certificate,
				Auth:      make([]string, 0),
			},
		}
		if len(v.Data.Auth) > 0 {
			data.Origin.Auth = v.Data.Auth
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(data))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func courseDetail(c *gin.Context) {
	type Course struct {
		Cover    string  `json:"cover" `
		Id       int64   `json:"id" `
		Title    string  `json:"title" `
		Desc     string  `json:"desc" `
		Video    string  `json:"video" `
		Plan     string  `json:"plan" `
		HomeWork []int64 `json:"homework" `
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
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseMsg("参数有误"))
		return
	}

	chapterResult, err := manager.GetChapterService().GetAdminChapterDetail(c, &pb.ChapterServiceRequest{Ids: []int64{params.Id}, Status: enum.StatusAll})
	if err == nil {
		v := chapterResult.List[0]
		result := Course{
			Cover:    v.Cover,
			Id:       v.Id,
			Title:    v.Title,
			Desc:     v.Desc,
			Plan:     v.Plan,
			Video:    v.Video,
			HomeWork: v.Homework,
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func itemDetail(c *gin.Context) {
	type Result struct {
		Cover   string  `json:"cover" `
		Video   string  `json:"video" `
		Id      int64   `json:"id" `
		Title   string  `json:"title" `
		Tags    []int64 `json:"tags" `
		Content string  `json:"content" `
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
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseMsg("参数有误"))
		return
	}

	result := Result{}
	videoResult, err := manager.GetVideoService().GetAdminItem(c, &pb.GetItemRequest{Id: []int64{params.Id}})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}

	v := videoResult.Data[0]
	tags := make([]int64, 0)
	for _, v := range v.Tag {
		tags = append(tags, v.Id)
	}
	result = Result{
		Cover:   v.Cover,
		Id:      v.Id,
		Title:   v.Title,
		Tags:    tags,
		Content: v.Content,
		Video:   v.Url,
	}

	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
}

func editLevel(c *gin.Context) {
	type query struct {
		Name string `form:"title" json:"title" `
		Id   int64  `form:"id" json:"id" `
		Type int64  `form:"type" json:"type" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Type == 0 {
		params.Type = 1
	}
	id, err := manager.GetCourseService().EditCourseLevel(c, &pb.EditCourseLevelRequest{Id: params.Id, Name: params.Name, Type: params.Type})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(id.Id))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func editSource(c *gin.Context) {
	type query struct {
		Id              int64  `form:"id" json:"id" `
		Title           string `form:"title" json:"title" `
		OriginTitle     string `form:"origin_title" json:"origin_title" `
		OriginDesc      string `form:"origin_desc" json:"origin_desc"`
		OriginName      string `form:"origin_name" json:"origin_name"`
		OriginInfoTitle string `form:"origin_info_title" json:"origin_info_title"`
		OriginInfo      string `form:"origin_info" json:"origin_info"`
		Certificate     string `form:"certificate" json:"certificate"`
		Level           int64  `form:"level_id" json:"level_id" `
		Type            int64  `form:"type" json:"type" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	userData, ok := c.MustGet(handler.AdminData).(*jwt.AdminClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	result, err := manager.GetCourseService().EditCourse(c, &pb.EditCourseRequest{
		Title:           params.Title,
		OriginTitle:     params.OriginTitle,
		OriginDesc:      params.OriginDesc,
		OriginName:      params.OriginName,
		OriginInfoTitle: params.OriginInfoTitle,
		OriginInfo:      params.OriginInfo,
		Certificate:     params.Certificate,
		Level:           params.Level,
		Id:              params.Id,
		CreateBy:        userData.Uid,
	})
	if err == nil {
		if params.Type==2{
			cId:=int64(0)
			if params.Id!=0{
				chapterResult, err := manager.GetChapterService().GetAdminChapter(c, &pb.ChapterServiceRequest{Page: 1, Size: 100, CourseId: params.Id, Status: enum.StatusAll})
				if err==nil&&len(chapterResult.List)>0{
					cId=chapterResult.List[0].Id
				}
			}
			_, err := manager.GetChapterService().EditChapter(c, &pb.EditChapterRequest{Title: params.Title, Desc:"", Source:result.Id,Id:cId, CreateBy: userData.Uid})
			if err==nil{
				c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result.Id))
			}else{
				c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
			}
		}else{
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result.Id))
		}
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func editTagCover(c *gin.Context) {
	type query struct {
		Id     int64  `form:"id" json:"id" `
		Cover  string `form:"cover" json:"cover" `
		Prefix string `form:"prefix" json:"prefix" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetCourseService().EditTag(c, &pb.EditTagRequest{Id: params.Id, Cover: params.Cover, Prefix: params.Prefix, Type: enum.EditTagCover})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func editCourseCover(c *gin.Context) {
	type query struct {
		Id     int64  `form:"id" json:"id" `
		Cover  string `form:"cover" json:"cover" `
		Prefix string `form:"prefix" json:"prefix" `
		Type int64 `form:"type" json:"type" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetCourseService().EditCourseCover(c, &pb.EditCourseRequest{Cover: params.Cover, Id: params.Id, Prefix: params.Prefix,Type:params.Type})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func editCourseOriginAuth(c *gin.Context) {
	type Data struct {
		Cover  string `form:"cover" json:"cover" `
		Prefix string `form:"prefix" json:"prefix" `
	}
	type query struct {
		Id   int64  `form:"id" json:"id" `
		List []Data `form:"list" json:"list" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	list := make([]*pb.EditCourseOriginData, 0)
	for _, v := range params.List {
		list = append(list, &pb.EditCourseOriginData{
			Cover:  v.Cover,
			Prefix: v.Prefix,
		})
	}
	xlog.Info(params.List)
	_, err := manager.GetCourseService().EditCourseCover(c, &pb.EditCourseRequest{Type: 4, Id: params.Id, AuthList: list})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func editCourseInfo(c *gin.Context) {
	type query struct {
		Id      int64  `form:"id" json:"id" `
		Type    int64  `form:"type" json:"type" `
		Content string `form:"content" json:"content" `
		Prefix  string `form:"prefix" json:"prefix" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetChapterService().EditChapterCoverInfo(c, &pb.EditChapterRequest{InfoContent: params.Content, Id: params.Id, Prefix: params.Prefix, Type: params.Type})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func editItemInfo(c *gin.Context) {
	type query struct {
		Id      int64  `form:"id" json:"id" `
		Type    int64  `form:"type" json:"type" `
		Content string `form:"content" json:"content" `
		Prefix  string `form:"prefix" json:"prefix" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	list := make([]*pb.AddVideoData, 0)

	if params.Type == enum.EDIT_VIDEO {
		videoResult, err := manager.GetQiniuService().GetVideoInfo(c, &pb.GetVideoInfoRequest{
			Url:    params.Content,
			Prefix: params.Prefix,
		})
		if err == nil {
			list = append(list, &pb.AddVideoData{
				Cover:         videoResult.Cover,
				CoverPrefix:   videoResult.CoverPrefix,
				Duration:      videoResult.Duration,
				Url:           params.Content,
				UrlPrefix:     params.Prefix,
				DownUrl:       params.Content,
				DownUrlPrefix: params.Prefix,
			})
		}
	}
	_, err := manager.GetVideoService().EditItemInfo(c, &pb.EditItemInfoRequest{Content: params.Content, Id: params.Id, Prefix: params.Prefix, Type: params.Type, List: list})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func editCourse(c *gin.Context) {
	type query struct {
		Title    string  `form:"title" json:"title" `
		Desc     string  `form:"desc" json:"desc" `
		Id       int64   `form:"id" json:"id" `
		Homework []int64 `form:"homework" json:"homework" `
		SourceId int64   `form:"source_id" json:"source_id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	userData, ok := c.MustGet(handler.AdminData).(*jwt.AdminClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	result, err := manager.GetChapterService().EditChapter(c, &pb.EditChapterRequest{Title: params.Title, Desc: params.Desc, Id: params.Id, Source: params.SourceId, HomeWork: params.Homework, CreateBy: userData.Uid})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result.Id))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func editItem(c *gin.Context) {
	type query struct {
		Title    string  `form:"title" json:"title" `
		Id       int64   `form:"id" json:"id" `
		CourseId int64   `form:"course_id" json:"course_id" `
		Tag      []int64 `form:"tag" json:"tag" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	result, err := manager.GetVideoService().EditItem(c, &pb.EditItemRequest{Title: params.Title, Id: params.Id, ChapterId: params.CourseId, TagId: params.Tag})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result.Id))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func courseLevel(c *gin.Context) {
	type Level struct {
		Id    int64  `json:"id" `
		Title string `json:"title" `
	}
	type query struct {
		Type int64 `form:"type" json:"type" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	courseResult, err := manager.GetCourseService().GetCourseLevel(c, &pb.CourseLevelRequest{Type: params.Type})
	if err == nil {
		data := make([]Level, 0)
		for _, v := range courseResult.List {
			data = append(data, Level{
				Id:    v.Id,
				Title: v.Name,
			})
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(data))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func tagList(c *gin.Context) {
	type Tag struct {
		Id      int64  `json:"id" `
		Cover   string `json:"cover" `
		Title   string `json:"title" `
		Content string `json:"content" `
	}
	type query struct {
		Status int64 `form:"status" json:"status" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Status == 0 {
		params.Status = enum.StatusAll
	}
	videoResult, err := manager.GetCourseService().TagList(c, &pb.TagRequest{Status: params.Status})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	list := make([]interface{}, 0)
	for _, temp := range videoResult.List {
		list = append(list, Tag{
			Id:      temp.Id,
			Title:   temp.Title,
			Cover:   temp.Cover,
			Content: temp.Content,
		})
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
}

func tagDetail(c *gin.Context) {
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
	type Tag struct {
		Id      int64  `json:"id" `
		Cover   string `json:"cover" `
		Title   string `json:"title" `
		Content string `json:"content" `
	}

	result, err := manager.GetCourseService().TagList(c, &pb.TagRequest{Id: params.Id, Status: enum.StatusAll})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if len(result.List) == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	temp := result.List[0]
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(Tag{
		Id:      temp.Id,
		Title:   temp.Title,
		Cover:   temp.Cover,
		Content: temp.Content,
	}))
}

func editTag(c *gin.Context) {
	type query struct {
		Title   string `form:"title" json:"title" `
		Content string `form:"content" json:"content" `
		Id      int64  `form:"id" json:"id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	result, err := manager.GetCourseService().EditTag(c, &pb.EditTagRequest{Id: params.Id, Title: params.Title, Content: params.Content, Type: enum.AddTag})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result.Id))
}

func addAppContent(c *gin.Context) {
	type query struct {
		Id          int64 `form:"id" json:"id" `
		ContentType int64 `form:"content_type" json:"content_type" `
		Type        int64 `form:"type" json:"type" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetRecommendService().EditRecommend(c, &pb.EditRecommendRequest{ContentId: params.Id, ContentType: params.ContentType, Type: params.Type})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
}

func delAppContent(c *gin.Context) {
	type query struct {
		Id int64 `form:"id" json:"id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetRecommendService().EditRecommend(c, &pb.EditRecommendRequest{Id: params.Id, IsDel: true})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
}

func appContent(c *gin.Context) {
	type query struct {
		Type int64 `form:"type" json:"type" `
		Page int64 `form:"page" json:"page" `
		Size int64 `form:"size" json:"size" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Size == 0 {
		params.Size = 10
	}
	result, err := manager.GetRecommendService().GetAdminRecommendList(c, &pb.RecommendRequest{Page: params.Page, Size: params.Size, Ids: []int64{params.Type}})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	list1 := make([]int64, 0)
	list2 := make([]int64, 0)
	list3 := make([]int64, 0)
	type Result struct {
		Id        int64  `json:"id" `
		Cover     string `json:"cover" `
		Title     string `json:"title" `
		ContentId int64  `json:"content_id" `
		Type      int64  `json:"type" `
	}
	resultList := make([]*Result, 0)
	resultMap := make(map[int64]*Result)
	for _, v := range result.List {
		temp := &Result{
			Id:        v.Id,
			Type:      v.InfoId,
			ContentId: v.ContentId,
		}
		if v.ContentType == enum.RECOMMEND_LEVEL {

		} else if v.ContentType == enum.RECOMMEND_SOURCE {
			list1 = append(list1, v.ContentId)
			resultList = append(resultList, temp)
			resultMap[v.ContentId] = temp
		} else if v.ContentType == enum.RECOMMEND_COURSE {
			list2 = append(list2, v.ContentId)
			resultList = append(resultList, temp)
			resultMap[v.ContentId] = temp
		} else if v.ContentType == enum.RECOMMEND_ITEM {
			list3 = append(list3, v.ContentId)
			resultList = append(resultList, temp)
			resultMap[v.ContentId] = temp
		}
	}
	if len(list1) > 0 {
		courseResult, err := manager.GetCourseService().GetAdminCourse(c, &pb.CourseServiceRequest{Status: enum.StatusAll, Ids: list1})
		if err == nil {
			for _, v := range courseResult.List {
				temp, ok := resultMap[v.Id]
				if ok {
					temp.Title = v.Title
					temp.Cover = v.Cover
				}
			}
		}
	}
	if len(list2) > 0 {
		chapterResult, err := manager.GetChapterService().GetAdminChapter(c, &pb.ChapterServiceRequest{Page: 1, Size: 100, Ids: list2, Status: enum.StatusAll})
		if err == nil {
			for _, v := range chapterResult.List {
				temp, ok := resultMap[v.Id]
				if ok {
					temp.Title = v.Title
					temp.Cover = v.Cover
				}
			}
		}
	}
	if len(list3) > 0 {
		videoResult, err := manager.GetVideoService().GetAdminItem(c, &pb.GetItemRequest{Id: list3})
		if err == nil {
			for _, v := range videoResult.Data {
				temp, ok := resultMap[v.Id]
				if ok {
					temp.Title = v.Title
					temp.Cover = v.Cover
				}
			}
		}
	}

	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(resultList))
}

func homeworkList(c *gin.Context) {
	type query struct {
		Course int64 `form:"course" json:"course" `
		Page   int64 `form:"page" json:"page" `
		Size   int64 `form:"size" json:"size" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Course == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Size == 0 {
		params.Size = 10
	}
	type HomeWork struct {
		Id    int64    `json:"id" `
		Title string   `json:"title" `
		Cover string   `json:"cover" `
		Tag   []string `json:"tag" `
	}
	list := make([]HomeWork, 0)
	homeworkResult, err := manager.GetHomeWorkService().HomeWork(c, &pb.HomeWorkRequest{ContentId: []int64{params.Course}, Status: enum.StatusNormal, Page: params.Page, Size: params.Size})
	if err == nil {
		for _, v1 := range homeworkResult.List {
			tags := make([]string, 0)
			for _, v := range v1.Tag {
				tags = append(tags, v.Title)
			}
			list = append(list, HomeWork{
				Id:    v1.Id,
				Title: v1.Title,
				Cover: v1.Cover,
				Tag:   tags,
			})
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
}

func unBind(c *gin.Context) {
	type query struct {
		Id int64 `form:"id" json:"id" `
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetHomeWorkService().UnbindHomeWork(c, &pb.UnBindHomeWorkRequest{Id: params.Id})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
}
