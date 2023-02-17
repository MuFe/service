package main

import (
	"os"
	_ "mufe_service/api/app/home"
	_ "mufe_service/api/app/testapi"
	"mufe_service/camp/server"
)

func main() {
	server.Start(os.Getenv("PORT"))
}
