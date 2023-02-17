package main

import (
	"os"
	_ "mufe_service/api/app/school"
	_ "mufe_service/api/app/homework"
	"mufe_service/camp/server"
)

func main() {
	server.Start(os.Getenv("PORT"))
}
