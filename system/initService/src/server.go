package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"mufe_service/camp/xlog"
	"strconv"
	"time"
)

type Consul struct {
	ServiceAddress string
	ServicePort    int
}

func main() {

	list:=make([]string,0)
	name := os.Getenv("HTTPS_NAME")
	list=append(list,name)
	name1:=os.Getenv("HTTP_NAME")
	if name1!=""{
		list=append(list,name1)
	}
	for _,value:= range list{
		str, _ := get("http://"+os.Getenv("CONSUL_IP")+"/v1/catalog/service/"+value, "GET")
		var data []Consul
		json.Unmarshal([]byte(str), &data)
		for _, v := range data {
			get("http://"+os.Getenv("CONSUL_IP")+"/v1/kv/upstreams/"+value+"/"+v.ServiceAddress+":"+strconv.Itoa(v.ServicePort), "PUT")
		}
	}


}

func get(url string, method string) (response string, err error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, _ := http.NewRequest(method, url, nil)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return "", xlog.Error("请求token失败")
		}
	}

	response = result.String()
	return response, nil
}
