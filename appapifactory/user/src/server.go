package main

import (
	"os"
	_ "mufe_service/api/app/system"
	_ "mufe_service/api/app/user"
	"mufe_service/camp/server"
)

func main() {
	server.Start(os.Getenv("PORT"))
}
