package main

import (
	"mufe_service/camp/service"
	_ "mufe_service/service/homework"
)

func main() {
	service.StartService()
}
