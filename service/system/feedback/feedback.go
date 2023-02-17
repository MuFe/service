package feedback

import (
	"context"
	"mufe_service/camp/service"
	pb "mufe_service/jsonRpc"
	"mufe_service/model/feedback"
)

func init() {
	nSer := &rpcServer{}
	pb.RegisterFeedbackServiceServer(service.GetRegisterRpc(), nSer)

}

type rpcServer struct {
}

func (rpc *rpcServer) AddFeedback(ctx context.Context, request *pb.FeedbackRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, feedbackmodel.AddFeedback(request.Uid,request.Content)
}

func (rpc *rpcServer) FeedbackList(ctx context.Context, request *pb.FeedbackRequest) (*pb.FeedbackResponse, error) {
	list,err:=feedbackmodel.GetFeedback(request.Status)
	if err!=nil{
		return nil,err
	}
	result:=make([]*pb.FeedbackData,0)
	for _,v:=range list{
		result=append(result,&pb.FeedbackData{
			Content:v.Content,
			Uid:v.Uid,
			Name:v.Name,
			CreateTime:v.Time,
			Id:v.ID,
		})
	}
	return &pb.FeedbackResponse{List:result},nil
}

func (rpc *rpcServer) EditFeedback(ctx context.Context, request *pb.FeedbackRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, feedbackmodel.EditFeedback(request.Id,request.Status)
}
