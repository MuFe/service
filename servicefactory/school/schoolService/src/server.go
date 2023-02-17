package main

import (
	"mufe_service/camp/service"
	_ "mufe_service/service/school"
)

func main() {
	service.StartService()
}
