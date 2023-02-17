package news

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/handler"
	"mufe_service/camp/server"
	"mufe_service/camp/utils"
	"mufe_service/jsonRpc"
	"mufe_service/manager"
)

func init() {
	server.Post("/adminNews/list", handler.AdminLogin, news)
	server.Post("/adminNews/editRe", handler.AdminLogin, editRe)
	server.Post("/adminNews/editContent", handler.AdminLogin, editContent)
	server.Post("/adminNews/editCover", handler.AdminLogin, editCover)
	server.Post("/adminNews/edit", handler.AdminLogin, edit)
	server.Post("/adminNews/del", handler.AdminLogin, del)
}

type News struct {
	Id      int64  `json:"id" `
	Time    int64  `json:"time" `
	Recommend    bool  `json:"recommend" `
	Title   string `json:"title" `
	Cover   string `json:"cover" `
	Source  string `json:"source" `
	Content string `json:"content" `
}

func news(c *gin.Context) {
	type query struct {
		Id   int64 `form:"id" json:"id" `
		Page int64 `form:"page" json:"page" `
		Size int64 `form:"size" json:"size" `
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	result, err := manager.GetNewsService().GetNews(c, &app.GetNewsRequest{Id: params.Id, Page: params.Page, Size: params.Size})
	if err == nil {
		list := make([]interface{}, 0)
		for _, v := range result.List {
			list = append(list, News{
				Id:      v.Id,
				Title:   v.Title,
				Content: v.Content,
				Cover:   v.Cover,
				Time:    v.Time,
				Source:  v.Source,
				Recommend:    v.Type==enum.NewsAppHome,
			})
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(utils.CreateListResultReturn(0, list)))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func editRe(c *gin.Context){
	type query struct {
		Id   int64 `form:"id" json:"id" `
		Recommend bool `form:"recommend" json:"recommend" `
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	typeInt:=enum.NewsNone
	if params.Recommend{
		typeInt=enum.NewsAppHome
	}
	_, err := manager.GetNewsService().EditNewType(c, &app.EditNewTypeRequest{Id: params.Id, Type:typeInt})
	if err==nil{
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("修改成功"))
	}else{
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func editContent(c *gin.Context){
	type query struct {
		Id   int64 `form:"id" json:"id" `
		Content string `form:"content" json:"content" `
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetNewsService().EditContent(c, &app.EditNewRequest{Id: params.Id, Content:params.Content})
	if err==nil{
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("修改成功"))
	}else{
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func edit(c *gin.Context){
	type query struct {
		Id   int64 `form:"id" json:"id" `
		Title string `form:"title" json:"title" `
		Source string `form:"source" json:"source" `
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	result, err := manager.GetNewsService().EditNew(c, &app.EditNewRequest{Id: params.Id, Title:params.Title,Source:params.Source})
	if err==nil{
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result.Id))
	}else{
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func editCover(c *gin.Context){
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
	_, err := manager.GetNewsService().EditCover(c, &app.EditNewCoverRequest{Id: params.Id,Cover:params.Cover,Prefix:params.Prefix})
	if err==nil{
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("修改成功"))
	}else{
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func del(c *gin.Context){
	type query struct {
		Id     int64  `form:"id" json:"id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetNewsService().DelNews(c, &app.EditNewRequest{Id: params.Id})
	if err==nil{
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("删除成功"))
	}else{
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}
