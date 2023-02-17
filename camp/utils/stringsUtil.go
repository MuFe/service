package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

//mysql in写法
func MysqlStringInUtils(buf *strings.Builder, list []int64, preString string) {
	list = RemoveDupsInt64(list)
	if len(list) == 1 && list[0] == 0 {
		return
	}
	for k, info := range list {
		if len(list) == 1 {
			buf.WriteString(preString)
			buf.WriteString(" =")
			buf.WriteString(strconv.FormatInt(info, 10))
		} else {
			if k == 0 {
				buf.WriteString(preString)
				buf.WriteString(" in (")
			}
			buf.WriteString(strconv.FormatInt(info, 10))
			if k < len(list)-1 {
				buf.WriteString(",")
			} else {
				buf.WriteString(") ")
			}
		}

	}
}

//mysql in写法
func MysqlStringInUtilsWithZero(buf *strings.Builder, list []int64, preString string) {
	list = RemoveDupsInt64(list)
	for k, info := range list {
		if len(list) == 1 {
			buf.WriteString(preString)
			buf.WriteString(" =")
			buf.WriteString(strconv.FormatInt(info, 10))
		} else {
			if k == 0 {
				buf.WriteString(preString)
				buf.WriteString(" in (")
			}
			buf.WriteString(strconv.FormatInt(info, 10))
			if k < len(list)-1 {
				buf.WriteString(",")
			} else {
				buf.WriteString(") ")
			}
		}

	}
}

//mysql in写法
func MysqlInUtils(buf *strings.Builder, list []string, preString string) {
	RemoveDupsString(list)
	if len(list) == 1 && list[0] == "" {
		return
	}
	for k, info := range list {
		if k == 0 {
			buf.WriteString(preString)
		}
		buf.WriteString(fmt.Sprintf("'%s'", info))
		if k < len(list)-1 {
			buf.WriteString(",")
		} else {
			buf.WriteString(") ")
		}
	}
}

func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}


var baseStr= "0123456789abcdedfhigklmnopqrstuvwxyz"

func BaseCode6(inviteMap map[string]bool)string{
	bytes := []byte(baseStr)
	rand.Seed(time.Now().Unix())
	for {
		id := make([]byte, 6)
		for i := 0; i < 6; i++ {
			id[i] = bytes[rand.Int()%len(bytes)]
		}
		idstr := string(id)
		_,ok:=inviteMap[idstr]
		if !ok{
			return idstr
		}
	}
}


func BaseCode4(inviteMap map[string]bool)string{
	bytes := []byte(baseStr)
	rand.Seed(time.Now().Unix())
	for {
		id := make([]byte, 4)
		for i := 0; i < 4; i++ {
			id[i] = bytes[rand.Int()%len(bytes)]
		}
		idstr := string(id)
		_,ok:=inviteMap[idstr]
		if !ok{
			return idstr
		}
	}
}
