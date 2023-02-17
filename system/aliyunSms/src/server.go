package main

import (
	"mufe_service/camp/service"
	_ "mufe_service/service/system/aliyunSms"
	_ "mufe_service/service/system/football"
)

func main() {
	service.StartService()
}
