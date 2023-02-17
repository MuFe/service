package main

import (
	"mufe_service/camp/service"
	_ "mufe_service/service/system/qiniu"
)

func main() {
	service.StartService()
}
