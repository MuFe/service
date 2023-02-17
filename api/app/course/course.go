package course

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
	server.Post("/appApi/course/levelList", levelList)
	server.Post("/appApi/course/courseList", courseList)
	server.Post("/appApi/course/courseDetail", courseDetail)
	server.Post("/appApi/course/courseDetail/v2", courseDetailV2)
	server.Post("/appApi/course/noticeDetail", noticeDetail)
	server.Post("/appApi/course/finish", handler.UserLogin, handler.TeacherCheck, finish)

}

func levelList(c *gin.Context) {
	type query struct {
		Type int64 `form:"type" json:"type" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	type Result struct {
		Id     int64  `json:"id"`
		Name   string `json:"name"`
		Number int64  `json:"number"`
	}
	if params.Type == 0 {
		params.Type = 1
	}
	result, err := manager.GetCourseService().GetCourseLevel(c, &pb.CourseLevelRequest{Type: params.Type})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	list := make([]Result, 0)
	for _, v := range result.List {
		list = append(list, Result{
			Id:     v.Id,
			Name:   v.Name,
			Number: v.Number,
		})
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
}

func courseList(c *gin.Context) {
	type Course struct {
		Cover   string         `json:"cover" `
		Bg      string         `json:"bg" `
		Id      int64          `json:"id" `
		Title   string         `json:"title" `
		Desc    string         `json:"desc" `
		Tag     []dataUtil.Tag `json:"tag" `
		FirstId int64          `json:"child_id" `
		Section int64          `json:"section" `
	}
	type query struct {
		Id int64 `form:"id" json:"id" `
	}
	result := make([]Course, 0)
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if params.Id == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseMsg("参数有误"))
		return
	}
	courseResult, err := manager.GetCourseService().GetCourse(c, &pb.CourseServiceRequest{LevelId: []int64{params.Id}, Status: enum.StatusNormal})
	if err == nil {
		for _, v := range courseResult.List {
			result = append(result, Course{
				Cover:   v.Cover,
				Id:      v.Id,
				Title:   v.Title,
				Desc:    v.Desc,
				Section: v.Section,
				Bg:      v.Bg,
				Tag:     dataUtil.ParseTag(v.Tag),
			})
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))

}

func courseDetail(c *gin.Context) {
	type query struct {
		Id int64 `form:"id" json:"id" `
	}
	type Notice struct {
		Id    int64  `json:"id"`
		Title string `json:"title"`
		Time  int64  `json:"time"`
	}

	type Video struct {
		Id      int64  `json:"id" `
		Title   string `json:"name" `
		Section int64  `json:"index" `
		Video   bool   `json:"video" `
	}

	type HomeWork struct {
		Id     int64  `json:"id" `
		Title  string `json:"name" `
		Video  bool   `json:"video" `
		InfoId int64  `json:"info" `
	}

	type Chapter struct {
		Id       int64      `json:"id" `
		Title    string     `json:"name" `
		Section  int64      `json:"index" `
		List     []Video    `json:"list" `
		HomeWork []HomeWork `json:"homework" `
	}
	type CourseDetail struct {
		NoticeList  []Notice   `json:"noticeList"`
		ChapterList []*Chapter `json:"courseList"`
	}
	result := CourseDetail{
		NoticeList:  make([]Notice, 0),
		ChapterList: make([]*Chapter, 0),
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
	courseResult, err := manager.GetCourseService().Notice(c, &pb.NoticeRequest{Id: params.Id})
	if err == nil {
		for _, v := range courseResult.List {
			result.NoticeList = append(result.NoticeList, Notice{
				Id:    v.Id,
				Title: v.Title,
				Time:  v.Time,
			})
		}
	}
	chapterResult, err := manager.GetChapterService().GetChapterWithVideo(c, &pb.ChapterVideoServiceRequest{CourseId: params.Id, Status: enum.StatusNormal})
	contentIdList := make([]int64, 0)
	chapterMap := make(map[int64]*Chapter)
	if err == nil {
		for _, v := range chapterResult.List {
			temp := &Chapter{
				Id:       v.Id,
				Title:    v.Title,
				Section:  v.Section,
				List:     make([]Video, 0),
				HomeWork: make([]HomeWork, 0),
			}
			chapterMap[v.Id] = temp
			contentIdList = append(contentIdList, v.Id)
			index := int64(0)
			for _, v1 := range v.VideoList {
				index++
				temp.List = append(temp.List, Video{
					Id:      v1.Id,
					Title:   v1.Title,
					Section: index,
					Video:   true,
				})
			}
			result.ChapterList = append(result.ChapterList, temp)
		}
	}
	homeworkResult, err := manager.GetHomeWorkService().HomeWork(c, &pb.HomeWorkRequest{ContentId: contentIdList, Status: enum.StatusNormal})
	if err == nil {
		for _, v1 := range homeworkResult.List {
			temp, ok := chapterMap[v1.ContentId]
			if ok {
				temp.HomeWork = append(temp.HomeWork, HomeWork{
					Id:     v1.Id,
					Title:  v1.Title,
					Video:  true,
					InfoId: v1.InfoId,
				})
			}
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
}

func courseDetailV2(c *gin.Context) {
	type query struct {
		Id int64 `form:"id" json:"id" `
	}

	type HomeWork struct {
		Id     int64  `json:"id" `
		Index  int64  `json:"index" `
		Level  int64  `json:"level" `
		Title  string `json:"name" `
		Cover  string `json:"cover" `
		InfoId int64  `json:"info" `
	}

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

	type Chapter struct {
		Id      int64  `json:"id" `
		Title   string `json:"name" `
		Section int64  `json:"index" `
		Cover   string `json:"cover" `
		VideoId int64  `json:"video" `
	}
	type CourseDetail struct {
		Origin      Origin     `json:"origin" `
		Bg          string     `json:"bg" `
		Cover       string     `json:"cover" `
		ChapterList []*Chapter `json:"courseList"`
	}
	result := CourseDetail{
		ChapterList: make([]*Chapter, 0),
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
	courseResult, err := manager.GetCourseService().GetCourse(c, &pb.CourseServiceRequest{Ids: []int64{params.Id}, Status: enum.StatusNormal})
	if err == nil {
		for _, v := range courseResult.List {
			result.Bg = v.Bg
			result.Cover = v.Cover
		}
	}
	originResult, err := manager.GetCourseService().GetOrigin(c, &pb.GetOriginRequest{Id: params.Id})
	if err == nil {
		result.Origin = Origin{
			Photo:     originResult.Photo,
			Name:      originResult.Name,
			Info:      originResult.Info,
			InfoTitle: originResult.InfoTitle,
			Title:     originResult.Title,
			Desc:      originResult.Desc,
			AuthTitle: originResult.Certificate,
			Auth:      make([]string, 0),
		}
		if len(originResult.Auth) > 0 {
			result.Origin.Auth = originResult.Auth
		}
	}
	chapterResult, err := manager.GetChapterService().GetChapterWithVideo(c, &pb.ChapterVideoServiceRequest{CourseId: params.Id,Size:30, Status: enum.StatusNormal})
	if err == nil {
		for _, v := range chapterResult.List {
			vId := int64(0)
			if len(v.VideoList) > 0 {
				vId = v.VideoList[0].Id
			}
			temp := &Chapter{
				Id:      v.Id,
				Title:   v.Title,
				Section: v.Section,
				Cover:   v.Cover,
				VideoId: vId,
			}
			result.ChapterList = append(result.ChapterList, temp)
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
}

func noticeDetail(c *gin.Context) {
	type query struct {
		Id int64 `form:"id" json:"id" `
	}
	type Notice struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
		Time    int64  `json:"time"`
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
	result := Notice{}
	v, err := manager.GetCourseService().NoticeDetail(c, &pb.NoticeRequest{Id: params.Id})
	if err == nil {
		result = Notice{
			Id:    v.Id,
			Title: v.Title,
			Time:  v.Time,
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func finish(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		ClassId int64 `form:"class_id" json:"class_id" `
		VideoId int64 `form:"id" json:"id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetSchoolService().EditCourseProgress(c, &pb.EditCourseProgressRequest{Uid: userData.Uid, ClassId: params.ClassId, VideoId: params.VideoId})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("更新成功"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}
