package common

import (
	"github.com/google/uuid"
	"net"
	"server/src/game"
)

type Client struct {
	Uuid   uuid.UUID
	Conn   net.Conn
	Player game.Player
}
