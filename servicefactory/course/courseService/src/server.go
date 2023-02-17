package main

import (
	"mufe_service/camp/service"
	_ "mufe_service/service/course"
	_ "mufe_service/service/chapter"
	_ "mufe_service/service/collection"
	_ "mufe_service/service/collection"
	_ "mufe_service/service/video"
	_ "mufe_service/service/recommend"
	_ "mufe_service/service/search"
)

func main() {
	service.StartService()
}
