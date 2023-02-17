package system

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"mufe_service/camp/errcode"
	"mufe_service/camp/jwt"
	"mufe_service/camp/server"
	"mufe_service/manager"
	pb "mufe_service/jsonRpc"
)

func init() {
	server.Post("/appApi/update", update)
	server.Post("/appApi/feedback", feedback)
	server.Get("/appApi/about", about)
	server.Get("/appApi/privacy", privacy)
}

func update(c *gin.Context) {
	type Result struct {
		Version     string `json:"version" `
		VersionCode int64  `json:"code" `
		Android     string `json:"android" `
		Ios         string `json:"ios" `
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(Result{
		Version:     "1.0",
		VersionCode: 1,
		Android:     "",
		Ios:         "",
	}))
}

func feedback(c *gin.Context) {
	userData, err := jwt.CheckUserJwt(c.GetHeader(jwt.AuthHeader))
	uid:=int64(0)
	if err == nil {
		uid=userData.Uid
	}
	type query struct {
		Content   string `form:"content" json:"content" `
	}
	params := query{}
	if err := c.ShouldBind(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	_,err=manager.GetFeedbackService().AddFeedback(c,&pb.FeedbackRequest{Content:params.Content,Uid:uid})
	if  err != nil {
		c.AbortWithStatusJSON(http.StatusOK, errcode.ParseError(err))
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, errcode.ParseOK(""))
}

func about(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "https://www.sgsportsgroup.com/agreement.html?v3")
}

func privacy(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "https://www.sgsportsgroup.com/privacy.html?v3")
}


