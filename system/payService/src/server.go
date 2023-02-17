package main

import (
	"mufe_service/camp/service"
	_ "mufe_service/service/system/pay"
)

func main() {
	service.StartService()
}
