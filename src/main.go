package main

import (
	"github.com/ByteArena/box2d"
	"github.com/joho/godotenv"
	"log"
	"server/src/common"
	"server/src/game"
	"server/src/game/physic"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	gravity := box2d.MakeB2Vec2(0.0, -10.0)
	world := box2d.MakeB2World(gravity)
	collisionSystem := physic.CollisionSystem{}
	collisionSystem.NewListener(&world)

	GameController := game.NewGameController(&world, &collisionSystem)

	authServer := new(common.AuthServer)
	authServer.Init("127.0.0.1")
	go authServer.ListenAndServe()

	logicServer := new(common.LogicServer)
	logicServer.Init("127.0.0.1", 1024)
	go logicServer.ListenAndServe()

	GameController.Run()
}
