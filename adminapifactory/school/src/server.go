package main

import (
	"os"
	_ "mufe_service/api/admin/school"
	"mufe_service/camp/server"
)

func main() {
	server.Start(os.Getenv("PORT"))
}
