package feedback

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"mufe_service/camp/errcode"
	"mufe_service/camp/server"
	"mufe_service/manager"
	pb "mufe_service/jsonRpc"
)

func init() {
	server.Post("/adminApi/feedback", feedback)
	server.Post("/adminApi/editFeedbackStatus", editFeedbackStatus)

}

type Result struct {
	Id      int64  `json:"id"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
	Uid     int64  `json:"uid"`
	Name    string `json:"name"`
}

func feedback(c *gin.Context) {
	type query struct {
		Status     int64  `form:"status" json:"status" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	result, err := manager.GetFeedbackService().FeedbackList(c, &pb.FeedbackRequest{Status:params.Status})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	list := make([]Result, 0)
	for _, v := range result.List {
		list = append(list, Result{
			Id:      v.Id,
			Uid:     v.Uid,
			Content: v.Content,
			Time:    v.CreateTime,
			Name:    v.Name,
		})
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(list))
}

func editFeedbackStatus(c *gin.Context) {
	type query struct {
		Status     int64  `form:"status" json:"status" `
		Id     int64  `form:"id" json:"id" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_, err := manager.GetFeedbackService().EditFeedback(c, &pb.FeedbackRequest{Status:params.Status,Id:params.Id})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
}
