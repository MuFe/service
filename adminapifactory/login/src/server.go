package main

import (
	"os"
	_ "mufe_service/api/admin/login"
	_ "mufe_service/api/admin/brand"
	"mufe_service/camp/server"
)

func main() {
	server.Start(os.Getenv("PORT"))
}
