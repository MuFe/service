package main

import (
	"os"
	_ "mufe_service/api/admin/course"
	_ "mufe_service/api/admin/home"
	_ "mufe_service/api/admin/feedback"
	"mufe_service/camp/server"
)

func main() {
	server.Start(os.Getenv("PORT"))
}
