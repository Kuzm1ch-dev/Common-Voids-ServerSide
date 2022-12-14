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
		log.Println("No .env file found")
	}

	gravity := box2d.MakeB2Vec2(0.0, -10.0)
	world := box2d.MakeB2World(gravity)
	collisionSystem := physic.CollisionSystem{}
	collisionSystem.NewListener(&world)

	game := game.NewGameController(&world, &collisionSystem)

	server := new(common.Server)
	server.Init("127.0.0.1", 1024)
	go server.ListenAndServe()

	game.Run()
}
