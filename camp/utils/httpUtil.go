package utils

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"mufe_service/camp/consul"
)

func IsMiniProgram(c *gin.Context) bool {
	if c.GetHeader("ProgramTag") == "miniProgram" {
		return true
	} else {
		return false
	}
}

func GetHeaderFromKey(c *gin.Context, key string) string {
	return c.GetHeader(key)
}


func GetInt32ValueFromReq(c *gin.Context, key string) int32 {
	return FormatStrToInt32(c.Query(key))
}

func GetInt64ValueFromReq(c *gin.Context, key string) int64 {
	return FormatStrToInt64(c.Query(key))
}

func GetFloat64ValueFromReq(c *gin.Context, key string) float64 {
	value, err := strconv.ParseFloat(c.Query(key), 64)
	if err != nil {
		value = 0
	}
	return value
}

func FormatStrToInt32(str string) int32 {
	intStr, err := strconv.Atoi(str)
	if err != nil {
		intStr = 0
	}
	return int32(intStr)
}

func FormatStrToInt64(str string) int64 {
	intStr, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		intStr = 0
	}
	return intStr
}

type BaseResult struct {
	Total int64         `json:"total"`
	Current int64         `json:"current"`
	List  []interface{} `json:"list"`
}

func CreateListResultReturn(total int64, list []interface{}) BaseResult {
	return BaseResult{
		Total: total,
		List:  list,
	}
}
func CreateListCurrentResultReturn(total ,current int64, list []interface{}) BaseResult {
	return BaseResult{
		Total: total,
		Current: current,
		List:  list,
	}
}

//获取rpc服务(服务发现)
func GetRPCService(name string, tag string, consulIp string) (*grpc.ClientConn, error) {
	r := consul.NewResolver(name, tag)
	b := grpc.RoundRobin(r)
	conn, err := grpc.Dial(consulIp, grpc.WithInsecure(), grpc.WithBalancer(b), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	return conn, err
}

func LocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
