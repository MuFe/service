package football

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"mufe_service/camp/service"
	"mufe_service/camp/xlog"
	pb "mufe_service/jsonRpc"
	"mufe_service/model/foot"
	"strconv"
	"time"
)

type rpcServer struct {
}



func init() {
	pb.RegisterFootballServiceServer(service.GetRegisterRpc(), &rpcServer{})
}

func (rpc *rpcServer) GetData(ctx context.Context, request *pb.GetFootDataRequest) (result *pb.GetDataResponse, err error) {
	list:=make([]*pb.GetDataResult,0)
	for _,v:=range request.List{
		result,err:=send(v.Mac, v.Date)
		if err==nil{
			list=append(list,result...)
		}
	}

	return &pb.GetDataResponse{List:list},nil
}

func (rpc *rpcServer) GetFoot(ctx context.Context, request *pb.GetFootRequest) (*pb.GetDataResponse,  error) {
	list:=footModel.GetFoot(request.Uid)
	result:=&pb.GetDataResponse{
		List:make([]*pb.GetDataResult,0),
	}
	for _,v:=range list{
		result.List=append(result.List,&pb.GetDataResult{
			Mac:v.Mac,
			Uid:v.UID,
		})
	}
	return result,nil
}

func (rpc *rpcServer) GetFootBallData(ctx context.Context, request *pb.GetFootBallRequest) (*pb.GetFootBallResponseData,  error) {
	return &pb.GetFootBallResponseData{},nil
}

func (rpc *rpcServer) Bind(ctx context.Context, request *pb.BindFootRequest) (result *pb.EmptyResponse, err error) {
	return &pb.EmptyResponse{},footModel.Bind(request.Uid,request.Mac)
}

func (rpc *rpcServer) Unbind(ctx context.Context, request *pb.BindFootRequest) (result *pb.EmptyResponse, err error) {
	return &pb.EmptyResponse{},footModel.UnBind(request.Uid,request.Mac)
}

func (rpc *rpcServer) GetSchool(ctx context.Context, request *pb.FootBallSchoolRequest) (result *pb.FootBallSchoolResponse, err error) {
	id,err:=footModel.GetSchool(request.Uid,request.Type)
	return &pb.FootBallSchoolResponse{SchoolId:id},err
}


func send(mac, date string) ([]*pb.GetDataResult, error) {
	url:=fmt.Sprintf("http://web.ostrichfc.cn:8000/api/ft/movement/bymac/?macid=%s&type=1&date=%s",mac,date)
	req, err:= http.NewRequest("GET",url ,nil)
	if err != nil {
		return nil,err
	}
	req.Header = http.Header{}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	type Result struct {
		Mac string `json:"ft_mac"`
		Cnt string `json:"ft_cnt"`
		Duration string `json:"ft_duration"`
	}
	type Data struct {
		Code  int64    `json:"code"`
		Data []Result `json:"data"`
		Message string `json:"message"`
	}
	tempData := &Data{}
	xlog.Info(string(result))
	err = json.Unmarshal(result, &tempData)
	if err != nil {
		return nil,xlog.Error("请联系管理员")
	}
	xlog.Info(tempData)

	tempResult:=make([]*pb.GetDataResult,0)
	if tempData.Code==0{
		for _,v:=range tempData.Data{
			tInt,_:=strconv.ParseInt(v.Cnt,10,64)
			tFloat,_:=strconv.ParseFloat(v.Duration,10)
			tempResult=append(tempResult,&pb.GetDataResult{
				Mac:v.Mac,
				Score:tInt,
				Duration:tFloat,
				Date:date,
			})
		}
	}
	return tempResult, nil
}

