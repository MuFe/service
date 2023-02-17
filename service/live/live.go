package live

import (
	"context"
	"mufe_service/camp/enum"
	"mufe_service/camp/service"
	pb "mufe_service/jsonRpc"
	"mufe_service/manager"
	liveModel "mufe_service/model/live"
)

func init() {
	nSer := &rpcServer{}
	pb.RegisterLiveServiceServer(service.GetRegisterRpc(), nSer)

}

type rpcServer struct {
}

func (rpc *rpcServer) GetList(ctx context.Context, request *pb.GetLiveDataRequest) (*pb.GetLiveDataResponse, error) {
	list,err:=liveModel.GetLiveList([]int64{},request.Uid)
	result:=make([]*pb.GetLiveDataResult,0)
	resultMap:=make(map[int64]*pb.GetLiveDataResult)
	if err==nil{
		for _,v:=range list{
			temp,ok:=resultMap[v.Type]
			if !ok{
				temp=&pb.GetLiveDataResult{
					Type:v.Type,
					List:make([]*pb.GetLiveData,0),
				}
				result=append(result,temp)
				resultMap[v.Type]=temp
			}
			temp.List=append(temp.List,&pb.GetLiveData{
				Id:v.Id,
				Address:v.Address,
				Start:v.StartTime,
				End:v.EndTime,
				Status:v.Status,
				Home:&pb.TeamData{
					Id:v.Home.Id,
					Name:v.Home.Name,
					Head:v.Home.Head,
				},
				Visiting:&pb.TeamData{
					Id:v.Visiting.Id,
					Name:v.Visiting.Name,
					Head:v.Visiting.Head,
				},
				Type:v.Type,
				HomeScore:v.HomeScore,
				VisitingScore:v.VisitingScore,
			})
		}
	}
	idList,err:=liveModel.GetUserWatch(request.Uid)
	if err==nil{
		temp:=&pb.GetLiveDataResult{
			Type:enum.USERLIVE,
			List:make([]*pb.GetLiveData,0),
		}
		result=append(result,temp)
		if len(idList)>0{
			tempResult,_:=liveModel.GetLiveList(idList,0)


			for _,v:=range tempResult{
				temp.List=append(temp.List,&pb.GetLiveData{
					Id:v.Id,
					Address:v.Address,
					Start:v.StartTime,
					End:v.EndTime,
					Status:v.Status,
					Home:&pb.TeamData{
						Id:v.Home.Id,
						Name:v.Home.Name,
						Head:v.Home.Head,
					},
					Visiting:&pb.TeamData{
						Id:v.Visiting.Id,
						Name:v.Visiting.Name,
						Head:v.Visiting.Head,
					},
					Type:v.Type,
					HomeScore:v.HomeScore,
					VisitingScore:v.VisitingScore,
				})
			}
		}
	}
	return &pb.GetLiveDataResponse{
		List:result,
	},err
}



func (rpc *rpcServer) Create(ctx context.Context, request *pb.CreateLiveRequest) (*pb.CreateLiveResponse, error) {
	homeList:=make([]*liveModel.LiveTeamMemberData,0)
	visitingList:=make([]*liveModel.LiveTeamMemberData,0)
	for _,v:=range request.HomeMember{
		homeList=append(homeList,&liveModel.LiveTeamMemberData{
			Name:v.Name,
			Number:v.Number,
			Uid:v.Uid,
		})
	}
	for _,v:=range request.VisitingMember{
		visitingList=append(visitingList,&liveModel.LiveTeamMemberData{
			Name:v.Name,
			Number:v.Number,
			Uid:v.Uid,
		})
	}
	id,pass,err:=liveModel.CreateLive(request.Type, request.ClassId,request.MatchId,request.Uid,request.PackageId,request.Home,request.Visiting,request.Address,homeList,visitingList)
	return &pb.CreateLiveResponse{
		Id:id,
		Pass:pass,
	}, err
}


func (rpc *rpcServer) Start(ctx context.Context, request *pb.StartLivRequest) (*pb.StartLivResponse, error) {
	address,end,err:=liveModel.UpdateLivePass(request.Uid,request.Pass)
	return &pb.StartLivResponse{Address:address,End:end},err
}


func (rpc *rpcServer) End(ctx context.Context, request *pb.EndLivRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{},liveModel.End(request.Id,request.Uid)
}


func (rpc *rpcServer) UpdateScore(ctx context.Context, request *pb.UpdateScoreRequest) (*pb.EmptyResponse, error) {
	list:=make([]*liveModel.ScoreInfo,0)
	for _,v:=range request.List{
		list=append(list,&liveModel.ScoreInfo{
			MemberId:v.Id,
			Time:v.Time,
		})
	}
	return &pb.EmptyResponse{}, liveModel.UpdateLiveScore(request.Id,request.HomeScore,request.VisitingScore,list)
}

func (rpc *rpcServer) TeamMember(ctx context.Context, request *pb.TeamMemberRequest) (*pb.TeamMemberResponse, error) {
	homeList,visitingList,err:=getMemberInfo(ctx,request.Id)
	return &pb.TeamMemberResponse{
		HomeMember:homeList,
		VisitingMember:visitingList,
	}, err
}

func (rpc *rpcServer) LiveInfo(ctx context.Context, request *pb.LiveRequest) (*pb.GetLiveData, error) {
	result,err:=liveModel.GetLiveList([]int64{request.Id},0)
	re:=&pb.GetLiveData{}
	if err==nil&&len(result)>0{
		v:=result[0]
		re=&pb.GetLiveData{
			Id:v.Id,
			Address:v.Address,
			Start:v.StartTime,
			End:v.EndTime,
			Status:v.Status,
			Home:&pb.TeamData{
				Id:v.Home.Id,
				Name:v.Home.Name,
				Head:v.Home.Head,
			},
			Visiting:&pb.TeamData{
				Id:v.Visiting.Id,
				Name:v.Visiting.Name,
				Head:v.Visiting.Head,
			},
			Info:make([]*pb.ScoreData,0),
			Type:v.Type,
			HomeScore:v.HomeScore,
			VisitingScore:v.VisitingScore,
		}
		homeList,visitingList,err:=getMemberInfo(ctx,request.Id)
		if err==nil{
			list,err:=liveModel.GetScoreInfo(request.Id)
			if err==nil{
				homeMap:=make(map[int64]*pb.TeamMemberData)
				visitingMap:=make(map[int64]*pb.TeamMemberData)
				for _,v:=range homeList{
					homeMap[v.Id]=v
				}
				for _,v:=range visitingList{
					visitingMap[v.Id]=v
				}
				for _,v:=range list{
					temp,ok:=homeMap[v.MemberId]
					if ok{
						re.Info=append(re.Info,&pb.ScoreData{
							Time:v.Time,
							IsHome:true,
							Data:temp,
						})
					}
					temp,ok=visitingMap[v.MemberId]
					if ok{
						re.Info=append(re.Info,&pb.ScoreData{
							Time:v.Time,
							IsHome:false,
							Data:temp,
						})
					}
				}
			}
		}
	}

	return re, nil
}


func (rpc *rpcServer) EndLive(ctx context.Context, request *pb.StartLivRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, liveModel.EndLive(request.Id,request.Uid)
}

func (rpc *rpcServer) Watch(ctx context.Context, request *pb.WatchLiveRequest) (*pb.StartLivResponse, error) {
	address,id,err:=liveModel.StartWatch(request.Uid,request.Pass)
	return &pb.StartLivResponse{
		Address:address,
		Id:id,
	}, err
}

func (rpc *rpcServer) GetWatchInfo(ctx context.Context, request *pb.LiveRequest) (*pb.GetLiveData, error) {
	result,err:=liveModel.GetWatchNumber(request.Id)
	re:=&pb.GetLiveData{}
	if err==nil{
		re.User=result
	}
	return re, nil
}

func (rpc *rpcServer) EndWatch(ctx context.Context, request *pb.WatchLiveRequest) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, liveModel.OutWatch(request.Uid)
}

func (rpc *rpcServer) Package(ctx context.Context, request *pb.EmptyRequest) (*pb.PackageDataResponse, error) {
	list,err:=liveModel.GetPackage()
	resultList:=make([]*pb.PackageData,0)
	if err==nil{
		for _,v:=range list{
			resultList=append(resultList,&pb.PackageData{
				Id:v.Id,
				Name:v.Name,
				Max:v.Max,
				Price:v.Price,
				Duration:v.Duration,
			})
		}
	}
	return &pb.PackageDataResponse{List:resultList}, err
}



func getMemberInfo(ctx context.Context, id int64)([]*pb.TeamMemberData,[]*pb.TeamMemberData,error){
	home,visiting,err:=liveModel.GetLiveTeamMember(id)
	homeList:=make([]*pb.TeamMemberData,0)
	visitingList:=make([]*pb.TeamMemberData,0)
	if err==nil{
		uidList:=make([]int64,0)
		dataMap:=make(map[int64]*pb.TeamMemberData)
		for _,v:=range home{

			temp:=&pb.TeamMemberData{
				Id:v.Id,
				Number:v.Number,
				Name:v.Name,
			}
			if v.Uid!=0{
				uidList=append(uidList,v.Uid)
				dataMap[v.Uid]=temp
			}

			homeList=append(homeList,temp)
		}
		for _,v:=range visiting{
			temp:=&pb.TeamMemberData{
				Id:v.Id,
				Number:v.Number,
				Name:v.Name,
			}
			if v.Uid!=0{
				uidList=append(uidList,v.Uid)
				dataMap[v.Uid]=temp
			}
			visitingList=append(visitingList,temp)
		}
		if len(uidList)>0{
			userResult,err:=manager.GetUserService().GetUserList(ctx,&pb.GetUserListRequest{IdList:uidList})
			if err==nil{
				for _,v:=range userResult.List{
					temp,ok:=dataMap[v.Uid]
					if ok{
						temp.Head=v.Head
					}
				}
			}
		}
	}
	return homeList,visitingList,err
}
