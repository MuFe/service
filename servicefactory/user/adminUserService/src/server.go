package main

import (
	"mufe_service/camp/service"
	_ "mufe_service/service/admin/user"
)

func main() {
	service.StartService()
}
