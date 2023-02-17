package main

import (
	"mufe_service/camp/service"
	_ "mufe_service/service/live"
)

func main() {
	service.StartService()
}
