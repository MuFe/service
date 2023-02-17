package main

import (
	"mufe_service/camp/service"
	_ "mufe_service/service/coach"
)

func main() {
	service.StartService()
}
