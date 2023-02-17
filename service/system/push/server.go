package push

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/ylywyn/jpush-api-go-client"
	"io/ioutil"
	"net/http"
	"os"
	"mufe_service/camp/errcode"
	"mufe_service/camp/service"
	"mufe_service/camp/xlog"
	pb "mufe_service/jsonRpc"
	"strings"
)

type rpcServer struct {
}



func init() {
	pb.RegisterPushServiceServer(service.GetRegisterRpc(), &rpcServer{})
}

func (rpc *rpcServer) PushMessage(ctx context.Context, request *pb.PushRequest) (result *pb.PushResponse, err error) {
	list:=make([]string,0)
	for _,v:=range request.DeviceList{
		if v!=""{
			list=append(list,v)
		}
	}
	if len(list)==0{
		return nil,xlog.Error(errcode.HttpErrorWringParam.Msg)
	}
	return send(request.DeviceList,request.Content,request.Message)
}

func (rpc *rpcServer) GetPhone(ctx context.Context, request *pb.GetPhoneRequest) (result *pb.GetPhoneResponse, err error) {
	return getPhone(request.Token)
}

func getPhone(token string)(*pb.GetPhoneResponse, error){
	key:=os.Getenv("PUSH_APP_KEY")
	secret:=os.Getenv("PUSH_APP_SECRET")
	c := jpushclient.NewPushClient(secret, key)
	type Data struct {
		Token   string  `json:"loginToken" `
	}
	c.BaseUrl="https://api.verification.jpush.cn/v1/web/loginTokenVerify"
	data:=Data{
		Token:token,
	}
	content, err := json.Marshal(data)
	str, err := SendPostBytes2("https://api.verification.jpush.cn/v1/web/loginTokenVerify",content,c.AuthCode)
	if err != nil {
		xlog.Info(fmt.Sprintf("err:%s", err.Error()))
		return nil,xlog.Error(err)
	} else {
		xlog.Info(fmt.Sprintf("ok:%s", str))
	}
	type Result struct {
		Phone   string  `json:"phone" `
	}
	result:=Result{}
	err = json.Unmarshal([]byte(str), &result)
	jKey,err:= ioutil.ReadFile("./jpush_private.txt")
	jKeyStr:=FormatAlipayPrivateKey(string(jKey))
	if err != nil {
		xlog.Info("invalid encrypted")
		return nil,err
	}
	encryptedB, err := base64.StdEncoding.DecodeString(result.Phone)
	if err != nil {
		xlog.Info("invalid encrypted")
		return nil,err
	}
	temp, err := RsaDecrypt(encryptedB, []byte(jKeyStr))
	if err != nil {
		xlog.Info("invalid encrypted")
		return nil,err
	}
	return &pb.GetPhoneResponse{
		Phone:string(temp),
	}, nil
}

func send(device  []string,content string,message string) (*pb.PushResponse, error) {
	key:=os.Getenv("PUSH_APP_KEY")
	secret:=os.Getenv("PUSH_APP_SECRET")
	var pf jpushclient.Platform
	pf.All()

	//Audience
	var ad jpushclient.Audience
	ad.SetID(device)
	payload := jpushclient.NewPushPayLoad()
	payload.SetPlatform(&pf)
	payload.SetAudience(&ad)

	if message!=""{
		msg:=&jpushclient.Message{
			Content:message,
		}
		payload.SetMessage(msg)
	}else{
		//Notice
		var notice jpushclient.Notice
		notice.SetAlert(content)
		notice.SetAndroidNotice(&jpushclient.AndroidNotice{Alert: content})
		notice.SetIOSNotice(&jpushclient.IOSNotice{Alert:content})
		payload.SetNotice(&notice)
	}


	bytes, _ := payload.ToBytes()


	//push
	c := jpushclient.NewPushClient(secret, key)
	str, err := c.Send(bytes)
	if err != nil {
		xlog.Info(fmt.Sprintf("err:%s", err.Error()))
		return nil,xlog.Error(err)
	}
	return &pb.PushResponse{
		Message:str,
	}, err
}

func SendPostBytes2(url string, data []byte, authCode string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Add("Authorization", authCode)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return "", err
	}
	if resp == nil {
		return "", nil
	}

	defer resp.Body.Close()
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(r), nil
}

func RsaDecrypt(encrypted, prikey []byte)([]byte, error) {var data []byte
	block, _ := pem.Decode(prikey)
	if block == nil {return data, errors.New("private key error")
	}
	rsaKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {return data, err}
	key, ok := rsaKey.(*rsa.PrivateKey)
	if !ok {return data, errors.New("invalid private key")
	}
	data, err = rsa.DecryptPKCS1v15(rand.Reader, key, encrypted)
	return data, err
}

func FormatAlipayPrivateKey(privateKey string) (pKey string) {
	var buffer strings.Builder
	buffer.WriteString("-----BEGIN RSA PRIVATE KEY-----\n")
	rawLen := 64
	keyLen := len(privateKey)
	raws := keyLen / rawLen
	temp := keyLen % rawLen
	if temp > 0 {
		raws++
	}
	start := 0
	end := start + rawLen
	for i := 0; i < raws; i++ {
		if i == raws-1 {
			buffer.WriteString(privateKey[start:])
		} else {
			buffer.WriteString(privateKey[start:end])
		}
		buffer.WriteByte('\n')
		start += rawLen
		end = start + rawLen
	}
	buffer.WriteString("-----END RSA PRIVATE KEY-----\n")
	pKey = buffer.String()
	return
}
