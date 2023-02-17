package qiniu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"mufe_service/camp/service"
	"mufe_service/camp/xlog"
	pb "mufe_service/jsonRpc"
	"strconv"
	"time"
)

type rpcServer struct {
}

var accessKey string
var secretKey string

func (rpc *rpcServer) UploadBytes(ctx context.Context, request *pb.UploadBytesReq) (result *pb.UploadBytesRsp, err error) {
	result = &pb.UploadBytesRsp{}
	putPolicy := storage.PutPolicy{
		Scope: request.Bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	regionResult, err := storage.GetZone(accessKey, request.Bucket)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = regionResult
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}

	dataLen := int64(len(request.Data))
	err = formUploader.Put(
		ctx,
		&ret,
		upToken,
		request.Key,
		bytes.NewReader(request.Data),
		dataLen,
		&putExtra)
	if err != nil {
		return result, xlog.Error(err)
	}
	return result, nil
}
func (rpc *rpcServer) GetToken(ctx context.Context, request *pb.QiniuServiceRequest) (result *pb.QiniuServiceResponse, err error) {
	putPolicy := storage.PutPolicy{
		Scope: request.Bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	regionResult, err := storage.GetZone(accessKey, request.Bucket)
	if err != nil {
		return nil, err
	}
	upToken := putPolicy.UploadToken(mac)


	return &pb.QiniuServiceResponse{Token: upToken, UploadHost:regionResult.SrcUpHosts[0], Base64UploadHost: regionResult.CdnUpHosts[0]}, err
}

func (rpc *rpcServer) GetVideoInfo(ctx context.Context, request *pb.GetVideoInfoRequest) (*pb.VideoInfoResponse, error) {
	resp, err := http.Get(request.Prefix + request.Url + "?avinfo")
	if err != nil {
		xlog.ErrorP(err)
		return nil, xlog.Error(err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return nil, xlog.Error(err)
		}
	}
	type format struct {
		Duration string `json:"duration"`
	}
	type streams struct {
		Duration string `json:"duration"`
	}
	type data struct {
		Format  *format    `json:"format"`
		Streams []*streams `json:"streams"`
	}
	tempData := &data{}
	err = json.Unmarshal([]byte(result.String()), &tempData)
	if err != nil {
		return nil, xlog.Error(err)
	}
	var duration int64
	if tempData.Format != nil {
		dInt, err := strconv.ParseFloat(tempData.Format.Duration,  64)
		if err==nil{
			duration=int64(dInt)
		} else {
			duration =0
		}
	} else if len(tempData.Streams)>0{
		dInt, err := strconv.ParseFloat(tempData.Streams[0].Duration, 64)
		if err==nil{
			duration=int64(dInt)
		} else {
			duration =0
		}
	}
	return &pb.VideoInfoResponse{Duration:duration,CoverPrefix:request.Prefix,Cover:request.Url+"?vframe/jpg/offset/1"}, err
}

func (rpc *rpcServer) WordInspect(ctx context.Context, request *pb.WordInspectRequest) (result *pb.WordInspectResponse, err error) {
	can,err:=canUse(request.Content)
	typeInt:=int64(0)
	if can{
		typeInt=1
	}
	return &pb.WordInspectResponse{Type:typeInt}, err
}

func canUse(content string)(bool,error){
	data:=`{
    "data": {
        "text": "%s"
    },
    "params": {
        "scenes": [
            "antispam"
        ]
    }
}`

	req, err:= http.NewRequest("POST", "http://ai.qiniuapi.com/v3/text/censor",bytes.NewBufferString(fmt.Sprintf(data,content)))
	if err != nil {
		return false,xlog.Error(err)
	}
	req.Header = http.Header{}
	req.Header.Add("Content-Type", "application/json")
	mac := qbox.NewMac(accessKey, secretKey)
	token, signErr := mac.SignRequestV2(req)
	if signErr != nil {
		err = signErr
		 	return false,xlog.Error(err)
	}
	req.Header.Add("Authorization", "Qiniu "+token)
	xlog.Info(token)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	type Result struct {
		Suggestion string `json:"suggestion"`
	}
	type Data struct {
		Code  int64    `json:"code"`
		Result Result `json:"result"`
	}
	tempData := &Data{}
	err = json.Unmarshal(result, &tempData)
	if err != nil {
		return false, xlog.Error(err)
	}
	if tempData.Code==200{
		if tempData.Result.Suggestion=="pass"{
			return true,nil
		}
	}
	return false,nil
}

func init() {
	accessKey = os.Getenv("AK")
	secretKey = os.Getenv("SK")
	nSer := &rpcServer{}
	pb.RegisterQiniuServiceServer(service.GetRegisterRpc(), nSer)
}
