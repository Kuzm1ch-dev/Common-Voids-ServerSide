package main

import (
	"github.com/joho/godotenv"
	"log"
	"server/src/common"
)

const (
	UUIDPackage      int32 = 101
	pNewPlayer             = 102
	pMessage               = 103
	pBroadcast             = 104
	pUpdateEquuiment       = 105
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		log.Println("No .env file found")
	}

	server := new(common.Server)
	server.Init("127.0.0.1", 1024)
	server.ListenAndServe()
}
