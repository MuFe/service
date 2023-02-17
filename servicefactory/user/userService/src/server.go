package main

import (
	"mufe_service/camp/service"
	_ "mufe_service/service/user"
	_ "mufe_service/service/system/feedback"
)

func main() {
	service.StartService()
}
