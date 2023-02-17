package homework

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"mufe_service/camp/enum"
	"mufe_service/camp/errcode"
	"mufe_service/camp/handler"
	"mufe_service/camp/server"
	"mufe_service/jsonRpc"
	"mufe_service/manager"
)

func init() {
	server.Post("/adminHomeWork/editHomeWork", handler.AdminLogin, homeWorkAdd)
	server.Post("/adminHomeWork/editCover", handler.AdminLogin, editCover)
	server.Post("/adminHomeWork/editVideo", handler.AdminLogin, editVideo)
	server.Post("/adminHomeWork/editContent", handler.AdminLogin, editContent)

	server.Post("/adminHomeWork/detail", handler.AdminLogin, detail)
	server.Post("/adminHomeWork/list", handler.AdminLogin, list)

}

func homeWorkAdd(c *gin.Context) {
	type query struct {
		Title string  `form:"title" json:"title" `
		Id    int64   `form:"id" json:"id" `
		Level int64   `form:"level" json:"level" `
		Tag   []int64 `form:"tag" json:"tag" `
	}
	params := query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}

	result, err := manager.GetHomeWorkService().EditHomeWork(c, &app.EditHomeWorkRequest{Title: params.Title, Tag: params.Tag, Id: params.Id, Level: params.Level})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result.Id))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func list(c *gin.Context) {
	type Query struct {
		Page int64 `form:"page" json:"page"`
		Size int64 `form:"size" json:"size"`
	}
	params := Query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	result, err := manager.GetHomeWorkService().HomeWorkInfo(c, &app.HomeWorkRequest{Page: params.Page, Size: params.Size, Status: enum.StatusNormal})
	if err == nil {
		type Result struct {
			Title string   `json:"title"`
			Cover string   `json:"cover"`
			Tag   []string `json:"tag"`
			Id    int64    `json:"id"`
		}
		list := make([]Result, 0)
		for _, v := range result.List {
			tagList := make([]string, 0)
			for _, vv := range v.Tag {
				tagList = append(tagList, vv.Title)
			}
			list = append(list, Result{
				Id:    v.Id,
				Title: v.Title,
				Cover: v.Cover,
				Tag:   tagList,
			})
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func editCover(c *gin.Context) {
	type Query struct {
		Id     int64  `form:"id" json:"id"`
		Cover  string `form:"cover" json:"cover"`
		Prefix string `form:"prefix" json:"prefix"`
	}
	params := Query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetHomeWorkService().EditHomeworkInfo(c, &app.EditHomeWorkInfoRequest{Content: params.Cover, Prefix: params.Prefix, Type: enum.HomeWork_Edit_Cover, Id: params.Id})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("修改成功"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func editVideo(c *gin.Context) {
	type Query struct {
		Id     int64  `form:"id" json:"id"`
		Video  string `form:"video" json:"video"`
		Prefix string `form:"prefix" json:"prefix"`
	}
	params := Query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetHomeWorkService().EditHomeworkInfo(c, &app.EditHomeWorkInfoRequest{Content: params.Video, Prefix: params.Prefix, Type: enum.HomeWork_Edit_Video, Id: params.Id})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("修改成功"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func editContent(c *gin.Context) {
	type Query struct {
		Id      int64  `form:"id" json:"id"`
		Content string `form:"content" json:"content"`
	}
	params := Query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetHomeWorkService().EditHomeworkInfo(c, &app.EditHomeWorkInfoRequest{Content: params.Content, Type: enum.HomeWork_Edit_Content, Id: params.Id})
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK("修改成功"))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
	}
}

func detail(c *gin.Context) {
	type Query struct {
		Id int64 `form:"id" json:"id"`
	}
	params := Query{}
	if err := c.Bind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	detailResult, err := manager.GetHomeWorkService().HomeWorkDetail(c, &app.HomeWorkDetailRequest{InfoId: params.Id})
	if err == nil {
		type Result struct {
			Title   string  `json:"title"`
			Cover   string  `json:"cover"`
			Video   string  `json:"video"`
			Content string  `json:"content"`
			Level   int64   `json:"level"`
			Tag     []int64 `json:"tag"`
		}
		tagList := make([]int64, 0)
		for _, v := range detailResult.Tag {
			tagList = append(tagList, v.Id)
		}
		result := Result{
			Title:   detailResult.Title,
			Cover:   detailResult.Cover,
			Video:   detailResult.Video,
			Content: detailResult.Content,
			Tag:     tagList,
			Level:   detailResult.Level,
		}
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(result))
	} else {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))

	}
}
