package main

import (
	"mufe_service/camp/service"
	_ "mufe_service/service/good/category"
	_ "mufe_service/service/good/delivery"
	_ "mufe_service/service/good/good"
)

func main() {
	service.StartService()
}
