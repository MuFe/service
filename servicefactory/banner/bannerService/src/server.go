package main

import (
	"mufe_service/camp/service"
	_ "mufe_service/service/banner"
	_ "mufe_service/service/web"
)

func main() {
	service.StartService()
}
