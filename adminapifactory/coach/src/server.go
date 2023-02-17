package main

import (
	"os"
	_ "mufe_service/api/admin/coach"
	"mufe_service/camp/server"
)

func main() {
	server.Start(os.Getenv("PORT"))
}
