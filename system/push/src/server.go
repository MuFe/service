package main

import (
	"mufe_service/camp/service"
	_ "mufe_service/service/system/push"
)

func main() {
	service.StartService()
}
