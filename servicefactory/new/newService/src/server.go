package main

import (
	"mufe_service/camp/service"
	_ "mufe_service/service/new"
)

func main() {
	service.StartService()
}
