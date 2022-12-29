package main

import (
	"github.com/joho/godotenv"
	"log"
	"server/src/common"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	authServer := new(common.AuthServer)
	authServer.Init("127.0.0.1")
	go authServer.ListenAndServe()

	logicServer := new(common.LogicServer)
	logicServer.Init("127.0.0.1", 1024)
	go logicServer.ListenAndServe()

	logicServer.GameController.Run()
}
