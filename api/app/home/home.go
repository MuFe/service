package home

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
	server.Post("/appApi/mainPage", mainPage)
	server.Post("/appApi/courseDetail", detail)
	server.Post("/appApi/videoDetail", videoDetail)
	server.Post("/appApi/chapterDetail", chapterDetail)
	server.Post("/appApi/collection", handler.UserLogin, collection)
	server.Post("/appApi/videoHistory", history)
	server.Post("/appApi/searchHistory", searchHistory)
	server.Post("/appApi/searchHint", searchHint)
	server.Post("/appApi/search", search)
	server.Post("/appApi/tagVideo", tagVideo)
	server.Post("/appApi/share", share)
	server.Post("/appApi/newsDetail", newDetail)
	server.Post("/appApi/news", newList)
}

type News struct {
	Id      int64  `json:"id" `
	Time    int64  `json:"time" `
	Title   string `json:"title" `
	Cover   string `json:"cover" `
	Source  string `json:"source" `
	Content string `json:"content" `
}

func mainPage(c *gin.Context) {
	type BannerResult struct {
		Photo string `json:"photo" `
		Type  int64  `json:"type" `
		Url   string `json:"url" `
		Other int64  `json:"content_id" `
	}

	type Course struct {
		Cover   string `json:"cover" `
		Id      int64  `json:"id" `
		Type    int64  `json:"type" `
		Title   string `json:"title" `
		Desc    string `json:"desc" `
		Section int64  `json:"section" `
	}

	type CourseResult struct {
		List  []Course `json:"list" `
		Title string   `json:"title" `
		Type  int64    `json:"type" `
		Icon  string   `json:"icon" `
	}

	type HomeResult struct {
		Banner []BannerResult  `json:"ad" `
		News   []News          `json:"new" `
		List   []*CourseResult `json:"list" `
	}
	result := HomeResult{
		List:   make([]*CourseResult, 0),
		Banner: make([]BannerResult, 0),
		News:   make([]News, 0),
	}

	bannerResult, err := manager.GetBannerService().GetAds(c, &pb.AdServiceRequest{Status: enum.StatusNormal})
	if err == nil {
		for _, v := range bannerResult.Result {
			result.Banner = append(result.Banner, BannerResult{
				Photo: v.Photo,
				Type:  v.Type,
				Url:   v.Url,
				Other: v.LinkId,
			})
		}
	}

	dataMap := make(map[int64]*CourseResult)
	reIds := make([]int64, 0)
	userData, err := jwt.CheckUserJwt(c.GetHeader(jwt.AuthHeader))
	typeInt := enum.STUDENT_TYPE
	if err == nil && userData.Identity == enum.TEACHER_TYPE {
		typeInt = userData.Identity
	}
	infoResult, err := manager.GetRecommendService().GetRecommendInfoList(c, &pb.RecommendInfoRequest{Type: typeInt})
	if err == nil {
		xlog.Info(infoResult)
		for _, v := range infoResult.List {
			temp := &CourseResult{
				Title: v.Title,
				Icon:  v.Icon,
				Type:  v.Type,
				List:  make([]Course, 0),
			}
			dataMap[v.Id] = temp
			reIds = append(reIds, v.Id)
			result.List = append(result.List, temp)
		}
	}

	list1 := make([]int64, 0)
	list1Map := make(map[int64]int64)
	list2 := make([]int64, 0)
	list2Map := make(map[int64]int64)
	list3 := make([]int64, 0)
	list3Map := make(map[int64]int64)
	if len(reIds) > 0 {
		reResult, err := manager.GetRecommendService().GetRecommendList(c, &pb.RecommendRequest{Ids: reIds})
		if err == nil {
			for _, v := range reResult.List {
				if v.ContentType == enum.RECOMMEND_LEVEL {

				} else if v.ContentType == enum.RECOMMEND_SOURCE {
					list1 = append(list1, v.ContentId)
					list1Map[v.ContentId] = v.InfoId
				} else if v.ContentType == enum.RECOMMEND_COURSE {
					list2 = append(list2, v.ContentId)
					list2Map[v.ContentId] = v.InfoId
				} else if v.ContentType == enum.RECOMMEND_ITEM {
					list3 = append(list3, v.ContentId)
					list3Map[v.ContentId] = v.InfoId
				}

			}

			if len(list1) > 0 {
				courseResult, err := manager.GetCourseService().GetCourse(c, &pb.CourseServiceRequest{Status: enum.StatusAll, Ids: list1})
				if err == nil {

					for _, v := range courseResult.List {
						infoId, ok := list1Map[v.Id]
						if ok {
							info, ok := dataMap[infoId]
							if ok {
								info.List = append(info.List, Course{
									Cover:   v.Cover,
									Id:      v.Id,
									Title:   v.Title,
									Desc:    v.Desc,
									Section: v.Section,
									Type:    enum.RECOMMEND_SOURCE,
								})
							}
						}
					}
				}
			}
			if len(list2) > 0 {
				chapterResult, err := manager.GetChapterService().GetAdminChapter(c, &pb.ChapterServiceRequest{Page: 1, Size: 100, Ids: list2, Status: enum.StatusAll})
				if err == nil {
					for _, v := range chapterResult.List {
						infoId, ok := list2Map[v.Id]
						if ok {
							info, ok := dataMap[infoId]
							if ok {
								info.List = append(info.List, Course{
									Cover:   v.Cover,
									Id:      v.Id,
									Title:   v.Title,
									Desc:    v.Desc,
									Section: v.Section,
									Type:    enum.RECOMMEND_COURSE,
								})
							}
						}
					}
				}
			}
			if len(list3) > 0 {

				videoResult, err := manager.GetVideoService().GetAdminItem(c, &pb.GetItemRequest{Id: list3})
				if err == nil {
					for _, v := range videoResult.Data {
						infoId, ok := list3Map[v.Id]
						if ok {
							info, ok := dataMap[infoId]
							if ok {
								info.List = append(info.List, Course{
									Cover: v.Cover,
									Id:    v.Id,
									Title: v.Title,
									Type:  enum.RECOMMEND_ITEM,
								})
							}
						}
					}
				}
			}
		}
	}

	newsResult, err := manager.GetNewsService().GetNews(c, &pb.GetNewsRequest{Status: enum.StatusNormal, Type: enum.NewsAppHome})
	if err == nil {
		for _, v := range newsResult.List {
			result.News = append(result.News, News{
				Id:      v.Id,
				Title:   v.Title,
				Content: v.Content,
				Cover:   v.Cover,
				Time:    v.Time,
				Source:  v.Source,
			})
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
}

func detail(c *gin.Context) {
	type Chapter struct {
		Cover   string `json:"cover" `
		Id      int64  `json:"id" `
		Title   string `json:"title" `
		Desc    string `json:"desc" `
		Section int64  `json:"section" `
		Time    int64  `json:"time" `
		Type    int64  `json:"type" `
	}
	type Course struct {
		Cover   string    `json:"cover" `
		Id      int64     `json:"id" `
		Title   string    `json:"title" `
		Desc    string    `json:"desc" `
		Section int64     `json:"section" `
		Study   []string  `json:"study" `
		Chapter []Chapter `json:"chapter" `
	}
	type query struct {
		Id   int64 `form:"id" json:"id" `
		Page int64 `form:"page" json:"page" `
		Size int64 `form:"size" json:"size" `
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

	result := Course{}
	if params.Page == 1 {
		courseResult, err := manager.GetCourseService().GetCourse(c, &pb.CourseServiceRequest{Ids: []int64{params.Id}, Status: enum.StatusNormal})
		if err == nil {
			if len(courseResult.List) == 0 {
				c.AbortWithStatusJSON(http.StatusOK, errcode.ParseMsg("课程参数有误"))
				return
			}
			v := courseResult.List[0]
			result = Course{
				Cover:   v.Cover,
				Id:      v.Id,
				Title:   v.Title,
				Desc:    v.Desc,
				Section: v.Section,
				Study:   v.Study,
			}
		}
	}
	chapterResult, err := manager.GetChapterService().GetChapter(c, &pb.ChapterServiceRequest{Page: params.Page, Size: params.Size, CourseId: params.Id, Status: enum.StatusNormal})
	if err == nil {
		haveSchool := false
		userData, err := jwt.CheckUserJwt(c.GetHeader(jwt.AuthHeader))
		if err == nil {
			collectionResult, err := manager.GetSchoolService().MySchool(c, &pb.SchoolRequest{Uid: userData.Uid})
			if err == nil {
				haveSchool = collectionResult.Total > 0
			}
		}
		max := 2

		for k, v := range chapterResult.List {
			typeInt := enum.FreeCanSee
			if v.Price > 0 {
				typeInt = enum.PriceNoSee
			} else if k > max && !haveSchool {
				typeInt = enum.FreeSchoolNoSee
			}

			result.Chapter = append(result.Chapter, Chapter{
				Cover:   v.Cover,
				Id:      v.Id,
				Title:   v.Title,
				Desc:    v.Desc,
				Time:    v.Time,
				Section: v.Section,
				Type:    typeInt,
			})
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
}

func videoDetail(c *gin.Context) {
	type Chapter struct {
		Cover   string `json:"cover" `
		Id      int64  `json:"id" `
		Section int64  `json:"index" `
		Title   string `json:"title" `
	}
	type HomeWork struct {
		Id    int64  `json:"id" `
		Index int64  `json:"index" `
		Level int64  `json:"level" `
		Title string `json:"name" `
		Cover string `json:"cover" `
	}
	type Result struct {
		Cover      string         `json:"cover" `
		Id         int64          `json:"id" `
		Title      string         `json:"title" `
		Tags       []dataUtil.Tag `json:"tag" `
		Collection bool           `json:"collection" `
		PlayUrl    string         `json:"play_url" `
		DownUrl    string         `json:"down_url" `
		Content    string         `json:"content" `
		List       []Chapter      `json:"list" `
		Hot        []Chapter      `json:"hot" `
		Plan       string         `json:"plan" `
		HomeWork   []HomeWork     `json:"homework" `
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

	result := &Result{}
	videoResult, err := manager.GetVideoService().GetVideo(c, &pb.VideoRequest{VideoId: params.Id})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if len(videoResult.VideoList) == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseMsg("参数有误"))
		return
	}
	collection := false
	userData, err := jwt.CheckUserJwt(c.GetHeader(jwt.AuthHeader))
	if err == nil {
		collectionResult, err := manager.GetCollectionService().IsCollection(c, &pb.CollectionServiceRequest{ContentId: params.Id, Type: 3, Uid: userData.Uid})
		if err == nil {
			collection = collectionResult.Collection
		}
	}
	v := videoResult.VideoList[0]

	result = &Result{
		Cover:      v.Cover,
		Id:         params.Id,
		Title:      v.Title,
		Collection: collection,
		Tags:       dataUtil.ParseTag(v.Tag),
		PlayUrl:    v.Url,
		DownUrl:    v.DownUrl,
		Content:    v.Content,
		List:       make([]Chapter, 0),
		Hot:       make([]Chapter, 0),
		HomeWork:   make([]HomeWork, 0),
	}
	contentIdList := make([]int64, 0)
	chapterResult, err := manager.GetChapterService().GetChapterWithVideo(c, &pb.ChapterVideoServiceRequest{Page: 1, Size: 10, ChapterId: []int64{v.ChapterId}, Status: enum.StatusNormal})
	chapterMap := make(map[int64]*Result)
	if err == nil {
		for _, v := range chapterResult.List {
			result.Plan = v.Plan
			chapterMap[v.Id] = result
			contentIdList = append(contentIdList, v.Id)
			index := int64(0)
			for _, v1 := range v.VideoList {
				index++
				result.List = append(result.List, Chapter{
					Id:      v1.Id,
					Title:   v1.Title,
					Section: index,
					Cover:   v1.Cover,
				})
				result.Hot = append(result.Hot, Chapter{
					Id:      v1.Id,
					Title:   v1.Title,
					Section: index,
					Cover:   v1.Cover,
				})
			}
		}
	}
	homeworkResult, err := manager.GetHomeWorkService().HomeWork(c, &pb.HomeWorkRequest{ContentId: contentIdList, Status: enum.StatusNormal})
	if err == nil {
		for _, v1 := range homeworkResult.List {
			temp, ok := chapterMap[v1.ContentId]
			if ok {
				temp.HomeWork = append(temp.HomeWork, HomeWork{
					Id:    v1.Id,
					Title: v1.Title,
					Cover: v1.Cover,
					Level: v1.Level,
					Index: v1.Index,
				})
			}
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
}

func collection(c *gin.Context) {
	userData, ok := c.MustGet(handler.UserData).(*jwt.UserClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusOK, errcode.CommErrorNotLogin)
		return
	}
	type query struct {
		Id         int64 `form:"id" json:"id" `
		Collection bool  `form:"collection" json:"collection" `
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
	_, err := manager.GetCollectionService().EditCollection(c, &pb.CollectionServiceRequest{ContentId: params.Id, Type: 3, Del: !params.Collection, Uid: userData.Uid})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(true))
}

func history(c *gin.Context) {
	type query struct {
		Id int64 `form:"id" json:"id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}

	userData, err := jwt.CheckUserJwt(c.GetHeader(jwt.AuthHeader))
	if err == nil {
		_, err = manager.GetVideoService().AddHistoryVideo(c, &pb.AddVideoHistoryRequest{Uid: userData.Uid, Id: params.Id})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		}
	}

	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
}

func searchHistory(c *gin.Context) {
	type Data struct {
		Content string `json:"content" `
		Hot     bool   `json:"hot" `
	}
	type Result struct {
		List []Data `json:"history" `
		Hot  []Data `json:"hot" `
	}
	list := make([]Data, 0)
	hot := make([]Data, 0)
	uid := int64(0)
	userData, err := jwt.CheckUserJwt(c.GetHeader(jwt.AuthHeader))
	if err == nil {
		uid = userData.Uid
	}
	result, err := manager.GetSearchService().GetSearchHistory(c, &pb.SearchRequest{Uid: uid})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	for _, v := range result.List {
		list = append(list, Data{
			Content: v.Content,
			Hot:     v.TodayNumber > 10,
		})
	}
	for _, v := range result.Hot {
		hot = append(hot, Data{
			Content: v.Content,
			Hot:     v.TodayNumber > 10,
		})
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(Result{
		List: list,
		Hot:  hot,
	}))
}

func searchHint(c *gin.Context) {
	type query struct {
		Content string `form:"content"`
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	type Result struct {
		Content string `json:"content" `
	}
	list := make([]Result, 0)

	result, err := manager.GetSearchService().SearchHint(c, &pb.SearchRequest{Content: params.Content})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	for _, v := range result.List {
		list = append(list, Result{
			Content: v.Content,
		})
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
}

func search(c *gin.Context) {
	type Video struct {
		Cover string `json:"cover" `
		Id    int64  `json:"id" `
		Title string `json:"title" `
	}
	type Result struct {
		VideoList []Video `json:"items" `
	}

	type query struct {
		Content string `form:"content"`
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	returnResult := Result{
		VideoList: make([]Video, 0),
	}
	videoResult, err := manager.GetVideoService().GetVideo(c, &pb.VideoRequest{Page: 1, Pagesize: 50, Key: params.Content})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	uid := int64(0)
	userData, err := jwt.CheckUserJwt(c.GetHeader(jwt.AuthHeader))
	if err == nil {
		uid = userData.Uid
	}
	_, _ = manager.GetSearchService().AddSearch(c, &pb.SearchRequest{Content: params.Content, Uid: uid})
	for _, temp := range videoResult.VideoList {
		returnResult.VideoList = append(returnResult.VideoList, Video{
			Id:    temp.Id,
			Title: temp.Title,
			Cover: temp.Cover,
		})
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(returnResult))
}

func tagVideo(c *gin.Context) {

	type Video struct {
		Cover string `json:"cover" `
		Id    int64  `json:"id" `
		Title string `json:"title" `
		Time  int64  `json:"time" `
		Users int64  `json:"users" `
	}
	type Result struct {
		List    []Video `json:"list" `
		Title   string  `json:"title" `
		Content string  `json:"content" `
		Cover   string  `json:"cover" `
	}

	type query struct {
		Id   int64  `form:"id" json:"id" `
		Name string `form:"name" json:"name" `
		Sort int64  `form:"sort" json:"sort" `
		Page int64  `form:"page" json:"page" `
		Size int64  `form:"size" json:"size" `
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}

	tagResult, err := manager.GetCourseService().TagList(c, &pb.TagRequest{Id: params.Id, Status: enum.StatusNormal})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	if len(tagResult.List) == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
		return
	}
	temp := tagResult.List[0]
	result := Result{
		Title:   temp.Title,
		Content: temp.Content,
		Cover:   temp.Cover,
		List:    make([]Video, 0),
	}
	videoResult, err := manager.GetVideoService().GetVideo(c, &pb.VideoRequest{TagId: params.Id, Sort: params.Sort, Page: params.Page, Pagesize: params.Size})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}

	for _, temp := range videoResult.VideoList {
		result.List = append(result.List, Video{
			Id:    temp.Id,
			Time:  temp.Duration,
			Title: temp.Title,
			Cover: temp.Cover,
			Users: temp.User,
		})
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
}

func share(c *gin.Context) {
	type Result struct {
		Url     string `json:"url" `
		Icon    string `json:"icon" `
		Content string `json:"content" `
		Title   string `json:"title" `
	}
	result := Result{}
	resultInfo, err := manager.GetWebService().GetWebInfo(c, &pb.GetWebInfoRequest{List: []int64{enum.WebShareTitle, enum.WebShareContent, enum.WebShareIcon, enum.WebShareUrl}})
	if err == nil {
		result.Title = resultInfo.Content[enum.WebShareTitle]
		result.Content = resultInfo.Content[enum.WebShareContent]
		result.Icon = resultInfo.Content[enum.WebShareIcon]
		result.Url = resultInfo.Content[enum.WebShareUrl]
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))

}

func newDetail(c *gin.Context) {
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
	result := News{}
	resultInfo, err := manager.GetNewsService().GetNews(c, &pb.GetNewsRequest{Id: params.Id, Status: enum.StatusNormal})
	if err == nil {
		if len(resultInfo.List) == 0 {
			c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
			return
		} else {
			v := resultInfo.List[0]
			result = News{
				Id:      v.Id,
				Title:   v.Title,
				Content: v.Content,
				Cover:   v.Cover,
				Time:    v.Time,
				Source:  v.Source,
			}
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
}

func newList(c *gin.Context) {
	type query struct {
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
		params.Size = 20
	}
	result := make([]News, 0)
	resultInfo, err := manager.GetNewsService().GetNews(c, &pb.GetNewsRequest{Page: params.Page, Size: params.Size})
	if err == nil {
		for _, v := range resultInfo.List {
			result = append(result, News{
				Id:     v.Id,
				Title:  v.Title,
				Cover:  v.Cover,
				Time:   v.Time,
				Source: v.Source,
			})
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
}

func chapterDetail(c *gin.Context) {
	type Video struct {
		Id      int64          `json:"id" `
		Cover   string         `json:"cover" `
		Title   string         `json:"title" `
		Time    int64          `json:"time" `
		PlayUrl string         `json:"play_url" `
		DownUrl string         `json:"down_url" `
		Tags    []dataUtil.Tag `json:"tag" `

		Canplay bool `json:"can_play" `
	}

	type Result struct {
		Id    int64   `json:"id" `
		Title string  `json:"title" `
		Desc  string  `json:"desc" `
		Video []Video `json:"video" `
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

	videoResult, err := manager.GetChapterService().GetChapterWithVideo(c, &pb.ChapterVideoServiceRequest{ChapterId: []int64{params.Id}, Status: enum.StatusNormal})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	} else if len(videoResult.List) == 0 {
		c.AbortWithStatusJSON(http.StatusOK, errcode.HttpErrorWringParam)
	}
	v := videoResult.List[0]
	result := Result{
		Id:    v.Id,
		Title: v.Title,
		Desc:  v.Desc,
		Video: make([]Video, 0),
	}

	for _, temp := range v.VideoList {
		result.Video = append(result.Video, Video{
			Id:      temp.Id,
			Time:    temp.Duration,
			Title:   temp.Title,
			Cover:   temp.Cover,
			PlayUrl: temp.Url,
			DownUrl: temp.DownUrl,
			Tags:    dataUtil.ParseTag(v.Tag),
			Canplay: true,
		})
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
}
